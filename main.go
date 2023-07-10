package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fiatjaf/go-lnurl"
	"github.com/tidwall/gjson"
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

func decode(url string, amount int) (invoice string) {
	fmt.Println("Decoding LNURL: ", url)
	//lnstr := "lnurl1dp68gurn8ghj7ampd3kx2ar0veekzar0wd5xjtnrdakj7tnhv4kxctttdehhwm30d3h82unvwqhhw6tndpn82mrxd3jhx6phxs38f45a"
	// Decode LN Url
	decodedLnUrl, _ := lnurl.LNURLDecode(url)
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
	callbackUrl := response.Callback + "?amount=" + strconv.Itoa(amount)
	resp, err = http.Get(callbackUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Println(string(body))
	var resInvoice LNServiceInvoice
	err = json.Unmarshal(body, &resInvoice)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println("Invoice: ", resInvoice.Pr)
	// Decode the ln invoice and check amount to verify
	data := map[string]string{"pay_req": resInvoice.Pr}
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	decodeInvoiceUtl := "http://localhost:8080/api/v1/lightning/invoice/decodepayreq"
	req, err := http.NewRequest("POST", decodeInvoiceUtl, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	parsed := gjson.ParseBytes(b)
	if resp.StatusCode != 200 {
		fmt.Println("Error: ", parsed.Get("error").String())
	}

	//fmt.Println("response Body:", parsed)
	return parsed.Get("paymentHash").String()
}

func payInvoice(invoice string) {
	//TODO: Pay invoice
}

func generate(url string) {
	fmt.Println("Encoding into LNURL: ", url)
	encodedUrl, err := lnurl.Encode(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("LNURL: ", encodedUrl)

}

func main() {
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "lnurl=") {
		lnurl := strings.TrimPrefix(os.Args[1], "lnurl=")
		amount, err := strconv.Atoi(strings.TrimPrefix(os.Args[2], "amount="))
		if err != nil {
			fmt.Println("Invalid amount")
			return
		}
		invoice := decode(lnurl, amount)
		if invoice != "" {
			payInvoice(invoice)
		}
	} else if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "url=") {
		url := strings.TrimPrefix(os.Args[1], "url=")
		generate(url)
	}
}
