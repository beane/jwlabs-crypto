package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "math"
  "os"
  "strconv"
)

const BTCPercentage = 70
const ETHPercentage = 30

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

type CryptoHoldings struct {
  BTCHoldings float64
  ETHHoldings float64
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
  // http request
  rates := getRates()

  // parse rates from strings
  if rates.BTC == "" || rates.ETH == "" {
    log.Fatalln("BTC or ETH rates not found")
  }

  btcRate, err:= strconv.ParseFloat(rates.BTC, 64)
  if err != nil {
    log.Fatalln(err)
  }

  ethRate, err:= strconv.ParseFloat(rates.ETH, 64)
  if err != nil {
    log.Fatalln(err)
  }

  // parse spending money from first argument
  spendingMoneyInput := os.Args[1]
  spendingMoney, err := strconv.ParseFloat(spendingMoneyInput, 64)
  if err != nil {
    log.Fatalln(err)
  }

  // compute amount of USD to spend on each kind of asset
  // ideally we would convert to ints to avoid floating point errors (multiplying by 10^18)
  // but it seems like the numbers we're dealing with would be too big to avoid errors
  // according to this article:
  // https://betterprogramming.pub/compute-without-floating-point-errors-7b92695bde4 
  btcSpendingMoney := BTCPercentage * spendingMoney / 100
  ethSpendingMoney := spendingMoney - btcSpendingMoney

  // Round to 2 digits of precision
  btcSpendingMoney = math.Round(btcSpendingMoney * 100)/100
  ethSpendingMoney = math.Round(ethSpendingMoney * 100)/100

  // spending money usd * (btc / usd) = btc
  btcPurchases := btcSpendingMoney * btcRate
  ethPurchases := ethSpendingMoney * ethRate

  totalHoldings := CryptoHoldings{BTCHoldings: btcPurchases, ETHHoldings: ethPurchases}
  resultJson, err := json.Marshal(totalHoldings)

  // TODO: write tests
  fmt.Println(string(resultJson[:]))
}
