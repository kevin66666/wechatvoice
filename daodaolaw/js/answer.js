var recorder;
var audio = document.querySelector('audio');

var Answer=React.createClass({
	getInitialState:function(){
		return {
			orderId:'',
			answer:'',
			typeId:'',           
			typeName:'婚姻类型'
		}
	},
	componentDidMount:function(){
		var orderId=location.search?location.search.slice(1).split('=')[1]:''
		this.setState({orderId:orderId})
		// $.ajax({
		// 	url:'',
		// 	type:'POST',
		// 	data:JSON.stringify(data),
		// 	success:function(data){
		// 		if(data.code===10000){
		// 			this.setState({
		// 				typeId:data.typeId,
		// 				typeName:data.typeName
		// 			})
		// 		}
		// 	}.bind(this),
		// 	error:function(data){
		// 		console.log('初始化信息失败:',data)
		// 	}
		// })
	},
	start:function(){
    HZRecorder.get(function (rec) {
        recorder = rec;
        recorder.start();
    })
	},
	stop:function(){
    recorder.stop()
	},
	reset:function(){

	},
	play:function(){
		var audio = document.querySelector('audio');
    recorder.play(audio)

	},
	save:function(){

	},
	render:function(){
		return (
			<div className="answer">
				<p>类型: {this.state.typeName}</p>
				<p className="content">这是一个什么测试测试这是一个什么测试测试这是一个什么测试测试这是一个什么测试测试测试测试这是一个什么测试测试</p>
				<p><img src="img/luyin.png" onTouchStart={this.start} onTouchEnd={this.stop}/></p>
				<p>(按住开始回答)</p>
				<audio controls autoplay></audio>
				<div className="save">
					<span className="margin-md-r" onTouchEnd={this.stop}>停止录音</span>
					<span onTouchEnd={this.play}>播放录音</span>
					<p onTouchEnd={this.save}>确认发送</p>
				</div>
			</div>
		)
	}
})
React.render(<Answer/>,document.getElementById('answer'))