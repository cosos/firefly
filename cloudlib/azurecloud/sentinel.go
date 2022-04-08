package azurecloud

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights"
)

func ListIndicators(client *armsecurityinsights.ThreatIntelligenceIndicatorsClient, ctx context.Context, resourcegroupName, workspaceName string) {
	pager := client.List(resourcegroupName, workspaceName, &armsecurityinsights.ThreatIntelligenceIndicatorsClientListOptions{})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Printf("failed to get next page: %s", err.Error())
			break
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			result := v.GetThreatIntelligenceInformation()
			var temp []byte
			err := result.UnmarshalJSON(temp)
			if err != nil {
				log.Println("UnmarshalJosn Failed: ", err.Error())
				continue
			}
			log.Panicln(string(temp))
		}
	}
}

func CheckIndicator(client *armsecurityinsights.ThreatIntelligenceIndicatorClient, ctx context.Context, resourcegroupName, workspaceName, indicatorName string) bool {
	resp, err := client.Get(ctx, resourcegroupName, workspaceName, indicatorName, nil)
	if err != nil {
		log.Println(err.Error())
	}
	result := resp.ThreatIntelligenceIndicatorClientGetResult.GetThreatIntelligenceInformation()
	log.Println(*result.Name)
	return false
}
