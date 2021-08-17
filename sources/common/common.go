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
	ValMin      float64
	ValMax      float64
	TimeCol     string
	TimeFmt     string
}

type DataSource interface {
	GetDataPts()
}

func (s Source) GetDataPts(ptsc chan *write.Point, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(s.Url)
	if err != nil {
		fmt.Println(err.Error())
	}
	df := qframe.ReadCSV(resp.Body, csv.Delimiter(s.Delimiter), csv.Types(s.Dtypes))
	defer resp.Body.Close()
	log.Printf("Downloaded %s containing %d data points.\n", s.Url, df.Len())
	value := df.MustFloatView(s.ValCol)
	datetime := df.MustStringView(s.TimeCol)
	station := df.MustStringView(s.IdCol)
	for i := 0; i < df.Len(); i++ {
		t, err := time.Parse(s.TimeFmt, *datetime.ItemAt(i))
		if err != nil {
			log.Printf("%s; Cannot parse %s to fromat %s", err, s.TimeFmt, *datetime.ItemAt(i))
		}
		value := value.ItemAt(i)
		if (value >= s.ValMin) && (value <= s.ValMax) {
			p := influxdb2.NewPoint(
				s.Measurement,
				map[string]string{
					"stationId": (s.Prefix + *station.ItemAt(i)),
				},
				map[string]interface{}{
					"value": value,
				},
				t)
			ptsc <- p
		}
	}
}
