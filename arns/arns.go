package arns

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type ArNS struct {
	DreUrl      string
	ArNSAddress string
	HttpClient  *http.Client
}

func NewArNS(dreUrl string, arNSAddr string, timout time.Duration) *ArNS {

	// default timeout is 5s
	if timout == 0 {
		timout = 5 * time.Second
	}

	httpClient := &http.Client{
		Timeout: timout, // Set the timeout for HTTP requests
	}
	return &ArNS{
		DreUrl:      dreUrl,
		ArNSAddress: arNSAddr,
		HttpClient:  httpClient,
	}
}

func (a *ArNS) QueryLatestRecord(domain string) (txId string, err error) {
	spliteDomains := strings.Split(domain, "_")
	if len(spliteDomains) > 2 { // todo now only support level-2 subdomain
		return "", errors.New("current arseeding gw not support over level-2 subdomain")
	}
	rootDomain := spliteDomains[len(spliteDomains)-1]
	// step1 query NameCA address
	caAddress, err := a.QueryNameCa(rootDomain)
	if err != nil {
		return "", err
	}
	// step2 query latest txId
	// Currently, only level-1 domain name resolution is queried
	subdomain := spliteDomains[0]
	if subdomain == rootDomain {
		subdomain = "@"
	}
	txId, err = a.GetArNSTxID(caAddress, subdomain)
	return
}

func (a *ArNS) QueryNameCa(domain string) (caAddress string, err error) {
	baseURL := a.DreUrl + "/contract/"

	// Construct the complete URL
	url := baseURL + "?id=" + a.ArNSAddress

	// Make the HTTP request using the custom HTTP client
	response, err := a.HttpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", response.Status)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	value := gjson.Get(string(body), "state.records."+domain+".contractTxId")

	if !value.Exists() {
		return "", fmt.Errorf("domain %s not exist", domain)
	}

	return value.String(), nil

}

func (a *ArNS) GetArNSTxID(caAddress string, domain string) (txId string, err error) {

	baseURL := a.DreUrl + "/contract/"

	// Construct the complete URL
	url := baseURL + "?id=" + caAddress

	// Make the HTTP request using the custom HTTP client
	response, err := a.HttpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GetArNSTxID: unexpected response status: %s", response.Status)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	value := gjson.Get(string(body), "state.records."+domain+".transactionId")

	/** The ArNS interface will return two types of data for compatibility processing
	data type 1: https://dre-1.warp.cc/contract?id=jr4P6Y_Olv3QGho0uo7p9DpvSn33mUC_XgJSKB3JDZ4
	records:{
	@:{
	transactionId:"wQk7txuMvlrlYlVozj6aeF7E9dlwar8nNtfs3iNTpbQ"
	ttlSeconds:900
	}
	}
	data type 2: https://dre-3.warp.cc/contract?id=Vx4bW_bh7nXMyq-Jy24s9EiCyY_BXZuToshhSqabc9o
	records:{
	@:"wQk7txuMvlrlYlVozj6aeF7E9dlwar8nNtfs3iNTpbQ"
	}
	*/
	if !value.Exists() {
		value = gjson.Get(string(body), "state.records."+domain)
	}

	if !value.Exists() {
		return "", fmt.Errorf("GetArNSTxID: domain %s not exist", domain)
	}

	return value.String(), nil

}
