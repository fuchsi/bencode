/*
 * Copyright (c) 2017 Daniel Müller
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package bencode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type decoder struct {
	bufio.Reader
	offset uint64
}

func (d *decoder) ReadByte() (byte, error) {
	p, err := d.Reader.ReadByte()
	d.offset++
	if err == nil {

		return p, nil
	}

	return 0, err
}

func (d *decoder) UnreadByte() (error) {
	d.offset--
	return d.Reader.UnreadByte()
}

func (d *decoder) decode() (interface{}, error) {
	b, err := d.ReadByte()
	if err != nil {
		return nil, err
	}
	switch b {
	case 'd':
		return d.readDictionary()
	case 'l':
		return d.readList()
	case 'i':
		return d.readInteger()
	}

	if b >= '0' && b <= '9' {
		d.Reader.UnreadByte()
		return d.readString()
	}

	return nil, errors.New("bencode: invalid type " + fmt.Sprintf("'%c' at offset 0x%X", b, d.offset))
}

func (d *decoder) readDictionary() (map[string]interface{}, error) {
	var dict map[string]interface{}
	dict = make(map[string]interface{})

	b, err := d.ReadByte()
	if err != nil {
		return nil, err
	}
	d.Reader.UnreadByte()

	for b != 'e' {
		key, err := d.readString()
		if err != nil {
			return nil, err
		}
		value, err := d.decode()
		if err != nil {
			return nil, err
		}
		dict[key] = value
		b, err = d.ReadByte()
		if err != nil {
			return nil, err
		}
		if b != 'e' {
			d.Reader.UnreadByte()
		}
	}

	return dict, nil
}

func (d *decoder) readList() ([]interface{}, error) {
	var l []interface{}
	l = make([]interface{}, 0)

	b, err := d.ReadByte()
	if err != nil {
		return nil, err
	}
	d.Reader.UnreadByte()

	for b != 'e' {
		value, err := d.decode()
		if err != nil {
			return nil, err
		}
		l = append(l, value)
		b, err = d.ReadByte()
		if err != nil {
			return nil, err
		}
		if b != 'e' {
			d.Reader.UnreadByte()
		}
	}

	return l, nil
}

func (d *decoder) readString() (string, error) {
	size, err := d.readIntegerUntilDelimiter(':')

	var strLen int64
	var ok bool
	if strLen, ok = size.(int64); !ok {
		return "", errors.New("string length may not exceed the size of int64")
	}
	if strLen < 0 {
		return "", errors.New("string length can not be a negative number")
	}

	buf := make([]byte, strLen)
	_, err = io.ReadFull(d, buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (d *decoder) readInteger() (interface{}, error) {
	return d.readIntegerUntilDelimiter('e')
}

func (d *decoder) readIntegerUntilDelimiter(delim byte) (interface{}, error) {
	var buf []byte
	b, err := d.ReadByte()
	if err != nil {
		return 0, err
	}
	if !(('0' <= b && b <= '9') || b == '-') { // b must be a number or '-'
		return 0, errors.New("bencode: integers must begin with a digit or negative sign")
	}
	for b != delim {
		if !(('0' <= b && b <= '9') || b == '-') { // b must be a number or '-'
			return 0, errors.New("bencode: invalid character for integer")
		}
		buf = append(buf, b)

		if len(buf) == 2 { // check for invalid leading zeroes and negative zeroes
			if buf[0] == '0' && b != '0' {
				return 0, errors.New("bencode: invalid leading zero in integer")
			}
			if buf[0] == '-' && b == '0' {
				return 0, errors.New("bencode: invalid negative zero")
			}
		}

		b, err = d.ReadByte()
		if err != nil {
			return 0, err
		}
	}

	if data, err := strconv.ParseInt(fmt.Sprintf("%s", buf), 10, 64); err == nil {
		return data, nil
	} else if data, err := strconv.ParseUint(fmt.Sprintf("%s", buf), 10, 64); err == nil {
		return data, nil
	}

	return 0, err
}

func Decode(reader io.Reader) (map[string]interface{}, error) {
	decoder := decoder{*bufio.NewReader(reader), 0}

	if firstByte, err := decoder.ReadByte(); err != nil {
		return make(map[string]interface{}), nil
	} else if firstByte != 'd' {
		return nil, errors.New("bencode: data must begin with a dictionary")
	}

	return decoder.readDictionary()
}
