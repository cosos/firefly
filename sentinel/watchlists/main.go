package main

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights"
	"github.com/cosos/firefly/cloudlib/azurecloud"
)

func main() {
	subscription := os.Getenv("Subscription")
	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Panic(err.Error())
	}
	client := armsecurityinsights.NewWatchlistsClient(subscription, creds, &arm.ClientOptions{})
	err = azurecloud.SentinelMaliciousIPWatchlist(client, "core", "sentinel")
	if err != nil {
		log.Panic(err.Error())
	}
}
