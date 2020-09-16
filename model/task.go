package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	UUID      string `gorm:"column:uuid;type:char(32);not null;unique;primaryKey"`
	Subject   string `gorm:"column:subject;type:varchar(512);not null"`
	Body      string `gorm:"column:body;type:text;not null"`
	Workflow  string `gorm:"column:workflow;type:varchar(256);not null"`
	State     int    `gorm:"column:state"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrTaskExists = errors.New("task exists")
var ErrTaskNotFound = errors.New("task not found")

func (Task) TableName() string {
	return "msa_approval_task"
}

type TaskQuery struct {
	UUID string
    Subject string
    Body string
    Workflow string
    State int
}

type TaskDAO struct {
	conn *Conn
}

func NewTaskDAO(_conn *Conn) *TaskDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &TaskDAO{
		conn: conn,
	}
}

func (this *TaskDAO) Count() (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Task{}).Count(&count).Error
	return count, err
}

func (this *TaskDAO) Insert(_task *Task) error {
	var count int64
	err := this.conn.DB.Model(&Task{}).Where("uuid = ?", _task.UUID).Count(&count).Error
	if nil != err {
		return err
	}

	if count > 0 {
		return ErrTaskExists
	}

	return this.conn.DB.Create(_task).Error
}

func (this *TaskDAO) Update(_task *Task) error {
	var count int64
	err := this.conn.DB.Model(&Task{}).Where("uuid = ?", _task.UUID).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrTaskNotFound
	}

	return this.conn.DB.Updates(_task).Error
}

func (this *TaskDAO) Delete(_uuid string) error {
	var count int64
	err := this.conn.DB.Model(&Task{}).Where("uuid = ?", _uuid).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrTaskNotFound
	}

	return this.conn.DB.Where("uuid = ?", _uuid).Delete(&Task{}).Error
}

func (this *TaskDAO) List(_offset int64, _count int64, _query *TaskQuery) (_total int64, _task []*Task, _err error) {
    _err = nil
    _total = int64(0)
    _task = make([]*Task, 0)

	db := this.conn.DB.Model(&Task{})
	if "" != _query.Workflow{
		db = db.Where("workflow = ?", _query.Workflow)
	}
	if 0 != _query.State{
		db = db.Where("state = ?", _query.State)
	}

	_err = db.Count(&_total).Error
    if nil != _err {
        return
    }
	_err = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&_task).Error
	return
}

func (this *TaskDAO) QueryOne(_query *TaskQuery) (*Task, error) {
	db := this.conn.DB.Model(&Task{})
	hasWhere := false
	if "" != _query.UUID {
		db = db.Where("uuid = ?", _query.UUID)
		hasWhere = true
	}
	if !hasWhere {
		return nil, ErrTaskNotFound
	}

	var Task Task
	err := db.First(&Task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}
	return &Task, err
}
