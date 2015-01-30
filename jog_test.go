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
				if obj.Type() != TypeNumber {
					t.Fatalf("Expected TypeNumber, got %v\n", obj.Type())
				}
				got, err = obj.GetInt()
			} else {
				if obj.Type(*(test.path)...) != TypeNumber {
					t.Fatalf("Expected TypeNumber, got %v\n", obj.Type(*(test.path)...))
				}
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
				if obj.Type() != TypeNumber {
					t.Fatalf("Expected TypeNumber, got %v\n", obj.Type())
				}
				got, err = obj.GetUInt()
			} else {
				if obj.Type(*(test.path)...) != TypeNumber {
					t.Fatalf("Expected TypeNumber, got %v\n", obj.Type(*(test.path)...))
				}
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
				if obj.Type() != TypeBool {
					t.Fatalf("Expected TypeBool, got %v\n", obj.Type())
				}
				got, err = obj.GetBool()
			} else {
				if obj.Type(*(test.path)...) != TypeBool {
					t.Fatalf("Expected TypeBool, got %v\n", obj.Type(*(test.path)...))
				}
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
				if obj.Type() != TypeNumber {
					t.Fatalf("Expected TypeNumber, got %v\n", obj.Type())
				}
				got, err = obj.GetFloat()
			} else {
				if obj.Type(*(test.path)...) != TypeNumber {
					t.Fatalf("Expected TypeNumber, got %v\n", obj.Type(*(test.path)...))
				}
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
				if obj.Type() != TypeString {
					t.Fatalf("Expected TypeString, got %v\n", obj.Type())
				}
				got, err = obj.GetString()
			} else {
				if obj.Type(*(test.path)...) != TypeString {
					t.Fatalf("Expected TypeString, got %v\n", obj.Type(*(test.path)...))
				}
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
		TestCase{&[]string{"isActive"}, true},
		TestCase{&[]string{"balance"}, "$1,750.21"},
		TestCase{&[]string{"registered"}, "2014-10-12T09:38:08 +07:00"},
		TestCase{&[]string{"latitude"}, float64(-59.816976)},
	}
	DoTests(t, GetSample(t), cases)
}

func TestGetObject(t *testing.T) {
	o := GetSample(t)
	obj, err := o.GetObject()
	if err != nil {
		t.Fatalf("Expected no errors, got %v\n", err)
	}
	if len(obj) != 10 {
		t.Fatalf("Expected 10 keys in root object, got %d\n", len(obj))
	}
	details, _ := o.GetObject("details")
	if len(details) != 3 {
		t.Fatalf("Expected 3 keys in details object, got %d\n", len(details))
	}
	age, _ := details["age"].GetInt()
	if age != 36 {
		t.Fatalf("Expected details/age to be 36, got %d\n", age)
	}
	if details["longitude"].Type() != TypeNumber {
		t.Fatalf("Expected details/longitude to be a number\n")
	}
}

func TestNestedObject(t *testing.T) {
	cases := []TestCase{
		TestCase{&[]string{"details", "age"}, 36},
		TestCase{&[]string{"details", "eyeColor"}, "brown"},
		TestCase{&[]string{"details", "longitude"}, float64(102.563977)},
	}
	DoTests(t, GetSample(t), cases)
}

func TestNestedArray(t *testing.T) {
	obj := GetSample(t)

	if obj.Type("tags") != TypeArray {
		t.Fatalf("Expected TypeArray, got %v\n", obj.Type("tags"))
	}

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

func TestNestedArrayObject(t *testing.T) {
	obj := GetSample(t)
	friends, err := obj.GetArray("friends")
	if err != nil {
		t.Fatalf("Couldn't get nested array: %v\n", err)
	}

	values := [][]TestCase{
		[]TestCase{
			TestCase{&[]string{"id"}, 0},
			TestCase{&[]string{"name"}, "Case Gross"},
		},
		[]TestCase{
			TestCase{&[]string{"id"}, 1},
			TestCase{&[]string{"name"}, "Gilbert Rasmussen"},
		},
		[]TestCase{
			TestCase{&[]string{"id"}, 2},
			TestCase{&[]string{"name"}, "Harris Huff"},
		},
	}
	if len(friends) != len(values) {
		t.Fatalf("Got length %d, expected %d\n", len(friends), len(values))
	}
	for i, v := range values {
		DoTests(t, friends[i], v)
	}
}
