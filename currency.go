package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
)

var ErrInvalidRate error = errors.New("incorrect exchange rate")

type LatestRates struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

// "currencyConv" convert a user_balance in other currency,
// using external API. If user entered the wrong currency name
// this func return the error.
func currencyConv(balance float64, currency string) (float64, error) {
	url := "https://api.exchangeratesapi.io/latest?base=RUB"
	resp, err := http.Get(url)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0.0, err
	}

	var lr *LatestRates
	json.Unmarshal(body, &lr)

	if rate, ok := lr.Rates[currency]; ok {
		balance = balance * rate
		balance = math.Round(balance*100) / 100
		return balance, nil
	}
	return 0.0, ErrInvalidRate
}
