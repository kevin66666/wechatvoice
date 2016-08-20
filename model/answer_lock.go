package model

import (
	"github.com/jinzhu/gorm"
	"wechatvoice/tool/db"
)

type AnswerLockInfo struct {
	gorm.Model
	Uuid         string
	QuestionId   string
	OpenIdFirst  string //律师ID
	OpenIdSecond string //律师ID
	IsLocked     string //是否已经被锁 0 1
	LockedTimes  string //被锁几次
	LockedTime   string //被锁时间
}

func init() {
	info := new(Category)
	info.GetConn().AutoMigrate(&AnswerLockInfo{})
}

func (this *AnswerLockInfo) GetConn() *gorm.DB {
	db := dbpool.OpenConn()
	return db.Model(&AnswerLockInfo{})
}

func (this *AnswerLockInfo) CloseConn(db *gorm.DB) {
	dbpool.CloseConn(db)
}

func GetLockListById(questionId string) ([]AnswerLockInfo, error) {
	conn := dbpool.OpenConn()
	defer conn.Close()
	list := make([]AnswerLockInfo, 0)
	err := conn.Where("question_id = ?", questionId).Find(&list).Error
	return list, err
}
