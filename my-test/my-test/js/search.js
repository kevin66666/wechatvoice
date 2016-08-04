var Search=React.createClass({
	getInitialState:function(){
		return {
			keywords:'',
			isShowList:false,
			isShowAsk:false,
			isAddMore:true,
			searchList:[]
		}
	},
	getVal:function(name,val){
		var newState={}
		newState[name]=val
		this.setState(newState)
	},
	changeDisp:function(name){
		if(name==='isShowAsk'){
			this.setState({
				isShowList:false,
				isShowAsk:true
			})
		}
	},
	search:function(){
		var keywords=this.state.keywords
		this.setState({
			isShowList:true,
			isShowAsk:false
		})
		if(keywords){
			this.setState({
				searchList:[],
				isAddMore:true
			})
			this.getSearchList()
		}
	},
	getSearchList:function(){
		if(this.state.isAddMore){
			var keywords=this.state.keywords
			var searchList=this.state.searchList
			var data={
				keywords:keywords,
				startNum:searchList.length,
				endNum:searchList.length+10
			}
			$.ajax({
				url:'json/search.json',
				type:'GET',
				// data:JSON.stringify(data),
				dataType:'json',
				success:function(data){
					if(data.code===10000){
						if(data.list.length>0){
							this.setState({searchList:this.state.searchList.concat(data.list)})
						}else{
							this.setState({isAddMore:false})
						}
					}
				}.bind(this),
				error:function(data){
					console.log('搜索问题失败:',data)
				}
			})
		}
	},
	render:function(){
		return (
			<div>
				<SearchBar getVal={this.getVal} search={this.search} changeDisp={this.changeDisp}/>
				<SearchList isShowList={this.state.isShowList} searchList={this.state.searchList} getSearchList={this.getSearchList} isAddMore={this.state.isAddMore}/>
				<Ask isShowAsk={this.state.isShowAsk} getVal={this.getVal}/>
			</div>
		)
	}
})
var SearchBar=React.createClass({
	handleChange:function(e){
		this.props.getVal('keywords',e.target.value)
	},
	render:function(){
		return (
			<div className="search-bar">
				<label>
					<input type="text" placeholder="我要搜索" onChange={this.handleChange}/>
					<span className="icon" onTouchEnd={this.props.search}><i className="fa fa-search" aria-hidden="true"></i></span>
				</label>
				<span className="ask" onTouchEnd={this.props.changeDisp.bind(this,'isShowAsk')}>我要提问</span>
			</div>
		)
	}
})
var SearchList=React.createClass({
	render:function(){
		var isAddMore=this.props.isAddMore?'点击加载更多':'没有相关信息了'
		var isShow=this.props.isShowList?'quest-list':'dispN'
		var searchList=this.props.searchList
		var everyInfo=<p className="text-center">没有相关信息</p>
		if(searchList&&searchList.length>0){
			everyInfo=searchList.map(function(dom,index){
				var url="ask.html?laywerId="+dom.laywerId;
				return  <div className="media">
								  <div className="media-left">{index+1}.</div>
								  <div className="media-body">
								    <p>{dom.question}</p>
								    <p>{dom.name}&nbsp;&nbsp;|&nbsp;&nbsp;{dom.selfIntr}</p>
									  <p className="pull-left"><a href={url}><img src={dom.pic}/></a></p>
								    <p className="voice pull-left">
									    <audio src={dom.answer} controls="controls"/>
									    <span className="price">1元偷偷听</span>
									    <img src="img/xiaoxi.png"/>
								    </p>
								  </div>
								</div>
			}.bind(this))
		}
		return (
			<div className={isShow}>
				{everyInfo}
				<p className="text-center margin-lg-t padding-bottom-20" onTouchEnd={this.props.getSearchList}>点击加载更多</p>
			</div>
		)
	}
})
var Ask=React.createClass({
	getInitialState:function(){
		return {
			type:'',
			content:'',
			text:'',
			isShowType:false
		}
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
		var isShow=this.props.isShowAsk?'question':'dispN'
		var isShowType=this.state.isShowType?'':'dispN'
		var text=this.state.text?this.state.text:'选择提问类型'
		return (
			<div className={isShow}>
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
				<p className="price">￥1元</p>
				<div className="btn-ask"><p onTouchEnd={this.doAsk}>写好了</p></div>
			</div>
		)
	}
})
React.render(<Search/>,document.getElementById('search'))