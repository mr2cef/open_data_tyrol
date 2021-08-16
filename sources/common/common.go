package common

import (
	"fmt"
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
	value := df.MustFloatView(s.ValCol)
	datetime, err := df.StringView(s.TimeCol)
	if err != nil {
		// TODO
	}
	station, err := df.StringView(s.IdCol)
	if err != nil {
		// TODO
	}
	for i := 0; i < df.Len(); i++ {
		t, err := time.Parse(s.TimeFmt, *datetime.ItemAt(i))
		if err != nil {
			// TODO
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
