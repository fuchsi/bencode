package bencode

import (
	"fmt"
	"math"
	"testing"
)

func TestEncodeSingleFileTorrent(t *testing.T) {
	dict := make(map[string]interface{})
	dict["announce"] = "http://tracker.archlinux.org:6969/announce"
	dict["comment"] = "Arch Linux 2017.11.01 (www.archlinux.org)"
	dict["created by"] = "mktorrent 1.1"
	dict["creation date"] = 1509525415
	info := make(map[string]interface{})
	info["length"] = 548405248
	info["name"] = "archlinux-2017.11.01-x86_64.iso"
	info["piece length"] = 524288
	info["pieces"] = ""
	dict["info"] = info

	result := string(Encode(dict))
	expected := "d8:announce42:http://tracker.archlinux.org:6969/announce7:comment41:Arch Linux 2017.11.01 (www.archlinux.org)10:created by13:mktorrent 1.113:creation datei1509525415e4:infod6:lengthi548405248e4:name31:archlinux-2017.11.01-x86_64.iso12:piece lengthi524288e6:pieces0:ee"

	if result != expected {
		t.Errorf("expected \"%s\" got \"%s\" instead", expected, result)
	}
}

func TestEncodeListOfIntegers(t *testing.T) {
	dict := make(map[string]interface{})
	var list []interface{}
	list = append(list, math.MinInt8)
	list = append(list, math.MaxInt8)
	list = append(list, math.MinInt16)
	list = append(list, math.MaxInt16)
	list = append(list, math.MinInt32)
	list = append(list, math.MaxInt32)
	list = append(list, math.MinInt64)
	list = append(list, math.MaxInt64)
	list = append(list, -1)
	list = append(list, 0)
	list = append(list, 1)
	dict["integers"] = list

	result := string(Encode(dict))
	expected := "d8:integersl"
	expected += fmt.Sprintf("i%de", math.MinInt8)
	expected += fmt.Sprintf("i%de", math.MaxInt8)
	expected += fmt.Sprintf("i%de", math.MinInt16)
	expected += fmt.Sprintf("i%de", math.MaxInt16)
	expected += fmt.Sprintf("i%de", math.MinInt32)
	expected += fmt.Sprintf("i%de", math.MaxInt32)
	expected += fmt.Sprintf("i%de", math.MinInt64)
	expected += fmt.Sprintf("i%de", math.MaxInt64)
	expected += fmt.Sprintf("i%de", -1)
	expected += fmt.Sprintf("i%de", 0)
	expected += fmt.Sprintf("i%de", 1)
	expected += "ee"

	if result != expected {
		t.Errorf("expected \"%s\" got \"%s\" instead", expected, result)
	}
}
