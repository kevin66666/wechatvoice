package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/macaron"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"wechatvoice/model"
	"wechatvoice/tool/util"

	"github.com/robfig/cron"
)

func init() {
	c := cron.New()
	c.AddFunc("@every 1m", func() { UpdateAllQuestions() })
	c.Start()
}

func UpdateAllQuestions() {
	list, err := model.GetAllLocked()
	if err != nil && !strings.Contains(err.Error(), RNF) {
		fmt.Println(err.Error())
	}
	if len(list) > 0 {

		for _, k := range list {
			log.Println("=====>>>", k.Uuid)
			go UpdateInfo(k)
		}
	}
}
func UpdateInfo(info model.WechatVoiceQuestions) {
	now := time.Now().Unix()
	// fmt.Println(no
	b4 := info.LockTime
	a := now - b4
	times := a / 60
	fmt.Println(now)
	fmt.Println(b4)
	fmt.Println(a)
	fmt.Println(times)
	if times > 10 {
		info.IsLocked = "0"
		info.LockTime = 0
		err := info.GetConn().Update(&info).Error
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

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
	CanDelete    bool   `json:"canDelete"`
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

var userLawList = "http://www.mylvfa.com/daodaolaw/laywer-order.html"
var front = "http://www.mylvfa.com/voice/front/toindex"

func ToLawOrders(ctx *macaron.Context) {
	fmt.Println("=================进入方法")
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	if cookieStr != "" {
		openId := strings.Split(cookieStr, "|")[0]
		fmt.Println(openId)
		u := new(model.LawyerInfo)
		erru := u.GetConn().Where("open_id = ?", openId).Find(&u).Error
		if erru != nil && !strings.Contains(erru.Error(), RNF) {
			fmt.Println(erru.Error)
		}
		if u.Uuid == "" {
			//如果有cookie 可能是用户 可能是律师 需要做区分
			ctx.Resp.Write([]byte("<script type=\"text/javascript\">window.location.href=\"" + front + "\"</script>"))
		} else {
			ctx.Resp.Write([]byte("<script type=\"text/javascript\">window.location.href=\"" + userLawList + "\"</script>"))
		}
	} else {
		if cookieStr == "" && ctx.Query("code") == "" {
			re := "http://www.mylvfa.com/voice/order/tolaworder"
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
			openId := res1.OpenId
			u := new(model.LawyerInfo)
			erru := u.GetConn().Where("open_id = ?", openId).Find(&u).Error
			if erru != nil && !strings.Contains(erru.Error(), RNF) {
				fmt.Println(erru.Error)
			}
			if u.Uuid == "" {
				//说明这个人已经注册进来了
				list1, err1 := model.GetUserInfoByOpenId(openId)
				if err1 != nil {
					fmt.Print(err1.Error())
				}
				var userid string
				fmt.Println(len(list1))
				if len(list1) == 0 {
					fmt.Println("律师表中没有这个人")
					ctx.Resp.Write([]byte("<script type=\"text/javascript\">window.location.href=\"" + front + "\"</script>"))
				}
				for _, k := range list1 {
					userid = string(k["userID"])
				}

				fmt.Println("userid is ====>>>", userid)

				list2, err2 := model.GetLawerInfoById(userid)
				if err2 != nil {
					fmt.Println(err2.Error())
				}
				if len(list2) == 0 {
					fmt.Println("律师表中没有这个人")
					ctx.Resp.Write([]byte("<script type=\"text/javascript\">window.location.href=\"" + front + "\"</script>"))
				} else {
					//把律师信息拉到我这边
					var lawerId, lawyerName, lawyerPhone, cert, groupPhoto, singlePhoto, province, city, lawFirm, business, desc, createDate string
					for _, k := range list2 {
						lawerId = string(k["lawyerId"])
						lawyerName = string(k["lawyerName"])
						lawyerPhone = string(k["lawyerPhone"])
						cert = string(k["lawyerCertificateNo"])
						groupPhoto = string(k["groupPhoto"])
						singlePhoto = string(k["singlePhoto"])
						province = string(k["selProvince"])
						city = string(k["selCity"])
						lawFirm = string(k["lawFirm"])
						business = string(k["goodAtBusiness"])
						desc = string(k["description"])
						createDate = string(k["createDate"])

					}
					fmt.Println(list2)
					fmt.Println(province, city)
					lawInfo := new(model.LawyerInfo)
					lawInfo.Uuid = lawerId
					lawInfo.RegistTime = createDate
					lawInfo.OpenId = openId
					lawInfo.PhoneNumber = lawyerPhone
					fmt.Println("==============>>>>>")
					fmt.Println(singlePhoto)
					fmt.Println("====================>>>>")
					photo := singlePhoto
					var a string
					if photo != "" {
						list := strings.Split("lawyer", photo)
						a = photo[14:]
						//fmt.Println(b)
						fmt.Println(list)
					}
					lawInfo.Uuid = util.GenerateUuid()
					lawInfo.HeadImgUrl = "images/" + photo
					lawInfo.Name = lawyerName
					lawInfo.FirstCategory = business
					lawInfo.Cet = cert
					lawInfo.GroupPhoto = groupPhoto
					lawInfo.LawFirm = lawFirm
					lawInfo.Desc = desc
					from := "/usr/local/apache-tomcat-6.0.32/webapps/mylawyerfriend" + singlePhoto
					to := "/home/workspace_go/src/wechatvoice/daodaolaw/images"
					cmd := exec.Command("cp", from, to)
					errcmd := cmd.Run()
					if errcmd != nil {
						fmt.Println(errcmd.Error())
					}
					toP := "images" + a
					lawInfo.HeadImgUrl = toP
					err := lawInfo.GetConn().Create(&lawInfo)
					if err != nil {
						fmt.Println(err.Error)
					}
					ctx.SetSecureCookie("userloginstatus", openId+"|0")
				}
			} else {
				ctx.SetSecureCookie("userloginstatus", openId+"|0")
			}
			// ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
			// member := new(model.MemberInfo)
			// memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
			// if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
			// 	fmt.Println(memberErr.Error(), "=====会员出错")
			// }
			// if member.Uuid == "" {
			// 	fmt.Println("新的用户")
			// 	user := GetUserInfo(res1.OpenId, res1.AccessToken)
			// 	member.Uuid = util.GenerateUuid()
			// 	member.HeadImgUrl = user.HeadImgUrl
			// 	member.OpenId = user.OpenId
			// 	member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
			// 	member.NickName = user.NickName
			// 	err := member.GetConn().Create(&member).Error
			// 	if err != nil {
			// 		fmt.Println(err.Error(), "xxxxx")
			// 	}
			// }
			//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
		}
		fmt.Println(cookieStr)
		openId := strings.Split(cookieStr, "|")[0]
		// userType := strings.Split(cookieStr, "|")[1]
		fmt.Println(openId)
		ctx.Resp.Write([]byte("<script type=\"text/javascript\">window.location.href=\"" + userLawList + "\"</script>"))
	}
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
	//status: //0 代表抢答（直接提问的） 1代表指定提问  2代表追问
	Content   string `json:"content"`
	Type      string `json:"type"`
	Time      string `json:"time"`
	Price     int64  `json:"price"`
	Answer    string `json:"answer"`
	IsPlay    bool   `json:"isPlay"`
	CanDelete bool   `json:"canDelete"`
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
	law := new(model.LawyerInfo)
	lawErr := law.GetConn().Where("open_id = ?", openId).Find(&law).Error
	if lawErr != nil && !strings.Contains(lawErr.Error(), RNF) {
		fmt.Println(lawErr)
	}
	switch req.OrderType {
	case "0":
		//带解答
		// list, err = model.GetLawyerQs(req.OrderType, law.Uuid, req.StartLine, req.EndLine)
		// if err != nil && !strings.Contains(err.Error(), RNF) {
		// 	response.Code = CODE_ERROR
		// 	response.Msg = err.Error()
		// 	ret_str, _ := json.Marshal(response)
		// 	return string(ret_str)
		// }
		// if int64(len(list)) != (req.EndLine - req.StartLine + 1) {
		// 	a := req.EndLine - int64(len(list))
		// 	list1, list1Err := model.GetNotSpectial(lawyer.FirstCategoryId, req.OrderType, req.StartLine, a)
		// 	if list1Err != nil && !strings.Contains(list1Err.Error(), RNF) {
		// 		response.Code = CODE_ERROR
		// 		response.Msg = list1Err.Error()
		// 		ret_str, _ := json.Marshal(response)
		// 		return string(ret_str)
		// 	}
		// 	for _, k := range list1 {
		// 		list = append(list, k)
		// 	}
		// }
		//指定问题部分
		directList, dirErr := model.GetLawerDirectInfo(openId, req.StartLine, req.EndLine)
		if dirErr != nil && !strings.Contains(dirErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = dirErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		log.Println(len(directList)) //这里打印输出下 有多少条
		dirLen := len(directList)
		if dirLen > 0 {
			//先把这些数据丢到list中
			for _, k := range directList {
				list = append(list, k)
			}
		}
		aList := make([]model.WechatVoiceQuestions, 0)
		var aErr error
		if int64(dirLen) < (req.EndLine - req.StartLine) {
			//说明数量不够 需要后期去补充
			need := req.EndLine - int64(dirLen)
			fmt.Println("=======>>>>>>")
			aList, aErr = model.GetLaerOther(req.StartLine, need)
		}
		if aErr != nil && !strings.Contains(aErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = aErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}

		fmt.Println("-===asdasdasd")
		fmt.Println(aList)
		fmt.Println("asdasdasdasd")
		for _, k := range aList {
			list = append(list, k)
		}
	case "2":
		list, _, err = model.QueryLawyerQuestions(req.StartLine, req.EndLine, openId)
	}
	if err != nil && !strings.Contains(err.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	fmt.Println("========>>>>>>>>>>订单l列表")
	fmt.Println(len(list))
	log.Println(list)
	fmt.Println("========>>>>>>>>>>订单l列表")
	retList := make([]LawOrder, 0)
	for _, k := range list {
		single := new(LawOrder)
		single.OrderId = k.Uuid
		single.Status = k.IsSolved
		single.Content = k.Description
		single.Type = k.Category
		single.Time = k.CreateTime[0:10]

		var flag bool
		if k.IsSolved == "2" {
			flag = true
		} else {
			flag = false
		}
		single.CanDelete = flag
		price, _ := strconv.ParseInt(k.PaymentInfo, 10, 64)
		single.Price = price

		single.Answer = k.VoicePath
		single.IsPlay = true
		var status string
		//status: //0 代表抢答（直接提问的） 1代表指定提问  2代表追问

		if k.QType == "1" {
			//追加
			status = "2"
		} else if k.QType == "2" {
			//指定
			status = "1"
		} else {
			status = "0"
		}
		single.Status = status
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
	OrderId   string `json:"orderId"`
	Status    string `json:"status"`
	Content   string `json:"content"`
	Type      string `json:"typeName"`
	TypeId    string `json:"typeId"`
	Time      string `json:"time"`
	Price     int64  `json:"price"`
	AddNum    int64  `json:"addNum"`
	Answer    string `json:"answer"`
	CanEval   bool   `json:"canEval"`
	LawyerId  string `json:"laywerId"`
	IsPlay    bool   `json:"isPlay"`
	CanDelete bool   `json:"canDelete"`
}

func GetMemberOrderList(ctx *macaron.Context) string {
	response := new(MemberListReponse)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	// if cookieStr == "" && ctx.Query("code") == "" {
	// 	re := "http://www.mylvfa.com/voice/ucenter/userlist"
	// 	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
	// 	//cookieStr = "1|2"
	// 	ctx.Redirect(url)
	// }
	// code := ctx.Query("code")
	// if code != "" {
	// 	url := "http://60.205.4.26:22334/getOpenid?code=" + code
	// 	res, err := http.Get(url)
	// 	if err != nil {
	// 		fmt.Println("=========xxxxx")
	// 		fmt.Println(err.Error())
	// 	}
	// 	resBody, _ := ioutil.ReadAll(res.Body)
	// 	fmt.Println(string(resBody))
	// 	defer res.Body.Close()
	// 	fmt.Println("==========>>>>")
	// 	res1 := new(OpenIdResponse)
	// 	json.Unmarshal(resBody, res1)
	// 	ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
	// 	member := new(model.MemberInfo)
	// 	memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
	// 	if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
	// 		response.Code = CODE_ERROR
	// 		response.Msg = memberErr.Error()
	// 		ret_str, _ := json.Marshal(res)
	// 		return string(ret_str)
	// 	}
	// 	if member.Uuid == "" {
	// 		fmt.Println("新的用户")
	// 		user := GetUserInfo(res1.OpenId, res1.AccessToken)
	// 		member.Uuid = util.GenerateUuid()
	// 		member.HeadImgUrl = user.HeadImgUrl
	// 		member.OpenId = user.OpenId
	// 		member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
	// 		member.NickName = user.NickName
	// 		err := member.GetConn().Create(&member).Error
	// 		if err != nil {
	// 			response.Code = CODE_ERROR
	// 			response.Msg = err.Error()
	// 			ret_str, _ := json.Marshal(response)
	// 			return string(ret_str)
	// 		}
	// 	}
	// 	//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
	// }
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
	if req.OrderType == "-1" {
		req.OrderType = "2"
	}
	//用户的逻辑来说 需要将用户删除的订单排除
	logList, logListErr := model.GetUserDeletedList(openId)
	if logListErr != nil && !strings.Contains(logListErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = logListErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	log.Println(len(logList))
	logIdList := make([]string, 0)
	for _, k := range logList {
		logIdList = append(logIdList, k.OrderId)
	}
	log.Println(logIdList)
	//备用
	//list, err := model.GetCustomerInfoNew(openId, req.OrderType, logIdList, req.StartNum, req.EndNum)
	// list, err := model.GetCustomerInfo(openId, req.OrderType, req.StartNum, req.EndNum)
	//未完成订单 我就不管他   已完成 的 需要进行筛选
	// fmt.Println(len(list))
	// if err != nil && !strings.Contains(err.Error(), RNF) {
	// 	response.Code = CODE_ERROR
	// 	response.Msg = err.Error()
	// 	ret_str, _ := json.Marshal(response)
	// 	return string(ret_str)
	// }
	list := make([]model.WechatVoiceQuestions, 0)
	var err error
	retList := make([]MemberOrder, 0)
	switch req.OrderType {
	case "0":
		list, err = model.GetCustomerInfo(openId, req.OrderType, req.StartNum, req.EndNum)

	case "2":
		idList, idErr := model.GetPaymentQuery(openId)
		if idErr != nil && !strings.Contains(idErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = idErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		id := make([]string, 0)
		for _, k := range idList {
			id = append(id, k.QuestionId)
		}
		list, err = model.GetCustomerPaiedInfo(openId, id, logIdList, req.StartNum, req.EndNum)

	}
	fmt.Println(len(list))
	if err != nil && !strings.Contains(err.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	for _, k := range list {
		single := new(MemberOrder)
		single.OrderId = k.Uuid
		single.Status = k.IsSolved
		single.Content = k.Description
		single.TypeId = k.CategoryId
		single.Type = k.Category
		single.Time = k.CreateTime[0:10]
		single.Answer = k.VoicePath
		single.IsPlay = true

		l, errs := model.GetInfos(openId, k.Uuid)
		if errs != nil && !strings.Contains(errs.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = errs.Error()
			ret_Str, _ := json.Marshal(response)
			return string(ret_Str)
		}
		var can bool
		if k.IsSolved == "0" {
			//没有解答的 肯定可以删除
			can = true
		} else {
			//已解答的
			if k.CustomerOpenId == openId {
				//用户自己的订单 可以删除
				can = true
			} else {
				can = false
			}
		}
		single.CanDelete = can
		if k.ParentQuestionId != "" {
			single.AddNum = int64(0)
		} else {
			single.AddNum = int64(2) - int64(len(l))
		}
		price, _ := strconv.ParseInt(k.PaymentInfo, 10, 64)
		single.Price = price
		single.LawyerId = k.AnswerId
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

type EvaluateAnswers struct {
	OrderId string `json:"orderId"`
	Number  string `json:"start"`
}

type RedpacketResponse struct {
	Code      int64  `json:"code"`
	Msg       string `json:"msg"`
	RedPacket string `json:"redPacket"`
}

func EvalAnswers(ctx *macaron.Context) string {
	body, _ := ctx.Req.Body().String()
	req := new(EvaluateAnswers)
	fmt.Println("-========----pingjia request ", body)
	json.Unmarshal([]byte(body), req)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	openId := strings.Split(cookieStr, "|")[0]
	log.Println(openId)
	response := new(RedpacketResponse)

	orderInfo := new(model.WechatVoiceQuestions)
	orderErr := orderInfo.GetConn().Where("uuid = ?", req.OrderId).Where("customer_open_id  = ?", openId).Find(&orderInfo).Error

	if orderErr != nil && !strings.Contains(orderErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = orderErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	if orderInfo.Uuid == "" || orderInfo.IsRanked == "1" || orderInfo.IsSolved == "0" || orderInfo.IsSolved == "1" {
		response.Code = CODE_ERROR
		response.Msg = "error"
		ret_Str, _ := json.Marshal(response)
		return string(ret_Str)
	}
	cateId := orderInfo.CategoryId
	setting := new(model.WechatVoiceQuestionSettings)
	settingErr := setting.GetConn().Where("category_id = ?", cateId).Find(&setting).Error

	if settingErr != nil && !strings.Contains(settingErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = settingErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	amountF, _ := strconv.ParseFloat(orderInfo.PaymentInfo, 64)
	amountF = amountF * 100
	lp, _ := strconv.ParseFloat(setting.LawyerFeePercent, 64)
	red := 100.00 - lp
	amountLeft := (amountF * red) / 100
	amount := int64(amountLeft)
	redint := rand.Int63n(amount)
	// log.Println(redStr)
	redIntStr := strconv.FormatInt(redint, 10)
	redF, _ := strconv.ParseFloat(redIntStr, 64)
	redFr := redF / 100
	redStr := strconv.FormatFloat(redFr, 'f', 2, 64)

	orderInfo.IsRanked = "1"
	orderInfo.RankInfo = req.Number
	updateErr := orderInfo.GetConn().Save(&orderInfo).Error
	if updateErr != nil && !strings.Contains(updateErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = updateErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	//给律师发红包
	a := amountF * lp
	astr := strconv.FormatFloat(a, 'f', 2, 64)
	reds := new(RedPackages)
	reds.Act_name = "发送红包"
	reds.Client_ip = "127.0.0.1"
	reds.Remark = ""
	reds.Re_openid = orderInfo.AnswerOpenId
	reds.Nick_name = "叨叨律法"
	reds.SendNickName = orderInfo.AnswerName
	reds.Wishing = "您的订单已完成"
	reds.Amount = int64(a)
	reds.MpId = ""
	fmt.Println(reds)
	// suc, strsuc := SendRedPacket(reds)
	// fmt.Println("===========================", suc, strsuc)
	//记录律师信息
	law := new(model.LawyerInfo)
	lawErr := law.GetConn().Where("uuid = ?", orderInfo.AnswerId).Find(&law).Error
	if lawErr != nil && !strings.Contains(lawErr.Error(), RNF) {
		fmt.Println(lawErr.Error())
	}
	switch req.Number {
	case "1":
		law.RankFirst = law.RankFirst + 1
	case "2":
		law.RankSecond = law.RankSecond + 1
	case "3":
		law.RankThird = law.RankThird + 1
	case "4":
		law.RankFouth = law.RankFouth + 1
	case "5":
		law.RankLast = law.RankLast + 1
	}
	lawErr = law.GetConn().Save(&law).Error
	if lawErr != nil && !strings.Contains(lawErr.Error(), RNF) {
		fmt.Println(lawErr)
	}
	//记录钱的信息
	pay := new(model.OrderPaymentInfo)
	pay.GetConn().Where("order_number = ?", req.OrderId).Where("open_id = ?", openId).Where("is_first = 1").Find(&pay)
	payment := new(model.WechatVoicePaymentInfo)
	payment.Uuid = util.GenerateUuid()
	payment.SwiftNumber = pay.WeixinSwiftNumber
	payment.MemberId = orderInfo.CustomerId
	payment.OpenId = openId
	payment.RedPacketAmount = redStr
	payment.LawyerAmount = astr
	payment.OrderId = req.OrderId
	errPay := payment.GetConn().Create(&payment).Error
	if errPay != nil {
		fmt.Println(errPay)
	}
	// payment.SwiftNumber = orderInfo.

	//保存orderInfo
	orderInfo.IsRanked = "1"
	orderInfo.IsSolved = "2"
	star := req.Number
	orderInfo.RankInfo = star
	orderUpdateErr := orderInfo.GetConn().Save(&orderInfo).Error
	if orderUpdateErr != nil && !strings.Contains(orderUpdateErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = orderUpdateErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	response.Code = CODE_SUCCESS
	response.Msg = "ok"
	response.RedPacket = redStr
	ret_str, _ := json.Marshal(response)

	return string(ret_str)
}

type QuestionId struct {
	OrderId string `json:"orderId"`
}

type CheckResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func GetQuestionToAnswer(ctx *macaron.Context) string {
	body, _ := ctx.Req.Body().String()
	req := new(QuestionId)
	json.Unmarshal([]byte(body), req)
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	openId := strings.Split(cookieStr, "|")[0]
	log.Println(openId)

	response := new(CheckResponse)
	lock := new(model.AnswerLockInfo)
	lockerr := lock.GetConn().Where("question_id = ?", req.OrderId).Where("open_id = ?", openId).Find(&lock).Error

	if lockerr != nil && !strings.Contains(lockerr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = lockerr.Error()
		ret_str, _ := json.Marshal(response)
		fmt.Print(string(ret_str), "1")
		return string(ret_str)
	}
	if lock.Uuid != "" {
		response.Code = CODE_SUCCESS
		response.Msg = "ok"
		ret_str, _ := json.Marshal(response)
		fmt.Print(string(ret_str), "2")

		return string(ret_str)
	}

	lockList, lockListErr := model.GetLockListById(req.OrderId)
	if lockListErr != nil && !strings.Contains(lockListErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = lockListErr.Error()
		ret_str, _ := json.Marshal(response)
		fmt.Print(string(ret_str), "3")

		return string(ret_str)
	}

	if len(lockList) == 0 || len(lockList) == 1 {
		question := new(model.WechatVoiceQuestions)
		qErr := question.GetConn().Where("uuid =?", req.OrderId).Find(&question).Error
		if qErr != nil && !strings.Contains(qErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = qErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		if question.IsLocked == "0" {
			response.Code = CODE_SUCCESS
			response.Msg = "ok"

			question.IsLocked = "1"
			question.LockTime = time.Now().Unix()
			errsss := question.GetConn().Save(&question).Error
			if errsss != nil {
				response.Msg = errsss.Error()
				response.Code = CODE_ERROR
				ret_str, _ := json.Marshal(response)
				return string(ret_str)
			}
			ret_str, _ := json.Marshal(response)
			fmt.Print(string(ret_str), "4")

			return string(ret_str)
		} else {
			response.Code = CODE_ERROR
			response.Msg = ""
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}

	} else {
		response.Code = CODE_ERROR
		response.Msg = "error"
		ret_str, _ := json.Marshal(response)
		fmt.Print(string(ret_str), "5")

		return string(ret_str)
	}

}

// func GetMemberOrderListNew(ctx *macaron.Context) string {
// 	response := new(MemberListReponse)
// 	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
// 	if cookieStr == "" && ctx.Query("code") == "" {
// 		re := "http://www.mylvfa.com/voice/ucenter/userlist"
// 		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxac69efc11c5e182f&redirect_uri=" + re + "&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
// 		//cookieStr = "1|2"
// 		ctx.Redirect(url)
// 	}
// 	code := ctx.Query("code")
// 	if code != "" {
// 		url := "http://60.205.4.26:22334/getOpenid?code=" + code
// 		res, err := http.Get(url)
// 		if err != nil {
// 			fmt.Println("=========xxxxx")
// 			fmt.Println(err.Error())
// 		}
// 		resBody, _ := ioutil.ReadAll(res.Body)
// 		fmt.Println(string(resBody))
// 		defer res.Body.Close()
// 		fmt.Println("==========>>>>")
// 		res1 := new(OpenIdResponse)
// 		json.Unmarshal(resBody, res1)
// 		ctx.SetSecureCookie("userloginstatus", res1.OpenId+"|0")
// 		member := new(model.MemberInfo)
// 		memberErr := member.GetConn().Where("open_id = ?", res1.OpenId).Find(&member).Error
// 		if memberErr != nil && !strings.Contains(memberErr.Error(), RNF) {
// 			response.Code = CODE_ERROR
// 			response.Msg = memberErr.Error()
// 			ret_str, _ := json.Marshal(res)
// 			return string(ret_str)
// 		}
// 		if member.Uuid == "" {
// 			fmt.Println("新的用户")
// 			user := GetUserInfo(res1.OpenId, res1.AccessToken)
// 			member.Uuid = util.GenerateUuid()
// 			member.HeadImgUrl = user.HeadImgUrl
// 			member.OpenId = user.OpenId
// 			member.RegistTime = time.Unix(time.Now().Unix(), 0).String()[0:19]
// 			member.NickName = user.NickName
// 			err := member.GetConn().Create(&member).Error
// 			if err != nil {
// 				response.Code = CODE_ERROR
// 				response.Msg = err.Error()
// 				ret_str, _ := json.Marshal(response)
// 				return string(ret_str)
// 			}
// 		}
// 		//ctx.Redirect("http://www.mylvfa.com/voice/front/getcatList")
// 	}
// 	fmt.Println(cookieStr)
// 	//fmt.Println(cookieStr)
// 	openId := strings.Split(cookieStr, "|")[0]
// 	//userType := strings.Split(cookieStr, "|")[1]

// 	log.Println("=========>>>>>>,用户OPENID 为", openId)
// 	//log.Println("=========>>>>>>,用户类型为", userType)

// 	body, _ := ctx.Req.Body().String()

// 	req := new(MemberRequest)
// 	fmt.Println("=======>>>>>>请求数据wei", body)
// 	marshallErr := json.Unmarshal([]byte(body), req)

// 	if marshallErr != nil {
// 		response.Code = CODE_ERROR
// 		response.Msg = marshallErr.Error()
// 		ret_str, _ := json.Marshal(response)
// 		return string(ret_str)
// 	}
// 	if req.OrderType == "-1" {
// 		req.OrderType = "2"
// 	}
// 	retList := make([]MemberOrder, 0)
// 	list := make([]model.WechatVoiceQuestions, 0)
// 	logList, logListErr := model.GetUserDeletedList(openId)
// 	if logListErr != nil && !strings.Contains(logListErr.Error(), RNF) {
// 		response.Code = CODE_ERROR
// 		response.Msg = logListErr.Error()
// 		ret_str, _ := json.Marshal(response)
// 		return string(ret_str)
// 	}
// 	log.Println(len(logList))
// 	logIdList := make([]string, 0)
// 	for _, k := range logList {
// 		logIdList = append(logIdList, k.OrderId)
// 	}
// 	log.Println(logIdList)
// 	//备用
// 	//list, err := model.GetCustomerInfoNew(openId, req.OrderType, logIdList, req.StartNum, req.EndNum)
// 	/**
// 	用户删除  userdelete 1
// 	用户未完成的订单中  userdelete  1

// 	如果订单完成了     userdelete  1      然后加一条log
// 	如果订单未完成   那么删就删吧
// 	*/
// 	var err error
// 	switch req.OrderType {
// 	case "0":
// 		list, err = model.GetCustomerInfo(openId, req.OrderType, req.StartNum, req.EndNum)

// 	case "2":
// 		idList, idErr := model.GetPaymentQuery(openId)
// 		if idErr != nil && !strings.Contains(idErr.Error(), RNF) {
// 			response.Code = CODE_ERROR
// 			response.Msg = idErr.Error()
// 			ret_str, _ := json.Marshal(response)
// 			return string(ret_str)
// 		}
// 		id := make([]string, 0)
// 		for _, k := range idList {
// 			id = append(id, k.QuestionId)
// 		}
// 		list, err = model.GetCustomerPaiedInfo(openId, id, req.StartNum, req.EndNum)

// 	}
// 	fmt.Println(len(list))
// 	if err != nil && !strings.Contains(err.Error(), RNF) {
// 		response.Code = CODE_ERROR
// 		response.Msg = err.Error()
// 		ret_str, _ := json.Marshal(response)
// 		return string(ret_str)
// 	}

// 	for _, k := range list {
// 		single := new(MemberOrder)
// 		single.OrderId = k.Uuid
// 		single.Status = k.IsSolved
// 		single.Content = k.Description
// 		single.TypeId = k.CategoryId
// 		single.Type = k.Category
// 		single.Time = k.CreateTime
// 		single.Answer = k.VoicePath
// 		single.IsPlay = true
// 		l, errs := model.GetInfos(openId, k.Uuid)
// 		if errs != nil && !strings.Contains(errs.Error(), RNF) {
// 			response.Code = CODE_ERROR
// 			response.Msg = errs.Error()
// 			ret_Str, _ := json.Marshal(response)
// 			return string(ret_Str)
// 		}

// 		if k.ParentQuestionId != "" {
// 			single.AddNum = int64(0)
// 		} else {
// 			single.AddNum = int64(2) - int64(len(l))
// 		}
// 		price, _ := strconv.ParseInt(k.PaymentInfo, 10, 64)
// 		single.Price = price
// 		single.LawyerId = k.AnswerId
// 		var a bool
// 		if k.IsRanked == "1" {
// 			a = false
// 		} else {
// 			a = true
// 		}
// 		single.CanEval = a
// 		retList = append(retList, *single)
// 	}
// 	response.Code = CODE_SUCCESS
// 	response.Msg = "ok"
// 	response.List = retList
// 	ret_str, _ := json.Marshal(response)
// 	return string(ret_str)
// }

type DeleteOrderRequest struct {
	OrderId string `json:"orderId"`
}
type Response struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func DeleteOrderInfo(ctx *macaron.Context) string {
	body, _ := ctx.Req.Body().String()
	req := new(DeleteOrderRequest)
	json.Unmarshal([]byte(body), req)
	orderId := req.OrderId
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	cookie := strings.Split(cookieStr, "|")[0]
	fmt.Println(cookie)
	response := new(Response)
	question := new(model.WechatVoiceQuestions)
	questionErr := question.GetConn().Where("uuid = ?", orderId).Find(&question).Error
	if questionErr != nil && !strings.Contains(questionErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = questionErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	if cookie != question.CustomerOpenId {
		response.Code = CODE_ERROR
		response.Msg = "不能删除不是自己的订单"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	if question.IsSolved == "0" {
		//如果订单未完成 删就删吧
		deleteErr := question.GetConn().Delete(&question).Error
		if deleteErr != nil && strings.Contains(deleteErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = deleteErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		response.Code = CODE_ERROR
		response.Msg = "ok"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	} else if question.IsSolved == "2" && question.IsRanked == "0" {
		//订单已有人回答  还没评价 那么不能删除
		response.Code = CODE_ERROR
		response.Msg = "未评价的订单不能删除"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	} else {
		question.UserDelete = "1"
		err := question.GetConn().Save(&question).Error
		if err != nil && !strings.Contains(err.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = err.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		deleteLog := new(model.UserDeleteLogs)
		deleteLog.Uuid = util.GenerateUuid()
		deleteLog.UserOpenId = cookie
		deleteLog.OrderNumber = question.OrderNumber
		deleteLog.Uuid = question.Uuid
		cErr := deleteLog.GetConn().Create(&deleteLog).Error
		if cErr != nil && !strings.Contains(cErr.Error(), RNF) {
			response.Code = CODE_ERROR
			response.Msg = cErr.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		response.Code = CODE_SUCCESS
		response.Msg = "ok"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

}
func LawyerDeleteOrderInfo(ctx *macaron.Context) string {
	body, _ := ctx.Req.Body().String()
	req := new(DeleteOrderRequest)
	json.Unmarshal([]byte(body), req)
	orderId := req.OrderId
	cookieStr, _ := ctx.GetSecureCookie("userloginstatus")
	cookie := strings.Split(cookieStr, "|")[0]
	fmt.Println(cookie)
	response := new(Response)
	question := new(model.WechatVoiceQuestions)
	questionErr := question.GetConn().Where("uuid = ?", orderId).Find(&question).Error
	if questionErr != nil && !strings.Contains(questionErr.Error(), RNF) {
		response.Code = CODE_ERROR
		response.Msg = questionErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	if cookie != question.AnswerOpenId {
		response.Code = CODE_ERROR
		response.Msg = "不能删除不是自己的订单"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	question.LawyerDelete = "1"
	err := question.GetConn().Save(&question).Error
	if err != nil && !strings.Contains(err.Error(), RNF) {
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

func EvalAnswersTest(ctx *macaron.Context) {

	// amountF, _ := strconv.ParseFloat(orderInfo.PaymentInfo, 64)
	// amountF = amountF * 100
	// lp, _ := strconv.ParseFloat(setting.LawyerFeePercent, 64)
	// red := 100.00 - lp
	// amountLeft := (amountF * red) / 100
	// amount := int64(amountLeft)
	// redint := rand.Int63n(amount)
	// // log.Println(redStr)
	// redIntStr := strconv.FormatInt(redint, 10)
	// redF, _ := strconv.ParseFloat(redIntStr, 64)
	// redFr := redF / 100
	// redStr := strconv.FormatFloat(redFr, 'f', 2, 64)

	// orderInfo.IsRanked = "1"
	// orderInfo.RankInfo = req.Number
	// updateErr := orderInfo.GetConn().Save(&orderInfo).Error
	// if updateErr != nil && !strings.Contains(updateErr.Error(), RNF) {
	// 	response.Code = CODE_ERROR
	// 	response.Msg = updateErr.Error()
	// 	ret_str, _ := json.Marshal(response)
	// 	return string(ret_str)
	// }

	//给律师发红包
	// a := amountF * lp
	// astr := strconv.FormatFloat(a, 'f', 2, 64)
	// reds := new(RedPackages)
	reds := new(RedPackages)
	reds.Act_name = "发送红包"
	reds.Client_ip = "127.0.0.1"
	reds.Remark = "ahahahah"
	reds.Re_openid = "o-u0Nv5Rjxrw2EdmYXqzLXi_uTVo"
	reds.Nick_name = "叨叨律法"
	reds.SendNickName = "差不多先生"
	reds.Wishing = "您的订单已完成"
	reds.Amount = int64(100)
	reds.MpId = ""
	fmt.Println(reds)
	suc, strsuc := SendRedPacket(reds)
	fmt.Println("================红包红包红包===============>>>>")
	log.Println(suc)
	log.Println(strsuc)
	fmt.Println("===============================>>>>")
	// fmt.Println("===========================", suc, strsuc)
	//记录律师信息
	// law := new(model.LawyerInfo)
	// lawErr := law.GetConn().Where("uuid = ?", orderInfo.AnswerId).Find(&law).Error
	// if lawErr != nil && !strings.Contains(lawErr.Error(), RNF) {
	// 	fmt.Println(lawErr.Error())
	// }
	// switch req.Number {
	// case "1":
	// 	law.RankFirst = law.RankFirst + 1
	// case "2":
	// 	law.RankSecond = law.RankSecond + 1
	// case "3":
	// 	law.RankThird = law.RankThird + 1
	// case "4":
	// 	law.RankFouth = law.RankFouth + 1
	// case "5":
	// 	law.RankLast = law.RankLast + 1
	// }
	// lawErr = law.GetConn().Save(&law).Error
	// if lawErr != nil && !strings.Contains(lawErr.Error(), RNF) {
	// 	fmt.Println(lawErr)
	// }
	// //记录钱的信息
	// pay := new(model.OrderPaymentInfo)
	// pay.GetConn().Where("order_number = ?", req.OrderId).Where("open_id = ?", openId).Where("is_first = 1").Find(&pay)
	// payment := new(model.WechatVoicePaymentInfo)
	// payment.Uuid = util.GenerateUuid()
	// payment.SwiftNumber = pay.WeixinSwiftNumber
	// payment.MemberId = orderInfo.CustomerId
	// payment.OpenId = openId
	// payment.RedPacketAmount = redStr
	// payment.LawyerAmount = astr
	// payment.OrderId = req.OrderId
	// errPay := payment.GetConn().Create(&payment).Error
	// if errPay != nil {
	// 	fmt.Println(errPay)
	// }
	// // payment.SwiftNumber = orderInfo.

	// //保存orderInfo
	// orderInfo.IsRanked = "1"
	// orderInfo.IsSolved = "2"
	// star := req.Number
	// orderInfo.RankInfo = star
	// orderUpdateErr := orderInfo.GetConn().Save(&orderInfo).Error
	// if orderUpdateErr != nil && !strings.Contains(orderUpdateErr.Error(), RNF) {
	// 	response.Code = CODE_ERROR
	// 	response.Msg = orderUpdateErr.Error()
	// 	ret_str, _ := json.Marshal(response)
	// 	return string(ret_str)
	// }
	// response.Code = CODE_SUCCESS
	// response.Msg = "ok"
	// //response.RedPacket = redStr
	// ret_str, _ := json.Marshal(response)

	// return string(ret_str)
}
