package gomd

import (
	"testing"
)

func TestOne(t *testing.T) {
	b := Builder{}
	original := []*Element{
		b.H1("Header"),
		b.UL(
			b.Textln("hi"),
		),
	}
	// Generate should trim nils
	trimmed := []*Element{}
	for _, el := range original {
		if el != nil {
			trimmed = append(trimmed, el)
		}
	}

	result := ParseMD(b.Generate(trimmed...), "")

	deepEquals(t, trimmed, result, 0)
}

func deepEquals(t *testing.T, original, result []*Element, depth int) bool {
	if len(original) != len(result) {
		t.Errorf("unequal amount of members at depth %d: original: %d result: %d", depth, len(original), len(result))
		return false
	}
	ok := true
	// loop through all and check
	for i := 0; i < len(original); i++ {
		if original[i].name != result[i].name {
			t.Errorf("name missmatch at depth %d, i %d: original: %s result: %s ", depth, i, original[i].name, result[i].name)
			return false
		}
		if original[i].content != result[i].content {
			t.Errorf("content missmatch at depth %d, i %d: original %s result %s", depth, i, original[i].content, result[i].content)
			return false
		}
		isEqual := deepEquals(t, original[i].children, result[i].children, depth+1)
		if !isEqual {
			ok = false
			break
		}
	}
	return ok
}
