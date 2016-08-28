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
var LaywerOrder=React.createClass({
	getInitialState:function(){
    return {
      load:false,
      tips:'加载中,请稍等'
    }
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
				<OrderList changeLoad={this.changeLoad}/>
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
				<UnsolvedList changeLoad={this.props.changeLoad}/>
				<ResolvedList changeLoad={this.props.changeLoad}/>
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
    	this.req(this,'http://www.mylvfa.com/voice/ucenter/lawyerlist','0')
    }
  },
  toAnswer:function(orderId){
  	var data={orderId:orderId}
  	$.ajax({
  		url:'http://www.mylvfa.com/voice/order/getorder', //抢答接口
  		type:'POST',
  		data:JSON.stringify(data),
  		dataType:'json',
  		contentType: "application/json",
  		success:function(data){
  			if(data.code===10000){
  				location.href='answer.html?orderId='+orderId;
  			}else if(data.code===10001){
  				this.tips('问题已被锁定')
  			}
  		}.bind(this)
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
		var list=<p className="no-info">没有相关信息</p>
		var orderInfo=this.state.orderInfo
		var isAddMore=this.state.isAddMore?'点击加载更多':'没有相关信息了'
		if(orderInfo&&orderInfo.length>0){
			list=orderInfo.map(function(dom){
				return  <div className="laywer-order-list">
									<p className="over-hidden">
										<span className="pull-left">订单号: {dom.orderId}</span>
										<span className="pull-right status" onTouchEnd={this.toAnswer.bind(this,dom.orderId)}>抢答</span>
									</p>
									<p>{dom.content}</p>
									<p className="over-hidden">
										<span className="pull-left">类型: {dom.type}</span>
										<span className="pull-right">时间: {dom.time}</span>
									</p>
									<p className="text-right">￥{dom.price}</p>
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
    //点击加载更多
    if(this.state.isAddMore){
    	this.req(this,'http://www.mylvfa.com/voice/ucenter/lawyerlist','2')
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
	render:function(){
		var list=<p className="no-info">没有相关信息</p>
		var orderInfo=this.state.orderInfo
		var isAddMore=this.state.isAddMore?'点击加载更多':'没有相关信息了'
		if(orderInfo&&orderInfo.length>0){
			list=orderInfo.map(function(dom){
				return <PerOrder dom={dom} changePlay={this.changePlay} end={this.end} changeLoad={this.props.changeLoad} getOrderId={this.props.getOrderId} changeEvaluate={this.props.changeEvaluate}/>
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
  getAnswer:function(answer,isPlay,orderId,e){
  	//听完语音后显示评价框
    var $audio=$(e.target).prev()
    var timer=''
    var _this=this
    $audio.on('play',function(){
      timer=setInterval(function(){
        var imgIndex=_this.state.imgIndex;
        if(imgIndex<=1){
          _this.setState({imgIndex:imgIndex+1})
        }else{
          _this.setState({imgIndex:0})
        }
      },500)
    })
    $audio.on('ended',function(){
      clearInterval(timer)
      _this.setState({imgIndex:0})
      _this.props.end(orderId)
    })
    $audio.on('pause',function(){
      clearInterval(timer)
      _this.setState({imgIndex:0})
    })
    if(isPlay){
      $audio.prop({src:answer,autoplay:'autoplay'})
    }else{
      clearInterval(timer)
      $audio[0].pause()
    }
    this.props.changePlay(orderId)
  },
	render:function(){
		var dom=this.props.dom;
		var src=['img/xiaoxi.png','img/half.png'][this.state.imgIndex]
		return (
			  <div className="laywer-order-list">
					<p className="over-hidden">
						<span className="pull-left">订单号: {dom.orderId}</span>
					</p>
					<p>{dom.content}</p>
					<p className="over-hidden">
						<span className="pull-left">类型: {dom.type}</span>
						<span className="pull-right">时间: {dom.time}</span>
					</p>
					<div className="over-hidden padding-md-b">
						<p className="voice pull-right">
					    <audio src={dom.answer} controls="controls"/>
					    <span className="price" onTouchEnd={this.getAnswer.bind(this,dom.answer,dom.isPlay,dom.orderId)}>查听</span>
					    <img src={src}/>
				    </p>
			    </div>
				</div>
		)
	}
})

var Per
React.render(<LaywerOrder/>,document.getElementById('laywer-order'))