package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

func TestGetRates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/exchange-rates" {
			t.Errorf("Expected to request '/v2/exchange-rates', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"rates":{"BTC":"1","ETH":"0.5"}}}`))
	}))
	defer server.Close()

	rates := getRates(server.URL)
	expected := CryptoRates{BTC: "1", ETH: "0.5"}

	if rates != expected {
		t.Errorf("Expected BTC => 1, ETH => 0.5, got BTC => %s, ETH => %s", rates.BTC, rates.ETH)
	}
}

func TestComputeCryptoHoldings(t *testing.T) {
	tests := []struct {
		rates  CryptoRates
		usd    decimal.Decimal
		result CryptoHoldings
	}{
		{
			rates: CryptoRates{BTC: "1", ETH: "1"},
			usd:   decimal.NewFromInt(100),
			result: CryptoHoldings{
				BTCHoldings: decimal.NewFromInt(70),
				ETHHoldings: decimal.NewFromInt(30),
			},
		},
		{
			rates: CryptoRates{BTC: "5", ETH: "3"},
			usd:   decimal.NewFromInt(100),
			result: CryptoHoldings{
				BTCHoldings: decimal.NewFromInt(350),
				ETHHoldings: decimal.NewFromInt(90),
			},
		},
		{
			rates: CryptoRates{BTC: "0.0123", ETH: "0.00045"},
			usd:   decimal.NewFromInt(100),
			result: CryptoHoldings{
				BTCHoldings: decimal.NewFromFloat(0.861),
				ETHHoldings: decimal.NewFromFloat(0.0135),
			},
		},
		{
			rates: CryptoRates{BTC: "0.0123", ETH: "0.00045"},
			usd:   decimal.NewFromFloat(200.58),
			result: CryptoHoldings{
				BTCHoldings: decimal.NewFromFloat(1.727043),
				ETHHoldings: decimal.NewFromFloat(0.0270765),
			},
		},
	}

	for _, test := range tests {
		assets := computeCryptoHoldings(test.usd, test.rates)
		expectedAssets := test.result
		if !expectedAssets.BTCHoldings.Equal(assets.BTCHoldings) {
			t.Errorf(
				"Spending $%s (BTC: %s): Expected BTC => %s, got BTC => %s",
				test.usd,
				test.rates.BTC,
				expectedAssets.BTCHoldings,
				assets.BTCHoldings,
			)
		}

		if !expectedAssets.ETHHoldings.Equal(assets.ETHHoldings) {
			t.Errorf(
				"Spending $%s (ETH: %s): Expected ETH => %s, got ETH => %s",
				test.usd,
				test.rates.ETH,
				expectedAssets.ETHHoldings,
				assets.ETHHoldings,
			)
		}
	}
}
