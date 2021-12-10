package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/leekchan/accounting"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Quote struct {
	Name                     string  `json:"name"`
	Price                    float64 `json:"price"`
	Volume24h                float64 `json:"volume24h"`
	Volume7d                 float64 `json:"volume7d"`
	Volume30d                float64 `json:"volume30d"`
	MarketCap                float64 `json:"marketCap"`
	SelfReportedMarketCap    float64 `json:"selfReportedMarketCap"`
	PercentChange1h          float64 `json:"percentChange1h"`
	PercentChange24h         float64 `json:"percentChange24h"`
	PercentChange7d          float64 `json:"percentChange7d"`
	LastUpdated              string  `json:"lastUpdated"`
	PercentChange30d         float64 `json:"percentChange30d"`
	PercentChange60d         float64 `json:"percentChange60d"`
	PercentChange90d         float64 `json:"percentChange90d"`
	FullyDilluttedMarketCap  float64 `json:"fullyDilluttedMarketCap"`
	MarketCapByTotalSupply   float64 `json:"marketCapByTotalSupply"`
	Dominance                float64 `json:"dominance"`
	Turnover                 float64 `json:"turnover"`
	YtdPriceChangePercentage float64 `json:"ytdPriceChangePercentage"`
}

type CryptoItem struct {
	ID                            int     `json:"id"`
	Name                          string  `json:"name"`
	Symbol                        string  `json:"symbol"`
	Slug                          string  `json:"slug"`
	CMCRank                       int     `json:"cmcRank"`
	MarketPairCount               int     `json:"marketPairCount"`
	CirculatingSupply             float64 `json:"circulatingSupply"`
	SelfReportedCirculatingSupply float64 `json:"selfReportedCirculatingSupply"`
	TotalSupply                   float64 `json:"totalSupply"`
	ATH                           float64 `json:"ath"`
	ATL                           float64 `json:"atl"`
	High24h                       float64 `json:"high24h"`
	Low24h                        float64 `json:"low24h"`

	LastUpdated string `json:"lastUpdated"`
	DateAdded   string `json:"dateAdded"`

	Quote []Quote `json:"quotes"`
}

type CoinMarketCapResponse struct {
	Data struct {
		CryptoCurrencyList []CryptoItem `json:"cryptoCurrencyList"`
		TotalCount         int          `json:"totalCount,string"`
	} `json:"data"`

	Status struct {
		Timestamp    string `json:"timestamp"`
		ErrorCode    int    `json:"error_code,string"`
		ErrorMessage string `json:"error_message"`
		Elapsed      int    `json:"elapsed,string"`
		CreditCount  int    `json:"credit_count"`
	} `json:"status"`
}

func main() {

	endpoint := "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/listing?start=1&limit=1000&sortBy=market_cap&sortType=desc&convert=USD,BTC,ETH&cryptoType=all&tagType=all&audited=false&aux=ath,atl,high24h,low24h,num_market_pairs,cmc_rank,date_added,max_supply,circulating_supply,total_supply,volume_7d,volume_30d,self_reported_circulating_supply,self_reported_market_cap"

	resp, err := http.Get(endpoint)

	if err != nil {
		log.Fatal("Error fetching endpoint.")
	}

	defer resp.Body.Close()
	fmt.Println("Fetching API response...")
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Failed to read response.")
	}

	var result CoinMarketCapResponse
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	p := message.NewPrinter(language.English)

	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err.Error())
	}

	var dataFormatted [][]string

	dataFormatted = append(dataFormatted, []string{
		"Name",
		"Symbol",
		"Price",
		"24H%",
		"7D%",
		"Market Cap",
		"Volume (24H)",
		"Circulating Supply",
	})

	for _, crypto := range result.Data.CryptoCurrencyList {
		quote := Quote{}
		for _, v := range crypto.Quote {
			if v.Name == "USD" {
				quote = v
			}
		}

		if quote != (Quote{}) {
			dataFormatted = append(dataFormatted, []string{
				crypto.Name,
				crypto.Symbol,
				fmt.Sprintf("$%.2f", quote.Price),
				fmt.Sprintf("%.2f%%", quote.PercentChange24h),
				fmt.Sprintf("%.2f%%", quote.PercentChange7d),
				ac.FormatMoney(quote.MarketCap),
				ac.FormatMoney(quote.Volume24h),
				p.Sprintf("%f", crypto.CirculatingSupply),
			})
		}

	}

	fmt.Println("Saving to export.csv...")
	// Create file
	file, err := os.Create("./export.csv")
	if err != nil {
		log.Fatal("Failed to create report file.")
	}

	defer file.Close()
	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(dataFormatted)
	fmt.Println("Press the enter key to exit.")
	fmt.Scanln()
}
