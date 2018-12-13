package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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

// IP转成字符串 大端解析
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

func Long2IpBig(intIp int) (string, bool) {
	i1 := intIp & 255
	i2 := intIp >> 8 & 255
	i3 := intIp >> 16 & 255
	i4 := intIp >> 24 & 255
	if i1 > 255 || i2 > 255 || i3 > 255 || i4 > 255 {
		return "", false
	}
	ipstring := fmt.Sprintf("%d.%d.%d.%d", i1, i2, i3, i4)

	return ipstring, true
}

// IP转成字符串 (小端解析)
func Long2IpLittle(intIp int) (string, bool) {
	i1 := intIp & 255
	i2 := intIp >> 8 & 255
	i3 := intIp >> 16 & 255
	i4 := intIp >> 24 & 255
	if i1 > 255 || i2 > 255 || i3 > 255 || i4 > 255 {
		return "", false
	}
	ipString := fmt.Sprintf("%d.%d.%d.%d", i4, i3, i2, i1)

	return ipString, true
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

// 截取字符串
func subStr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func hashmd5(signString string) string {
	hash := md5.New()
	hash.Write([]byte(signString))
	cipherStr := hash.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// Format 跟 PHP 中 date 类似的使用方式，如果 ts 没传递，则使用当前时间
func Format(format string, ts ...time.Time) string {
	patterns := []string{
		// 年
		"Y", "2006", // 4 位数字完整表示的年份
		"y", "06", // 2 位数字表示的年份

		// 月
		"m", "01", // 数字表示的月份，有前导零
		"n", "1", // 数字表示的月份，没有前导零
		"M", "Jan", // 三个字母缩写表示的月份
		"F", "January", // 月份，完整的文本格式，例如 January 或者 March

		// 日
		"d", "02", // 月份中的第几天，有前导零的 2 位数字
		"j", "2", // 月份中的第几天，没有前导零

		"D", "Mon", // 星期几，文本表示，3 个字母
		"l", "Monday", // 星期几，完整的文本格式;L的小写字母

		// 时间
		"g", "3", // 小时，12 小时格式，没有前导零
		"G", "15", // 小时，24 小时格式，没有前导零
		"h", "03", // 小时，12 小时格式，有前导零
		"H", "15", // 小时，24 小时格式，有前导零

		"a", "pm", // 小写的上午和下午值
		"A", "PM", // 小写的上午和下午值

		"i", "04", // 有前导零的分钟数
		"s", "05", // 秒数，有前导零
	}
	replacer := strings.NewReplacer(patterns...)
	format = replacer.Replace(format)

	t := time.Now()
	if len(ts) > 0 {
		t = ts[0]
	}
	return t.Format(format)
}

func StrToLocalTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	zoneName, offset := time.Now().Zone()

	zoneValue := offset / 3600 * 100
	if zoneValue > 0 {
		value += fmt.Sprintf(" +%04d", zoneValue)
	} else {
		value += fmt.Sprintf(" -%04d", zoneValue)
	}

	if zoneName != "" {
		value += " " + zoneName
	}
	return StrToTime(value)
}

func StrToTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	layouts := []string{
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05 -0700 MST",
		"2006/01/02 15:04:05 -0700",
		"2006/01/02 15:04:05",
		"2006-01-02 -0700 MST",
		"2006-01-02 -0700",
		"2006-01-02",
		"2006/01/02 -0700 MST",
		"2006/01/02 -0700",
		"2006/01/02",
		"2006-01-02 15:04:05 -0700 -0700",
		"2006/01/02 15:04:05 -0700 -0700",
		"2006-01-02 -0700 -0700",
		"2006/01/02 -0700 -0700",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, value)
		if err == nil {
			return t
		}
	}
	panic(err)
}

func Shuffle(vals []interface{}) []interface{} {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]interface{}, len(vals))
	for i, randIndex := range r.Perm(len(vals)) {
		ret[i] = vals[randIndex]
	}
	return ret
}

// 将内容转码成GBK
func Encode(src string) string {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		return ""
	}
	return string(data)
}

// 将内容转码成UTF-8
func Decode(src string) string {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return ""
	}
	return string(data)
}
