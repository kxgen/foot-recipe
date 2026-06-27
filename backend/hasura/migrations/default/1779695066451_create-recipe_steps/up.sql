CREATE TABLE recipe_steps (
    id               INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    recipe_id        INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    step_number      INT NOT NULL CHECK (step_number > 0),
    instruction      TEXT NOT NULL,
    duration_minutes INT,
    image_url        TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (recipe_id, step_number)
);
