package handler

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/Unknwon/macaron"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"wechatvoice/tool/util"
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

func GetUserString(openId, accessToken string) string {
	fmt.Println("新用户------->>>>>>")
	url := "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openId + "&lang=zh_CN"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("=========150errs")
		fmt.Println(err.Error())
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println("===============================23333====")
	fmt.Println(string(resBody))
	fmt.Println("===============================223333====")

	defer res.Body.Close()
	fmt.Println("===============================111122====")
	fmt.Println(string(resBody))
	fmt.Println("===============================111122====")
	return string(resBody)
}

func GetUserTest(ctx *macaron.Context) string {
	Print("进入index页面方法")
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")

	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/user/getusertest"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	var str string
	Print("获取到的code为==========>>>", code)
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			Print("获取openId 出错", err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		//ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")

		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
		str = GetUserString(res1.OpenId, res1.AccessToken)

	}
	fmt.Println("===================================")
	fmt.Println(str)
	fmt.Println("===================================")
	userInfo := new(UserInfo)
	errs := json.Unmarshal([]byte(str), userInfo)
	if errs != nil {
		fmt.Println(errs, "==========errs")
	}
	fmt.Println(userInfo, "user after unmashall")
	return str
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
func JsapiTicker12() string {
	appid := "wxac69efc11c5e182f"
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
	fmt.Println(resss)
	return resss.Ticket
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

type RedPackages struct {
	Act_name     string `json:"act_name"`
	Client_ip    string `json:"cliend_ip"`
	Remark       string `json:"remark"`
	Re_openid    string `json:"re_openid"`
	Nick_name    string `json:"nick_name"`
	SendNickName string `json:"send_nick_name"`
	Wishing      string `json:"wishing"`
	Amount       int64  `json:"amount"`
	MpId         string `json:"mpid"`
}

func EncodeMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func FormatBizQueryParaMap(m map[string]string) string {
	ms := NewMapSorter(m)
	sort.Sort(ms)

	var buff string

	for _, item := range ms {
		fmt.Printf("%s:%s\n", item.Key, item.Val)
		buff += item.Key + `=` + item.Val + "&"
	}

	return SubstrForCn(buff, ShowStrLenForCn(buff)-1)
}
func SendRedPacket(red *RedPackages) (int64, string) {
	nonceStr := util.GenerateUuid()
	mchid := ""
	key := ""
	appid := ""
	amountStr := strconv.FormatInt(red.Amount, 10)
	mch_billno := mchid + time.Now().Format("20060102") + strconv.FormatInt(time.Now().Unix(), 10)

	m := map[string]string{
		"wxappid":      appid,
		"act_name":     red.Act_name,
		"wishing":      red.Wishing,
		"client_ip":    red.Client_ip,
		"re_openid":    red.Re_openid,
		"mch_id":       mchid,
		"mch_billno":   mch_billno,
		"remark":       red.Remark,
		"nonce_str":    nonceStr,
		"nick_name":    red.Nick_name,
		"send_name":    red.SendNickName,
		"min_value":    amountStr,
		"max_value":    amountStr,
		"total_num":    "1",
		"total_amount": amountStr,
		// "key":          "weibo868weixinpay686chs8mzf8hjh8", //先写到了getsign里
	}
	var v = struct {
		XMLName    xml.Name `xml:"xml"`
		Wxappid    string   `xml:"wxappid"`    //说明：公众号商户appid
		Act_name   string   `xml:"act_name"`   //说明：活动名称
		Wishing    string   `xml:"wishing"`    //说明：红包祝福语
		Client_ip  string   `xml:"client_ip"`  //说明：调用接口的机器Ip地址
		Re_openid  string   `xml:"re_openid"`  //说明：接受收红包的用户在wxappid下的openid
		Mch_id     string   `xml:"mch_id"`     //说明：微信支付分配的商户号
		Mch_billno string   `xml:"mch_billno"` //说明：商户订单号(每个订单号必须唯一) 组成: mch_id + yyyymmdd + 10位一天内不重复的数字。 接口根据商户订单号支持重入, 如出现 超时可再调用。
		Remark     string   `xml:"remark"`     //说明：备注信息
		Nonce_str  string   `xml:"nonce_str"`  //说明：随机字符串,不长于 32 位
		Nick_name  string   `xml:"nick_name"`  //说明：提供方名称
		Send_name  string   `xml:"send_name"`  //说明：红包发送者名称
		Min_value  int64    `xml:"min_value"`  //说明：最小红包金额，单位分
		Max_value  int64    `xml:"max_value"`  //说明：最大红包金额,单位分(最小金额等于最大金额: min_value=max_value =total_amount)
		//Min_value float64 `xml:"min_value"` //说明：最小红包金额，单位分
		//Max_value float64 `xml:"max_value"` //说明：最大红包金额,单位分(最小金额等于最大金额: min_value=max_value =total_amount)
		Total_num    int64 `xml:"total_num"`    //说明：红包发放总人数，total_num=1
		Total_amount int64 `xml:"total_amount"` //说明：付款金额，单位分
		//Total_amount float64 `xml:"total_amount"` //说明：付款金额，单位分
		Sign string `xml:"sign"` //说明：签名
	}{
		Wxappid:      appid,
		Act_name:     red.Act_name,
		Wishing:      red.Wishing,
		Client_ip:    red.Client_ip,
		Re_openid:    red.Re_openid,
		Mch_id:       mchid,
		Mch_billno:   mch_billno,
		Remark:       red.Remark,
		Nonce_str:    nonceStr,
		Nick_name:    red.Nick_name,
		Send_name:    red.SendNickName,
		Min_value:    red.Amount,
		Max_value:    red.Amount,
		Total_num:    1,
		Total_amount: red.Amount,
		Sign:         strings.ToUpper(EncodeMd5(FormatBizQueryParaMap(m) + "&key=" + key)),
	}
	output, _ := xml.Marshal(v)
	vbody := bytes.NewBuffer([]byte(output))
	SSLKEY_PATH := "key/apiclient_key.pem"
	SSLCERT_PATH := "key/apiclient_cert.pem"
	SSLROOTCA_PATH := "key/rootca.pem"
	cert, err := tls.LoadX509KeyPair(SSLCERT_PATH, SSLKEY_PATH)
	if err != nil {
		fmt.Println(err.Error())
	}

	caData, err1 := ioutil.ReadFile(SSLROOTCA_PATH)
	if err != nil {
		fmt.Println(err1.Error())
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	var tlsConfig *tls.Config
	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}
	//beego.Info("========================")
	//beego.Info(vbody)
	//beego.Info("========================")
	res, err := client.Post("https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack", "text/xml", vbody)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	//beego.Info(body)
	//fmt.Println(string(body))
	type result struct {
		XMLName      xml.Name `xml:"xml"`
		Return_Code  string   `xml:"return_code"`
		Return_Msg   string   `xml:"return_msg"`
		Result_Code  string   `xml:"result_code"`
		Err_Code     string   `xml:"err_code"`
		Err_Code_Des string   `xml:"err_code_des"`
		Mch_BillNo   string   `xml:"mch_billno"`
		Mch_Id       string   `xml:"mch_id"`
		WxAppid      string   `xml:"wxappid"`
		Re_OpenId    string   `xml:"re_openid"`
		Total_Amount string   `xml:"total_amount"`
	}
	r := result{}
	errs := xml.Unmarshal(body, &r)
	fmt.Println(errs.Error())
	//fmt.Println(r)
	var success int64
	fmt.Println("=======")
	fmt.Println(r)
	fmt.Println("=======")
	if r.Result_Code == "SUCCESS" {
		success = 1
	} else {
		if r.Err_Code == "NOTENOUGH" {
			fmt.Println("====余额不足")
			success = 2
		} else {
			fmt.Println("发送失败")
			success = 0
		}
		// success = 2

	}

	msg := r.Return_Msg

	return success, msg
}

type MapSorter []MapItem

type MapItem struct {
	Key string
	Val string
}

func (ms MapSorter) Len() int {
	return len(ms)
}

func (ms MapSorter) Less(i, j int) bool {
	// return ms[i].Val < ms[j].Val // 按值排序
	return ms[i].Key < ms[j].Key // 按键排序
}

func (ms MapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func NewMapSorter(m map[string]string) MapSorter {
	ms := make(MapSorter, 0, len(m))

	for k, v := range m {
		ms = append(ms, MapItem{k, v})
	}

	return ms
}
func SubstrForCn(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl+rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
}
func ShowStrLenForCn(s string) int {
	sl := 0
	rs := []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			sl++
		} else {
			sl += 2
		}
	}
	return sl
}
