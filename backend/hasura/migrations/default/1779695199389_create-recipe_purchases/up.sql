CREATE TABLE recipe_purchases (
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    user_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id   INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,

    payment_provider        VARCHAR(50),
    provider_transaction_id VARCHAR(255),

    amount_paid  NUMERIC(10,2) NOT NULL, -- ETB only

    status       VARCHAR(20)
        NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),

    purchased_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, recipe_id)
);
