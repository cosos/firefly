package azurecloud

import (
	"bytes"
	"context"
	"encoding/csv"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/cosos/firefly/cloudlib"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights"
)

const IntelUrl = "https://urlhaus.abuse.ch/downloads/csv_online/"

func SentinelMaliciousIPWatchlist(client *armsecurityinsights.WatchlistsClient, resourcegroup, workspace string) error {
	maliciousIPs := [][]string{{"IP", "Threat", "Tags"}}
	maliciousDomains := [][]string{{"Domain", "Threat", "Tags"}}
	validID := regexp.MustCompile(`[a-zA-Z]+`)
	resp, err := http.DefaultClient.Get(IntelUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	csvReader.Comma = ','
	csvReader.Comment = '#'
	allRecords, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for _, item := range allRecords {
		urldata, err := url.Parse(item[2])
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if validID.MatchString(urldata.Host) {
			maliciousDomains = append(maliciousDomains, []string{urldata.Host, item[4], item[5]})
			continue
		}
		maliciousIPs = append(maliciousIPs, []string{urldata.Hostname(), item[4], item[5]})
	}
	if len(maliciousIPs) > 1 {
		buf := new(bytes.Buffer)
		csvWriter := csv.NewWriter(buf)
		err = csvWriter.WriteAll(maliciousIPs)
		if err != nil {
			return err
		}
		_, err := client.CreateOrUpdate(
			context.TODO(),
			resourcegroup,
			workspace,
			"urlhaus-ip",
			armsecurityinsights.Watchlist{
				Properties: &armsecurityinsights.WatchlistProperties{
					DisplayName:    cloudlib.String("MaliciousIPs-urlhaus"),
					ContentType:    cloudlib.String("text/csv"),
					Provider:       cloudlib.String("Microsoft"),
					ItemsSearchKey: cloudlib.String("IP"),
					Source:         armsecurityinsights.Source("Local file").ToPtr(),
					RawContent:     cloudlib.String(buf.String()),
				},
			},
			&armsecurityinsights.WatchlistsClientCreateOrUpdateOptions{},
		)
		if err != nil {
			return err
		}
	}
	if len(maliciousDomains) > 1 {
		domainbuf := new(bytes.Buffer)
		csvWriter := csv.NewWriter(domainbuf)
		err = csvWriter.WriteAll(maliciousDomains)
		if err != nil {
			return err
		}
		_, err := client.CreateOrUpdate(
			context.TODO(),
			resourcegroup,
			workspace,
			"urlhaus-domain",
			armsecurityinsights.Watchlist{
				Properties: &armsecurityinsights.WatchlistProperties{
					DisplayName:    cloudlib.String("MaliciousDomains-urlhaus"),
					ContentType:    cloudlib.String("text/csv"),
					Provider:       cloudlib.String("Microsoft"),
					ItemsSearchKey: cloudlib.String("Domain"),
					Source:         armsecurityinsights.Source("Local file").ToPtr(),
					RawContent:     cloudlib.String(domainbuf.String()),
				},
			},
			&armsecurityinsights.WatchlistsClientCreateOrUpdateOptions{},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
