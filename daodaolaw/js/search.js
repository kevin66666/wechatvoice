var Search=React.createClass({
	getInitialState:function(){
		return {
			keywords:'',
			isShowList:false,
			isShowAsk:false,
			isAddMore:true,
			searchList:[],
			load:false,
      tips:'加载中,请稍等'
		}
	},
	componentDidMount:function(){
		var data={
			keywords:'',
			startNum:this.state.searchList.length,
			endNum:this.state.searchList.length+10
		}
		$.ajax({
			//url:'json/init.json',
			url:'/voice/front/questionquery',
			//http://www.mylvfa.com/wechatvoice/ 搜索页面初始化接口
			type:'POST',
			data:JSON.stringify(data),
			contentType: "application/json",
			dataType:'json',
			success:function(data){
				if(data.code===10000){
					if(data.list.length>0){
						this.setState({
							searchList:this.state.searchList.concat(data.list),
							isShowList:true
						})
					}else{
						this.setState({isAddMore:false})
					}
				}
			}.bind(this),
			error:function(data){
				console.log('搜索问题失败:',data)
			}
		})
	},
	getVal:function(name,val){
		var newState={}
		newState[name]=val
		this.setState(newState)
	},
	changeDisp:function(name){
		if(name==='isShowAsk'){
			this.setState({
				isShowList:false,
				isShowAsk:true
			})
		}
	},
	search:function(){
		var keywords=this.state.keywords
		this.setState({
			isShowList:true,
			isShowAsk:false
		})
		if(keywords){
			this.setState({
				searchList:[],
				isAddMore:true
			})
			this.getSearchList()
		}
	},
	getSearchList:function(){
		if(this.state.isAddMore){
			var keywords=this.state.keywords
			var searchList=this.state.searchList
			var data={
				keywords:keywords,
				startNum:searchList.length,
				endNum:searchList.length+10
			}
			$.ajax({
				url:'json/search.json',
				//关键字搜索接口
				type:'POST',
			  data:JSON.stringify(data),
				dataType:'json',
				contentType: "application/json",
				success:function(data){
					if(data.code===10000){
						if(data.list.length>0){
							this.setState({searchList:this.state.searchList.concat(data.list)})
						}else{
							this.setState({isAddMore:false})
						}
					}
				}.bind(this),
				error:function(data){
					console.log('搜索问题失败:',data)
				}
			})
		}
	},
	changeFold:function(index){
		var newList=this.state.searchList
		newList[index]['isShow']=!newList[index]['isShow']
		this.setState({searchList:newList})
	},
	resetList:function(index){
		var newList=this.state.searchList
		newList[index].isPay=true
		this.setState({searchList:newList})
	},
	changeLoad:function(name,val){
    var newData={}
    newData[name]=val
    this.setState(newData)
  },
	render:function(){
		return (
			<div>
				<SearchBar getVal={this.getVal} search={this.search} changeDisp={this.changeDisp}/>
				<SearchList isShowList={this.state.isShowList} searchList={this.state.searchList} getSearchList={this.getSearchList} changeFold={this.changeFold} isAddMore={this.state.isAddMore} resetList={this.resetList} changeLoad={this.changeLoad}/>
				<Ask isShowAsk={this.state.isShowAsk} getVal={this.getVal}/>
				<Loading load={this.state.load} tips={this.state.tips}/>
			</div>
		)
	}
})
var SearchBar=React.createClass({
	handleChange:function(e){
		this.props.getVal('keywords',e.target.value)
	},
	render:function(){
		return (
			<div className="search-bar">
				<label>
					<input type="text" placeholder="我要搜索" onChange={this.handleChange}/>
					<span className="icon" onTouchEnd={this.props.search}><i className="fa fa-search" aria-hidden="true"></i></span>
				</label>
				<span className="ask" onTouchEnd={this.props.changeDisp.bind(this,'isShowAsk')}>我要提问</span>
			</div>
		)
	}
})
var SearchList=React.createClass({
	render:function(){
		var isAddMore=this.props.isAddMore?'点击加载更多':'没有相关信息了'
		var isShow=this.props.isShowList?'quest-list':'dispN'
		var searchList=this.props.searchList
		var isShowAdd=(searchList&&searchList.length>0)?'text-center margin-lg-t padding-bottom-20':'dispN'
		var everyInfo=<p className="text-center padding-vertical-10">没有相关信息</p>
		if(searchList&&searchList.length>0){
			everyInfo=searchList.map(function(dom,index){
				return  <EverySearch info={dom} index={index} changeFold={this.props.changeFold} resetList={this.props.resetList} changeLoad={this.props.changeLoad}/>
			}.bind(this))
		}
		return (
			<div className={isShow}>
				{everyInfo}
				<p className={isShowAdd} onTouchEnd={this.props.getSearchList}>{isAddMore}</p>
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
				url:'',//支付点击
				type:'POST',
				data:JSON.stringify(data),
				contentType: "application/json",
				dataType:'json',
				success:function(data){
					if(data.code===10000){
						this.props.resetList(index)

						//调取支付接口
						wx.config({
							debug: false,
							appId: data.page_appid,
							timestamp: data.page_appid,
							nonceStr: data.page_appid,
							signature: data.page_appid,
							jsApiList: ['chooseWXPay']
						});
						wx.ready(function(){
						  wx.chooseWXPay({
							timestamp: data.pay_timeStamp,
							nonceStr: data.pay_nonceStr,
							package: data.pay_package,
							signType: data.pay_signType,
							paySign: data.pay_paySign,
							success: function (res) {
							  // 支付成功
								location.href = 'order-detail.html?orderId='+orderId
							},
							fail: function (res) {
							  // 支付失败
							  window.location.replace="pay-fail?r=0&orderId="+orderId
							},
							cancel: function (res) {
							  // 用户取消
							  window.location.replace="pay-fail?r=1&orderId="+orderId
							}
						  });
						});
						wx.error(function(res){
						  window.location.replace="pay-fail?r=2&orderId="+orderId
						});
					}
				}.bind(this),
				error:function(data){
					console.log('搜索问题失败:',data)
				}
			})
		}else{
			// 提取录音
			location.href = 'order-detail.html?orderId='+orderId
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
						  		<p>{dom.question}</p>
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
			typePrice:'',
			isShowType:false,
			allType:[]
		}
	},
	componentDidMount:function(){
 		$.ajax({
			url:'http://www.mylvfa.com/voice/front/getcatList',
			//获取所有问题类型接口
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
  		content:this.state.content
  	}
  	$.ajax({
				url:'http://www.mylvfa.com/voice/front/createquestion',
				//搜索页面提问接口
				type:'POST',
				data:JSON.stringify(data),
				dataType:'json',
				contentType: "application/json",
				success:function(data){
					if(data.code===10000){
						//调取支付
						wx.config({
							debug: false,
							appId: data.page_appid,
							timestamp: data.page_appid,
							nonceStr: data.page_appid,
							signature: data.page_appid,
							jsApiList: ['chooseWXPay']
						});
						wx.ready(function(){
						  wx.chooseWXPay({
							timestamp: data.pay_timeStamp,
							nonceStr: data.pay_nonceStr,
							package: data.pay_package,
							signType: data.pay_signType,
							paySign: data.pay_paySign,
							success: function (res) {
							  // 支付成功
								location.href = 'user-order.html'
							},
							fail: function (res) {
							  // 支付失败
							  window.location.replace="pay-fail?r=0&orderId="+orderId
							},
							cancel: function (res) {
							  // 用户取消
							  window.location.replace="pay-fail?r=1&orderId="+orderId
							}
						});
						wx.error(function(res){
						  window.location.replace="pay-fail?r=2&orderId="+orderId
						});
					}
				}.bind(this),
				error:function(data){
					console.log('提交问题失败:',data)
				}
			})
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
