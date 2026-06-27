CREATE TABLE users (
    id         INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username   VARCHAR(50)  NOT NULL UNIQUE,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   TEXT         NOT NULL,
    avatar_url TEXT,
    bio        TEXT,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);


CREATE TABLE categories (
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    sort_order  INT NOT NULL DEFAULT 0
);

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

    price             NUMERIC(10,2)
        NOT NULL DEFAULT 0
        CHECK (price >= 0),

    avg_rating        NUMERIC(3,2) NOT NULL DEFAULT 0,
    rating_count      INT NOT NULL DEFAULT 0,
    view_count        INT NOT NULL DEFAULT 0,

    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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


CREATE TABLE recipe_ingredients (
    id         INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    recipe_id  INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    name       VARCHAR(150) NOT NULL,
    quantity   NUMERIC(8,2),
    unit       VARCHAR(50),
    notes      TEXT,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE recipe_likes (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, recipe_id)
);

CREATE TABLE recipe_bookmarks (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, recipe_id)
);

CREATE TABLE recipe_comments (
    id         INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id  INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    parent_id  INT REFERENCES recipe_comments(id) ON DELETE CASCADE,

    body       TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE recipe_ratings (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,

    score SMALLINT NOT NULL CHECK (score BETWEEN 1 AND 5),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, recipe_id)
);

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
