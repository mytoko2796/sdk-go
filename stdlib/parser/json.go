package parser

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"
)

type JSONParser interface {
	// Marshal go structs into bytes
	Marshal(orig interface{}) ([]byte, error)
	// Unmarshal bytes into go structs
	Unmarshal(blob []byte, dest interface{}) error
}

type jsonConfig string

const (
	// JSONConfigDefault set :
	//	EscapeHTML :					true
	JSONConfigDefault jsonConfig = `default`

	// JSONConfigCompatibleWithStdLibrary
	//  EscapeHTML:             		true
	//  SortMapKeys:					true
	//  ValidateJsonRawMessage:			true
	JSONConfigCompatibleWithStdLibrary jsonConfig = `standard`

	// JSONConfigFastest
	//  EscapeHTML:                    	false
	//  MarshalFloatWith6Digits:       	true
	//  ObjectFieldMustBeSimpleString: 	true
	JSONConfigFastest jsonConfig = `fastest`

	// JSONConfigCustom
	//	Custom Configuration which is set in JSONOptions
	JSONConfigCustom jsonConfig = `custom`
)

type jsonparser struct {
	schema map[string]*gojsonschema.Schema
	API    jsoniter.API
	opt    JSONOptions
}

func (p *jsonparser) Marshal(orig interface{}) ([]byte, error) {
	stream := p.API.BorrowStream(nil)
	defer p.API.ReturnStream(stream)
	stream.WriteVal(orig)
	result := make([]byte, stream.Buffered())
	if stream.Error != nil {
		return nil, stream.Error
	}
	copy(result, stream.Buffer())
	return result, nil
}
func (p *jsonparser) Unmarshal(blob []byte, dest interface{}) error {
	iter := p.API.BorrowIterator(blob)
	defer p.API.ReturnIterator(iter)
	iter.ReadVal(dest)
	if iter.Error != nil {
		return iter.Error
	}
	return nil
}
