package authentication

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"main/pkg/envvar"
	"net/http"
	"os"
	"strings"

	oidc "github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator(host string) (*Authenticator, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://login.microsoftonline.com/"+os.Getenv("TENANT_ID")+"/v2.0")
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	scheme := envvar.GetEnvVar("SCHEME", "https")

	conf := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  fmt.Sprint(scheme, "://", host, "/login/azure/callback"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}

func VerifyAccessToken(r *http.Request) (*jwt.Token, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return nil, fmt.Errorf("no authorization header found")
	}
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	keySet, err := jwk.Fetch(r.Context(), "https://login.microsoftonline.com/common/discovery/v2.0/keys")
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {

		}

		publickey := &rsa.PublicKey{}
		err = keys.Raw(publickey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey")
		}

		return publickey, nil
	})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return token, nil
}
