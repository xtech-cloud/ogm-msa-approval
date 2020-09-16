package handler

import (
	"context"
	"github.com/micro/go-micro/v2/metadata"
)

func buildNotifyContext(_ctx context.Context, _operator string) context.Context {
	optLabel := _operator
	optType := "None"

	// 使用metadata覆盖默认值
	meta, ok := metadata.FromContext(_ctx)
	if ok {
		if _, exists := meta["Operator-Label"]; exists {
			optLabel = meta["Operator-Label"]
		}
		if _, exists := meta["Operator-Type"]; exists {
			optType = meta["Operator-Type"]
		}
	}

	ctx := metadata.NewContext(context.Background(), map[string]string{
		"Operator-Label": optLabel,
		"Operator-Type":  optType,
	})

	return ctx
}
