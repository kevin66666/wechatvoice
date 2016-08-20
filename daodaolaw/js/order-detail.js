var OrderDetail=React.createClass({
	getInitialState:function(){
		return {
			info:'',
			isShow:false
		}
	},
	componentDidMount:function(){
		var orderId=location.search?location.search.slice(1).split('=')[1]:'';
		var data={orderId:orderId}
		$.ajax({
			url:'http://localhost:8000/ucenter/orderdetail',
			type:'POST',
			 data:JSON.stringify(data),
			dataType:'json',
			success:function(data){
				if(data.code===10000){
					this.setState({info:data})
				}
			}.bind(this),
			error:function(data){
				console.log('搜索问题失败:',data)
			}
		})
	},
	changeFold:function(){
		this.setState({isShow:!this.state.isShow})
	},
	play:function(e){
		$(e.target).prev()[0].play()
	},
	render:function(){
		var info=this.state.info;
		var isAddNum=info.addNum>0?'text-center padding-md-t':'dispN';
		var isShow=this.state.isShow?'padding-md-t add-Info':'dispN';
		var url="ask.html?laywerId="+info.laywerId+'&typeId='+info.typeId+'&orderId=-1&isAdd=0';
		var addInfo=''
		if(info.addInfo.length>0){
			addInfo=info.addInfo.map(function(dom){
				return 	<li>
						  		<p>{dom.question}</p>
						  		<p className="add-voice">
								    <audio src={dom.answer} controls="controls" ref="record"/>
								    <span className="price" onTouchEnd={this.play}>免费听取</span>
								    <img src="img/xiaoxi.png"/>
							    </p>
						  	</li>
			}.bind(this))
		}
		var star=[]
		for(var i=0;i<info.star;i++){
			star.push(<i className="fa fa-star col-yellow"></i>)
		}
		return (
			<div className="media">
			  <div className="media-body">
			    <p>{info.question}</p>
			    <p className="over-hidden">
			    	<span className="pull-left">{info.typeName}&nbsp;|&nbsp;{info.name}&nbsp;|&nbsp;{info.selfIntr}</span>
			    	<span className="pull-right">{star}</span>
			    </p>
				  <p className="pull-left"><a href={url}><img src={info.pic}/></a></p>
			    <p className="voice pull-left">
				    <audio src={info.answer} controls="controls"/>
				    <span className="price" onTouchEnd={this.play}>免费听取</span>
				    <img src="img/xiaoxi.png"/>
			    </p>
			  </div>
			  <p className={isAddNum} onTouchEnd={this.changeFold}>有{info.addNum}次追问<i className="fa fa-angle-down"></i></p>
			  <ul className={isShow}>
			  	{addInfo}
			  </ul>
			</div>
		)
	}
})

