package router

import (
	"github.com/Unknwon/macaron"
	// c "wechatvoice/handler"
)

func SetRouter(m *macaron.Macaron) {
	m.Post("/pay/decodewechatpayinfo")
}
