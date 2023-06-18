package auth

import (
	"context"
	"encoding/base64"
	"fmt"
)

type basicAuth struct {
	username string
	password string
}

func New(username, password string) basicAuth {
	return basicAuth{
		username: username,
		password: password,
	}
}

func (b basicAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	auth := fmt.Sprintf("%s:%s", b.username, b.password)
	encAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + encAuth,
	}, nil
}

func (b basicAuth) RequireTransportSecurity() bool {
	return true
}
