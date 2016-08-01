package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type WechatVoiceQuestions struct {
	gorm.Model
	Uuid              string
	Category          string //分类名称
	CategoryId        string //分类名称ID
	CategoryIdInt     int64  //分类名称ID int
	Description       string //问题描述
	CreateTime        string //创建时间
	CustomerId        string //用户ID
	CustomerName      string //用户姓名
	CustomerOpenId    string //用户OPENID
	AskTime           string //提问时间
	AskerHeadImg      string //用户头像
	AnswerId          string
	AnswerName        string
	AnswerOpenId      string
	AnswerHeadImg     string
	IsAnswerd         string //是否已经做出回答
	VoicePath         string //服务器保留MP3文件路径
	AnswerdTime       string //回答时间
	Pv                int64  //问题浏览数
	IsRanked          string //0 未评价  1 已评价
	RankInfo          string //1 ~5 分
	IsSendToBack      string //0 微推送后台进行审核 1 已推送
	BackRankPoint     int64  //后台审核后进行评价的分数
	SolvedTime        string //问题解决时间
	AnswerdCount      int64  //听取数量
	IsSolved          string //是否解决 0 未解决 1 有人在解决 2 已解决
	OrderNumber       string //
	PaymentInfo       string //分记录
	PaymentInfoInt    int64  //int
	QuestionType      string //0 丢出去 1 指定人回答
	Important         string //如果是指定回答 那么这个需要在律师显示出来靠前 权重增加
	ParentQuestionId  string //追问的上一层问题ID
	AppenQuestionTime int64  //追问次数
	HaveAppendChild   string //是否有追问问题
}

func init() {
	info := new(WechatVoiceQuestions)
	info.GetConn().AutoMigrate(&WechatVoiceQuestions{})
}

func (this *WechatVoiceQuestions) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&WechatVoiceQuestions{})
}

func (this *WechatVoiceQuestions) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}

type QuestionQuery struct {
	KeyWord    string `json:"keyWord"`
	CategoryId string `json:"categoryId"`
	StartLine  int64  `json:"startLine"`
	EndLine    int64  `json:"endLine"`
}

func GetQuestionQuery(req QuestionQuery) ([]WechatVoiceQuestions, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	list1 := make([]WechatVoiceQuestions, 0)
	query := conn.Where("is_solved = 1")
	var err error
	var count int64
	if req.KeyWord != "" {
		query = query.Where("description LIKE ?", "%"+req.KeyWord+"%")
	}

	if req.CategoryId != "" {
		query = query.Where("category_id = ?", req.CategoryId)
	}
	query = query.Order("rank_info asc")
	err = query.Find(&list1).Count(&count).Error
	err = query.Offset(req.StartLine - 1).Limit(req.EndLine - req.StartLine + 1).Find(&list).Error
	return list, count, err
}

func QueryUserQuestions(userOpenId string, status []string, startLine, endLine int64) ([]WechatVoiceQuestions, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	l1 := make([]WechatVoiceQuestions, 0)
	var err error
	var count int64
	err = conn.Where("customer_open_id = ?", userOpenId).Where("is_solved in (?)", status).Find(&l1).Count(&count).Error
	err = conn.Where("customer_open_id = ?", userOpenId).Where("is_solved in (?)", status).Offset(startLine - 1).Limit(endLine - startLine + 1).Find(&list).Error
	return list, count, err
}

func QueryLawyerQuestions(userOpenId string, startLine int64, endLine int64) ([]WechatVoiceQuestions, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	l1 := make([]WechatVoiceQuestions, 0)
	var err error
	var count int64
	err = conn.Where("answer_open_id = ?", userOpenId).Where("is_solved =?", "2").Find(&l1).Count(&count).Error
	err = conn.Where("answer_open_id = ?", userOpenId).Where("is_solved =?", "2").Offset(startLine - 1).Limit(endLine - startLine + 1).Find(&list).Error
	return list, count, err
}

func QueryLawerNotSolvedQuestions(catList []string, startLine, endLine int64) ([]WechatVoiceQuestions, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	l1 := make([]WechatVoiceQuestions, 0)
	var err error
	var count int64
	err = conn.Where("category_id in (?)", catList).Where("is_solved =?", "1").Find(&l1).Count(&count).Error
	err = conn.Where("category_id in (?)", catList).Where("is_solved =?", "1").Order("important desc").Offset(startLine - 1).Limit(endLine - startLine + 1).Find(&list).Error
	return list, count, err
}

func QueryBadAnswers(statusList []string, startLine, endLine int64) ([]WechatVoiceQuestions, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	l1 := make([]WechatVoiceQuestions, 0)
	var err error
	var count int64
	err = conn.Where("category_id in (?)", catList).Where("is_solved =?", "1").Where("rank_info in (?)", statusList).Find(&l1).Count(&count).Error
	err = conn.Where("category_id in (?)", catList).Where("is_solved =?", "1").Where("rank_info in (?)", statusList).Offset(startLine - 1).Offset(endLine - startLine + 1).Find(&l1).Count(&count).Error
	return list, count, err
}
