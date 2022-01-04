package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/seenark/poc-svg/helpers"
	"github.com/seenark/poc-svg/repository"
	"github.com/seenark/poc-svg/routine"
)

func MSvgCache(c *fiber.Ctx) error {
	symbol := c.Query("symbol")
	found := false
	coin := repository.CoinKLine{}
	for _, cn := range routine.SvgCaches {
		if cn.Symbol == symbol {
			found = true
			coin = cn
		}
	}
	if found {
		fmt.Println("found in caches and return it")
		helpers.PrintMemUsage()
		return c.Type(".svg").Status(http.StatusOK).Send([]byte(coin.Svg))
	}
	return c.Next()
}
