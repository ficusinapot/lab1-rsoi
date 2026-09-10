package version

type AppVersionInfo struct {
	Version   string `json:"version" doc:"Application version"`
	BuildTime string `json:"build_time" doc:"Application build time"`
	Branch    string `json:"branch" doc:"Application build branch"`
	Commit    string `json:"commit" doc:"Application build commit"`
}

type AppVersionOutput struct {
	Body AppVersionBody
}

type AppVersionBody struct {
	Version AppVersionInfo `json:"data" doc:"Application version"`
}
