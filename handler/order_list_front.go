package handler

import (
	"encoding/json"
	"github.com/Unknwon/macaron"
	"log"
	"strconv"
	"strings"
	"wechatvoice/model"
)

type OrderListFrontRequest struct {
	StartLine   int64  `json:"startLine"`
	EndLine     int64  `json:"endLine"`
	OrderStatus string `json:"status"` // 0 未解决 1 已解决
}

type OrderListResponse struct {
	Code  int64       `json:"code"`
	Msg   string      `json:"msg"`
	Total int64       `json:"total"`
	List  []OrderInfo `json:"list"`
}

type OrderInfo struct {
	OrderId      string `json:"orderId"`
	Destription  string `json:"des"`
	AskerName    string `json:"askerName"`
	AskerOpenId  string `json:"askerOpenId"`
	AskerHeadImg string `json:"askerHeadImg"`
	IsSolved     string `json:"isSolved"`
	AnsweredTime string `json:"answeredTime"`
	AskTime      string `json:"askTime"`

	LawyerId      string `json:"lawyerId"`
	LawyerName    string `json:"lawyerName"`
	LawyerHeadImg string `json:"lawyerHeadImg"`
	LawyerOpenId  string `json:"lawyerOpenId"`

	VoicePath string `json:"voicePath"`

	RankInfo     int64  `json:"rank"`
	Pv           int64  `json:"pv"`
	QuestionType string `json:"questionType"`
	HasChild     bool   `json:"haveChild"`
}

func GetOrderList(ctx *macaron.Context) string {
	//设置cookie  第一段为openId 第二段为类型 1 用户 2律师
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" {
		//这里直接调取util重新过一次绿叶 获取openId 等信息
	}
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]

	log.Println("=========>>>>>>,用户OPENID 为", openId)
	log.Println("=========>>>>>>,用户类型为", userType)

	body, _ := ctx.Req.Body().String()

	req := new(OrderListFrontRequest)

	response := new(OrderListResponse)

	marshallErr := json.Unmarshal([]byte(body), req)

	if marshallErr != nil {
		response.Code = CODE_ERROR
		response.Msg = marshallErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	a := []string{"0"}
	b := []string{"1", "2"}

	list := make([]model.WechatVoiceQuestions, 0)
	var count int64
	var err error

	retList := make([]OrderInfo, 0)
	switch userType {
	case "1":
		if req.OrderStatus == "0" {
			list, count, err = model.QueryUserQuestions(openId, a, req.StartLine, req.EndLine)
		} else {
			list, count, err = model.QueryUserQuestions(openId, b, req.StartLine, req.EndLine)
		}
		if err != nil && !strings.Contains(err.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = err.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		for _, k := range list {
			//里面分别是订单
			single := new(OrderInfo)
			single.OrderId = k.Uuid
			single.Destription = k.Description
			single.AskerName = k.CustomerName
			single.AskerOpenId = k.CustomerOpenId
			single.AskerHeadImg = k.AskerHeadImg
			if k.IsSolved == "2" {
				single.LawyerId = k.AnswerId
				single.LawyerName = k.AnswerName
				single.LawyerHeadImg = k.AnswerHeadImg
				single.LawyerOpenId = k.AnswerOpenId
				single.VoicePath = k.VoicePath
				rank := k.RankInfo
				rankInt, _ := strconv.ParseInt(rank, 10, 64)
				single.RankInfo = rankInt
				single.Pv = k.Pv
				single.QuestionType = k.Category
			}
			retList = append(retList, *single)
		}
	case "2":
		if req.OrderStatus == "1" {
			list, count, err = model.QueryLawyerQuestions(req.StartLine, req.EndLine, openId)
			if err != nil && !strings.Contains(err.Error(), RNF) {
				response.Code = CODE_ERROR
				response.Msg = err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}

		} else {
			userInfo := new(model.LawyerInfo)
			layerErr := userInfo.GetConn().Where("open_id = ?", openId).Find(&userInfo).Error
			if layerErr != nil && !strings.Contains(layerErr.Error(), RNF) {
				response.Code = CODE_ERROR
				response.Msg = layerErr.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
			catList := make([]string, 0)
			catList = append(catList, userInfo.FirstCategory)
			catList = append(catList, userInfo.SecondCategory)
			catList = append(catList, userInfo.ThridCategory)

			list, count, err = model.QueryLawerNotSolvedQuestions(catList, req.StartLine, req.EndLine)
		}
		for _, k := range list {
			//里面分别是订单
			single := new(OrderInfo)
			single.OrderId = k.Uuid
			single.Destription = k.Description
			single.AskerName = k.CustomerName
			single.AskerOpenId = k.CustomerOpenId
			single.AskerHeadImg = k.AskerHeadImg
			if k.IsSolved == "2" {
				single.LawyerId = k.AnswerId
				single.LawyerName = k.AnswerName
				single.LawyerHeadImg = k.AnswerHeadImg
				single.LawyerOpenId = k.AnswerOpenId
				single.VoicePath = k.VoicePath
				rank := k.RankInfo
				rankInt, _ := strconv.ParseInt(rank, 10, 64)
				single.RankInfo = rankInt
				single.Pv = k.Pv
				single.QuestionType = k.Category
			}
			retList = append(retList, *single)
		}
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.Total = count
	response.List = retList
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}
type OrderDetailInfo struct {
	QuestionId string `json:"questionId"`
}
type OrderDetailResponse struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
	OrderInfo `json:"orderInfo"`
}
func GetOrderDetailById(ctx *macaron.Context)string{
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" {
		//这里直接调取util重新过一次绿叶 获取openId 等信息
	}
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]

	log.Println("=========>>>>>>,用户OPENID 为", openId)
	log.Println("=========>>>>>>,用户类型为", userType)

	req :=new(OrderDetailInfo)
	body, _ := ctx.Req.Body().String()
	json.Unmarshal([]byte(body),req)
	response :=new(OrderDetailResponse)

	k :=new(model.WechatVoiceQuestions)
	quesionInfoErr:=k.GetConn().Where("uuid = ?",req.QuestionId).Find(&k).Error

	if quesionInfoErr!=nil&&!strings.Contains(quesionInfoErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = quesionInfoErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}
	var single OrderInfo
	single.OrderId = k.Uuid
	single.Destription = k.Description
	single.AskerName = k.CustomerName
	single.AskerOpenId = k.CustomerOpenId
	single.AskerHeadImg = k.AskerHeadImg
	if k.IsSolved == "2" {
		single.LawyerId = k.AnswerId
		single.LawyerName = k.AnswerName
		single.LawyerHeadImg = k.AnswerHeadImg
		single.LawyerOpenId = k.AnswerOpenId
		single.VoicePath = k.VoicePath
		rank := k.RankInfo
		rankInt, _ := strconv.ParseInt(rank, 10, 64)
		single.RankInfo = rankInt
		single.Pv = k.Pv
		single.QuestionType = k.Category
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.OrderInfo = single
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}