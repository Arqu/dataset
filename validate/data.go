package validate

import (
	"fmt"
	"regexp"

	"github.com/qri-io/dataset"
	"github.com/qri-io/dataset/dsio"
)

var alphaNumericRegex = regexp.MustCompile(`^[a-z0-9_-]{1-144}$`)

// truthCount returns the number of arguments that are true
func truthCount(args ...bool) (count int) {
	for _, arg := range args {
		if arg {
			count++
		}
	}
	return
}

type ErrFormat int

const (
	ErrFmtUnknown ErrFormat = iota
	ErrFmtOneHotMatrix
	ErrFmtErrStrings
)

type ValidateDataOpt struct {
	ErrorFormat ErrFormat
	DataFormat  DataFormat
}

func ValidateData(r dsio.RowReader, options ...func(*ValidateDataOpt)) (errors dsio.RowReader, count int, err error) {
	vst := &dataset.Structure{
		Format: CsvDataFormat,
		Schema: &dataset.Schema{
			Fields: []*dataset.Field{
				&dataset.Field{Name: "row_index", Type: datatype.Integer},
			},
		},
	}
	for _, f := range r.Structure().Schema.Fields {
		vst.Schema.Fields = append(vst.Schema.Fields, &Field{Name: f.Name + "_error", Type: datatype.String})
	}

	buf := dsio.NewBuffer(vst)

	err = dsio.EachRow(r, func(num int, row [][]byte, err error) error {
		if err != nil {
			return err
		}

		errData, errNum, err := validateRow(ds.Fields, num, row)
		if err != nil {
			return err
		}

		count += errNum
		if errNum != 0 {
			buf.WriteRow(errData)
		}

		return nil
	})

	if err = buf.Close(); err != nil {
		err = fmt.Errorf("error closing valdation buffer: %s", err.Error())
		return
	}

	errors = buf
	return
}

func validateRow(fields []*Field, num int, row [][]byte) ([][]byte, int, error) {
	count := 0
	errors := make([][]byte, len(fields)+1)
	errors[0] = []byte(strconv.FormatInt(int64(num), 10))
	if len(row) != len(fields) {
		return errors, count, fmt.Errorf("column mismatch. expected: %d, got: %d", len(fields), len(row))
	}

	for i, f := range fields {
		_, e := f.Type.Parse(row[i])
		if e != nil {
			count++
			errors[i+1] = []byte(e.Error())
		} else {
			errors[i+1] = []byte("")
		}
	}

	return errors, count, nil
}

// func (ds *Resource) ValidateDeadLinks(store fs.Store) (validation *Resource, data []byte, count int, err error) {
// 	proj := map[int]int{}
// 	validation = &Resource{
// 		Address: NewAddress(ds.Address.String(), "errors"),
// 		Format:  CsvDataFormat,
// 	}

// 	for i, f := range ds.Fields {
// 		if f.Type == datatype.Url {
// 			proj[i] = len(validation.Fields)
// 			validation.Fields = append(validation.Fields, f)
// 			validation.Fields = append(validation.Fields, &Field{Name: f.Name + "_dead", Type: datatype.Integer})
// 		}
// 	}

// 	dsData, e := ds.FetchBytes(store)
// 	if e != nil {
// 		err = e
// 		return
// 	}
// 	ds.Data = dsData

// 	buf := &bytes.Buffer{}
// 	cw := csv.NewWriter(buf)

// 	err = ds.EachRow(func(num int, row [][]byte, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		result := make([][]byte, len(validation.Fields))
// 		for l, r := range proj {
// 			result[r] = row[l]
// 			if err := checkUrl(string(result[r])); err != nil {
// 				count++
// 				result[r+1] = []byte("1")
// 			} else {
// 				result[r+1] = []byte("0")
// 			}
// 		}

// 		csvRow := make([]string, len(result))
// 		for i, d := range result {
// 			csvRow[i] = string(d)
// 		}
// 		if err := cw.Write(csvRow); err != nil {
// 			fmt.Sprintln(err)
// 		}

// 		return nil
// 	})

// 	cw.Flush()
// 	data = buf.Bytes()
// 	return
// }

// func checkUrl(rawurl string) error {
// 	u, err := url.Parse(rawurl)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	if u.Scheme == "" {
// 		u.Scheme = "http"
// 	}
// 	res, err := http.Get(u.String())
// 	if err != nil {
// 		return err
// 	}
// 	res.Body.Close()
// 	fmt.Println(u.String(), res.StatusCode)
// 	if res.StatusCode > 399 {
// 		return fmt.Errorf("non-200 status code")
// 	}
// 	return nil
// }
