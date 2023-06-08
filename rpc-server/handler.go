package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	getRoomID(chat string) (string, error) {
		senders := strings.Split(chat, ":")
		if len(senders) != 2 {
			return "", fmt.Errorf("invalid Chat ID '%s', should be in the format of user1:user2", chat)
		}

		sort.Strings(senders)
		roomID := strings.Join(senders, ":")

		return roomID, nil
	}
	resp := rpc.NewSendResponse()
	resp.Code, resp.Messages = 0, "success"
	return resp, nil
}

func getRoomID(chat string{}) {
	
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	roomID, err := getRoomID(req.GetChat())
	if err != nil {
		return nil, err
	}

	start := req.GetCursor()
	end := start + int64(req.GetLimit())

	messages, err := rdb.GetMessagesByRoomID(ctx, roomID, start, end, req.GetReverse())
	if err != nil {
		return nil, err
	}

	respMessages := make([]*rpc.Message, 0)
	hasMore := len(messages) > int(req.GetLimit())
	nextCursor := end
	for _, msg := range messages {
		temp := &rpc.Message{
			Chat:     req.GetChat(),
			Text:     msg.Message,
			Sender:   msg.Sender,
			SendTime: msg.Timestamp,
		}
		respMessages = append(respMessages, temp)
	}

	resp := &rpc.PullResponse{
		Messages:    respMessages,
		HasMore:     &hasMore,
		NextCursor:  &nextCursor,
	}
	return resp, nil

	resp := rpc.NewPullResponse()
	resp.Code, resp.Messages = 0, "success"
	return resp, nil
}

func areYouLucky() (int32, string) {
	if rand.Int31n(2) == 1 {
		return 0, "success"
	} else {
		return 500, "oops"
	}
}
