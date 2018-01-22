package dataset

import (
	"encoding/json"
	"fmt"

	"github.com/ipfs/go-datastore"
	"github.com/qri-io/dataset/compression"
	"github.com/qri-io/jsonschema"
)

// Structure designates a deterministic definition for working with a discrete dataset.
// Structure is a concrete handle that provides precise details about how to interpret a given
// piece of data (the reference to the data itself is provided elsewhere, specifically in the dataset struct )
// These techniques provide mechanisms for joining & traversing multiple structures.
// This example is shown in a human-readable form, for storage on the network the actual
// output would be in a condensed, non-indented form, with keys sorted by lexographic order.
type Structure struct {
	// private storage for reference to this object
	path datastore.Key
	// Checksum is a bas58-encoded multihash checksum of the data
	// file this structure points to. This is different from IPFS
	// hashes, which are calculated after breaking the file into blocks
	Checksum string `json:"checksum,omitempty"`
	// Compression specifies any compression on the source data,
	// if empty assume no compression
	Compression compression.Type `json:"compression,omitempty"`
	// Encoding specifics character encoding
	// should assume utf-8 if not specified
	Encoding string `json:"encoding,omitempty"`
	// Entries is number of top-level entries in the dataset. With tablular data
	// this is the same as the number of rows
	// required when structure is concrete, and must match underlying dataset.
	Entries int `json:"entries,omitempty"`
	// Format specifies the format of the raw data MIME type
	Format DataFormat `json:"format"`
	// FormatConfig removes as much ambiguity as possible about how
	// to interpret the speficied format.
	FormatConfig FormatConfig `json:"formatConfig,omitempty"`
	// Kind should always be KindStructure
	Kind Kind `json:"kind"`
	// Length is the length of the data object in bytes.
	// must always match & be present
	Length int `json:"length,omitempty"`
	// Schema contains the schema definition for the underlying data
	Schema *jsonschema.RootSchema `json:"schema,omitempty"`
}

// Path gives the internal path reference for this structure
func (s *Structure) Path() datastore.Key {
	return s.path
}

// NewStructureRef creates an empty struct with it's
// internal path set
func NewStructureRef(path datastore.Key) *Structure {
	return &Structure{Kind: KindStructure, path: path}
}

// Abstract returns this structure instance in it's "Abstract" form
// stripping all nonessential values &
// renaming all schema field names to standard variable names
func (s *Structure) Abstract() *Structure {
	a := &Structure{
		Format:       s.Format,
		FormatConfig: s.FormatConfig,
		Encoding:     s.Encoding,
	}
	if s.Schema != nil {
		// TODO - Fix meeeeeeee
		// a.Schema = &Schema{
		// 	PrimaryKey: s.Schema.PrimaryKey,
		// 	Fields:     make([]*Field, len(s.Schema.Fields)),
		// }
		// for i, f := range s.Schema.Fields {
		// 	a.Schema.Fields[i] = &Field{
		// 		Name:         AbstractColumnName(i),
		// 		Type:         f.Type,
		// 		MissingValue: f.MissingValue,
		// 		Format:       f.Format,
		// 		Constraints:  f.Constraints,
		// 	}
		// }
	}
	return a
}

// Hash gives the hash of this structure
func (s *Structure) Hash() (string, error) {
	return JSONHash(s)
}

// separate type for marshalling into & out of
// most importantly, struct names must be sorted lexographically
type _structure struct {
	Checksum     string                 `json:"checksum,omitempty"`
	Compression  compression.Type       `json:"compression,omitempty"`
	Encoding     string                 `json:"encoding,omitempty"`
	Entries      int                    `json:"entries,omitempty"`
	Format       DataFormat             `json:"format"`
	FormatConfig map[string]interface{} `json:"formatConfig,omitempty"`
	Kind         Kind                   `json:"kind"`
	Length       int                    `json:"length,omitempty"`
	Schema       *jsonschema.RootSchema `json:"schema,omitempty"`
}

