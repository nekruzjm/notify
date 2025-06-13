package firebase

import (
	"context"

	"firebase.google.com/go/v4/messaging"
)

func (f *fb) SendPush(ctx context.Context, message *messaging.Message) (string, error) {
	return f.client.Send(ctx, message)
}

func (f *fb) SendPushDryRun(ctx context.Context, message *messaging.Message) (string, error) {
	return f.client.SendDryRun(ctx, message)
}

func (f *fb) SendEach(ctx context.Context, messages []*messaging.Message) (*messaging.BatchResponse, error) {
	return f.client.SendEach(ctx, messages)
}

func (f *fb) SendEachDryRun(ctx context.Context, messages []*messaging.Message) (*messaging.BatchResponse, error) {
	return f.client.SendEachDryRun(ctx, messages)
}

func (f *fb) SendMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
	return f.client.SendEachForMulticast(ctx, message)
}

func (f *fb) SendMulticastDryRun(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
	return f.client.SendEachForMulticastDryRun(ctx, message)
}
