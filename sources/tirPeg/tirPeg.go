package tirPeg

import (
	"github.com/mr2cef/open_data_tyrol/sources/common"
)

func GetSource() common.Source {

	pegelDTypes := map[string]string{
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

	pegel := common.Source{
		Url:         "https://wiski.tirol.gv.at/hydro/ogd/OGD_W.csv",
		Dtypes:      pegelDTypes,
		Delimiter:   ';',
		Bucket:      "pegel",
		Measurement: "pegel",
		Prefix:      "tirPeg",
		IdCol:       "Stationsnummer",
		ValCol:      "Wert",
		ValMin:      0.,
		ValMax:      10000.,
		TimeFmt:     "2006-01-02T15:04:05-0700",
		TimeCol:     "Zeitstempel in ISO8601",
	}
	return pegel
}
