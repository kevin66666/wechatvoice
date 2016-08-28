var getDataMixin={
  getInitialState:function(){
    return {
      orderInfo:[],
      isAddMore:true
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
          alert(data.msg)
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
	render:function(){	
		return (
			<div>
				<OrderNav/>
				<OrderList changeEvaluate={this.changeEvaluate} getOrderId={this.getOrderId} changeLoad={this.changeLoad}/>
				<Evaluate isShowEvaluate={this.state.isShowEvaluate} changeEvaluate={this.changeEvaluate} changeMoney={this.changeMoney} changeLoad={this.changeLoad} orderId={this.state.orderId}/>
				<Money changeMoney={this.changeMoney} isShowMoney={this.state.isShowMoney} money={this.state.money}/>
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
				<UnsolvedList />
				<ResolvedList changeEvaluate={this.props.changeEvaluate} getOrderId={this.props.getOrderId} changeLoad={this.props.changeLoad}/>
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
	render:function(){
		var list=<p className="no-info">没有相关信息</p>
		var orderInfo=this.state.orderInfo
		var isAddMore=this.state.isAddMore?'点击加载更多':'没有相关信息了'
		if(orderInfo&&orderInfo.length>0){
			list=orderInfo.map(function(dom){
				return  <div className="laywer-order-list padding-bottom-20">
									<p className="over-hidden">
										<span className="pull-left">订单号: {dom.orderId}</span>
									</p>
									<p>详情: {dom.content}</p>
									<p className="over-hidden">
										<span className="pull-left">类型: {dom.typeName}</span>
										<span className="pull-right">时间: {dom.time}</span>
									</p>
								</div>
			}.bind(this))
		}
		return (
			<div className="tab-pane active unsolved">
				{list}
				<a href="javascript:void(0)" className="wBtn-showMore clr-gray" onTouchEnd={this.addMore}>{isAddMore}</a>
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
	render:function(){
		var list=<p className="no-info">没有相关信息</p>
		var orderInfo=this.state.orderInfo
		var isAddMore=this.state.isAddMore?'点击加载更多':'没有相关信息了'
		if(orderInfo&&orderInfo.length>0){
			list=orderInfo.map(function(dom){
				return  <PerOrder dom={dom} changeLoad={this.props.changeLoad} getOrderId={this.props.getOrderId} changeEvaluate={this.props.changeEvaluate}/>
			}.bind(this))
		}
		return (
			<div className="tab-pane resolved">
				{list}
				<a href="javascript:void(0)" className="wBtn-showMore clr-gray" onTouchEnd={this.addMore}>{isAddMore}</a>
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
  		this.tips('不能再追问')
  	}
  },
  getAnswer:function(orderId,canEval,answer,e){
  	this.props.getOrderId(orderId)
  	//听完语音后显示评价框
  	var _this=this
  	var $audio=$(e.target).prev()
	$('img').prop('src','img/xiaoxi.png')
	$('audio').prop('src','')
  	$audio.prop({src:answer,autoplay:'autoplay'})
  	var timer=''
  	$audio.on('play',function(){
		alert($audio[0].duration)
  		timer=setInterval(function(){
  			var imgIndex=_this.state.imgIndex;
  			if(imgIndex<=2){
  				_this.setState({imgIndex:imgIndex+1})
  			}else{
  				_this.setState({imgIndex:0})
  			}
  		},500)
  	})
  	$audio.on('ended',function(){
  		clearInterval(timer)
  		_this.setState({imgIndex:0})
  		if(canEval){
  			_this.props.changeEvaluate(true)
  		}
  	})
  },
  tips:function(text){
		this.props.changeLoad('load',true)
    this.props.changeLoad('tips',text)
    setTimeout(function(){
      this.props.changeLoad('load',false)
    }.bind(this),2000)
	},
	render:function(){
		var dom=this.props.dom;
		var src=['img/xiaoxi.png','img/dian.png','img/half.png'][this.state.imgIndex]
		return (
			<div className="laywer-order-list user-order">
				<p className="over-hidden">
					<span className="pull-left">订单号: {dom.orderId}</span>
				</p>
				<p>详情: {dom.content}</p>
				<p className="over-hidden">
					<span className="pull-left">类型: {dom.typeName}</span>
					<span className="pull-right">时间: {dom.time}</span>
				</p>
				<div className="over-hidden padding-md-b">
					<span className="user-add-num" onTouchEnd={this.addOne.bind(this,dom)}>可追问{dom.addNum}次</span>
					<p className="voice pull-right">
				    <audio src={dom.answer} controls="controls"/>
				    <span className="price" onTouchEnd={this.getAnswer.bind(this,dom.orderId,dom.canEval,dom.answer)}>收听</span>
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
						this.props.changeMoney(true,data.redPacket)
			 		}else{
				    this.tips(data.msg)
			   }
				}.bind(this),
				error:function(data){
					console.log('评价失败:',data)
				}
			})
		}else{	
			this.tips('您还没有进行评分')
		}
	},
	tips:function(text){
		this.props.changeLoad('load',true)
    this.props.changeLoad('tips',text)
    setTimeout(function(){
      this.props.changeLoad('load',false)
    }.bind(this),2000)
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
          <p className="title">*请对本次服务做出评价<br/>评价后可领取随机红包</p>
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
				.animate({top:'30%'},800)
				.animate({top:"5%"})
				.animate({top:'30%'},800)
				.animate({top:"20%"},function(){
					$('.pic>img').attr('src','img/money_open.jpg');
					$('.price').css('display','inline');
					setTimeout(function(){
						_this.props.changeMoney(false,'')	

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