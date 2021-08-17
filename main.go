package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"github.com/mr2cef/open_data_tyrol/influx"
	"github.com/mr2cef/open_data_tyrol/sources/common"
	"github.com/mr2cef/open_data_tyrol/sources/tirPeg"
	"github.com/mr2cef/open_data_tyrol/sources/tirPrec"
	"github.com/mr2cef/open_data_tyrol/sources/tirTemp"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func handleCollect(w http.ResponseWriter, r *http.Request) {
	ptsc := make(chan *write.Point, 10000)
	donec := make(chan string)
	var wg sync.WaitGroup
	sources := []common.Source{tirPeg.Source, tirPrec.Source, tirTemp.Source}

	for _, s := range sources {
		wg.Add(1)
		go s.GetDataPts(ptsc, &wg)
	}
	go influx.WriteDb(ptsc, donec)
	wg.Wait()

	close(ptsc)
	fmt.Fprintf(w, <-donec)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wellcome.\n\nCurrently there are the following endpoints available:\n\n")
	fmt.Fprintf(w, "/collect: Collect data points and write to InfluxDB.")
}

func main() {

	godotenv.Load(".env")

	http.HandleFunc("/collect", handleCollect)
	http.HandleFunc("/", handleRoot)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
