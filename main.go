package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
)

const (
	// AzureErrorMsgDefaultAzureCredential is a substring included in messages of errors returned by
	// the Azure SDK when any problem occurs when it tries to authenticate using
	// DefaultAzureCredential.
	AzureErrorMsgDefaultAzureCredential = "DefaultAzureCredential: "

	// AzureErrorCodeInvalidClientSecretProvided is a substring included in the message of the error
	// returned by the Azure SDK when the Azure SDK knew how to authenticate because the environment
	// variables required by EnvironmentCredential were provided, but CLIENT_SECRET had an
	// incorrect value.
	AzureErrorCodeInvalidClientSecretProvided = "AADSTS7000215"
)

func isDefaultAzureCredentialError(err error) bool {
	return strings.Contains(err.Error(), AzureErrorMsgDefaultAzureCredential)
}

func isInvalidClientSecretProvidedError(err error) bool {
	return strings.Contains(err.Error(), AzureErrorCodeInvalidClientSecretProvided)
}

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	subID := os.Getenv("SUB_ID")
	client, err := armauthorization.NewRoleAssignmentsClient(subID, cred, nil)
	if err != nil {
		panic(err)
	}

	pager := client.NewListForSubscriptionPager(nil)
	var ras []*armauthorization.RoleAssignment
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			if isDefaultAzureCredentialError(err) {
				panic(fmt.Errorf("DefaultAzureCredential error%s: %w", os.Getenv("DEFAULT_AZURE_CREDENTIAL_ERROR_HELP_MSG"), err))
			}

			if isInvalidClientSecretProvidedError(err) {
				panic(fmt.Errorf("invalid client secret provided%s: %w", os.Getenv("INVALID_CLIENT_SECRET_PROVIDED_HELP_MSG"), err))
			}

			var rerr *azcore.ResponseError
			if errors.As(err, &rerr) {
				if rerr.ErrorCode == "AuthorizationFailed" {
					panic(fmt.Errorf("client unauthorized%s: %w", os.Getenv("CLIENT_UNAUTHORIZED_HELP_MSG"), rerr))
				}
			}

			panic(fmt.Errorf("unknown error: %w", err))
		}
		ras = append(ras, page.Value...)
	}

	fmt.Println("Done!")
}
