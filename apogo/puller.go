package apogo

import (
	"context"
	"encoding/json"
	"github.com/Bishoptylaor/go-toolkit/xnet/xhttp"
	"sync/atomic"
	"time"
)

// this is a static check
var _ puller = (*longPuller)(nil)

// puller fetch confi updates
type puller interface {
	// start poll updates
	start()
	// preload fetch all config to local cache, and update all notifications
	preload() error
	// stop poll updates
	stop()
	// addNamespaces add new namespace and pump config data
	addNamespaces(namespaces ...string) error
}

// notificationHandler handle namespace update notification
type notificationHandler func(namespace string) error

// longPuller implement puller interface
type longPuller struct {
	conf *Conf

	pullerInterval time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
	version        uint64
	hclient        xhttp.HttpClientWrapper

	notifications *notificationRepo
	handler       notificationHandler
}

// newLongPuller create a Puller
func newLongPuller(conf *Conf, interval time.Duration, handler notificationHandler) puller {
	puller := &longPuller{
		conf:           conf,
		pullerInterval: interval,
		hclient:        NewHclient(longPollTimeout),
		notifications:  new(notificationRepo),
		handler:        handler,
	}

	puller.ctx, puller.cancel = context.WithCancel(context.Background())

	for _, namespace := range conf.NameSpaceNames {
		puller.notifications.setNotificationID(namespace, defaultNotificationID)
	}

	return puller
}

func (p *longPuller) start() {
	go p.watchUpdates()
}

func (p *longPuller) preload() error {
	return p.pumpUpdates()
}

// addNamespaces subscribe to new namespaces and pull all config data to local
func (p *longPuller) addNamespaces(namespaces ...string) error {
	var update bool
	for _, namespace := range namespaces {
		if p.notifications.addNotificationID(namespace, defaultNotificationID) {
			update = true
		}
	}
	if update {
		return p.pumpUpdates()
	}
	return nil
}

func (p *longPuller) watchUpdates() {
	timer := time.NewTimer(p.pullerInterval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := p.pumpUpdates(); err != nil {
				defaultLogger.Printf("module:agollo method:watchUpdates err:%v", err)
			}
			timer.Reset(p.pullerInterval)

		case <-p.ctx.Done():
			return
		}
	}
}

func (p *longPuller) stop() {
	p.cancel()
}

func (p *longPuller) updateNotificationConf(notification *notification) {
	p.notifications.setNotificationID(notification.NamespaceName, notification.NotificationID)
}

// pumpUpdates fetch updated namespace, handle updated namespace then update notification id
func (p *longPuller) pumpUpdates() error {
	// serialize pumpUpdates request

	version := atomic.AddUint64(&p.version, 1)

	var ret error

	updates, err := p.poll()
	if err != nil {
		return err
	}

	if atomic.LoadUint64(&p.version) != version {
		return nil
	}

	for _, update := range updates {
		if err := p.handler(update.NamespaceName); err != nil {
			ret = err
			continue
		}
		p.updateNotificationConf(update)
	}

	return ret
}

// poll until a update or timeout
func (p *longPuller) poll() ([]*notification, error) {
	notifications := p.notifications.toString()
	url := notificationURL(p.conf, notifications)
	defaultLogger.Printf("module:agollo method:longPuller.poll url:%s start", url)
	bts, err := Request(p.hclient, url)
	defaultLogger.Printf("module:agollo method:longPuller.poll url:%s finish with data:%s err:%v", url, bts, err)
	if err != nil || len(bts) == 0 {
		return nil, err
	}
	var ret []*notification
	if err := json.Unmarshal(bts, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
