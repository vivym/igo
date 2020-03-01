package igo

// AccountInfo holds basic account info
type AccountInfo struct {
	ResponseError
	DsInfo                       DsInfo `json:"dsInfo"`
	HasMinimumDeviceForPhotosWeb bool   `json:"hasMinimumDeviceForPhotosWeb"`
	ICDPEnabled                  bool   `json:"iCDPEnabled"`
	Webservices                  map[string]struct {
		URL    string `json:"url"`
		Status string `json:"status"`
	} `json:"webservices"`
	PcsEnabled        bool `json:"pcsEnabled"`
	TermsUpdateNeeded bool `json:"termsUpdateNeeded"`
	ConfigBag         struct {
		URLs                 map[string]string `json:"urls"`
		AccountCreateEnabled string            `json:"accountCreateEnabled"`
	} `json:"configBag"`
	HsaTrustedBrowser            bool     `json:"hsaTrustedBrowser"`
	AppsOrder                    []string `json:"appsOrder"`
	Version                      int      `json:"version"`
	IsExtendedLogin              bool     `json:"isExtendedLogin"`
	PcsServiceIdentitiesIncluded bool     `json:"pcsServiceIdentitiesIncluded"`
	IsRepairNeeded               bool     `json:"isRepairNeeded"`
	HsaChallengeRequired         bool     `json:"hsaChallengeRequired"`
	RequestInfo                  struct {
		Country  string `json:"country"`
		TimeZone string `json:"timeZone"`
		Region   string `json:"region"`
	} `json:"requestInfo"`
	PcsDeleted bool `json:"pcsDeleted"`
	ICloudInfo struct {
		SafariBookmarksHasMigratedToCloudKit bool `json:"SafariBookmarksHasMigratedToCloudKit"`
	} `json:"iCloudInfo"`
}

// DsInfo holds basic userinfo
type DsInfo struct {
	LastName                  string   `json:"lastName"`
	ICDPEnabled               bool     `json:"iCDPEnabled"`
	TantorMigrated            bool     `json:"tantorMigrated"`
	DsID                      string   `json:"dsid"`
	HSAEnabled                bool     `json:"hsaEnabled"`
	IroncadeMigrated          bool     `json:"ironcadeMigrated"`
	Locale                    string   `json:"locale"`
	BrZoneConsolidated        bool     `json:"brZoneConsolidated"`
	IsManagedAppleID          bool     `json:"isManagedAppleID"`
	GilliganInvited           string   `json:"gilligan-invited"`
	AppleIDAliases            []string `json:"appleIdAliases"`
	HSAVersion                int      `json:"hsaVersion"`
	IsPaidDeveloper           bool     `json:"isPaidDeveloper"`
	CountryCode               string   `json:"countryCode"`
	NotificationID            string   `json:"notificationId"`
	PrimaryEmailVerified      bool     `json:"primaryEmailVerified"`
	ADsID                     string   `json:"aDsID"`
	Locked                    bool     `json:"locked"`
	HasICloudQualifyingDevice bool     `json:"hasICloudQualifyingDevice"`
	PrimaryEmail              string   `json:"primaryEmail"`
	AppleIDEntries            []struct {
		IsPrimary bool   `json:"isPrimary"`
		Type      string `json:"type"`
		Value     string `json:"value"`
	} `json:"appleIdEntries"`
	GilliganEnabled    string `json:"gilligan-enabled"`
	FullName           string `json:"fullName"`
	LanguageCode       string `json:"languageCode"`
	AppleID            string `json:"appleId"`
	FirstName          string `json:"firstName"`
	ICloudAppleIDAlias string `json:"iCloudAppleIdAlias"`
	NotesMigrated      bool   `json:"notesMigrated"`
	HasPaymentInfo     bool   `json:"hasPaymentInfo"`
	PcsDeleted         bool   `json:"pcsDeleted"`
	AppleIDAlias       string `json:"appleIdAlias"`
	BrMigrated         bool   `json:"brMigrated"`
	StatusCode         int    `json:"statusCode"`
	FamilyEligible     bool   `json:"familyEligible"`
}

// AppInfo holds app info
type AppInfo struct {
	AppName   string `json:"app_name"`
	PushTopic string `json:"push_topic"`
}

// ResponseError holds response error
type ResponseError struct {
	// Error  int    `json:"error"`
	Reason string `json:"reason,omitempty"`
}
