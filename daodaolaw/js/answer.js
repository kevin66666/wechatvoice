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
			load:false,
      tips:'加载中,请稍等',
      config:{}
		}
	},
	componentDidMount:function(){
		var orderId=location.search?location.search.slice(1).split('=')[1]:''
		// this.setState({orderId:orderId})
		// $.ajax({
		// 	url:'http://www.mylvfa.com/voice/',  //回答问题页面初始化订单信息
		// 	type:'POST',
		// 	data:JSON.stringify(orderId),
		// 	contentType: "application/json",
		// 	success:function(data){
		// 		if(data.code===10000){
		// 			this.setState({
		// 				info:data,
		// 				typeId:data.typeId,
		// 				typeName:data.typeName
		// 			})
		// 		}else{
		// 			this.tips(data.msg)
		// 		}
		// 	}.bind(this),
		// 	error:function(data){
		// 		console.log('初始化信息失败:',data)
		// 	}
		// })
		//初始化时获取微信config
		$.ajax({
			url:'http://www.mylvfa.com/voice/front/getconfig',  //回答问题页面初始化订单信息
			type:'GET',
			// data:JSON.stringify(data),
			// contentType: "application/json",
			success:function(data){
				if(data.code===10000){  
					this.setState({
						config:config
					})
				//返回的config,该代码只做提示功能
				data={
						debug: true,
						appId: data.appId,
						timestamp: data.timestamp,
						nonceStr: data.nonceStr,
						signature: data.signature,
						jsApiList: ['translateVoice','startRecord', 'stopRecord',  'onRecordEnd',
        'playVoice',
        'pauseVoice',
        'stopVoice',
        'uploadVoice',
        'downloadVoice',] 
						//这个是固定的api
				}
				}else{
					this.tips(data.msg)
				}
			}.bind(this),
			error:function(data){
				console.log('初始化信息失败:',data)
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
	start:function(){
    // HZRecorder.get(function (rec) {
    //     recorder = rec;
    //     recorder.start();
    // })
		//开始录音
		var _this=this
		wx.config(this.state.config)
		wx.ready(function(){
			wx.startRecord()
			//监听录音自动停止接口
			wx.onVoiceRecordEnd({
			    // 录音时间超过一分钟没有停止的时候会执行 complete 回调
			    complete: function (res) {
			      var localId = res.localId
			      _this.setState({answer:localId})
			      _this.tips('录音时长超过1分钟,关闭录音')
			    }
			})
		})
		wx.error(function(res){
		  _this.tips('微信录音接口调取失败')
		})
	},
	stop:function(){
    // recorder.stop()
    //停止录音
    var _this=this
    wx.config(this.state.config)
		wx.ready(function(){
			wx.stopRecord({
		    success: function (res) {
		      var localId = res.localId
		      _this.setState({answer:localId})
		      _this.tips('结束录音')
		    }
			})
		})
		wx.error(function(res){
		  _this.tips('微信录音接口调取失败')
		})
	},
	reset:function(){
		var _this=this
		wx.config(this.state.config)
		wx.ready(function(){
			wx.startRecord()
			//监听录音自动停止接口
			wx.onVoiceRecordEnd({
			    // 录音时间超过一分钟没有停止的时候会执行 complete 回调
			    complete: function (res) {
			      var localId = res.localId
			      _this.setState({answer:localId})
			      _this.tips('录音时长超过1分钟,关闭录音')
			    }
			})
			wx.stopRecord({
		    success: function (res) {
		      var localId = res.localId
		      _this.setState({answer:localId})
		      _this.tips('结束录音')
		    }
			})
		})
		wx.error(function(res){
		  _this.tips('微信录音接口调取失败')
		})
	},
	play:function(){
		// var audio = document.querySelector('audio');
   //   recorder.play(audio)
   	var _this=this
		wx.config(this.state.config)
    wx.ready(function(){
			wx.playVoice({
			  localId: _this.state.answer 
		  })
		})
	  wx.error(function(res){
		  _this.tips('微信播放录音接口调取失败')
		})
	},
	save:function(){
		var data={};
		var _this=this;

		wx.config(this.state.config)
    wx.ready(function(){
			wx.uploadVoice({
		    localId: _this.state.answer, 
		    isShowProgressTips: 1, // 默认为1，显示进度提示
	      success: function (res) {
	      	data.serverId = res.serverId; // 返回音频的服务器端ID
	      	_this.doSave(data)
	    	}
			})
		})
	},
	doSave:function(data){
		$.ajax({
			url:'', //保存录音--服务器端ID{serverId:"serverId"}
			type:'POST',
			data:JSON.stringify(data),
			contentType: "application/json",
			success:function(data){
				if(data.code===10000){
					this.tips('音频保存成功')
				}else{
					this.tips(data.msg)
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
				<p><img src="img/luyin.png" onTouchStart={this.start} onTouchEnd={this.stop}/></p>
				<p>(按住开始回答)</p>
				<audio controls autoplay></audio>
				<div className="save">
					<span className="margin-md-r" onTouchEnd={this.reset}>重新录音</span>
					<span onTouchEnd={this.play}>播放录音</span>
					<p onTouchEnd={this.save}>确认发送</p>
				</div>
				<Loading load={this.state.load} tips={this.state.tips}/>
			</div>
		)
	}
})
React.render(<Answer/>,document.getElementById('answer'))