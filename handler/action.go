package handler

import (
	"context"
	"omo-msa-approval/model"

	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-approval/proto/approval"
)

type Action struct{}

func (this *Action) Query(_ctx context.Context, _req *proto.ActionQueryRequest, _rsp *proto.ActionQueryResponse) error {
	logger.Infof("Received Action.Query, req is %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	count := int64(100)

	if _req.Offset > 0 {
		offset = _req.Offset
	}

	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewActionDAO(nil)
	query := &model.ActionQuery{
		State:    int(_req.State),
		Workflow: _req.Workflow,
		Operator: _req.Operator,
		Task:     _req.Task,
	}
	total, actions, err := dao.List(offset, count, query)
	if nil != err {
		return nil
	}

	_rsp.Total = uint64(total)
	_rsp.Entity = make([]*proto.ActionEntity, len(actions))
	for i, action := range actions {
		_rsp.Entity[i] = &proto.ActionEntity{
			Uuid:      action.UUID,
			Task:      action.Task,
			Operator:  action.Operator,
			State:     proto.ActionStatus(action.State),
			UpdatedAt: action.UpdatedAt.UTC().Unix(),
		}
	}
	return nil
}
