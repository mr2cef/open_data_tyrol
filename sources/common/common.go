package common

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
