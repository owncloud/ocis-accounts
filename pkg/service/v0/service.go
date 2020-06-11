package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/CiscoM31/godata"
	"github.com/golang/protobuf/ptypes/empty"
	mclient "github.com/micro/go-micro/v2/client"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-accounts/pkg/provider"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/protobuf/field_mask"
	"gopkg.in/ldap.v2"
)

// New returns a new instance of Service
func New(cfg *config.Config) Service {
	s := Service{
		Config: cfg,
	}

	return s
}

// Service implements the AccountsServiceHandler interface
type Service struct {
	Config *config.Config
}

func (s Service) getBoundConnection(binddn string, password string) (l *ldap.Conn, err error) {
	l, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", s.Config.LDAP.Hostname, s.Config.LDAP.Port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}

	err = l.Bind(binddn, password)
	if err != nil {
		l.Close()
		return nil, err
	}

	return
}

func (s Service) lookupDN(login string) (binddn string, err error) {
	l, err := s.getBoundConnection(s.Config.LDAP.BindDN, s.Config.LDAP.BindPassword)
	if err != nil {
		return "", err
	}
	defer l.Close()

	filter := fmt.Sprintf("(%s=%s)", s.Config.LDAP.Schema.Username, ldap.EscapeFilter(login))

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		s.Config.LDAP.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return "", err
	}

	switch len(sr.Entries) {
	case 0: // TODO return not found error
	case 1:
		return sr.Entries[0].DN, nil
	default: // TODO return too many results error?
	}
	return "", fmt.Errorf("dn not found for %s", filter)
}

// the auth request is currently hardcoded and has to macth this regex
// userName eq \"teddy\" and password eq \"F&1!b90t111!\"
// TODO allow email to check password?
var authQuery = regexp.MustCompile(`^username eq '(.*)' and password eq '(.*)'$`) // TODO how is ' escaped in the password?

func (s Service) attributeForField(fieldname string) string {
	switch fieldname {
	case "id":
		return s.Config.LDAP.Schema.AccountID
	case "displayname":
		return s.Config.LDAP.Schema.DisplayName
	case "username":
		return s.Config.LDAP.Schema.Username
	// password
	case "uid":
		return s.Config.LDAP.Schema.UID
	case "gid":
		return s.Config.LDAP.Schema.GID
	case "mail":
		return s.Config.LDAP.Schema.Mail
	case "description":
		return s.Config.LDAP.Schema.Description
	}
	// memberof
	return ""
}

var convertAttribute = map[string]func(*ldap.Entry, *config.LDAPSchema, *proto.Account) error{
	"id": func(e *ldap.Entry, s *config.LDAPSchema, a *proto.Account) error {
		a.Id = e.GetAttributeValue(s.AccountID)
		return nil
	},
	"displayname": func(e *ldap.Entry, s *config.LDAPSchema, a *proto.Account) error {
		a.DisplayName = e.GetAttributeValue(s.DisplayName)
		return nil
	},
	"username": func(e *ldap.Entry, s *config.LDAPSchema, a *proto.Account) error {
		a.Username = e.GetAttributeValue(s.Username)
		return nil
	},
	// password
	"uid": func(e *ldap.Entry, s *config.LDAPSchema, a *proto.Account) (err error) {
		a.Uid, err = strconv.ParseInt(e.GetAttributeValue(s.UID), 10, 64)
		return
	},
	"gid": func(e *ldap.Entry, s *config.LDAPSchema, a *proto.Account) (err error) {
		a.Gid, err = strconv.ParseInt(e.GetAttributeValue(s.GID), 10, 64)
		return
	},
	"mail": func(e *ldap.Entry, s *config.LDAPSchema, a *proto.Account) error {
		a.Mail = e.GetAttributeValue(s.Mail)
		return nil
	},
	"description": func(e *ldap.Entry, s *config.LDAPSchema, a *proto.Account) error {
		a.Description = e.GetAttributeValue(s.Description)
		return nil
	},
	//case "memberof":
	//	a.MemberOf = e.GetAttributeValue(s.Config.LDAP.Schema.Description)

	// Future mapping ideas
	/*
		PasswordProfile: &proto.PasswordProfile{
			// Password: write only
			ForceChangePasswordNextSignIn:        false,                    // TODO map to what?
			ForceChangePasswordNextSignInWithMfa: false,                    // TODO map to what?
			LastPasswordChangeDateTime:           &timestamppb.Timestamp{}, // TODO map to what?
			PasswordPolicies:                     []string{},               // TODO map to what?
			// https://www.arctiq.ca/our-blog/2018/9/4/implementing-a-password-policy-in-an-ldap-directory/
			// openldap has a ppolicy: https://www.openldap.org/software/man.cgi?query=slapo-ppolicy&apropos=0&sektion=0&manpath=OpenLDAP+2.4-Release&format=html
			// it uses a pwdpolicy objectclass

			// last login is available in AD https://ldapwiki.com/wiki/Last%20Login%20Time#:~:text=The%20Last%20Login%20Time%20feature,with%20a%20user%2Ddefined%20format.
			// AD: lastLogonTimestamp (is replicated batween ad controllers after two weeks) or lastLogon
			// oracle: pwdLastAuthTime
			// openldap https://www.openldap.org/doc/admin24/overlays.html#Access%20Logging
		},
	*/

	// this needs to be implemented by the ldap server
	// AccountEnabled: len(sr.Entries[i].GetAttributeValue(s.Config.LDAP.Schema.AccountID))==0,
	// OpenLdap: ppolicy overlay and then pwdAccountLockedTime
	// if pwdAccountLockedTime is present the user is disabled
	// or shadowexpire: 0
	// AD uses this filter (&(objectCategory=person)(objectClass=user)(userAccountControl:1.2.840.113556.1.4.803:=2))
	// ... well it is more complicated: https://ldapwiki.com/wiki/AD%20Determining%20Password%20Expiration

	//IsResourceAccount: false, // TODO could be represented by an attribute or an objectclass or even by a svc_ prefix

	//CreationType:      "", // TODO could be represented by an attribute or an objectclass

	// TODO identities
	//Identities: nil, // TODO map to what?

	// on premise attributes? no longer needed. in the graph api we can add them when implementing sync
	//ExternalUserState:               "",                       // TODO could be represented by an attribute or an objectclass
	//ExternalUserStateChangeDateTime: &timestamppb.Timestamp{}, // TODO needs new attribute

	//CreatedDateTime: &timestamppb.Timestamp{}, // TODO -> operational attribute createTimestamp
	//DeletedDateTime: &timestamppb.Timestamp{}, // TODO map to what?
}

