package session

import (
	"context"
	"time"

	"github.com/morya/sms-gw/protocol/unified"
	"github.com/sirupsen/logrus"
)

func (c *Conn) Loop() error {
	ctx, cancel := context.WithCancel(context.Background())

	switch c.tcpMode {
	case CONN_MODE_CLIENT:
		go c.DoLogin(ctx)
	}

	go c.RecvLoop(ctx, cancel)
	go c.SendLoop(ctx, cancel)

	func() {
		func() {
			defer func() {
				recover()
			}()
			c.sock.Close()
		}()

		<-ctx.Done()
	}()

	return c.lastErr
}

func (c *Conn) DoLogin(ctx context.Context) {
	logrus.Debugf("before send login msg")
	time.Sleep(time.Millisecond * 100)

	var bind = &unified.MsgBindReq{Account: c.opt.Account, Password: c.opt.Password}
	data, err := c.coder.Encode(bind)
	if err != nil {
		ctx.Done()
		logrus.Error(err)
		return
	}

	logrus.Debug("send login msg to remote")
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
				logrus.Errorf("encode heartbeat msg failed %v", err)
				return
			}

			c.Send(data)

			logrus.Debug("heartbeat sent")
		}
	}
}

func (c *Conn) RecvLoop(ctx context.Context, cancel context.CancelFunc) {
	for {
		data, err := c.coder.NextMsg(c.reader)
		if err != nil {
			cancel()
			logrus.Error(err)
			return
		}

		msg, err := c.coder.Decode(data)
		if err != nil {
			cancel()
			logrus.Errorf("decode remote msg failed, %v", err)
			return
		}
		logrus.Debugf("got msg %#v", msg)
	}
}

func (c *Conn) SendLoop(ctx context.Context, cancel context.CancelFunc) {
	for {
		select {
		case data := <-c.chanMsgSend:
			n, err := c.writer.Write(data)

			if err != nil {
				cancel()
				logrus.Info(err, "send data failed")
				return
			}

			if len(data) != n {
				logrus.Infof("sent data too short, sent length = %v, real length = %v", n, len(data))
			}
			c.writer.Flush()

		case <-ctx.Done():
			return

		}
	}
}
