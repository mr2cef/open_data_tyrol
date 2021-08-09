package tirTemp

import "github.com/mr2cef/open_data_tyrol/sources/common"

func GetSource() common.Source {

	tempDTypes := map[string]string{
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

	temp := common.Source{
		Url:         "https://wiski.tirol.gv.at/hydro/ogd/OGD_LT.csv",
		Dtypes:      tempDTypes,
		Delimiter:   ';',
		Bucket:      "pegel",
		Measurement: "temp",
		Prefix:      "tirTemp",
		IdCol:       "Stationsnummer",
		ValCol:      "Wert",
		ValMin:      -60.,
		ValMax:      100.,
		TimeFmt:     "2006-01-02T15:04:05-0700",
		TimeCol:     "Zeitstempel in ISO8601",
	}

	return temp
}