package main

import (
	"fmt"
	http "net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tobgu/qframe"
	"github.com/tobgu/qframe/config/csv"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {

	godotenv.Load(".env")

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
	client := influxdb2.NewClientWithOptions(os.Getenv("INFLUX_DB_HOST"), os.Getenv("INFLUX_DB_TOCKEN"),
		influxdb2.DefaultOptions().SetBatchSize(1000))
	// Get non-blocking write client
	writeAPI := client.WriteAPI(os.Getenv("INFLUX_DB_ORG"), "pegel")

	value := df.MustFloatView("Wert")
	datetime, _ := df.StringView("Zeitstempel in ISO8601")
	station, _ := df.StringView("Stationsnummer")
	for i := 0; i < df.Len(); i++ {
		layout := "2006-01-02T15:04:05-0700"
		t, _ := time.Parse(layout, *datetime.ItemAt(i))
		value := value.ItemAt(i)
		if value >= 0 {
			p := influxdb2.NewPoint(
				"system",
				map[string]string{
					"stationId": *station.ItemAt(i),
				},
				map[string]interface{}{
					"value": value,
				},
				t)
			writeAPI.WritePoint(p)
		}
	}
	// Force all unwritten data to be sent
	writeAPI.Flush()
	// Ensures background processes finishes
	client.Close()

}
