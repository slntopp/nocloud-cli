package chats

import (
	"context"
	"fmt"

	"github.com/jroimartin/gocui"
	proto "github.com/slntopp/nocloud-proto/cc"
	regpb "github.com/slntopp/nocloud-proto/registry"
	"github.com/slntopp/nocloud-proto/registry/accounts"
)

type UI struct {
	*gocui.Gui
	ctx       context.Context
	ctxAcc    context.Context
	stream    proto.ChatService_StreamClient
	client    proto.ChatServiceClient
	clientAcc regpb.AccountsServiceClient
	chat      string
}

func NewUI(
	ctx context.Context,
	ctxAcc context.Context,
	client proto.ChatServiceClient,
	clientAcc regpb.AccountsServiceClient,
	stream proto.ChatService_StreamClient, chat string) (*UI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	ui := &UI{Gui: g, stream: stream, client: client, chat: chat, ctx: ctx, ctxAcc: ctxAcc, clientAcc: clientAcc}

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

			from := "anon"
			resp, err := ui.clientAcc.Get(ui.ctxAcc, &accounts.GetRequest{Uuid: message.From})
			if err == nil {
				from = resp.GetTitle()
			}

			fmt.Fprintf(view, "%s: %s", from, message.Message)
			return nil
		})
	}
}
