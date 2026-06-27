CREATE OR REPLACE FUNCTION recipe_total_time(recipe_row recipes)
RETURNS INT AS $$
SELECT COALESCE(recipe_row.prep_time_minutes, 0) + COALESCE(SUM(duration_minutes), 0)::INT
FROM recipe_steps
WHERE recipe_id = recipe_row.id;
$$ LANGUAGE sql STABLE;
