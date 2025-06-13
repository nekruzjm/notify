package user

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"

	"notifications/internal/db"
	"notifications/internal/lib/ctxman"
	"notifications/internal/repo/repomodel"
)

func (r *repo) GetActiveByPersonExternalRef(ctx context.Context, personExternalRef string) (*User, error) {
	return r.selectUser(ctx, "person_external_ref = $1 AND status != 'deleted'", personExternalRef)
}

func (r *repo) GetByUserID(ctx context.Context, userID int) (*User, error) {
	return r.selectUser(ctx, "user_id = $1 AND status != 'deleted'", userID)
}

func (r *repo) GetActiveByPhone(ctx context.Context, phone string) (*User, error) {
	return r.selectUser(ctx, "phone = $1 AND status != 'deleted'", phone)
}

func (r *repo) GetByUserIDs(ctx context.Context, userIDs []int) ([]*User, error) {
	return r.selectUsers(ctx, "user_id = ANY($1) AND status != 'deleted'", userIDs)
}

func (r *repo) GetByPhones(ctx context.Context, phones []string) ([]*User, error) {
	return r.selectUsers(ctx, "phone = ANY($1) AND status != 'deleted'", phones)
}

func (r *repo) GetTokensByUserIDs(ctx context.Context, userIDs []string) ([]User, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var query = `SELECT user_id, token, language FROM users WHERE user_id = ANY($1) AND status != 'deleted'`

	rows, err := r.db.Query(ctx, query, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = make([]User, 0, len(userIDs))

	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserID, &user.Token, &user.Language)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, repomodel.ErrNotFound
	}

	return users, nil
}

func (r *repo) GetTokensWithLimit(ctx context.Context, lastID int) ([]User, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var query = `SELECT user_id, token, language FROM users WHERE user_id > $1 AND status != 'deleted' ORDER BY user_id LIMIT 1000`

	rows, err := r.db.Query(ctx, query, lastID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = make([]User, 0, 1000)

	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserID, &user.Token, &user.Language)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, repomodel.ErrNotFound
	}

	return users, nil
}

func (r *repo) GetUserIDsByEventID(ctx context.Context, eventID int) ([]int, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	rows, err := r.db.Query(ctx, "SELECT user_id FROM user_event_relations WHERE event_id = $1 ", eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids = make([]int, 0)

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, repomodel.ErrNotFound
	}

	return ids, nil
}

func (r *repo) GetRelationsByEventID(ctx context.Context, eventID int) ([]EventRelation, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	rows, err := r.db.Query(ctx, "SELECT user_id, language FROM user_event_relations WHERE event_id = $1 ", eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations = make([]EventRelation, 0)

	for rows.Next() {
		var rel EventRelation
		err = rows.Scan(&rel.UserID, &rel.Lang)
		if err != nil {
			return nil, err
		}
		relations = append(relations, rel)
	}

	if len(relations) == 0 {
		return nil, repomodel.ErrNotFound
	}

	return relations, nil
}

func (r *repo) selectUser(ctx context.Context, c string, args ...any) (*User, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var (
		user    = new(User)
		builder strings.Builder
	)

	builder.WriteString("SELECT ")
	builder.WriteString(_cols)
	builder.WriteString(" FROM users WHERE ")
	builder.WriteString(c)

	query := builder.String()

	err := r.db.QueryRow(ctx, query, args...).Scan(fields(user)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repomodel.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *repo) selectUsers(ctx context.Context, c string, args ...any) ([]*User, error) {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	var builder strings.Builder
	builder.WriteString("SELECT ")
	builder.WriteString(_cols)
	builder.WriteString(" FROM users WHERE ")
	builder.WriteString(c)

	query := builder.String()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = make([]*User, 0)

	for rows.Next() {
		var user = new(User)
		err = rows.Scan(fields(user)...)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, repomodel.ErrNotFound
	}

	return users, nil
}
