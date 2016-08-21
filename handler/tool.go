package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/macaron"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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
