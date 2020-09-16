package storage

import (
	"encoding/json"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/rs/zerolog"
	"io/ioutil"
	"path/filepath"
	"sync"
)

var groupLock sync.Mutex

type Repo struct {
	serviceID string
	dataPath  string
	log       zerolog.Logger
}

func New(serviceID string, dataPath string, log zerolog.Logger) Repo {
	return Repo{
		serviceID: serviceID,
		dataPath:  dataPath,
		log:       log.With().Str("id", serviceID).Logger(),
	}
}

// WriteAccount to the storage
func (r Repo) WriteAccount(a *proto.Account) (err error) {
	// leave only the group id
	r.deflateMemberOf(a)

	var bytes []byte
	if bytes, err = json.Marshal(a); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not marshal account: %v", err.Error())
	}

	path := filepath.Join(r.dataPath, "accounts", a.Id)

	if err = ioutil.WriteFile(path, bytes, 0600); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not write account: %v", err.Error())
	}
	return
}

// LoadAccount from the storage
func (r Repo) LoadAccount(id string, a *proto.Account) (err error) {
	path := filepath.Join(r.dataPath, "accounts", id)

	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return merrors.NotFound(r.serviceID, "could not read account: %v", err.Error())
	}

	if err = json.Unmarshal(data, a); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not unmarshal account: %v", err.Error())
	}
	return
}

// WriteGroup persists a given group to the storage
func (r Repo) WriteGroup(g *proto.Group) (err error) {
	// leave only the member id
	r.deflateMembers(g)

	var bytes []byte
	if bytes, err = json.Marshal(g); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not marshal group: %v", err.Error())
	}

	path := filepath.Join(r.dataPath, "groups", g.Id)

	groupLock.Lock()
	defer groupLock.Unlock()
	if err = ioutil.WriteFile(path, bytes, 0600); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not write group: %v", err.Error())
	}
	return
}

// LoadGroup from the storage
func (r Repo) LoadGroup(id string, g *proto.Group) (err error) {
	path := filepath.Join(r.dataPath, "groups", id)

	groupLock.Lock()
	defer groupLock.Unlock()
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return merrors.NotFound(r.serviceID, "could not read group: %v", err.Error())
	}

	if err = json.Unmarshal(data, g); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not unmarshal group: %v", err.Error())
	}

	return
}

// deflateMemberOf replaces the groups of a user with an instance that only contains the id
func (r Repo) deflateMemberOf(a *proto.Account) {
	if a == nil {
		return
	}
	deflated := []*proto.Group{}
	for i := range a.MemberOf {
		if a.MemberOf[i].Id != "" {
			deflated = append(deflated, &proto.Group{Id: a.MemberOf[i].Id})
		} else {
			// TODO fetch and use an id when group only has a name but no id
			r.log.Error().Str("id", a.Id).Interface("group", a.MemberOf[i]).Msg("resolving groups by name is not implemented yet")
		}
	}
	a.MemberOf = deflated
}

// deflateMembers replaces the users of a group with an instance that only contains the id
func (r Repo) deflateMembers(g *proto.Group) {
	if g == nil {
		return
	}
	deflated := []*proto.Account{}
	for i := range g.Members {
		if g.Members[i].Id != "" {
			deflated = append(deflated, &proto.Account{Id: g.Members[i].Id})
		} else {
			// TODO fetch and use an id when group only has a name but no id
			r.log.Error().Str("id", g.Id).Interface("account", g.Members[i]).Msg("resolving members by name is not implemented yet")
		}
	}
	g.Members = deflated
}
