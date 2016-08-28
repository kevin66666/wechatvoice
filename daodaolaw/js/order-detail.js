var OrderDetail=React.createClass({
	getInitialState:function(){
		return {
			info:'',
			isShow:false,
			imgIndex:0,
			isPlay:true
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
	  	if(this.state.isPlay){
	      $audio.prop({src:answer,autoplay:'autoplay'})
	      $audio.on('play',function(){
	        timer=setInterval(function(){
	          var imgIndex=_this.state.imgIndex;
	          if(imgIndex<=1){
	            _this.setState({imgIndex:imgIndex+1})
	          }else{
	            _this.setState({imgIndex:0})
	          }
	        },100)
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
	        _this.setState({imgIndex:0})
	      })
	    }else{
	      clearInterval(timer)
	      $audio[0].pause()
	      _this.setState({imgIndex:0})
	    }
	    this.setState({isPlay:!this.state.isPlay})
	},
	play:function(answer,e){
	    $('img').prop('src','img/xiaoxi.png')
		$('audio').prop({'src':''})
		var $audio=$(e.target).prev()
		var $img=$(e.target).next()
  		$audio.prop({src:answer,autoplay:'autoplay'})
  		var timer=''
  		$audio.on('play',function(){
  			var imgIndex=0;
  			timer=setInterval(function(){
  				var src=['img/xiaoxi.png','img/half.png'][imgIndex]
  				if(imgIndex<=2){
  					$img.prop({src:src})
  					imgIndex+=1
  				}else{
  					imgIndex=0;
  				}
  			},1000)
  		})
	  	$audio.on('ended',function(){
	  		clearInterval(timer)
	  		$img.prop({src:'img/xiaoxi.png'})
	  	})
	},
	render:function(){
		var info=this.state.info;
		var isAddNum=info.addNum>0?'text-center padding-md-t':'dispN';
		var isShow=this.state.isShow?'padding-md-t add-Info':'dispN';
		var url="ask.html?laywerId="+info.laywerId+'&typeId='+info.typeId+'&orderId=-1&isAdd=0';
		// var addInfo=''
		// if(info.addInfo&&info.addInfo.length>0){
		// 	addInfo=info.addInfo.map(function(dom){
		// 		return 	<li>
		// 				  		<p>{dom.question}</p>
		// 				  		<p className="add-voice">
		// 						    <audio src={dom.answer} controls="controls" ref="record"/>
		// 						    <span className="price" onTouchEnd={this.play.bind(this,dom.answer)}>免费听取</span>
		// 						    <img src="img/xiaoxi.png"/>
		// 					    </p>
		// 				  	</li>
		// 	}.bind(this))
		// }
		var star=[]
		for(var i=0;i<info.star;i++){
			star.push(<i className="fa fa-star col-yellow"></i>)
		}
		var src=['img/xiaoxi.png','img/half.png'][this.state.imgIndex]
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
			</div>
		)
	}
})

React.render(<OrderDetail/>,document.getElementById('order-detail'))