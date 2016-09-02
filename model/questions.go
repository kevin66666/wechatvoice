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
	IsLocked          string
	LockTime          int64
	LockedOpenId      string
	AnswerName        string
	AnswerOpenId      string
	NeedId            string
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
	IsPaied           string
	UserDelete        string
	LawyerDelete      string
	QType             string //1 追加 2 指定
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
	KeyWord    string `json:"keyWords"`
	CategoryId string `json:"categoryId"`
	StartLine  int64  `json:"startNum"`
	EndLine    int64  `json:"endNum"`
}

func GetQuestionQuery(req QuestionQuery) ([]WechatVoiceQuestions, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	list1 := make([]WechatVoiceQuestions, 0)
	query := conn.Where("is_solved = 2")
	var err error
	var count int64
	if req.KeyWord != "" {
		query = query.Where("description LIKE ?", "%"+req.KeyWord+"%")
	}

	if req.CategoryId != "" {
		query = query.Where("category_id = ?", req.CategoryId)
	}
	query = query.Order("id desc")
	err = query.Find(&list1).Count(&count).Error
	err = query.Offset(req.StartLine).Limit(req.EndLine - req.StartLine).Find(&list).Error
	return list, count, err
}
func GetQuestionQueryNew(req QuestionQuery) ([]WechatVoiceQuestions, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	list1 := make([]WechatVoiceQuestions, 0)
	query := conn.Where("is_solved in (?)", []string{"2", "3"})
	var err error
	var count int64
	if req.KeyWord != "" {
		query = query.Where("description LIKE ?", "%"+req.KeyWord+"%")
	}

	if req.CategoryId != "" {
		query = query.Where("category_id = ?", req.CategoryId)
	}
	// query = query.Order("id desc")

	err = query.Not("q_type", "1").Find(&list1).Count(&count).Error
	err = query.Not("q_type", "1").Order("solved_time desc").Offset(req.StartLine - 1).Limit(req.EndLine - req.StartLine + 1).Find(&list).Error

	// err = query.Offset(req.StartLine).Limit(req.EndLine - req.StartLine).Find(&list).Error
	return list, count, err
}

func GetQueryList(startLine, endLine int64) ([]WechatVoiceQuestions, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	var count int64
	list := make([]WechatVoiceQuestions, 0)
	list1 := make([]WechatVoiceQuestions, 0)

	err := conn.Where("is_solved = 1").Order("id desc").Offset(startLine).Limit(endLine - startLine).Find(&list).Error
	err = conn.Where("is_solved = 1").Order("id desc").Find(&list1).Count(&count).Error

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
	err = conn.Where("customer_open_id = ?", userOpenId).Where("is_solved in (?)", status).Offset(startLine).Limit(endLine - startLine).Find(&list).Error
	return list, count, err
}

func QueryLawyerQuestions(startLine int64, endLine int64, userOpenId string) ([]WechatVoiceQuestions, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	l1 := make([]WechatVoiceQuestions, 0)
	var err error
	var count int64
	err = conn.Where("answer_open_id = ?", userOpenId).Where("is_paied = 1").Where("is_solved =?", "2").Order("solved_time desc").Find(&l1).Count(&count).Error
	err = conn.Where("answer_open_id = ?", userOpenId).Where("is_paied = 1").Where("is_solved =?", "2").Order("solved_time desc").Offset(startLine).Limit(endLine - startLine).Find(&list).Error
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
	err = conn.Where("category_id in (?)", catList).Where("is_solved =?", "1").Order("important desc").Offset(startLine).Limit(endLine - startLine).Find(&list).Error
	return list, count, err
}

func QueryBadAnswers(statusList []string, startLine, endLine int64) ([]WechatVoiceQuestions, int64, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	l1 := make([]WechatVoiceQuestions, 0)
	var err error
	var count int64
	err = conn.Where("category_id in (?)", statusList).Where("is_solved =?", "1").Where("rank_info in (?)", statusList).Find(&l1).Count(&count).Error
	err = conn.Where("category_id in (?)", statusList).Where("is_solved =?", "1").Where("rank_info in (?)", statusList).Offset(startLine).Offset(endLine - startLine).Find(&l1).Count(&count).Error
	return list, count, err
}

