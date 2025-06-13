package apiclient

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"notifications/internal/db"
	"notifications/internal/lib/ctxman"
	"notifications/internal/repo/repomodel"
)

func (r *repo) GetByUserID(ctx context.Context, userID string) (APIClient, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var client APIClient
	err := r.db.QueryRow(ctx, `SELECT client, api_key, permissions FROM api_clients WHERE client = $1`, userID).Scan(&client.Client, &client.APIKey, &client.Permissions)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return APIClient{}, repomodel.ErrNotFound
		}
		return APIClient{}, err
	}

	return client, nil
}
