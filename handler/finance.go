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
	"github.com/jinzhu/gorm"
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

