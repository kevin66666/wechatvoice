var recorder;
var audio = document.querySelector('audio');
var Answer=React.createClass({
	getInitialState:function(){
		return {
			info:'',
			orderId:'',
			answer:'',
			typeId:'',           
			typeName:'',
			content:'',
			isRecord:true,
			load:false,
      		tips:'加载中,请稍等'
		}
	},
	componentDidMount:function(){
		var orderId=location.search?location.search.slice(1).split('=')[1]:''
		 this.setState({orderId:orderId})
		$.ajax({
			url:'http://www.mylvfa.com/voice/order/getdetail',  //回答问题页面初始化订单信息
			type:'POST',
			data:JSON.stringify({orderId:orderId}),
			contentType: "application/json",
		 	dataType:'json',
			success:function(data){
				if(data.code===10000){
					this.setState({
						info:data,
						typeName:data.typeName,
						content:data.content
					})
				}else{
					this.tips(data.msg)
				}
			}.bind(this),
			error:function(data){
				console.log('初始化信息失败:',data)
			}
		})
		//初始化时获取微信config
		$.ajax({
			url:'http://www.mylvfa.com/voice/front/getconfig',  //回答问题页面初始化订单信息
			type:'POST',
			data:JSON.stringify({orderId:orderId}),
			contentType: "application/json",
			dataType:'json',
			success:function(data){
				if(data.code===10000){  
					wx.config({
							debug: false,
							appId: data.appId,
							timestamp: data.timestamp,
							nonceStr: data.nonceStr,
							signature: data.signature,
							jsApiList: ['translateVoice','startRecord', 'stopRecord', 'onRecordEnd','playVoice','pauseVoice','stopVoice','uploadVoice','downloadVoice',] 
					})
				}else{
					this.tips(data.msg)
				}
			}.bind(this),
			error:function(data){
				console.log('获取config失败:',data)
			}
		})
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
	record:function(e){
		alert(this.state.isRecord)
		e.stopPropagation()
		if(this.state.isRecord){
			this.start()
		}else{
			this.stop()
		}
	},
	start:function(){
		var _this=this
		_this.tips('开始录音')
		this.setState({isRecord:false})
		wx.startRecord({
			cancel: function () {
				_this.tips('用户拒绝授权录音')
      }
		})
		//监听录音自动停止接口
		wx.onVoiceRecordEnd({
		    // 录音时间超过一分钟没有停止的时候会执行 complete 回调
		    complete: function (res) {
		      var localId = res.localId
		      _this.setState({answer:localId})
		      _this.tips('录音时长超过1分钟,关闭录音')
		    }
		})
		wx.error(function(res){
		  _this.tips('微信录音接口调取失败')
		})
	},
	stop:function(){
        var _this=this
		this.setState({isRecord:true})
		wx.stopRecord({
	    success: function (res) {
	      var localId = res.localId
	      _this.setState({answer:localId})
	      _this.tips('结束录音')
	    }
		})
		wx.error(function(res){
		  _this.tips('微信录音接口调取失败')
		})
	},
	reset:function(){
		this.start()
		wx.stopRecord({
	    success: function (res) {
	      var localId = res.localId
	      _this.setState({answer:localId})
	      _this.tips('结束录音')
	    }
		})
	},
	play:function(e){
		e.stopPropagation()
		alert('play')
		var $play=$(e.target)
		$play.addClass('bg-answer')
	   	var _this=this
		wx.playVoice({
			localId: _this.state.answer 
		})
		wx.onVoicePlayEnd({
		    success: function (res) {
		    	$play.removeClass('bg-answer')
		        var localId = res.localId; // 返回音频的本地ID
		    }
		})
	    wx.error(function(res){
	      $play.removeClass('bg-answer')
		  _this.tips('微信播放录音接口调取失败')
		})
	},
	save:function(e){
		//var $save=$(e.target)
		//$save.addClass('bg-answer')
		e.stopPropagation()
		this.changeLoad('load',true)
    	this.changeLoad('tips','保存中，请稍后')
		var _this=this;
		wx.uploadVoice({
		    localId: _this.state.answer, 
		    isShowProgressTips: 1, // 默认为1，显示进度提示
	      	success: function (res) {
	      	var serverId = res.serverId; // 返回音频的服务器端ID
	      		_this.doSave(serverId) 
	    	}
		})
		wx.error(function(res){
	      //$save.removeClass('bg-answer')
		  _this.tips('微信播放录音接口调取失败')
		})
	},
	doSave:function(serverId){
		var data={
			orderId:this.state.orderId,
			mediaId:serverId
		}
		$.ajax({
			url:'http://www.mylvfa.com/voice/order/uploadmedia', //保存录音--服务器端ID{serverId:"serverId"}
			type:'POST',
			data:JSON.stringify(data),
			contentType: "application/json",
			dataType:'json',
			success:function(data){
				if(data.code===10000){
					this.tips('音频保存成功')
					location.href='laywer-order.html'
				}else{
					this.tips('音频保存失败')
				}
			}.bind(this),
			error:function(data){
				console.log('初始化信息失败:',data)
			}
		})
	},
	render:function(){
		var info=this.state.info;
		return (
			<div className="answer">
				<p>类型: {this.state.typeName}</p>
				<p className="content">{info.content}</p>
				<p><img src="img/luyin.png" onTouchStart={this.record}/></p>
				<p>(点击开始录音,再次点击结束录音)</p>
				<audio src="" controls></audio>
				<div className="save">
					<span className="margin-md-r dispN" onTouchEnd={this.reset}>重新录音</span>
					<span className="margin-md-r" onTouchEnd={this.play}>播放录音</span>
					<span className="margin-md-r" onTouchEnd={this.save}>确认发送</span>
				</div>
				<Loading load={this.state.load} tips={this.state.tips}/>
			</div>
		)
	}
})
React.render(<Answer/>,document.getElementById('answer'))