var Confirm=React.createClass({
  chose:function(status){
    this.props.changeChose(status)
    this.props.changeDisp('confirmDispN','dispN')
  },
  render:function(){
    var isDisp="mcover "+ this.props.confirmDispN;
    return (
      <div className={isDisp}>
        <div className="bgcolor-white chose-comfirm">
          <p className="title">确定删除该地址</p>
          <span>确定地址删除将无法恢复</span>
          <p className="chose-btn">
            <span className="pull-left" onTouchEnd={this.chose.bind(this,false)}>取消</span><span className="pull-right" onTouchEnd={this.chose.bind(this,true)}>确定</span>
          </p>
        </div>
      </div>
    )
  }
})