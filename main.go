package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fiatjaf/go-lnurl"
)

type LNServicePayResponse struct {
	Callback    string `json:"callback"`
	MaxSendable int    `json:"maxSendable"`
	MinSendable int    `json:"minSendable"`
	Metadata    string `json:"metadata"`
	Tag         string `json:"tag"`
}

type LNServiceInvoice struct {
	Pr     string     `json:"pr"`
	Routes []struct{} `json:"routes"`
}

func main() {

	lnstr := "lnurl1dp68gurn8ghj7ampd3kx2ar0veekzar0wd5xjtnrdakj7tnhv4kxctttdehhwm30d3h82unvwqhhw6tndpn82mrxd3jhx6phxs38f45a"
	// Decode LN Url
	decodedLnUrl, _ := lnurl.LNURLDecode(lnstr)
	fmt.Println(decodedLnUrl)
	// Get LN Service URL
	resp, err := http.Get(decodedLnUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// extract callback URL
	var response LNServicePayResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Make an HTTP GET request to callback URL with amount to be paid
	// LN Service will create and return a lnd invoice to be paid with this amount
	callbackUrl := response.Callback + "?amount=3333"
	resp, err = http.Get(callbackUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
	var resInvoice LNServiceInvoice
	err = json.Unmarshal(body, &resInvoice)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Invoice: ",resInvoice.Pr)
	// Decode the ln invoice and check amount to verify

	//"message": "invoice not for current active network 'mainnet'"
	// To decode the invoice, we need pld to be connected to the same network.
	// Then pay the invoice

	
}
