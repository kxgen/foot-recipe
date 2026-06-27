import { gql } from "@apollo/client";

/*
======================================================
RECIPE MUTATIONS
======================================================
*/
export const INSERT_RECIPE = gql`
  mutation InsertRecipe($object: recipes_insert_input!) {
    insert_recipes_one(object: $object) {
      id
      title
      slug
    }
  }
`;
export const DELETE_RECIPE = gql`
    mutation DeleteRecipe($id: Int!) {
        delete_recipes_by_pk(id: $id) {
            id
        }
    }
`;
export const UPDATE_RECIPE_WITH_RELATIONS = gql`
    mutation UpdateRecipeWithRelations(
        $id: Int!
        $set: recipes_set_input!
        $steps: [recipe_steps_insert_input!]!
        $ingredients: [recipe_ingredients_insert_input!]!
        $images: [recipe_images_insert_input!]!
    ) {
        update_recipes_by_pk(pk_columns: { id: $id }, _set: $set) {
            id
        }
        delete_recipe_steps(where: { recipe_id: { _eq: $id } }) {
            affected_rows
        }
        delete_recipe_ingredients(where: { recipe_id: { _eq: $id } }) {
            affected_rows
        }
        delete_recipe_images(where: { recipe_id: { _eq: $id } }) {
            affected_rows
        }
        insert_recipe_steps(objects: $steps) {
            affected_rows
        }
        insert_recipe_ingredients(objects: $ingredients) {
            affected_rows
        }
        insert_recipe_images(objects: $images) {
            affected_rows
        }
    }
`;
export const INSERT_RECIPE_LIKE = gql`
    mutation InsertRecipeLike($recipeId: Int!) {
        insert_recipe_likes_one(
            object: { recipe_id: $recipeId }
            on_conflict: { constraint: recipe_likes_pkey, update_columns: [] }
        ) {
            recipe_id
            user_id
            # Add this so Apollo automatically updates your UI cache safely
            recipe {
                id
                like_count
            }
        }
    }
`;
export const DELETE_RECIPE_LIKE = gql`
    mutation DeleteRecipeLike($recipeId: Int!, $userId: Int!) {
        delete_recipe_likes(
            where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }
        ) {
            affected_rows
            # Add this so Apollo automatically updates your UI cache safely
            returning {
                recipe {
                    id
                    like_count
                }
            }
        }
    }
`;
export const INSERT_RECIPE_BOOKMARK = gql`
    mutation InsertRecipeBookmark($recipeId: Int!) {
        insert_recipe_bookmarks_one(
            object: { recipe_id: $recipeId }
            on_conflict: {
                constraint: recipe_bookmarks_pkey
                update_columns: []
            }
        ) {
            recipe_id
            user_id
        }
    }
`;
export const DELETE_RECIPE_BOOKMARK = gql`
    mutation DeleteRecipeBookmark($recipeId: Int!, $userId: Int!) {
        delete_recipe_bookmarks(
            where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }
        ) {
            affected_rows
        }
    }
`;
export const INITIATE_PURCHASE = gql`
    mutation InitiatePurchase($recipeId: Int!, $userId: Int!) {
        initiatePurchase(recipeId: $recipeId, userId: $userId) {
            checkoutUrl
        }
    }
`;


export const INSERT_RECIPE_COMMENT = gql`
    mutation InsertRecipeComment($recipeId: Int!, $body: String!) {
        insert_recipe_comments_one(
            object: { recipe_id: $recipeId, body: $body }
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
    }
`;
// update, insert 
export const UPSERT_RECIPE_RATING = gql`
    mutation UpsertRecipeRating($recipeId: Int!, $score: smallint!) {
        insert_recipe_ratings_one(
            # create new with this data
            object: {
                recipe_id: $recipeId
                score: $score
            }
            # if exist, overwrite score old with new
            on_conflict: {
                constraint: recipe_ratings_pkey
                update_columns: [score]
            }
        ) {
            recipe_id
            user_id
            score
        }
    }
`;
// update recipe comment by PK
export const UPDATE_RECIPE_COMMENT = gql`
    mutation UpdateRecipeComment($id: Int!, $body: String!) {
        update_recipe_comments_by_pk(
            pk_columns: { id: $id }
            _set: { body: $body }
        ) {
            id
            body
            created_at
        }
    }
`;
export const DELETE_RECIPE_COMMENT = gql`
    mutation DeleteRecipeComment($id: Int!) {
        delete_recipe_comments_by_pk(id: $id) {
            id
        }
    }
`;

/*
======================================================
USER MUTATIONS
======================================================
*/

// Login
export const LOGIN_USER = gql`
    mutation LoginUser($username: String!, $password: String!) {
        loginUser(username: $username, password: $password) {
            id
            token
        }
    }
`;

// Register
export const REGISTER_USER = gql`
    mutation RegisterUser($username: String!, $email: String!, $password: String!) {
        registerUser(username: $username, email: $email, password: $password) {
            id
            token
        }
    }
`;

// Update
export const UPDATE_USER_PROFILE = gql`
  mutation UpdateUserProfile($id: Int!, $set: users_set_input!) {
    update_users_by_pk(pk_columns: { id: $id }, _set: $set) {
      id
      username
      email
      avatar_url
      bio
    }
  }
`;

// Update Password
export const UPDATE_PASSWORD = gql`
    mutation UpdatePassword($newPassword: String!) {
        updatePassword(newPassword: $newPassword) {
            success
            message
        }
    }
`;