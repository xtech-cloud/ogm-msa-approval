package handler

import (
	"context"
	"errors"
	"omo-msa-approval/model"
	"omo-msa-approval/publisher"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-approval/proto/approval"
)

type Operator struct{}

func (this *Operator) Join(_ctx context.Context, _req *proto.OperatorJoinRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Operator.Join, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Operator {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "operator is required"
		return nil
	}

	if "" == _req.Workflow {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "workflow is required"
		return nil
	}

	uuid := model.ToUUID(_req.Operator + _req.Workflow)

	operator := &model.Operator{
		UUID:     uuid,
		Name:     _req.Operator,
		Workflow: _req.Workflow,
	}

	dao := model.NewOperatorDAO(nil)
	err := dao.Insert(operator)
	if nil != err {
		if errors.Is(err, model.ErrOperatorExists) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}

	// 发布消息
	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "/operator/join", _req, _rsp)
	return nil
}

func (this *Operator) Leave(_ctx context.Context, _req *proto.OperatorLeaveRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Operator.Leave, req is %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Operator {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "operator is required"
		return nil
	}

	if "" == _req.Workflow {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "workflow is required"
		return nil
	}

	uuid := model.ToUUID(_req.Operator + _req.Workflow)
	dao := model.NewOperatorDAO(nil)
	err := dao.Delete(uuid)
	if nil != err {
		if errors.Is(err, model.ErrOperatorNotFound) {
			_rsp.Status.Code = 2
			_rsp.Status.Message = err.Error()
			return nil
		} else {
			return err
		}
	}
	// 发布消息
	ctx := buildNotifyContext(_ctx, "root")
	publisher.Publish(ctx, "/operator/leave", _req, _rsp)
	return nil
}

func (this *Operator) List(_ctx context.Context, _req *proto.OperatorListRequest, _rsp *proto.OperatorListResponse) error {
	logger.Infof("Received Operator.List, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(100)

	if "" == _req.Workflow {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "workflow is required"
		return nil
	}

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewOperatorDAO(nil)
	query := &model.OperatorQuery{
		Workflow: _req.Workflow,
	}
	total, operators, err := dao.List(offset, count, query)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]string, len(operators))
	for i, operator:= range operators{
		_rsp.Entity[i] = operator.Name
	}
	return nil
}
