package publisher

import (
	"context"
	"encoding/json"
	"omo-msa-approval/config"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-notification/proto/notification"
)

var (
	DefaultPublisher micro.Event
	filter           map[string]bool
)

func init() {
	filter = make(map[string]bool)
}

func Publish(_ctx context.Context, _action string, _req interface{}, _rsp interface{}) {

	if _, ok := filter[_action]; !ok {
		found := false
		for _, action := range config.Schema.Publisher {
			if action == _action {
				found = true
				break
			}
		}
		filter[_action] = !found
	}

	if filter[_action] {
		return
	}

    head, err := json.Marshal(_req)
	if nil != err {
		logger.Error(err)
        return
	}

    body, err := json.Marshal(_rsp)
	if nil != err {
		logger.Error(err)
        return
	}

	err = DefaultPublisher.Publish(_ctx, &proto.SimpleMessage{
		Action: _action,
		Head:   string(head),
		Body:   string(body),
	})
	if nil != err {
		logger.Error(err)
	}
}
