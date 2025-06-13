package main

import (
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

func main() {
	fx.New(
		transport.Module,
		handler.Module,
		service.Module,
		gateway.Module,
		repo.Module,
		db.Module,
		tx.Module,
		lib.Module,
	).Run()

	// Add some stats for admin panel
}
