package push

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"notifications/internal/repo/repomodel"
)

func (r *repo) GetByID(ctx context.Context, pushID int) (*Push, error) {
	return r.selectPush(ctx, "id = $1", pushID)
}

func (r *repo) GetActiveByID(ctx context.Context, pushID int) (*Push, error) {
	return r.selectPush(ctx, "id = $1 AND status = 'active'", pushID)
}

func (r *repo) selectPush(ctx context.Context, condition string, args ...any) (push *Push, err error) {
	var query = "SELECT " + _cols + " FROM push WHERE " + condition

	err = r.db.QueryRow(ctx, query, args...).Scan(fields(push)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repomodel.ErrNotFound
		}
		return nil, err
	}

	return push, nil
}

func (r *repo) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM push WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
