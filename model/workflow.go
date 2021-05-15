package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Workflow struct {
	UUID      string `gorm:"column:uuid;type:char(32);not null;unique;primaryKey"`
	Name      string `gorm:"column:name;type:varchar(256);not null;unique"`
	Mode      int    `gorm:"column:mode"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrWorkflowExists = errors.New("workflow exists")
var ErrWorkflowNotFound = errors.New("workflow not found")

func (Workflow) TableName() string {
	return "ogm_approval_workflow"
}

type WorkflowQuery struct {
	UUID string
	Name string
}

type WorkflowDAO struct {
	conn *Conn
}

func NewWorkflowDAO(_conn *Conn) *WorkflowDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &WorkflowDAO{
		conn: conn,
	}
}

func (this *WorkflowDAO) Count() (int64, error) {
	var count int64
	err := this.conn.DB.Model(&Workflow{}).Count(&count).Error
	return count, err
}

func (this *WorkflowDAO) Insert(_workflow *Workflow) error {
	var count int64
	err := this.conn.DB.Model(&Workflow{}).Where("uuid = ? OR name = ?", _workflow.UUID, _workflow.Name).Count(&count).Error
	if nil != err {
		return err
	}

	if count > 0 {
		return ErrWorkflowExists
	}

	return this.conn.DB.Create(_workflow).Error
}

func (this *WorkflowDAO) Update(_workflow *Workflow) error {
	var count int64
	err := this.conn.DB.Model(&Workflow{}).Where("uuid = ?", _workflow.UUID).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrWorkflowNotFound
	}

	return this.conn.DB.Updates(_workflow).Error
}

func (this *WorkflowDAO) Delete(_uuid string) error {
	var count int64
	err := this.conn.DB.Model(&Workflow{}).Where("uuid = ?", _uuid).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrWorkflowNotFound
	}

	return this.conn.DB.Where("uuid = ?", _uuid).Delete(&Workflow{}).Error
}

func (this *WorkflowDAO) List(_offset int64, _count int64) ([]*Workflow, error) {
	var workflows []*Workflow
	res := this.conn.DB.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&workflows)
	return workflows, res.Error
}

func (this *WorkflowDAO) QueryOne(_query *WorkflowQuery) (*Workflow, error) {
	db := this.conn.DB.Model(&Workflow{})
	hasWhere := false
	if "" != _query.UUID {
		db = db.Where("uuid = ?", _query.UUID)
		hasWhere = true
	}
	if "" != _query.Name{
		db = db.Where("name = ?", _query.Name)
		hasWhere = true
	}
	if !hasWhere {
		return nil, ErrWorkflowNotFound
	}

	var workflow Workflow
	err := db.First(&workflow).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrWorkflowNotFound
	}
	return &workflow, err
}
