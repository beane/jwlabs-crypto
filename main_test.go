package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
		t.Errorf("Expected 'BTC => 1, ETH => 0.5', got BTC => %s, ETH => %s", rates.BTC, rates.ETH)
	}
}

func TestGetCryptoHoldings(t *testing.T) {
	exampleRates := CryptoRates{BTC: "1", ETH: "1"}
	rates := computeCryptoHoldings(10000, exampleRates)
	expectedRates := CryptoHoldings{BTCHoldings: 70, ETHHoldings: 30}

	if rates != expectedRates {
		t.Errorf("blurp")
	}
}
