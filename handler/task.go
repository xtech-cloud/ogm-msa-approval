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
		Meta:     _req.Meta,
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

	daoAction := model.NewActionDAO(nil)
	// 写入记录
	uuid := model.ToUUID(_req.Uuid + _req.Operator)
	action := model.Action{
		UUID:     uuid,
		Task:     _req.Uuid,
		Operator: _req.Operator,
		State:    int(proto.ActionStatus_ACTION_STATUS_ACCEPTED),
	}

	err := daoAction.Upsert(&action)
	if nil != err {
		return err
	}

	err = this.updateTaskStatus(_req.Uuid)
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

	dao := model.NewActionDAO(nil)
	uuid := model.ToUUID(_req.Uuid + _req.Operator)
	action := model.Action{
		UUID:     uuid,
		Task:     _req.Uuid,
		Operator: _req.Operator,
		State:    int(proto.ActionStatus_ACTION_STATUS_REJECTED),
	}
	err := dao.Upsert(&action)
	if nil != err {
		return err
	}

	err = this.updateTaskStatus(_req.Uuid)
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
			Uuid:      task.UUID,
			Subject:   task.Subject,
			Body:      task.Body,
			State:     proto.TaskStatus(task.State),
			UpdatedAt: task.UpdatedAt.UTC().Unix(),
		}
	}
	return nil
}

func (this *Task) Search(_ctx context.Context, _req *proto.TaskSearchRequest, _rsp *proto.TaskSearchResponse) error {
	logger.Infof("Received Task.Search, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Workflow {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "workflow is required"
		return nil
	}

	offset := int64(0)
	count := int64(100)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewJoinsDAO(nil)
	query := &model.JoinsQuery{
		State:    int(_req.State),
		Subject:  _req.Subject,
		Body:     _req.Body,
		Meta:     _req.Meta,
		Workflow: _req.Workflow,
	}
	total, tasks, err := dao.SearchTask(offset, count, query)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.TaskEntity, len(tasks))
	for i, task := range tasks {
		_rsp.Entity[i] = &proto.TaskEntity{
			Uuid:      task.UUID,
			Subject:   task.Subject,
			Body:      task.Body,
			Meta:      task.Meta,
			State:     proto.TaskStatus(task.State),
			UpdatedAt: task.UpdatedAt.UTC().Unix(),
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
	task, err := dao.Get(_req.Uuid)
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
		Meta:    task.Meta,
		State:   proto.TaskStatus(task.State),
	}
	return nil
}

func (this *Task) updateTaskStatus(_task string) error {
	daoAction := model.NewActionDAO(nil)
	// 已通过的数量
	countAccepted, err := daoAction.CountWithState(_task, int(proto.ActionStatus_ACTION_STATUS_ACCEPTED))
	if nil != err {
		return err
	}
	// 已拒绝的数量
	countRejected, err := daoAction.CountWithState(_task, int(proto.ActionStatus_ACTION_STATUS_REJECTED))
	if nil != err {
		return err
	}

	// 获取任务实体
	daoTask := model.NewTaskDAO(nil)
	task, err := daoTask.Get(_task)
	if nil != err {
		return err
	}

	// 获取工作流中的操作员数量
	daoOperator := model.NewOperatorDAO(nil)
	countOperator, err := daoOperator.Count(&model.OperatorQuery{
		Workflow: task.Workflow,
	})
	if nil != err {
		return err
	}

	// 获取工作流实体
	daoWorkflow := model.NewWorkflowDAO(nil)
	workflow, err := daoWorkflow.QueryOne(&model.WorkflowQuery{
        UUID: task.Workflow,
    })
	if nil != err {
		return err
	}

	taskState := proto.TaskStatus_TASK_STATUS_PENDING
	if workflow.Mode == int(proto.WorkflowMode_WORKFLOW_MODE_ALL) {
		if countRejected > 0 {
			// 全票模式下只要有一个操作员拒绝，任务不通过
			taskState = proto.TaskStatus_TASK_STATUS_REJECTED
		} else if countAccepted == countOperator {
			// 全票模式所有操作员通过，任务通过
			taskState = proto.TaskStatus_TASK_STATUS_ACCEPTED
		} 
	} else if workflow.Mode == int(proto.WorkflowMode_WORKFLOW_MODE_ANY) {
		if countAccepted > 0 {
			// 单票模式下只要有一个操作员通过，任务通过
			taskState = proto.TaskStatus_TASK_STATUS_ACCEPTED
		} else if countRejected == countOperator {
			// 单票模式所有操作员拒绝，任务拒绝
			taskState = proto.TaskStatus_TASK_STATUS_REJECTED
		}
	} else if workflow.Mode == int(proto.WorkflowMode_WORKFLOW_MODE_MAJORITY) {
		if countAccepted > countOperator/2 {
			// 过半模式下只要有一半以上操作员通过，任务通过
			taskState = proto.TaskStatus_TASK_STATUS_ACCEPTED
		} else if countRejected > countOperator/2 {
			// 过半模式下只要有一半以上操作员，任务拒绝
			taskState = proto.TaskStatus_TASK_STATUS_REJECTED
		}
	}
    task.State = int(taskState)
    return daoTask.Update(task)
}
