package types

type Preferences struct {
	General  PreferencesGeneral  `json:"general"`
	Provider PreferencesProvider `json:"provider"`
}

type PreferencesGeneral struct {
	LastOpenedSecret string `json:"lastOpenedSecret"`
}

type PreferencesProvider struct {
	Current    string `json:"current"`
	AWSProfile string `json:"awsProfile"`
}