var defaultFields = []string{
	// "id", is always returned
	"displayname",
	"username",
	"mail",
	"description",
}

func (s Service) entryToAccount(c context.Context, e *ldap.Entry, m *field_mask.FieldMask, a *proto.Account) error {
	// id is always returned
	a.Id = e.GetAttributeValue(s.Config.LDAP.Schema.AccountID)

	var fields []string
	if m == nil || len(m.Paths) == 0 {
		fields = defaultFields
	} else {
		// TODO check every path: if it is about groups we need to fetch groups and the selected sub properties
		fields = m.Paths
	}
	for i := range fields {
		if err := convertAttribute[fields[i]](e, &s.Config.LDAP.Schema, a); err != nil {
			log.Error().Err(err).Interface("entry", e).Msg("skipping user")
			continue
		}
	}
	return nil
}

// ListAccounts implements the AccountsServiceHandler interface
// the query contains account properties
func (s Service) ListAccounts(ctx context.Context, in *proto.ListAccountsRequest, res *proto.ListAccountsResponse) (err error) {

	var binddn string
	var password string

	// check if this looks like an auth request
	match := authQuery.FindStringSubmatch(in.Query)
	if len(match) == 3 {

		binddn, err = s.lookupDN(match[1])
		if err != nil {
			log.Error().Err(err).Msg("ListAccounts with auth request")
			return
		}
		log.Debug().Str("username", match[1]).Str("binddn", binddn).Msg("ListAccounts with auth request")

		password = match[2]
		// remove password from query
		in.Query = fmt.Sprintf("username eq '%s'", match[1])
	} else {
		log.Debug().Str("query", in.Query).Uint32("page-size", in.PageSize).Str("page-token", in.PageToken).Msg("ListAccounts")
		binddn = s.Config.LDAP.BindDN
		password = s.Config.LDAP.BindPassword
	}

	filter := "(&)" // see Absolute True and False Filters in https://tools.ietf.org/html/rfc4526#section-2

	if in.Query != "" {
		// parse the query like an odata filter
		var q *godata.GoDataFilterQuery
		if q, err = godata.ParseFilterString(in.Query); err != nil {
			return
		}

		// convert to ldap filter
		filter, err = provider.BuildLDAPFilter(q, &s.Config.LDAP.Schema)
		if err != nil {
			return
		}
	}

	log.Debug().Str("filter", filter).Msg("using filter")

	var l *ldap.Conn
	l, err = s.getBoundConnection(binddn, password)
	if err != nil {
		return
	}
	defer l.Close()

	// TODO combine the parsed query with a query filter from the config, eg. fmt.Sprintf(s.Config.LDAP.UserFilter, clientID)

	attributes := []string{"dn"}
	if in.FieldMask != nil && len(in.FieldMask.Paths) > 0 {
		for i := range in.FieldMask.Paths {
			attributes = append(attributes, s.attributeForField(in.FieldMask.Paths[i]))
		}
	} else {
		for i := range defaultFields {
			attributes = append(attributes, s.attributeForField(defaultFields[i]))
		}
	}

	var controls []ldap.Control
	if in.PageSize != 0 {
		paging := ldap.NewControlPaging(in.PageSize)
		// TODO cookies only work on a connection basis, so we need to keep the connection alive to resume pagination with that cookie.
		// we cun build our own cookie to store and retrieve the connection in a map
		// paging.Cookie = base64.decode(in.PageToken)
		controls = []ldap.Control{paging}
	}

	// Search for the given clientID
	searchRequest := ldap.NewSearchRequest(
		s.Config.LDAP.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		controls,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return err
	}

	log.Debug().Interface("entries", sr.Entries).Msg("entries")

	res.Accounts = make([]*proto.Account, 0)
	for i := range sr.Entries {
		a := &proto.Account{}
		err := s.entryToAccount(ctx, sr.Entries[i], in.FieldMask, a)
		if err != nil {
			continue
		}
		res.Accounts = append(res.Accounts, a)
	}

	return nil
}

