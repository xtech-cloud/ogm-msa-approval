package handler

import (
	"context"
	"errors"
	"omo-msa-approval/config"
	"omo-msa-approval/model"
	"omo-msa-approval/publisher"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-approval/proto/approval"
)

type Workflow struct{}

func (this *Workflow) Make(_ctx context.Context, _req *proto.WorkflowMakeRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Workflow.Make, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	if proto.WorkflowMode_WORKFLOW_MODE_INVALID == _req.Mode {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "mode is required"
		return nil
	}

	// 本地数据库使用存储桶名生成UUID，方便测试和开发
	uuid := model.NewUUID()
	if config.Schema.Database.Lite {
		uuid = model.ToUUID(_req.Name)
	}

	workflow := &model.Workflow{
		UUID: uuid,
		Name: _req.Name,
		Mode: int(_req.Mode),
	}

	dao := model.NewWorkflowDAO(nil)
	err := dao.Insert(workflow)
	if nil != err {
		if errors.Is(err, model.ErrWorkflowExists) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

	// 发布消息
	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "/workflow/make", _req, _rsp)
	return nil
}

func (this *Workflow) List(_ctx context.Context, _req *proto.WorkflowListRequest, _rsp *proto.WorkflowListResponse) error {
	logger.Infof("Received Workflow.List, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(100)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewWorkflowDAO(nil)

	total, err := dao.Count()
	if nil != err {
		return nil
	}
	workflows, err := dao.List(offset, count)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.WorkflowEntity, len(workflows))
	for i, workflow := range workflows {
		_rsp.Entity[i] = &proto.WorkflowEntity{
			Uuid: workflow.UUID,
			Name: workflow.Name,
			Mode: proto.WorkflowMode(workflow.Mode),
		}
	}
	return nil
}

func (this *Workflow) Remove(_ctx context.Context, _req *proto.WorkflowRemoveRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Workflow.Remove, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewWorkflowDAO(nil)
	err := dao.Delete(_req.Uuid)
	if nil != err {
		if errors.Is(err, model.ErrWorkflowNotFound) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}
	// 发布消息
	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "/Workflow/remove", _req, _rsp)
	return nil
}

func (this *Workflow) Get(_ctx context.Context, _req *proto.WorkflowGetRequest, _rsp *proto.WorkflowGetResponse) error {
	logger.Infof("Received Workflow.Get, req is %v", _req)
	_rsp.Status = &proto.Status{}

	dao := model.NewWorkflowDAO(nil)
	query := model.WorkflowQuery{
		UUID: _req.Uuid,
        Name: _req.Name,
	}

	workflow, err := dao.QueryOne(&query)
	if nil != err {
		if errors.Is(err, model.ErrWorkflowNotFound) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}
	_rsp.Entity = &proto.WorkflowEntity{
		Uuid: workflow.UUID,
		Name: workflow.Name,
		Mode: proto.WorkflowMode(workflow.Mode),
	}
	return nil
}
