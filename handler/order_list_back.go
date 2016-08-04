package handler

import (
	"encoding/json"

	"github.com/Unknwon/macaron"
	"strings"
	"wechatvoice/model"
	//"wechatvoice/tool/util"
	//"time"
	//"strconv"
	//"strconv"
)

type QuestionInfoBackResponse struct {
	Code  int64              `json:"code"`
	Msg   string             `json:"msg"`
	List  []QuestionInfoBack `json:"list"`
	Total int64              `json:"total"`
}
type QuestionInfoBack struct {
	QuestionId         string `json:"questionId"`
	QuestionTopic      string `json:"questionName"`
	QuestionCategoryId string `json:"questionCateId"`
	QuestionCateName   string `json:"cateName"`
	LawyerName         string `json:"lawyerName"`
	LawyerId           string `json:"lawyerId"`
	LawyerOpenId       string `json:"lOpenId"`
	LawyerHeadImg      string `json:"lHead"`
	HeadImg            string `json:"headImg"`
	VoicePath          string `json:"path"`

	AskerName   string `json:"askerName"`
	AskerOpenId string `json:"askerOpenId"`
}

type QuestionInfoReq struct {
	StartLine int64 `json:"startLine"`
	EndLine   int64 `json:"endLine"`
}

func GetBadAnswerList(ctx *macaron.Context) string {
	body, _ := ctx.Req.Body().String()

	req := new(QuestionInfoReq)
	response := new(QuestionInfoBackResponse)

	unmarshallErr := json.Unmarshal([]byte(body), req)

	if unmarshallErr != nil {
		response.Code = CODE_ERROR
		response.Msg = unmarshallErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	stastusList := make([]string, 0)
	stastusList = append(stastusList, "1")
	stastusList = append(stastusList, "2")

	list, count, err := model.QueryBadAnswers(stastusList, req.StartLine, req.EndLine)

	if err != nil && !strings.Contains(err.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	retList := make([]QuestionInfoBack, 0)
	for _, k := range list {
		single := new(QuestionInfoBack)
		single.QuestionId = k.Uuid
		single.QuestionTopic = k.Description
		single.QuestionCategoryId = k.CategoryId
		single.QuestionCateName = k.Category
		single.LawyerName = k.AnswerName
		single.LawyerOpenId = k.AnswerOpenId
		single.LawyerId = k.AnswerId
		single.LawyerHeadImg = k.AnswerHeadImg
		single.VoicePath = k.VoicePath
		single.AskerName = k.CustomerName
		single.AskerOpenId = k.CustomerOpenId
		retList = append(retList, *single)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.List = retList
	response.Total = count
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

type EvaluateRequest struct {
	Stars      string `json:"starts"`
	PassStatus string `json:"status"`
	QuestionId string `json:"questionId"`
}

func ReEvaluatBadAnswers(ctx *macaron.Context) string {
	body, _ := ctx.Req.Body().String()

	req := new(EvaluateRequest)
	response := new(GeneralResponse)

	json.Unmarshal([]byte(body), req)

	question := new(model.WechatVoiceQuestions)

	err := question.GetConn().Where("uuid = ?", req.QuestionId).Find(&question).Error
	if err != nil && !strings.Contains(err.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	question.RankInfo = req.Stars
	question.IsRanked = req.PassStatus

	err = question.GetConn().Save(&question).Error

	if err != nil {
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
