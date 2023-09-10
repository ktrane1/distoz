package counter

import (
	"context"
	"encoding/json"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type AddStruct struct {
	Delta int
}

/*
Node:
	each node keeps local copy of counter

	write queue?

Add:
	broadcast delta to all nodes
		- wait for broadcast to finish
		- delta or final value?
		- keep track of log of writes
	then update counter

Read:
	- guarantee is only the node that performed the write will have the latest value


*/

func Counter(n *maelstrom.Node) {
	store := maelstrom.NewSeqKV(n)

	n.Handle("init", func(msg maelstrom.Message) error {
		store.Write(context.Background(), n.ID(), 0)
		return nil
	})

	n.Handle("add", func(msg maelstrom.Message) error {
		var body AddStruct

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		current, _ := store.ReadInt(context.Background(), n.ID())
		store.Write(context.Background(), n.ID(), current+body.Delta)

		return n.Reply(msg, map[string]any{
			"type": "add_ok",
		})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		count := 0

		for _, id := range n.NodeIDs() {
			if id == n.ID() {
				value, _ := store.ReadInt(context.Background(), n.ID())
				count += value
			} else {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				res, err := n.SyncRPC(ctx, id, map[string]any{
					"type": "store_value",
				})
				defer cancel()
				if err != nil {
					continue
				}

				var body map[string]any
				if err := json.Unmarshal(res.Body, &body); err != nil {
					return err
				}

				value := int(body["value"].(float64))
				count += value
			}
		}

		return n.Reply(msg, map[string]any{
			"type":  "read_ok",
			"value": count,
		})
	})

	n.Handle("store_value", func(msg maelstrom.Message) error {
		value, err := store.ReadInt(context.Background(), n.ID())
		if err != nil {
			return err
		}

		return n.Reply(msg, map[string]any{
			"type":  "store_value_ok",
			"value": value,
		})
	})
}
