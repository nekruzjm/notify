package main

import (
	"testing"

	"go.uber.org/fx"

	"notifications/internal/api/transport"
	"notifications/internal/db"
	"notifications/internal/db/tx"
	"notifications/internal/gateway"
	"notifications/internal/handler"
	"notifications/internal/repo"
	"notifications/internal/service"
	"notifications/pkg/lib"
)

func Test_Deps(t *testing.T) {
	if err := fx.ValidateApp(deps()); err != nil {
		t.Error("err occurred during dependency injection:", err)
		return
	}
}

func deps() fx.Option {
	return fx.Options(
		transport.Module,
		handler.Module,
		service.Module,
		gateway.Module,
		repo.Module,
		db.Module,
		tx.Module,
		lib.Module,
	)
}
