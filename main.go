package main

import (
	"fmt"
	http "net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/mr2cef/open_data_tyrol/sources/common"
	"github.com/mr2cef/open_data_tyrol/sources/tirPeg"
	"github.com/mr2cef/open_data_tyrol/sources/tirPrec"
	"github.com/mr2cef/open_data_tyrol/sources/tirTemp"
	"github.com/tobgu/qframe"
	"github.com/tobgu/qframe/config/csv"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func getDataPts(s common.Source, ptsc chan *write.Point, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(s.Url)
	if err != nil {
		fmt.Println(err.Error())
	}
	df := qframe.ReadCSV(resp.Body, csv.Delimiter(s.Delimiter), csv.Types(s.Dtypes))
	defer resp.Body.Close()
	value := df.MustFloatView(s.ValCol)
	datetime, _ := df.StringView(s.TimeCol)
	station, _ := df.StringView(s.IdCol)
	for i := 0; i < df.Len(); i++ {
		t, _ := time.Parse(s.TimeFmt, *datetime.ItemAt(i))
		value := value.ItemAt(i)
		if (value >= s.ValMin) && (value <= s.ValMax) {
			p := influxdb2.NewPoint(
				s.Measurement,
				map[string]string{
					"stationId": *station.ItemAt(i),
				},
				map[string]interface{}{
					"value": value,
				},
				t)
			ptsc <- p
		}
	}
}

func writeDb(ptsc chan *write.Point, donec chan bool) {
	i := 0
	defer func() {
		fmt.Println("Done. ", i, " points written.")
		donec <- true
	}()
	// Create a new client using an InfluxDB server base URL and an authentication token
	// and set batch size to 20
	client := influxdb2.NewClientWithOptions(
		os.Getenv("INFLUX_DB_HOST"),
		os.Getenv("INFLUX_DB_TOCKEN"),
		influxdb2.DefaultOptions().SetBatchSize(1000))
	// Ensures background processes finishes
	defer client.Close()
	// Get non-blocking write client
	writeAPI := client.WriteAPI(
		os.Getenv("INFLUX_DB_ORG"),
		"pegel")
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
	var wg sync.WaitGroup
	wg.Add(1)
	go getDataPts(tirPeg.GetSource(), ptsc, &wg)
	wg.Add(1)
	go getDataPts(tirPrec.GetSource(), ptsc, &wg)
	wg.Add(1)
	go getDataPts(tirTemp.GetSource(), ptsc, &wg)

	go writeDb(ptsc, donec)
	wg.Wait()
	close(ptsc)
	<-donec
}
