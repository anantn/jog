package test

import (
	"testing"

	"github.com/anantn/jog"
	"github.com/anantn/jog/rapid"
	"github.com/anantn/jog/yajl"
)

var SAMPLE = `{"index":0,"_id":"54c7fff8e3268528239d9cb1","guid":"b4940c5c-82ee-4f5e-bd02-f847fe2b9fc6","isActive":true,"balance":"$1,750.21","details":{"age":36,"eyeColor":"brown","longitude":102.563977},"registered":"2014-10-12T09:38:08 +07:00","latitude":-59.816976,"tags":["nisi","sint","aute","tempor","sit","esse","in"],"friends":[{"id":0,"name":"Case Gross"},{"id":1,"name":"Gilbert Rasmussen"},{"id":2,"name":"Harris Huff"}]}`

type TestCase struct {
	path  *[]string
	value interface{}
}

func GetSamples(t *testing.T) []jog.Value {
	robj, err := rapid.New(SAMPLE)
	if err != nil {
		t.Fatalf("Couldn't parse sample JSON with Rapid: %v\n", err)
	}
	yobj, err := yajl.New(SAMPLE)
	if err != nil {
		t.Fatalf("Couldn't parse sample JSON with YAJL: %v\n", err)
	}
	return []jog.Value{robj, yobj}
}

func DoTests(t *testing.T, objs []jog.Value, cases []TestCase) {
	var err error
	for _, obj := range objs {
		for _, test := range cases {
			switch val := test.value.(type) {
			case int:
				var got int
				if test.path == nil {
					if obj.Type() != jog.TypeNumber {
						t.Fatalf("Expected jog.TypeNumber, got %v\n", obj.Type())
					}
					got, err = obj.GetInt()
				} else {
					if obj.Type(*(test.path)...) != jog.TypeNumber {
						t.Fatalf("Expected jog.TypeNumber, got %v\n", obj.Type(*(test.path)...))
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
					if obj.Type() != jog.TypeNumber {
						t.Fatalf("Expected jog.TypeNumber, got %v\n", obj.Type())
					}
					got, err = obj.GetUInt()
				} else {
					if obj.Type(*(test.path)...) != jog.TypeNumber {
						t.Fatalf("Expected jog.TypeNumber, got %v\n", obj.Type(*(test.path)...))
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
					if obj.Type() != jog.TypeBool {
						t.Fatalf("Expected jog.TypeBool, got %v\n", obj.Type())
					}
					got, err = obj.GetBool()
				} else {
					if obj.Type(*(test.path)...) != jog.TypeBool {
						t.Fatalf("Expected jog.TypeBool, got %v\n", obj.Type(*(test.path)...))
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
					if obj.Type() != jog.TypeNumber {
						t.Fatalf("Expected jog.TypeNumber, got %v\n", obj.Type())
					}
					got, err = obj.GetFloat()
				} else {
					if obj.Type(*(test.path)...) != jog.TypeNumber {
						t.Fatalf("Expected jog.TypeNumber, got %v\n", obj.Type(*(test.path)...))
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
					if obj.Type() != jog.TypeString {
						t.Fatalf("Expected jog.TypeString, got %v\n", obj.Type())
					}
					got, err = obj.GetString()
				} else {
					if obj.Type(*(test.path)...) != jog.TypeString {
						t.Fatalf("Expected jog.TypeString, got %v\n", obj.Type(*(test.path)...))
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
	DoTests(t, GetSamples(t), cases)
}

func TestGetObject(t *testing.T) {
	for _, o := range GetSamples(t) {
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
		if details["longitude"].Type() != jog.TypeNumber {
			t.Fatalf("Expected details/longitude to be a number\n")
		}
	}
}

func TestNestedObject(t *testing.T) {
	cases := []TestCase{
		TestCase{&[]string{"details", "age"}, 36},
		TestCase{&[]string{"details", "eyeColor"}, "brown"},
		TestCase{&[]string{"details", "longitude"}, float64(102.563977)},
	}
	DoTests(t, GetSamples(t), cases)
}

func TestNestedArray(t *testing.T) {
	for _, obj := range GetSamples(t) {
		if obj.Type("tags") != jog.TypeArray {
			t.Fatalf("Expected jog.TypeArray, got %v\n", obj.Type("tags"))
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
			DoTests(t, []jog.Value{friend}, []TestCase{TestCase{nil, values[i]}})
		}
	}
}

func TestNestedArrayObject(t *testing.T) {
	for _, obj := range GetSamples(t) {
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
			DoTests(t, []jog.Value{friends[i]}, v)
		}
	}
}

func TestStringify(t *testing.T) {
	for _, obj := range GetSamples(t) {
		val, err := obj.Stringify()
		if err != nil {
			t.Fatalf("Did not stringify object correctly: %v\n", err)
		}
		if val != SAMPLE {
			t.Fatalf("Did not stringify object correctly: %v\n", val)
		}
		tags, _ := obj.Stringify("tags")
		if tags != `["nisi","sint","aute","tempor","sit","esse","in"]` {
			t.Fatalf("Did not stringify array correctly! %#v", tags)
		}
	}
}
