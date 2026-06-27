CREATE OR REPLACE FUNCTION update_recipe_rating_stats()
RETURNS TRIGGER AS $$
DECLARE
    target_recipe_id INT;
BEGIN
    IF TG_OP = 'DELETE' THEN
        target_recipe_id := OLD.recipe_id;
    ELSE
        target_recipe_id := NEW.recipe_id;
    END IF;

    UPDATE recipes
    SET 
        avg_rating = COALESCE((SELECT ROUND(AVG(score), 2) FROM recipe_ratings WHERE recipe_id = target_recipe_id), 0.00),
        rating_count = (SELECT COUNT(*) FROM recipe_ratings WHERE recipe_id = target_recipe_id),
        updated_at = NOW()
    WHERE id = target_recipe_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
