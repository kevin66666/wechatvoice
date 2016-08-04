var Answer=React.createClass({
	getInitialState:function(){
		return {
			orderId:'',
			answer:''
		}
	},
	componentDidMount:function(){
		var orderId=location.search?location.search.slice(1).split('=')[1]:''
		this.setState({orderId:orderId})
	},
	render:function(){
		return (
			<div className="answer">
				<p className="content">这是一个什么测试测试这是一个什么测试测试这是一个什么测试测试这是一个什么测试测试测试测试这是一个什么测试测试</p>
				<p><img src="img/luyin.png"/></p>
				<p>(按住开始回答)</p>
				<div className="save">
					<span>重新开始</span>
					<span>确认发送</span>
				</div>
			</div>
		)
	}
})
React.render(<Answer/>,document.getElementById('answer'))