import { gql } from "@apollo/client";
import { RECIPE_CARD_FIELDS } from "./fragments.js";

/*
======================================================
RECIPE QUERIES
======================================================
*/

export const GET_RECIPES = gql`
  query GetRecipes(
    $where: recipes_bool_exp
    $limit: Int
    $offset: Int
    $orderBy: [recipes_order_by!]
  ) {
    recipes(where: $where, limit: $limit, offset: $offset, order_by: $orderBy) {
      ...RecipeCardFields
    }
  }
  ${RECIPE_CARD_FIELDS}
`;

export const GET_LATEST_RECIPES = gql`
  query GetRecipes(
    $where: recipes_bool_exp
    $limit: Int
    $offset: Int
    $orderBy: [recipes_order_by!]
  ) {
    recipes(where: $where, limit: $limit, offset: $offset, order_by: $orderBy) {
      ...RecipeCardFields
    }
  }
  ${RECIPE_CARD_FIELDS}
`;

export const GET_RECIPE_BY_ID = gql`
  query GetRecipeById($id: Int!) {
    recipes_by_pk(id: $id) {
      id
      title
      slug
      description
      total_time_minutes
      difficulty
      avg_rating
      rating_count
      user_id
      category_id
      price
      like_count
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
      recipe_images(order_by: { sort_order: asc }) {
        id
        url
        alt_text
        is_featured
        sort_order
      }
      recipe_steps(order_by: { step_number: asc }) {
        id
        step_number
        instruction
        duration_minutes
        image_url
      }
      recipe_ingredients(order_by: { sort_order: asc }) {
        id
        name
        quantity
        unit
        notes
        sort_order
      }
    }
  }
`;

export const GET_CATEGORIES = gql`
  query GetCategories($limit: Int) {
    categories(order_by: { sort_order: asc }, limit: $limit) {
      id
      name
      slug
      description
      recipes_aggregate {
        aggregate {
          count
        }
      }
    }
  }
`;

export const GET_CATEGORY_BY_SLUG = gql`
  query GetCategoryBySlug($slug: String!) {
    categories(where: { slug: { _eq: $slug } }, limit: 1) {
      id
      name
      slug
      description
      recipes(order_by: { created_at: desc }) {
        ...RecipeCardFields
      }
    }
  }
  ${RECIPE_CARD_FIELDS}
`;

export const GET_FILTER_OPTIONS = gql`
  query GetFilterOptions {
    categories(order_by: { name: asc }) {
      id
      name
      slug
    }
    users(where: { recipes: {} }, order_by: { username: asc }) {
      id
      username
    }
    recipe_ingredients(distinct_on: name, order_by: { name: asc }, limit: 50) {
      name
    }
  }
`;

export const GET_RECIPE_LIKE_COUNT = gql`
    query GetRecipeLikeCount($id: Int!) {
        recipes_by_pk(id: $id) {
            id
            recipe_likes_aggregate {
                aggregate {
                    count
                }
            }
        }
    }
`;


export const CHECK_USER_RECIPE_FLAGS = gql`
    query CheckUserRecipeFlags($recipeId: Int!, $userId: Int!) {
        recipe_likes(
            where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }
        ) {
            recipe_id
        }
        recipe_bookmarks(
            where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }
        ) {
            recipe_id
        }
        recipe_ratings(
            where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }
        ) {
            score
        }
        recipe_purchases(
            where: {
                recipe_id: { _eq: $recipeId }
                user_id: { _eq: $userId }
                status: { _eq: "completed" }
            }
        ) {
            id
        }
    }
`;

export const GET_RECIPE_COMMENTS = gql`
    query GetRecipeComments($recipeId: Int!) {
        recipe_comments(
            where: { recipe_id: { _eq: $recipeId } }
            order_by: { created_at: desc }
        ) {
            id
            body
            created_at
            user {
                id
                username
                avatar_url
            }
        }
        recipe_ratings(where: { recipe_id: { _eq: $recipeId } }) {
            user_id
            score
        }
    }
`;

export const CHECK_USER_BOOKMARK = gql`
  query CheckUserBookmark($recipeId: Int!, $userId: Int!) {
    recipe_bookmarks(
      where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }
    ) {
      recipe_id
    }
  }
`;


/*
======================================================
USER QUERIES
======================================================
*/

export const GET_CREATORS = gql`
  query GetCreators {
    users(
      where: { recipes: {} }
      order_by: { username: asc }
    ) {
      id
      username
      slug
      avatar_url
      bio
      recipes_aggregate {
        aggregate {
          count
        }
      }
    }
  }
`;

export const GET_USER_PROFILE = gql`
  query GetUserProfile($id: Int!) {
    users_by_pk(id: $id) {
      id
      username
      slug
      email
      avatar_url
      bio
    }
  }
`;

export const GET_MY_RECIPES = gql`
  query GetMyRecipes($userId: Int!) {
    recipes(
      where: { user_id: { _eq: $userId } }
      order_by: { created_at: desc }
    ) {
      ...RecipeCardFields
    }
  }
  ${RECIPE_CARD_FIELDS}
`;

export const GET_MY_BOOKMARKS = gql`
  query GetMyBookmarks($userId: Int!) {
    recipe_bookmarks(
      where: { user_id: { _eq: $userId } }
      order_by: { created_at: desc }
    ) {
      recipe {
        ...RecipeCardFields
      }
    }
  }
  ${RECIPE_CARD_FIELDS}
`;

export const GET_MY_PURCHASES = gql`
  query GetMyPurchases($userId: Int!) {
    recipe_purchases(
      where: { user_id: { _eq: $userId }, status: { _eq: "completed" } }
      order_by: { purchased_at: desc }
    ) {
      id
      amount_paid
      purchased_at
      recipe {
        ...RecipeCardFields
      }
    }
  }
  ${RECIPE_CARD_FIELDS}
`;

// input: slug of creator, output: creators info
export const GET_CREATOR_BY_SLUG = gql`
  query GetCreatorBySlug($slug: String!) {
    users(where: { slug: { _eq: $slug } }, limit: 1) {
      id
      username
      slug
      avatar_url
      bio
      recipes(
        order_by: { created_at: desc }
      ) {
        ...RecipeCardFields
      }
    }
  }
  ${RECIPE_CARD_FIELDS}
`;