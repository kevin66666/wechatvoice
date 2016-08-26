var OrderDetail=React.createClass({
	getInitialState:function(){
		return {
			info:'',
			isShow:false,
			imgIndex:0
		}
	},
	componentDidMount:function(){
		var orderId=location.search?location.search.slice(1).split('=')[1]:'';
		var data={orderId:orderId}
		$.ajax({
			url:'http://www.mylvfa.com/voice/front/getdetailbyid',
			type:'POST',
			data:JSON.stringify(data),
			dataType:'json',
			contentType: "application/json",
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
	getAnswer:function(answer,e){
  	var _this=this
  	var imgIndex=this.state.imgIndex;
  	var $audio=$(e.target).prev()
  	$audio.prop({src:answer,autoplay:'autoplay'})
  	var timer=''
  	$audio.on('play',function(){
  		_this.setState({imgIndex:0})
  		timer=setInterval(function(){
  			if(imgIndex<=2){
  				_this.setState({imgIndex:imgIndex+1})
  			}else{
  				_this.setState({imgIndex:0})
  			}
  		},1000)
  	})
  	$audio.on('ended',function(){
  		clearInterval(timer)
  		_this.setState({imgIndex:0})
  	})
  },
	play:function(answer,e){
		var $audio=$(e.target).prev()
  	$audio.prop({src:answer,autoplay:'autoplay'})
  	var timer=''
  	var imgIndex=0;
  	$audio.on('play',function(){
  		timer=setInterval(function(){
  			var src=['img/xiaoxi.png','img/dian.png','img/half.png'][imgIndex]
  			if(imgIndex<=2){
  				$(e.target).next().prop({src:src})
  				imgIndex+=1
  			}else{
  				imgIndex=0;
  			}
  		},1000)
  	})
  	$audio.on('ended',function(){
  		clearInterval(timer)
  		_this.setState({imgIndex:0})
  	})
	},
	render:function(){
		var info=this.state.info;
		var isAddNum=info.addNum>0?'text-center padding-md-t':'dispN';
		var isShow=this.state.isShow?'padding-md-t add-Info':'dispN';
		var url="ask.html?laywerId="+info.laywerId+'&typeId='+info.typeId+'&orderId=-1&isAdd=0';
		var addInfo=''
		if(info.addInfo&&info.addInfo.length>0){
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
		var src=['img/xiaoxi.png','img/dian.png','img/half.png'][imgIndex]
		return (
			<div className="media quest-list margin-lg-t padding-vertical-md">
			  <div className="media-body">
			    <p>{info.question}</p>
			    <p className="over-hidden">
			    	<span className="pull-left">{info.typeName}&nbsp;|&nbsp;{info.name}&nbsp;|&nbsp;{info.selfIntr}</span>
			    	<span className="pull-right">{star}</span>
			    </p>
				  <p className="pull-left"><a href={url}><img src={info.pic}/></a></p>
			    <p className="voice pull-left">
				    <audio src={info.answer} controls="controls"/>
				    <span className="price" onTouchEnd={this.getAnswer.bind(this,info.answer)}>免费听取</span>
				    <img src={src}/>
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

React.render(<OrderDetail/>,document.getElementById('order-detail'))