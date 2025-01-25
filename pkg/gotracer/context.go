package gotracer

import (
	"context"
)

type contextKey string

var FieldsKey = contextKey("loggerFields")

func FieldsFromContext(ctx context.Context) map[string]interface{} {
	if f, ok := ctx.Value(FieldsKey).(map[string]interface{}); ok {
		return f
	}
	return nil
}
