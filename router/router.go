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
}
