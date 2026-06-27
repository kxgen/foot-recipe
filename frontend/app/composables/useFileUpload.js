import { gql } from "@apollo/client/core";

const MAX_IMAGE_BYTES = 2 * 1024 * 1024;

export function useFileUpload() {
  function validateImageFile(file) {
    if (!file) return "No file selected.";
    if (!file.type.startsWith("image/")) return "Please choose an image file.";
    if (file.size > MAX_IMAGE_BYTES) return "Image must be 2 MB or smaller.";
    return null;
  }

  function fileToDataUrl(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(reader.result);
      reader.onerror = () => reject(new Error("Could not read the image file."));
      reader.readAsDataURL(file);
    });
  }

  async function readImageFile(file) {
    const validationError = validateImageFile(file);
    if (validationError) throw new Error(validationError);
    return fileToDataUrl(file);
  }

  async function uploadRecipeImage(file) {
    const dataUrl = await readImageFile(file);
    const { $apollo } = useNuxtApp();
    
    const { data } = await $apollo.mutate({
      mutation: gql`
        mutation UploadRecipeImage(
          $imageName: String!
          $base64Image: String!
          $recipeId: Int
          $tempFolder: String
          $filename: String
        ) {
          uploadRecipeImage(
            imageName: $imageName
            base64Image: $base64Image
            recipeId: $recipeId
            tempFolder: $tempFolder
            filename: $filename
          ) {
            success
            imageUrl
          }
        }
      `,
      variables: {
        imageName: file.name,
        base64Image: dataUrl,
        recipeId: null,
        tempFolder: null,
        filename: null,
      },
    });

    if (!data?.uploadRecipeImage?.success) {
      throw new Error("Failed to upload image to server.");
    }

    return data.uploadRecipeImage.imageUrl; // The public URL from Go server
  }

  async function uploadAvatar(file, { userId, slug }) {
    const dataUrl = await readImageFile(file);
    const { $apollo } = useNuxtApp();

    const { data } = await $apollo.mutate({
      mutation: gql`
        mutation UploadAvatar($userId: bigint!, $slug: String!, $imageName: String!, $base64Image: String!) {
          uploadAvatar(userId: $userId, slug: $slug, imageName: $imageName, base64Image: $base64Image) {
            success
            avatarUrl
          }
        }
      `,
      variables: {
        userId: userId,
        slug: slug,
        imageName: file.name,
        base64Image: dataUrl,
      },
    });

    if (!data?.uploadAvatar?.success) {
      throw new Error("Failed to upload avatar to server.");
    }

    return data.uploadAvatar.avatarUrl;
  }

  return {
    validateImageFile,
    readImageFile,
    uploadRecipeImage,
    uploadAvatar,
    maxImageBytes: MAX_IMAGE_BYTES,
  };
}

