package storage

import (
	"context"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
)

type Repo interface {
	WriteAccount(ctx context.Context, a *proto.Account) (err error)
	LoadAccount(ctx context.Context, id string, a *proto.Account) (err error)
	WriteGroup(ctx context.Context, g *proto.Group) (err error)
	LoadGroup(ctx context.Context, id string, g *proto.Group) (err error)
}
