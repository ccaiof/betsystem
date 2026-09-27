-- +goose Up
CREATE TYPE wager_kind AS ENUM (
    'OPENING',
    'BET',
    'WIN',
    'LOSS',
    'REFUND',
    'ROLLBACK'
);

CREATE TYPE wager_status AS ENUM (
    'PENDING',
    'PENDING_REFERENCE',
    'PROCESSED',
    'REJECTED',
    'FAILED'
);

CREATE TABLE
    IF NOT EXISTS wager_transactions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        provider_id TEXT,
        external_transaction_id TEXT,
        idempotency_key TEXT,
        payload_hash TEXT,
        round_id TEXT,
        game_id TEXT,
        wallet_id UUID NOT NULL REFERENCES wallets (id),
        player_id UUID NOT NULL,
        kind wager_kind NOT NULL,
        amount_minor BIGINT NOT NULL,
        currency CHAR(3) NOT NULL,
        -- for REFUND/ROLLBACK
        reference_external_transaction_id TEXT,
        reference_transaction_id UUID REFERENCES wager_transactions (id),
        --
        status wager_status NOT NULL DEFAULT 'PENDING',
        failure_code TEXT,
        retry_count INT NOT NULL DEFAULT 0,
        next_retry_at TIMESTAMPTZ,
        financial_result_balance_minor BIGINT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        processed_at TIMESTAMPTZ,
        -- CHECKS
        CONSTRAINT ck_wager_external_metadata CHECK (
            (
                kind = 'OPENING'
                AND provider_id IS NULL
                AND external_transaction_id IS NULL
                AND idempotency_key IS NULL
                AND payload_hash IS NULL
                AND round_id IS NULL
                AND game_id IS NULL
                AND reference_external_transaction_id IS NULL
            )
            OR (
                kind <> 'OPENING'
                AND provider_id IS NOT NULL
                AND external_transaction_id IS NOT NULL
                AND idempotency_key IS NOT NULL
                AND payload_hash IS NOT NULL
                AND round_id IS NOT NULL
                AND game_id IS NOT NULL
            )
        ),
        CONSTRAINT ck_wager_reference_required CHECK (
            (
                kind IN ('REFUND', 'ROLLBACK')
                AND reference_external_transaction_id IS NOT NULL
                AND reference_transaction_id IS NOT NULL
            )
            OR (
                kind NOT IN ('REFUND', 'ROLLBACK')
                AND reference_external_transaction_id IS NULL
                AND reference_transaction_id IS NULL
            )
        ),
        CONSTRAINT ck_wager_amount_by_kind CHECK (
            (
                kind = 'LOSS'
                AND amount_minor = 0
            )
            OR (
                kind <> 'LOSS'
                AND amount_minor > 0
            )
        )
    );

-- INDEXES
CREATE UNIQUE INDEX uq_wager_provider_external ON wager_transactions (provider_id, external_transaction_id)
WHERE
    provider_id IS NOT NULL;

CREATE UNIQUE INDEX uq_wager_idempotency_key ON wager_transactions (idempotency_key)
WHERE
    idempotency_key IS NOT NULL;

CREATE UNIQUE INDEX uq_wager_opening_per_wallet ON wager_transactions (wallet_id)
WHERE
    kind = 'OPENING';

CREATE UNIQUE INDEX uq_wager_single_successful_reversal ON wager_transactions (reference_transaction_id, kind)
WHERE
    status = 'PROCESSED'
    AND kind IN ('REFUND', 'ROLLBACK');

CREATE INDEX idx_wager_wallet_id ON wager_transactions (wallet_id);

CREATE INDEX idx_wager_reference_transaction_id ON wager_transactions (reference_transaction_id);

CREATE INDEX idx_wager_pending_reference_retry ON wager_transactions (next_retry_at)
WHERE
    status = 'PENDING_REFERENCE';

CREATE TRIGGER trg_wager_transactions_updated_at BEFORE
UPDATE ON wager_transactions FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_wager_currency_matches_wallet()
RETURNS TRIGGER AS $$
DECLARE
    wallet_currency CHAR(3);
BEGIN
    SELECT currency
    INTO wallet_currency
    FROM wallets
    WHERE id = NEW.wallet_id;

    IF wallet_currency IS NULL THEN
        RAISE EXCEPTION 'wallet % not found', NEW.wallet_id;
    END IF;

    IF wallet_currency <> NEW.currency THEN
        RAISE EXCEPTION
            'transaction currency % does not match wallet currency %',
            NEW.currency,
            wallet_currency;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_wager_currency_matches_wallet
    BEFORE INSERT OR UPDATE ON wager_transactions
    FOR EACH ROW EXECUTE FUNCTION check_wager_currency_matches_wallet();

CREATE TRIGGER trg_wager_forbid_terminal_transition
    BEFORE UPDATE ON wager_transactions
    FOR EACH ROW EXECUTE FUNCTION forbid_terminal_transition();

-- +goose Down
DROP TABLE IF EXISTS wager_transactions;

DROP FUNCTION IF EXISTS forbid_terminal_transition();
DROP FUNCTION IF EXISTS check_wager_currency_matches_wallet();

DROP TYPE IF EXISTS wager_status;
DROP TYPE IF EXISTS wager_kind;