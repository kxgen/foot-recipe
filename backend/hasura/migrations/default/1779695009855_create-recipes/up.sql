CREATE TABLE recipes (
    id                INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id           INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id       INT REFERENCES categories(id) ON DELETE SET NULL,

    title             VARCHAR(255) NOT NULL,
    slug              VARCHAR(255) NOT NULL UNIQUE,
    description       TEXT,

    prep_time_minutes INT CHECK (prep_time_minutes > 0),

    difficulty        VARCHAR(20)
        CHECK (difficulty IN ('easy', 'medium', 'hard')),

    status            VARCHAR(20)
        NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'published', 'archived')),

    price             NUMERIC(10,2)
        NOT NULL DEFAULT 0
        CHECK (price >= 0),

    avg_rating        NUMERIC(3,2) NOT NULL DEFAULT 0,
    rating_count      INT NOT NULL DEFAULT 0,
    view_count        INT NOT NULL DEFAULT 0,

    published_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
