package main

import (
	"flag"

	"github.com/morya/utils/log"
	"github.com/morya/sms/protocol"
)

var (
	flagAccount    = flag.String("account", "acc", "account name to login remote")
	flagPassword   = flag.String("password", "", "password to login remote")
	flagRemoteAddr = flag.String("remoteaddr", "127.0.0.1:7890", "remote ISMG address")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
}

func main() {
	log.Info("OK")

	protocol.Test(*flagAccount, *flagPassword, *flagRemoteAddr)
}
