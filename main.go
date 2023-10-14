package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/celestix/whatsapp-userbot/core"
	"github.com/celestix/whatsapp-userbot/core/sql"
	"github.com/celestix/whatsapp-userbot/ext"
	"github.com/celestix/whatsapp-userbot/logger"

	"github.com/mdp/qrterminal/v3"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	dbLog := waLog.Stdout("Database", "INFO", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:waub.session?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	client := whatsmeow.NewClient(deviceStore, waLog.Noop)
	core.LOGGER.ChangeLevel(logger.LevelInfo)
	core.LOGGER.Println("Created new client")
	dispatcher := ext.NewDispatcher(core.LOGGER)
	core.LOGGER.Println("Created new disptacher")
	dispatcher.InitialiseProcessing(ctx, client)
	db := sql.LoadDB(core.LOGGER)
	core.Load(dispatcher)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(ctx)
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}
	core.LOGGER.ChangeLevel(logger.LevelMain)
	core.LOGGER.Println("Whatsapp Userbot", "has been started...")
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1) // buffered channel
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	//_ = client.SendPresence(types.PresenceUnavailable)
	client.Disconnect()
	if err := db.Close(); err != nil {
		core.LOGGER.ChangeLevel(logger.LevelError).Panicln("failed to close db:", err.Error())
	}
}
