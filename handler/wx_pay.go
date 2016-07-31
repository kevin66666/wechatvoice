package handler

import (
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"strings"
	"wechatvoice/model"
	"wechatvoice/tool/util"
	"strconv"
	"fmt"
	"net/http"
	"time"
	"bytes"
	"log"
	"github.com/Unknwon/macaron"
	"github.com/henrylee2cn/teleport/example"
)

func DoWechatPay(ctx *macaron.Context)string{
	info := new(ReqDoPay)
	result := new(RespDoPay)

	finalResult := make(map[string]interface{}, 0)
	head := make(map[string]interface{}, 0)
	finalResult["head"] = head
	finalResult["body"] = result

	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")

	if cookieStr==""{
		//这里直接调取util重新过一次绿叶 获取openId 等信息
	}
	openId :=strings.Split(cookieStr,"|")[0]
	userType :=strings.Split(cookieStr,"|")[1]

	log.Println(openId)
	log.Println(userType)

	memberInfo :=new(model.MemberInfo)
	memberErr :=memberInfo.GetConn().Where("open_id = ?",openId).Find(&memberInfo).Error

	if memberErr!=nil&&!strings.Contains(memberErr.Error(),RNF){
		result.Code = CODE_ERROR
		result.Msg = memberErr.Error()
		ret_str,_:=json.Marshal(result)
		return string(ret_str)
	}

	reqData, _ := ctx.Req.Body().String()

	// 解析请求体
	err := json.Unmarshal([]byte(reqData), info)
	if err != nil {
		log.Println("[errorInfo]: error when unmarshal request body")
		result.Code = CODE_ERROR
		result.Msg = "json解析异常"
	}else{
		order :=new(model.WechatVoiceQuestions)

		orderErr :=order.GetConn().Where("order_number = ?",info.OrderId).Find(&order).Error

		if orderErr!=nil&&!strings.Contains(orderErr.Error(),RNF){
			result.Code = CODE_ERROR
			result.Msg = orderErr.Error()
			ret_str,_:=json.Marshal(result)
			return string(ret_str)
		}
		price :=order.PaymentInfo
		priceF :=strconv.FormatFloat(price,'f',2,64)
		return doWxPay(info.OrderId,"1",priceF)
	}
	resByte, _ := json.Marshal(finalResult)
	return string(resByte)
}

func doWxPay(orderId string ,payType,price float64)string{
	result := make(map[string]interface{}, 0)
	head := make(map[string]interface{}, 0)
	body := make(map[string]interface{}, 0)
	result["head"] = head
	result["body"] = body
	orderInfo :=new(model.WechatVoiceQuestions)

	orderInfoErr :=orderInfo.GetConn().Where("order_number = ?",orderId).Find(&orderInfo).Error
	if orderInfoErr!=nil&&!strings.Contains(orderInfoErr.Error(),RNF){
		head["code"] = CODE_ERROR
		head["msg"] = "支付信息初始化失败!"
	}
	priceInt := int64(price)
	//WECHAT_PREPAY_URL             = "/shangqu-3rdparty/pay/unifiedorder?appid=%s&mch_id=%s&body=%s&out_trade_no=%s&total_fee=%d&spbill_create_ip=%s&key=%s&openid=%s&url=%s&notify_url=%s"
	url := fmt.Sprintf(WECHAT_PREPAY_URL,APPID, MCHID, MCHNAME, orderId, priceInt, SERVER_IP, KEY, orderInfo.CustomerOpenId, PAY_PAGE_URL, AFTER_PAY_ORDER_URL)

	res, err := http.Get(url)
	if err != nil {
		log.Println("[errorInfo][DoPay]:error when do prepay:" + err.Error())
		head["code"] = CODE_ERROR
		head["msg"] = "支付信息初始化失败!"
	} else {
		HTTPResult, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			log.Println("[errorInfo][DoPay]:error when read server resp:" + err.Error())
			head["code"] = CODE_ERROR
			head["msg"] = "支付信息初始化失败!"
		} else {
			serverResp := new(UnifiedOrderResp)
			json.Unmarshal(HTTPResult, serverResp)
			if serverResp.Head.Code != 200 {
				log.Println("[errorInfo][DoPay]:error when prepay order:msg is :" + serverResp.Head.Msg)
				log.Println("[errorInfo][DoPay]:error when prepay order:request is :" + string(HTTPResult))
				log.Println("[errorInfo][DoPay]:error when prepay order:url is :" + url)
				head["code"] = CODE_ERROR
				head["msg"] = "预支付失败!"
			} else {
				head["code"] = CODE_SUCCESS
				head["msg"] = "成功!"

				body["type"] = payType
				body["payFailed"] = AFTER_PAY_JUMP_PAGE_FAILD + orderId
				body["paySuccess"] = AFTER_PAY_JUMP_PAGE_SUCCESS + orderId
				body["nonceStr"] = serverResp.Body.NonceStr
				body["package"] = serverResp.Body.Package
				body["paySign"] = serverResp.Body.PaySign
				body["signType"] = serverResp.Body.SignType
				body["appId"] = serverResp.Body.AppId
				body["configSign"] = serverResp.Body.ConfigSign
				timeStamp, _ := strconv.ParseInt(serverResp.Body.TimeStamp, 10, 64)
				body["timestamp"] = timeStamp
			}
		}
	}
	resByte, _ := json.Marshal(result)
	return string(resByte)
}

