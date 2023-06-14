package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func getRoomID(chat string) (string, error) {
	var roomID string

	lowercase := strings.ToLower(chat)
	senders := strings.Split(lowercase, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("chat room %s should have 2 users exactly, Chat ID should be in the format user1:user2", chat)
		return "", err
	}

	sender1, sender2 := senders[0], senders[1]
	// Compare the sender and receiver alphabetically, and sort them asc to form the room ID
	if comp := strings.Compare(sender1, sender2); comp == 1 {
		roomID = fmt.Sprintf("%s:%s", sender2, sender1)
	} else {
		roomID = fmt.Sprintf("%s:%s", sender1, sender2)
	}

	return roomID, nil
}

func validateSendRequest(req *rpc.SendRequest) error {
	senders := strings.Split(req.Message.Chat, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("chat room %s should have 2 users exactly, should be in the format of user1:user2", req.Message.GetChat())
		return err
	}
	sender1, sender2 := senders[0], senders[1]

	// Avoid user sending empty string for user
	if sender1 == "" || sender2 == "" {
		err := fmt.Errorf("both users %s and %s must be strings", sender1, sender2)
		return err
	}

	if req.Message.GetSender() != sender1 && req.Message.GetSender() != sender2 {
		err := fmt.Errorf("sender '%s' not in the chat room", req.Message.GetSender())
		return err
	}

	return nil
}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	if err := validateSendRequest(req); err != nil {
		return nil, err
	}

	timestamp := time.Now().Unix()
	message := &Message{
		Message:   req.Message.GetText(),
		Sender:    req.Message.GetSender(),
		Timestamp: timestamp,
	}

	roomID, _ := getRoomID(req.Message.GetChat())

	err := rdb.SaveMessage(ctx, roomID, message)
	if err != nil {
		return nil, err
	}

	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = 0, "success"
	//resp.Code, resp.Msg = areYouLucky()
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	roomID, err := getRoomID(req.GetChat())
	if err != nil {
		return nil, err
	}

	start := req.GetCursor()
	end := start + int64(req.GetLimit()) //get 1 more than required to check hasMore

	messages, err := rdb.GetMessagesByRoomID(ctx, roomID, start, end, req.GetReverse())
	if err != nil {
		return nil, err
	}

	respMessages := make([]*rpc.Message, 0)
	var counter int32 = 0
	var nextCursor int64 = 0
	hasMore := false
	for _, msg := range messages {
		if counter+1 > req.GetLimit() {
			hasMore = true
			nextCursor = end //assign nextCursor to next message after limit
			break            // do not return the message after the limit
		}
		temp := &rpc.Message{
			Chat:     req.GetChat(),
			Text:     msg.Message,
			Sender:   msg.Sender,
			SendTime: msg.Timestamp,
		}
		respMessages = append(respMessages, temp)
		counter += 1
	}

	resp := rpc.NewPullResponse()
	//resp.Code, resp.Msg = areYouLucky()
	resp.Messages = respMessages
	resp.Code = 0
	resp.Msg = "success"
	resp.HasMore = &hasMore
	resp.NextCursor = &nextCursor
	return resp, nil
}

//func areYouLucky() (int32, string) {
//	if rand.Int31n(2) == 1 {
//		return 0, "success"
//	} else {
//		return 500, "oops"
//	}
//}
