## Prompt

This endpoint provides up-to-the-minute crypto exchange rates relative to US dollars:
```
https://api.coinbase.com/v2/exchange-rates?currency=USD
```

Crypto is all the rage these days and I don't want to miss out! I want to keep 70% of my crypto holdings in BTC and 30% in ETH. Write a function that takes the amount I have to spend in USD as a parameter and returns the number of Bitcoin and Ethereum to buy.

I have $X I want to keep in BTC and ETH, 70/30 split. How many of each should I buy?

I'd say you should take the input amount as a command line parameter and print a json response with the allocations back out of your program.

## Design decisions

### Code
All the code lives in `main.go` and `main_test.go`. I think the solution is simple enough that this makes sense and does not require more refactoring.

In addition to the `main` function that controls the input logic, the two main functions are 
- `getRates` which retrieves the BTC and ETH rates from the Coinbase API and parses the response. I constructed types to collect only the data we need (BTC and ETH rates), while ignoring data that we don't (eg. the currency field, which is always "USD" or the rates of other currencies)
- `computeCryptoHoldings` which takes the spending money in USD and the parsed rates of BTC and ETH and computes how to allocate our purchases according to the expected 70/30 split.

I used the [popular `Decimal` library](https://github.com/shopspring/decimal) to avoid floating point errors.

### Tests
The test for `getRates` creates a test http server and spoofs the data it returns to confirm that we parse the response correctly.

The test for `getCryptoHoldings` is simpler, since it only requires us to pass in the rates we which to test. TODO more

## How to run

### Install dependencies

```
go mod tidy
```

### Run program

```
go run . 20.33 # put in whatever dollar amount you want
```

### Run tests

```
go test
```
