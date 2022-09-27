package chats

import (
	"context"
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/slntopp/nocloud-cc/pkg/chats/proto"
)

type UI struct {
	*gocui.Gui
	ctx    context.Context
	stream proto.ChatService_StreamClient
	client proto.ChatServiceClient
	chat   string
}

func NewUI(ctx context.Context,
	client proto.ChatServiceClient, stream proto.ChatService_StreamClient, chat string) (*UI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	ui := &UI{Gui: g, stream: stream, client: client, chat: chat, ctx: ctx}

	return ui, nil
}

func (ui *UI) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true

	if messages, err := g.SetView("messages", 0, 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Title = "messages"
		messages.Autoscroll = true
		messages.Wrap = true
	}

	if input, err := g.SetView("input", 0, maxY-5, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		input.Title = "send"
		input.Autoscroll = false
		input.Wrap = true
		input.Editable = true
	}

	g.SetCurrentView("input")

	return nil
}

func (ui *UI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *UI) sendMsg(g *gocui.Gui, v *gocui.View) error {
	if len(v.Buffer()) == 0 {
		v.SetCursor(0, 0)
		v.Clear()
		return nil
	}

	_, err := ui.client.SendChatMessage(ui.ctx, &proto.SendChatMessageRequest{
		Message: &proto.ChatMessage{
			Message: v.Buffer(),
			To:      ui.chat,
		},
	})
	if err != nil {
		return err
	}

	v.SetCursor(0, 0)
	v.Clear()
	return nil
}

func (ui *UI) receiveMsg() {
	for {
		message, err := ui.stream.Recv()
		if err != nil {
			return
		}

		ui.Update(func(g *gocui.Gui) error {
			view, _ := ui.View("messages")
			fmt.Fprintf(view, "%s: %s", message.From, message.Message)
			return nil
		})
	}
}
