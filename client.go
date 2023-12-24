// client.go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fix44"
	"github.com/quickfixgo/quickfix"
)

type FixClientApp struct{}

func (a FixClientApp) OnCreate(sessionID quickfix.SessionID) {
	fmt.Println("Session Created:", sessionID)
}

func (a FixClientApp) OnLogon(sessionID quickfix.SessionID) {
	fmt.Println("Logon Received:", sessionID)
}

func (a FixClientApp) OnLogout(sessionID quickfix.SessionID) {
	fmt.Println("Logout Received:", sessionID)
}

func (a FixClientApp) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a FixClientApp) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
}

func (a FixClientApp) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	return nil
}

func (a FixClientApp) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	fmt.Println("Received FIX Message:", msg.String())
	return nil
}

func main() {
	app := FixClientApp{}

	settings, err := quickfix.ParseSettingsFile("client.cfg")
	if err != nil {
		fmt.Println("Error reading client configuration:", err)
		os.Exit(1)
	}

	storeFactory := quickfix.NewMemoryStoreFactory()
	logFactory := quickfix.NewScreenLogFactory(settings)
	initiator, err := quickfix.NewInitiator(app, storeFactory, settings, logFactory)
	if err != nil {
		fmt.Println("Error creating initiator:", err)
		os.Exit(1)
	}

	err = initiator.Start()
	if err != nil {
		fmt.Println("Error starting initiator:", err)
		os.Exit(1)
	}

	defer initiator.Stop()

	// Send a sample FIX message
	sendSampleFIXMessage()

	select {}
}

func sendSampleFIXMessage() {
	sessionID := quickfix.NewSessionID("FIX4.2", "CLIENT", "SERVER")

	msg := fix44.NewExecutionReport(
		field.NewOrderID("12345"),
		field.NewExecID("56789"),
		field.NewExecType(field.ExecType_NEW),
		field.NewOrdStatus(field.OrdStatus_FILLED),
		field.NewSymbol("AAPL"),
		field.NewSide(field.Side_BUY),
		field.NewOrderQty(100),
		field.NewLastShares(100),
		field.NewLastPx(150.25),
		field.NewCumQty(100),
		field.NewAvgPx(150.25),
	)

	header := msg.Header
	header.SetField(field.NewSenderCompID("CLIENT"))
	header.SetField(field.NewTargetCompID("SERVER"))
	header.SetField(field.NewMsgSeqNum(1))
	header.SetField(field.NewSendingTime(time.Now().UTC()))

	trailer := msg.Trailer
	trailer.SetField(field.NewCheckSum(quickfix.NewCheckSum(msg.String())))

	err := quickfix.SendToTarget(msg, sessionID)
	if err != nil {
		fmt.Println("Error sending FIX message:", err)
	}
}
