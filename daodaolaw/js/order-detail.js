var OrderDetail=React.createClass({
	getInitialState:function(){
		return {
			info:'',
			isShow:false,
			imgIndex:0,
			imgOne:0,
			imgTwo:0,
			isPlay:true,
			isOnePlay:true,
			isTwoPlay:true
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
      },300)
    })
    $audio.on('ended',function(){
      clearInterval(timer)
      _this.setState({
      	imgIndex:0,
      	isPlay:true
      })
    })
    $audio.on('pause',function(){
      clearInterval(timer)
      _this.setState({
      	imgIndex:0,
      	isPlay:true
      })
    })
    if(this.state.isPlay){
      $audio[0].play()
    }else{
      clearInterval(timer)
      $audio[0].pause()
    }
    this.setState({isPlay:!this.state.isPlay})
	},
	play:function(answer,index,e){
	 	var $audio=$(e.target).prev()
  	var timer=''
  	var _this=this
  	$audio.on('play',function(){
      timer=setInterval(function(){
        var imgIndex=index==0?_this.state.imgOne:_this.state.imgTwo;
        if(imgIndex<=1){
        	index==0?_this.setState({imgOne:imgIndex+1}):_this.setState({imgTwo:imgIndex+1})
        }else{
        	index==0?_this.setState({imgOne:0}):_this.setState({imgTwo:0})
        }
      },300)
    })
    $audio.on('ended',function(){
      clearInterval(timer)
      if(index==0){
      	_this.setState({
	      	imgOne:0,
	      	isOnePlay:true
	      })
      }else{
      	_this.setState({
	      	imgTwo:0,
	      	isTwoPlay:true
	      })
      }
    })
    $audio.on('pause',function(){
      clearInterval(timer)
      if(index==0){
      	_this.setState({
	      	imgOne:0,
	      	isOnePlay:true
	      })
      }else{
      	_this.setState({
	      	imgTwo:0,
	      	isTwoPlay:true
	      })
      }
    })
    if(index==0){
    	if(this.state.isOnePlay){
	      $audio[0].play()
	    }else{
	      clearInterval(timer)
	      $audio[0].pause()
	    }
    }else{
    	if(this.state.isTwoPlay){
	      $audio[0].play()
	    }else{
	      clearInterval(timer)
	      $audio[0].pause()
	    }
    }
    index==0?this.setState({isOnePlay:!this.state.isOnePlay}):this.setState({isTwoPlay:!this.state.isTwoPlay})
	},
	render:function(){
		var info=this.state.info;
		var isAddNum=info.addNum>0?'text-center padding-md-t':'dispN';
		var isShow=this.state.isShow?'padding-md-t add-Info':'dispN';
		var url="ask.html?laywerId="+info.laywerId+'&typeId='+info.typeId+'&orderId=-1&isAdd=0';
		var addInfo=''
		if(info.addInfo&&info.addInfo.length>0){
			addInfo=info.addInfo.map(function(dom,index){
				var src='img/xiaoxi.png'
				if(index==0){
					src=['img/xiaoxi.png','img/half.png'][this.state.imgOne]
				}else{
					src=['img/xiaoxi.png','img/half.png'][this.state.imgTwo]
				}
				var img=index==0?imgOne
				return 	<li>
						  		<p>{dom.question}</p>
						  		<div className="over-hidden">
						  			<p className="pull-left"><a href={url}><img src={info.pic}/></a></p>
							  		<p className="add-voice pull-left">
									    <audio src={dom.answer} controls="controls" ref="record"/>
									    <span className="price" onTouchEnd={this.play.bind(this,dom.answer,index)}>点击听取</span>
									    <img src={src}/>
								    </p>
							    </div>
						  	</li>
			}.bind(this))
		}
		var star=[]
		for(var i=0;i<info.star;i++){
			star.push(<i className="fa fa-star col-yellow"></i>)
		}
		var src=['img/xiaoxi.png','img/half.png'][this.state.imgIndex]
		return (
			<div className="quest-list">
				<div className="media margin-lg-t padding-vertical-md">
				  <div className="media-body">
				    <p>{info.question}</p>
				    <p className="over-hidden">
				    	<span className="pull-left">{info.typeName}&nbsp;|&nbsp;{info.name}律师&nbsp;|&nbsp;{info.selfIntr}</span>
				    	<span className="pull-right">{star}</span>
				    </p>
					  <p className="pull-left"><a href={url}><img src={info.pic}/></a></p>
				    <p className="voice pull-left">
					    <audio src={info.answer} controls="controls"/>
					    <span className="price" onTouchEnd={this.getAnswer.bind(this,info.answer)}>点击听取</span>
					    <img src={src}/>
				    </p>
				    <p className="pull-right">{info.time}</p>
				  </div>
				  <p className={isAddNum} onTouchEnd={this.changeFold}>有{info.addNum}次追问<i className="fa fa-angle-down"></i></p>
				  <ul className={isShow}>
				  	{addInfo}
				  </ul>
				</div>
			</div>
		)
	}
})

React.render(<OrderDetail/>,document.getElementById('order-detail'))