// MarshalJSON satisfies the json.Marshaler interface
func (s Structure) MarshalJSON() (data []byte, err error) {
	if s.path.String() != "" && s.Encoding == "" && s.Schema == nil {
		return s.path.MarshalJSON()
	}

	kind := s.Kind
	if kind == "" {
		kind = KindStructure
	}

	var opt map[string]interface{}
	if s.FormatConfig != nil {
		opt = s.FormatConfig.Map()
	}

	return json.Marshal(&_structure{
		Checksum:     s.Checksum,
		Compression:  s.Compression,
		Encoding:     s.Encoding,
		Entries:      s.Entries,
		Format:       s.Format,
		FormatConfig: opt,
		Kind:         kind,
		Length:       s.Length,
		Schema:       s.Schema,
	})
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (s *Structure) UnmarshalJSON(data []byte) (err error) {
	var (
		str    string
		fmtCfg FormatConfig
	)
	if err := json.Unmarshal(data, &str); err == nil {
		*s = Structure{path: datastore.NewKey(str)}
		return nil
	}

	_s := &_structure{}
	if err := json.Unmarshal(data, _s); err != nil {
		return fmt.Errorf("error unmarshaling dataset structure from json: %s", err.Error())
	}

	if _s.FormatConfig != nil {
		fmtCfg, err = ParseFormatConfigMap(_s.Format, _s.FormatConfig)
		if err != nil {
			return fmt.Errorf("error parsing structure formatConfig: %s", err.Error())
		}

	}

	*s = Structure{
		Checksum:     _s.Checksum,
		Compression:  _s.Compression,
		Encoding:     _s.Encoding,
		Entries:      _s.Entries,
		Format:       _s.Format,
		FormatConfig: fmtCfg,
		Kind:         _s.Kind,
		Length:       _s.Length,
		Schema:       _s.Schema,
	}
	return nil
}

// IsEmpty checks to see if structure has any fields other than the internal path
func (s *Structure) IsEmpty() bool {
	return s.Checksum == "" &&
		s.Compression == compression.None &&
		s.Encoding == "" &&
		s.Entries == 0 &&
		s.Format == UnknownDataFormat &&
		s.FormatConfig == nil &&
		s.Length == 0 &&
		s.Schema == nil
}

// Assign collapses all properties of a group of structures on to one
// this is directly inspired by Javascript's Object.assign
func (s *Structure) Assign(structures ...*Structure) {
	for _, st := range structures {
		if st == nil {
			continue
		}

		if st.path.String() != "" {
			s.path = st.path
		}
		if st.Checksum != "" {
			s.Checksum = st.Checksum
		}
		if st.Compression != compression.None {
			s.Compression = st.Compression
		}
		if st.Encoding != "" {
			s.Encoding = st.Encoding
		}
		if st.Entries != 0 {
			s.Entries = st.Entries
		}
		if st.Format != UnknownDataFormat {
			s.Format = st.Format
		}
		if st.FormatConfig != nil {
			s.FormatConfig = st.FormatConfig
		}
		if st.Kind != "" {
			s.Kind = st.Kind
		}
		if st.Length != 0 {
			s.Length = st.Length
		}
		// TODO - fix me
		if st.Schema != nil {
			// if s.Schema == nil {
			// 	s.Schema = &RootSchema{}
			// }
			// s.Schema.Assign(st.Schema)
			s.Schema = st.Schema
		}
	}
}

// UnmarshalStructure tries to extract a structure type from an empty
// interface. Pairs nicely with datastore.Get() from github.com/ipfs/go-datastore
func UnmarshalStructure(v interface{}) (*Structure, error) {
	switch r := v.(type) {
	case *Structure:
		return r, nil
	case Structure:
		return &r, nil
	case []byte:
		structure := &Structure{}
		err := json.Unmarshal(r, structure)
		return structure, err
	default:
		return nil, fmt.Errorf("couldn't parse structure, value is invalid type")
	}
}

// AbstractTableName prepends a given index with "t"
// t1, t2, t3, ...
func AbstractTableName(i int) string {
	return fmt.Sprintf("t%d", i+1)
}

// AbstractColumnName is the "base26" value of a column name
// to make short, sql-valid, deterministic column names
func AbstractColumnName(i int) string {
	return base26(i)
}

// b26chars is a-z, lowercase
const b26chars = "abcdefghijklmnopqrstuvwxyz"

// base26 maps the set of natural numbers
// to letters, using repeating characters to handle values
// greater than 26
func base26(d int) (s string) {
	var cols []int
	if d == 0 {
		return "a"
	}

	for d != 0 {
		cols = append(cols, d%26)
		d = d / 26
	}
	for i := len(cols) - 1; i >= 0; i-- {
		if i != 0 && cols[i] > 0 {
			s += string(b26chars[cols[i]-1])
		} else {
			s += string(b26chars[cols[i]])
		}
	}
	return s
}
