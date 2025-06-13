package user

import (
	"context"

	"notifications/internal/db"
	"notifications/internal/lib/ctxman"
)

func (r *repo) Create(ctx context.Context, user User) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `
		INSERT INTO users (user_id, 
		                   phone, 
		                   person_external_ref, 
		                   token, 
		                   status, 
		                   language,
		                   country_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE
			SET phone = $2,
				person_external_ref = $3,
			    token        = $4,
				status       = $5,
				language	 = $6,
				country_id   = $7,
				updated_at    = now()`,
		user.UserID,
		user.Phone,
		user.PersonExternalRef,
		user.Token,
		user.Status,
		user.Language,
		user.CountryID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) UpdateToken(ctx context.Context, userID int, token string) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `UPDATE users SET token = $1, updated_at = now() WHERE user_id = $2`, token, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) UpdateStatus(ctx context.Context, userID int, status string) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `UPDATE users SET status = $1, updated_at = now() WHERE user_id = $2`, status, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) UpdatePhone(ctx context.Context, userID int, phone string) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `UPDATE users SET phone = $1, updated_at = now() WHERE user_id = $2`, phone, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) UpdatePersonExternalRef(ctx context.Context, userID int, personExternalRef string) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `UPDATE users SET person_external_ref = $1, updated_at = now() WHERE user_id = $2`, personExternalRef, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) UpdateUserSettings(ctx context.Context, userID int, language string, isEnabled bool) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `
			UPDATE users SET 
				push_enabled = $1, 
				language = $2, 
				updated_at = now() 
			WHERE user_id = $3`, isEnabled, language, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) DeleteByUserID(ctx context.Context, userID int) error {
	ctx = ctxman.Save(ctx, ctxman.Info{
		DBName:    db.Notifications,
		IsReplica: false,
	})

	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	return nil
}
