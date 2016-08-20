package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/macaron"
	"io/ioutil"
	"net/http"
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

func GetOpenCodeInfo(ctx *macaron.Context) string {
	fmt.Println("aaaaaaaaaa")
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=http://www.mylvfa.com&response_type=code&scope=SCOPE&state=STATE#wechat_redirect"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(resBody))

	return ""
}