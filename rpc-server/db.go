package main

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	cli *redis.Client
}

// cap first letter to make method/fn/field public
type Message struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (c *RedisClient) InitClient(ctx context.Context, address, password string) error {
	r := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	// test connection
	//Ping from commands.go
	if err := r.Ping(ctx).Err(); err != nil {
		return err
	}

	c.cli = r
	return nil
}

func (c *RedisClient) SaveMessage(ctx context.Context, roomID string, message *Message) error {
	// Marshal the Go struct into JSON bytes
	text, err := json.Marshal(message)
	if err != nil {
		return err
	}

	//redis Z represents sorted set member.
	member := &redis.Z{
		Score:  float64(message.Timestamp), // The sort key
		Member: text,                       // Data
	}

	//roomID used as score
	_, err = c.cli.ZAdd(ctx, roomID, *member).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) GetMessagesByRoomID(ctx context.Context, roomID string, start, end int64, reverse bool) ([]*Message, error) {
	var (
		dbMessages []string
		messages   []*Message
		err        error
	)

	if reverse {
		// Desc order with time -> first message is the latest message
		dbMessages, err = c.cli.ZRevRange(ctx, roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	} else {
		// Asc order with time -> first message is the earliest message
		dbMessages, err = c.cli.ZRange(ctx, roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	}

	for _, msg := range dbMessages {
		temp := &Message{}
		//prase json data msg and store in temp
		err := json.Unmarshal([]byte(msg), temp)
		if err != nil {
			return nil, err
		}
		//append temp to messages
		messages = append(messages, temp)
	}

	return messages, nil
}
