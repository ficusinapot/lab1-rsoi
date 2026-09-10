package version

type Handler struct {
	version AppVersionInfo
}

func NewHandler(version, buildTime, branch, commit string) *Handler {
	return &Handler{
		version: AppVersionInfo{
			Version:   version,
			BuildTime: buildTime,
			Branch:    branch,
			Commit:    commit,
		},
	}
}
