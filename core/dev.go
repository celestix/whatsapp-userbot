package core

import (
	"fmt"
	"log"
	"strings"

	"github.com/celestix/whatsapp-userbot/ext"
	"github.com/celestix/whatsapp-userbot/ext/context"
	"github.com/celestix/whatsapp-userbot/ext/handlers"
	waLogger "github.com/celestix/whatsapp-userbot/logger"

	piston "github.com/milindmadhukar/go-piston"

	"go.mau.fi/whatsmeow"
)

var (
	pClient *piston.Client

	languages      *[]string
	availableLangs string
)

func initLanguagesString() {
	availableLangs = "Here is the list of available languages:"
	for _, lang := range *languages {
		availableLangs += fmt.Sprintf("\n- ```%s```", lang)
	}
	availableLangs += "\n\n*Usage*: ```.eval python print(\"hello\")```"
}

func languagesC(client *whatsmeow.Client, ctx *context.Context) error {
	_, _ = ctx.Message.Edit(client, availableLangs)
	return ext.EndGroups
}

func eval(client *whatsmeow.Client, ctx *context.Context) error {
	args := ctx.Message.ArgsN(3)
	if len(args) < 3 {
		_, _ = ctx.Message.Edit(client, "Invalid amount of arguments.\nExpected: ```.eval <language> <code>```")
		return ext.EndGroups
	}
	lang := args[1]
	text := args[2]

	if text == "" {
		_, _ = ctx.Message.Edit(client, "You need to provide me some code to eval.")
		return ext.EndGroups
	}

	output, err := pClient.Execute(strings.ToLower(lang), "",
		[]piston.Code{
			{Content: text},
		},
	)
	if err != nil {
		_, _ = ctx.Message.Edit(client, fmt.Sprintf("failed to eval: %s", err.Error()))
		return ext.EndGroups
	}

	out := output.GetOutput()

	if out == "" {
		out = "No Output"
	}

	reply := fmt.Sprintf("*Language*: ```%s```\n\n*Input*: ```%s```\n\n*Output*: ```%s```", lang, text, output.GetOutput())

	_, err = ctx.Message.Edit(client, reply)
	if err != nil {
		log.Println("failed to send message:", err.Error())
	}
	return ext.EndGroups
}

func (*Module) LoadDev(dispatcher *ext.Dispatcher) {
	ppLogger := LOGGER.Create("dev")
	defer ppLogger.Println("Loaded Dev module")
	pClient = piston.CreateDefaultClient()
	languages = pClient.GetLanguages()
	initLanguagesString()
	dispatcher.AddHandler(
		handlers.NewCommand("eval", authorizedOnly(eval), ppLogger.Create("eval-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription(`Execute codes using piston engine`),
	)
	dispatcher.AddHandler(
		handlers.NewCommand("langs", authorizedOnly(languagesC), ppLogger.Create("langs-cmd").
			ChangeLevel(waLogger.LevelInfo),
		).AddDescription(`Get list of supported languages`),
	)
}
