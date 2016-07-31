package model
import(
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type MsMerchantWechatInfo struct {
	gorm.Model
	Uuid              string // 唯一主键
	MerchantId        string // 商户ID
	Wid               string // 老微播系统商户id号
	Appid             string // 商户APPID
	NickName          string // 商户昵称
	HeadImg           string // 商户头像地址
	ServiceTypeInfo   string // 授权公众号类型
	VerifyTypeInfo    string // 认证类型
	UserName          string // 商户原始ID
	Alias             string // 授权方设置的微信号,可能为空
	BusinessInfo      string // 功能开通情况
	QrcodeUrl         string // 二维码图片连接
	AuthorizationInfo string // 授权信息
	GuideLink         string // 引导关注 连接
	Mchid             string // 商户MCHID
	PayKey            string // 商户payKey
	CertPath          string // 商户证书路径
	KeyPath           string // 商户密钥路径
	RootCaPath        string // 商户CA证书路径
	PlatformName      string // 商户托管第三方平台名称
	AppSecret         string //商户秘钥
	Token             string // 商户Token
	EncodingAesKey    string // 商户EncodingAESKey
}
func init() {
	// info := new(MsMerchantWechatInfo)
	// info.GetConn().AutoMigrate(&MsMerchantWechatInfo{})
}

func (this *MsMerchantWechatInfo) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&MsMerchantWechatInfo{})
}

func (this *MsMerchantWechatInfo) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}
