CREATE OR REPLACE FUNCTION update_recipe_like_stats()
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
        like_count = (SELECT COUNT(*) FROM recipe_likes WHERE recipe_id = target_recipe_id),
        updated_at = NOW()
    WHERE id = target_recipe_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
