package handler

import (
	"context"
	"errors"
	"omo-msa-approval/model"
	"omo-msa-approval/publisher"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-approval/proto/approval"
)

type Task struct{}

func (this *Task) Submit(_ctx context.Context, _req *proto.TaskSubmitRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Task.Submit, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Subject {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "subject is required"
		return nil
	}

	if "" == _req.Body {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "body is required"
		return nil
	}

	if "" == _req.Workflow {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "workflow is required"
		return nil
	}

	daoWorkflow := model.NewWorkflowDAO(nil)
	_, err := daoWorkflow.QueryOne(&model.WorkflowQuery{
		UUID: _req.Workflow,
	})
	if nil != err {
		if errors.Is(err, model.ErrWorkflowNotFound) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

	uuid := model.ToUUID(_req.Subject + _req.Body + _req.Workflow)

	Task := &model.Task{
		UUID:     uuid,
		Subject:  _req.Subject,
		Body:     _req.Body,
		Workflow: _req.Workflow,
		State:    int(proto.TaskStatus_TASK_STATUS_PENDING),
	}

	dao := model.NewTaskDAO(nil)
	err = dao.Insert(Task)
	if nil != err {
		if errors.Is(err, model.ErrTaskExists) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

	// 发布消息
	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "/task/submit", _req, _rsp)
	return nil
}

func (this *Task) Accept(_ctx context.Context, _req *proto.TaskAcceptRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Task.Accept, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	if "" == _req.Operator {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "operator is required"
		return nil
	}

	dao := model.NewTaskDAO(nil)
	query := model.TaskQuery{
		UUID: _req.Uuid,
	}
	task, err := dao.QueryOne(&query)
	if nil != err {
		if errors.Is(err, model.ErrTaskNotFound) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

	task.State = int(proto.TaskStatus_TASK_STATUS_ACCEPTED)
	err = dao.Update(task)
	if nil != err {
		return err
	}

	// 发布消息
	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "/task/accept", _req, _rsp)
	return nil
}

func (this *Task) Reject(_ctx context.Context, _req *proto.TaskRejectRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Task.Reject, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	if "" == _req.Operator {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "operator is required"
		return nil
	}

	dao := model.NewTaskDAO(nil)
	query := model.TaskQuery{
		UUID: _req.Uuid,
	}
	task, err := dao.QueryOne(&query)
	if nil != err {
		if errors.Is(err, model.ErrTaskNotFound) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

	task.State = int(proto.TaskStatus_TASK_STATUS_REJECTED)
	err = dao.Update(task)
	if nil != err {
		return err
	}

	// 发布消息
	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "/task/reject", _req, _rsp)
	return nil
}

func (this *Task) List(_ctx context.Context, _req *proto.TaskListRequest, _rsp *proto.TaskListResponse) error {
	logger.Infof("Received Task.List, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(100)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewTaskDAO(nil)
	query := &model.TaskQuery{
		State:    int(_req.State),
		Workflow: _req.Workflow,
	}
	total, tasks, err := dao.List(offset, count, query)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.TaskEntity, len(tasks))
	for i, task := range tasks {
		_rsp.Entity[i] = &proto.TaskEntity{
			Uuid:    task.UUID,
			Subject: task.Subject,
			Body:    task.Body,
			State:   proto.TaskStatus(task.State),
		}
	}
	return nil
}

func (this *Task) Get(_ctx context.Context, _req *proto.TaskGetRequest, _rsp *proto.TaskGetResponse) error {
	logger.Infof("Received Task.Get, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewTaskDAO(nil)
	query := model.TaskQuery{
		UUID: _req.Uuid,
	}
	task, err := dao.QueryOne(&query)
	if nil != err {
		if errors.Is(err, model.ErrTaskNotFound) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}
	_rsp.Entity = &proto.TaskEntity{
		Uuid:    task.UUID,
		Subject: task.Subject,
		Body:    task.Body,
		State:   proto.TaskStatus(task.State),
	}
	return nil
}
