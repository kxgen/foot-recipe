// maps incoming data to the ui recipe card
export function mapRecipeCard(recipe) {
    // if no data then exit 
    if (!recipe) return null;
    
    // find feature image or fallback to the first image
    const featured =
        recipe.recipe_images?.find((img) => img.is_featured) ??
        recipe.recipe_images?.[0];

    return {
        id: recipe.id,
        title: recipe.title,
        slug: recipe.slug,
        creator: recipe.user?.username ?? "Unknown",
        creatorId: recipe.user?.id ?? recipe.user_id,
        creatorSlug: recipe.user?.slug,
        ratingCount: recipe.rating_count ?? 0,
        avgRating: Number(recipe.avg_rating ?? 0),
        timeMinutes: recipe.total_time_minutes ?? 0,
        imageSrc: featured?.url,
        category: recipe.category?.name,
        categoryId: recipe.category?.id ?? recipe.category_id,
        categorySlug: recipe.category?.slug,
        price: Number(recipe.price ?? 0),
        isFree: Number(recipe.price ?? 0) === 0,
        recipeLink: `/recipes/${recipe.id}`,
    };
}
