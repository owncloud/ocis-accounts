package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type CS3Repo struct {
	serviceID string
	dataPath  string
	rootPath  string

	tm            token.Manager
	storageClient provider.ProviderAPIClient
}

func NewCS3Repo(secret string) (Repo, error) {
	tokenManager, err := jwt.New(map[string]interface{}{
		"secret": "Pive-Fumkiu4",
	})

	if err != nil {
		return nil, err
	}

	client, err := pool.GetStorageProviderServiceClient("localhost:9185")
	if err != nil {
		return nil, err
	}

	return CS3Repo{tm: tokenManager, storageClient: client}, nil

}

func (r CS3Repo) WriteAccount(ctx context.Context, a *proto.Account) (err error) {
	t, err := r.authenticate(ctx)
	if err != nil {
		return err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)
	if err := r.makeRootDirIfNotExist(ctx); err != nil {
		return err
	}

	var by []byte
	if by, err = json.Marshal(a); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not marshal account: %v", err.Error())
	}

	ureq, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:9187/data/accounts/%s", a.Id), bytes.NewReader(by))
	if err != nil {
		return err
	}

	ureq.Header.Add("x-access-token", t)
	cl := http.Client{
		Transport: http.DefaultTransport,
	}

	if _, err := cl.Do(ureq); err != nil {
		return err
	}

	return nil
}

func (r CS3Repo) LoadAccount(ctx context.Context, id string, a *proto.Account) (err error) {
	t, err := r.authenticate(ctx)
	if err != nil {
		return err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)

	ureq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:9187/data/accounts/%s", id), nil)
	if err != nil {
		return err
	}

	ureq.Header.Add("x-access-token", t)
	cl := http.Client{
		Transport: http.DefaultTransport,
	}

	if _, err = cl.Do(ureq); err != nil {
		return err
	}

	return nil
}

func (r CS3Repo) WriteGroup(ctx context.Context, g *proto.Group) (err error) {
	panic("implement me")
}

func (r CS3Repo) LoadGroup(ctx context.Context, id string, g *proto.Group) (err error) {
	panic("implement me")
}

func (r CS3Repo) authenticate(ctx context.Context) (token string, err error) {
	return r.tm.MintToken(ctx, &user.User{
		Id:     &user.UserId{},
		Groups: []string{},
	})
}

func (r CS3Repo) makeRootDirIfNotExist(ctx context.Context) error {
	var rootPathRef = &provider.Reference{
		Spec: &provider.Reference_Path{Path: "/meta/accounts"},
	}

	resp, err := r.storageClient.Stat(ctx, &provider.StatRequest{
		Ref: rootPathRef,
	})

	if err != nil {
		return err
	}

	if resp.Status.Code == v1beta11.Code_CODE_NOT_FOUND {
		_, err := r.storageClient.CreateContainer(ctx, &provider.CreateContainerRequest{
			Ref: rootPathRef,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
