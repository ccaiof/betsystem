-- +goose Up
CREATE TABLE
    IF NOT EXISTS wallets (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        player_id UUID NOT NULL,
        currency CHAR(3) NOT NULL,
        balance_minor BIGINT NOT NULL DEFAULT 0 CHECK (balance_minor >= 0),
        version BIGINT NOT NULL DEFAULT 1 CHECK (balance_minor >= 1),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        CONSTRAINT uq_wallets_player_currency UNIQUE (player_id, currency)
    );

CREATE TRIGGER trg_wallets_updated_at BEFORE
UPDATE ON wallets FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TABLE IF EXISTS wallets;