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
	exampleRates := CryptoRates{BTC: "1", ETH: "1"}
	assets := computeCryptoHoldings(decimal.NewFromInt(100), exampleRates)
	expectedAssets := CryptoHoldings{BTCHoldings: decimal.NewFromInt(70), ETHHoldings: decimal.NewFromInt(30)}

	if !expectedAssets.BTCHoldings.Equal(assets.BTCHoldings) {
		t.Errorf("Expected BTC => 70, got BTC => %s", assets.BTCHoldings)
	}

	if !expectedAssets.ETHHoldings.Equal(assets.ETHHoldings) {
		t.Errorf("Expected ETH => 30, got ETH => %s", assets.ETHHoldings)
	}
}
