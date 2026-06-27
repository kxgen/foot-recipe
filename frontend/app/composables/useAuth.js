export const useAuth = () => {
  const token = useCookie("token", {
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    default: () => null,
  });

  const decodeJWT = (t) => {
    try {
      const base64 = t.split(".")[1];
      return JSON.parse(atob(base64));
    } catch {
      return null;
    }
  };

  const isLoggedIn = computed(() => {
    // does token exist?
    if (!token.value) return false;
    // is token.payload malformed?
    const payload = decodeJWT(token.value);
    if (!payload) return false;
    // is token expired?
    return payload.exp * 1000 > Date.now();
  });

  const login = (newToken) => {
    token.value = newToken;
    return navigateTo("/me");
  };

  const logout = () => {
    token.value = null;
    return navigateTo("/auth/login");
  };

  const checkToken = () => {
    if (!token.value) return;

    const payload = decodeJWT(token.value);
    if (!payload || payload.exp * 1000 <= Date.now()) {
      logout();
    }
  };

  return {
    token,
    isLoggedIn,
    login,
    logout,
    checkToken,
  };
};
