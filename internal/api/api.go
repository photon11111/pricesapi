package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const coingeckoAPI = "https://api.coingecko.com/api/v3/simple"
const binanceAPI = "https://api.binance.com/api/v3/ticker"

func GetPriceFromCoingecko(baseCurrency string, quotedCurrency string) (string, error) {
	url := fmt.Sprintf("%s/price?ids=%s&vs_currencies=%s", coingeckoAPI, baseCurrency, quotedCurrency)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("getting data from Coingecko's API error: %s", resp.Status)
	}

	result := make(map[string]map[string]string)

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result[baseCurrency][quotedCurrency], nil
}

type binanceResponse struct {
	Symbol string
	Price  string
}

func GetPriceFromBinance(baseCurrency string, quotedCurrency string) (string, error) {
	url := fmt.Sprintf("%s/price?symbol=%s%s", binanceAPI, baseCurrency, quotedCurrency)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("getting data from Binance's API error: %s", resp.Status)
	}

	var priceResp binanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return "", err
	}

	return priceResp.Price, nil
}