// 预支付返回结构体
type UnifiedOrderResp struct {
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
type TicketServerResp struct {
	Head struct {
		     Code int64  `json:"code"`
		     Msg  string `json:"msg"`
	     } `json:"head"`
	Body struct {
		     Ticket      string `json:"jsapiTicket"`
		     Accesstoken string `json:"accesstoken"`
	     } `json:"body"`
}

func UnifiedOrder(ctx *macaron.Context)string{
	result := new(UnifiedOrderResp)
	paramsMap := make(map[string]string, 0)
	paramsList := []string{"appid", "mch_id", "body", "out_trade_no", "total_fee", "spbill_create_ip", "device_info", "nonce_str", "fee_type", "time_start", "notify_url", "trade_type"}

	// 获取请求中需要带的必要信息
	Appid := ctx.Query("appid")
	MchId := ctx.Query("mch_id")
	Body := ctx.Query("body")
	OutTradeNo := ctx.Query("out_trade_no")
	TotalFee := ctx.Query("total_fee")
	SpbillCreateIp := ctx.Query("spbill_create_ip")
	Key := ctx.Query("key")
	PageUrl := ctx.Query("url")

	// 如果有空的字段,直接返回
	if Appid == "" || MchId == "" || Body == "" || OutTradeNo == "" || TotalFee == "" || SpbillCreateIp == "" || Key == "" || PageUrl == "" {
		result.Head.Code = 400
		result.Head.Msg = "参数不全"
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

		if ctx.Query("device_info") != "" {
			paramsMap["device_info"] = ctx.Query("device_info")
		}
		if ctx.Query("nonce_str") != "" {
			paramsMap["nonce_str"] = ctx.Query("nonce_str")
		}
		if ctx.Query("fee_type") != "" {
			paramsMap["fee_type"] = ctx.Query("fee_type")
		}
		if ctx.Query("time_start") != "" {
			paramsMap["time_start"] = ctx.Query("time_start")
		}
		if ctx.Query("notify_url") != "" {
			paramsMap["notify_url"] = ctx.Query("notify_url")
		}
		if ctx.Query("trade_type") != "" {
			paramsMap["trade_type"] = ctx.Query("trade_type")
		}

		if ctx.Query("detail") != "" {
			paramsMap["detail"] = ctx.Query("detail")
			paramsList = append(paramsList, "detail")
		}
		if ctx.Query("attach") != "" {
			paramsMap["attach"] = ctx.Query("attach")
			paramsList = append(paramsList, "attach")
		}
		if ctx.Query("time_expire") != "" {
			paramsMap["time_expire"] = ctx.Query("time_expire")
			paramsList = append(paramsList, "time_expire")
		}
		if ctx.Query("goods_tag") != "" {
			paramsMap["goods_tag"] = ctx.Query("goods_tag")
			paramsList = append(paramsList, "goods_tag")
		}
		if ctx.Query("product_id") != "" {
			paramsMap["product_id"] = ctx.Query("product_id")
			paramsList = append(paramsList, "product_id")
		}
		if ctx.Query("limit_pay") != "" {
			paramsMap["limit_pay"] = ctx.Query("limit_pay")
			paramsList = append(paramsList, "limit_pay")
		}
		if ctx.Query("openid") != "" {
			paramsMap["openid"] = ctx.Query("openid")
			paramsList = append(paramsList, "openid")
		}
		sign := util.GenerateSign(paramsMap, paramsList, ctx.Query("key"))
		paramsMap["sign"] = sign
		xmlStr := util.GenerateXMLStr(paramsMap)
		HTTPBody := bytes.NewBuffer([]byte(xmlStr))

		res, err := http.Post(UNIFIEDORDER_URL, "application/xml", HTTPBody)
		if err != nil {
			log.Println("[UnifiedOrder]:error when post to " + UNIFIEDORDER_URL + ":" + err.Error())
			result.Head.Code = CODE_ERROR
			result.Head.Msg = "HTTP通信错误!"
		} else {
			HTTPResult, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err!=nil{
				log.Println("[UnifiedOrder]:error when read responce body:" + err.Error())
				result.Head.Code = CODE_ERROR
				result.Head.Msg = "读取返回体错误!"
			}else{
				// 开始签名前端页面发起支付所用参数
				var wechatResult WechatRespUnifiedOrder
				err := xml.Unmarshal(HTTPResult, &wechatResult)
				if err != nil {
					log.Println("[UnifiedOrder]:error when unmarshal http result body:" + err.Error())
					result.Head.Code = CODE_ERROR
					result.Head.Msg = "解析返回体错误!"
				} else if strings.ToUpper(wechatResult.ReturnCode) != "SUCCESS" {
					log.Println("[UnifiedOrder]:error in wechant server ...")
					log.Println("[UnifiedOrder]:" + string(HTTPResult))
					result.Head.Code = CODE_ERROR
					result.Head.Msg = wechatResult.ReturnMsg
				} else if strings.ToUpper(wechatResult.ResultCode) != "SUCCESS" {
					log.Println("[UnifiedOrder]:error in wechant server ...")
					log.Println("[UnifiedOrder]:" + string(HTTPResult))
					result.Head.Code = CODE_ERROR
					result.Head.Msg = "errCode:" + wechatResult.ErrCode + "errMsg:" + wechatResult.ErrCodeDes
				} else{
					prepayMap := make(map[string]string, 0)
					prepayList := []string{"appId", "timeStamp", "nonceStr", "package", "signType"}

					nTimeStr := strconv.FormatInt(time.Now().Unix(), 10)
					newNonceStr := util.GenerateUuid()

					prepayMap["appId"] = Appid
					prepayMap["timeStamp"] = nTimeStr
					prepayMap["nonceStr"] = newNonceStr
					prepayMap["package"] = DEFAULT_PACKAGE_PRE_STR + wechatResult.PrepayId
					prepayMap["signType"] = DEFAULT_SIGN_TYPE

					prepaySign := util.GenerateSign(prepayMap, prepayList, ctx.Query("key"))

					// 生成页面配置签名(参与签名的字段包括参与签名的字段包括noncestr, 有效的jsapi_ticket, timestamp, url)
					// 首先获取当前商户的JSAPI-TICKET
					res, err := http.Get(TICKET_SERVER_URL + Appid)
					if err != nil {
						log.Println("[UnifiedOrder]:error when get jsapi-ticket:" + err.Error())
						result.Head.Code = CODE_ERROR
						result.Head.Msg = "生成JSAPI错误!"
					} else {

						SERVERResult, err := ioutil.ReadAll(res.Body)
						defer res.Body.Close()
						if err != nil {
							log.Println("[UnifiedOrder]:error when read responce body:" + err.Error())
							result.Head.Code = CODE_ERROR
							result.Head.Msg = "读取返回体错误!"
						} else {
							var serverResult TicketServerResp
							err := json.Unmarshal(SERVERResult, &serverResult)
							if err != nil {
								log.Println("[UnifiedOrder]:error when unmarshal http result body:" + err.Error())
								result.Head.Code = CODE_ERROR
								result.Head.Msg = "解析返回体错误!"
							} else {
								configMap := make(map[string]string, 0)
								configList := []string{"jsapi_ticket", "timestamp", "noncestr", "url"}
								configMap["jsapi_ticket"] = serverResult.Body.Ticket
								configMap["timestamp"] = nTimeStr
								configMap["noncestr"] = newNonceStr
								configMap["url"] = PageUrl

								configSign := util.GeneratePageSign(configMap, configList)

								result.Head.Code = CODE_SUCCESS
								result.Head.Msg = "SUCCESS"
								result.Body.CodeUrl = wechatResult.CodeUrl
								result.Body.NonceStr = newNonceStr
								result.Body.Package = DEFAULT_PACKAGE_PRE_STR + wechatResult.PrepayId
								result.Body.PaySign = prepaySign
								result.Body.PrepayId = wechatResult.PrepayId
								result.Body.SignType = DEFAULT_SIGN_TYPE
								result.Body.TimeStamp = nTimeStr
								result.Body.AppId = Appid
								result.Body.ConfigSign = configSign
							}
						}
					}
				}
			}

		}
	}
	resByte, _ := json.Marshal(result)
	return string(resByte)
}
// 支付成功通知返回结构体
type AfterPayResp struct {
	Head struct {
		     Code int64  `json:"code"`
		     Msg  string `json:"msg"`
	     } `json:"head"`
	Body struct {
		     ReturnCode     string   `json:"return_code"`
		     ReturnMsg      string   `json:"return _msg"`
		     AppId          string   `json:"appid"`
		     MchId          string   `json:"mch_id"`
		     DeviceInfo     string   `json:"device_info"`
		     ResultCode     string   `json:"result_code"`
		     ErrCode        string   `json:"err_code"`
		     ErrCodeDes     string   `json:"err_code_des"`
		     OpenId         string   `json:"openid"`
		     TradeType      string   `json:"trade_type"`
		     BankType       string   `json:"bank_type"`
		     TotalFee       int64    `json:"total_fee"`
		     FeeType        string   `json:"fee_type"`
		     CashFee        int64    `json:"cash_fee"`
		     CashFeeType    string   `json:"cash_fee_type"`
		     CouponFee      int64    `json:"coupon_fee"`
		     CouponCount    int64    `json:"coupon_count"`
		     CouponIds      []string `json:"coupon_ids"`
		     CouponFees     []string `json:"coupon_fees"`
		     TransactionId  string   `json:"transaction_id"`
		     OutTradeNo     string   `json:"out_trade_no"`
		     Attach         string   `json:"attach"`
		     TimeEnd        string   `json:"time_end"`
		     ReturnToWechat string   `json:"returnToWechat"`
	     } `json:"body"`
}

// 解析完微信推送信息后应该返回给微信的结构体
type AfterPayRespToWechat struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

func DecodeWechatPayInfo(ctx *macaron.Context) string {
	result := new(AfterPayResp)
	resultToWechat := new(AfterPayRespToWechat)

	// 读取请求体
	bodyByte, err := ioutil.ReadAll(ctx.Req.Body().ReadCloser())
	body := string(bodyByte)
	if err != nil {
		log.Println("[DecodeWechatPayInfo]:error when reade request body:" + err.Error())
		result.Head.Code = CODE_ERROR
		result.Head.Msg = "读取请求体错误!"
		resultToWechat.ReturnCode = "FAIL"
		resultToWechat.ReturnMsg = "参数格式校验错误" // 这里随便写一个原因,因为是自己的问题  - - |||
	} else {
		wechatRespMap := make(map[string]string, 0)
		wechatRespList := make([]string, 0)

		// 读取微信返回内容
		inputReader := strings.NewReader(body)
		decoder := xml.NewDecoder(inputReader)
		var t xml.Token
		var err error
		name := ""
		for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
			switch token := t.(type) {

			// 处理元素开始（标签）
			case xml.StartElement:
				name = token.Name.Local
				if name != "xml" && name != "sign" {
					// 不是开始标签,将属性加入list中
					// 总标签(xml)和签名标签(sign)不参与签名计算
					wechatRespList = append(wechatRespList, name)
				}

			// 处理元素结束（标签）
			case xml.EndElement:

			// 处理字符数据（这里就是元素的文本）
			case xml.CharData:
				// 将数据存入对应的map
				content := string([]byte(token))
				if strings.TrimSpace(content) != "" {
					fmt.Println(name, ",", content)
					wechatRespMap[name] = content
				}
			default:
			}
		}

		//key := models.GetMerchantPayKey(wechatRespMap["appid"])

		// 对微信返回数据进行签名
		sign := util.GenerateSign(wechatRespMap, wechatRespList, KEY)
		if sign == wechatRespMap["sign"] {
		} else {
			log.Println("==================订单号" + wechatRespMap["out_trade_no"] + "验签失败!!!!!====================")
		}
		result.Head.Code = CODE_SUCCESS
		result.Head.Msg = "SUCCESS"
		result.Body.ReturnCode = wechatRespMap["return_code"]
		result.Body.ResultCode = wechatRespMap["return_msg"]
		result.Body.AppId = wechatRespMap["appid"]
		result.Body.MchId = wechatRespMap["mch_id"]
		result.Body.ReturnCode = wechatRespMap["result_code"]
		result.Body.OpenId = wechatRespMap["openid"]
		result.Body.TradeType = wechatRespMap["trade_type"]
		result.Body.BankType = wechatRespMap["bank_type"]
		result.Body.TransactionId = wechatRespMap["transaction_id"]
		result.Body.OutTradeNo = wechatRespMap["out_trade_no"]
		result.Body.Attach = wechatRespMap["attach"]
		result.Body.TimeEnd = wechatRespMap["time_end"]

		totalFee, _ := strconv.ParseInt(wechatRespMap["total_fee"], 10, 64)
		cashFee, _ := strconv.ParseInt(wechatRespMap["cash_fee"], 10, 64)
		result.Body.TotalFee = totalFee
		result.Body.CashFee = cashFee

		result.Body.DeviceInfo = wechatRespMap["device_info"]
		result.Body.ErrCode = wechatRespMap["err_code"]
		result.Body.ErrCodeDes = wechatRespMap["err_code_des"]

		if wechatRespMap["coupon_fee"] != "" {
			couponFee, _ := strconv.ParseInt(wechatRespMap["coupon_fee"], 10, 64)
			result.Body.CouponFee = couponFee
		}

		if wechatRespMap["coupon_count"] != "" {
			couponCount, _ := strconv.ParseInt(wechatRespMap["coupon_count"], 10, 64)
			result.Body.CouponCount = couponCount
		}

		result.Body.FeeType = wechatRespMap["fee_type"]
		if wechatRespMap["fee_type"] == "" {
			result.Body.FeeType = DEFAULT_FEE_TYPE
		}

		result.Body.CashFeeType = wechatRespMap["cash_fee_type"]
		if wechatRespMap["cash_fee_type"] == "" {
			result.Body.CashFeeType = DEFAULT_FEE_TYPE
		}

		couponIdList := make([]string, 0)
		couponFeeList := make([]string, 0)
		// 循环wechatRespMap 筛选是否存在优惠券或立减优惠信息
		for k, _ := range wechatRespMap {
			if strings.Index(k, "coupon_id_") != -1 {
				// 存在
				count := strings.Trim(k, "coupon_id_")
				couponFeeList = append(couponFeeList, wechatRespMap["coupon_fee_"+count])
				couponIdList = append(couponIdList, wechatRespMap["coupon_id_"+count])
			}
		}
		result.Body.CouponFees = couponFeeList
		result.Body.CouponIds = couponIdList

		resultToWechat.ReturnCode = "SUCCESS"
		resultToWechat.ReturnMsg = "OK"

	}

	msgToWechat, _ := xml.Marshal(resultToWechat)
	result.Body.ReturnToWechat = string(msgToWechat)
	resBytes, _ := json.Marshal(result)
	return string(resBytes)
}

