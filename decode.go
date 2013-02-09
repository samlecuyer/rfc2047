package rfc2047

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"strconv"
)

type Encoding int
type Charset int

// Represents encoding 'Q'
const Q  Encoding = 0

// Represents charset UTF-8
const Utf8 Charset = 0

type Decoder struct {
	e Encoding
	c Charset
}

func (d *Decoder) Decode(in *bufio.Reader, out *bytes.Buffer) (err error) {
	for {
		b, err := in.ReadByte()
		if err != nil { break }
		switch b {
		case '=':
			char := make([]byte,2)
			in.Read(char)
			code, _ := strconv.ParseUint(string(char), 16, 8)
			out.WriteByte(byte(code))
		case '_':
			err = out.WriteByte(' ')
		case ' ':
			return errors.New("Invalid char ` `")
		case '?':
			b, err = in.ReadByte()
			if b == '=' {
				return nil
			}
			return errors.New("Unexpected `?`")
		default:
			err = out.WriteByte(b)
		}
	}
	return err
}

func DecodeString(in string) (string, error) {
	return Decode(bufio.NewReader(bytes.NewBufferString(in)))
}

func Decode(input *bufio.Reader) (_ string, err error) {
	buf := bytes.NewBufferString("")
	for {
		b, err := input.ReadByte()
		if err != nil { break }
		if b != '=' {
			err = buf.WriteByte(b)
		} else {
			input.UnreadByte()
			charset, err := readCharset(input)
			if err != nil { return "", err }
			encoding, err := readEncoding(input)
			if err != nil { return "", err }
			(&Decoder{encoding, charset}).Decode(input, buf)
		}
	}
	if err == io.EOF {
		return buf.String(), nil
	}
	return buf.String(), err
}

func readCharset(input *bufio.Reader) (Charset, error) {
	start, err := input.ReadString('?')
	if err != nil { return -1, err }
	if start != "=?" {
		return 0, errors.New("Invalid encoded string start")
	}
	charset, err := input.ReadString('?') 
	if err != nil { return -1, err }
	if strings.EqualFold(charset, "UTF-8?") {
		return Utf8, nil
	}
	return -1, errors.New("Unknown characterset")
}
func readEncoding(input *bufio.Reader) (Encoding, error) {
	encoding, err := input.ReadString('?')
	if err != nil { return -1, err }
	if strings.EqualFold(encoding, "Q?") {
		return Q, nil
	}
	return -1, errors.New("Unknown encoding")
}
