package util

import (

	"github.com/satori/go.uuid"
	"strings"
	"time"
	"strconv"
	"sort"
	"crypto/md5"
	"fmt"
	"crypto/sha1"
	mr "math/rand"

)



//生成uuid
func GenerateUuid() string {
	uid := uuid.NewV1()
	uids := strings.Split(uid.String(), "-")
	return uids[0] + uids[1] + uids[2] + uids[4] + uids[3]
}

func GenerateOrderNumber()string{
	it :=time.Now().Unix()

	itStr := strconv.FormatInt(it,10)

	return itStr
}

// 生成签名
func GenerateSign(paramsMap map[string]string, paramsList []string, key string) string {
	paramsStr := ""

	// 首先进行字典序排序
	sort.Strings(paramsList)

	for _, param := range paramsList {
		if paramsMap[param] != "" {
			if paramsStr == "" {
				paramsStr = param + "=" + paramsMap[param]
			} else {
				paramsStr = paramsStr + "&" + param + "=" + paramsMap[param]
			}
		}
	}

	if key != "" {
		paramsStr += "&key=" + key
	}

	return strings.ToUpper(Md5(paramsStr))
}

// 生成XML串
func GenerateXMLStr(params map[string]string) string {
	result := "<xml>"

	for k, v := range params {
		result = result + "<" + k + "><![CDATA[" + v + "]]></" + k + ">"
	}

	result += "</xml>"
	return result
}
// 对串进行md5加密
func Md5(str string) string {
	result := ""
	if str == "" {
		return result
	}
	b := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", b)
}



func GeneratePageSign(paramsMap map[string]string, paramsList []string) string {
	paramsStr := ""

	// 首先进行字典序排序
	sort.Strings(paramsList)

	for _, param := range paramsList {
		if paramsStr == "" {
			paramsStr = param + "=" + paramsMap[param]
		} else {
			paramsStr = paramsStr + "&" + param + "=" + paramsMap[param]
		}
	}

	return Sha1(paramsStr)
}

// 对串进行sha1加密
func Sha1(str string) string {
	result := ""
	if str == "" {
		return result
	}
	b := sha1.Sum([]byte(str))
	return fmt.Sprintf("%x", b)
}

func RandomNumberLen(length int) string {
	var chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	r := mr.New(mr.NewSource(time.Now().UnixNano()))
	var str = ""
	for i := 0; i < length; i++ {
		number := r.Intn(length)
		str += string(chars[number])
	}
	return str
}

func ConfigSign(jsapiTicket, nonceStr, timestamp, url string) string {
	str := "jsapi_ticket=" + jsapiTicket + "&noncestr=" + nonceStr + "&timestamp=" + timestamp + "&url=" + url
	fmt.Println("string1=", str)

	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}