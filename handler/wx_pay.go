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

func DecodeWechatPayInfo(ctx *macaron.Context)string{
	return ""
}
