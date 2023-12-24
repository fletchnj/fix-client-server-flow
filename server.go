// server.go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fix44"
	"github.com/quickfixgo/quickfix"
)

type FixServerApp struct{}

func (a FixServerApp) OnCreate(sessionID quickfix.SessionID) {
	fmt.Println("Session Created:", sessionID)
}

func (a FixServerApp) OnLogon(sessionID quickfix.SessionID) {
	fmt.Println("Logon Received:", sessionID)
}

func (a FixServerApp) OnLogout(sessionID quickfix.SessionID) {
	fmt.Println("Logout Received:", sessionID)
}

func (a FixServerApp) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a FixServerApp) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
}

func (a FixServerApp) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	return nil
}

func (a FixServerApp) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	fmt.Println("Received FIX Message:", msg.String())
	return nil
}

func main() {
	app := FixServerApp{}

	settings, err := quickfix.ParseSettingsFile("server.cfg")
	if err != nil {
		fmt.Println("Error reading server configuration:", err)
		os.Exit(1)
	}

	storeFactory := quickfix.NewMemoryStoreFactory()
	logFactory := quickfix.NewScreenLogFactory(settings)
	acceptor, err := quickfix.NewAcceptor(app, storeFactory, settings, logFactory)
	if err != nil {
		fmt.Println("Error creating acceptor:", err)
		os.Exit(1)
	}

	err = acceptor.Start()
	if err != nil {
		fmt.Println("Error starting acceptor:", err)
		os.Exit(1)
	}

	defer acceptor.Stop()

	select {}
}
