CREATE TABLE recipe_comments (
    id         INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id  INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    parent_id  INT REFERENCES recipe_comments(id) ON DELETE CASCADE,

    body       TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
