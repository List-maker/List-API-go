package users

import "encoding/json"

type UserSettings struct {
	// TODO
}

var defaultSettings = UserSettings{}
var jsonDefaultSettings, _ = json.Marshal(defaultSettings)

func parseUserSettings(rawSettings string) (settings UserSettings) {
	err := json.Unmarshal([]byte(rawSettings), &settings)
	if err != nil {
		return defaultSettings
	}
	return settings
}
