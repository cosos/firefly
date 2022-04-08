package azurecloud

import (
	"context"
	"log"

	"go.uber.org/zap"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights"
)

func CheckIndicatorStatus(client *armsecurityinsights.ThreatIntelligenceIndicatorClient, ctx context.Context, resourcegroupName, indicatorName, workspaceName, paternType, patern, source string) bool {
	pager := client.QueryIndicators(resourcegroupName, workspaceName, armsecurityinsights.ThreatIntelligenceFilteringCriteria{
		IncludeDisabled: to.BoolPtr(false),
		PatternTypes:    []*string{to.StringPtr(paternType)},
		Keywords:        []*string{to.StringPtr(patern)},
		Sources:         []*string{to.StringPtr(source)},
	}, &armsecurityinsights.ThreatIntelligenceIndicatorClientQueryIndicatorsOptions{})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			zap.S().Error(err.Error())
		}
		if !nextResult {
			break
		}
		if len(pager.PageResponse().Value) == 0 {
			return false
		}
		for _, item := range pager.PageResponse().Value {
			indicator := item.GetThreatIntelligenceInformation()
			zap.S().Infof("%s has been found", *indicator.Name)
		}
	}
	return true
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
