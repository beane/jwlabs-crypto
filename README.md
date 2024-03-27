## Prompt

This endpoint provides up-to-the-minute crypto exchange rates relative to US dollars:
```
https://api.coinbase.com/v2/exchange-rates?currency=USD
```

Crypto is all the rage these days and I don't want to miss out! I want to keep 70% of my crypto holdings in BTC and 30% in ETH. Write a function that takes the amount I have to spend in USD as a parameter and returns the number of Bitcoin and Ethereum to buy.

I have $X I want to keep in BTC and ETH, 70/30 split. How many of each should I buy?

I'd say you should take the input amount as a command line parameter and print a json response with the allocations back out of your program.

## Design decisions

### Overview

The solution needs to
  - fetch and parse data from the Coinbase API
  - take an amount in USD, split it 70/30
  - spend the 70% on BTC and 30% on ETH

Conceptually this is pretty straightforward. I think the main issue to worry about is floating point errors compounding when multiplying by highly precise crypto rates. Because this problem deals with money we must be much more precise than float types permit. Luckily, there is an [excellent `Decimal` library](https://github.com/shopspring/decimal) that allows precise operations with decimal numbers.

The other potential concern is when we are confronted with spending a fraction of a cent. For example, 70% of $200.58 is $140.406, which we cannot accurately spend. I chose to simply round to 2 digits of precision to ensure that we spend only sensible amounts of money. In this case, we treat 70% of $200.58 as $140.41 and 30% as $60.17.

In the case where the user passes in an number with too much decimal precision (fractional cents), we report an error explaining correct usage.

### Code
All the code lives in `main.go` and `main_test.go`. I think the solution is simple enough that this makes sense and does not require more refactoring.

In addition to the `main` function that controls the input logic, the two main functions are 
- `getRates` which retrieves the BTC and ETH rates from the Coinbase API and parses the response. I constructed types to collect only the data we need (BTC and ETH rates), while ignoring data that we don't (eg. the currency field, which is always "USD" or the rates of other currencies)
- `computeCryptoHoldings` which takes the spending money in USD and the parsed rates of BTC and ETH and computes how to allocate our purchases according to the expected 70/30 split.

### Tests
The test for `getRates` creates a test http server and spoofs the data it returns to confirm that we parse the response correctly.

The tests for `getCryptoHoldings` are simpler, since they only require us to pass in the rates we want test and the dollar amounts we intend to spend. In this test I iterate over a few different cases to check the system's ability to correctly divide up the original USD allocation as well as its capacity to deal with decimal money amounts and correctly compute BTC and ETH holdings from decimal rates.

## How to run

### Install dependencies

```
go mod tidy
```

### Run program

```
go run . 20.33 # put in whatever USD amount you want
```

### Run tests

```
go test
```

### A small sorry!

I wrote the code from two different computers. When I originally set up my account on my desktop I thought it would be fun to put my email as "tachyons". So it looks like there are commits from two different people, but it's just me!
