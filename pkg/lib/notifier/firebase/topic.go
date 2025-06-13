package firebase

import (
	"context"

	"firebase.google.com/go/v4/messaging"
)

func (f *fb) SubscribeTokens(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error) {
	return f.client.SubscribeToTopic(ctx, tokens, topic)
}

func (f *fb) UnsubscribeTokens(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error) {
	return f.client.UnsubscribeFromTopic(ctx, tokens, topic)
}

func (f *fb) Subscribe(ctx context.Context, token string, topic string) (*messaging.TopicManagementResponse, error) {
	return f.client.SubscribeToTopic(ctx, []string{token}, topic)
}

func (f *fb) Unsubscribe(ctx context.Context, token string, topic string) (*messaging.TopicManagementResponse, error) {
	return f.client.UnsubscribeFromTopic(ctx, []string{token}, topic)
}
