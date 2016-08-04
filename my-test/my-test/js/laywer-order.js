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
var LaywerOrder=React.createClass({
	render:function(){
		return (
			<div>
				<OrderNav/>
				<OrderList/>
			</div>
		)
	}
})
var OrderNav=React.createClass({
	render:function(){
		return (
			<ul className="nav nav-tabs order-nav">
        <li className="active"><a href=".unsolved" data-toggle="tab">待收货</a></li>
        <li className=""><a href=".resolved" data-toggle="tab">已收货</a></li>
      </ul>
		)
	}
})
var OrderList=React.createClass({
	render:function(){
		return (
			<div className="tab-content">
				<UnsolvedList/>
				<ResolvedList/>
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
				return  <div className="laywer-order-list">
									<p className="over-hidden">
										<span className="pull-left">订单号: {dom.orderId}</span>
										<span className="pull-right status" onTouchEnd={this.toAnswer.bind(this,dom.orderId)}>{dom.status}</span>
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
    	this.req(this,'json/laywerOrder.json','-1')
    }
  },
	render:function(){
		var list=<p>没有相关信息</p>
		var orderInfo=this.state.orderInfo
		var isAddMore=this.state.isAddMore?'点击加载更多':'没有相关信息了'
		if(orderInfo&&orderInfo.length>0){
			list=orderInfo.map(function(dom){
				return  <div className="laywer-order-list">
									<p className="over-hidden">
										<span className="pull-left">订单号: {dom.orderId}</span>
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
			<div className="tab-pane resolved">
				{list}
				<a href="javascript:void(0)" className="wBtn-showMore clr-gray" onTouchEnd={this.addMore}>{isAddMore}</a>
			</div>
		)
	}
})
React.render(<LaywerOrder/>,document.getElementById('laywer-order'))