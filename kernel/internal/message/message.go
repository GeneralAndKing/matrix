package message

import (
	"context"
	"slices"
	"sync"
)

type Type uint

type Message struct {
	Type    Type
	Content string
}

const (
	WS       = "ws"
	BUSINESS = "business"
)
const (
	DEBUG Type = iota
	INFO
	WARN
	ERROR
)

var chanMap = make(map[string]*Chan)

type Chan struct {
	name      string
	internals []chan Message
	ctx       context.Context
	mutex     sync.Mutex
}

func (c *Chan) Subscribe(ctx context.Context, fn func(message Message) error) error {
	subscribeChan := make(chan Message)
	c.mutex.Lock()
	c.internals = append(c.internals, subscribeChan)
	c.mutex.Unlock()
	var err error
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			goto done
		case <-c.ctx.Done():
			err = c.ctx.Err()
			goto done
		case message := <-subscribeChan:
			if err = fn(message); err != nil {
				goto done
			}
		}
	}
done:
	c.mutex.Lock()
	slices.DeleteFunc(c.internals, func(messages chan Message) bool {
		return messages == subscribeChan
	})
	c.mutex.Unlock()
	return nil
}

func (c *Chan) Publish(message Message) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, internal := range c.internals {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case internal <- message:
		default:
			continue
		}
	}
	return nil
}

func Init(ctx context.Context) {
	chanMap[WS] = &Chan{ctx: ctx, name: WS, internals: []chan Message{}}
	chanMap[BUSINESS] = &Chan{ctx: ctx, name: BUSINESS, internals: []chan Message{}}

}

func Subscribe(ctx context.Context, namespace string, fn func(message Message) error) error {
	c, ok := chanMap[namespace]
	if !ok {
		return nil
	}
	return c.Subscribe(ctx, fn)
}

func Publish(namespace string, m Message) error {
	c, ok := chanMap[namespace]
	if !ok {
		return nil
	}
	return c.Publish(m)
}
