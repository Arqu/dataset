package dsfs

import (
	"fmt"
	"io"

	"github.com/ipfs/go-ipfs/commands/files"
	"github.com/qri-io/cafs"
	"github.com/qri-io/dataset"
	"github.com/qri-io/dataset/dsio"
)

// LoadDatasetData loads the data this dataset points to from the store
func LoadData(store cafs.Filestore, ds *dataset.Dataset) (files.File, error) {
	return store.Get(ds.Data)
}

// ReadRows loads a slice of raw bytes inside a limit/offset row range
func LoadRows(store cafs.Filestore, ds *dataset.Dataset, limit, offset int) ([]byte, error) {

	datafile, err := LoadData(store, ds)
	if err != nil {
		return nil, fmt.Errorf("error loading dataset data: %s", err.Error())
	}

	added := 0
	buf := dsio.NewBuffer(ds.Structure)
	rr := dsio.NewRowReader(ds.Structure, datafile)
	err = dsio.EachRow(rr, func(i int, row [][]byte, err error) error {
		if err != nil {
			return err
		}

		if i < offset {
			return nil
		} else if limit > 0 && added == limit {
			return io.EOF
		}

		buf.WriteRow(row)
		added++
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating through dataset data: %s", err.Error())
	}

	err = buf.Close()
	return buf.Bytes(), err
}
