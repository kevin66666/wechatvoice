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
				<OrderList changeLoad={this.changeLoad} tips={this.tips}/>
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
				<ResolvedList tips={this.props.tips}/>
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
  				this.props.tips('问题已被锁定')
  			}
  		}.bind(this)
  	})
  },
	render:function(){
		var list=<p className="no-info">没有相关信息</p>
		var orderInfo=this.state.orderInfo
		var isAddMore=this.state.isAddMore?'点击加载更多':'没有相关信息了'
		if(orderInfo&&orderInfo.length>0){
			list=orderInfo.map(function(dom){
        var status=['抢答','解答','解答'][dom.questionType];
        var questionType=['普通咨询','定向咨询','追问咨询'][dom.questionType];
        var color=dom.questionType==0?'pull-right status clr-orange':'pull-right status color-skyblue';
				return  <div className="laywer-order-list">
									<p className="over-hidden">
										<span className="pull-left">订单号:{dom.orderId}</span>
										<span className={color} onTouchEnd={this.toAnswer.bind(this,dom.orderId)}>{status}</span>
									</p>
									<p>{dom.content}</p>
									<p className="over-hidden">
										<span className="pull-left">类型:{dom.type}&nbsp;|&nbsp;{questionType}</span>
										<span className="pull-right">{dom.time}</span>
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
      url:'http://www.mylvfa.com/voice/order/deletLawOrder',
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
				return <PerOrder dom={dom} changePlay={this.changePlay} end={this.end} delet={this.delet} tips={this.props.tips} changeLoad={this.props.changeLoad} />
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
  getAnswer:function(answer,isPlay,orderId,e){
  	//听完语音后显示评价框
    var $audio=$(e.target).prev()
    var _this=this
    $audio.on('play',function(){
      _this.setState({imgIndex:1})
    })
    $audio.on('ended',function(){
      _this.setState({imgIndex:0})
      _this.props.end(orderId)
    })
    $audio.on('pause',function(){
      _this.setState({imgIndex:0})
      _this.props.end(orderId)
    })
    if(isPlay){
      // $audio.prop({src:answer,autoplay:'autoplay'})
      $audio[0].play()
      this.props.changePlay(orderId)
    }else{
      $audio[0].pause()
    }
  },
	render:function(){
		var dom=this.props.dom;
		var style=['voice pull-right','voice-bg pull-right'][this.state.imgIndex]
    alert(style)
    var questionType=['普通咨询','定向咨询','追问咨询'][dom.questionType];
		return (
			  <div className="laywer-order-list">
					<p className="over-hidden">
						<span className="pull-left">订单号:{dom.orderId}</span>
            <span className="pull-right del-order" onTouchEnd={this.props.delet.bind(this,dom.orderId)}>删除</span>
					</p>
					<p>{dom.content}</p>
					<p className="over-hidden">
						<span className="pull-left">类型:{dom.type}&nbsp;|&nbsp;{questionType}</span>
						<span className="pull-right">{dom.time}</span>
					</p>
					<div className="over-hidden padding-md-b">
						<p className={style}>
					    <audio src={dom.answer} controls="controls"/>
					    <span className="price" onTouchEnd={this.getAnswer.bind(this,dom.answer,dom.isPlay,dom.orderId)}>查听</span>
					    <img src="img/xiaoxi.png"/>
				    </p>
			    </div>
				</div>
		)
	}
})

React.render(<LaywerOrder/>,document.getElementById('laywer-order'))