package main

import (
	"io/ioutil"

	mychart "github.com/seenark/poc-svg/myChart"
)

func main() {
	series := []float64{8, 5, 6, 5, 6, 6, 7, 5, 6, 7, 5, 6, 5, 4, 6, 6, 5, 5, 4, 5, 6, 4, 6, 5}

	buffer := mychart.GenerateSVG(series)

	// fmt.Println(buffer.String())
	ioutil.WriteFile("test.svg", buffer.Bytes(), 0644)

}
