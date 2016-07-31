package handler

import(
	"encoding/json"

	"github.com/Unknwon/macaron"
	"net/http"
	"wechatvoice/model"
	"io/ioutil"
	"wechatvoice/tool/util"
	"strings"
)

const URL  = "http://file.api.weixin.qq.com/cgi-bin/media/get?access_token="

type MediaId struct {
	MId string `json:"mediaId"`
	QuestionId string `json:"questionId"`
}
func GetWechatVoiceInfoFromWechatServer(ctx *macaron.Context)string{
	response :=new(model.GeneralResponse)

	token,_ :=GetAccesstoken(APPID)
	body,_:=ctx.Req.Body().String()

	req :=new(MediaId)
	json.Unmarshal([]byte(body),req)
	http://file.api.weixin.qq.com/cgi-bin/media/get?access_token=ACCESS_TOKEN&media_id=MEDIA_ID

	url :=URL+token+"&media_id="+req.MId

	resp ,err:=http.Get(url)
	defer resp.Body.Close()

	if err!=nil{
		response.Code = CODE_ERROR
		response.Msg = err.Error()
		ret_str,_:=json.Marshal(resp)
		return string(ret_str)
	}
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1!=nil{
		response.Code = CODE_ERROR
		response.Msg = err1.Error()
		ret_str,_:=json.Marshal(resp)
		return string(ret_str)
	}
	fileName :=util.GenerateUuid()+".mp3"
	DirName :="voiceFiles"
	fileName = DirName + fileName

	errWrite := ioutil.WriteFile(fileName, body, 0777)
	if errWrite!=nil{
		response.Code = CODE_ERROR
		response.Msg = errWrite.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	//这里需要跟问题进行关联

	question:=new(model.WechatVoiceQuestions)
	questionErr :=question.GetConn().Where("uuid = ?",req.QuestionId).Find(&question).Error

	if questionErr!=nil&&!strings.Contains(questionErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = questionErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	question.VoicePath = fileName
	questionErr = question.GetConn().Save(&question).Error

	if questionErr!=nil&&!strings.Contains(questionErr.Error(),RNF){
		response.Code = CODE_ERROR
		response.Msg = questionErr.Error()
		ret_str,_:=json.Marshal(response)
		return string(ret_str)
	}

	response.Code = CODE_SUCCESS
	response.Msg = MSG_SUCCESS
	ret_str,_:=json.Marshal(response)
	return string(ret_str)
}
