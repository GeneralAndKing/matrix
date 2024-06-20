package message

import "context"

const WS = "ws"

var chanMap map[string]*Chan = make(map[string]*Chan)

type Chan struct {
	name     string
	internal chan Message
	ctx      context.Context
}

func (c *Chan) Subscribe(ctx context.Context, fn func(message Message) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.ctx.Done():
			return c.ctx.Err()
		case message := <-c.internal:
			if err := fn(message); err != nil {
				return err
			}
		}
	}
}

func (c *Chan) Publish(message Message) error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case c.internal <- message:
		return nil
	}
}

type Message struct {
	Type    uint
	Content string
}

func InitAndRegister(ctx context.Context, name string) *Chan {
	c := &Chan{ctx: ctx, name: name, internal: make(chan Message)}
	chanMap[name] = c
	return c
}
func Fetch(name string) *Chan {
	return chanMap[name]
}
