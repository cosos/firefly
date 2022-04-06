package threatfox

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const threatfox_url = "https://threatfox-api.abuse.ch/api/v1/"

//{
//    "id": "492653",
//    "ioc": "f07764455422ee25fda00ae3714a05f55a9f342cc652c1ac6938b4fed350c642",
//    "ioc_type": "sha256_hash",
//    "threat_type": "payload",
//    "malware": "win.emotet",
//    "malware_printable": "Emotet",
//    "malware_alias": "Geodo,Heodo",
//    "confidence_level": "75",
//    "first_seen": "2022-04-05 13:51:11 UTC",
//    "last_seen": "2022-04-05 13:54:34 UTC",
//    "reporter": "Cryptolaemus1",
//    "reference": null,
//    "threatfox_link": "https:\/\/threatfox\/ioc\/492653",
//    "tags": [
//        "epoch5",
//        "exe"
//    ]
//}

// IOC define the ioc data
type IOC struct {
	Id               *string   `json:"id"`
	IoC              *string   `json:"ioc"`
	IoCType          *string   `json:"ioc_type"`
	ThreatType       *string   `json:"threat_type"`
	Malware          *string   `json:"malware"`
	MalwarePrintable *string   `json:"malware_printable"`
	MalwareAlias     *string   `json:"malware_alias"`
	ConfidenceLevel  *string   `json:"confidence_level"`
	FirstSeen        time.Time `json:"first_seen"`
	LastSeen         time.Time `json:"last_seen"`
	Reporter         *string   `json:"reporter"`
	Reference        *string   `json:"reference"`
	ThreatFoxLink    *string   `json:"threatfox_link"`
	Tags             []string  `json:"tags"`
}

// ThreatFoxIOCSet type
type ThreatFoxIOCSet struct {
	QueryStatus *string `json:"query_status"`
	Data        []IOC   `json:"data"`
}

type GetIOCRequestInput struct {
	Query string `json:"query"`
	Days  int    `json:"days"`
}

func GetThreatFoxIoCSet(client http.Client, days int) (*ThreatFoxIOCSet, error) {
	reqBody := GetIOCRequestInput{Query: "get_iocs", Days: days}
	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(reqBody)
	req, err := http.NewRequest(http.MethodPost, threatfox_url, data)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var result ThreatFoxIOCSet
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
