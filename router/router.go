package router

import (
	"github.com/Unknwon/macaron"
	c "wechatvoice/handler"
)

func SetRouter(m *macaron.Macaron) {
	m.Post("/pay/decodewechatpayinfo", c.CreateNewQuestion)
	m.Post("/back/addcate", c.CreateCateList)
	m.Post("/back/getcatelist", c.GetCateList)
	m.Post("/back/deletecate", c.DeleteCateInfo)
	m.Post("/back/editcate", c.EditCateInfo)

	m.Post("/back/questionsettinglist", c.GetQuestionSettingList)
	m.Post("/back/getsettingbyid", c.GetQuestionSettingsById)
	m.Post("/back/deletesetting", c.DeleteQuestionSettingsById)
	m.Post("/back/editsetting", c.EditWechatVoiceQuestionSettings)
	m.Post("/back/addsetting", c.AddQuestionSetting)

	m.Post("/back/getquestionlist", c.GetBadAnswerList)
	m.Post("/back/getquestionbyid", c.GetAnswerInfoById)
	m.Post("/back/editquestion", c.ReEvaluatBadAnswers)

	m.Post("/front/questionquery", c.QuestionQuery)
	m.Get("/toindex", c.ToIndex)

	m.Get("/front/getcatList", c.GetQuestionCateList)
	m.Post("/front/createquestion", c.CreateNewQuestion)
	m.Post("/front/createnewspecialquestion", c.CreateNewSpecialQuestion)
	m.Post("/front/appendquestion", c.AppendQuestion)
	m.Post("/front/peekavalable", c.PeekAvalable)
	m.Post("/front/answerquestioninit", c.AnswerQuestionInit)
	m.Post("/front/getorderdetailbyid", c.GetOrderDetailById)
	m.Post("/front/doanswerquestion", c.DoAnswerQuestion)
	m.Post("/front/ranktheanswer", c.RankTheAnswer)
	m.Post("/front/checklock", c.CheckAnswerIsLocked)
	m.Post("/front/initspecialinfo", c.InitSpecialInfo)

	m.Post("/ucenter/lawyerlist", c.GetLayerOrderList)
	m.Post("/ucenter/userlist", c.GetMemberOrderList)
	m.Post("/ucenter/orderdetail", c.GetOrderDetailById)

	m.Get("/tool/sign", c.GetSign)
	m.Get("/tool/code", c.GetOpenCodeInfo)
	m.Get("/tool/info", c.GetAllInfo)

	m.Post("/front/initpay", c.InitPay)
	m.Post("/front/dopay", c.DoPayNew)
	m.Get("/front/uni", c.UniFi)
	m.Get("/front/toindex", c.ToIndex)
}
