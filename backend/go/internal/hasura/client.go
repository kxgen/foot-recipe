package hasura

import (
	"net/http"
	"server/internal/config"
	graphql "github.com/hasura/go-graphql-client"
)

/*
 * It's 1 Job: It automatically attach a password (an admin secret) onto every
 * single request the server sends to Hasura
 */
 
func NewClient(cfg config.Config) *graphql.Client {
	httpClient := &http.Client{
		Transport: &transport{
			underlying: http.DefaultTransport,
			adminSecret: cfg.HasuraAdminSecret,
		},
	}
	client := graphql.NewClient(
		cfg.HasuraURL,
		httpClient,
	)
	return client
}


type transport struct {
	underlying http.RoundTripper
	adminSecret string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.adminSecret != "" {
		req.Header.Set("x-hasura-admin-secret", t.adminSecret)
	}
	return t.underlying.RoundTrip(req)
}

/*
 * It initializes and returns a brand new Hasura GraphQL client
 * that is fully authenticated
 */