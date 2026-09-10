package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ficusinapot/ds/cmd/service/commands"
	serviceconfig "github.com/ficusinapot/ds/cmd/service/config"
	"github.com/ficusinapot/ds/cmd/service/version"
	"github.com/ficusinapot/ds/internal/app"
	"github.com/ficusinapot/ds/internal/rest/restsvc"

	"github.com/joomcode/errorx"
)

func main() {
	if err := commands.Execute(runApp); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "%v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func runApp(cfgPath string) error {
	cfg, err := serviceconfig.LoadConfig(cfgPath)
	if err != nil {
		return errorx.InitializationFailed.Wrap(err, "load config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, *cfg, appInfo()); err != nil {
		slog.Error("application stopped with error", "error", err)
		return errorx.InternalError.Wrap(err, "run application")
	}

	return nil
}

func appInfo() restsvc.AppInfo {
	info := version.GetInfo()

	return restsvc.AppInfo{
		Name:      info.Name,
		Version:   info.Version,
		BuildTime: info.BuildTime,
		Branch:    info.Branch,
		Commit:    info.Commit,
	}
}
