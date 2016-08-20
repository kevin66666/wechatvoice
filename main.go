package main

import (
	//"shangqu-finance/router"
	"github.com/Unknwon/macaron"
	"wechatvoice/router"
)

func main() {
	m := macaron.New()
	m.SetDefaultCookieSecret("git.yanzhilu.org")
	router.SetRouter(m)
	m.Run(8000)
}
