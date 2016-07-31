package handler

import (
	"encoding/json"


	"github.com/Unknwon/macaron"

	"wechatvoice/model"
	"strings"
	"wechatvoice/tool/util"
	"time"
	"strconv"
	"log"

)

const  (
	CODE_SUCCESS = 10000
	CODE_ERROR = 10001
	CODE_REDIRECT = 10002
	MSG_SUCCESS = "ok"
	RNF  = "record not found"
	DEFAULT_DEVICE_INFO = "WEB"
	DEFAULT_FEE_TYPE    = "CNY"
	DEFAULT_TRADE_TYPE  = "JSAPI"
	DEFAULT_NOTIFY_URL = "/wechatvoice/pay/decodewechatpayinfo"
	UNIFIEDORDER_URL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
 	DEFAULT_PACKAGE_PRE_STR = "prepay_id="
  	DEFAULT_SIGN_TYPE       = "MD5"
	TICKET_SERVER_URL  = "/shangqu-3rdparty/token/jsapi_ticket?appid="
	APPID = ""
	APPSECRET = ""
	MCHID = ""
	MCHNAME = ""
	KEY = ""
	SERVER_IP    = "127.0.0.1"

	PAY_PAGE_URL                = "orderSubmit.html"
	AFTER_PAY_ORDER_URL         = "/shangqu-shop/afterpay/wx"
	AFTER_PAY_JUMP_PAGE_FAILD   = "payFailed.html?"
	AFTER_PAY_JUMP_PAGE_SUCCESS = "paySuccess.html?"

	WECHAT_PREPAY_URL             = "/wechatvoice/pay/unifiedorder?appid=%s&mch_id=%s&body=%s&out_trade_no=%s&total_fee=%d&spbill_create_ip=%s&key=%s&openid=%s&url=%s&notify_url=%s"
)

//查询问题返回
type QuestionQueryResponse struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	List []QuestionInfo `json:"list"`
	Total int64 `json:"total"`
}

