package structtrace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type model struct {
	IgnoreStr    string `trace:"ignore"`
	Str          string
	StrPtr       *string
	Nil          *string
	Int          int
	IntPtr       *int
	StrCustomTag string `trace:"key=strCustomTag"`
	Struct       modelChildA
	StructTag    modelChildA `trace:"key=struct_custom_tag"`
	NestedIgnore modelChildB
}

type modelChildA struct {
	Str          string
	StrCustomTag string `trace:"key=nested_custom_tag"`
}

type modelChildB struct {
	Ignore string `trace:"ignore"`
}

type mockSpan struct {
	tags map[string]interface{}
}

func (t mockSpan) SetTag(key string, value interface{}) {
	t.tags[key] = value
}

func TestStructTrace(t *testing.T) {
	span := mockSpan{
		tags: map[string]interface{}{},
	}

	strPtr := "test string ptr"
	intPtr := 100

	value := model{
		IgnoreStr:    "ignored",
		Str:          "test string",
		StrPtr:       &strPtr,
		Nil:          nil,
		Int:          1,
		IntPtr:       &intPtr,
		StrCustomTag: "custom str",
		Struct: modelChildA{
			Str:          "value",
			StrCustomTag: "ct",
		},
		StructTag: modelChildA{
			Str:          "value",
			StrCustomTag: "ct",
		},
	}

	StructTrace(span, value)

	require.EqualValues(t, map[string]interface{}{
		"int":                                 int64(1),
		"int_ptr":                             int64(100),
		"nil":                                 "nil",
		"str":                                 "test string",
		"strCustomTag":                        "custom str",
		"str_ptr":                             "test string ptr",
		"struct.nested_custom_tag":            "ct",
		"struct.str":                          "value",
		"struct_custom_tag.nested_custom_tag": "ct",
		"struct_custom_tag.str":               "value",
	}, span.tags)
}
