package bencode

import (
	"reflect"
	"testing"
)

func TestHandleString(t *testing.T) {
	t.Log("Testing handle String")
	s := "7:example5:hello"
	globalIndex := 0
	t.Log("s is", s)
	resultString, newGlobalIndex := handleString([]byte(s), globalIndex)
	if resultString != "example" || newGlobalIndex != 9 {
		t.Error("Expected: example, example, got", resultString, " expected: 9 as INDEX, got", newGlobalIndex)
	}
}

func TestHandleInteger(t *testing.T) {
	t.Log("Testing handle Integer")
	s := "i75234e5:hello"
	globalIndex := 0
	result, newGlobalIndex := handleInt([]byte(s), globalIndex)
	if result != 75234 || newGlobalIndex != 7 {
		t.Error("Expected 75234 got", result, " as result expected: 7 got", newGlobalIndex, " as global index")
	}
}

func TestHandleList(t *testing.T) {
	t.Log("Testing handle List")
	s := "li75234e5:helloe"
	result, newGlobalIndex := handleList([]byte(s), 0)
	if result[0].(int64) != 75234 || result[1] != "hello" || newGlobalIndex != 16 {
		t.Error("Expected 75234 GOT ", result[0], " as first member, as second member expected hello GOT ", result[1], " as global index expected 16 GOT", newGlobalIndex)
	}
}

func TestHandleDictionary(t *testing.T) {
	t.Log("Testing handle Dictionary")
	s := "d4:listli75234e5:helloe1:a1:be"
	value, globalIndex := handleDictionary([]byte(s), 0)
	subList, _ := handleList([]byte("li75234e5:helloe"), 0)
	if !compareSlices(value["list"].([]interface{}), subList) {
		t.Error("Expected:", subList, "got:", value["list"])
	}
	if value["a"] != "b" || globalIndex != len(s) {
		t.Error("Expected value with key a to be b and globalIndex to be", len(s), " instead got ", value["a"], globalIndex)
	}
}

func compareSlices(a, b []interface{}) bool {
	return reflect.DeepEqual(a, b)
}
