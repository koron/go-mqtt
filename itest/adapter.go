package itest

import (
	"sync"

	"github.com/koron/go-mqtt/mqtopic"
	"github.com/koron/go-mqtt/packet"
	"github.com/koron/go-mqtt/server"
)

type onPingFunc func() (bool, error)

// Adapter provides simple MQTT server behavior
type Adapter struct {
	mu  sync.Mutex
	cas map[string]*clientAdapter

	onPing onPingFunc
}

// Connect is called when new client is connected.
func (a *Adapter) Connect(srv *server.Server, c server.Client, p *packet.Connect) (server.ClientAdapter, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	ca := &clientAdapter{
		id: p.ClientID,
		a:  a,
		c:  c,
	}
	if a.cas == nil {
		a.cas = make(map[string]*clientAdapter)
	}
	a.cas[ca.id] = ca
	return ca, nil
}

// Disconnect is called when a known client is disconnected some reason.
func (a *Adapter) Disconnect(srv *server.Server, ca server.ClientAdapter, err error) {
	ca2, ok := ca.(*clientAdapter)
	if !ok {
		return
	}
	a.mu.Lock()
	delete(a.cas, ca2.id)
	a.mu.Unlock()
}

func (a *Adapter) dispatch(src *clientAdapter, m *server.Message) {
	topic, err := mqtopic.Parse(m.Topic)
	if err != nil {
		return
	}
	a.mu.Lock()
	for _, dst := range a.cas {
		if dst == src {
			continue
		}
		go dst.dispatch(topic, m)
	}
	a.mu.Unlock()
}

var _ server.Adapter = (*Adapter)(nil)

type clientAdapter struct {
	id string
	a  *Adapter
	c  server.Client
	fm map[string]mqtopic.Filter
}

func (ca *clientAdapter) ID() string {
	return ca.id
}

func (ca *clientAdapter) IsSessionPresent() bool {
	return false
}

func (ca *clientAdapter) OnDisconnect() error {
	return nil
}

func (ca *clientAdapter) OnPing() (bool, error) {
	if ca.a.onPing != nil {
		return ca.a.onPing()
	}
	return true, nil
}

func (ca *clientAdapter) OnSubscribe(topics []server.Topic) ([]server.QoS, error) {
	q := make([]server.QoS, len(topics))
	if len(topics) > 0 && ca.fm == nil {
		ca.fm = make(map[string]mqtopic.Filter)
	}
	for i, topic := range topics {
		f, err := mqtopic.ParseFilter(topic.Filter)
		if err != nil {
			q[i] = server.Failure
			continue
		}
		q[i] = server.AtMostOnce
		ca.fm[topic.Filter] = f
	}
	return q, nil
}

func (ca *clientAdapter) OnUnsubscribe(filters []string) error {
	if len(ca.fm) == 0 || len(filters) == 0 {
		return nil
	}
	for _, f := range filters {
		delete(ca.fm, f)
	}
	return nil
}

func (ca *clientAdapter) OnPublish(m *server.Message) error {
	ca.a.dispatch(ca, m)
	return nil
}

func (ca *clientAdapter) dispatch(topic mqtopic.Topic, m *server.Message) {
	for _, f := range ca.fm {
		if f.Match(topic) {
			_ = ca.c.Publish(m.QoS, m.Retain, m.Topic, m.Body)
			return
		}
	}
}

var _ server.ClientAdapter = (*clientAdapter)(nil)
