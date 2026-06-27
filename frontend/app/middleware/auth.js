// redirect to login if not loggedin
// runs everytime route nav happens
export default defineNuxtRouteMiddleware(() => {
  const { isLoggedIn } = useAuth();

  if (!isLoggedIn.value) {
    return navigateTo("/auth/login");
  }
});