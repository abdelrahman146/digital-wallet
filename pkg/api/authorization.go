package api

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
)

func IsAuthorizedUser(ctx context.Context, recordOwner string) error {
	actor := GetActor(ctx)
	actorId := GetActorID(ctx)
	if actor != "" && actorId != "" && (actor == AppActorAdmin || actor == AppActorSystem || actorId == recordOwner) {
		return nil
	}
	return errs.NewUnauthorizedError("Unauthorized", "", nil)
}

func IsAdmin(ctx context.Context) error {
	actor := GetActor(ctx)
	if actor != "" && (actor == AppActorAdmin || actor == AppActorSystem) {
		return nil
	}
	return errs.NewUnauthorizedError("Unauthorized", "", nil)
}

func IsSystem(ctx context.Context) error {
	actor := GetActor(ctx)
	if actor == AppActorSystem {
		return nil
	}
	return errs.NewUnauthorizedError("Unauthorized", "", nil)
}
