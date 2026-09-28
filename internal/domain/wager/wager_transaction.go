package wager

import (
	"time"
	"uuid"

	"github.com/betsystem/internal/domain/money"
)

type WagerTransaction struct {
	ID                     uuid.UUID
	ProviderID             *string
	ExternalTransactionID  *string
	IdempotencyKey         *string
	PayloadHash            *string
	WalletID               uuid.UUID
	PlayerID               uuid.UUID
	RoundID                *string
	GameID                 *string
	Kind                   Kind
	Amount                 money.Money
	ReferenceExternalID    *string
	ReferenceTransactionID *uuid.UUID
	Status                 Status
	FailureCode            *string
	RetryCount             int
	NextRetryAt            *time.Time
	FinancialResultBalance *money.Money
	CreatedAt              time.Time
	UpdatedAt              time.Time
	ProcessedAt            *time.Time
}

func NewExternal(
	id uuid.UUID,
	providerID string,
	externalTransactionID string,
	idempotencyKey string,
	payloadHash string,
	walletID uuid.UUID,
	playerID uuid.UUID,
	roundID string,
	gameID string,
	kind Kind,
	amount money.Money,
	now time.Time,
) (WagerTransaction, error) {
	if err := validateExternalKind(kind); err != nil {
		return WagerTransaction{}, err
	}

	return WagerTransaction{
		ID:                    id,
		ProviderID:            new(providerID),
		ExternalTransactionID: new(externalTransactionID),
		IdempotencyKey:        new(idempotencyKey),
		PayloadHash:           new(payloadHash),
		WalletID:              walletID,
		PlayerID:              playerID,
		RoundID:               new(roundID),
		GameID:                new(gameID),
		Kind:                  kind,
		Amount:                amount,
		Status:                StatusPending,
		CreatedAt:             now,
		UpdatedAt:             now,
	}, nil
}

func NewOpening(
	id uuid.UUID,
	walletID uuid.UUID,
	playerID uuid.UUID,
	amount money.Money,
	now time.Time,
) (WagerTransaction, error) {
	return WagerTransaction{
		ID:        id,
		WalletID:  walletID,
		PlayerID:  playerID,
		Kind:      KindOpening,
		Amount:    amount,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (w *WagerTransaction) MarkProcessed(
	now time.Time,
) error {
	if err := w.transitionTo(StatusProcessed); err != nil {
		return err
	}

	w.ProcessedAt = &now
	w.UpdatedAt = now

	return nil
}

func (w *WagerTransaction) MarkPendingReference(
	now time.Time,
) error {
	if err := w.transitionTo(StatusPendingReference); err != nil {
		return err
	}

	w.UpdatedAt = now

	return nil
}

func (w *WagerTransaction) Reject(
	failureCode string,
	now time.Time,
) error {
	if failureCode == "" {
		return ErrFailureCodeRequired
	}

	if err := w.transitionTo(StatusRejected); err != nil {
		return err
	}

	w.FailureCode = new(failureCode)
	w.UpdatedAt = now

	return nil
}

func (w *WagerTransaction) Fail(
	failureCode string,
	now time.Time,
) error {
	if failureCode == "" {
		return ErrFailureCodeRequired
	}

	if err := w.transitionTo(StatusFailed); err != nil {
		return err
	}

	w.FailureCode = new(failureCode)
	w.UpdatedAt = now

	return nil
}
