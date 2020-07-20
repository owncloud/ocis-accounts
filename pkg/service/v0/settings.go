package service

import (
	"context"

	mclient "github.com/micro/go-micro/v2/client"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
	ssvc "github.com/owncloud/ocis-settings/pkg/service/v0"
)

const (
	settingUuidProfileLanguage = "aa8cfbe5-95d4-4f7e-a032-c3c01f5f062f"
)

// RegisterSettingsBundles pushes the settings bundle definitions for this extension to the ocis-settings service.
func RegisterSettingsBundles(l *olog.Logger) {
	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	service := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)

	bundleRequests := []settings.SaveSettingsBundleRequest{
		generateSettingsBundleProfileRequest(),
	}

	for i := range bundleRequests {
		res, err := service.SaveSettingsBundle(context.Background(), &bundleRequests[i])
		if err != nil {
			l.Err(err).Str("bundle", res.SettingsBundle.Id).Msg("Error registering bundle")
		} else {
			l.Info().Str("bundle", res.SettingsBundle.Id).Msg("Successfully registered bundle")
		}
	}

	permissionRequests := generateProfilePermissionsRequests()
	for i := range permissionRequests {
		res, err := service.AddSettingToSettingsBundle(context.Background(), &permissionRequests[i])
		bundleId := permissionRequests[i].BundleId
		if err != nil {
			l.Err(err).Str("bundle", bundleId).Str("setting", res.Setting.Id).Msg("Error adding setting to bundle")
		} else {
			l.Info().Str("bundle", bundleId).Str("setting", res.Setting.Id).Msg("Successfully added setting to bundle")
		}
	}
}

var languageSetting = settings.Setting_SingleChoiceValue{
	SingleChoiceValue: &settings.SingleChoiceListSetting{
		Options: []*settings.ListOption{
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "cs",
					},
				},
				DisplayValue: "Czech",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "de",
					},
				},
				DisplayValue: "Deutsch",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "en",
					},
				},
				DisplayValue: "English",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "es",
					},
				},
				DisplayValue: "Español",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "fr",
					},
				},
				DisplayValue: "Français",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "gl",
					},
				},
				DisplayValue: "Galego",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "it",
					},
				},
				DisplayValue: "Italiano",
			},
		},
	},
}

func generateSettingsBundleProfileRequest() settings.SaveSettingsBundleRequest {
	return settings.SaveSettingsBundleRequest{
		SettingsBundle: &settings.SettingsBundle{
			Id:        "2a506de7-99bd-4f0d-994e-c38e72c28fd9",
			Name:      "profile",
			Extension: "ocis-accounts",
			Type:      settings.SettingsBundle_TYPE_DEFAULT,
			Resource: &settings.Resource{
				Type: settings.Resource_TYPE_SYSTEM,
			},
			DisplayName: "Profile",
			Settings: []*settings.Setting{
				{
					Id:          settingUuidProfileLanguage,
					Name:        "language",
					DisplayName: "Language",
					Description: "User language",
					Resource: &settings.Resource{
						Type: settings.Resource_TYPE_USER,
					},
					Value: &languageSetting,
				},
			},
		},
	}
}

func generateProfilePermissionsRequests() []settings.AddSettingToSettingsBundleRequest {
	// TODO: we don't want to set up permissions for settings manually in the future. Instead each setting should come with
	// a set of default permissions for the default roles (guest, user, admin).
	return []settings.AddSettingToSettingsBundleRequest{
		{
			BundleId: ssvc.BundleUuidRoleAdmin,
			Setting: &settings.Setting{
				Id:          "7d81f103-0488-4853-bce5-98dcce36d649",
				Name:        "language-create",
				DisplayName: "Permission to set the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUuidProfileLanguage,
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.PermissionSetting{
						Operation:  settings.PermissionSetting_OPERATION_CREATE,
						Constraint: settings.PermissionSetting_CONSTRAINT_OWN,
					},
				},
			},
		},
		{
			BundleId: ssvc.BundleUuidRoleAdmin,
			Setting: &settings.Setting{
				Id:          "04ef2fd3-e724-48f6-a411-129dd461c820",
				Name:        "language-read",
				DisplayName: "Permission to read the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUuidProfileLanguage,
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.PermissionSetting{
						Operation:  settings.PermissionSetting_OPERATION_READ,
						Constraint: settings.PermissionSetting_CONSTRAINT_OWN,
					},
				},
			},
		},
		{
			BundleId: ssvc.BundleUuidRoleAdmin,
			Setting: &settings.Setting{
				Id:          "30ac1e63-10e2-4ef8-bf0a-941cd5b56c5c",
				Name:        "language-update",
				DisplayName: "Permission to update the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUuidProfileLanguage,
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.PermissionSetting{
						Operation:  settings.PermissionSetting_OPERATION_UPDATE,
						Constraint: settings.PermissionSetting_CONSTRAINT_OWN,
					},
				},
			},
		},
		{
			BundleId: ssvc.BundleUuidRoleUser,
			Setting: &settings.Setting{
				Id:   "dcaeb961-da25-46f2-9892-731603a20d3b",
				Name: "language-read",
				DisplayName: "Permission to read the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUuidProfileLanguage,
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.PermissionSetting{
						Operation:  settings.PermissionSetting_OPERATION_READ,
						Constraint: settings.PermissionSetting_CONSTRAINT_OWN,
					},
				},
			},
		},
		{
			BundleId: ssvc.BundleUuidRoleGuest,
			Setting: &settings.Setting{
				Id:   "ca878636-8b1a-4fae-8282-8617a4c13597",
				Name: "language-read",
				DisplayName: "Permission to read the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUuidProfileLanguage,
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.PermissionSetting{
						Operation:  settings.PermissionSetting_OPERATION_READ,
						Constraint: settings.PermissionSetting_CONSTRAINT_OWN,
					},
				},
			},
		},
	}
}
