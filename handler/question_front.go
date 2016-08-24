package handler

import (
	"bytes"
	"encoding/json"
	"encoding/xml"

	"github.com/Unknwon/macaron"

	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	//"sort"
	"strconv"
	"strings"
	"time"
	"wechatvoice/model"
	"wechatvoice/tool/util"
)

const (
	CODE_SUCCESS            = 10000
	CODE_ERROR              = 10001
	CODE_REDIRECT           = 10002
	MSG_SUCCESS             = "ok"
	RNF                     = "record not found"
	DEFAULT_DEVICE_INFO     = "WEB"
	DEFAULT_FEE_TYPE        = "CNY"
	DEFAULT_TRADE_TYPE      = "JSAPI"
	DEFAULT_NOTIFY_URL      = "/wechatvoice/pay/decodewechatpayinfo"
	UNIFIEDORDER_URL        = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	DEFAULT_PACKAGE_PRE_STR = "prepay_id="
	DEFAULT_SIGN_TYPE       = "MD5"
	TICKET_SERVER_URL       = "/shangqu-3rdparty/token/jsapi_ticket?appid="
	APPID                   = ""
	APPSECRET               = ""
	MCHID                   = ""
	MCHNAME                 = ""
	KEY                     = ""
	SERVER_IP               = "127.0.0.1"

	PAY_PAGE_URL                = "orderSubmit.html"
	AFTER_PAY_ORDER_URL         = "/shangqu-shop/afterpay/wx"
	AFTER_PAY_JUMP_PAGE_FAILD   = "payFailed.html?"
	AFTER_PAY_JUMP_PAGE_SUCCESS = "paySuccess.html?"

	WECHAT_PREPAY_URL = "/wechatvoice/pay/unifiedorder?appid=%s&mch_id=%s&body=%s&out_trade_no=%s&total_fee=%d&spbill_create_ip=%s&key=%s&openid=%s&url=%s&notify_url=%s"
)

var merchantIndexUrl = "http://www.mylvfa.com/daodaolaw/search.html"

//查询问题返回
type QuestionQueryResponse struct {
	Code  int64          `json:"code"`
	Msg   string         `json:"msg"`
	List  []QuestionInfo `json:"list"`
	Total int64          `json:"total"`
}

type QuestionInfo struct {
	OrderId   string     `json:"orderId"`
	LaywerId  string     `json:"laywerId"`
	Question  string     `json:"question"`
	Name      string     `json:"name"`
	SelfIntr  string     `json:"selfIntr"`
	LawerPic  string     `json:"pic"`
	Answer    string     `json:"answer"`
	TypeId    string     `json:"typeId"`
	TypeName  string     `json:"typeName"`
	TypePrice string     `json:"typePrice"`
	Star      int64      `json:"star"`
	IsPay     bool       `json:"isPay"`
	AddNum    int64      `json:"addNum"`
	IsShow    bool       `json:"isShow"`
	AddInfo   []AddInfos `json:""`
}

type AddInfos struct {
	QuestionInfo string `json:"question"`
	OrderId      string `json:"orderId"`
	Answer       string `json:"answer"`
}

//查询问题方法

func ToIndex(ctx *macaron.Context) {
	fmt.Println("=================进入方法")
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")

	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/toindex"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	fmt.Println("============code is ")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			fmt.Println(memberErr.Error(), "=====会员出错")
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				fmt.Println(err.Error(), "xxxxx")
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	openId := strings.Split(cookieStr, "|")[0]
	// userType := strings.Split(cookieStr, "|")[1]
	fmt.Println(openId)
	ctx.Resp.Write([]byte("<script type=\"text/javascript\">window.location.href=\"" + merchantIndexUrl + "\"</script>"))
}

