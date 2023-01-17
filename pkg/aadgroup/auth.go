package aadgroup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/dfds/manage-aadgroup-members/pkg/config"
	"github.com/dfds/manage-aadgroup-members/pkg/logging"
	msazureauth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	// GraphScope The OAuth2 scope used for dealing with MS Graph API
	GraphScope string = "https://graph.microsoft.com/.default"
)

func getGraphClient() (*msgraphsdk.GraphServiceClient, error) {
	log := logging.GetLogger()
	cred, err := azidentity.NewClientSecretCredential(config.TenantId(), config.ClientId(),
		config.ClientSecret(), &azidentity.ClientSecretCredentialOptions{})

	if err != nil {
		log.Error(err)
		return nil, err
	}

	auth, err := msazureauth.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{GraphScope})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return msgraphsdk.NewGraphServiceClient(adapter), nil
}

// GetBearerToken returns a token from Azure AD
func GetBearerToken() (string, error) {
	log := logging.GetLogger()
	address := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", config.TenantId())

	data := url.Values{
		"client_id":     {config.ClientId()},
		"client_secret": {config.ClientSecret()},
		"grant_type":    {"client_credentials"},
		"scope":         {GraphScope},
	}

	resp, err := http.PostForm(address, data)

	if err != nil {
		log.Error(err)
		return "", err
	}

	var res map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return fmt.Sprintf("Bearer %v", res["access_token"]), nil
}
