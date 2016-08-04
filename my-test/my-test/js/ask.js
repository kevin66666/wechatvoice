var Ask=React.createClass({
	getInitialState:function(){
		return {
			init:'',
			type:'',
			content:'',
			text:'',
			isShowType:false
		}
	},
	componentDidMount:function(){
		var laywerId=location.search?location.search.slice(1).split('=')[1]:''
		var data={laywerId:laywerId}
		$.ajax({
			url:'json/ask.json',
			type:'GET',
			// data:JSON.stringify(data),
			dataType:'json',
			success:function(data){
				if(data.code===10000){
					this.setState({init:data})
				}
			}.bind(this),
			error:function(data){
				console.log('获取初始化信息失败:',data)
			}
		})
	},
	handleChange:function(e){
		this.setState({content:e.target.value})
	},
	limitNum:function(e){
    var value=e.target.value
    if(value.length>100){
      e.preventDefault()
    }
  },
  changeType:function(){
  	this.setState({isShowType:!this.state.isShowType})
  },
  getType:function(val){
  	this.setState({
  		isShowType:false,
  		type:val,
  		text:'各种类型'
  	})
  },
  doAsk:function(){
  	var data={
  		type:this.state.type,
  		content:this.state.content
  	}
  	$.ajax({
				url:'json/search.json',
				type:'POST',
				data:JSON.stringify(data),
				dataType:'json',
				success:function(data){
					if(data.code===10000){
						this.setState({searchList:data.list})
					}
				}.bind(this),
				error:function(data){
					console.log('提交问题失败:',data)
				}
			})
  },
	render:function(){
		var init=this.state.init
		var isShowType=this.state.isShowType?'':'dispN'
		var text=this.state.text?this.state.text:'选择提问类型'
		return (
			<div className="question">
				<p className="laywer-info">
					<img src={init.pic}/><br/>
					{init.name}&nbsp;&nbsp;|&nbsp;&nbsp;{init.selfIntr}
				</p>
				<div className="type">
					<span onTouchEnd={this.changeType}>{text}</span>
					<ul className={isShowType}>
						<li onTouchEnd={this.getType.bind(this,'1')}>各种类型</li>
						<li onTouchEnd={this.getType.bind(this,'2')}>各种类型</li>
						<li onTouchEnd={this.getType.bind(this,'3')}>各种类型</li>
						<li onTouchEnd={this.getType.bind(this,'4')}>各种类型</li>
					</ul>
				</div>
				<div className="content"><textarea rows="8" placeholder="最多100个字" onChange={this.handleChange} onKeyPress={this.limitNum}></textarea></div>
				<p className="price">￥{init.price}元</p>
				<div className="btn-ask"><p onTouchEnd={this.doAsk}>写好了</p></div>
			</div>
		)
	}
})
React.render(<Ask/>,document.getElementById('ask'))