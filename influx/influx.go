package influx

import (
	"fmt"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func WriteDb(ptsc chan *write.Point, donec chan string) {
	i := 0
	defer func() {
		response := fmt.Sprintf("Done. %d points written.", i)
		log.Println(response)
		donec <- response
	}()
	// Create a new client using an InfluxDB server base URL and an authentication token
	// and set batch size to 20
	client := influxdb2.NewClientWithOptions(
		os.Getenv("INFLUX_DB_HOST"),
		os.Getenv("INFLUX_DB_TOCKEN"),
		influxdb2.DefaultOptions().SetBatchSize(1000).SetPrecision(time.Second),
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
