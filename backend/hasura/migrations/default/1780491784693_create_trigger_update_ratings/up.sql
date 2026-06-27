CREATE TRIGGER trigger_update_recipe_ratings
AFTER INSERT OR UPDATE OR DELETE ON recipe_ratings
FOR EACH ROW
EXECUTE FUNCTION update_recipe_rating_stats();
