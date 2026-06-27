CREATE TABLE recipe_ingredients (
    id         INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    recipe_id  INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    name       VARCHAR(150) NOT NULL,
    quantity   NUMERIC(8,2),
    unit       VARCHAR(50),
    notes      TEXT,
    sort_order INT NOT NULL DEFAULT 0
);
