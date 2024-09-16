package api

import (
	"context"
	"digital-wallet/pkg/errs"
)

func IsAuthorizedUser(ctx context.Context, recordOwner string) error {
	actor := GetActor(ctx)
	actorId := GetActorID(ctx)
	if actor != nil && actorId != nil && (*actor == AppActorAdmin || *actor == AppActorSystem || *actorId == recordOwner) {
		return nil
	}
	return errs.NewUnauthorizedError("Unauthorized", "", nil)
}

func IsAdmin(ctx context.Context) error {
	actor := GetActor(ctx)
	if actor != nil && (*actor == AppActorAdmin || *actor == AppActorSystem) {
		return nil
	}
	return errs.NewUnauthorizedError("Unauthorized", "", nil)
}

func IsSystem(ctx context.Context) error {
	actor := GetActor(ctx)
	if actor != nil && *actor == AppActorSystem {
		return nil
	}
	return errs.NewUnauthorizedError("Unauthorized", "", nil)
}
