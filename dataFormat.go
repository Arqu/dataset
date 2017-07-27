// TODO - consider placing this in a subpackage: dataformats
package dataset

import (
	"encoding/json"
	"fmt"
)

// DataFormat represents different types of data
type DataFormat int

const (
	UnknownDataFormat DataFormat = iota
	CsvDataFormat
	JsonDataFormat
	JsonArrayDataFormat
	XmlDataFormat
	XlsDataFormat
	// TODO - make this list more exhaustive
)

// String implements stringer interface for DataFormat
func (f DataFormat) String() string {
	s, ok := map[DataFormat]string{
		UnknownDataFormat:   "",
		CsvDataFormat:       "csv",
		JsonDataFormat:      "json",
		JsonArrayDataFormat: "jsona",
		XmlDataFormat:       "xml",
		XlsDataFormat:       "xls",
	}[f]

	if !ok {
		return ""
	}

	return s
}

// ParseDataFormatString takes a string representation of a data format
func ParseDataFormatString(s string) (df DataFormat, err error) {
	df, ok := map[string]DataFormat{
		"":       UnknownDataFormat,
		".csv":   CsvDataFormat,
		"csv":    CsvDataFormat,
		".json":  JsonDataFormat,
		"json":   JsonDataFormat,
		".jsona": JsonArrayDataFormat,
		"jsona":  JsonArrayDataFormat,
		".xml":   XmlDataFormat,
		"xml":    XmlDataFormat,
		".xls":   XlsDataFormat,
		"xls":    XlsDataFormat,
	}[s]
	if !ok {
		err = fmt.Errorf("invalid DataFormat %q", s)
		df = UnknownDataFormat
	}

	return
}

// MarshalJSON satisfies the json.Marshaler interface
func (f DataFormat) MarshalJSON() ([]byte, error) {
	if f == UnknownDataFormat {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (f *DataFormat) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Filed type should be a string, got %s", data)
	}

	if df, err := ParseDataFormatString(s); err != nil {
		return err
	} else {
		*f = df
	}

	return nil
}
