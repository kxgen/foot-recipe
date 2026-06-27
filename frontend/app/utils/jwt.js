export function getUserIdFromToken(token) {
  if (!token) return null;

  try {
    const payload = token.split(".")[1];
    if (!payload) return null;

    const decoded = JSON.parse(atob(payload.replace(/-/g, "+").replace(/_/g, "/")));
    const hasuraClaims = decoded["https://hasura.io/jwt/claims"];
    const id =
      hasuraClaims?.["x-hasura-user-id"] ??
      decoded.sub ??
      decoded.user_id;

    if (id == null || id === "") return null;
    return Number(id);
  } catch {
    return null;
  }
}
