package job

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type AvailableDate struct {
	Date        string `json:"date"`
	BusinessDay bool   `json:"business_day"`
}

type AvailableTime struct {
	AvailableTimes []string `json:"available_times"`
	BusinessTimes  []string `json:"business_times"`
}

type Listener interface {
	Start()
	Stop()
}

type RespListener struct {
	Url       string
	Ctx       context.Context
	Cancel    context.CancelFunc
	ReqIDChan chan network.RequestID
}

func RespListenerInit(url string, c context.Context) *RespListener {
	ch := make(chan network.RequestID, 1)
	ctx, cancel := context.WithCancel(c)
	return &RespListener{
		Url:       url,
		Ctx:       ctx,
		Cancel:    cancel,
		ReqIDChan: ch,
	}
}

func (l *RespListener) Start() {
	chromedp.ListenTarget(l.Ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			resp := ev.Response
			if resp.URL == l.Url {
				fmt.Println("response received, ", ev.RequestID)
				l.ReqIDChan <- ev.RequestID
			}
		}
	})

}

func (l *RespListener) Stop() {
	l.Cancel()
	close(l.ReqIDChan)
}
