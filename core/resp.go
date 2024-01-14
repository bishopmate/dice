package core

import (
	"errors"
	"strconv"
)

// reads the length typically the first integer of the string
// until hit by an non-digit byte and returns
// the integer and the delta = length + 2 (CRLF)
func readLength(data []byte) (int, int) {
	length := 0
	for i := range data {
		if data[i] == '\r' {
			return length, i + 2
		}
		cur, _ := strconv.Atoi(string(data[i]))
		length = (length * 10) + (cur - 0)
	}
	return 0, 0
}

/*
Reads a Simple string encoded in RESP and returns the string, the next position in data after the RESP simple string ends, and the error
Simple strings are encoded as a plus (+) character, followed by a string. The string mustn't
contain a CR (\r) or LF (\n) character and is terminated by CRLF (i.e., \r\n).
*/
func readSimpleString(data []byte) (string, int, error) {

	pos := 1
	for ; data[pos] != '\r'; pos++ {

	}

	return string(data[1:pos]), pos + 2, nil
}

// reads a RESP encoded error from data and returns
// the error string, the delta, and the error
func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

// reads a RESP encoded integer from data and returns
// the intger value, the delta, and the error
func readInt64(data []byte) (int64, int, error) {
	// first character :
	pos := 1
	var value int64 = 0

	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}

	return value, pos + 2, nil
}

// reads a RESP encoded string from data and returns
// the string, the delta, and the error
func readBulkString(data []byte) (string, int, error) {

	pos := 1

	length, delta := readLength(data[pos:])
	pos += delta

	return string(data[pos : pos+length]), pos + length + 2, nil
}

// reads a RESP encoded array from data and returns
// the array, the delta, and the error
func readArray(data []byte) (interface{}, int, error) {
	pos := 1
	numberOfElements, delta := readLength(data[pos:])
	pos += delta

	var elements []interface{} = make([]interface{}, numberOfElements)

	for i := range elements {
		element, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elements[i] = element
		pos += delta
	}

	return elements, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}

	return nil, 0, nil
}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	value, _, err := DecodeOne(data)
	return value, err
}
