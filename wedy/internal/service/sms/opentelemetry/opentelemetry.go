package opentelemetry

import (
	"GoProj/wedy/internal/service/sms"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Decorater struct {
	svc    sms.Service
	tracer trace.Tracer
}

func NewDecorater(svc sms.Service, tracer trace.Tracer) *Decorater {
	return &Decorater{
		svc:    svc,
		tracer: tracer,
	}
}
func (d *Decorater) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	ctx, span := d.tracer.Start(ctx, "SMS Service Send")
	defer span.End()
	span.SetAttributes(attribute.String("tplId", tplId))
	err := d.svc.Send(ctx, tplId, args, numbers...)
	if err != nil {
		span.RecordError(err)
	}
	return err
}
