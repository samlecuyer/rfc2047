package rfc2047

import (
	"testing"
	"bytes"
	"bufio"
)

type decodeTest struct {
	input    string
	expected string
}

func (test decodeTest) Run(t *testing.T) {
	r := bufio.NewReader(bytes.NewBufferString(test.input))
	ps, err := Decode(r)
	if err != nil {
		t.Fatalf("parsing %s: %s", test.input, err)
	}

	if ps != test.expected {
		t.Fatalf("Expected, Actual(%#v, %#v)", test.expected, ps)
	}
}

func TestDecode_InvalidSpaces(t *testing.T) {
	_, err := Decode(
		bufio.NewReader(
			bytes.NewBufferString("=?utf-8?q?this is some text?=")))
	if err == nil {
		t.Fatalf("Should not have decoded text with illegal spaces")
	}
}

func TestDecode_HappyPath(t *testing.T) {
	decodeTest {
		"this is some text",
		"this is some text",
	}.Run(t)
}


func TestDecode_HappyPathUtf8(t *testing.T) {
	decodeTest {
		"=?utf-8?q?this_is_some_text?=",
		"this is some text",
	}.Run(t)
}

func TestDecode_ValidUtf8(t *testing.T) {
	decodeTest {
		"=?UTF-8?Q?You_are_now_friends_with_K=C3=A4rsten_Jaynes?=",
		"You are now friends with Kärsten Jaynes",
	}.Run(t)
}
