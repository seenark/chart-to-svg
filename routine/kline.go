package routine

import (
	"fmt"
	"net/http"
	"time"

	"github.com/seenark/poc-svg/binance"
	"github.com/seenark/poc-svg/helpers"
	mychart "github.com/seenark/poc-svg/myChart"
	"github.com/seenark/poc-svg/repository"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/svg"
)

var IsRoutineRunning = false
var SvgCaches = map[string]repository.CoinKLine{}

func FetchKlineRoutine(kRepo repository.ICoinKLineRepository, exitChan chan bool) {
	for {
		select {
		case <-exitChan:
			return
		default:
			now := time.Now()
			// mil := now.Second()
			// if mil%10 == 0 {
			min := now.Minute()
			// fmt.Printf("min: %v\n", min)
			if min == 0 {
				all, err := kRepo.GetMultiple([]string{})
				if err != nil {
					fmt.Println(err)
				}
				for _, sb := range all {
					err = updateKLineForSymbol(sb.Symbol, kRepo)
					if err != nil {
						fmt.Println("Error in Fetch routine", err)
						continue
					}
				}
				helpers.PrintMemUsage()
				// time.Sleep(55 * time.Second)
				time.Sleep(57 * time.Minute)
			}
		}

	}
}

func StoreHourKLineForSymbol(symbol string, klineCollection repository.ICoinKLineRepository) (*repository.CoinKLine, error) {
	now := time.Now()
	min := int64(now.Minute())
	minTime := time.Duration(min)
	end := now.Add(-minTime * time.Minute)
	start := end.AddDate(0, 0, -1)
	bClient := new(binance.BinanceClient)
	bClient.HttpClient = &http.Client{}
	klines := bClient.GetKLine(symbol, "1h", ToMilliseconds(start), ToMilliseconds(end), 24)
	if len(klines) == 0 {
		return nil, fmt.Errorf("error not found symbol on binance server")
	}
	closePrices := []float64{}
	for _, k := range klines {
		closePrices = append(closePrices, k.Close)
	}
	coinKl := repository.CoinKLine{
		Symbol:      symbol,
		ClosePrices: closePrices,
	}
	genSvg(&coinKl)
	err := klineCollection.Create(coinKl)
	if err != nil {
		return nil, err
	}
	SvgCaches[coinKl.Symbol] = coinKl
	return &coinKl, nil
}

func updateKLineForSymbol(symbol string, klineCollection repository.ICoinKLineRepository) error {
	now := time.Now()
	min := int64(now.Minute())
	minTime := time.Duration(min)
	end := now.Add(-minTime * time.Minute)
	start := end.AddDate(0, 0, -1)
	bClient := new(binance.BinanceClient)
	bClient.HttpClient = &http.Client{}
	klines := bClient.GetKLine(symbol, "1h", ToMilliseconds(start), ToMilliseconds(end), 24)
	closePrices := []float64{}
	for _, k := range klines {
		closePrices = append(closePrices, k.Close)
	}
	coinKl := repository.CoinKLine{
		Symbol:      symbol,
		ClosePrices: closePrices,
	}
	genSvg(&coinKl)
	err := klineCollection.Update(symbol, coinKl)
	if err != nil {
		return err
	}
	SvgCaches[coinKl.Symbol] = coinKl
	return nil
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

// helpers
func minifySVG(b []byte) ([]byte, error) {
	m := minify.New()
	m.AddFunc("image/svg+xml", svg.Minify)
	newBytes, err := m.Bytes("image/svg+xml", b)
	if err != nil {
		return nil, err
	}
	return newBytes, nil
}

func genSvg(ck *repository.CoinKLine) {
	buffer := mychart.GenerateSVG(ck.ClosePrices)
	b, err := minifySVG(buffer.Bytes())
	if err != nil {
		fmt.Println(err)
	}
	ck.Svg = string(b)
}
