package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/macaron"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wechatvoice/model"
	"wechatvoice/tool/util"
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
	response := new(OrderListResponse)
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

	log.Println("=========>>>>>>,用户OPENID 为", openId)
	log.Println("=========>>>>>>,用户类型为", userType)

	body, _ := ctx.Req.Body().String()

	req := new(OrderListFrontRequest)

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
	QuestionId string `json:"orderId"`
}
type OrderDetailResponse struct {
	Code      int64  `json:"code"`
	Msg       string `json:"msg"`
	OrderInfo `json:"orderInfo"`
}

var userOrderList = "http://www.mylvfa.com/daodaolaw/user-order.html"

func ToUserOrders(ctx *macaron.Context) {
	fmt.Println("=================进入方法")
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")

	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/order/touserorder"
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
	ctx.Resp.Write([]byte("<script type=\"text/javascript\">window.location.href=\"" + userOrderList + "\"</script>"))
}

func GetOrderDetailById(ctx *macaron.Context) string {
	response := new(OrderDetailResponse)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/orderdetail"
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
	fmt.Println(cookieStr)
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]

	log.Println("=========>>>>>>,用户OPENID 为", openId)
	log.Println("=========>>>>>>,用户类型为", userType)

	req := new(OrderDetailInfo)
	body, _ := ctx.Req.Body().String()
	json.Unmarshal([]byte(body), req)

	k := new(model.WechatVoiceQuestions)
	quesionInfoErr := k.GetConn().Where("uuid = ?", req.QuestionId).Find(&k).Error

	if quesionInfoErr != nil && !strings.Contains(quesionInfoErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = quesionInfoErr.Error()
		ret_str, _ := json.Marshal(response)
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
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

type LawyerOrderListReq struct {
	StartLine int64  `json:"startNum"`
	EndLine   int64  `json:"endNum"`
	OrderType string `json:"orderType"`
}
type LawyerOrderListResponse struct {
	Code int64      `json:"code"`
	Msg  string     `json:"msg"`
	List []LawOrder `json:"list"`
}

type LawOrder struct {
	OrderId string `json:"orderId"`
	Status  string `json:"status"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Time    string `json:"time"`
	Price   int64  `json:"price"`
	Answer  string `json:"answer"`
}

func GetLayerOrderList(ctx *macaron.Context) string {
	response := new(LawyerOrderListResponse)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/front/lawyerlist"
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
	fmt.Println(cookieStr)
	openId := strings.Split(cookieStr, "|")[0]
	userType := strings.Split(cookieStr, "|")[1]

	log.Println("=========>>>>>>,用户OPENID 为", openId)
	log.Println("=========>>>>>>,用户类型为", userType)

	body, _ := ctx.Req.Body().String()

	req := new(LawyerOrderListReq)

	marshallErr := json.Unmarshal([]byte(body), req)

	if marshallErr != nil {
		response.Code = CODE_ERROR
		response.Msg = marshallErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	//这里要区分下
	list := make([]model.WechatVoiceQuestions, 0)
	var err error
	lawyer := new(model.LawyerInfo)
	lawyerErr := lawyer.GetConn().Where("open_id = ?", openId).Find(&lawyer).Error
	if lawyerErr != nil && !strings.Contains(lawyerErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = lawyerErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	switch req.OrderType {
	case "0":
		//带解答
		list, err = model.GetLawyerQs(lawyer.FirstCategoryId, req.OrderType, req.StartLine, req.EndLine)
		if err != nil && !strings.Contains(err.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = err.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		if int64(len(list)) != (req.EndLine - req.StartLine + 1) {
			a := req.EndLine - int64(len(list))
			list1, list1Err := model.GetNotSpectial(lawyer.FirstCategoryId, req.OrderType, req.StartLine, a)
			if list1Err != nil && !strings.Contains(list1Err.Error(), RNF) {
				response.Code = CODE_ERROR
				response.Msg = list1Err.Error()
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
			for _, k := range list1 {
				list = append(list, k)
			}
		}
	case "1":
		list, _, err = model.QueryLawyerQuestions(req.StartLine, req.EndLine, openId)
	}
	retList := make([]LawOrder, 0)
	for _, k := range list {
		single := new(LawOrder)
		single.OrderId = k.Uuid
		single.Status = k.IsSolved
		single.Content = k.Description
		single.Type = k.Category
		single.Time = k.CreateTime

		price, _ := strconv.ParseInt(k.PaymentInfo, 10, 64)
		single.Price = price

		single.Answer = k.VoicePath
		/**
			OrderId string `json:"orderId"`
		Status string `json:"status"`
		Content string `json:"content"`
		Type string `json:"type"`
		Time string `json:"time"`
		Price int64 `json:"price"`
		Answer string `json:"answer"`
		*/
		retList = append(retList, *single)
	}
	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	response.List = retList
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

type MemberRequest struct {
	StartNum  int64  `json:"startNum"`
	EndNum    int64  `json:"endNum"`
	OrderType string `json:"orderType"`
}

type MemberListReponse struct {
	Code int64         `json:"code"`
	Msg  string        `json:"msg"`
	List []MemberOrder `json:"list"`
}
type MemberOrder struct {
	OrderId string `json:"orderId"`
	Status  string `json:"status"`
	Content string `json:"content"`
	Type    string `json:"type"`
	TypeId  string `json:"typeId"`
	Time    string `json:"time"`
	Price   int64  `json:"price"`
	AddNum  int64  `json:"addNum"`
	Answer  string `json:"answer"`
	CanEval bool   `json:"canEval"`
}

func GetMemberOrderList(ctx *macaron.Context) string {
	response := new(MemberListReponse)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr == "" && ctx.Query("code") == "" {
		re := "http://www.mylvfa.com/voice/ucenter/userlist"
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
	//fmt.Println(cookieStr)
	openId := strings.Split(cookieStr, "|")[0]
	//userType := strings.Split(cookieStr, "|")[1]

	log.Println("=========>>>>>>,用户OPENID 为", openId)
	//log.Println("=========>>>>>>,用户类型为", userType)

	body, _ := ctx.Req.Body().String()

	req := new(MemberRequest)
	fmt.Println("=======>>>>>>请求数据wei", body)
	marshallErr := json.Unmarshal([]byte(body), req)

	if marshallErr != nil {
		response.Code = CODE_ERROR
		response.Msg = marshallErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	list, err := model.GetCustomerInfo(openId, req.OrderType, req.StartNum, req.EndNum)
	fmt.Println(len(list))
	if err != nil && !strings.Contains(err.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	retList := make([]MemberOrder, 0)
	for _, k := range list {
		single := new(MemberOrder)
		single.OrderId = k.Uuid
		single.Status = k.IsSolved
		single.Content = k.Description
		single.TypeId = k.CategoryId
		single.Type = k.Category
		single.Time = k.CreateTime
		single.Answer = k.VoicePath
		single.AddNum = k.AppenQuestionTime
		price, _ := strconv.ParseInt(k.PaymentInfo, 10, 64)
		single.Price = price
		var a bool
		if k.IsRanked == "1" {
			a = false
		} else {
			a = true
		}
		single.CanEval = a
		retList = append(retList, *single)
	}
	response.Code = CODE_SUCCESS
	response.Msg = "ok"
	response.List = retList
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

func EvalAnswers(ctx *macaron.Context) string {
	return ""
}
