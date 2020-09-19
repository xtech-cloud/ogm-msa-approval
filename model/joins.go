package model

type JoinsDAO struct {
	conn *Conn
}

func NewJoinsDAO(_conn *Conn) *JoinsDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &JoinsDAO{
		conn: conn,
	}
}

type JoinsQuery struct {
    Subject string
    Body string
    Meta string
    Workflow string
    State int
}


func (this *JoinsDAO) SearchTask(_offset int64, _count int64, _query *JoinsQuery) (_total int64, _task []*Task, _err error) {
    _err = nil
    _total = int64(0)
    _task = make([]*Task, 0)

	db := this.conn.DB
    db = db.Joins("JOIN msa_approval_workflow ON msa_approval_workflow.uuid = msa_approval_task.workflow")
	db = db.Where("msa_approval_workflow.name LIKE ?", "%" + _query.Workflow + "%")
	if "" != _query.Subject{
        db = db.Where("msa_approval_task.subject LIKE ?", "%" + _query.Subject+ "%")
	}
	if "" != _query.Body{
        db = db.Where("msa_approval_task.body LIKE ?", "%" + _query.Body + "%")
	}
	if "" != _query.Meta{
        db = db.Where("msa_approval_task.meta LIKE ?", "%" + _query.Meta + "%")
	}
	if 0 != _query.State{
		db = db.Where("msa_approval_task.state = ?", _query.State)
	}
    db = db.Model(&Task{})

	_err = db.Count(&_total).Error
    if nil != _err {
        return
    }
	_err = db.Offset(int(_offset)).Limit(int(_count)).Order("msa_approval_task.created_at desc").Find(&_task).Error
	return
}
