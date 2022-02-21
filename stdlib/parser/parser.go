package parser

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"
)

// Options hold configurations of respective parser type.
type Options struct {
	JSON  JSONOptions
}

type Parser interface {
	// JSONParser return JSONParser Object
	JSONParser() JSONParser
}

type JSONOptions struct {
	// Config based on jsoniterator-go default config
	Config jsonConfig
	// Custom config to be frozen for go-jsoniterator Config
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	EscapeHTML                    bool
	SortMapKeys                   bool
	UseNumber                     bool
	DisallowUnknownFields         bool
	TagKey                        string
	OnlyTaggedField               bool
	ValidateJSONRawMessage        bool
	ObjectFieldMustBeSimpleString bool
	CaseSensitive                 bool
	// Schema contains schema definitions with key as schema name and value as source path
	// schema sources can be file or URL. Schema definition will be iniatialized during
	// JSON parser object initialization.
	Schema map[string]string
}

type parser struct {
	json  JSONParser
	opt   Options
}


// Init Main Parser Object.
func Init(opt Options) Parser {
	return &parser{
		json:  initJSONP(opt.JSON),
	}
}

func (p *parser) JSONParser() JSONParser {
	return p.json
}

// initJSONP initialize JSON parser with logger and declared options
func initJSONP(opt JSONOptions) JSONParser {
	var jsonAPI jsoniter.API
	switch opt.Config {

	case JSONConfigDefault:
		jsonAPI = jsoniter.ConfigDefault

	case JSONConfigFastest:
		jsonAPI = jsoniter.ConfigFastest

	case JSONConfigCompatibleWithStdLibrary:
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary

	case JSONConfigCustom:
		jsonAPI = jsoniter.Config{
			IndentionStep:                 opt.IndentionStep,
			MarshalFloatWith6Digits:       opt.MarshalFloatWith6Digits,
			EscapeHTML:                    opt.EscapeHTML,
			SortMapKeys:                   opt.SortMapKeys,
			UseNumber:                     opt.UseNumber,
			DisallowUnknownFields:         opt.DisallowUnknownFields,
			TagKey:                        opt.TagKey,
			OnlyTaggedField:               opt.OnlyTaggedField,
			ValidateJsonRawMessage:        opt.ValidateJSONRawMessage,
			ObjectFieldMustBeSimpleString: opt.ObjectFieldMustBeSimpleString,
			CaseSensitive:                 opt.CaseSensitive,
		}.Froze()

	default:
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary
	}
	p := &jsonparser{
		API:    jsonAPI,
		opt:    opt,
		schema: make(map[string]*gojsonschema.Schema),
	}
	//init defined schema
	p.initSchema(opt.Schema)
	return p
}

func (p *jsonparser) initSchema(sources map[string]string) {
	for sch, src := range sources {
		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(src))
		if err != nil {
			panic(err)
			return
		}
		p.schema[sch] = schema
	}
}
