package main

import (
	"fmt"
	http "net/http"

	"github.com/tobgu/qframe"
)

func main() {
	resp, err := http.Get("https://wiski.tirol.gv.at/hydro/ogd/OGD_W.csv")
	if err != nil {
		fmt.Println(err.Error())
	}

	f := qframe.ReadCSV(resp.Body)
	fmt.Println(f)

	defer resp.Body.Close()
}
