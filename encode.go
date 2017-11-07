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
	"bytes"
	"reflect"
	"sort"
	"strconv"
)

type encoder struct {
	bytes.Buffer
}

func (encoder *encoder) writeDictionary(d map[string]interface{}) {
	// Sort keys
	list := make(sort.StringSlice, len(d))
	i := 0
	for key := range d {
		list[i] = key
		i++
	}
	list.Sort()

	encoder.WriteByte('d')
	for _, key := range list {
		encoder.writeString(key)           // Key
		encoder.writeInterfaceType(d[key]) // Value
	}
	encoder.WriteByte('e')
}

func (encoder *encoder) writeList(l []interface{}) {
	encoder.WriteByte('l')

	for _, value := range l {
		encoder.writeInterfaceType(value)
	}

	encoder.WriteByte('e')
}

func (encoder *encoder) writeString(str string) {
	encoder.WriteString(strconv.Itoa(len(str)))
	encoder.WriteByte(':')
	encoder.WriteString(str)
}

func (encoder *encoder) writeInteger(i int64) {
	encoder.WriteByte('i')
	encoder.WriteString(strconv.FormatInt(i, 10))
	encoder.WriteByte('e')
}

func (encoder *encoder) writeUInteger(i uint64) {
	encoder.WriteByte('i')
	encoder.WriteString(strconv.FormatUint(i, 10))
	encoder.WriteByte('e')
}

func (encoder *encoder) writeInterfaceType(v interface{}) {
	switch v := v.(type) {
	case string:
		encoder.writeString(v)
	case []interface{}:
		encoder.writeList(v)
	case map[string]interface{}:
		encoder.writeDictionary(v)
	case int, int8, int16, int32, int64:
		encoder.writeInteger(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		encoder.writeUInteger(reflect.ValueOf(v).Uint())
	}
}

func Encode(dict map[string]interface{}) []byte {
	encoder := encoder{}
	encoder.writeInterfaceType(dict)
	return encoder.Bytes()
}