func GetChildAnsers(questionId string) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	err := conn.Where("parent_question_id = ?", questionId).Find(&list).Error
	return list, err
}
func GetLawyerQs(status, openId string, startLine, endLien int64) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	err := conn.Order("created_at desc").Where("is_solved = ?", status).Where("is_paied = 1").Where("lawyer_delete is not 1").Where("need_id = ?", openId).Offset(startLine).Limit(endLien - startLine).Find(&list).Find(&list).Error
	return list, err
}

func GetNotSpectial(cateId, status string, start, end int64) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	str := make([]string, 0)
	str = append(str, cateId)
	err := conn.Where("category_id  not in (?)", str).Where("is_solved = ?", status).Where("is_locked = 0").Where("lawyer_delete is not 1").Offset(start).Limit(end - start).Find(&list).Error
	return list, err
}
func GetCustomerInfo(openId, status string, start, end int64) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	err := conn.Where("customer_open_id = ?", openId).Where("is_paied = 1").Not("user_delete", "1").Where("is_solved = ?", status).Order("id desc").Offset(start).Limit(end - start).Find(&list).Error
	return list, err
}
func GetCustomerInfoNew(openId, status string, idList []string, start, end int64) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	var err error
	if len(idList) == 0 {
		err = conn.Where("customer_open_id = ?", openId).Where("is_paied = 1").Where("is_solved = ?", status).Order("id desc").Offset(start).Limit(end - start).Find(&list).Error
	} else {
		err = conn.Where("customer_open_id = ?", openId).Where("is_paied = 1").Where("uuid not in (?)", idList).Where("is_solved = ?", status).Order("id desc").Offset(start).Limit(end - start).Find(&list).Error
	}
	return list, err
}
func GetAllLocked() ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	err := conn.Where("is_locked = ?", "1").Find(&list).Error
	return list, err
}
func GetInfos(openId, parentId string) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	err := conn.Where("customer_open_id = ?", openId).Where("parent_question_id =?", parentId).Find(&list).Error
	return list, err
}

func GetCustomerPaiedInfo(openid string, orderIdList, deleteList []string, startLine, endLine int64) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	var err error
	if len(deleteList) > 0 {
		err = conn.Where("customer_open_id = ?", openid).Where("uuid in (?)", orderIdList).Not("uuid", deleteList).Where("is_solved = 2").Order("id desc").Offset(startLine).Limit(endLine - startLine).Find(&list).Error

	} else {
		err = conn.Where("customer_open_id = ?", openid).Where("uuid in (?)", orderIdList).Where("is_solved = 2").Order("id desc").Offset(startLine).Limit(endLine - startLine).Find(&list).Error
	}
	return list, err
}

func GetLawerDirectInfo(openId string, startLine, endLine int64) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	err := conn.Where("answer_open_id = ?", openId).Where("is_locked = 0").Where("is_solved = 0").Order("id desc").Offset(startLine).Limit(endLine - startLine).Find(&list).Error
	return list, err
}

func GetLaerOther(startLine, endLine int64) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	err := conn.Where("is_solved = 0").Where("answer_open_id = ?", "1").Order("id desc").Offset(startLine).Limit(endLine - startLine).Find(&list).Error
	return list, err
}
func GetLockedInfo(startLine, endLine int64, openId string) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	err := conn.Where("is_solved = 0").Where("locked_open_id = ?", openId).Where("is_locked =1").Order("id desc").Offset(startLine).Limit(endLine - startLine).Find(&list).Error
	return list, err
}
func GetAppendInfo(questionId string) ([]WechatVoiceQuestions, error) {
	conn := dbpool.OpenConn()
	defer dbpool.CloseConn(&conn)
	list := make([]WechatVoiceQuestions, 0)
	err := conn.Where("parent_question_id = ?", questionId).Order("id desc").Find(&list).Error
	return list, err
}
