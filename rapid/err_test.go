package rapid

import (
	"testing"
)

func DoErrorTest(t *testing.T, input string, output string) {
	_, err := New(input)
	if err == nil || err.Error() != output {
		t.Fatalf("Expected '%v', got '%v'\n", output, err)
	}
}

func TestEmpty(t *testing.T) {
	DoErrorTest(t, "", "[0] The document is empty.")
}

func TestNonSimple(t *testing.T) {
	DoErrorTest(t, "{\"foo\":2}10", "[9] The document root must not follow by other values.")
}

func TestInvalid(t *testing.T) {
	DoErrorTest(t, "nulL", "[3] Invalid value.")
}

func TestMissingName(t *testing.T) {
	DoErrorTest(t, "{:3.14}", "[1] Missing a name for object member.")
	DoErrorTest(t, "{null:1}", "[1] Missing a name for object member.")
}

func TestMissingColon(t *testing.T) {
	DoErrorTest(t, "{\"name\"\"jog\"}", "[8] Missing a colon after a name of object member.")
	DoErrorTest(t, "{\"name\",\"jog\"}", "[8] Missing a colon after a name of object member.")
}

func TestMissingCommaObject(t *testing.T) {
	DoErrorTest(t, "{\"name\":\"jog\"\"foo\":\"bar\"}", "[14] Missing a comma or '}' after an object member.")
}

func TestMissingCommaArray(t *testing.T) {
	DoErrorTest(t, "[{\"name\":\"jog\"}{\"foo\":\"bar\"}]", "[16] Missing a comma or ']' after an array element.")
}

func TestInvalidUnicode(t *testing.T) {
	DoErrorTest(t, "[\"\\uABCG\"]", "[7] Incorrect hex digit after \\u escape in string.")
}

func TestInvalidSurrogate(t *testing.T) {
	DoErrorTest(t, "[\"\\uD800X\"]", "[7] The surrogate pair in string is invalid.")
	DoErrorTest(t, "[\"\\uD800\\uFFFF\"]", "[12] The surrogate pair in string is invalid.")
}

func TestInvalidEscape(t *testing.T) {
	DoErrorTest(t, "[\"\\a\"]", "[3] Invalid escape character in string.")
}

func TestMissingQuotation(t *testing.T) {
	DoErrorTest(t, "[\"Test]", "[6] Missing a closing quotation mark in string.")
}

func TestNumberTooBig(t *testing.T) {
	bigNumber := "1"
	for i := 0; i < 310; i++ {
		bigNumber += "0"
	}
	DoErrorTest(t, bigNumber, "[309] Number too big to be stored in double.")
}

func TestMissingFraction(t *testing.T) {
	DoErrorTest(t, "1.", "[2] Miss fraction part in number.")
	DoErrorTest(t, "1.a", "[2] Miss fraction part in number.")
}

func TestMissingExponent(t *testing.T) {
	DoErrorTest(t, "1e", "[2] Miss exponent in number.")
	DoErrorTest(t, "1e_", "[2] Miss exponent in number.")
}