var EverySearch=React.createClass({
	changeFold:function(){
		if(this.props.info.isPay){
			this.props.changeFold(this.props.index)
		}else{
			this.tips('只有付款后才能点开哦！')
		}
	},
	pay:function(info,index){
		var data={orderId:info.orderId}
		//调取支付接口
		if(!info.isPay){
			$.ajax({
				url:'',
				type:'POST',
				data:JSON.stringify(data),
				dataType:'json',
				success:function(prepayInfo){
					if(prepayInfo.code===10000){
						this.props.resetList(index)
						//调取支付接口
						WeixinJSBridge.invoke(
              'getBrandWCPayRequest', {
                   "appId": prepayInfo.body.appId,     //公众号名称，由商户传入
                   "timeStamp":prepayInfo.body.timestamp.toString(),         //时间戳，自1970年以来的秒数
                   "nonceStr":prepayInfo.body.nonceStr, //随机串
                   "package":prepayInfo.body.package,
                   "signType":prepayInfo.body.signType,         //微信签名方式：
                   "paySign":prepayInfo.body.paySign, //微信签名
               },
               function(res){
                if(res.err_msg == "get_brand_wcpay_request:ok" ) {     // 使用以上方式判断前端返回,微信团队郑重提示：res.err_msg将在用户支付成功后返回    ok，但并不保证它绝对可靠。
                  location.href = 'order-detail.html?orderId='+orderId
                }else{
                  alert('支付失败')
                }
               }
            )
					}
				}.bind(this),
				error:function(data){
					console.log('搜索问题失败:',data)
				}
			})
		}else{
			// 提取录音
			// var audio=this.refs['record'].getDOMNode()
			// audio.play()
		}
	},
	tips:function(text){
		this.props.changeLoad('load',true)
    this.props.changeLoad('tips',text)
    setTimeout(function(){
      this.props.changeLoad('load',false)
    }.bind(this),2000)
	},
	toDetail:function(orderId){
		location.href="order-detail.html?orderId="+orderId
	},
	render:function(){
		var info=this.props.info;
		var index=this.props.index;
		var url="ask.html?laywerId="+info.laywerId+'&typeId='+info.typeId+'&orderId=-1&isAdd=0';
		var isAddNum=info.addNum>0?'text-center padding-md-t':'dispN';
		var isShow=info.isShow?'padding-md-t add-Info':'dispN';
		var addInfo=''
		if(info.addInfo.length>0){
			addInfo=info.addInfo.map(function(dom){
				return 	<li>
						  		<p>这是追问</p>
						  		<p className="add-voice">
								    <audio src={dom.answer} controls="controls" ref="record"/>
								    <span className="price" onTouchEnd={this.toDetail.bind(this,dom.orderId)}>免费听取</span>
								    <img src="img/xiaoxi.png"/>
							    </p>
						  	</li>
			}.bind(this))
		}
		var star=[]
		for(var i=0;i<info.star;i++){
			star.push(<i className="fa fa-star col-yellow"></i>)
		}
		return (
			<div className="media">
			  <div className="media-left">{index+1}.</div>
			  <div className="media-body">
			    <p>{info.question}</p>
			    <p className="over-hidden">
			    	<span className="pull-left">{info.typeName}&nbsp;|&nbsp;{info.name}&nbsp;|&nbsp;{info.selfIntr}</span>
			    	<span className="pull-right">{star}</span>
			    </p>
				  <p className="pull-left"><a href={url}><img src={info.pic}/></a></p>
			    <p className="voice pull-left">
				    <audio src={info.answer} controls="controls"/>
				    <span className="price" onTouchEnd={this.pay.bind(this,info,index)}>{info.typePrice}元听取</span>
				    <img src="img/xiaoxi.png"/>
			    </p>
			  </div>
			  <p className={isAddNum} onTouchEnd={this.changeFold}>有{info.addNum}次追问<i className="fa fa-angle-down"></i></p>
			  <ul className={isShow}>
			  	{addInfo}
			  </ul>
			</div>
		)
	}
})
var Ask=React.createClass({
	getInitialState:function(){
		return {
			typeId:'',
			content:'',
			typeName:'',
			typePrice:1,
			isShowType:false,
			allType:[]
		}
	},
	componentDidMount:function(){
 		$.ajax({
			url:'json/allType.json',
			type:'GET',
			// data:JSON.stringify(data),
			dataType:'json',
			success:function(data){
				if(data.code===10000){
					this.setState({allType:data.list})
				}
			}.bind(this),
			error:function(data){
				console.log('获取类型失败:',data)
			}
		})
	},
	handleChange:function(e){
		this.setState({content:e.target.value})
	},
	limitNum:function(e){
    var value=e.target.value
    if(value.length>100){
      e.preventDefault()
    }
  },
  changeType:function(){
  	this.setState({isShowType:!this.state.isShowType})
  },
  getType:function(id,name,price){
  	this.setState({
  		isShowType:false,
  		typeId:id,
  		typeName:name,
  		typePrice:price
  	})
  },
  doAsk:function(){
  	var data={
  		typeId:this.state.typeId,
  		typePrice:this.state.typePrice,
  		content:this.state.content,
  		type:'0' //0 表示直接提问 1 指定律师提问
  	}
  	// $.ajax({
			// 	url:'json/search.json',
			// 	type:'POST',
			// 	data:JSON.stringify(data),
			// 	dataType:'json',
			// 	success:function(data){
			// 		if(data.code===10000){
			// 			调微信支付
			// 		}
			// 	}.bind(this),
			// 	error:function(data){
			// 		console.log('提交问题失败:',data)
			// 	}
			// })
  },
	render:function(){
		var isShow=this.props.isShowAsk?'question':'dispN'
		var isShowType=this.state.isShowType?'type-select':'dispN'
		var typeName=this.state.typeName?this.state.typeName:'选择类型'
		var allType=this.state.allType;
		var list=''
		if(allType.length>0){
			list=allType.map(function(dom){
				return <li onTouchEnd={this.getType.bind(this,dom.typeId,dom.typeName,dom.typePrice)}>{dom.typeName}</li>
			}.bind(this))
		}
		return (
			<div className={isShow}>
				<div className="type">
					<span onTouchEnd={this.changeType}>{typeName}<i className="fa fa-caret-down"></i></span>
					<ul className={isShowType}>
						{list}
					</ul>
				</div>
				<div className="content"><textarea rows="8" placeholder="最多100个字" onChange={this.handleChange} onKeyPress={this.limitNum}></textarea></div>
				<p className="price">￥{this.state.typePrice}元</p>
				<div className="btn-ask"><p onTouchEnd={this.doAsk}>写好了</p></div>
			</div>
		)
	}
})
React.render(<Search/>,document.getElementById('search'))