func QuestionQuery(ctx *macaron.Context) string {
	response := new(QuestionQueryResponse)
	fmt.Println("=================进入方法")
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")

	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/questionquery"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	fmt.Println("============code is ")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(res)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			log.Println(user)
			log.Println(user.HeadImgUrl)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	openId := strings.Split(cookieStr, "|")[0]
	// userType := strings.Split(cookieStr, "|")[1]
	fmt.Println(openId)
	body, _ := ctx.Req.Body().String()
	req := new(model.QuestionQuery)

	marshallErr := json.Unmarshal([]byte(body), req)

	if marshallErr != nil {
		response.Code = CODE_ERROR
		response.Msg = marshallErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	questionList, count, queryErr := model.GetQuestionQuery(*req)
	fmt.Println(questionList)
	if queryErr != nil && !strings.Contains(queryErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = queryErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	retList := make([]QuestionInfo, 0)
	for _, k := range questionList {
		single := new(QuestionInfo)
		lawyer := new(model.LawyerInfo)
		lawyerErr := lawyer.GetConn().Where("uuid = ?", k.AnswerId).Find(&lawyer).Error
		if lawyerErr != nil && !strings.Contains(lawyerErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = lawyerErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		single.OrderId = k.Uuid
		single.LaywerId = k.AnswerId
		single.Question = k.Description
		single.Name = k.AnswerName
		single.SelfIntr = lawyer.FirstCategory
		single.LawerPic = lawyer.HeadImgUrl
		single.Answer = k.VoicePath
		single.TypeId = k.CategoryId
		single.TypeName = k.Category
		cateInfo := new(model.WechatVoiceQuestionSettings)
		cateErr := cateInfo.GetConn().Where("category_id = ?", k.CategoryId).Find(&cateInfo).Error
		if cateErr != nil && !strings.Contains(cateErr.Error(), RNF) {
			fmt.Println("=================>>>>>cateErr")
			fmt.Println(cateErr.Error())
			fmt.Println("=================>>>>>")
			response.Code = CODE_ERROR
			response.Msg = cateErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		single.TypePrice = cateInfo.PayAmount
		rank, _ := strconv.ParseInt(k.RankInfo, 10, 64)
		single.Star = rank
		payment := new(model.WechatVoicePaymentInfo)
		payErr := payment.GetConn().Where("question_id = ?", k.Uuid).Where("open_id = ?", openId).Find(&payment).Error

		if payErr != nil && !strings.Contains(payErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = payErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		var payAble bool
		if payment.Uuid != "" {
			//说明有支付记录
			payAble = true
		} else {
			payAble = false
		}
		single.IsPay = payAble
		childList, childErr := model.GetChildAnsers(k.Uuid)
		if childErr != nil && !strings.Contains(childErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = childErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		single.AddNum = int64(len(childList))
		single.IsShow = false
		addInfo := make([]AddInfos, 0)
		if len(childList) > 0 {
			for _, v := range childList {
				singleChild := new(AddInfos)
				singleChild.OrderId = v.Uuid
				singleChild.QuestionInfo = v.Description
				singleChild.Answer = v.VoicePath
				addInfo = append(addInfo, *singleChild)
			}
		}
		single.AddInfo = addInfo
		retList = append(retList, *single)
	}

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.List = retList
	response.Total = count
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

//提问新的问题
/**
提问新的问题

params


*/
//init的那个方法在哪里
type NewQuestionRequest struct {
	CateId    string `json:"typeId"`
	TypePrice string `json:"typePrice"`
	Content   string `json:"content"`
}

type NewQuestionResponse struct {
	Code        int64  `json:"code"`
	Msg         string `json:"msg"`
	OrderNumber string `json:"orderId"`
	Payment     string `json:"price"`
	IsAdd       string `json:"isAdd"`
}

type OrderResponse struct {
	Code      int64  `json:"code"`
	Msg       string `json:"msg"`
	Appid     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
	OrderId   string `json:"orderId"`
}

func CreateNewQuestion(ctx *macaron.Context) string {
	response := new(OrderResponse)

	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/createquestion"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(res)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]
	fmt.Println(openId, userType)
	body, _ := ctx.Req.Body().String()

	req := new(NewQuestionRequest)

	unmarshallErr := json.Unmarshal([]byte(body), req)
	fmt.Println("发问请求提")
	fmt.Println(body)
	fmt.Println("发问请求提")
	if unmarshallErr != nil {
		fmt.Println(unmarshallErr.Error())
		response.Code = CODE_ERROR
		response.Msg = unmarshallErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	//fmt.Println(a)
	cateId := req.CateId
	typePrice := req.TypePrice
	content := req.Content

	cate := new(model.Category)
	cateErr := cate.GetConn().Where("uuid = ?", cateId).Find(&cate).Error

	if cateErr != nil && !strings.Contains(cateErr.Error(), RNF) {
		fmt.Println(cateErr)
		response.Code = CODE_ERROR
		response.Msg = cateErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	customer := new(model.MemberInfo)

	customerErr := customer.GetConn().Where("open_id = ?", openId).Find(&customer).Error

	if customerErr != nil && !strings.Contains(customerErr.Error(), RNF) {
		response.Code = CODE_ERROR
		fmt.Println(customerErr.Error())
		response.Msg = customerErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	fmt.Println("here.....")
	orderNumber := util.GenerateOrderNumber()
	question := new(model.WechatVoiceQuestions)
	question.Uuid = util.GenerateUuid()
	question.CategoryId = cateId
	question.Category = cate.CategoryName
	question.CategoryIdInt = int64(cate.Model.ID)
	question.Description = content
	today := time.Unix(time.Now().Unix(), 0).String()[0:19]
	question.CreateTime = today
	question.CustomerId = customer.Uuid
	question.CustomerName = customer.Name
	question.CustomerOpenId = openId
	question.PaymentInfo = typePrice
	question.IsSolved = "0"
	question.AskerHeadImg = customer.HeadImgUrl
	payInt, transferErr := strconv.ParseInt(typePrice, 10, 64)
	question.OrderNumber = orderNumber
	if transferErr != nil && !strings.Contains(transferErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = transferErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	question.PaymentInfoInt = payInt
	fmt.Println("here............")
	createErr := question.GetConn().Create(&question).Error

	if createErr != nil {
		response.Code = CODE_ERROR
		response.Msg = createErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	nstr := util.GenerateUuid()
	nSt := util.GenerateUuid()
	timeStamp := time.Now().Unix()
	fmt.Println(timeStamp)
	tStr := strconv.FormatInt(timeStamp, 10)
	sign, prepayId, sings, signErr := PayBill(nstr, nSt, openId, orderNumber, "1", tStr)
	if signErr != nil {
		fmt.Println(signErr.Error())
		response.Code = CODE_ERROR
		response.Msg = signErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	fmt.Println(sings)
	// paramsList := []string{"appid", "mch_id", "body", "out_trade_no", "total_fee", "spbill_create_ip", "device_info", "nonce_str", "fee_type", "time_start", "notify_url", "trade_type"}
	// paramsMap := make(map[string]string, 0)
	// paramsMap["appid"] = "wxac69efc11c5e182f"
	// paramsMap["mch_id"] = "1344737201"
	// paramsMap["body"] = "my_pay_test"
	// paramsMap["out_trade_no"] = orderNumber
	// paramsMap["total_fee"] = "1"
	// paramsMap["spbill_create_ip"] = "127.0.0.1"

	// // 有默认值的字段处理
	// paramsMap["device_info"] = DEFAULT_DEVICE_INFO
	// paramsMap["nonce_str"] = nstr
	// paramsMap["fee_type"] = DEFAULT_FEE_TYPE
	// paramsMap["time_start"] = tStr
	// paramsMap["notify_url"] = "http://www.mylvfa.com/wxpay/config/"
	// paramsMap["trade_type"] = DEFAULT_TRADE_TYPE
	// //key := "C4CA4238A0B923820DCC509A6F75849B"
	// signself := GeneratePageSign(paramsMap, paramsList)
	// fmt.Println("keyself", signself)
	// fmt.Println(sings)
	signSelf := GetSigns(tStr)
	/**

	Appid     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
	*/
	// signnew := GetSigns(tStr)
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.Appid = "wxac69efc11c5e182f"
	response.NonceStr = nstr
	response.Signature = signSelf
	response.SignType = "MD5"
	response.Package = "prepay_id=" + prepayId
	response.TimeStamp = tStr
	response.PaySign = sign
	response.OrderId = orderNumber
	ret_str, _ := json.Marshal(response)
	fmt.Println("=====================================>>>>>")
	fmt.Println(string(ret_str))
	fmt.Println("=====================================>>>>>")
	return string(ret_str)
}

type DoPayReq struct {
	OrderId string `json:"orderId"`
}

// func DoPay(ctx *macaron.Context) string {
// 	//http://60.205.4.26:22334/prepayId\?appid\=wxac69efc11c5e182f\
// 	//http://60.205.4.26:22334/prepayId?appid=wxac69efc11c5e182f&body=payTest&mch_id=1344737201&nonce_str=W1471365761W&
// 	//notify_url=http://www.mylvfa.com/wxpay/config/&openid=o-u0Nv-Hd01lO6wpKyb8VGQaNGpg&out_trade_no=20150806125346&spbill_create_ip=127.0.0.1&total_fee=1&trade_type=JSAPI"	//
// 	req := new(DoPayReq)
// 	body, _ := ctx.Req.Body().String()
// 	json.Unmarshal([]byte(body), req)
// 	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
// 	if cookieStr == "" && ctx.Query("code") == "" {
// 		re := "http://www.mylvfa.com/voice/front/createnewspecialquestion"
// 		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
// 		//cookieStr = "1|2"
// 		ctx.Redirect(url)
// 	}
// 	code := ctx.Query("code")
// 	if code != "" {
// 		url := "http://60.205.4.26:22334/getOpenid?code=" + code
// 		res, err := http.Get(url)
// 		if err != nil {
// 			fmt.Println("=========xxxxx")
// 			fmt.Println(err.Error())
// 		}
// 		resBody, _ := ioutil.ReadAll(res.Body)
// 		fmt.Println(string(resBody))
// 		defer res.Body.Close()
// 		fmt.Println("==========>>>>")
// 		res1 := new(OpenIdResponse)
// 		json.Unmarshal(resBody, res1)
// 		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
// 		member := new(model.MemberInfo)
// 		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
// 		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
// 			response.Code = CODE_ERROR
// 			response.Msg = memberErr.Error()
// 			ret_str, _ := json.Marshal(res)
// 			return string(ret_str)
// 		}
// 		if member.Uuid == "" {
// 			fmt.Println("新的用户")
// 			user := GetUserInfo(res1.OpenId, res1.AccessToken)
// 			member.Uuid = util.GenerateUuid()
// 			member.HeadImgUrl = user.HeadImgUrl
// 			member.OpenId = user.OpenId
// 			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
// 			member.NickName = user.NickName
// 			err := member.GetConn().Create(&member).Error
// 			if err != nil {
// 				response.Code = CODE_ERROR
// 				response.Msg = err.Error()
// 				ret_str, _ := json.Marshal(response)
// 				return string(ret_str)
// 			}
// 		}
// 		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
// 	}
// 	order := new(model.WechatVoiceQuestions)
// 	orderErr := order.GetConn().Where("order_number = ?", req.OrderId).Find(&order).Error
// 	if orderErr != nil && !strings.Contains(orderErr.Error(), RNF) {
// 		fmt.Println(orderErr.Error())
// 	}
// 	openId := strings.Split(cookieStr, "|")[0]
// 	userType := strings.Split(cookieStr, "|")[1]
// 	fmt.Println(cookieStr)
// 	// appId := "wxac69efc11c5e182f"
// 	// mchId := "1344737201"
// 	nonceStr := util.GenerateUuid()
// 	notify_url := "http://www.mylvfa.com/voice/afterpay"
// 	openId1 := openId
// 	orderId := req.OrderId
// 	pay := order.PaymentInfo

// 	url := "http://60.205.4.26:22334/prepayId?appid=wxac69efc11c5e182f&body=pay&mch_id=1344737201&nonce_str=" + nonceStr + "&notify_url=" + notify_url + "&openid=" + openId + "&out_trade_no=" + req.OrderId + "&spbill_create_ip=127.0.0.1&total_fee=" + pay + "&trade_type=JSAPI"
// 	res, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	a := new(PayJson)
// 	resBody, _ := ioutil.ReadAll(res.Body)
// 	fmt.Println(string(resBody))
// 	defer res.Body.Close()
// 	json.Unmarshal(resBody, a)

// 	str := a.Sign

// }

type SpecialQuestions struct {
	CateId     string `json:"typeId"`
	TypePrice  string `json:"typePrice"`
	Content    string `json:"content"`
	QuestionId string `json:"quesionId"`
	LawyerId   string `json:"lawyerId"`
}

type PayJson struct {
	Sign string `json:"sign"`
}

func CreateNewSpecialQuestion(ctx *macaron.Context) string {
	response := new(NewQuestionResponse)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/createnewspecialquestion"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(res)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]
	fmt.Println(openId, userType)
	body, _ := ctx.Req.Body().String()

	req := new(SpecialQuestions)

	unmarshallErr := json.Unmarshal([]byte(body), req)
	if unmarshallErr != nil {
		response.Code = CODE_ERROR
		response.Msg = unmarshallErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	cate := new(model.Category)
	cateErr := cate.GetConn().Where("uuid = ?", req.CateId).Find(&cate).Error

	if cateErr != nil && !strings.Contains(cateErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = cateErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	customer := new(model.MemberInfo)

	customerErr := customer.GetConn().Where("open_id = ?", openId).Find(&customer).Error

	if customerErr != nil && !strings.Contains(customerErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = customerErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	orderNumber := util.GenerateOrderNumber()
	question := new(model.WechatVoiceQuestions)
	question.Uuid = util.GenerateUuid()
	question.CategoryId = req.CateId
	question.Category = cate.CategoryName
	question.CategoryIdInt = int64(cate.Model.ID)
	question.Description = req.Content
	today := time.Unix(time.Now().Unix(), 0).String()[0:19]
	question.CreateTime = today
	question.CustomerId = customer.Uuid
	question.CustomerName = customer.Name
	question.CustomerOpenId = openId
	question.PaymentInfo = req.TypePrice
	payInt, transferErr := strconv.ParseInt(req.TypePrice, 10, 64)
	question.OrderNumber = orderNumber
	if transferErr != nil && !strings.Contains(transferErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = transferErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	question.PaymentInfoInt = payInt
	if req.QuestionId != "" {
		//追加
		question1 := new(model.WechatVoiceQuestions)
		questionErr := question1.GetConn().Where("uuid = ?", req.QuestionId).Find(&question1).Error
		if questionErr != nil && !strings.Contains(questionErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = questionErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		question.ParentQuestionId = req.QuestionId
		question.AppenQuestionTime = question.AppenQuestionTime + 1
		question.HaveAppendChild = "1"
	}
	if req.LawyerId != "" {
		//指定问题
		question.QuestionType = "1"
		question.Important = "1"
	}

	createErr := question.GetConn().Create(&question).Error

	if createErr != nil {
		response.Code = CODE_ERROR
		response.Msg = createErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.OrderNumber = orderNumber
	response.Payment = req.TypePrice
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

//这里获取分类问题的配置选项
type GetConfigRequest struct {
	CateGoryId string `json:"categoryId"`
}

type QuestionConfig struct {
	Code         int64  `json:"code"`
	Msg          string `json:"msg"`
	PayAmountStr string `json:"amountStr"`
	PayAmountInt int64  `json:"amountInt"`
}

func GetQuestionConfig(ctx *macaron.Context) string {
	body, _ := ctx.Req.Body().String()
	req := new(GetConfigRequest)
	unmarshallErr := json.Unmarshal([]byte(body), req)

	response := new(QuestionConfig)

	if unmarshallErr != nil {
		response.Code = CODE_ERROR
		response.Msg = unmarshallErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	setting := new(model.WechatVoiceQuestionSettings)
	settingErr := setting.GetConn().Where("category_id = ?", req.CateGoryId).Find(&setting).Error
	if settingErr != nil && !strings.Contains(settingErr.Error(), "record not found") {
		response.Code = CODE_ERROR
		response.Msg = settingErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = "ok"
	response.PayAmountInt = setting.PayAmountInt
	response.PayAmountStr = setting.PayAmount
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

//获取分类列表1
type QuestionCateList struct {
	Code int64      `json:"code"`
	Msg  string     `json:"msg"`
	List []CateInfo `json:"list"`
}

type CateInfo struct {
	CateName        string `json:"typeId"`
	CateId          string `json:"typeName"`
	CatePaymentInfo string `json:"typePrice"`
}

func GetQuestionCateList(ctx *macaron.Context) string {
	response := new(QuestionCateList)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")

	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/getcatList"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(res)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	list := make([]CateInfo, 0)
	fmt.Println("========进入这个方法啦")
	cateList, cateErr := model.GetCateLists()

	fmt.Println(len(cateList))
	if cateErr != nil && !strings.Contains(cateErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = cateErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	for _, k := range cateList {
		single := new(CateInfo)
		single.CateName = k.Uuid
		single.CateId = k.CategoryName
		price := new(model.WechatVoiceQuestionSettings)
		priceErr := price.GetConn().Where("category_id = ?", k.Uuid).Find(&price).Error
		if priceErr != nil {
			fmt.Println(priceErr.Error())
		}
		if priceErr != nil && !strings.Contains(priceErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = priceErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		single.CatePaymentInfo = price.PayAmount
		fmt.Println(single)
		list = append(list, *single)
	}

	response.Code = CODE_SUCCESS
	response.Msg = "ok"
	response.List = list
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

//将问题进行追问
type QuestionAppendRequest struct {
	QuestionId   string `json:"parentQId"`
	CateId       string `json:"cateId"`
	CateName     string `json:"cateName"`
	AskerOpenId  string `json:"askerOpenId"`
	Description  string `json:"description"`
	QuestionType string `json:"type"` //0 直接丢出去 1 指定人提问
	TargetOpenId string `json:"targetOpenId"`
	Payment      string `json:"payment"`
}

type QuestionNewResponse struct {
	Code        int64  `json:"code"`
	Msg         string `json:"msg"`
	OrderNumber string `json:"orderNumber"`
}

func AppendQuestion(ctx *macaron.Context) string {
	response := new(QuestionNewResponse)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/appendquestion"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(res)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)

	body, _ := ctx.Req.Body().String()
	req := new(QuestionAppendRequest)

	json.Unmarshal([]byte(body), req)

	questionInfoOld := new(model.WechatVoiceQuestions)
	questionOldErr := questionInfoOld.GetConn().Where("uuid = ?", req.QuestionId).Find(&questionInfoOld).Error

	if questionOldErr != nil && !strings.Contains(questionOldErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = questionOldErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if questionInfoOld.IsSolved != "2" {
		//问题没有解决 不能进行追问
		response.Code = CODE_ERROR
		response.Msg = "问题没有解决 不能进行追问"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	customer := new(model.MemberInfo)

	customerErr := customer.GetConn().Where("open_id = ?", req.AskerOpenId).Find(&customer).Error

	if customerErr != nil && !strings.Contains(customerErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = customerErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	questionNew := new(model.WechatVoiceQuestions)

	questionNew.Uuid = util.GenerateUuid()
	questionNew.Category = req.CateName
	questionNew.CategoryId = req.CateId
	questionNew.Description = req.Description
	today := time.Unix(time.Now().Unix(), 0).String()[0:19]
	questionNew.CreateTime = today
	questionNew.CustomerId = customer.Uuid
	questionNew.CustomerName = customer.Name
	questionNew.CustomerOpenId = req.AskerOpenId
	questionNew.PaymentInfo = req.Payment
	questionNew.Important = "1"
	questionNew.QuestionType = "1"
	questionNew.ParentQuestionId = req.QuestionId
	on := util.GenerateOrderNumber()
	questionNew.OrderNumber = on
	questNewErr := questionNew.GetConn().Create(&questionNew).Error

	if questNewErr != nil {
		response.Code = CODE_ERROR
		response.Msg = questNewErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.OrderNumber = on
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

//偷听业务
type PeekAnswerRequest struct {
	OrderId string `json:"orderId"`
	CateId  string `json:"cateId"`
}

/**
偷听应该是这样用户点击 先看是否有权限去听 如果有 直接放 如果不能 调起微信支付 支付结束后 可以看
*/
type PeekResponse struct {
	Code     int64  `json:"code"`
	Msg      string `json:"msg"`
	PlayAble bool   `json:"playAble"`
}

//这里需要看下微信支付业务 点击获取之后

func PeekAvalable(ctx *macaron.Context) string {
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	response := new(PeekResponse)

	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/peekavalable"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(res)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	body, _ := ctx.Req.Body().String()
	req := new(PeekAnswerRequest)
	json.Unmarshal([]byte(body), req)

	//设置cookie  第一段为openId 第二段为类型 1 用户 2律师
	// cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" {
		//这里直接调取util重新过一次绿叶 获取openId 等信息
	}
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]

	log.Println(openId)
	log.Println(userType)

	pay := new(model.OrderPaymentInfo)
	payErr := pay.GetConn().Where("question_id = ?", req.OrderId).Where("open_id = ?", openId).Find(&pay).Error

	if payErr != nil && !strings.Contains(payErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = payErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if pay.Uuid == "" {
		response.Code = CODE_SUCCESS
		response.Msg = MSG_SUCCESS
		response.PlayAble = false
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.PlayAble = true
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

func DoPay(ctx *macaron.Context) string {
	//这里去看下微信支付功能
	return ""
}

//开始回答部分

//一个文件上传服务

func VoiceUpLoad(ctx *macaron.Context) string {
	//这里返回一个路径
	return ""
}

type AnswerQuestion1 struct {
	QuestionId string `json:"questionId"`
}

type AnswerQuestion1Response struct {
	Code           int64  `json:"code"`
	Msg            string `json:"msg"`
	QuestionInfoss `json:"question"`
}
type QuestionInfoss struct {
	QuestionId       string `json:"quesiontId"`
	QuestionCateInfo string `json:"cateInfo"`
	QuestionCateId   string `json:"cateId"`
	QuestionDesc     string `json:"desc"`

	AskerName    string `json:"askerName"`
	AskerId      string `json:"askerId"`
	AskerHeadImg string `json:"askerHeadImg"`
	LawyerId     string `json:"lawyerId"`
	LawerName    string `json:"name"`
}

func AnswerQuestionInit(ctx *macaron.Context) string {
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	response := new(AnswerQuestion1Response)
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/answerquestioninit"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(res)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	//点击回答问题  显示问题
	body, _ := ctx.Req.Body().String()
	req := new(AnswerQuestion1)
	json.Unmarshal([]byte(body), req)

	//设置cookie  第一段为openId 第二段为类型 1 用户 2律师
	// cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" {
		//这里直接调取util重新过一次绿叶 获取openId 等信息
	}
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]

	log.Println(openId)
	log.Println(userType)

	question := new(model.WechatVoiceQuestions)
	qErr := question.GetConn().Where("uuid = ?", req.QuestionId).Find(&question).Error

	if qErr != nil && !strings.Contains(qErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = qErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if question.IsSolved == "1" {
		response.Code = CODE_ERROR
		response.Msg = "已经回答完毕 请勿重复做大"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if question.Important == "1" && question.AnswerOpenId != openId {
		response.Code = CODE_ERROR
		response.Msg = "这个为指定问题"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if userType == "1" {
		response.Code = CODE_ERROR
		response.Msg = "无权限"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	q := new(QuestionInfoss)

	// q.QuestionId = question.Uuid
	// q.HeadImg = question.AskerHeadImg
	// q.QuestionCategoryId = question.CategoryId
	// q.QuestionCateName = question.Category
	// q.VoicePath = ""
	// q.QuestionTopic = question.Description
	q.QuestionDesc = question.Uuid
	q.AskerHeadImg = question.AskerHeadImg
	q.QuestionCateId = question.CategoryId
	q.QuestionCateInfo = question.Category
	q.QuestionDesc = question.Description

	q.AskerName = question.CustomerName
	q.AskerId = question.CustomerId

	lawerInfo := new(model.LawyerInfo)

	lawerInfoErr := lawerInfo.GetConn().Where("open_id = ?", openId).Find(&lawerInfo).Error

	if lawerInfoErr != nil && !strings.Contains(lawerInfoErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = lawerInfoErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	q.LawyerId = lawerInfo.Uuid
	q.LawerName = lawerInfo.Name
	//在这里锁住
	lock := new(model.AnswerLockInfo)
	/*lockErr:=lock.GetConn().Where("question_id = ?",question.Uuid).Find(&lock).Error
	if lockErr!=nil&&!strings.Contains(lockErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = lockErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}*/
	lock.Uuid = util.GenerateUuid()
	lock.QuestionId = question.Uuid
	lock.OpenIdFirst = openId
	lock.LockedTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
	err := lock.GetConn().Create(&lock).Error
	if err != nil {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.QuestionInfoss = *q
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

type DoAnsweQuestion struct {
	FilPath    string `json:"filePath"`
	QuestionId string `json:"questionId"`
}

func DoAnswerQuestion(ctx *macaron.Context) string {
	response := new(model.GeneralResponse)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/doanswerquestion"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(res)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	body, _ := ctx.Req.Body().String()
	req := new(DoAnsweQuestion)
	json.Unmarshal([]byte(body), req)

	question := new(model.WechatVoiceQuestions)

	//设置cookie  第一段为openId 第二段为类型 1 用户 2律师
	// cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	// if cookieStr == "" {
	// 	//这里直接调取util重新过一次绿叶 获取openId 等信息
	// }
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]

	log.Println(openId)
	log.Println(userType)

	questionErr := question.GetConn().Where("uuid = ?", req.QuestionId).Find(&question).Error
	if questionErr != nil && !strings.Contains(questionErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = questionErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if question.IsSolved == "1" {
		response.Code = CODE_ERROR
		response.Msg = "已经回答完毕 请勿重复做大"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if question.Important == "1" && question.AnswerOpenId != openId {
		response.Code = CODE_ERROR
		response.Msg = "这个为指定问题"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if userType == "1" {
		response.Code = CODE_ERROR
		response.Msg = "无权限"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	lawerInfo := new(model.LawyerInfo)

	lawerInfoErr := lawerInfo.GetConn().Where("open_id = ?", openId).Find(&lawerInfo).Error

	if lawerInfoErr != nil && !strings.Contains(lawerInfoErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = lawerInfoErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	question.IsAnswerd = "1"
	question.VoicePath = req.FilPath
	question.AnswerId = lawerInfo.Uuid
	question.AnswerOpenId = openId
	question.AnswerName = lawerInfo.Name
	question.AnswerHeadImg = lawerInfo.HeadImgUrl

	err := question.GetConn().Save(&question).Error

	if err != nil && !strings.Contains(err.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

func AddPeekCount(questionId string) error {
	quesionInfo := new(model.WechatVoiceQuestions)
	quesionInfoErr := quesionInfo.GetConn().Where("uuid = ?", questionId).Find(&quesionInfo).Error
	if quesionInfoErr != nil {
		return quesionInfoErr
	}
	quesionInfo.AnswerdCount = quesionInfo.AnswerdCount + 1
	err := quesionInfo.GetConn().Save(&quesionInfo).Error
	return err
}

type SingleQuestionInfo struct {
	Code         int64  `json:"code"`
	Msg          string `json:"msg"`
	QuestionInfo `json:"questionInfo"`
}

//根据ID获取问题详情
func GetQuestionInfoById(ctx *macaron.Context) string {

	qId := ctx.Query("id")
	questionInfo := new(model.WechatVoiceQuestions)
	response := new(SingleQuestionInfo)
	qErr := questionInfo.GetConn().Where("uuid = ?", qId).Find(&questionInfo).Error
	if qErr != nil && !strings.Contains(qErr.Error(), "record not found") {
		response.Code = CODE_ERROR
		response.Msg = qErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	var single QuestionInfo
	// single.LawerPic = questionInfo.AskerHeadImg
	// single.QuestionId = questionInfo.Uuid
	// single.QuestionTopic = questionInfo.Description
	// single.LawyerId = questionInfo.AnswerId
	// single.LawyerName = questionInfo.AnswerName
	// single.QuestionCategoryId = questionInfo.CategoryId
	// single.QuestionCateName = questionInfo.Category
	// single.IsSolved = questionInfo.IsSolved

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.QuestionInfo = single
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

//做出评价
type RankAnswerReq struct {
	QuestionId string `json:"questionId"`
	RankInfo   int64  `json:"rank"`
}
type RankResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Red  string `json:"redPacket"`
}

func RankTheAnswer(ctx *macaron.Context) string {
	response := new(RankResponse)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/ranktheanswer"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(res)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	// fmt.Println(cookieStr)
	// cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	// if cookieStr == "" {
	// 	//这里直接调取util重新过一次绿叶 获取openId 等信息
	// }
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]
	fmt.Println(openId, userType)
	body, _ := ctx.Req.Body().String()
	req := new(RankAnswerReq)

	json.Unmarshal([]byte(body), req)
	questionInfo := new(model.WechatVoiceQuestions)

	questionInfoErr := questionInfo.GetConn().Where("uuid = ?", req.QuestionId).Find(&questionInfo).Error
	if questionInfoErr != nil && !strings.Contains(questionInfoErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = questionInfoErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	questionInfo.IsRanked = "1"
	r := strconv.FormatInt(req.RankInfo, 10)
	questionInfo.RankInfo = r

	errUpdate := questionInfo.GetConn().Save(&questionInfo).Error
	if errUpdate != nil {
		response.Code = CODE_ERROR
		response.Msg = errUpdate.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	lId := questionInfo.AnswerId
	rankLog := new(model.RankInfoLogs)
	logUuid := util.GenerateUuid()
	rankLog.Uuid = logUuid
	rankLog.QuestionId = questionInfo.Uuid
	rankLog.LawyerId = lId
	rankLog.LawyerName = questionInfo.AnswerName
	rankLog.AskerName = questionInfo.CustomerName
	rankLog.AskerId = questionInfo.CustomerId
	rankLog.RankInfo = r
	today := time.Unix(time.Now().Unix(), 0).String()[0:19]
	rankLog.RankTime = today
	rankLog.RankPerson = "0"
	err := rankLog.GetConn().Create(&rankLog).Error
	if err != nil {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	lInfo := new(model.LawyerInfo)
	linfoErr := lInfo.GetConn().Where("uuid = ?", lId).Find(&lInfo).Error
	if linfoErr != nil && !strings.Contains(linfoErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = linfoErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	switch req.RankInfo {
	case 1:
		lInfo.RankLast = lInfo.RankLast + 1
	case 2:
		lInfo.RankFouth = lInfo.RankFouth + 1
	case 3:
		lInfo.RankThird = lInfo.RankThird + 1
	case 4:
		lInfo.RankSecond = lInfo.RankSecond + 1
	case 5:
		lInfo.RankFirst = lInfo.RankFirst + 1
	}
	lInfoUErr := lInfo.GetConn().Save(&lInfo).Error

	if lInfoUErr != nil {
		response.Code = CODE_ERROR
		response.Msg = lInfoUErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	setting := new(model.WechatVoiceQuestionSettings)
	settingErr := setting.GetConn().Where("category_id = ?", questionInfo.CategoryId).Find(&setting).Error
	if settingErr != nil && !strings.Contains(settingErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = settingErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	amount := setting.PayAmountInt
	// amount, _ := strconv.ParseInt(setting.PayAmountInt, 10, 64)
	lfp, _ := strconv.ParseInt(setting.LawyerFeePercent, 10, 64)
	bb := (100 - lfp) / 100
	cc := bb * amount
	urpp, _ := strconv.ParseInt(setting.UserRedPacketPercent, 10, 64)
	a := (rand.Int63n(urpp)) / 100
	packet := cc * a
	b := strconv.FormatInt(packet, 10)
	percent, _ := strconv.ParseInt(setting.LawyerFeePercent, 10, 64)
	fees := amount * percent / 100
	bf, _ := strconv.ParseFloat(b, 64)
	//记录余
	cost := new(model.MemberInfo)
	costErr := cost.GetConn().Where("open_id = ?", openId).Find(&cost).Error
	if costErr != nil {
		response.Code = CODE_ERROR
		response.Msg = costErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	balance := cost.Balance
	balanseF, _ := strconv.ParseFloat(balance, 64)
	balanseF = balanseF + bf
	balanseStr := strconv.FormatFloat(balanseF, 'f', 2, 64)
	cost.Balance = balanseStr
	updaerr := cost.GetConn().Save(&cost).Error
	if updaerr != nil {
		response.Code = CODE_ERROR
		response.Msg = updaerr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	//记录律师余额
	lawyer := new(model.LawyerInfo)
	lawyerErr := lawyer.GetConn().Where("open_id = ?", openId).Find(&lawyer).Error
	if lawyerErr != nil {
		response.Code = CODE_ERROR
		response.Msg = lawyerErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	balancel := lawyer.Balance
	balanseFl, _ := strconv.ParseFloat(balancel, 64)
	bbbb := float64(fees)
	balanseFl = balanseF + bbbb
	balanseStrl := strconv.FormatFloat(balanseFl, 'f', 2, 64)
	lawyer.Balance = balanseStrl
	updaerr1 := cost.GetConn().Save(&cost).Error
	if updaerr1 != nil {
		response.Code = CODE_ERROR
		response.Msg = updaerr1.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.Red = b
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

//检查问题是否被锁
type CheckIsLocked struct {
	QuestionId string `json:"questionId"`
}

func CheckAnswerIsLocked(ctx *macaron.Context) string {
	// cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	body, _ := ctx.Req.Body().String()
	req := new(CheckIsLocked)

	json.Unmarshal([]byte(body), req)

	questionInfo := new(model.WechatVoiceQuestions)
	response := new(model.GeneralResponse)

	questionInfoErr := questionInfo.GetConn().Where("uuid = ?", req.QuestionId).Find(&questionInfo).Error
	if questionInfoErr != nil && !strings.Contains(questionInfoErr.Error(), RNF) {
		response.Msg = questionInfoErr.Error()
		response.Code = CODE_ERROR
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if questionInfo.IsSolved == "1" {
		response.Msg = "问题已被解答"
		response.Code = CODE_ERROR
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	list, err := model.GetLockListById(req.QuestionId)
	if err != nil && !strings.Contains(err.Error(), RNF) {
		response.Msg = err.Error()
		response.Code = CODE_ERROR
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if len(list) == 2 {
		response.Code = CODE_ERROR
		response.Msg = "当前人数过多"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

func CreatePvInfo(ctx *macaron.Context) string {

	body, _ := ctx.Req.Body().String()
	req := new(CheckIsLocked)

	json.Unmarshal([]byte(body), req)
	questionInfo := new(model.WechatVoiceQuestions)
	response := new(model.GeneralResponse)
	questionInfoErr := questionInfo.GetConn().Where("uuid = ?", req.QuestionId).Find(&questionInfo).Error
	if questionInfoErr != nil && !strings.Contains(questionInfoErr.Error(), RNF) {
		fmt.Println(questionInfoErr.Error())
	}
	questionInfo.Pv = questionInfo.Pv + 1
	err := questionInfo.GetConn().Save(&questionInfo).Error
	if err != nil && !strings.Contains(err.Error(), RNF) {
		fmt.Println(err.Error())
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

type InitSpecail struct {
	LawyerId string `json:"lawyerId"`
	TypeId   string `json:"typeId"`
	OrderId  string `json:"orderId"`
}
type InitSpecailResponse struct {
	Code          int64  `json:"code"`
	Msg           string `json:"msg"`
	Name          string `json:"name"`
	SelfIntr      string `json:"selfIntr"`
	Pic           string `json:"pic"`
	TypePrice     string `json:"typePrice"`
	TypeId        string `json:"typeId"`
	ParentOrderId string `json:"parentOrderId"`
	TypeName      string `json:"typeName"`
}

func InitSpecialInfo(ctx *macaron.Context) string {
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	response := new(InitSpecailResponse)
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/initspecialinfo"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(res)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]
	fmt.Println(openId, userType)
	body, _ := ctx.Req.Body().String()
	req := new(InitSpecail)

	json.Unmarshal([]byte(body), req)
	if req.OrderId == "-1" {
		lawer := new(model.LawyerInfo)
		lawerErr := lawer.GetConn().Where("open_id = ?", req.LawyerId).Find(&lawer).Error
		if lawerErr != nil && !strings.Contains(lawerErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = lawerErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		cateInfo := new(model.Category)
		cateErr := cateInfo.GetConn().Where("uuid = ?", lawer.FirstCategory).Find(&cateInfo).Error
		if cateErr != nil && !strings.Contains(cateErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = cateErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		pay := new(model.WechatVoiceQuestionSettings)
		payErr := pay.GetConn().Where("category_id = ?", req.TypeId).Find(&pay).Error
		if payErr != nil && !strings.Contains(cateErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = payErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}

		response.Code = CODE_SUCCESS
		response.Msg = MSG_SUCCESS
		response.Name = lawer.Name
		response.SelfIntr = cateInfo.CategoryName
		response.Pic = lawer.HeadImgUrl
		response.TypePrice = pay.PayAmount
		response.TypeId = req.TypeId
		response.ParentOrderId = ""
		response.TypeName = cateInfo.CategoryName

	} else {
		//追问
		question := new(model.WechatVoiceQuestions)
		qErr := question.GetConn().Where("uuid = ?", req.OrderId).Find(&question).Error
		if qErr != nil && !strings.Contains(qErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = qErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		response.Code = CODE_SUCCESS
		response.Msg = MSG_SUCCESS
		response.Name = question.AnswerName
		response.Pic = question.AnswerHeadImg
		response.TypePrice = question.PaymentInfo
		response.TypeId = question.CategoryId
		response.TypeName = question.Category
		response.ParentOrderId = req.OrderId
		lawer := new(model.LawyerInfo)
		lawerErr := lawer.GetConn().Where("open_id = ?", question.AnswerOpenId).Find(&lawer).Error
		if lawerErr != nil && !strings.Contains(lawerErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = lawerErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		cateInfo := new(model.Category)
		cateErr := cateInfo.GetConn().Where("uuid = ?", lawer.FirstCategory).Find(&cateInfo).Error
		if cateErr != nil && !strings.Contains(cateErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = cateErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		response.SelfIntr = cateInfo.CategoryName

	}
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

// 初始化预支付界面请求结构体
type ReqInitPrePay struct {
	OrderId string `json:"orderId"`
}

// 初始化预支付界面响应结构体
type RespInitPrePay struct {
	Code       int64  `json:"code"`
	Msg        string `json:"msg"`
	OrderId    string `json:"orderId"`
	TotalMoney string `json:"totalMoney"`
}

func InitPay(ctx *macaron.Context) string {
	info := new(ReqInitPrePay)
	result := new(RespInitPrePay)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	//response := new(InitSpecailResponse)
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/initspecialinfo"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			result.Code = CODE_ERROR
			result.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(result)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				result.Code = CODE_ERROR
				result.Msg = err.Error()
				ret_str, _ := json.Marshal(result)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	openId := strings.Split(cookieStr, "|")[0]
	fmt.Println(openId)
	userType := strings.Split(cookieStr, "|")[1]
	fmt.Println(userType)
	reqData, _ := ctx.Req.Body().String()
	err := json.Unmarshal([]byte(reqData), info)
	if err != nil {
		//log.Error("[errorInfo]: error when unmarshal request body")
		result.Code = CODE_ERROR
		result.Msg = "json解析异常"
	} else {

		result.Code = CODE_SUCCESS
		result.Msg = "ok"
		orderInfo := new(model.WechatVoiceQuestions)
		orderErr := orderInfo.GetConn().Where("order_number = ?", info.OrderId).Find(&orderInfo).Error
		if orderErr != nil && !strings.Contains(orderErr.Error(), RNF) {
			result.Code = CODE_ERROR
			result.Msg = orderErr.Error()
			ret_str, _ := json.Marshal(result)
			return string(ret_str)
		}

		result.OrderId = orderInfo.OrderNumber
		result.TotalMoney = orderInfo.PaymentInfo
	}
	ret_str, _ := json.Marshal(result)
	return string(ret_str)

}

// 进行支付请求
type ReqDoPay struct {
	OrderId string `json:"orderId"`
	Type    string `json:"type"` // 0 余额支付   1 微信支付   2 货到付款
}

// 进行支付响应
type RespDoPay struct {
	Code           int64  `json:"code"`
	Msg            string `json:"msg"`
	Type           string `json:"type"`
	JumpSuccessUrl string `json:"paySuccess"`
	JumpFailedUrl  string `json:"payFailed"`
	JumpSubmitUrl  string `json:"submitSuccess"`
	Timestamp      int64  `json:"timestamp"`
	NonceStr       string `json:"nonceStr"`
	Package        string `json:"package"`
	SignType       string `json:"signType"`
	PaySign        string `json:"paySign"`
	AppId          string `json:"appId"`
	ConfigSign     string `json:"configSign"`
}

func DoPayNew(ctx *macaron.Context) string {
	result := new(RespDoPay)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	//response := new(InitSpecailResponse)
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/initspecialinfo"
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
		//cookieStr = "1|2"
		ctx.Redirect(url)
	}
	code := ctx.Query("code")
	if code != "" {
		url := "http://60.205.4.26:22334/getOpenid?code=" + code
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("=========xxxxx")
			fmt.Println(err.Error())
		}
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(resBody))
		defer res.Body.Close()
		fmt.Println("==========>>>>")
		res1 := new(OpenIdResponse)
		json.Unmarshal(resBody, res1)
		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
		member := new(model.MemberInfo)
		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			result.Code = CODE_ERROR
			result.Msg = memberErr.Error()
			ret_str, _ := json.Marshal(result)
			return string(ret_str)
		}
		if member.Uuid == "" {
			fmt.Println("新的用户")
			user := GetUserInfo(res1.OpenId, res1.AccessToken)
			member.Uuid = util.GenerateUuid()
			member.HeadImgUrl = user.HeadImgUrl
			member.OpenId = user.OpenId
			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			member.NickName = user.NickName
			err := member.GetConn().Create(&member).Error
			if err != nil {
				result.Code = CODE_ERROR
				result.Msg = err.Error()
				ret_str, _ := json.Marshal(result)
				return string(ret_str)
			}
		}
		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	}
	fmt.Println(cookieStr)
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]
	fmt.Println(openId, userType)
	info := new(ReqDoPay)
	reqData, _ := ctx.Req.Body().String()
	fmt.Println("===============================do pay new 中获取到的参数===============================")
	fmt.Println(reqData)
	fmt.Println("===============================do pay new 中获取到的参数===============================")
	// 解析请求体
	orderInfo := new(model.WechatVoiceQuestions)
	orderErr := orderInfo.GetConn().Where("order_number = ?", info.OrderId).Find(&orderInfo).Error
	if orderErr != nil && !strings.Contains(orderErr.Error(), RNF) {
		fmt.Println(orderErr.Error())
	}
	err := json.Unmarshal([]byte(reqData), info)
	if err != nil {
		// log.Error("[errorInfo]: error when unmarshal request body")
		result.Code = CODE_ERROR
		result.Msg = "json解析异常"
	} else {
		appId := ""
		mchId := ""
		merchantName := ""
		orderNumber := info.OrderId
		priceInt := orderInfo.PaymentInfo
		serverIp := "127.0.0.1"
		key := ""
		payPayge := ""
		afterPay := "http://www.mylvfa.com/voice/afterpay"
		url := "http://www.mylvfa.com/voice/uni?appid=" + appId + "&mchid=" + mchId + "&name=" + merchantName + "&ordernumber=" + orderNumber + "&price=" + priceInt + "&serverIp=" + serverIp + "&key=" + key + "&payPayge=" + payPayge + "&afterpay=" + afterPay
		res, reserr := http.Get(url)
		if reserr != nil {
			fmt.Println(reserr)
		}
		defer res.Body.Close()
		HTTPResult, _ := ioutil.ReadAll(res.Body)
		res111 := new(UnifiedOrderResp)
		json.Unmarshal(HTTPResult, res111)
		result.Code = CODE_SUCCESS
		result.Msg = "成功!"
		result.Type = info.Type
		result.JumpFailedUrl = AFTER_PAY_JUMP_PAGE_FAILD + info.OrderId
		result.JumpSuccessUrl = AFTER_PAY_JUMP_PAGE_SUCCESS + info.OrderId
		result.NonceStr = res111.NonceStr
		result.Package = res111.Package
		result.PaySign = res111.PaySign
		result.SignType = res111.SignType
		result.AppId = res111.AppId
		result.ConfigSign = res111.ConfigSign
		timeStamp, _ := strconv.ParseInt(res111.TimeStamp, 10, 64)
		result.Timestamp = timeStamp
		//这里说明 是支付成功的 说明开始进入请求福记接口阶段
		fmt.Println("预支付成功")
	}
	ret_str, _ := json.Marshal(result)
	return string(ret_str)
}

type UnifiedOrderResp struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`

	AppId      string `json:"appId"`
	PrepayId   string `json:"prepayId"`
	CodeUrl    string `json:"codeUrl"`
	TimeStamp  string `json:"timestamp"`
	NonceStr   string `json:"nonceStr"`
	Package    string `json:"package"`
	SignType   string `json:"MD5"`
	PaySign    string `json:"paySign"`
	ConfigSign string `json:"configSign"`
}

// 微信预支付推送结构体
type WechatRespUnifiedOrder struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	Appid      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	TradeType  string `xml:"trade_type"`
	PrepayId   string `xml:"prepay_id"`
	CodeUrl    string `xml:"code_url"`
}

func UniFi(ctx *macaron.Context) string {
	result := new(UnifiedOrderResp)
	paramsMap := make(map[string]string, 0)
	paramsList := []string{"appid", "mch_id", "body", "out_trade_no", "total_fee", "spbill_create_ip", "device_info", "nonce_str", "fee_type", "time_start", "notify_url", "trade_type"}

	// 获取请求中需要带的必要信息
	Appid := ctx.Query("appid")
	MchId := ctx.Query("mchid")
	Body := ctx.Query("name")
	OutTradeNo := ctx.Query("ordernumber")
	TotalFee := ctx.Query("price")
	SpbillCreateIp := ctx.Query("serverIp")
	Key := ctx.Query("key")
	PageUrl := ctx.Query("payPage")
	if Appid == "" || MchId == "" || Body == "" || OutTradeNo == "" || TotalFee == "" || SpbillCreateIp == "" || Key == "" || PageUrl == "" {
		result.Code = 400
		result.Msg = "参数不全"
	} else {
		// 所有必须字段不为空,加入必须字段到map中
		paramsMap["appid"] = Appid
		paramsMap["mch_id"] = MchId
		paramsMap["body"] = Body
		paramsMap["out_trade_no"] = OutTradeNo
		paramsMap["total_fee"] = TotalFee
		paramsMap["spbill_create_ip"] = SpbillCreateIp

		// 有默认值的字段处理
		paramsMap["device_info"] = DEFAULT_DEVICE_INFO
		paramsMap["nonce_str"] = util.GenerateUuid()
		paramsMap["fee_type"] = DEFAULT_FEE_TYPE
		paramsMap["time_start"] = time.Now().Format("20060102150405")
		paramsMap["notify_url"] = DEFAULT_NOTIFY_URL
		paramsMap["trade_type"] = DEFAULT_TRADE_TYPE

		sign := GenerateSign(paramsMap, paramsList, ctx.Query("key"))
		paramsMap["sign"] = sign

		xmlStr := GenerateXMLStr(paramsMap)

		HTTPBody := bytes.NewBuffer([]byte(xmlStr))

		res1, err1 := http.Post(UNIFIEDORDER_URL, "application/xml", HTTPBody)
		if err1 != nil {
			result.Code = CODE_ERROR
			result.Msg = err1.Error()
			ret_str, _ := json.Marshal(result)
			return string(ret_str)
		} else {
			HTTPResult, err := ioutil.ReadAll(res1.Body)
			defer res1.Body.Close()
			if err != nil {
				// log.Error("[UnifiedOrder]:error when read responce body:" + err.Error())
				result.Code = CODE_ERROR
				result.Msg = "读取返回体错误!"
				ret_str, _ := json.Marshal(result)
				return string(ret_str)
			} else {
				var wechatResult WechatRespUnifiedOrder
				fmt.Println("=========>>>>>>>>>>")
				fmt.Println(string(HTTPResult))
				fmt.Println("=========>>>>>>>>>>")
				err := xml.Unmarshal(HTTPResult, &wechatResult)
				if err != nil {
					//log.Error("[UnifiedOrder]:error when unmarshal http result body:" + err.Error())
					result.Code = CODE_ERROR
					result.Msg = "解析返回体错误!"
				} else if strings.ToUpper(wechatResult.ResultCode) != "SUCCESS" {
					result.Code = CODE_ERROR
					result.Msg = wechatResult.ReturnMsg
				} else if strings.ToUpper(wechatResult.ReturnCode) != "SUCCESS" {
					result.Code = CODE_ERROR
					result.Msg = "errCode:" + wechatResult.ErrCode + "errMsg:" + wechatResult.ErrCodeDes
				} else {

					// 开始签名前端页面发起支付所用参数
					prepayMap := make(map[string]string, 0)
					prepayList := []string{"appId", "timeStamp", "nonceStr", "package", "signType"}

					nTimeStr := strconv.FormatInt(time.Now().Unix(), 10)
					newNonceStr := util.GenerateUuid()

					prepayMap["appId"] = Appid
					prepayMap["timeStamp"] = nTimeStr
					prepayMap["nonceStr"] = newNonceStr
					prepayMap["package"] = DEFAULT_PACKAGE_PRE_STR + wechatResult.PrepayId
					prepayMap["signType"] = DEFAULT_SIGN_TYPE

					prepaySign := GenerateSign(prepayMap, prepayList, ctx.Query("key"))

					rest, errT := http.Get("http://www.mylvfa.com/voice/ticket?appid=wxac69efc11c5e182f")
					if errT != nil {
						result.Code = CODE_ERROR
						result.Msg = errT.Error()
						ret_Str, _ := json.Marshal(result)
						return string(ret_Str)
					} else {
						defer rest.Body.Close()
						resBody2, _ := ioutil.ReadAll(rest.Body)
						tic := new(RespJsapiTicket)
						err11 := json.Unmarshal(resBody2, tic)
						if err11 != nil {
							result.Code = CODE_ERROR
							result.Msg = err11.Error()
							ret_str, _ := json.Marshal(result)
							return string(ret_str)
						} else {
							configMap := make(map[string]string, 0)
							configList := []string{"jsapi_ticket", "timestamp", "noncestr", "url"}
							configMap["jsapi_ticket"] = tic.Ticket
							configMap["timestamp"] = nTimeStr
							configMap["noncestr"] = newNonceStr
							configMap["url"] = PageUrl
							configSign := GeneratePageSign(configMap, configList)
							result.Code = CODE_SUCCESS
							result.Msg = "SUCCESS"
							result.CodeUrl = wechatResult.CodeUrl
							result.NonceStr = newNonceStr
							result.Package = DEFAULT_PACKAGE_PRE_STR + wechatResult.PrepayId
							result.PaySign = prepaySign
							result.PrepayId = wechatResult.PrepayId
							result.SignType = DEFAULT_SIGN_TYPE
							result.TimeStamp = nTimeStr
							result.AppId = Appid
							result.ConfigSign = configSign
						}

					}

				}
			}

		}
	}
	resByte, _ := json.Marshal(result)
	fmt.Println(string(resByte))
	return string(resByte)
}

func SendRedPacketToLaw() string {
	return ""
}

type ResponsePay struct {
	Sign string `json:"sign"`
}

// <xml>
// 	<return_code><![CDATA[SUCCESS]]></return_code>
// 	<return_msg><![CDATA[OK]]></return_msg>
// 	<appid><![CDATA[wxac69efc11c5e182f]]></appid>
// 	<mch_id><![CDATA[1344737201]]></mch_id>
// 	<nonce_str><![CDATA[aMDA4RftWtlZXt9N]]></nonce_str>
// 	<sign><![CDATA[156CF9C13F8F85E6FB89A9958D97DC6D]]></sign>
// 	<result_code><![CDATA[SUCCESS]]></result_code>
// 	<prepay_id><![CDATA[wx20160823013657f0f64bfe2d0925319748]]></prepay_id>
// 	<trade_type><![CDATA[JSAPI]]></trade_type>
// </xml>

type PaySignResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	Appid      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	PrepayId   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
}

type PayFinal struct {
	Sign string `json:"sign"`
}

func PayBill(nstr, nSt, openId, orderNumber, fee, timeStamp string) (string, string, string, error) {
	// var nstr string
	// var openId string
	// var orderNumber string
	// var fee string
	var sign string
	var prepayId string
	var sings string
	url := "http://60.205.4.26:22334/prepayId?appid=wxac69efc11c5e182f&mch_id=1344737201&nonce_str=" + nstr + "&notify_url=http://www.mylvfa.com/voice/front/afterpay&openid=" + openId + "&out_trade_no=" + orderNumber + "&spbill_create_ip=127.0.0.1&total_fee=" + fee + "&trade_type=JSAPI&body=my_pay_test"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error(), "1")
		return sign, prepayId, sings, err
	} else {
		body, bodyErr := ioutil.ReadAll(res.Body)
		if bodyErr != nil {
			fmt.Println("bodyerr ", bodyErr.Error(), "2")
			return sign, prepayId, sings, bodyErr
		}
		responseSign := new(ResponsePay)
		unmarErr := json.Unmarshal(body, responseSign)
		if unmarErr != nil {
			fmt.Println(unmarErr.Error(), "3")
			return sign, prepayId, sings, unmarErr
		} else {
			sign = responseSign.Sign
			//sings = responseSign
			fmt.Println(sign)

			resp := new(PaySignResponse)
			unmarErr = xml.Unmarshal([]byte(sign), resp)
			if unmarErr != nil {
				fmt.Println(unmarErr.Error(), "4")
				return sign, prepayId, sings, unmarErr
			} else {
				prepayId = resp.PrepayId
				sings = resp.Sign
				//var nSt string
				url1 := "http://60.205.4.26:22334/prepaySign?appId=wxac69efc11c5e182f&nonceStr=" + nstr + "&package=prepay_id=" + prepayId + "&signType=MD5&timeStamp=" + timeStamp
				res2, res2err := http.Get(url1)
				if res2err != nil {
					fmt.Println(res2err.Error(), "5")
					return sign, prepayId, sings, res2err
				} else {
					r := new(PayFinal)
					bodyF, errF := ioutil.ReadAll(res2.Body)
					if errF != nil {
						fmt.Println(errF.Error(), "6")
						return sign, sings, prepayId, errF
					} else {
						json.Unmarshal(bodyF, r)
						sign = r.Sign
					}
				}
			}
		}
	}
	fmt.Println(sign, prepayId, sings)
	return sign, prepayId, sings, nil
}
func GetSigns(timeStr string) string {
	// signs := time.Now().Unix()
	// signsStr := strconv.FormatInt(signs, 10)
	url := "http://60.205.4.26:22334/configSign?noncestr=W1471365761W&timestamp=" + timeStr + "&url=http://www.mylvfa.com/daodaolaw/search.html"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer res.Body.Close()
	resa := new(AResponse)
	resBody, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(resBody, resa)
	fmt.Println(string(resBody))
	return resa.Sign
}

type ConfigResponssss struct {
	Code      int64  `json:"code"`
	Msg       string `json:"msg"`
	Appid     string `json:"appId"`
	TimeStamp string `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Sing      string `json:"signature"`
}

func GetWxVoiceConfig(ctx *macaron.Context) string {
	ticker := JsapiTicker12()
	nstr := util.GenerateUuid()

	timeStamp := time.Now().Format("20060102150405")
	str := "jsapi_ticket=" + ticker + "&noncestr=" + nstr + "&timestamp=" + timeStamp + "&url=http://mp.weixin.qq.com?params=value"
	str = Sha1(str)
	response := new(ConfigResponssss)
	response.Code = CODE_SUCCESS
	response.Msg = "ok"
	response.Appid = "wxac69efc11c5e182f"
	response.TimeStamp = timeStamp
	response.NonceStr = nstr
	response.Sing = str
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

type UnifiedOrderResps struct {
	Head struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
	} `json:"head"`
	Body struct {
		AppId      string `json:"appId"`
		PrepayId   string `json:"prepayId"`
		CodeUrl    string `json:"codeUrl"`
		TimeStamp  string `json:"timestamp"`
		NonceStr   string `json:"nonceStr"`
		Package    string `json:"package"`
		SignType   string `json:"MD5"`
		PaySign    string `json:"paySign"`
		ConfigSign string `json:"configSign"`
	} `json:"body"`
}

// func UnifieOrder(appid, mchId, body, outNo, total, ip, deviceInfo, nonceStr, feetype, times, notiu, tradetype string) UnifiedOrderResps {
// 	paramsMap := make(map[string]string, 0)
// 	paramsList := []string{"appid", "mch_id", "body", "out_trade_no", "total_fee", "spbill_create_ip", "device_info", "nonce_str", "fee_type", "time_start", "notify_url", "trade_type"}
// 	paramsMap["appid"] = appid
// 	paramsMap["mch_id"] = mchId
// 	paramsMap["body"] = body
// 	paramsMap["out_trade_no"] = outNo
// 	paramsMap["total_fee"] = total
// 	paramsMap["spbill_create_ip"] = ip

// 	// 有默认值的字段处理
// 	paramsMap["device_info"] = deviceInfo
// 	paramsMap["nonce_str"] = nonceStr
// 	paramsMap["fee_type"] = feetype
// 	paramsMap["time_start"] = times
// 	paramsMap["notify_url"] = notiu
// 	paramsMap["trade_type"] = tradetype
// }
type PekReq struct {
	OrderId string `json:"orderId"`
}

type PeekResponses struct {
	Code      int64  `json:"code"`
	Msg       string `json:"msg"`
	Appid     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
	OrderId   string `json:"orderId"`
}

func PayPeekAnswer(ctx *macaron.Context) string {
	req := new(PekReq)
	body, _ := ctx.Req.Body().String()
	openId := ""
	response := new(PeekResponses)

	json.Unmarshal([]byte(body), req)
	nstr := util.GenerateUuid()
	nSt := util.GenerateUuid()
	timeStamp := time.Now().Unix()
	fmt.Println(timeStamp)
	tStr := strconv.FormatInt(timeStamp, 10)

	orderNumber := util.GenerateOrderNumber()
	orderInfo := new(model.WechatVoiceQuestions)
	orderInfoErr := orderInfo.GetConn().Where("order_number = ?", req.OrderId).Find(&orderInfo).Error

	if orderInfoErr != nil && !strings.Contains(orderInfoErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = orderInfoErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	signSelf := GetSigns(tStr)
	sign, prepayId, sings, signErr := PayBill(nstr, nSt, openId, orderNumber, "1", tStr)
	fmt.Println(sings)
	if signErr != nil {
		fmt.Println(signErr.Error())
	}
	response.Code = CODE_SUCCESS
	//response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.Appid = "wxac69efc11c5e182f"
	response.NonceStr = nstr
	response.Signature = signSelf
	response.SignType = "MD5"
	response.Package = "prepay_id=" + prepayId
	response.TimeStamp = tStr
	response.PaySign = sign
	response.OrderId = orderNumber
	ret_str, _ := json.Marshal(response)
	fmt.Println("=======================>>>")
	fmt.Println(string(ret_str))
	fmt.Println("=======================>>>>")
	return string(ret_str)
}

type AfterPayInfo struct {
	Appid         string `xml:"appid"`
	BankType      string `xml:"bank_type"`
	CashFee       string `xml:"cash_fee"`
	FeeType       string `xml:"fee_type"`
	IsSubscribe   string `xml:"is_subscribe"`
	MchId         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	OpenId        string `xml:"openId"`
	OutTradeNum   string `xml:"out_trade_no"`
	ResultCode    string `xml:"result_code"`
	ReturnCode    string `xml:"return_code"`
	Sign          string `xml:"sign"`
	TimeEnd       string `xml:"time_end"`
	TotalFee      string `xml:"total_fee"`
	TradeType     string `xml:"trade_type"`
	TransactionId string `xml:"transaction_id"`
}
type AfterPayRespToWechat struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

func AfterPay(ctx *macaron.Context) string {
	req, _ := ctx.Req.Body().String()
	fmt.Println("===========")
	fmt.Println(req)
	fmt.Println("===========")
	a := new(AfterPayInfo)
	unmarErr := xml.Unmarshal([]byte(req), a)
	if unmarErr != nil {
		fmt.Println("=====>>>", unmarErr.Error())
	}
	response := new(AfterPayRespToWechat)
	if a.ResultCode == "SUCCESS" {
		fmt.Println("支付回调成功")
		//修改订单状态
		orderNumber := a.OutTradeNum
		order := new(model.WechatVoiceQuestions)
		orderErr := order.GetConn().Where("order_number =?", orderNumber).Find(&order).Error
		if orderErr != nil && !strings.Contains(orderErr.Error(), RNF) {
			fmt.Println(orderErr.Error())
		}
		order.IsPaied = "1"
		// order.
		orderErr = order.GetConn().Save(&order).Error
		if orderErr != nil && !strings.Contains(orderErr.Error(), RNF) {
			fmt.Println("update err", orderErr.Error())
		}
		pay := new(model.WechatVoicePaymentInfo)
		pay.Uuid = util.GenerateUuid()
		pay.SwiftNumber = a.TransactionId
		pay.QuestionId = order.Uuid
		pay.OpenId = a.OpenId
		pay.OrderId = a.OutTradeNum
		payErr := pay.GetConn().Save(&pay).Error

		if payErr != nil && !strings.Contains(payErr.Error(), RNF) {
			log.Println(payErr.Error())
		}
		/*
			type WechatVoicePaymentInfo struct {
				gorm.Model
				Uuid            string
				SwiftNumber     string
				QuestionId      string
				MemberId        string
				OpenId          string
				RedPacketAmount string
				LawyerAmount    string
				Left            string
				OrderId         string
			}
					**/
	} else {
		fmt.Println("失败")
		//response
	}
	response.ReturnCode = "SUCCESS"
	response.ReturnMsg = "OK"
	ret_str, _ := xml.Marshal(response)
	return string(ret_str)
}

type GetInfo struct {
	OrderId string `json:"orderId"`
}
type GetInfoResponse struct {
	Code          int64  `json:"code"`
	Msg           string `json:"msg"`
	Name          string `json:"name"`
	SelfIntr      string `json:"selfIntr"`
	Pic           string `json:"pic"`
	TypePrice     string `json:"typePrice"`
	TypeId        string `json:"typeId"`
	TypeName      string `json:"typeName"`
	ParentOrderId string `json:"parentOrderId"`
}

func GetOrderInfoById(ctx *macaron.Context) string {
	body, _ := ctx.Req.Body().String()
	req := new(GetInfo)
	json.Unmarshal([]byte(body), req)
	response := new(GetInfoResponse)
	order := new(model.WechatVoiceQuestions)
	orderErr := order.GetConn().Where("uuid = ?", req.OrderId).Find(&order).Error
	if orderErr != nil && !strings.Contains(orderErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = orderErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	// lId := order.AnswerId
	law := new(model.LawyerInfo)
	lc := new(model.LawCatgory)
	ctSet := new(model.WechatVoiceQuestionSettings)
	err := law.GetConn().Where("uuid = ?", order.AnswerId).Find(&law).Error
	err = lc.GetConn().Where("uuid = ?", order.CategoryId).Find(&lc).Error
	err = ctSet.GetConn().Where("category_id = ?", order.CategoryId).Find(&ctSet).Error
	if err != nil && !strings.Contains(err.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	response.Code = CODE_SUCCESS
	response.Msg = "ok"
	response.Name = law.Name
	response.SelfIntr = law.FirstCategory
	response.Pic = law.HeadImgUrl
	response.TypePrice = ctSet.PayAmount
	response.TypeId = order.CategoryId
	response.TypeName = lc.CategoryName
	response.ParentOrderId = ""
	ret_str, _ := json.Marshal(response)
	// fmt.Print(string)
	log.Println(string(ret_str))
	return string(ret_str)
}

// func GetOrderDetailById(ctx)
