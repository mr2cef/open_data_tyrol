package tirPrec

import "github.com/mr2cef/open_data_tyrol/sources/common"

func GetSource() common.Source {

	precDTypes := map[string]string{
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

	prec := common.Source{
		Url:         "https://wiski.tirol.gv.at/hydro/ogd/OGD_N.csv",
		Dtypes:      precDTypes,
		Delimiter:   ';',
		Bucket:      "pegel",
		Measurement: "prec",
		Prefix:      "tirPrec",
		IdCol:       "Stationsnummer",
		ValCol:      "Wert",
		ValMin:      0.,
		ValMax:      10000.,
		TimeFmt:     "2006-01-02T15:04:05-0700",
		TimeCol:     "Zeitstempel in ISO8601",
	}

	return prec
}
