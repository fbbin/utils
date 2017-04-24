package utils

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func Ip2Long(ip string) (uint32, error) {
	ipSplit := strings.Split(ip, ".")
	if len(ipSplit) != 4 {
		return 0, errors.New("IP格式错误。")
	}
	var intIp uint32
	for key, value := range ipSplit {
		i, err := strconv.Atoi(value)
		if err != nil || i > 255 {
			return 0, errors.New("IP格式错误。")
		}
		intIp = intIp | uint32(i)<<uint((3-key)*8)
	}
	return intIp, nil
}

func Long2Ip(intIp uint32) (string, error) {
	ip4 := intIp >> 0 & 255
	ip3 := intIp >> 8 & 255
	ip2 := intIp >> 16 & 255
	ip1 := intIp >> 24 & 255
	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return "", errors.New("不是一个整型的IP数据")
	}
	ipString := fmt.Sprintf("%d.%d.%d.%d", ip1, ip2, ip3, ip4)
	return ipString, nil
}

func Int64ToByte(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func ByteToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func Int32ToByte(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func ByteToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

func ByteToInt(buf []byte) int {
	value, _ := strconv.Atoi(string(buf))
	return value
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func Exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ReadFileByte(path string) ([]byte, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	} else {
		defer fi.Close()
		return ioutil.ReadAll(fi)
	}
}

func ReadFileStr(path string) (string, error) {
	raw, err := ReadFileByte(path)
	return string(raw), err
}

// url encode string, is + not %20
func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

// url decode string
func UrlDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// base64 encode
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// base64 decode
func Base64Decode(str string) (string, error) {
	s, e := base64.StdEncoding.DecodeString(str)
	return string(s), e
}
