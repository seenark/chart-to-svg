package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/seenark/poc-svg/binance"
	"github.com/seenark/poc-svg/coins"
	mychart "github.com/seenark/poc-svg/myChart"
)

func main() {
	// series := []float64{8, 5, 6, 5, 6, 6, 7, 5, 6, 7, 5, 6, 5, 4, 6, 6, 5, 5, 4, 5, 6, 4, 6, 5}

	// buffer := mychart.GenerateSVG(series)

	// WriteFileSvg("btc.svg", buffer)

	createAllSvg()

}

func createAllSvg() {
	binanceClient := binance.NewBinanceClient(&http.Client{})
	start, end := timeLast24H()
	for _, c := range coins.CoinList {
		btc := binanceClient.GetKLine(c, "1h", start, end, 500)
		series := []float64{}
		for _, v := range btc {
			series = append(series, v.Close)
		}
		buffer := mychart.GenerateSVG(series)

		WriteFileSvg(c, buffer)
	}
}

func WriteFileSvg(name string, buffer *bytes.Buffer) {
	newName := fmt.Sprintf("svg/%s.svg", name)
	err := ioutil.WriteFile(newName, buffer.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func timeLast24H() (start, end int) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -1)
	start = ToMilliseconds(startDate)
	end = ToMilliseconds(now)
	return
}

// helpers
func ToMilliseconds(t time.Time) int {
	return int(t.UnixNano()) / 1e6
}
