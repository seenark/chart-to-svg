package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const BINANCE_API = "https://api.binance.com"

type KLine struct {
	// OpenTime                float64
	// Open                    float64
	// High                    float64
	// Low                     float64
	Close float64
	// Volume                  float64
	// CloseTime float64
	// QuoteAssetVolume        float64
	// NumberOfTrades          float64
	// TakerBuyAssetVolume     float64
	// TakerBuyQuoteAssetVolum float64
	// Ignore                  float64
}

type BinanceClient struct {
	HttpClient *http.Client
}

func NewBinanceClient(client *http.Client) *BinanceClient {
	return &BinanceClient{
		HttpClient: client,
	}
}

func (b *BinanceClient) GetKLine(symbol string, interval string, startTime int, endTime int, limit int) []KLine {
	uri := "/api/v3/klines"
	url := fmt.Sprintf("%s%s", BINANCE_API, uri)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	queryString := req.URL.Query()
	symbolUpperCase := strings.ToUpper(symbol)
	queryString.Add("symbol", symbolUpperCase)
	queryString.Add("interval", interval)
	if startTime > 0 {

		queryString.Add("startTime", fmt.Sprintf("%d", startTime))
	}
	if endTime > 0 {

		queryString.Add("endTime", fmt.Sprintf("%d", endTime))
	}
	queryString.Add("limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = queryString.Encode()

	fmt.Println(req.URL)

	if err != nil {
		fmt.Println(err)
	}
	resp, err := b.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var bd [][]interface{}
	err = json.Unmarshal(body, &bd)
	if err != nil {
		fmt.Println(err)
	}

	newKlines := NewKLine(bd)
	// fmt.Println("newKlines", newKlines)
	return newKlines

}

func NewKLine(klineSlice [][]interface{}) []KLine {
	// fmt.Println("v", klineSlice)
	klines := []KLine{}

	for _, kl := range klineSlice {
		kline := KLine{}
		for i, k := range kl {
			str := fmt.Sprintf("%v", k)
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				appendKline(i, f, &kline)
			}
		}

		// fmt.Println("kline", kline)
		klines = append(klines, kline)
	}
	return klines
}

func appendKline(indexNumber int, f float64, kl *KLine) *KLine {
	switch indexNumber {
	// case 0:
	// 	kl.OpenTime = f
	// case 1:
	// 	kl.Open = f
	// case 2:
	// 	kl.High = f
	// case 3:
	// 	kl.Low = f
	case 4:
		kl.Close = f
		// case 5:
		// 	kl.Volume = f
		// case 6:
		// 	kl.CloseTime = f
		// case 7:
		// 	kl.QuoteAssetVolume = f
		// case 8:
		// 	kl.NumberOfTrades = f
		// case 9:
		// 	kl.TakerBuyAssetVolume = f
		// case 10:
		// 	kl.TakerBuyQuoteAssetVolum = f
		// case 11:
		// 	kl.Ignore = f

	}

	return kl
}
