package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Operator struct {
	UUID      string `gorm:"column:uuid;type:char(32);not null;unique;primaryKey"`
	Name      string `gorm:"column:name;type:varchar(256);not null"`
	Workflow  string `gorm:"column:workflow;type:char(32);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrOperatorExists = errors.New("operator exists")
var ErrOperatorNotFound = errors.New("operator not found")

func (Operator) TableName() string {
	return "ogm_approval_operator"
}

type OperatorQuery struct {
	UUID     string
	Workflow string
    Name string
}

type OperatorDAO struct {
	conn *Conn
}

func NewOperatorDAO(_conn *Conn) *OperatorDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &OperatorDAO{
		conn: conn,
	}
}

func (this *OperatorDAO) Count(_query *OperatorQuery) (int64, error) {
	var count int64
	db := this.conn.DB.Model(&Operator{})
    if "" != _query.Workflow {
        db = db.Where("workflow = ?", _query.Workflow)
    }
	err := db.Count(&count).Error
	return count, err
}

func (this *OperatorDAO) Insert(_operator *Operator) error {
	var count int64
	err := this.conn.DB.Model(&Operator{}).Where("uuid = ?", _operator.UUID).Count(&count).Error
	if nil != err {
		return err
	}

	if count > 0 {
		return ErrOperatorExists
	}

	return this.conn.DB.Create(_operator).Error
}

func (this *OperatorDAO) Update(_operator *Operator) error {
	var count int64
	err := this.conn.DB.Model(&Operator{}).Where("uuid = ?", _operator.UUID).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrOperatorNotFound
	}

	return this.conn.DB.Updates(_operator).Error
}

func (this *OperatorDAO) Delete(_uuid string) error {
	var count int64
	err := this.conn.DB.Model(&Operator{}).Where("uuid = ?", _uuid).Count(&count).Error
	if nil != err {
		return err
	}

	if 0 == count {
		return ErrOperatorNotFound
	}

	return this.conn.DB.Where("uuid = ?", _uuid).Delete(&Operator{}).Error
}

func (this *OperatorDAO) List(_offset int64, _count int64, _query *OperatorQuery) (_total int64, _operator []*Operator, _err error) {
	_err = nil
	_total = int64(0)
	_operator = make([]*Operator, 0)

	db := this.conn.DB.Model(&Operator{})
	if "" != _query.Workflow {
		db = db.Where("workflow = ?", _query.Workflow)
	}
	if "" != _query.Name{
		db = db.Where("name = ?", _query.Name)
	}

	_err = db.Count(&_total).Error
	if nil != _err {
		return
	}
	_err = db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&_operator).Error
	return
}

func (this *OperatorDAO) QueryOne(_query *OperatorQuery) (*Operator, error) {
	db := this.conn.DB.Model(&Operator{})
	hasWhere := false
	if "" != _query.UUID {
		db = db.Where("uuid = ?", _query.UUID)
		hasWhere = true
	}
	if !hasWhere {
		return nil, ErrOperatorNotFound
	}

	var Operator Operator
	err := db.First(&Operator).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOperatorNotFound
	}
	return &Operator, err
}