// GetAccount implements the AccountsServiceHandler interface
func (s Service) GetAccount(c context.Context, req *proto.GetAccountRequest, res *proto.Account) (err error) {

	l, err := s.getBoundConnection(s.Config.LDAP.BindDN, s.Config.LDAP.BindPassword)
	if err != nil {
		return err
	}
	defer l.Close()

	// TODO combine the query with a query filter from the config, eg. fmt.Sprintf(s.Config.LDAP.UserFilter, clientID)
	filter := fmt.Sprintf("(%s=%s)", s.Config.LDAP.Schema.AccountID, ldap.EscapeFilter(req.Id))

	attributes := []string{"dn"}
	if req.FieldMask != nil && len(req.FieldMask.Paths) > 0 {
		for i := range req.FieldMask.Paths {
			attributes = append(attributes, s.attributeForField(req.FieldMask.Paths[i]))
		}
	} else {
		for i := range defaultFields {
			attributes = append(attributes, s.attributeForField(defaultFields[i]))
		}
	}

	// Search for the given clientID
	searchRequest := ldap.NewSearchRequest(
		s.Config.LDAP.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		[]ldap.Control{ldap.NewControlPaging(2)}, // at max two results
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return err
	}

	log.Debug().Interface("entries", sr.Entries).Msg("entries")

	switch len(sr.Entries) {
	case 0: // TODO return not found error
	case 1:
		err = s.entryToAccount(c, sr.Entries[0], req.FieldMask, res)
	default: // TODO return too many results error?
	}

	return nil
}

// CreateAccount implements the AccountsServiceHandler interface
func (s Service) CreateAccount(c context.Context, req *proto.CreateAccountRequest, res *proto.Account) (err error) {
	return errors.New("not implemented")
}

// UpdateAccount implements the AccountsServiceHandler interface
func (s Service) UpdateAccount(c context.Context, req *proto.UpdateAccountRequest, res *proto.Account) (err error) {
	return errors.New("not implemented")
}

// DeleteAccount implements the AccountsServiceHandler interface
func (s Service) DeleteAccount(c context.Context, req *proto.DeleteAccountRequest, res *empty.Empty) (err error) {
	return errors.New("not implemented")
}

// ListGroups implements the AccountsServiceHandler interface
func (s Service) ListGroups(c context.Context, req *proto.ListGroupsRequest, res *proto.ListGroupsResponse) (err error) {
	return errors.New("not implemented")
}

// GetGroup implements the AccountsServiceHandler interface
func (s Service) GetGroup(c context.Context, req *proto.GetGroupRequest, res *proto.Group) (err error) {
	return errors.New("not implemented")
}

// CreateGroup implements the AccountsServiceHandler interface
func (s Service) CreateGroup(c context.Context, req *proto.CreateGroupRequest, res *proto.Group) (err error) {
	return errors.New("not implemented")
}

// UpdateGroup implements the AccountsServiceHandler interface
func (s Service) UpdateGroup(c context.Context, req *proto.UpdateGroupRequest, res *proto.Group) (err error) {
	return errors.New("not implemented")
}

// DeleteGroup implements the AccountsServiceHandler interface
func (s Service) DeleteGroup(c context.Context, req *proto.DeleteGroupRequest, res *empty.Empty) (err error) {
	return errors.New("not implemented")
}

// RegisterSettingsBundles pushes the settings bundle definitions for this extension to the ocis-settings service.
func RegisterSettingsBundles(l *olog.Logger) {
	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	service := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)

	requests := []settings.SaveSettingsBundleRequest{
		generateSettingsBundleProfileRequest(),
	}

	for i := range requests {
		res, err := service.SaveSettingsBundle(context.Background(), &requests[i])
		if err != nil {
			l.Err(err).
				Msg("Error registering settings bundle")
		} else {
			l.Info().
				Str("bundle key", res.SettingsBundle.Identifier.BundleKey).
				Msg("Successfully registered settings bundle")
		}
	}
}
