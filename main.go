package main

import (
	"fmt"
	http "net/http"
	"time"

	"github.com/tobgu/qframe"
	"github.com/tobgu/qframe/config/csv"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	resp, err := http.Get("https://wiski.tirol.gv.at/hydro/ogd/OGD_W.csv")
	if err != nil {
		fmt.Println(err.Error())
	}
	dtypes := map[string]string{
		"Stationsname":           "string",
		"Stationsnummer":         "string",
		"Gew�sser":               "string",
		"Parameter":              "string",
		"Zeitstempel in ISO8601": "string",
		"Wert":                   "float",
		"Einheit":                "string",
		"Seeh�he":                "float",
		"Rechtswert":             "float",
		"Hochwert":               "float",
		"EPSG-Code":              "string",
	}

	df := qframe.ReadCSV(resp.Body, csv.Delimiter(';'), csv.Types(dtypes))
	defer resp.Body.Close()

	// Create a new client using an InfluxDB server base URL and an authentication token
	// and set batch size to 20
	client := influxdb2.NewClientWithOptions("http://localhost:8086", "my-token",
		influxdb2.DefaultOptions().SetBatchSize(20))
	// Get non-blocking write client
	writeAPI := client.WriteAPI("my-org", "my-bucket")

	value := df.MustFloatView("Wert")
	datetime, _ := df.StringView("Zeitstempel in ISO8601")
	for i := 0; i < df.Len(); i++ {
		layout := "2006-01-02T15:04:05-0700"
		t, _ := time.Parse(layout, *datetime.ItemAt(i))
		p := influxdb2.NewPoint(
			"system",
			map[string]string{
				"id": "hello",
			},
			map[string]interface{}{
				"value": value.ItemAt(i),
			},
			t)
		fmt.Println(p)
	}
	// Force all unwritten data to be sent
	writeAPI.Flush()
	// Ensures background processes finishes
	client.Close()

}
