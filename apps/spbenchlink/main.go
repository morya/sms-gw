package main

import (
	"flag"

	options "github.com/morya/sms-gw/options"
	session "github.com/morya/sms-gw/session"
	log "github.com/sirupsen/logrus"
)

var (
	flagAccount    = flag.String("account", "sp", "account name")
	flagPassword   = flag.String("password", "", "account password")
	flagRemoteAddr = flag.String("remoteaddr", "127.0.0.1:7890", "socket address to remote")
	flagConcurrent = flag.Int("concurrent", 100, "concurrent connections to make")

	flagLogLevel = flag.String("loglevel", "INFO", "set log level")
)

type App struct {
	opt *options.Options
}

func (a *App) Main() {
	c, err := session.NewConn(a.opt, "cmpp")
	if err != nil {
		log.Error(err)
		return
	}
	c.Client()
}

func main() {
	var opt = &options.Options{
		Account:         *flagAccount,
		Password:        *flagPassword,
		RemoteAddr:      *flagRemoteAddr,
		ConcurrentCount: *flagConcurrent,
	}
	var app = &App{
		opt: opt,
	}

	app.Main()
}
