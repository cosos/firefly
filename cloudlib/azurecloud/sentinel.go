package azurecloud

import (
	"bytes"
	"context"
	"encoding/csv"
	"log"
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights"
)

const IntelUrl = "https://urlhaus.abuse.ch/downloads/csv_online/"
const DefaultProvider = "Microsoft"

func GetMaliciousCSV() ([]byte, error) {
	maliciousIPs := [][]string{{"IP", "Threat", "Tags"}}
	resp, err := http.DefaultClient.Get(IntelUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	csvReader.Comma = ','
	csvReader.Comment = '#'
	allRecords, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, item := range allRecords {
		urldata, err := url.Parse(item[2])
		if err != nil {
			log.Println(err.Error())
			continue
		}
		maliciousIPs = append(maliciousIPs, []string{urldata.Host, item[4], item[5]})
	}
	buf := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buf)
	err = csvWriter.WriteAll(maliciousIPs)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func CreateOrUpdateWatchListFromFile(client *armsecurityinsights.WatchlistsClient, resourcegroup, workspace, alias, name, contenttype string) error {
	resp, err := client.CreateOrUpdate(
		context.TODO(),
		resourcegroup,
		workspace, alias,
		armsecurityinsights.Watchlist{
			Properties: &armsecurityinsights.WatchlistProperties{
				DisplayName: &name,
				ContentType: &contenttype,
				Provider:    ,
				Source:      armsecurityinsights.Source("Local file").ToPtr(),
			},
		},
		&armsecurityinsights.WatchlistsClientCreateOrUpdateOptions{},
	)
}
