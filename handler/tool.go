package handler

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/Unknwon/macaron"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
http://60.205.4.26:22334/configSign?noncestr=W1471365761W&timestamp=1471627311&url=http://www.mylvfa.com/wxpay/config/pay.html
**/
type AResponse struct {
	Sign string `json:"sign"`
}
type ConfigResponse struct {
	Debug     bool     `json:"debug"`
	Appid     string   `json:"appId"`
	NonceStr  string   `json:"nonceStr"`
	Singature string   `json:"signature"`
	JsApiList []string `json:"jsApiList"`
	TimeStamp string   `json:"timeStamp"`
}

func GetSign(ctx *macaron.Context) string {
	signs := time.Now().Unix()
	signsStr := strconv.FormatInt(signs, 10)
	url := "http://60.205.4.26:22334/configSign?noncestr=W1471365761W&timestamp=" + signsStr + "&url=http://www.mylvfa.com/wxpay/config/pay.html"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer res.Body.Close()
	resa := new(AResponse)
	resBody, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(resBody, resa)
	fmt.Println(string(resBody))
	response := new(ConfigResponse)
	response.Debug = false
	response.Appid = "wxac69efc11c5e182f"
	response.NonceStr = "W1471365761W"
	response.Singature = resa.Sign
	response.TimeStamp = signsStr
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

func GetOpenCodeInfo(ctx *macaron.Context) {
	fmt.Println("aaaaaaaaaa")

	re := "http://www.mylvfa.com/voice/tool/info"
	a := url.QueryEscape(re)
	fmt.Println(a)
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
	fmt.Println(url)
	// "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=http://www.mylvfa.com&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
	// res, err := http.Get(url)
	// if err != nil {www.mylvfa.com&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect
	// 	fmt.Println(err.Error())
	//https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=http://www.mylvfa.com&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect
	// }

	ctx.Redirect(url)

}

func AuthCodeURL(appId, redirectURI, scope, state string) string {
	return "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=http://,ylvfa.com&response_type=code&scope=SCOPE&state=STATE#wechat_redirect"
}

type OpenIdResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expores_in"`
	RefressToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
}

func GetAllInfo(ctx *macaron.Context) {
	fmt.Println("=====================>>>>>")
	// ctx.Params(name)
	code := ctx.Query("code")
	fmt.Println(url.QueryEscape("code"))
	fmt.Println(ctx.Params("code"))
	url := "http://60.205.4.26:22334/getOpenid?code=" + code
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(resBody))
	defer res.Body.Close()

	res1 := new(OpenIdResponse)
	json.Unmarshal(resBody, res1)
	fmt.Println(res1)
	// code := url.QueryEscape("code")
	// fmt.Println(code)
}

type UserInfo struct {
	OpenId     string   `json:"openId"`
	NickName   string   `json:"nickName"`
	Sex        string   `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgUrl string   `json:"headImgUrl"`
	Privilege  []string `json:"privilege"`
	UnionId    string   `json:"unionid"`
}

func GetUserInfo(openId, accessToken string) *UserInfo {
	fmt.Println("新用户=================================》》》")
	url := "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openId + "&lang=zh_CN"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(resBody))
	defer res.Body.Close()

	res1 := new(UserInfo)
	json.Unmarshal(resBody, res1)
	fmt.Println(res1)
	return res1
}

// func GetAllUtil(code string) error {
// 	url := "http://60.205.4.26:22334/getOpenid?code=" + code
// 	res, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println("=========xxxxx")
// 		fmt.Println(err.Error())
// 		return err
// 	}
// 	resBody, _ := ioutil.ReadAll(res.Body)
// 	fmt.Println(string(resBody))
// 	defer res.Body.Close()
// 	fmt.Println("==========>>>>")
// 	res1 := new(OpenIdResponse)
// 	json.Unmarshal(resBody, res1)
// 	ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
// 	member := new(model.MemberInfo)
// 	memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
// 	if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
// 		return memberErr
// 	}
// 	if member.Uuid == "" {
// 		fmt.Println("新的用户")
// 		user := GetUserInfo(res1.OpenId, res1.AccessToken)
// 		member.Uuid = util.GenerateUuid()
// 		member.HeadImgUrl = user.HeadImgUrl
// 		member.OpenId = user.OpenId
// 		member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
// 		member.NickName = user.NickName
// 		err := member.GetConn().Create(&member).Error
// 		if err1 != nil {
// 			return err1
// 		}
// 	}
// }

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
func GenerateXMLStr(params map[string]string) string {
	result := "<xml>"

	for k, v := range params {
		result = result + "<" + k + "><![CDATA[" + v + "]]></" + k + ">"
	}

	result += "</xml>"
	return result
}

// 获取JSAPI_TICKET返回值
type RespJsapiTicket struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`

	Ticket string `json:"jsapiTicket"`
}

func JsapiTicker1(ctx *macaron.Context) string {
	appid := ctx.Query("appid")
	fmt.Println(appid)
	result := new(RespJsapiTicket)
	res, err := http.Get("http://www.mylvfa.com/getAccessToken")
	if err != nil {
		result.Code = CODE_ERROR
		result.Msg = err.Error()
		ret_str, _ := json.Marshal(result)
		return string(ret_str)
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(resBody))
	defer res.Body.Close()
	token := string(resBody)

	url1 := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + token + "&type=jsapi"
	res2, res2Err := http.Get(url1)
	if res2Err != nil {
		result.Code = CODE_ERROR
		result.Msg = res2Err.Error()
		ret_str, _ := json.Marshal(result)
		return string(ret_str)
	}
	resBody2, _ := ioutil.ReadAll(res2.Body)
	fmt.Println(string(resBody))
	defer res2.Body.Close()
	type ApiTk struct {
		ErrorCode int64  `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
		Ticket    string `json:"ticket"`
		Expires   int64  `json:"expires_in"`
	}

	resss := new(ApiTk)
	json.Unmarshal(resBody2, resss)
	result.Code = 10000
	result.Msg = "ok"
	result.Ticket = resss.Ticket
	str, _ := json.Marshal(result)
	return string(str)
}
func Sha1(str string) string {
	result := ""
	if str == "" {
		return result
	}
	b := sha1.Sum([]byte(str))
	return fmt.Sprintf("%x", b)
}

func Md5(str string) string {
	result := ""
	if str == "" {
		return result
	}
	b := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", b)
}
