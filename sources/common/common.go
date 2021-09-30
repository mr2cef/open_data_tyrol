package common

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/tobgu/qframe"
	"github.com/tobgu/qframe/config/csv"
	"golang.org/x/text/encoding/charmap"
)

type Source struct {
	Url         string
	Dtypes      map[string]string
	Delimiter   byte
	Bucket      string
	Measurement string
	Prefix      string
	IdCol       string
	ValCol      string
	NameCol     string
	RightCol    string
	HightCol    string
	StdCol      string
	ValMin      float64
	ValMax      float64
	TimeCol     string
	TimeFmt     string
}

type DataSource interface {
	GetDataPts()
}

func (s Source) GetDataPts(ptsc chan *write.Point, stationc chan map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(s.Url)
	d := charmap.ISO8859_1.NewDecoder()
	if err != nil {
		fmt.Println(err.Error())
	}
	df := qframe.ReadCSV(d.Reader(resp.Body), csv.Delimiter(s.Delimiter), csv.Types(s.Dtypes))
	defer resp.Body.Close()
	log.Printf("Downloaded %s containing %d data points.\n", s.Url, df.Len())
	value := df.MustFloatView(s.ValCol)
	datetime := df.MustStringView(s.TimeCol)
	station := df.MustStringView(s.IdCol)
	name := df.MustStringView(s.NameCol)
	rightVal := df.MustFloatView(s.RightCol)
	highVal := df.MustFloatView(s.HightCol)
	std := df.MustStringView(s.StdCol)
	for i := 0; i < df.Len(); i++ {
		t, err := time.Parse(s.TimeFmt, *datetime.ItemAt(i))
		if err != nil {
			log.Printf("%s; Cannot parse %s to fromat %s", err, s.TimeFmt, *datetime.ItemAt(i))
		}
		value := value.ItemAt(i)
		if (value >= s.ValMin) && (value <= s.ValMax) {
			id := (s.Prefix + *station.ItemAt(i))
			p := influxdb2.NewPoint(
				s.Measurement,
				map[string]string{
					"stationId": id,
				},
				map[string]interface{}{
					"value": value,
				},
				t)
			ptsc <- p
			m := map[string]string{
				"_id":    id,
				"name":   *name.ItemAt(i),
				"right":  fmt.Sprintf("%f", rightVal.ItemAt(i)),
				"high":   fmt.Sprintf("%f", highVal.ItemAt(i)),
				"format": *std.ItemAt(i),
			}
			stationc <- m
		}
	}
}
