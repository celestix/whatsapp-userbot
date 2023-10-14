package ext

import (
	stdctx "context"
	"errors"
	"fmt"
	"runtime/debug"
	"sort"

	"github.com/celestix/whatsapp-userbot/ext/context"
	"github.com/celestix/whatsapp-userbot/ext/handlers"
	"github.com/celestix/whatsapp-userbot/ext/helputil"
	waLogger "github.com/celestix/whatsapp-userbot/logger"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"go.mau.fi/whatsmeow"
)

var EndGroups = errors.New("group iteration ended")
var ContinueGroups error = nil

type Dispatcher struct {
	logger        *waLogger.Logger
	handlers      map[int][]handlers.Handler
	handlerGroups []int
	help          helputil.Help
}

func NewDispatcher(logger *waLogger.Logger) *Dispatcher {
	return &Dispatcher{
		handlers: make(map[int][]handlers.Handler),
		logger:   logger.Create("dispatcher").ChangeLevel(waLogger.LevelError),
		help:     make(helputil.Help),
	}
}

func (d *Dispatcher) InitialiseProcessing(ctx stdctx.Context, client *whatsmeow.Client) {
	client.AddEventHandler(
		func(evt interface{}) {
			// Run each update in a separate goroutine
			go func(evt interface{}) {
				defer func() {
					if r := recover(); r != nil {
						// Print reason for panic + stack for some sort of helpful log output
						d.logger.ChangeLevel(waLogger.LevelCritical)
						d.logger.Println(r)
						d.logger.Println(string(debug.Stack()))
						d.logger.ChangeLevel(waLogger.LevelError)
					}
				}()
				ctx := context.New(ctx, evt)
				for group := range d.handlerGroups {
					for _, x := range d.handlers[group] {
						if !x.CheckUpdate(ctx) {
							continue
						}
						err := x.HandleUpdate(client, ctx)
						if err == nil || errors.Is(err, ContinueGroups) {
							continue
						} else if errors.Is(err, EndGroups) {
							return
						} else {
							d.logger.Println(err.Error())
						}
					}
				}
			}(evt)
		},
	)
	d.runHelper()
}

func (d *Dispatcher) runHelper() {
	helper := func(client *whatsmeow.Client, ctx *context.Context) error {
		args := ctx.Message.Args()
		if len(args) > 1 {
			hstr := d.help.GetOne(args[1])
			if hstr != "" {
				ctx.Message.Edit(
					client,
					fmt.Sprintf(
						"*%s*\n\n%s", cases.Title(language.English).String(args[1]),
						hstr,
					))
				return EndGroups
			}
		}
		ctx.Message.Edit(client, "*Help Menu*\n"+d.help.Get())
		return EndGroups
	}

	d.AddHandler(handlers.NewCommand("help", helper, d.logger.Create("help").
		ChangeLevel(waLogger.LevelInfo),
	))
}

func (d *Dispatcher) AddHandler(handler handlers.Handler) {
	d.AddHandlerToGroup(handler, 0)
}

func (d *Dispatcher) AddHandlerToGroup(handler handlers.Handler, group int) {
	d.help.Add(handler.GetName(), handler.GetDescription())
	currHandlers, ok := d.handlers[group]
	if !ok {
		d.handlerGroups = append(d.handlerGroups, group)
		sort.Ints(d.handlerGroups)
	}
	d.handlers[group] = append(currHandlers, handler)
}
