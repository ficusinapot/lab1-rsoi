package version

var (
	applicationName    = "service"
	applicationVersion = "0.1.0"
	buildTime          = ""
	branch             = ""
	commit             = ""
)

type Info struct {
	Name      string
	Version   string
	BuildTime string
	Branch    string
	Commit    string
}

func GetApplicationVersion() string {
	return applicationVersion
}

func GetInfo() Info {
	return Info{
		Name:      applicationName,
		Version:   applicationVersion,
		BuildTime: buildTime,
		Branch:    branch,
		Commit:    commit,
	}
}
