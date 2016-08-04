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
      type:'GET',
      // data:JSON.stringify(data),
      dataType:'json',
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
			isShowEvaluate:false,
			orderId:''
		}
	},
	changeEvaluate:function(val){
		this.setState({isShowEvaluate:val})
	},
	getOrderId:function(orderId){
		this.setState({orderId:orderId})
	},
	render:function(){
		return (
			<div>
				<OrderNav/>
				<OrderList changeEvaluate={this.changeEvaluate} getOrderId={this.getOrderId}/>
				<Evaluate isShowEvaluate={this.state.isShowEvaluate} changeEvaluate={this.changeEvaluate} orderId={this.state.orderId}/>
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
    	this.req(this,'json/laywerOrder.json','-1')
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
    	this.req(this,'json/laywerOrder.json','-1')
    }
  },
  getAnswer:function(orderId){
  	this.props.changeEvaluate(true)
  	this.props.getOrderId(orderId)
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
									<p className="voice pull-right">
								    <audio src={dom.answer} controls="controls"/>
								    <span className="price">1元偷偷听</span>
								    <img src="img/xiaoxi.png"/>
							    </p>
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
	getEvaluate:function(val){
		this.props.changeEvaluate(false)
		var data={
			evaluate:val,
			orderId:''
		}
		// $.ajx({
		// 	url:'',
		// 	type:'POST',
		// 	data:JSON.stringify(data),
		// 	dataType:'json',
		// 	success:function(data){
		// 		if(data.code===10000){

		// 		}
		// 	}.bind(this),
		// 	error:function(data){
		// 		console.log('评价失败:',data)
		// 	}
		// })
	},
  render:function(){
  	var isShowEvaluate=this.props.isShowEvaluate?'mcover':'dispN';
    return (
      <div className={isShowEvaluate}>
        <div className="chose-evaluate">
          <p className="title">请对这次服务做出评价</p>
          <p className="chose-btn">
            <span className="pull-left" onTouchEnd={this.getEvaluate.bind(this,1)}>满意</span>
            <span className="pull-right" onTouchEnd={this.getEvaluate.bind(this,0)}>不满意</span>
          </p>
        </div>
      </div>
    )
  }
})
React.render(<UserOrder/>,document.getElementById('user-order'))