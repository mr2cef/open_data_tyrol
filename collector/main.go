package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"github.com/mr2cef/open_data_tyrol/collector/influx"
	"github.com/mr2cef/open_data_tyrol/collector/mongo"
	"github.com/mr2cef/open_data_tyrol/collector/sources/common"
	"github.com/mr2cef/open_data_tyrol/collector/sources/tirPeg"
	"github.com/mr2cef/open_data_tyrol/collector/sources/tirPrec"
	"github.com/mr2cef/open_data_tyrol/collector/sources/tirTemp"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func handleCollect(w http.ResponseWriter, r *http.Request) {
	ptsc := make(chan *write.Point, 10000)
	mdonec, idonec := make(chan string), make(chan string)
	stationc := make(chan map[string]string)
	var wg sync.WaitGroup
	sources := []common.Source{tirPeg.Source, tirPrec.Source, tirTemp.Source}

	for _, s := range sources {
		wg.Add(1)
		go s.GetDataPts(ptsc, stationc, &wg)
	}
	go influx.WriteDb(ptsc, idonec)
	go mongo.WriteDB(stationc, mdonec)
	wg.Wait()

	close(ptsc)
	close(stationc)
	fmt.Fprintf(w, <-idonec)
	fmt.Fprintf(w, <-mdonec)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wellcome.\n\nCurrently there are the following endpoints available:\n\n")
	fmt.Fprintf(w, "/collect: Collect data points and write to InfluxDB.")
}

func main() {

	godotenv.Load("../.env")

	http.HandleFunc("/collect", handleCollect)
	http.HandleFunc("/", handleRoot)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
