package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
)

// our json looks like this:
// {
//    "data"=> {
//      "currency"=>"USD",
//      "rates"=> {
//        "BTC"=>"0.000015409648952",
//        "ETH"=>"0.0002930428692413"
//      }
// }
//
// It contains the rates for many more currencies, but we can ignore them for this problem.
// We also ignore the currency field because we are always requesting USD.
type CoinbaseResponse struct {
  Data CurrencyInfo
}

// Also returns a Currency field that is always USD, which we can ignore
type CurrencyInfo struct {
  Rates CryptoRates
}

// The API returns rates of many currencies, and Golang is nice enough
// to parse out only the data we define in our structs.
// We only care about BTC and ETH for this problem
type CryptoRates struct {
  BTC string
  ETH string
}

func getRates() CryptoRates {
  // http request
  resp, err := http.Get("https://api.coinbase.com/v2/exchange-rates?currency=USD")
  if err != nil {
    log.Fatalln(err)
  }

  // retrieve body from successful request
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatalln(err)
  }

  // extract and parse data from the JSON byte string
  var responseData CoinbaseResponse
  err = json.Unmarshal([]byte(body), &responseData)
  if err != nil {
    log.Fatalln(err)
  }

  return responseData.Data.Rates
}

func main() {
  // http request!
  rates := getRates()
  if rates.BTC == "" || rates.ETH == "" {
    log.Fatalln("BTC or ETH rates not found")
  }

  spendingMoney := os.Args[1]

  fmt.Println(spendingMoney)

  fmt.Println(rates.BTC)
  fmt.Println(rates.ETH)
}
