var Loading=React.createClass({
	render:function(){
		var isHide=this.props.load?'mcover':'dispN'
		return (
			<div className={isHide}>
				<div className="load-tips">{this.props.tips}</div>
			</div>
		)
	}
})