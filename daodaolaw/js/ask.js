var Ask=React.createClass({
	getInitialState:function(){
		return {
			init:'',
			typeId:'',
			content:'',
			laywerId:'',
			typePrice:'2',
			//parentOrderId:'',
			// isShowType:false,
			allType:[],
			isAdd:'',
			orderId:''
		}
	},
	componentWillMount:function(){
		var params=location.search?location.search.slice(1):''
		var laywerId=params?params.split('&')[0].split('=')[1]:''
		var typeId=params?params.split('&')[1].split('=')[1]:''
		var orderId=params?params.split('&')[2].split('=')[1]:''
		var isAdd=params?params.split('&')[3].split('=')[1]:''

		var data={
			laywerId:laywerId,
			typeId:typeId,
			orderId:orderId //-1 是搜索  其他是追问orderId 
		}
		$.ajax({
			url:'http://www.mylvfa.com/voice/front/getbyid',
			type:'POST',
			data:JSON.stringify(data),
			dataType:'json',
			contentType: "application/json",
			success:function(data){
				if(data.code===10000){
					this.setState({
						init:data,
						typeId:data.typeId,
						laywerId:laywerId,
						orderId:orderId,
						typePrice:data.typePrice,
						parentOrderId:data.parentOrderId,
						isAdd:isAdd
					})
					if(isAdd==1){
						//追问免费
						this.setState({typePrice:''})
						//修改title
						$('title').html('追问')
						var $body = $('body');
			      document.title = '追问';
			      var $iframe = $('<iframe src="/favicon.ico" width="0" height="0"></iframe>').on('load', function() {
			      setTimeout(function() {
			        $iframe.off('load').remove()
			        $('iframe').remove()
			        }, 1)
			      }).appendTo($body)
					}
				}
			}.bind(this),
			error:function(data){
				console.log('获取初始化信息失败:',data)
			}
		})
	},
	// componentDidMount:function(){
 // 		$.ajax({
	// 		url:'json/allType.json',
	// 		type:'GET',
	// 		// data:JSON.stringify(data),
	// 		dataType:'json',
	// 		success:function(data){
	// 			if(data.code===10000){
	// 				this.setState({allType:data.list})
	// 			}
	// 		}.bind(this),
	// 		error:function(data){
	// 			console.log('获取类型失败:',data)
	// 		}
	// 	})
	// },
	handleChange:function(e){
		this.setState({content:e.target.value})
	},
	limitNum:function(e){
    var value=e.target.value
    if(value.length>300){
      e.preventDefault()
    }
  },
  changeType:function(){
  	this.setState({isShowType:!this.state.isShowType})
  },
  // getType:function(id,name,price){
  // 	this.setState({
  // 		isShowType:false,
  // 		typeId:id,
  // 		typeName:name,
  // 		typePrice:price
  // 	})
  // },
  doAsk:function(){
  	var data={
  		laywerId:this.state.laywerId,
  		typeId:this.state.typeId,
  		typePrice:this.state.typePrice,
  		content:this.state.content,
  		parentOrderId:this.state.parentOrderId,
  		orderId:this.state.orderId
  	}
  	$.ajax({
			url:'http://www.mylvfa.com/voice/front/createsquestion', 
			//提问接口
			type:'POST',
			data:JSON.stringify(data),
			contentType: "application/json",
			dataType:'json',
			success:function(data){
				if(data.code===10000){
					if(this.state.isAdd==0){
						WeixinJSBridge.invoke(
             	'getBrandWCPayRequest', {
                   "appId": data.appId,     //公众号名称，由商户传入
                   "timeStamp":data.timeStamp,         //时间戳，自1970年以来的秒数
                   "nonceStr":data.nonceStr, //随机串
                   "package":data.package,
                   "signType":data.signType,         //微信签名方式：
                   "paySign":data.paySign, //微信签名
              },
              function(res){
                if(res.err_msg == "get_brand_wcpay_request:ok" ) {     // 使用以上方式判断前端返回,微信团队郑重提示：res.err_msg将在用户支付成功后返回    ok，但并不保证它绝对可靠。
                  location.href = "user-order.html"
                }else if(res.err_msg == "get_brand_wcpay_request:cancel"){
                    // location.href = "pay-fail.html?r=1&orderId="+data.orderId
                }else{
                  location.href = "pay-fail.html?r=0&orderId="+data.orderId
                }
              }
            )
						wx.error(function(res){
						  window.location.replace="pay-fail.html?r=2&orderId="+data.orderId
						})
					}else{
						//追问
						location.href='user-order.html'
					}
				}
			}.bind(this),
			error:function(data){
				console.log('提交问题失败:',data)
			}
		})
  },
	render:function(){
		var init=this.state.init
		var isShowType=this.state.isShowType?'type-select':'dispN'
		var typeName=this.state.typeName?this.state.typeName:'选择类型'
		var typePrice=this.state.typePrice?<span>￥{this.state.typePrice}元</span>:'免费咨询'
		var allType=this.state.allType;
		var list=''
		// if(allType.length>0){
		// 	list=allType.map(function(dom){
		// 		return <li onTouchEnd={this.getType.bind(this,dom.typeId,dom.typeName,dom.typePrice)}>{dom.typeName}</li>
		// 	}.bind(this))
		// }
		return (
			<div className="question">
				<p className="laywer-info">
					<img src={init.pic}/><br/>
					{init.name}&nbsp;&nbsp;|&nbsp;&nbsp;{init.selfIntr}
				</p>
				<div className="content"><textarea rows="8" placeholder="最多300个字" onChange={this.handleChange} onKeyPress={this.limitNum}></textarea></div>
				<p className="price">{typePrice}</p>
				<div className="btn-ask"><p onTouchEnd={this.doAsk}>写好了</p></div>
			</div>
		)
	}
})
React.render(<Ask/>,document.getElementById('ask'))