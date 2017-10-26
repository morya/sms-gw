package main

import (
	"flag"

	"github.com/morya/sms/session"
	"github.com/morya/utils/log"
)

var (
	flagAccount    = flag.String("account", "sp", "account name")
	flagPassword   = flag.String("password", "", "account password")
	flagRemoteAddr = flag.String("remoteaddr", "127.0.0.1:7890", "socket address to remote")
)

func init() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	c, err := session.NewConn(*flagAccount, *flagPassword, *flagRemoteAddr)
	if err != nil {
		log.Error(err)
		return
	}
	c.Client()
}
