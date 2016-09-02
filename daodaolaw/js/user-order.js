var getDataMixin={
  getInitialState:function(){
    return {
      orderInfo:[],
      isAddMore:true,
      orderId:'',
      confirmDispN:'dispN',
      confirm:false
    }
  },
  req:function(that,url,type){
		var data={
      startNum:that.state.orderInfo.length,
      endNum:that.state.orderInfo.length+5,
      orderType:type
    }
    $.ajax({
      url:url,
      type:'POST',
       data:JSON.stringify(data),
      dataType:'json',
      contentType: "application/json",
      success:function(data){
        if(data.code===10000){
          if(data.list.length>0){
            that.setState({
              orderInfo:that.state.orderInfo.concat(data.list),
              isAddMore:true
            })
          }else{
            that.setState({isAddMore:false})
          }
        }else if(data.code===10003){
          location.href=data.msg
        }else{
          that.props.tips(data.msg)
        }
      }.bind(that),
      error:function(data){
        console.log('获取列表信息失败:',data)
      }
    })
  }
}
var UserOrder=React.createClass({
	getInitialState:function(){
		return {
			money:'',
			isShowEvaluate:false,
			isShowMoney:false,
			orderId:'',
			load:false,
      tips:'加载中,请稍等'
		}
	},
	changeEvaluate:function(val){
		this.setState({isShowEvaluate:val})
	},
	changeMoney:function(disp,val){
		this.setState({
			isShowMoney:disp,
			money:val
		})
	},
	getOrderId:function(orderId){
		this.setState({orderId:orderId})
	},
	changeLoad:function(name,val){
    var newData={}
    newData[name]=val
    this.setState(newData)
  },
  tips:function(text){
  	this.changeLoad('load',true)
    this.changeLoad('tips',text)
    setTimeout(function(){
      this.changeLoad('load',false)
    }.bind(this),2000)
  },
	render:function(){	
		return (
			<div>
				<OrderNav/>
				<OrderList changeEvaluate={this.changeEvaluate} getOrderId={this.getOrderId} tips={this.tips}/>
				<Evaluate isShowEvaluate={this.state.isShowEvaluate} changeEvaluate={this.changeEvaluate} changeMoney={this.changeMoney} tips={this.tips} orderId={this.state.orderId}/>
				<Money changeMoney={this.changeMoney} isShowMoney={this.state.isShowMoney} money={this.state.money} tips={this.tips}/>
				<Loading load={this.state.load} tips={this.state.tips}/>
			</div>
		)
	}
})
var OrderNav=React.createClass({
	render:function(){
		return (
			<ul className="nav nav-tabs order-nav">
        <li className="active"><a href=".unsolved" data-toggle="tab">待解答</a></li>
        <li className=""><a href=".resolved" data-toggle="tab">已解答</a></li>
      </ul>
		)
	}
})
var OrderList=React.createClass({
	render:function(){
		return (
			<div className="tab-content">
				<UnsolvedList tips={this.props.tips}/>
				<ResolvedList changeEvaluate={this.props.changeEvaluate} getOrderId={this.props.getOrderId} tips={this.props.tips}/>
			</div>
		)
	}
})
var UnsolvedList=React.createClass({
	mixins:[getDataMixin],
	componentWillMount:function(){
    this.addMore()
  },
  addMore:function(){
    //点击加载更多
    if(this.state.isAddMore){
    	this.req(this,'http://www.mylvfa.com/voice/ucenter/userlist','0')
    }
  },
  toAnswer:function(orderId){
  	location.href='answer.html?orderId='+orderId;
  },
  resetList:function(orderId){
  	var index='';
  	var newInfo=this.state.orderInfo;
  	newInfo.map(function(dom,i){
  		if(dom.orderId===orderId){
  			index=i
  			return 
  		}
  	})
  	newInfo.splice(index,1)
  	this.setState({orderInfo:newInfo})
  },
  delet:function(orderId){
    this.changeDisp('confirmDispN','dispB')
    this.setState({orderId:orderId})
  },
  doDelet:function(){
  	var data={orderId:this.state.orderId}
  	$.ajax({
  		url:'http://www.mylvfa.com/voice/order/delete',
  		type:'POST',
  		data:JSON.stringify(data),
  		dataType:'json',
  		contentType: "application/json",
  		success:function(data){
  			if(data.code===10000){
  				this.resetList(this.state.orderId)
  			}else{
  				this.props.tips('删除订单失败')
  			}
  		}.bind(this),
  		error:function(err){
  			console.log('删除订单失败:',err)
  		}
  	})
  },
  changeDisp:function(name,val){
    var newState={}
    newState[name]=val
    this.setState(newState)
  },
  changeChose:function(status){
    if(status){
      this.doDelet()
    }
  },
	render:function(){
		var list=<p className="no-info">没有相关信息</p>
		var orderInfo=this.state.orderInfo
		var isAddMore=this.state.isAddMore?'点击加载更多':'没有相关信息了'
		if(orderInfo&&orderInfo.length>0){
			list=orderInfo.map(function(dom){
        var questionType=['普通咨询','定向咨询','追问咨询'][dom.questionType];
				return  <div className="laywer-order-list padding-bottom-20">
									<p className="over-hidden">
										<span className="pull-left">订单号:{dom.orderId}</span>
										<span className="pull-right del-order" onTouchEnd={this.delet.bind(this,dom.orderId)}>删除</span>
									</p>
									<p>详情: {dom.content}</p>
									<p className="over-hidden">
										<span className="pull-left">类型:{dom.typeName}&nbsp;|&nbsp;{questionType}</span>
										<span className="pull-right">{dom.time}</span>
									</p>
								</div>
			}.bind(this))
		}
		return (
			<div className="tab-pane active unsolved">
				{list}
				<a href="javascript:void(0)" className="wBtn-showMore clr-gray" onTouchEnd={this.addMore}>{isAddMore}</a>
        <Confirm changeChose={this.changeChose} changeDisp={this.changeDisp} confirmDispN={this.state.confirmDispN}/>
			</div>
		)
	}
})
var ResolvedList=React.createClass({
	mixins:[getDataMixin],
	componentWillMount:function(){
    this.addMore()
  },
  addMore:function(){
    if(this.state.isAddMore){
    	this.req(this,'http://www.mylvfa.com/voice/ucenter/userlist','2')
    }
  },
  changePlay:function(orderId){
    var newInfo=this.state.orderInfo
    newInfo.map(function(dom){
      if(orderId=dom.orderId){
        dom.isPlay=!dom.isPlay
      }
    })
    this.setState({orderInfo:newInfo})
  },
  end:function(orderId){
    var newInfo=this.state.orderInfo
    newInfo.map(function(dom){
      if(orderId=dom.orderId){
        dom.isPlay=true
      }
    })
    this.setState({orderInfo:newInfo})
  },
  resetList:function(orderId){
  	var index='';
  	var newInfo=this.state.orderInfo;
  	newInfo.map(function(dom,i){
  		if(dom.orderId===orderId){
  			index=i
  			return 
  		}
  	})
  	newInfo.splice(index,1)
  	this.setState({orderInfo:newInfo})
  },
  delet:function(orderId){
    this.changeDisp('confirmDispN','dispB')
    this.setState({orderId:orderId})
  },
  doDelet:function(){
    var data={orderId:this.state.orderId}
    $.ajax({
      url:'http://www.mylvfa.com/voice/order/delete',
      type:'POST',
      data:JSON.stringify(data),
      dataType:'json',
      contentType: "application/json",
      success:function(data){
        if(data.code===10000){
          this.resetList(this.state.orderId)
        }else{
          this.props.tips('删除订单失败')
        }
      }.bind(this),
      error:function(err){
        console.log('删除订单失败:',err)
      }
    })
  },
  changeDisp:function(name,val){
    var newState={}
    newState[name]=val
    this.setState(newState)
  },
  changeChose:function(status){
    if(status){
      this.doDelet()
    }
  },
	render:function(){
		var list=<p className="no-info">没有相关信息</p>
		var orderInfo=this.state.orderInfo
		var isAddMore=this.state.isAddMore?'点击加载更多':'没有相关信息了'
		if(orderInfo&&orderInfo.length>0){
			list=orderInfo.map(function(dom){
				return  <PerOrder dom={dom} changePlay={this.changePlay} end={this.end} delet={this.delet} tips={this.props.tips} getOrderId={this.props.getOrderId} changeEvaluate={this.props.changeEvaluate}/>
			}.bind(this))
		}
		return (
			<div className="tab-pane resolved">
				{list}
				<a href="javascript:void(0)" className="wBtn-showMore clr-gray" onTouchEnd={this.addMore}>{isAddMore}</a>
        <Confirm changeChose={this.changeChose} changeDisp={this.changeDisp} confirmDispN={this.state.confirmDispN}/>
			</div>
		)
	}
})
var PerOrder=React.createClass({
  getInitialState:function(){
		return {
			imgIndex:0
		}
	},
  addOne:function(dom){
  	if(dom.addNum>0){
  		location.href="ask.html?laywerId="+dom.laywerId+'&typeId='+dom.typeId+'&orderId='+dom.orderId+'&isAdd=1';
  	}else{
  		this.props.tips('不能再追问')
  	}
  },
  getAnswer:function(orderId,canEval,isPlay,answer,e){
  	var $audio=$(e.target).prev()
    var timer=''
    var _this=this
    this.props.getOrderId(orderId)
    $audio.on('play',function(){
      timer=setInterval(function(){
        var imgIndex=_this.state.imgIndex;
        if(imgIndex<=1){
          _this.setState({imgIndex:imgIndex+1})
        }else{
          _this.setState({imgIndex:0})
        }
      },300)
    })
    $audio.on('ended',function(){
      clearInterval(timer)
      _this.setState({imgIndex:0})
      _this.props.end(orderId)
      if(canEval){
  		_this.props.changeEvaluate(true)
  	  }
    })
    $audio.on('pause',function(){
      clearInterval(timer)
      _this.setState({imgIndex:0})
      _this.props.end(orderId)
    //   if(canEval){
  		// _this.props.changeEvaluate(true)
  	 //  }
    })
    if(isPlay){
      // $audio.prop({src:answer,autoplay:'autoplay'})
      $audio[0].play()
    }else{
      clearInterval(timer)
      $audio[0].pause()
    }
    this.props.changePlay(orderId)
  },
	render:function(){
		var dom=this.props.dom;
		var src=['img/xiaoxi.png','img/half.png'][this.state.imgIndex]
    var questionType=['普通咨询','定向咨询','追问咨询'][dom.questionType];
    var isShowAdd=dom.questionType==2?'dispN':'user-add-num'
		return (
			<div className="laywer-order-list user-order">
				<p className="over-hidden">
					<span className="pull-left">订单号: {dom.orderId}</span>
					<span className="pull-right del-order" onTouchEnd={this.props.delet.bind(this,dom.orderId)}>删除</span>
				</p>
				<p>详情: {dom.content}</p>
				<p className="over-hidden">
					<span className="pull-left">类型:{dom.typeName}&nbsp;|&nbsp;{questionType}</span>
					<span className="pull-right">{dom.time}</span>
				</p>
				<div className="over-hidden padding-md-b">
					<span className={isShowAdd} onTouchEnd={this.addOne.bind(this,dom)}>可追问{dom.addNum}次</span>
					<p className="voice pull-right">
				    <audio src={dom.answer} controls="controls"/>
				    <span className="price" onTouchEnd={this.getAnswer.bind(this,dom.orderId,dom.canEval,dom.isPlay,dom.answer)}>收听</span>
				    <img src={src}/>  
			    </p>
				</div>
			</div>
		)
	}
})
var Evaluate=React.createClass({
	getInitialState:function(){
		return {number:0}
	},
	getNum:function(index){
		this.setState({number:index})
	},
	close:function(e){
		if($(e.target).hasClass('mcover')){
			this.props.changeEvaluate(false)
		}
	},
	getEvaluate:function(val){
		var data={
			star:this.state.number,
			orderId:this.props.orderId
		}
		if(this.state.number>0){
			 $.ajax({
			 	url:'http://www.mylvfa.com/voice/order/evaluate',
			 	type:'POST',
			 	data:JSON.stringify(data),
			 	dataType:'json',
			 	success:function(data){
			 		if(data.code===10000){
			 			//发红包 redPacket 接口
						this.props.changeEvaluate(false)
            this.props.tips('感谢您的评价')
						// this.props.changeMoney(true,data.redPacket)
			 		}else{
				    this.props.tips('评价失败')
			   }
				}.bind(this),
				error:function(err){
					console.log('评价失败:',err)
				}
			})
		}else{	
			this.props.tips('您还没有进行评分')
		}
	},
  render:function(){
  	var isShowEvaluate=this.props.isShowEvaluate?'mcover':'dispN';
  	var number=this.state.number;
  	var _this=this
  	var star=[]
		for(var i=1;i<=5;i++){
			var color=(0<number&&i<=number)?'fa fa-star col-yellow':'fa fa-star-o col-yellow'
			star.push(<i className={color} onTouchEnd={_this.getNum.bind(this,i)}></i>)
		}  	
    return (
      <div className={isShowEvaluate} onTouchEnd={this.close}>
        <div className="chose-evaluate">
          <p className="title">*请对本次服务做出评价<br/><span className="dispN">评价后可领取随机红包</span></p>
          <p className="chose-btn">
          	{star}
          </p>
          <span className="save-eval" onTouchEnd={this.getEvaluate}>提交评价</span>
        </div>
      </div>
    )
  }
})
var Money=React.createClass({
	componentDidUpdate:function(){
		if(this.props.isShowMoney){
			var _this = this
			$('.money .pic')
				.animate({top:'30%'},500)
				.animate({top:"5%"})
				.animate({top:'30%'},500)
				.animate({top:"20%"},function(){
					$('.pic>img').attr('src','img/money_open.jpg');
					$('.price').css('display','inline');
					setTimeout(function(){
						_this.props.changeMoney(false,'')
						location.reload()	
					},2000)
				})
		}
	},
  render:function(){
  	var isShowMoney=this.props.isShowMoney?'mcover':'dispN'; 
  	var money=this.props.money;	
    return (
      <div className={isShowMoney}>
        <div className="money">
          <p className="pic">
          	<img src="img/money_close.jpg"/>
          	<span className="price">￥{money}</span>
          </p>
          <p className="title">*恭喜获得一个红包<br/>红包已存入您的账号余额中</p>
        </div>
      </div>
    )
  }
})
React.render(<UserOrder/>,document.getElementById('user-order'))