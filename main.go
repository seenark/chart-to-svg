package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/seenark/poc-svg/binance"
	"github.com/seenark/poc-svg/coins"
	"github.com/seenark/poc-svg/config"
	"github.com/seenark/poc-svg/handlers"
	mychart "github.com/seenark/poc-svg/myChart"
	"github.com/seenark/poc-svg/repository"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/svg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	initTimeZone()
	// series := []float64{8, 5, 6, 5, 6, 6, 7, 5, 6, 7, 5, 6, 5, 4, 6, 6, 5, 5, 4, 5, 6, 4, 6, 5}

	// buffer := mychart.GenerateSVG(series)

	// WriteFileSvg("btc.svg", buffer)
	// createAllSvg()

	cf := config.GetConfig()

	ctx := context.TODO()

	mongoClient := connectMongo(cf.Mongo.Username, cf.Mongo.Password)

	klineDb := mongoClient.Database(cf.Mongo.KlineDbName)
	hourCollection := klineDb.Collection(cf.Mongo.HourKlineCollection)
	makeSymbolAsIndexes(hourCollection)

	klineRepository := repository.NewKLineRepository(hourCollection, ctx)
	app := fiber.New()
	app.Use(cors.New())
	klineGroup := app.Group("/kline")
	svgGroup := app.Group("/svg")
	handlers.NewKlineHandler(klineGroup, klineRepository)
	handlers.NewKlineSVGHandler(svgGroup, klineRepository)

	app.Listen(fmt.Sprintf(":%d", cf.Port))

}
func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}
func connectMongo(username string, password string) *mongo.Client {
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@hdgcluster.xmgsx.mongodb.net/myFirstDatabase?retryWrites=true&w=majority", username, password))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
func makeSymbolAsIndexes(collection *mongo.Collection) {
	indexName, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "symbol", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("index name:", indexName)
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
	b, err := minifySVG(buffer.Bytes())
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(newName, b, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func minifySVG(b []byte) ([]byte, error) {
	m := minify.New()
	m.AddFunc("image/svg+xml", svg.Minify)
	newBytes, err := m.Bytes("image/svg+xml", b)
	if err != nil {
		return nil, err
	}
	return newBytes, nil
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
