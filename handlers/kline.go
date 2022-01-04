package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/seenark/poc-svg/helpers"
	"github.com/seenark/poc-svg/repository"
	"github.com/seenark/poc-svg/routine"
)

type KLineHandler struct {
	Repo repository.ICoinKLineRepository
}

var exitChan = make(chan bool)

func NewKlineHandler(app fiber.Router, klineRepo repository.ICoinKLineRepository) {
	handler := KLineHandler{
		Repo: klineRepo,
	}

	// go func() {
	// 	all, err := handler.getMultiple(coins.CoinList)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	for _, v := range all {
	// 		fmt.Printf("v: %v\n", v.Symbol)
	// 	}
	// }()

	app.Get("/", handler.getMultipleCoinData)
	app.Get("/:symbol", handler.getKline)
	app.Post("/", handler.create)
	app.Put("/:symbol", handler.update)
	app.Delete("/:symbol", handler.delete)
}

func NewKlineSVGHandler(app fiber.Router, klineRepo repository.ICoinKLineRepository) {
	handler := KLineHandler{
		Repo: klineRepo,
	}
	all, _ := klineRepo.GetMultiple([]string{})
	for _, cn := range all {
		routine.SvgCaches[cn.Symbol] = cn
	}
	// fmt.Printf("svgCaches: %v\n", routine.SvgCaches)
	app.Get("/", MSvgCache, handler.getSvg)
	app.Post("/routine", func(c *fiber.Ctx) error {
		running := map[string]bool{
			"running": false,
		}
		err := c.BodyParser(&running)
		if err != nil {
			return err
		}
		fmt.Printf("running: %v\n", running["running"])

		if running["running"] {
			go routine.FetchKlineRoutine(klineRepo, exitChan)
			routine.IsRoutineRunning = running["running"]

		} else if !running["running"] && routine.IsRoutineRunning {
			routine.IsRoutineRunning = running["running"]
			exitChan <- true
		}

		return c.SendString(fmt.Sprintf("fetch kline routine is %v", running["running"]))
	})
}

func (kh KLineHandler) create(c *fiber.Ctx) error {
	kline := repository.CoinKLine{}
	err := c.BodyParser(&kline)
	if err != nil {
		return err
	}
	err = kh.Repo.Create(kline)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusCreated)

}

func (kh KLineHandler) getMultiple(symbols []string) ([]repository.CoinKLine, error) {
	// symbols := ctx.Query("symbols")
	// split := strings.Split(symbols, ",")

	// for index, v := range split {
	// 	if v == "" {
	// 		continue
	// 	}
	// 	split[index] = strings.Trim(v, " ")
	// }

	all, err := kh.Repo.GetMultiple(symbols)
	if err != nil {
		// ctx.SendStatus(http.StatusNotFound)
		return nil, err
	}
	notFoundList := []string{}
	symbolMap := map[string]bool{}
	for _, s := range all {
		symbolMap[s.Symbol] = true
	}
	fmt.Printf("symbolMap: %v\n", symbolMap)
	for _, s := range symbols {
		if _, ok := symbolMap[s]; !ok {
			notFoundList = append(notFoundList, s)
		}
	}

	fmt.Printf("notFoundList: %v\n", notFoundList)
	fmt.Printf("notFoundList: %v\n", len(notFoundList))
	for _, sb := range notFoundList {
		if sb == "" {
			continue
		}
		ck, err := routine.StoreHourKLineForSymbol(sb, kh.Repo)
		if err != nil {
			fmt.Println("some error", err.Error())
			continue
		}
		all = append(all, *ck)
	}
	helpers.PrintMemUsage()
	return all, nil
	// return ctx.Status(http.StatusOK).JSON(all)
}

func (kh KLineHandler) getMultipleCoinData(ctx *fiber.Ctx) error {
	symbols := ctx.Query("symbols")
	split := strings.Split(symbols, ",")

	for index, v := range split {
		if v == "" {
			continue
		}
		split[index] = strings.Trim(v, " ")
	}
	all, err := kh.getMultiple(split)
	if err != nil {
		ctx.SendStatus(http.StatusNotFound)
	}
	return ctx.Status(http.StatusOK).JSON(all)
}

func (kh KLineHandler) getSvg(ctx *fiber.Ctx) error {
	symbols := ctx.Query("symbol")
	split := strings.Split(symbols, ",")

	for index, v := range split {
		if v == "" {
			continue
		}
		split[index] = strings.Trim(v, " ")
	}
	all, err := kh.getMultiple(split)
	if err != nil {
		return ctx.SendStatus(http.StatusNotFound)
	}
	if len(all) <= 0 {
		return ctx.SendStatus(http.StatusNotFound)
	}
	coin := all[0]
	routine.SvgCaches[coin.Symbol] = coin
	return ctx.Type(".svg").Status(http.StatusOK).Send([]byte(coin.Svg))
}

func (kh KLineHandler) getKline(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	kl, err := kh.Repo.GetBySymbol(symbol)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return c.SendStatus(http.StatusNotFound)
		}

		return err
	}
	return c.JSON(kl)
}

func (kh KLineHandler) update(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	kline := repository.CoinKLine{}
	c.BodyParser(&kline)
	err := kh.Repo.Update(symbol, kline)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

func (kh KLineHandler) delete(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	err := kh.Repo.Delete(symbol)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}
