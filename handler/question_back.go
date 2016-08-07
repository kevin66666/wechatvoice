package handler

import (
	"encoding/json"

	"github.com/Unknwon/macaron"
	"strings"
	"wechatvoice/model"
	"wechatvoice/tool/util"
	//"time"
	//"strconv"
	"strconv"
)

type GeneralResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}
type AddCateRequest struct {
	CateName string `json:"name"`
}

func CreateCateList(ctx *macaron.Context) string {
	body, _ := ctx.Req.Body().String()

	req := new(AddCateRequest)
	response := new(GeneralResponse)
	json.Unmarshal([]byte(body), req)

	uuid := util.GenerateUuid()

	cate := new(model.LawCatgory)

	cate.Uuid = uuid
	cate.CategoryName = req.CateName

	err := cate.GetConn().Create(&cate).Error

	if err != nil {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = "ok"
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

//获取分类列表

type CateLisReq struct {
	StartLine int64 `json:"startLine"`
	EndLine int64 `json:"endLine"`
}

type CateListResponse struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	List CateDetail `json:"list"`
	Count int64 `json:"count"`
}

type CateDetail struct {
	CateId string `json:"cateId"`
	CateName string `json:"cateName"`
}

func GetCateList(ctx *macaron.Context)string{
	body,_:=ctx.Req.Body().String()
	req :=new(CateLisReq)

	response :=new(CateListResponse)

	err :=json.Unmarshal([]byte(body),req)
	if err!=nil{
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	list :=make([]CateDetail,0)

	lists,count,errs :=model.GetCateList(req.StartLine,req.EndLine)
	if errs!=nil&&!strings.Contains(errs.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = errs.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	for _,k:=range lists{
		single :=new(CateDetail)
		single.CateId = k.Uuid
		single.CateName = k.CategoryName
		list=  append(list,&single)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.Count = count
	response.List = list
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}
//删除分类

type DeleteReq struct {
	CateId string `json:"cateId"`
}
func DeleteCateInfo(ctx *macaron.Context)string{
	body,_:=ctx.Req.Body().String()
	req :=new(DeleteReq)

	json.Unmarshal([]byte(body),req)

	response :=new(model.GeneralResponse)

	cate:=new(model.Category)
	cateErr :=cate.GetConn().Where("uuid = ?",req.CateId).Delete(&cate).Error
	if cateErr!=nil&&!strings.Contains(cateErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = cateErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}

//编辑分类信息

type EditCate struct {
	CateName string `json:"cateName"`
	CateId string `json:"cateId"`
}

func EditCateInfo(ctx *macaron.Context)string{
	body,_:=ctx.Req.Body().String()
	req :=new(EditCate)

	json.Unmarshal([]byte(body),req)

	response :=new(model.GeneralResponse)

	cate:=new(model.Category)
	cateErr :=cate.GetConn().Where("uuid = ?",req.CateId).Find(&cate).Error
	if cateErr!=nil&&!strings.Contains(cateErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = cateErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	cate.CategoryName = req.CateName
	cateErr = cate.GetConn().Save(&cate).Error
	if cateErr!=nil&&!strings.Contains(cateErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = cateErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}

type SettingListReq struct {
	StartLine int64 `json:"startLine"`
	EndLine   int64 `json:"endLine"`
}

type SettingListResponse struct {
	Code int64     `json:"code"`
	Msg  string    `json:"msg"`
	List []Setting `json:"list"`
}

type Setting struct {
	SettingId        string `json:"settingId"`
	CateId           string `json:"cateId"`
	CateName         string `json:"cateName"`
	AmountInt        string `json:"amount"`
	LawyerPercent    string `json:"lawerP"`
	RedPacketPercent string `json:"redPacket"`
}

func GetQuestionSettingList(ctx *macaron.Context) string {

	req := new(SettingListReq)

	response := new(SettingListResponse)

	list := make([]Setting, 0)

	retList, err := model.GetSettingList(req.StartLine, req.EndLine)

	if err != nil && !strings.Contains(err.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	for _, k := range retList {
		single := new(Setting)
		single.SettingId = k.Uuid
		single.CateName = k.CateGoryName
		single.CateId = k.CategoryId
		single.AmountInt = k.PayAmount
		single.LawyerPercent = k.LawyerFeePercent
		single.RedPacketPercent = k.UserRedPacketPercent
		list = append(list, *single)
	}

	response.Code = CODE_SUCCESS
	response.Msg = "ok"
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

type IdRequest struct {
	SettingId string `json:"settingId"`
}

type SingleResponse struct {
	Code    int64  `json:"code"`
	Msg     string `json:"msg"`
	Setting `json:"setting"`
}

func GetQuestionSettingsById(ctx *macaron.Context) string {
	body, _ := ctx.Req.Body().String()

	req := new(IdRequest)

	json.Unmarshal([]byte(body), req)

	response := new(SingleResponse)

	setting := new(model.WechatVoiceQuestionSettings)

	settingErr := setting.GetConn().Where("uuid = ?", req.SettingId).Find(&setting).Error

	if settingErr != nil && !strings.Contains(settingErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = settingErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	single := new(Setting)

	single.SettingId = setting.Uuid
	single.CateName = setting.CateGoryName
	single.CateId = setting.CategoryId
	single.AmountInt = setting.PayAmount
	single.LawyerPercent = setting.LawyerFeePercent
	single.RedPacketPercent = setting.UserRedPacketPercent

	response.Code = CODE_SUCCESS
	response.Msg = "ok"
	response.Setting = *single
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

func EditWechatVoiceQuestionSettings(ctx *macaron.Context) string {
	req := new(Setting)

	response := new(GeneralResponse)

	body, _ := ctx.Req.Body().String()
	json.Unmarshal([]byte(body), req)

	setting := new(model.WechatVoiceQuestionSettings)

	err := setting.GetConn().Where("uuid = ?", req.SettingId).Find(&setting).Error

	if err != nil && !strings.Contains(err.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	setting.LawyerFeePercent = req.LawyerPercent
	setting.CategoryId = req.CateId
	setting.CateGoryName = req.CateName
	setting.PayAmount = req.AmountInt
	amtInt, _ := strconv.ParseInt(req.AmountInt, 10, 64)
	setting.PayAmountInt = amtInt
	setting.UserRedPacketPercent = req.RedPacketPercent

	err = setting.GetConn().Save(&setting).Error

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

//删除对应信息
func DeleteQuestionSettingsById(ctx *macaron.Context)string{
	body, _ := ctx.Req.Body().String()

	req := new(IdRequest)

	json.Unmarshal([]byte(body), req)

	response :=new(model.GeneralResponse)

	setting:=new(model.WechatVoiceQuestionSettings)
	settingErr :=setting.GetConn().Where("uuid = ?",req.SettingId).Delete(&setting).Error

	if settingErr!=nil&&!strings.Contains(settingErr.Error(),RNF){
		response.Code =CODE_ERROR
		response.Msg = settingErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}
func AddQuestionSetting(ctx *macaron.Context)string{
	req := new(Setting)

	response := new(GeneralResponse)

	body, _ := ctx.Req.Body().String()
	json.Unmarshal([]byte(body), req)

	setting := new(model.WechatVoiceQuestionSettings)
	setting.Uuid = util.GenerateUuid()
	setting.CategoryId=req.CateId
	setting.CateGoryName = req.CateName
	setting.LawyerFeePercent = req.LawyerPercent
	setting.CategoryId = req.CateId
	setting.CateGoryName = req.CateName
	setting.PayAmount = req.AmountInt
	amtInt, _ := strconv.ParseInt(req.AmountInt, 10, 64)
	setting.PayAmountInt = amtInt
	setting.UserRedPacketPercent = req.RedPacketPercent

	err :=setting.GetConn().Create(&setting).Error
	if err!=nil{
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}
