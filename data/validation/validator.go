package validation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"log"
)

func Validate(schemaDocument, document string) error {
	var schema map[string]interface{}
	err := compactAndUnmarshalJson(schemaDocument, &schema)
	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewGoLoader(schema)
	documentLoader := gojsonschema.NewStringLoader(document)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return errors.New(fmt.Sprintf("JSON schema validation failed: %v.", err))
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			log.Printf("Error %v.\n", desc)
		}
		return errors.New("JSON schema validation failed.")
	}

	return nil
}

func compactAndUnmarshalJson(rawJson string, v interface{}) error {
	buffer := new(bytes.Buffer)
	err := json.Compact(buffer, []byte(rawJson))
	if err != nil {
		return err
	}
	err = json.Unmarshal(buffer.Bytes(), v)
	if err != nil {
		return err
	}
	return nil
}
