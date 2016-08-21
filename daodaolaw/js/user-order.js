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
      url:"http://www.mylvfa.com/voice/ucenter/userlist",
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
        }else if(data.code=='10003'){
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
				<OrderList changeEvaluate={this.changeEvaluate} getOrderId={this.getOrderId}/>
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
				<ResolvedList changeEvaluate={this.props.changeEvaluate} getOrderId={this.props.getOrderId}/>
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
    	this.req(this,'json/userOrder.json','-1')
    }
  },
  toAnswer:function(orderId){
  	location.href='answer.html?orderId='+orderId;
  },
	render:function(){
		var list=<p>没有相关信息</p>
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
    	this.req(this,'json/userOrder.json','-1')
    }
  },
  getAnswer:function(orderId,canEval,e){
  	this.props.getOrderId(orderId)
  	//听完语音后显示评价框
  	var audio=$(e.target).prev()[0]
  	var ms=audio.duration*1000
  	// audio.play()
  	// setTimeout(function(){
  	// 	if(audio.ended&&canEval){
	  		this.props.changeEvaluate(true)
	  // 	}
  	// }.bind(this),ms)
  },
  addOne:function(dom){
  	if(dom.addNum>0){
  		location.href="ask.html?laywerId="+dom.laywerId+'&typeId='+dom.typeId+'&orderId='+dom.orderId+'&isAdd=1';
  	}else{
  		alert('不能再追问')
  	}
  },
	render:function(){
		var list=<p>没有相关信息</p>
		var orderInfo=this.state.orderInfo
		var isAddMore=this.state.isAddMore?'点击加载更多':'没有相关信息了'
		if(orderInfo&&orderInfo.length>0){
			list=orderInfo.map(function(dom){
				return  <div className="laywer-order-list user-order">
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
									    <span className="price" onTouchEnd={this.getAnswer.bind(this,dom.orderId,dom.canEval)}>收听</span>
									    <img src="img/xiaoxi.png"/>
								    </p>
									</div>
								</div>
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
			// $.ajx({
			// 	url:'',
			// 	type:'POST',
			// 	data:JSON.stringify(data),
			// 	dataType:'json',
			// 	success:function(data){
			// 		if(data.code===10000){
			// 			//发红包 redPacket
						this.props.changeEvaluate(false)
						this.props.changeMoney(true,data.redPacket)
			// 		}else{
			// 	    this.tips(data.msg)
			//    }
			// 	}.bind(this),
			// 	error:function(data){
			// 		console.log('评价失败:',data)
			// 	}
			// })
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
	// componentDidUpdate:function(){
	// 	if(!this.props.isShowMoney){
	// 		$('.money .pic')
	// 			.animate({top:'160px'},1000)
	// 			.animate({top:"70px"},{queue:false,duration:1000})
	// 	}
	// },
	componentDidMount:function(){
		if(!this.props.isShowMoney){
			$('.money .pic')
				.animate({top:'30%'},1000)
				.animate({top:"5%"})
				.animate({top:'30%'},1000)
				.animate({top:"20%"},function(){
					$('.pic>img').attr('src','img/money_open.jpg');
					$('.price').css('display','inline');
				})
		}
	},
  render:function(){
  	var isShowMoney=!this.props.isShowMoney?'mcover':'dispN'; 
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