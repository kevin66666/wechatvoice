package handler

import (
	"encoding/json"

	"github.com/Unknwon/macaron"
	"wechatvoice/model"
	"strings"
	"wechatvoice/tool/util"
	//"time"
	//"strconv"
	//"strconv"
	"strconv"


	"math/rand"
)
var RED_PACKET_SETTNG int64
var LAWYER_PERCENT_SETTING int64
var PAY_AMT_SETTING int64



func GetUserRedPacketSettings(cateId string)(string,string,string,error){
	setting :=new(model.WechatVoiceQuestionSettings)

	settingErr :=setting.GetConn().Where("category_id = ?",cateId).Find(&setting).Error

	if settingErr!=nil&&!strings.Contains(settingErr.Error(),RNF){
		return "","","",settingErr
	}

	return setting.PayAmount,setting.LawyerFeePercent,setting.UserRedPacketPercent,nil
}


type NewPaymentLogReq struct {
	OrderId string `json:"orderId"`
	SwiftNumber string `json:"sn"`
	QuestionId string `json:"questionId"`
	OpenId string `json:"openId"`
	MemberId string `json:"memberId"`
}
func CreatePaymentLog(ctx *macaron.Context)string{
	body,_:=ctx.Req.Body().String()

	req :=new(NewPaymentLogReq)
	response :=new(GeneralResponse)

	json.Unmarshal([]byte(body),req)

	pay :=new(model.WechatVoicePaymentInfo)

	pay.Uuid = util.GenerateUuid()

	pay.SwiftNumber = req.SwiftNumber
	pay.QuestionId = req.QuestionId
	pay.MemberId =req.MemberId
	pay.OpenId =req.OpenId

	err :=pay.GetConn().Create(&pay).Error
	if err!=nil{
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

type GetRedPacketRequest struct {
	OrderNumber string `json:"orderNumber"`
	QuestionId string `json:"questionId"`
}

type FinanceResponse struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	RedAmount string `json:"red"`
}

func GetFinanceInfoByOrderNumber(ctx *macaron.Context)string{
	body,_:=ctx.Req.Body().String()

	req :=new(GetRedPacketRequest)
	response :=new(FinanceResponse)
	json.Unmarshal([]byte(body),req)

	questionInfo :=new(model.WechatVoiceQuestions)
	questionSetting :=new(model.WechatVoiceQuestionSettings)

	questinfoErr:=questionInfo.GetConn().Where("uuid = ?",req.QuestionId).Find(&questionInfo).Error
	if questinfoErr!=nil&&!strings.Contains(questinfoErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = questinfoErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	cateId :=questionInfo.CategoryId

	questSetErr:=questionSetting.GetConn().Where("category_id = ?",cateId).Find(&questionSetting).Error
	if questSetErr!=nil&&!strings.Contains(questSetErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = questSetErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	payment :=new(model.OrderPaymentInfo)

	paymentErr :=payment.GetConn().Where("question_id = ?",req.QuestionId).Find(&payment).Error
	if paymentErr!=nil&&!strings.Contains(paymentErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = paymentErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	layFee :=questionSetting.LawyerFeePercent
	red :=questionSetting.UserRedPacketPercent
	total :=payment.UserPaiedAmountInt


	layFeeInt ,_:=strconv.ParseInt(layFee,10,64)
	layFeeAmt := (total*layFeeInt)/100
	layFeeStr :=strconv.FormatInt(layFeeAmt,10)

	redInt,_:=strconv.ParseInt(red,10,64)
	redAmountLeft :=(10-layFee)*redInt*total/100
	//100 *2 * 8 /100
	randomNumber :=rand.Int63n(10)

	redAmt :=redAmountLeft*randomNumber/10
	redAmtStr :=strconv.FormatInt(redAmt,10)

	amtSlice :=strings.Split(redAmtStr,".")
	decimal :=amtSlice[1]
	mainInfo :=amtSlice[0]

	var redInfo int64

	decimalInt,_:=strconv.ParseInt(decimal,10,64)
	mainInt ,_:=strconv.ParseInt(mainInfo,10,64)

	if decimalInt>=5{
		redInfo = mainInt+1
	}else{
		redInfo = mainInt
	}
	redStr :=strconv.FormatInt(redInfo,10)

	payment.LawyerFee = layFeeStr
	payment.LawyerFeeInt = layFeeInt

	payment.RedPacketAmountInt = redInfo
	payment.RedPacketAmount = redStr
	balance :=total-layFeeInt-redInfo
	payment.BalanceAmountInt =balance
	balanceStr :=strconv.FormatInt(balance,10)
	payment.BalanceAmount = balanceStr

	updatePayErr :=payment.GetConn().Save(&payment).Error
	if updatePayErr!=nil&&!strings.Contains(updatePayErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = updatePayErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	questionInfo.IsRanked = "1"
	questionInfo.IsRanked = "1"
	questionUpdateErr :=questionInfo.GetConn().Save(&questionInfo).Error
	if questionUpdateErr!=nil&&!strings.Contains(questionUpdateErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = questionUpdateErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	memberInfo :=new(model.MemberInfo)
	memberQueryErr :=memberInfo.GetConn().Where("open_id = ?",questionInfo.CustomerOpenId).Find(&memberInfo).Error

	if memberQueryErr!=nil&&!strings.Contains(memberQueryErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = memberQueryErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	userBalance :=memberInfo.Balance
	balanceInt,_:=strconv.ParseInt(userBalance,10,64)

	balanceInt = balanceInt + redInfo

	balanceNewStr :=strconv.FormatInt(balanceInt,10)

	memberInfo.Balance = balanceNewStr
	memberUpdateErr :=memberInfo.GetConn().Save(&memberInfo).Error

	if memberUpdateErr!=nil&&!strings.Contains(memberUpdateErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = memberUpdateErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.RedAmount = redStr
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}
type BackFinanceListRequest struct {
	StartLine int64 `json:"startLine"`
	EndLine int64 `json:"endLine"`
	Name string `json:"name"`
	//QuestionId string `json:"questionId"`
}

type BackFinanceListResponse struct{
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	Total int64 `json:"total"`
	List []LawyerInfo `json:"list"`
}

type LawyerInfo struct {
	Uuid string `json:"id"`
	Name string `json:"name"`
	OrderCount int64 `json:"orderCount"`
	PhoneNumber int64 `json:"phoneNumber"`
	NickName string `json:"nickName"`
	HeadImgUrl string `json:"headImgUrl"`
	OpenId string `json:"openId"`
	Balance string `json:"balance"`
}
func GetFinanceBackList(ctx *macaron.Context)string{
	body,_:=ctx.Req.Body().String()

	req :=new(BackFinanceListRequest)
	response :=new(BackFinanceListResponse)
	unmarShallErr  :=json.Unmarshal([]byte(body),req)

	if unmarShallErr!=nil{
		response.Code = CODE_ERROR
		response.Msg = unmarShallErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	list,count,err :=model.GetLayersFinanceQueryInfo(req.StartLine,req.EndLine,req.Name)

	if err!=nil&&!strings.Contains(err.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	retList :=make([]LawyerInfo,0)
	for _,k:=range list{
		single :=new(LawyerInfo)
		single.Uuid = k.Uuid
		single.Name = k.Name
		single.OrderCount = k.OrderCount
		single.PhoneNumber = k.PhoneNumber
		single.NickName = k.NickName
		single.HeadImgUrl  = k.HeadImgUrl
		balance:=k.Balance
		single.OpenId = k.OpenId
		balanceStr :=strconv.FormatFloat(balance,'f',2,64)
		single.Balance = balanceStr
		retList = append(retList,*single)
	}
	response.Code = CODE_SUCCESS
	response.Msg  ="ok"
	response.Total = count
	response.List = retList
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}
type LawyerInfoReq struct {
	StartLine int64 `json:"startLine"`
	EndLine int64 `json:"endLine"`
	Uuid string  `json:"openId"`
}
type LawyerOrderListInfo struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	Count int64 `json:"count"`
	List []LawerOrder `json:"list"`
}

type LawerOrder struct {
	Category string `json:"category"`
	LawerName string `json:"name"`
	QuetstionName string `json:"qName"`
	CreateTime string `json:"createTime"`
	RankInfo int64 `json:"rank"`
	Money string `json:"money"`
}

func GetLawyerOrderlistInfo(ctx *macaron.Context)string{
	body ,_:=ctx.Req.Body().String()
	req:=new(LawyerInfoReq)
	response :=new(LawyerOrderListInfo)

	unmarshallErr :=json.Unmarshal([]byte(body),req)
	if unmarshallErr!=nil{
		response.Code = CODE_ERROR
		response.Msg = unmarshallErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	list,count,err :=model.QueryLawyerQuestions(req.StartLine,req.EndLine,req.Uuid)
	retList :=make([]LawerOrder,0)
	if err!=nil&&!strings.Contains(err.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	for _,k:=range list{
		single :=new(LawerOrder)
		single.Category = k.Category
		single.LawerName = k.AnswerName
		single.QuetstionName = k.Description
		single.CreateTime = k.CreateTime
		single.RankInfo = k.RankInfo
		pay :=new(model.OrderPaymentInfo)

		payErr :=pay.GetConn().Where("question_id = ?",k.Uuid).Find(&pay).Error
		if payErr!=nil&&!strings.Contains(payErr.Error(),RNF){
			response.Code = CODE_ERROR
			response.Msg = payErr.Error()
			ret_str,_:=json.Marshal(response)
			return string(ret_str)
		}

		single.Money = pay.LawyerFee
		retList = append(retList,*single)
	}

	response.Code = CODE_SUCCESS
	response.Msg =MSG_SUCCESS
	response.Count = count
	response.List = retList
	ret_str,_:=json.Marshal(response)
	return string(retList)

}