// 获取JSAPI_TICKET返回值
type RespJsapiTicket struct {
	Head struct {
		     Code int64  `json:"code"`
		     Msg  string `json:"msg"`
	     } `json:"head"`
	Body struct {
		     Ticket string `json:"jsapiTicket"`
	     } `json:"body"`
}

/**
 * 获取JSAPI_TICKET
 * @Version 0.0.1
 * @Author  郭越洋
 * @Date    2015-12-18
 * @RequestType GET
 *
 * @param   appid(必填)											// 公众账号ID
 * @return
 * {
 * 	"head":{
 * 		"code":200,
 * 		"msg":"success",
 * 	},
 * 	"body":{
 *  	"jsapiTicket":"XXXXXXXXXX",
 * 	}
 * }
 */
func JsapiTicket(ctx *macaron.Context) string {
	appid := ctx.Query("appid")
	result := new(RespJsapiTicket)
	token, err := GetAccesstoken(appid)
	if err != nil {
		result.Head.Code = CODE_ERROR
		result.Head.Msg = err.Error()
	} else {
		ticket, err := GetJsapiTicket(appid, token)
		if err != nil {
			log.Println("[JsapiTicket]:error when get JsapiTicket :" + err.Error())
			result.Head.Code = CODE_ERROR
			result.Head.Msg = err.Error()
		} else {
			result.Head.Code = CODE_SUCCESS
			result.Head.Msg = "success"
			result.Body.Ticket = ticket
		}
	}
	resByte, _ := json.Marshal(result)
	return string(resByte)
}

type TokenResp struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	Token string `json:"token"`
}
type WXResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}


func GetAccesstoken(appid string) (string, error) {
	var err error

	wechantInfo := new(model.MsMerchantWechatInfo)
	wechantInfo.GetConn().Where("appid = ?", appid).First(wechantInfo)
	tokenUrl :="https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="+appid+"&secret="+wechantInfo.AppSecret
	// 改为通过apib1获取
	resp, err := http.Get(tokenUrl)
	if err!=nil{
		return "",err
	}
	res ,err1:=ioutil.ReadAll(resp)
	if err1!=nil{
		return "",err1
	}
	defer resp.Body.Close()
	wx :=new(WXResponse)
	json.Unmarshal(res,wx)
	token :=wx.AccessToken
	return token,nil
}

func GetJsapiTicket(appid, accesstoken string) (string, error) {
	result := ""
	var err error

	wechantInfo := new(model.MsMerchantWechatInfo)
	wechantInfo.GetConn().Where("appid = ?", appid).First(wechantInfo)

	//ticketType :="jsapi"


	return result, err
}

