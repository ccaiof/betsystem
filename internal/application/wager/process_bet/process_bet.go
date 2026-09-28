package process_bet

import (
	"context"
	"errors"
	"uuid"

	"github.com/betsystem/internal/application/wager"
	_ "github.com/betsystem/internal/application/wager"
	"github.com/betsystem/internal/domain/ledger"
	"github.com/betsystem/internal/domain/money"
	domainwager "github.com/betsystem/internal/domain/wager"
	domainwallet "github.com/betsystem/internal/domain/wallet"
)

type ProcessBet struct {
	wallets    wager.WalletRepository
	wagers     wager.WagerTransactionRepository
	ledger     wager.LedgerRepository
	transactor wager.Transactor
}

func NewProcessBet(
	wallets wager.WalletRepository,
	wagers wager.WagerTransactionRepository,
	ledger wager.LedgerRepository,
	transactor wager.Transactor,
) *ProcessBet {
	return &ProcessBet{
		wallets:    wallets,
		wagers:     wagers,
		ledger:     ledger,
		transactor: transactor,
	}
}

func (p *ProcessBet) Execute(
	ctx context.Context,
	cmd ProcessBetCommand,
) (ProcessBetResult, error) {
	amount, err := money.FromDecimalString(
		cmd.Amount,
		cmd.Currency,
		false,
	)
	if err != nil {
		return ProcessBetResult{}, err
	}

	transaction, err := domainwager.NewExternal(
		cmd.TransactionID,
		cmd.ProviderID,
		cmd.ExternalTransactionID,
		cmd.IdempotencyKey,
		cmd.PayloadHash,
		cmd.WalletID,
		cmd.PlayerID,
		cmd.RoundID,
		cmd.GameID,
		domainwager.KindBet,
		amount,
		cmd.Now,
	)
	if err != nil {
		return ProcessBetResult{}, err
	}

	var result ProcessBetResult

	err = p.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		wallet, err := p.wallets.GetForUpdate(txCtx, cmd.WalletID)
		if err != nil {
			return err
		}

		if wallet.PlayerID != cmd.PlayerID {
			return domainwager.ErrWalletPlayerMismatch
		}

		if wallet.Balance.Currency != cmd.Currency {
			return money.ErrCurrencyMismatch
		}

		balanceBefore := wallet.Balance

		if debitErr := wallet.Debit(amount); debitErr != nil {
			if !errors.Is(debitErr, domainwallet.ErrInsufficientFunds) {
				return debitErr
			}

			// Saldo insuficiente é um resultado de negócio, não uma falha técnica:
			// a transação rejeitada precisa ser persistida (idempotência e auditoria),
			// então não podemos devolver erro aqui, senão o rollback a descartaria.
			if err := transaction.Reject(
				domainwager.FailureCodeInsufficientFunds,
				cmd.Now,
			); err != nil {
				return err
			}

			if err := p.wagers.Save(txCtx, transaction); err != nil {
				return err
			}

			result = ProcessBetResult{
				TransactionID: transaction.ID,
				Status:        string(transaction.Status),
				Balance:       balanceBefore,
			}

			return nil
		}

		entry, err := ledger.NewEntry(
			uuid.New(),
			wallet.ID,
			transaction.ID,
			ledger.DirectionDebit,
			amount,
			balanceBefore,
			wallet.Balance,
			cmd.Now,
		)
		if err != nil {
			return err
		}

		if err := transaction.MarkProcessed(wallet.Balance, cmd.Now); err != nil {
			return err
		}

		if err := p.wallets.Save(txCtx, wallet); err != nil {
			return err
		}

		if err := p.wagers.Save(txCtx, transaction); err != nil {
			return err
		}

		if err := p.ledger.Save(txCtx, entry); err != nil {
			return err
		}

		result = ProcessBetResult{
			TransactionID: transaction.ID,
			Status:        string(transaction.Status),
			Balance:       wallet.Balance,
		}

		return nil
	})
	if err != nil {
		return ProcessBetResult{}, err
	}

	return result, nil
}
