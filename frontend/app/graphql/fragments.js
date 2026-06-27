import { gql } from "@apollo/client";

export const RECIPE_CARD_FIELDS = gql`
  fragment RecipeCardFields on recipes {
    id
    title
    slug
    total_time_minutes
    avg_rating
    rating_count
    user_id
    category_id
    price
    user {
      id
      username
      slug
      avatar_url
    }
    category {
      id
      name
      slug
    }
    recipe_images(where: { is_featured: { _eq: true } }, limit: 1) {
      url
      is_featured
    }
  }
`;
