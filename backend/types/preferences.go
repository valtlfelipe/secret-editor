package types

type Preferences struct {
	General PreferencesGeneral `json:"general"`
}

type PreferencesGeneral struct {
	LastOpenedSecret string `json:"lastOpenedSecret"`
}
