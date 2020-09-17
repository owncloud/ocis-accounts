package storage

import (
	"context"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFoo(t *testing.T) {
	r, err := NewCS3Repo("fooo")
	assert.NoError(t, err)

	err = r.WriteAccount(context.Background(), &proto.Account{
		Id:             "fefef-egegweg-gegeg",
		AccountEnabled: true,
		DisplayName:    "Mike Jones",
		Mail:           "mike@example.com",
	})

	assert.NoError(t, err)

}
