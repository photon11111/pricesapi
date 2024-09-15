package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const binanceAPI = "https://api.binance.com/api/v3/ticker"
const priceEndpoint = "/price?symbol="

type binanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func GetPriceFromBinance(baseCurrency string, quotedCurrency string) (float64, error) {
	url := fmt.Sprintf("%s%s%s%s", binanceAPI, priceEndpoint, baseCurrency, quotedCurrency)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error making request to Binance API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Binance API returned non-200 status: %s", resp.Status)
	}

	var priceResp binanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return 0, fmt.Errorf("error decoding response from Binance API: %w", err)
	}

	price, err := strconv.ParseFloat(priceResp.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting price to float64: %w", err)
	}

	return price, nil
}