type QuestionInfo struct {
	QuestionId string `json:"questionId"`
	QuestionTopic string `json:"questionName"`
	QuestionCategoryId string  `json:"questionCateId"`
	QuestionCateName string `json:"cateName"`
	LawyerName string `json:"lawyerName"`
	LawyerId string `json:"lawyerId"`
	HeadImg string `json:"headImg"`
	VoicePath string `json:"path"`
}
//查询问题方法
/**
首页进入 开始问题搜索

params
keyWord 关键词
categoryId 对应分类ID
startLine endLine 分页请求数量
*/
func QuestionQuery(ctx *macaron.Context)string{
	body,_ :=ctx.Req.Body().String()
	req :=new(model.QuestionQuery)

	marshallErr :=json.Unmarshal([]byte(body),req)

	response :=new(QuestionQueryResponse)

	if marshallErr!=nil{
		response.Code = CODE_ERROR
		response.Msg  = marshallErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	questionList,count,queryErr :=model.GetQuestionQuery(*req)

	if queryErr!=nil&&!strings.Contains(queryErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = queryErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	retList:=make([]QuestionInfo,0)

	for _,k:=range questionList{
		single :=new(QuestionInfo)
		single.QuestionId = k.Uuid
		single.QuestionCategoryId = k.CategoryId
		single.QuestionTopic = k.Category
		single.LawyerId = k.AnswerId
		single.LawyerName = k.AnswerName
		lawyer :=new(model.LawyerInfo)
		lawyerErr :=lawyer.GetConn().Where("uuid = ?",k.AnswerId).Find(&lawyer).Error
		if lawyerErr!=nil&&!strings.Contains(lawyerErr.Error(),RNF){
			response.Code = CODE_ERROR
			response.Msg = lawyerErr.Error()
			ret_str,_:=json.Marshal(response)
			return string(ret_str)
		}
		single.HeadImg = lawyer.HeadImgUrl
		single.VoicePath = k.VoicePath
		single.QuestionCateName = k.Category
		retList = append(single,*single)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.List  = retList
	response.Total = count
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}
//提问新的问题
/**
提问新的问题

params


*/
type NewQuestionRequest struct {
	CateId string `json:"cateId"`
	CateName string `json:"cateName"`
	AskerOpenId string `json:"askerOpenId"`
	Description string `json:"description"`
	QuestionType string `json:"type"`//0 直接丢出去 1 指定人提问
	TargetOpenId string `json:"targetOpenId"`
	Payment string `json:"payment"`
}

type NewQuestionResponse struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	OrderNumber string `json:"orderNumber"`
	Payment string `json:"payment"`
}
func CreateNewQuestion(ctx *macaron.Context)string{

	body,_:=ctx.Req.Body().String()

	req:=new(NewQuestionRequest)

	response :=new(NewQuestionResponse)


	unmarshallErr:=json.Unmarshal([]byte(body),req)
	if unmarshallErr!=nil{
		response.Code = CODE_ERROR
		response.Msg = unmarshallErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	cate:=new(model.Category)
	cateErr:=cate.GetConn().Where("uuid = ?",req.CateId).Find(&cate).Error

	if cateErr!=nil&&!strings.Contains(cateErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = cateErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	customer :=new(model.MemberInfo)

	customerErr :=customer.GetConn().Where("open_id = ?",req.AskerOpenId).Find(&customer).Error

	if customerErr!=nil&&!strings.Contains(customerErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = customerErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	orderNumber :=util.GenerateOrderNumber()
	question :=new(model.WechatVoiceQuestions)
	question.Uuid = util.GenerateUuid()
	question.CategoryId  = req.CateId
	question.Category = req.CateName
	question.CategoryIdInt = int64(cate.Model.ID)
	question.Description = req.Description
	today := time.Unix(time.Now().Unix(), 0).String()[0:19]
	question.CreateTime = today
	question.CustomerId = customer.Uuid
	question.CustomerName = customer.Name
	question.CustomerOpenId = req.AskerOpenId
	question.PaymentInfo = req.Payment
	payInt,transferErr :=strconv.ParseInt(req.Payment,10,64)
	question.OrderNumber = orderNumber
	if transferErr!=nil&&!strings.Contains(transferErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = transferErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	question.PaymentInfoInt = payInt
	switch req.QuestionType {
	case "0":
		question.QuestionType ="0"
		question.Important = "0"
	case "1":
		question.AnswerOpenId = req.TargetOpenId

		layer:=new(model.LawyerInfo)
		lawyerErr :=layer.GetConn().Where("open_id = ?",req.TargetOpenId).Find(&layer).Error

		if lawyerErr!=nil&&!strings.Contains(lawyerErr.Error(),RNF){
			response.Code = CODE_ERROR
			response.Msg = lawyerErr.Error()
			ret_str,_:=json.Marshal(response)
			return string(ret_str)
		}
		question.AnswerName = layer.Name
		question.AnswerId = layer.Uuid
		question.AnswerHeadImg = layer.HeadImgUrl

	}

	createErr :=question.GetConn().Create(&question).Error

	if createErr!=nil{
		response.Code = CODE_ERROR
		response.Msg = createErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg  = MSG_SUCCESS
	response.OrderNumber = orderNumber
	response.Payment = req.Payment
	ret_str,_:=json.Marshal(response)
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



//这里获取分类问题的配置选项
type GetConfigRequest struct{
	CateGoryId string `json:"cateId"`
}

type QuestionConfig struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	PayAmountStr string  `json:"amountStr"`
	PayAmountInt int64 `json:"amountInt"`
}
func GetQuestionConfig(ctx *macaron.Context)string{
	body,_:=ctx.Req.Body().String()
	req :=new(GetConfigRequest)
	unmarshallErr :=json.Unmarshal([]byte(body),req)

	response :=new(QuestionConfig)

	if unmarshallErr!=nil{
		response.Code = CODE_ERROR
		response.Msg = unmarshallErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	setting:=new(model.WechatVoiceQuestionSettings)
	settingErr :=setting.GetConn().Where("category_id = ?",req.CateGoryId).Find(&setting).Error
	if settingErr!=nil&&!strings.Contains(settingErr.Error(),"record not found"){
		response.Code = CODE_ERROR
		response.Msg = settingErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg  = "ok"
	response.PayAmountInt = setting.PayAmountInt
	response.PayAmountStr = setting.PayAmount
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}
//获取分类列表1
type QuestionCateList struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	List []CateInfo `json:"list"`
}

type CateInfo struct {
	CateName string `json:"cateName"`
	CateId string `json:"cateId"`
}

func GetQuestionCateList(ctx *macaron.Context)string{
	response :=new(QuestionCateList)

	list :=make([]CateInfo,0)

	cateList,cateErr :=model.GetCateList()

	if cateErr!=nil&&!strings.Contains(cateErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = cateErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	for _,k:=range cateList{
		single :=new(CateInfo)
		single.CateId = k.Uuid
		single.CateName = k.CategoryName
		list = append(list,*single)
	}

	response.Code = CODE_SUCCESS
	response.Msg  = "ok"
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}

//将问题进行追问
type QuestionAppendRequest struct {
	QuestionId string `json:"parentQId"`
	CateId string `json:"cateId"`
	CateName string `json:"cateName"`
	AskerOpenId string `json:"askerOpenId"`
	Description string `json:"description"`
	QuestionType string `json:"type"`//0 直接丢出去 1 指定人提问
	TargetOpenId string `json:"targetOpenId"`
	Payment string `json:"payment"`
}


type QuestionNewResponse struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	OrderNumber string `json:"orderNumber"`
}
func AppendQuestion(ctx *macaron.Context)string{
	body,_:=ctx.Req.Body().String()
	req :=new(QuestionAppendRequest)
	response :=new(QuestionNewResponse)
	json.Unmarshal([]byte(body),req)


	questionInfoOld :=new(model.WechatVoiceQuestions)
	questionOldErr :=questionInfoOld.GetConn().Where("uuid = ?",req.QuestionId).Find(&questionInfoOld).Error

	if questionOldErr.Error()!=nil&&!strings.Contains(questionOldErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = questionOldErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	if questionInfoOld.IsSolved!="2"{
		//问题没有解决 不能进行追问
		response.Code = CODE_ERROR
		response.Msg = "问题没有解决 不能进行追问"
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	customer :=new(model.MemberInfo)

	customerErr :=customer.GetConn().Where("open_id = ?",req.AskerOpenId).Find(&customer).Error

	if customerErr!=nil&&!strings.Contains(customerErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = customerErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}



	questionNew :=new(model.WechatVoiceQuestions)

	questionNew.Uuid = util.GenerateUuid()
	questionNew.Category = req.CateName
	questionNew.CategoryId  = req.CateId
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
	on :=util.GenerateOrderNumber()
	questionNew.OrderNumber = on
	questNewErr :=questionNew.GetConn().Create(&questionNew).Error

	if questNewErr!=nil{
		response.Code = CODE_ERROR
		response.Msg = questNewErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.OrderNumber =on
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}

//偷听业务
type PeekAnswerRequest struct {
	OrderId string `json:"orderId"`
	CateId string `json:"cateId"`
}

/**
偷听应该是这样用户点击 先看是否有权限去听 如果有 直接放 如果不能 调起微信支付 支付结束后 可以看
*/
type PeekResponse struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	PlayAble bool `json:"playAble"`
}
//这里需要看下微信支付业务 点击获取之后

func PeekAvalable(ctx *macaron.Context)string{
	body,_:=ctx.Req.Body().String()
	req :=new(PeekAnswerRequest)
	json.Unmarshal([]byte(body),req)
	response :=new(PeekResponse)


	//设置cookie  第一段为openId 第二段为类型 1 用户 2律师
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr==""{
		//这里直接调取util重新过一次绿叶 获取openId 等信息
	}
	openId :=strings.Split(cookieStr,"|")[0]
	userType :=strings.Split(cookieStr,"|")[1]

	log.Println(openId)
	log.Println(userType)


	pay :=new(model.OrderPaymentInfo)
	payErr:=pay.GetConn().Where("question_id = ?",req.OrderId).Where("open_id = ?",openId).Find(&pay).Error

	if payErr!=nil&&!strings.Contains(payErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = payErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	if pay.Uuid ==""{
		response.Code = CODE_SUCCESS
		response.Msg = MSG_SUCCESS
		response.PlayAble = false
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.PlayAble = true
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}

func DoPay(ctx *macaron.Context)string{
	//这里去看下微信支付功能
	return ""
}

//开始回答部分

//一个文件上传服务

func VoiceUpLoad(ctx *macaron.Context)string{
	//这里返回一个路径
	return ""
}

type AnswerQuestion1 struct {
	QuestionId string `json:"questionId"`
}

type AnswerQuestion1Response struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	QuestionInfo `json:"question"`
}
type QuestionInfo struct {
	QuestionId string `json:"quesiontId"`
	QuestionCateInfo string `json:"cateInfo"`
	QuestionCateId string `json:"cateId"`
	QuestionDesc string `json:"desc"`


	AskerName string `json:"askerName"`
	AskerId string `json:"askerId"`
	AskerHeadImg string `json:"askerHeadImg"`
}
func AnswerQuestionInit(ctx *macaron.Context)string{
	//点击回答问题  显示问题
	body,_:=ctx.Req.Body().String()
	req :=new(AnswerQuestion1)
	json.Unmarshal([]byte(body),req)
	response :=new(AnswerQuestion1Response)

	//设置cookie  第一段为openId 第二段为类型 1 用户 2律师
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr==""{
		//这里直接调取util重新过一次绿叶 获取openId 等信息
	}
	openId :=strings.Split(cookieStr,"|")[0]
	userType :=strings.Split(cookieStr,"|")[1]

	log.Println(openId)
	log.Println(userType)

	question:=new(model.WechatVoiceQuestions)
	qErr :=question.GetConn().Where("uuid = ?",req.QuestionId).Find(&question).Error

	if qErr!=nil&&!strings.Contains(qErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = qErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	if question.IsSolved == "1"{
		response.Code = CODE_ERROR
		response.Msg = "已经回答完毕 请勿重复做大"
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	if question.Important =="1"&&question.AnswerOpenId!=openId{
		response.Code = CODE_ERROR
		response.Msg = "这个为指定问题"
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	if userType=="1"{
		response.Code = CODE_ERROR
		response.Msg = "无权限"
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	q :=new(QuestionInfo)

	q.QuestionId = question.Uuid
	q.HeadImg = question.AskerHeadImg
	q.QuestionCategoryId = question.CategoryId
	q.QuestionCateName = question.Category
	q.VoicePath = ""
	q.QuestionTopic = question.Description
	lawerInfo :=new(model.LawyerInfo)

	lawerInfoErr :=lawerInfo.GetConn().Where("open_id = ?",openId).Find(&lawerInfo).Error

	if lawerInfoErr!=nil&&!strings.Contains(lawerInfoErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = lawerInfoErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	q.LawyerId = lawerInfo.Uuid
	q.LawyerName = lawerInfo.Name

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.QuestionInfo = *q
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}
type DoAnsweQuestion struct {
	FilPath string `json:"filePath"`
	QuestionId string `json:"questionId"`

}

func DoAnswerQuestion(ctx *macaron.Context)string{
	body,_:=ctx.Req.Body().String()
	req :=new(DoAnsweQuestion)
	json.Unmarshal([]byte(body),req)
	response :=new(model.GeneralResponse)

	question:=new(model.WechatVoiceQuestions)

	//设置cookie  第一段为openId 第二段为类型 1 用户 2律师
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr==""{
		//这里直接调取util重新过一次绿叶 获取openId 等信息
	}
	openId :=strings.Split(cookieStr,"|")[0]
	userType :=strings.Split(cookieStr,"|")[1]

	log.Println(openId)
	log.Println(userType)

	questionErr:=question.GetConn().Where("uuid = ?",req.QuestionId).Find(&question).Error
	if questionErr!=nil&&!strings.Contains(questionErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = questionErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	if question.IsSolved == "1"{
		response.Code = CODE_ERROR
		response.Msg = "已经回答完毕 请勿重复做大"
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	if question.Important =="1"&&question.AnswerOpenId!=openId{
		response.Code = CODE_ERROR
		response.Msg = "这个为指定问题"
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	if userType=="1"{
		response.Code = CODE_ERROR
		response.Msg = "无权限"
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	lawerInfo :=new(model.LawyerInfo)

	lawerInfoErr :=lawerInfo.GetConn().Where("open_id = ?",openId).Find(&lawerInfo).Error

	if lawerInfoErr!=nil&&!strings.Contains(lawerInfoErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = lawerInfoErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	question.IsAnswerd = "1"
	question.VoicePath = req.FilPath
	question.AnswerId = lawerInfo.Uuid
	question.AnswerOpenId = openId
	question.AnswerName = lawerInfo.Name
	question.AnswerHeadImg = lawerInfo.HeadImgUrl

	err :=question.GetConn().Save(&question).Error

	if err!=nil&&!strings.Contains(err.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}

func AddPeekCount(questionId string)error{
	quesionInfo:=new(model.WechatVoiceQuestions)
	quesionInfoErr:=quesionInfo.GetConn().Where("uuid = ?",questionId).Find(&quesionInfo).Error
	if quesionInfoErr!=nil{
		return quesionInfoErr
	}
	quesionInfo.AnswerdCount = quesionInfo.AnswerdCount+1
	err :=quesionInfo.GetConn().Save(&quesionInfo).Error
	return err
}



