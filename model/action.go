package model

import (
	"errors"
	"gorm.io/gorm/clause"
	"time"
)

type Action struct {
	UUID      string `gorm:"column:uuid;type:char(32);not null;unique;primaryKey"`
	Task      string `gorm:"column:task;type:char(32);not null"`
	Operator  string `gorm:"column:operator;type:char(32);not null"`
	State     int    `gorm:"column:state"`
	Reason    string `gorm:"column:reason;type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Action) TableName() string {
	return "msa_approval_action"
}

type ActionQuery struct {
	Workflow string
	Task     string
	Operator string
	State    int
}

var ErrActionExists = errors.New("action exists")
var ErrActionNotFound = errors.New("action not found")

type ActionDAO struct {
	conn *Conn
}

func NewActionDAO(_conn *Conn) *ActionDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &ActionDAO{
		conn: conn,
	}
}

func (this *ActionDAO) Count() (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Action{}).Count(&count).Error
	return count, err
}

func (this *ActionDAO) Upsert(_action *Action) error {
	db := this.conn.DB.Model(&Action{})
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.AssignmentColumns([]string{"state", "reason"}),
	}).Create(_action).Error
	return err
}

func (this *ActionDAO) Delete(_uuid string) error {
	var count int64
	err := this.conn.DB.Model(&Action{}).Where("uuid = ?", _uuid).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrActionNotFound
	}

	return this.conn.DB.Where("uuid = ?", _uuid).Delete(&Action{}).Error
}

func (this *ActionDAO) List(_offset int64, _count int64, _query *ActionQuery) (_total int64, _action []*Action, _err error) {
	_err = nil
	_total = int64(0)
	_action = make([]*Action, 0)

	db := this.conn.DB.Model(&Action{})
	if "" != _query.Task {
		db = db.Where("task = ?", _query.Task)
	}
	if "" != _query.Operator {
		db = db.Where("operator = ?", _query.Operator)
	}
	if 0 != _query.State {
		db = db.Where("state = ?", _query.State)
	}

	_err = db.Count(&_total).Error
	if nil != _err {
		return
	}
	_err = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&_action).Error
	return
}

func (this *ActionDAO) CountWithState(_task string, _state int) (_count int64, _err error) {
	_err = nil
	_count = int64(0)

	db := this.conn.DB.Model(&Action{})
	db = db.Where("task = ? AND state = ?", _task, _state)
	_err = db.Count(&_count).Error
	return
}
