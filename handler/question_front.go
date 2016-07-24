package handler

import (
	"encoding/json"

	"github.com/Unknwon/macaron"
	"wechatvoice/model"
	"strings"
	"wechatvoice/tool/util"
	"time"
	"strconv"
)

const  (
	CODE_SUCCESS = 10000
	CODE_ERROR = 10001
	CODE_REDIRECT = 10002
	MSG_SUCCESS = "ok"
	RNF  = "record not found"
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

type NewQuestionRequest struct {
	CateId string `json:"cateId"`
	CateName string `json:"cateName"`
	AskerOpenId string `json:"askerOpenId"`
	Description string `json:"description"`
	QuestionType string `json:"type"`//0 直接丢出去 1 指定人提问
	TargetOpenId string `json:"targetOpenId"`
	Payment string `json:"payment"`
}
func CreateNewQuestion(ctx *macaron.Context)string{

	body,_:=ctx.Req.Body().String()

	req:=new(NewQuestionRequest)

	response :=new(model.GeneralResponse)


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
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}
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
