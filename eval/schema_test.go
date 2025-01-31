package eval

import (
	"encoding/json"
	"fmt"
	"testing"

	"com.github/madbrain/gjq/lang"
)

func TestParser(t *testing.T) {
	content := "{ \"hello\": 10, \"foo\": \"bar\", \"bar\": [ 10 ], \"bool\": true }"
	var parsed any
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		t.Fatal("bad json")
	}

	result := InferSchema(parsed)
	if result == nil {
		t.Fatal("cannot get type from json")
	}
	switch result.(type) {
	case *lang.ObjectType:
		fmt.Printf("%+v\n", result)
	default:
		t.Fatal("expecting object type")
	}
}
