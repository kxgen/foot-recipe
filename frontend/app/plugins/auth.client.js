export default defineNuxtPlugin(() => {
  const { token, checkToken } = useAuth();
  let logoutTimer = null;

  if (process.client) {
    // Check immediately on startup
    checkToken();

    // Set up a watcher or a timer to handle expiration while the app is open
    watch(token, (newToken) => {
      if (logoutTimer) {
        clearTimeout(logoutTimer);
        logoutTimer = null;
      }

      if (newToken) {
        setupLogoutTimer(newToken);
      }
    }, { immediate: true });
  }

  function setupLogoutTimer(tokenValue) {
    try {
      const payload = JSON.parse(atob(tokenValue.split(".")[1]));
      const expiry = payload.exp * 1000;
      const delay = expiry - Date.now();

      if (delay > 0) {
        logoutTimer = setTimeout(() => {
          checkToken();
        }, delay + 1000); // Add 1s buffer
      } else {
        checkToken();
      }
    } catch (e) {
      console.error("Failed to parse token for auto-logout timer", e);
    }
  }
});