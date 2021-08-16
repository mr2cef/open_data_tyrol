package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/mr2cef/open_data_tyrol/sources/common"
	"github.com/mr2cef/open_data_tyrol/sources/tirPeg"
	"github.com/mr2cef/open_data_tyrol/sources/tirPrec"
	"github.com/mr2cef/open_data_tyrol/sources/tirTemp"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func writeDb(ptsc chan *write.Point, donec chan bool) {
	i := 0
	defer func() {
		fmt.Println("Done.", i, "points written.")
		donec <- true
	}()
	// Create a new client using an InfluxDB server base URL and an authentication token
	// and set batch size to 20
	client := influxdb2.NewClientWithOptions(
		os.Getenv("INFLUX_DB_HOST"),
		os.Getenv("INFLUX_DB_TOCKEN"),
		influxdb2.DefaultOptions().SetBatchSize(1000).SetPrecision(time.Minute),
	)
	// Ensures background processes finishes
	defer client.Close()
	// Get non-blocking write client
	writeAPI := client.WriteAPI(
		os.Getenv("INFLUX_DB_ORG"),
		"open_data")
	// Force all unwritten data to be sent
	defer writeAPI.Flush()
	for p := range ptsc {
		writeAPI.WritePoint(p)
		i++
	}
}

func main() {

	godotenv.Load(".env")

	ptsc := make(chan *write.Point, 10000)
	donec := make(chan bool)
	defer func() {
		<-donec
	}()
	var wg sync.WaitGroup
	sources := []common.Source{tirPeg.Source, tirPrec.Source, tirTemp.Source}

	for _, s := range sources {
		wg.Add(1)
		go s.GetDataPts(ptsc, &wg)
	}
	go writeDb(ptsc, donec)
	wg.Wait()

	close(ptsc)
}
