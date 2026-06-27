INSERT INTO recipes (
    user_id, category_id, title, slug, description,
    prep_time_minutes, difficulty, price,
    avg_rating, rating_count, view_count
) VALUES (
    (SELECT id FROM users WHERE username = 'chef_abebe'),
    (SELECT id FROM categories WHERE slug = 'breakfast'),
    'Simple Scrambled Eggs',
    'simple-scrambled-eggs',
    'Soft, creamy scrambled eggs done in under 10 minutes. A perfect quick breakfast.',
    5,
    'easy',
    0.00,
    4.50,
    30,
    200
);

-- RECIPE INGREDIENTS
INSERT INTO recipe_ingredients (recipe_id, name, quantity, unit, notes, sort_order) VALUES
    ((SELECT id FROM recipes WHERE slug = 'simple-scrambled-eggs'), 'Eggs',           3,  NULL,   'Large',              1),
    ((SELECT id FROM recipes WHERE slug = 'simple-scrambled-eggs'), 'Butter',         1,  'tbsp', NULL,                 2),
    ((SELECT id FROM recipes WHERE slug = 'simple-scrambled-eggs'), 'Milk',           2,  'tbsp', 'Or heavy cream',     3),
    ((SELECT id FROM recipes WHERE slug = 'simple-scrambled-eggs'), 'Salt',           1,  'pinch',NULL,                 4),
    ((SELECT id FROM recipes WHERE slug = 'simple-scrambled-eggs'), 'Black pepper',   1,  'pinch','Freshly ground',     5);

-- RECIPE STEPS
INSERT INTO recipe_steps (recipe_id, step_number, instruction, duration_minutes) VALUES
    (
        (SELECT id FROM recipes WHERE slug = 'simple-scrambled-eggs'),
        1,
        'Crack the eggs into a bowl, add the milk, salt and pepper. Whisk until fully combined.',
        2
    ),
    (
        (SELECT id FROM recipes WHERE slug = 'simple-scrambled-eggs'),
        2,
        'Melt the butter in a non-stick pan over low to medium heat.',
        1
    ),
    (
        (SELECT id FROM recipes WHERE slug = 'simple-scrambled-eggs'),
        3,
        'Pour in the egg mixture. Let it sit for a few seconds then gently fold with a spatula, pushing the eggs from the edges to the center. Remove from heat while still slightly soft.',
        3
    );
