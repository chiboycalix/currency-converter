package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/joho/godotenv"
)

func API_KEY() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("API_KEY")
}

type ExchangeRate struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

type CurrencyData struct {
	GBP ExchangeRate `json:"GBP"`
	EUR ExchangeRate `json:"EUR"`
	NGN ExchangeRate `json:"NGN"`
	USD ExchangeRate `json:"USD"`
}

type MetaData struct {
	LastUpdatedAt time.Time `json:"last_updated_at"`
}
type CurrencyExchange struct {
	Meta MetaData     `json:"meta"`
	Data CurrencyData `json:"data"`
}

func main() {
	url := "https://api.currencyapi.com/v3/latest?currencies=EUR%2CGBP%2CNGN%2CUSD"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("apikey", API_KEY())

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	apiResponse := string(body)

	var exchangeData CurrencyExchange
	error := json.Unmarshal([]byte(apiResponse), &exchangeData)
	if error != nil {
		fmt.Println("Error decoding JSON:", error)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Last_Updated_At", "USD", "NGN", "EUR", "GBP"})
	t.AppendRows([]table.Row{
		{1, exchangeData.Meta.LastUpdatedAt, exchangeData.Data.USD.Value, exchangeData.Data.NGN.Value, exchangeData.Data.EUR.Value, exchangeData.Data.GBP.Value},
	})
	t.Render()
}
