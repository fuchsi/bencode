package bencode

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func TestDecodeSingleFileTorrent(t *testing.T) {
	str := "d8:announce42:http://tracker.archlinux.org:6969/announce7:comment41:Arch Linux 2017.11.01 (www.archlinux.org)10:created by13:mktorrent 1.113:creation datei1509525415e4:infod6:lengthi548405248e4:name31:archlinux-2017.11.01-x86_64.iso12:piece lengthi524288e6:pieces0:ee"
	buf := bytes.NewBufferString(str)

	dict, err := Decode(buf)
	if err != nil {
		t.Fatal(err)
	}

	if dict["announce"] != "http://tracker.archlinux.org:6969/announce" {
		t.Error("announce mismatch")
	}
	if dict["comment"] != "Arch Linux 2017.11.01 (www.archlinux.org)" {
		t.Error("comment mismatch")
	}
	if dict["created by"] != "mktorrent 1.1" {
		t.Error("created by mismatch")
	}
	if dict["creation date"] != int64(1509525415) {
		t.Error("creation date mismatch")
	}

	info := dict["info"].(map[string]interface{})
	if info["length"] != int64(548405248) {
		t.Error("info.length mismatch")
	}
	if info["name"] != "archlinux-2017.11.01-x86_64.iso" {
		t.Error("info.name mismatch")
	}
	if info["piece length"] != int64(524288) {
		t.Error("info.piece length mismatch")
	}

	encoded := string(Encode(dict))
	if encoded != str {
		t.Error("decode(str).encode != str")
	}
}

func TestDecodeListOfIntegers(t *testing.T) {
	values := []int64{
		math.MinInt8,
		math.MaxInt8,
		math.MinInt16,
		math.MaxInt16,
		math.MinInt32,
		math.MaxInt32,
		math.MinInt64,
		math.MaxInt64,
		-1,
		0,
		1,
	}

	str := fmt.Sprintf("d8:integersli%dei%dei%dei%dei%dei%dei%dei%dei%dei%dei%deee",
			values[0], values[1], values[2], values[3], values[4], values[5],
			values[6], values[7], values[8], values[9], values[10])

	buf := bytes.NewBufferString(str)

	dict, err := Decode(buf)
	if err != nil {
		t.Fatal(err)
	}

	intList := dict["integers"].([]interface{})
	length := len(intList)
	if length != len(values) {
		t.Error("length mismatch")
	}

	for i := 0; i < length; i++ {
		if intList[i] != values[i] {
			t.Error("mismatch at index", i)
		}
	}

	encoded := string(Encode(dict))
	if encoded != str {
		t.Error("decode(str).encode != str")
	}
}

func TestDecodeUint64(t *testing.T) {
	values := []uint64{
		uint64(math.MaxInt64 + 1),
		math.MaxUint64,
	}

	str := fmt.Sprintf("d8:integersli%dei%deee",
		values[0], values[1])

	buf := bytes.NewBufferString(str)

	dict, err := Decode(buf)
	if err != nil {
		t.Fatal(err)
	}

	intList := dict["integers"].([]interface{})
	length := len(intList)
	if length != len(values) {
		t.Error("length mismatch")
	}

	for i := 0; i < length; i++ {
		if intList[i] != values[i] {
			t.Error("mismatch at index", i)
		}
	}

	encoded := string(Encode(dict))
	if encoded != str {
		t.Error("decode(str).encode != str")
	}
}

func TestDecodeNegativeString(t *testing.T) {
	str := "d3:key-1:"
	buf := bytes.NewBufferString(str)

	_, err := Decode(buf)
	if err.Error() != "string length can not be a negative number" {
	}
}

func TestDecodeUint64StringLength(t *testing.T) {
	str := fmt.Sprintf("d3:key%d:", uint64(math.MaxInt64 + 1337))
	buf := bytes.NewBufferString(str)

	_, err := Decode(buf)
	if err.Error() != "string length may not exceed the size of int64" {
	}
}