package dataset

import (
	"encoding/json"

	"github.com/qri-io/dataset/datatypes"
)

// Schema is analogous to a SQL schema definition
type Schema struct {
	Fields     []*Field `json:"fields,omitempty"`
	PrimaryKey FieldKey `json:"primaryKey,omitempty"`
}

// FieldNames gives a slice of field names defined in schema
func (s *Schema) FieldNames() (names []string) {
	if s.Fields == nil {
		return []string{}
	}
	names = make([]string, len(s.Fields))
	for i, f := range s.Fields {
		names[i] = f.Name
	}
	return
}

// FieldForName returns the field who's string name matches name
func (s *Schema) FieldForName(name string) *Field {
	for _, f := range s.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// FieldTypeStrings gives a slice of each field's type as a string
func (s *Schema) FieldTypeStrings() (types []string) {
	types = make([]string, len(s.Fields))
	for i, f := range s.Fields {
		types[i] = f.Type.String()
	}
	return
}

// Field is a field descriptor
type Field struct {
	Name         string            `json:"name"`
	Type         datatypes.Type    `json:"type,omitempty"`
	MissingValue interface{}       `json:"missingValue,omitempty"`
	Format       string            `json:"format,omitempty"`
	Constraints  *FieldConstraints `json:"constraints,omitempty"`
	Title        string            `json:"title,omitempty"`
	Description  string            `json:"description,omitempty"`
}

// field is a private struct for marshaling into and out of JSON
// most importantly, keys are sorted by lexographical order
type _field struct {
	Constraints  *FieldConstraints `json:"constraints,omitempty"`
	Description  string            `json:"description,omitempty"`
	Format       string            `json:"format,omitempty"`
	MissingValue interface{}       `json:"missingValue,omitempty"`
	Name         string            `json:"name"`
	Title        string            `json:"title,omitempty"`
	Type         datatypes.Type    `json:"type,omitempty"`
}

// MarshalJSON satisfies the json.Marshaler interface
func (f Field) MarshalJSON() ([]byte, error) {
	_f := &_field{
		Constraints:  f.Constraints,
		Description:  f.Description,
		Format:       f.Format,
		MissingValue: f.MissingValue,
		Name:         f.Name,
		Title:        f.Title,
		Type:         f.Type,
	}
	return json.Marshal(_f)
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (f *Field) UnmarshalJSON(data []byte) error {
	_f := &_field{}
	if err := json.Unmarshal(data, _f); err != nil {
		return err
	}

	*f = Field{
		Constraints:  _f.Constraints,
		Description:  _f.Description,
		Format:       _f.Format,
		MissingValue: _f.MissingValue,
		Name:         _f.Name,
		Title:        _f.Title,
		Type:         _f.Type,
	}
	return nil
}

// NewShortField is a shortcut for creating a field that only has a name and type
func NewShortField(name, dataType string) *Field {
	return &Field{
		Name: name,
		Type: datatypes.TypeFromString(dataType),
	}
}

// FieldKey allows a field key to be either a string or object
type FieldKey []string

// FieldConstraints is supposed to constrain the field,
// this is totally unfinished, unimplemented, and needs lots of work
// TODO - uh, finish this?
type FieldConstraints struct {
	Required  *bool         `json:"required,omitempty"`
	MinLength *int64        `json:"minLength,omitempty"`
	MaxLength *int64        `json:"maxLength,omitempty"`
	Unique    *bool         `json:"unique,omitempty"`
	Pattern   string        `json:"pattern,omitempty"`
	Minimum   interface{}   `json:"minimum,omitempty"`
	Maximum   interface{}   `json:"maximum,omitempty"`
	Enum      []interface{} `json:"enum,omitempty"`
}

// ForeignKey is supposed to be for supporting foreign key
// references. It's also totally unfinished.
// TODO - finish this
type ForeignKey struct {
	Fields FieldKey `json:"fields"`
	// Reference
}
