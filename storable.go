package dataset

import (
	"fmt"

	"github.com/ipfs/go-datastore"
	"github.com/qri-io/cafs"
)

// ErrNoPath is a canonical error for when an expected path is missing
var ErrNoPath = fmt.Errorf("missing path")

// storable is the internal interface for anything that can save / load from a
// content content-addressed store
type storable interface {
	Load(store cafs.Filestore) error
	Save(store cafs.Filestore) (datastore.Key, error)
}
