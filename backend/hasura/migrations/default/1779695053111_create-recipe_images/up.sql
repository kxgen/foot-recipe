CREATE TABLE recipe_images (
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    recipe_id   INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    url         TEXT NOT NULL,
    alt_text    VARCHAR(255),
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX one_featured_image_per_recipe
ON recipe_images(recipe_id)
WHERE is_featured = TRUE;
