package jog

import (
	"io/ioutil"
	"testing"
)

type TestCase struct {
	path  *[]string
	value interface{}
}

func GetSample(t *testing.T) *Jog {
	contents, err := ioutil.ReadFile("test/sample.json")
	if err != nil {
		t.Fatalf("Couldn't open sample JSON: %v\n", err)
	}
	obj, err := New(string(contents))
	if err != nil {
		t.Fatalf("Couldn't parse sample JSON: %v\n", err)
	}
	return obj
}

func DoTests(t *testing.T, obj *Jog, cases []TestCase) {
	var err error
	for _, test := range cases {
		switch val := test.value.(type) {
		case int:
			var got int
			if test.path == nil {
				got, err = obj.GetInt()
			} else {
				got, err = obj.GetInt(*(test.path)...)
			}
			if err != nil {
				t.Fatalf("Couldn't get value: %v\n", err)
			}
			if got != val {
				t.Fatalf("Expected %v, got %v\n", val, got)
			}
		case uint:
			var got uint
			if test.path == nil {
				got, err = obj.GetUInt()
			} else {
				got, err = obj.GetUInt(*(test.path)...)
			}
			if err != nil {
				t.Fatalf("Couldn't get value: %v\n", err)
			}
			if got != val {
				t.Fatalf("Expected %v, got %v\n", val, got)
			}
		case bool:
			var got bool
			if test.path == nil {
				got, err = obj.GetBool()
			} else {
				got, err = obj.GetBool(*(test.path)...)
			}
			if err != nil {
				t.Fatalf("Couldn't get value: %v\n", err)
			}
			if got != val {
				t.Fatalf("Expected %v, got %v\n", val, got)
			}
		case float64:
			var got float64
			if test.path == nil {
				got, err = obj.GetFloat()
			} else {
				got, err = obj.GetFloat(*(test.path)...)
			}
			if err != nil {
				t.Fatalf("Couldn't get value: %v\n", err)
			}
			if got != val {
				t.Fatalf("Expected %v, got %v\n", val, got)
			}
		case string:
			var got string
			if test.path == nil {
				got, err = obj.GetString()
			} else {
				got, err = obj.GetString(*(test.path)...)
			}
			if err != nil {
				t.Fatalf("Couldn't get value: %v\n", err)
			}
			if got != val {
				t.Fatalf("Expected %v, got %v\n", val, got)
			}
		}
	}
}

func TestRootObject(t *testing.T) {
	cases := []TestCase{
		TestCase{&[]string{"index"}, int(0)},
		TestCase{&[]string{"index"}, uint(0)},
		TestCase{&[]string{"_id"}, "54c7fff8e3268528239d9cb1"},
		TestCase{&[]string{"guid"}, "b4940c5c-82ee-4f5e-bd02-f847fe2b9fc6"},
		//TestCase{&[]string{"isActive"}, true},
		TestCase{&[]string{"balance"}, "$1,750.21"},
		TestCase{&[]string{"registered"}, "2014-10-12T09:38:08 +07:00"},
		TestCase{&[]string{"latitude"}, float64(-59.816976)},
	}
	DoTests(t, GetSample(t), cases)
}

func TestNestedObject(t *testing.T) {
	cases := []TestCase{
		TestCase{&[]string{"details", "age"}, 36},
		TestCase{&[]string{"details", "eyeColor"}, "brown"},
		TestCase{&[]string{"details", "longitude"}, float64(102.563977)},
	}
	DoTests(t, GetSample(t), cases)
}

func TestArray(t *testing.T) {
	obj := GetSample(t)
	tags, err := obj.GetArray("tags")
	if err != nil {
		t.Fatalf("Couldn't get nested array: %v\n", err)
	}

	values := []string{"nisi", "sint", "aute", "tempor", "sit", "esse", "in"}
	if len(values) != len(tags) {
		t.Fatalf("Got length %d, expected %d\n", len(tags), len(values))
	}
	for i, friend := range tags {
		DoTests(t, friend, []TestCase{TestCase{nil, values[i]}})
	}
}
