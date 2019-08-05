package auth

import (
	"context"
	"encoding/base64"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Authenticator exposes a function for authenticating requests.
type Authenticator struct {
	Username string
	Password string
	Token    string
}

// BasicAuth ..
func (a Authenticator) BasicAuth(ctx context.Context) (context.Context, error) {
	auth, err := ExtractHeader(ctx, "authorization")
	if err != nil {
		return ctx, err
	}

	a.Username = "apicentral"
	a.Password = "changeme"

	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return ctx, status.Error(codes.Unauthenticated, `missing "Basic " prefix in "Authorization" header`)
	}

	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, `invalid base64 in header`)
	}

	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return ctx, status.Error(codes.Unauthenticated, `invalid basic auth format`)
	}

	user, password := cs[:s], cs[s+1:]
	if user != a.Username || password != a.Password {
		return ctx, status.Error(codes.Unauthenticated, "invalid user or password")
	}

	// Remove token from headers from here on
	// return purgeHeader(ctx, "authorization"), nil
	return ctx, nil
}

// APIKeyAuth ..
func (a Authenticator) APIKeyAuth(ctx context.Context) (context.Context, error) {
	auth, err := ExtractHeader(ctx, "authorization")
	if err != nil {
		return ctx, err
	}

	a.Token = "S4Z3HJ4WdBLhtmmvP4aOQjdariGaSJxF"

	const prefix = "Apikey "
	if !strings.HasPrefix(auth, prefix) {
		return ctx, status.Error(codes.Unauthenticated, `missing "Apikey " prefix in "Authorization" header`)
	}

	key := string(auth[len(prefix):])

	if key != a.Token {
		return ctx, status.Error(codes.Unauthenticated, "invalid APIKey")
	}

	// Remove token from headers from here on
	// return purgeHeader(ctx, "authorization"), nil
	return ctx, nil
}

// ExtractHeader ..
func ExtractHeader(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no headers in request")
	}

	authHeaders, ok := md[header]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no authentication header in request")
	}

	if len(authHeaders) != 1 {
		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
	}

	return authHeaders[0], nil
}
