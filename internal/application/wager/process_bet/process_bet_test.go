package process_bet

import (
	"context"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/require"

	domainledger "github.com/betsystem/internal/domain/ledger"
	"github.com/betsystem/internal/domain/money"
	domainwager "github.com/betsystem/internal/domain/wager"
	domainwallet "github.com/betsystem/internal/domain/wallet"
)

var testNow = time.Date(2026, 9, 28, 10, 0, 0, 0, time.UTC)

type fakeWalletRepository struct {
	wallet domainwallet.Wallet

	getForUpdateCalls int
	saveCalls         int
}

func (f *fakeWalletRepository) GetForUpdate(
	ctx context.Context,
	walletID uuid.UUID,
) (domainwallet.Wallet, error) {
	f.getForUpdateCalls++

	return f.wallet, nil
}

func (f *fakeWalletRepository) Save(
	ctx context.Context,
	wallet domainwallet.Wallet,
) error {
	f.saveCalls++
	f.wallet = wallet

	return nil
}

type fakeWagerTransactionRepository struct {
	transactions []domainwager.WagerTransaction
}

func (f *fakeWagerTransactionRepository) Save(
	ctx context.Context,
	transaction domainwager.WagerTransaction,
) error {
	f.transactions = append(f.transactions, transaction)

	return nil
}

type fakeLedgerRepository struct {
	entries []domainledger.LedgerEntry
}

func (f *fakeLedgerRepository) Save(
	ctx context.Context,
	entry domainledger.LedgerEntry,
) error {
	f.entries = append(f.entries, entry)

	return nil
}

type fakeTransactor struct {
	calls int
}

func (f *fakeTransactor) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	f.calls++

	return fn(ctx)
}

type processBetTestDeps struct {
	wallets    *fakeWalletRepository
	wagers     *fakeWagerTransactionRepository
	ledger     *fakeLedgerRepository
	transactor *fakeTransactor
	useCase    *ProcessBet
}

// newProcessBetTestDeps cria as dependências com uma carteira BRL com saldo de 100.00.
func newProcessBetTestDeps(t *testing.T) processBetTestDeps {
	t.Helper()

	wallet, err := domainwallet.New(uuid.New(), uuid.New(), "BRL")
	require.NoError(t, err)

	balance, err := money.FromDecimalString("100.00", "BRL", false)
	require.NoError(t, err)

	require.NoError(t, wallet.Credit(balance))

	deps := processBetTestDeps{
		wallets:    &fakeWalletRepository{wallet: wallet},
		wagers:     &fakeWagerTransactionRepository{},
		ledger:     &fakeLedgerRepository{},
		transactor: &fakeTransactor{},
	}

	deps.useCase = NewProcessBet(
		deps.wallets,
		deps.wagers,
		deps.ledger,
		deps.transactor,
	)

	return deps
}

func newProcessBetCommand(
	wallet domainwallet.Wallet,
	amount string,
) ProcessBetCommand {
	return ProcessBetCommand{
		TransactionID:         uuid.New(),
		ProviderID:            "provider-acme",
		ExternalTransactionID: "external-123",
		IdempotencyKey:        "idempotency-123",
		PayloadHash:           "hash-123",

		WalletID: wallet.ID,
		PlayerID: wallet.PlayerID,

		RoundID: "round-123",
		GameID:  "game-123",

		Amount:   amount,
		Currency: "BRL",

		Now: testNow,
	}
}

func TestProcessBet_ProcessesBetSuccessfully(t *testing.T) {
	deps := newProcessBetTestDeps(t)
	cmd := newProcessBetCommand(deps.wallets.wallet, "30.00")

	result, err := deps.useCase.Execute(context.Background(), cmd)

	require.NoError(t, err)

	require.Equal(t, cmd.TransactionID, result.TransactionID)
	require.Equal(t, string(domainwager.StatusProcessed), result.Status)
	require.Equal(t, "70.00", result.Balance.String())

	require.Equal(t, 1, deps.transactor.calls)
	require.Equal(t, 1, deps.wallets.getForUpdateCalls)
	require.Equal(t, 1, deps.wallets.saveCalls)

	require.Len(t, deps.wagers.transactions, 1)
	require.Len(t, deps.ledger.entries, 1)

	transaction := deps.wagers.transactions[0]

	require.Equal(t, cmd.TransactionID, transaction.ID)
	require.Equal(t, domainwager.KindBet, transaction.Kind)
	require.Equal(t, domainwager.StatusProcessed, transaction.Status)
	require.Equal(t, int64(3000), transaction.Amount.AmountMinor)
	require.Equal(t, "BRL", transaction.Amount.Currency)
	require.NotNil(t, transaction.FinancialResultBalance)
	require.Equal(t, int64(7000), transaction.FinancialResultBalance.AmountMinor)

	entry := deps.ledger.entries[0]

	require.Equal(t, deps.wallets.wallet.ID, entry.WalletID)
	require.Equal(t, cmd.TransactionID, entry.TransactionID)
	require.Equal(t, domainledger.DirectionDebit, entry.Direction)
	require.Equal(t, int64(3000), entry.Amount.AmountMinor)
	require.Equal(t, int64(10000), entry.BalanceBefore.AmountMinor)
	require.Equal(t, int64(7000), entry.BalanceAfter.AmountMinor)
}

func TestProcessBet_RejectsWhenInsufficientFunds(t *testing.T) {
	deps := newProcessBetTestDeps(t)
	cmd := newProcessBetCommand(deps.wallets.wallet, "150.00")

	result, err := deps.useCase.Execute(context.Background(), cmd)

	require.NoError(t, err)

	require.Equal(t, cmd.TransactionID, result.TransactionID)
	require.Equal(t, string(domainwager.StatusRejected), result.Status)
	require.Equal(t, "100.00", result.Balance.String())

	require.Equal(t, 1, deps.transactor.calls)
	require.Equal(t, 1, deps.wallets.getForUpdateCalls)
	require.Equal(t, 0, deps.wallets.saveCalls)

	require.Len(t, deps.wagers.transactions, 1)

	transaction := deps.wagers.transactions[0]

	require.Equal(t, domainwager.StatusRejected, transaction.Status)
	require.NotNil(t, transaction.FailureCode)
	require.Equal(
		t,
		domainwager.FailureCodeInsufficientFunds,
		*transaction.FailureCode,
	)

	require.Empty(t, deps.ledger.entries)
}

func TestProcessBet_DoesNotChangeWalletOnInsufficientFunds(t *testing.T) {
	deps := newProcessBetTestDeps(t)
	originalBalance := deps.wallets.wallet.Balance

	_, err := deps.useCase.Execute(
		context.Background(),
		newProcessBetCommand(deps.wallets.wallet, "101.00"),
	)
	require.NoError(t, err)

	comparison, err := deps.wallets.wallet.Balance.Compare(originalBalance)
	require.NoError(t, err)

	require.Equal(t, 0, comparison)
}

func TestProcessBet_RejectsWalletFromAnotherPlayer(t *testing.T) {
	deps := newProcessBetTestDeps(t)

	cmd := newProcessBetCommand(deps.wallets.wallet, "10.00")
	cmd.PlayerID = uuid.New()

	_, err := deps.useCase.Execute(context.Background(), cmd)

	require.ErrorIs(t, err, domainwager.ErrWalletPlayerMismatch)

	require.Equal(t, 1, deps.transactor.calls)
	require.Equal(t, 1, deps.wallets.getForUpdateCalls)
	require.Equal(t, 0, deps.wallets.saveCalls)

	require.Empty(t, deps.wagers.transactions)
	require.Empty(t, deps.ledger.entries)
}
