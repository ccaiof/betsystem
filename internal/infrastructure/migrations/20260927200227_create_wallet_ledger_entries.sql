-- +goose Up
CREATE TYPE ledger_direction AS ENUM ('DEBIT', 'CREDIT');

CREATE TABLE
    wallet_ledger_entries (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        wallet_id UUID NOT NULL REFERENCES wallets (id),
        transaction_id UUID NOT NULL REFERENCES wager_transactions (id),
        direction ledger_direction NOT NULL,
        amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
        balance_before_minor BIGINT NOT NULL CHECK (balance_before_minor >= 0),
        balance_after_minor BIGINT NOT NULL CHECK (balance_after_minor >= 0),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        CONSTRAINT uq_ledger_wallet_transaction UNIQUE (wallet_id, transaction_id),
        CONSTRAINT ck_ledger_balance_math CHECK (
            (
                direction = 'CREDIT'
                AND balance_after_minor = balance_before_minor + amount_minor
            )
            OR (
                direction = 'DEBIT'
                AND balance_after_minor = balance_before_minor - amount_minor
            )
        )
    );

CREATE INDEX idx_ledger_wallet_id_created_at ON wallet_ledger_entries (wallet_id, created_at, id);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION forbid_ledger_mutation() RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'wallet_ledger_entries is append-only: % is not allowed', TG_OP;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_ledger_forbid_update BEFORE
UPDATE ON wallet_ledger_entries FOR EACH ROW EXECUTE FUNCTION forbid_ledger_mutation ();

CREATE TRIGGER trg_ledger_forbid_delete BEFORE DELETE ON wallet_ledger_entries FOR EACH ROW EXECUTE FUNCTION forbid_ledger_mutation ();

-- +goose Down
DROP TYPE IF EXISTS ledger_direction;

DROP TABLE IF EXISTS wallet_ledger_entries;
DROP FUNCTION IF EXISTS forbid_ledger_mutation();
