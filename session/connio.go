package session

import (
	"time"

	"github.com/morya/sms/protocol/unified"
	"github.com/morya/utils/log"
	"golang.org/x/net/context"
)

func (c *Conn) Loop() {
	ctx, cancel := context.WithCancel(context.Background())

	switch c.tcpMode {
	case CONN_MODE_CLIENT:
		go c.DoLogin(ctx)
	}

	go c.RecvLoop(ctx, cancel)
	go c.SendLoop(ctx, cancel)

	go func() {
		select {
		case <-ctx.Done():
			func() {
				defer func() {
					recover()
				}()
				c.sock.Close()
			}()
		}
	}()
}

func (c *Conn) DoLogin(ctx context.Context) {
	log.Debugf("before send login msg")
	time.Sleep(time.Millisecond * 100)

	var bind = &unified.MsgBindReq{Account: c.opt.Account, Password: c.opt.Password}
	data, err := c.coder.Encode(bind)
	if err != nil {
		ctx.Done()
		log.Error(err)
		return
	}

	log.Debug("send login msg to remote")
	c.Send(data)
	c.HeartBeatLoop(ctx)
}

func (c *Conn) HeartBeatLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 4)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return

		case <-ticker.C:
			var hb = &unified.MsgHeartBeat{}
			data, err := c.coder.Encode(hb)
			if err != nil {
				ctx.Done()
				log.Errorf("encode heartbeat msg failed %v", err)
				return
			}

			c.Send(data)

			log.Debug("heartbeat sent")
		}
	}
}

func (c *Conn) RecvLoop(ctx context.Context, cancel context.CancelFunc) {
	for {
		data, err := c.coder.NextMsg(c.reader)
		if err != nil {
			cancel()
			log.ErrorError(err)
			return
		}

		msg, err := c.coder.Decode(data)
		if err != nil {
			cancel()
			log.Errorf("decode remote msg failed, %v", err)
			return
		}
		log.Debugf("got msg %#v", msg)
	}
}

func (c *Conn) SendLoop(ctx context.Context, cancel context.CancelFunc) {
	for {
		select {
		case data := <-c.chanMsgSend:
			n, err := c.writer.Write(data)

			if err != nil {
				cancel()
				log.InfoError(err, "send data failed")
				return
			}

			if len(data) != n {
				log.Infof("sent data too short, sent length = %v, real length = %v", n, len(data))
			}
			c.writer.Flush()

		case <-ctx.Done():
			return

		}
	}
}
