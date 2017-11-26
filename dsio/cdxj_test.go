package dsio

import (
	"bytes"
	"testing"

	"github.com/qri-io/dataset"
	"github.com/qri-io/dataset/datatypes"
)

// TODO - vary up test input
const cdxjData = `!OpenWayback-CDXJ 1.0
(com,cnn,)/world 2015-09-03T13:27:52Z response {"a" : 0, "b" : "b", "c" : false }
(com,cnn,)/world 2015-09-03T13:27:52Z response {"a" : 0, "b" : "b", "c" : false }
(com,cnn,)/world 2015-09-03T13:27:52Z response {"a" : 0, "b" : "b", "c" : false }
(com,cnn,)/world 2015-09-03T13:27:52Z response {"a" : 0, "b" : "b", "c" : false }
`

var cdxjStruct = &dataset.Structure{
	Format: dataset.CDXJDataFormat,
	Schema: &dataset.Schema{
		Fields: []*dataset.Field{
			{Name: "surt_uri", Type: datatypes.String},
			// TODO - currently using string b/c date interface isn't fully implemented
			{Name: "timestamp", Type: datatypes.String},
			{Name: "record_type", Type: datatypes.String},
			{Name: "metadata", Type: datatypes.JSON},
		},
	},
}

func TestCdxjReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte(cdxjData))
	rdr, err := NewRowReader(cdxjStruct, buf)
	if err != nil {
		t.Errorf("error allocating rowReader: %s", err.Error())
		return
	}
	count := 0
	for {
		row, err := rdr.ReadRow()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Errorf("unexpected error: %s", err.Error())
			return
		}

		if len(row) != 4 {
			t.Errorf("invalid row length for row %d. expected %d, got %d", count, 4, len(row))
		}

		count++
	}
	if count != 4 {
		t.Errorf("expected: %d rows, got: %d", 4, count)
	}
}

func TestCdxjWriter(t *testing.T) {
	rows := [][][]byte{
		// TODO - vary up test input
		{[]byte("(com,cnn,)/world"), []byte("2015-09-03T13:27:52Z"), []byte("response"), []byte(`{"a" : 0, "b" : "b", "c" : false }`)},
		{[]byte("(com,cnn,)/world"), []byte("2015-09-03T13:27:52Z"), []byte("response"), []byte(`{"a" : 0, "b" : "b", "c" : false }`)},
		{[]byte("(com,cnn,)/world"), []byte("2015-09-03T13:27:52Z"), []byte("response"), []byte(`{"a" : 0, "b" : "b", "c" : false }`)},
		{[]byte("(com,cnn,)/world"), []byte("2015-09-03T13:27:52Z"), []byte("response"), []byte(`{"a" : 0, "b" : "b", "c" : false }`)},
	}

	buf := &bytes.Buffer{}
	rw, err := NewRowWriter(cdxjStruct, buf)
	if err != nil {
		t.Errorf("error allocating RowWriter: %s", err.Error())
		return
	}
	st := rw.Structure()
	if err := dataset.CompareStructures(&st, cdxjStruct); err != nil {
		t.Errorf("structure mismatch: %s", err.Error())
		return
	}

	for i, row := range rows {
		if err := rw.WriteRow(row); err != nil {
			t.Errorf("row %d write error: %s", i, err.Error())
		}
	}

	if err := rw.Close(); err != nil {
		t.Errorf("close reader error: %s", err.Error())
		return
	}

	if bytes.Equal(buf.Bytes(), []byte(cdxjData)) {
		t.Errorf("output mismatch. %s != %s", buf.String(), cdxjData)
	}
}
