package azurecloud

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

type authResult struct {
	accessToken    string
	expiresOn      time.Time
	grantedScopes  []string
	declinedScopes []string
}

func clientAssertionBearerAuthorizerCallback(tenantID, resource string) (*autorest.BearerAuthorizer, error) {
	clientID := os.Getenv("AZURE_CLIENT_ID")
	tokenFilePath := os.Getenv("AZURE_FEDERATED_TOKEN_FILE")
	authorityHost := os.Getenv("AZURE_AUTHORITY_HOST")

	signedAssertion, err := readJWTFromFS(tokenFilePath)
	if err != nil {
		return nil, errors.New(err.Error() + " failed to read service account token")
	}
	cred, err := confidential.NewCredFromAssertion(signedAssertion)
	if err != nil {
		return nil, errors.New(err.Error() + " failed to create confidential creds")
	}
	confidentialClientApp, err := confidential.New(
		clientID,
		cred,
		confidential.WithAuthority(fmt.Sprintf("%s%s/oauth2/token", authorityHost, tenantID)))
	if err != nil {
		return nil, errors.New(err.Error() + " failed to create confidential client app")
	}

	resource = strings.TrimSuffix(resource, "/")
	if !strings.HasSuffix(resource, ".default") {
		resource += "/.default"
	}

	result, err := confidentialClientApp.AcquireTokenByCredential(context.Background(), []string{resource})
	if err != nil {
		return nil, errors.New(err.Error() + " failed to acquire token")
	}

	return autorest.NewBearerAuthorizer(authResult{
		accessToken:    result.AccessToken,
		expiresOn:      result.ExpiresOn,
		grantedScopes:  result.GrantedScopes,
		declinedScopes: result.DeclinedScopes,
	}), nil
}

// OAuthToken implements the OAuthTokenProvider interface.  It returns the current access token.
func (ar authResult) OAuthToken() string {
	return ar.accessToken
}

// readJWTFromFS reads the jwt from file system
func readJWTFromFS(tokenFilePath string) (string, error) {
	token, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func Authorizer() *autorest.BearerAuthorizerCallback {
	return autorest.NewBearerAuthorizerCallback(nil, clientAssertionBearerAuthorizerCallback)
}
