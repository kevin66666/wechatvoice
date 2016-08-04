/*

1.search.html
        初始化
        res: init.json----推送最新问题
        req:
            var data={
                    keywords:'', //不需要可删
                    startNum:this.state.searchList.length,
                    endNum:this.state.searchList.length+10
                }
        搜索问题
            res: search.json
            req: 
                var data={
                    keywords:'keywords',
                    startNum:0,
                    endNum:10
                }

        提问类型
            res: allType.json
            req: GET

        提交问题
            var data={
            typeId:this.state.typeId, //类型
            typePrice:this.state.typePrice, //类型价格
            content:this.state.content  //提问内容
        }

2.指定律师提问 ask.html
        初始化
            res: ask.json
            req:
                var data={
                    laywerId:laywerId,
                    typeId:typeId, //类型
                    orderId:orderId  //orderId为-1时表示这个是search.html的问题列表id ，否则为用户追问的orderId
                }
            提交问题
            req:
                var data={
                laywerId:this.state.laywerId,
                typeId:this.state.typeId,
                typePrice:this.state.typePrice,
                content:this.state.content
            }
3.回答问题answer.html
    初始化
        res: answer.json
        req:
            var data={orderId:orderId}
    提交回答
        req:
        var data={
            orderId:this.state.orderId,
            serverId:serverId //音频的服务器端ID
        }
4.律师订单
    未解答/已解答
        res: laywerOrder.json
        req:
            var data={
          startNum:that.state.orderInfo.length,
          endNum:that.state.orderInfo.length+5,
          orderType:type  //type为0未解答，1已解答
        }
5.用户订单
    未解答/已解答
        res: userOrder.json( answer为回答的录音url,canEval为是否可以评价，值都为true)
        req:
            var data={
          startNum:that.state.orderInfo.length,
          endNum:that.state.orderInfo.length+5,
          orderType:type  //type为0未解答，1已解答
        }
            






Xcd123456      114.55.9.219
18616182109



*/ 