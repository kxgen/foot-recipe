import { 
  ApolloClient, 
  InMemoryCache, 
  HttpLink, 
  ApolloLink 
} from "@apollo/client";
import { SetContextLink } from "@apollo/client/link/context";
import { ErrorLink } from "@apollo/client/link/error";

export default defineNuxtPlugin(() => {
  const { token, logout } = useAuth();

  const httpLink = new HttpLink({
    uri: "http://localhost:8080/v1/graphql",
  });

  const authLink = new SetContextLink((prevContext, operation) => {
    if (!token.value) {
      return { 
        headers: prevContext.headers 
      };
    }
    return {
      headers: {
        ...prevContext.headers,
        authorization: `Bearer ${token.value}`,
      },
    };
  });

  const errorLink = new ErrorLink(({ graphQLErrors, networkError }) => {
    if (graphQLErrors) {
      for (const err of graphQLErrors) {
        if (
          err.extensions?.code === "invalid-jwt" ||
          err.message?.includes("JWTExpired") ||
          err.extensions?.status === 401
        ) {
          logout();
          return; 
        }
      }
    }
    if (networkError && 'statusCode' in networkError && networkError.statusCode === 401) {
      logout();
    }
  });

  const apolloClient = new ApolloClient({
    link: ApolloLink.from([errorLink, authLink, httpLink]),
    cache: new InMemoryCache(),
    connectToDevTools: process.dev,
  });

  return {
    provide: {
      apollo: apolloClient,
    },
  };
});