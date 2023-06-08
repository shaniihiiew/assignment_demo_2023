package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
   	"github.com/go-redis/redis"
)

type IMServiceImpl struct {
    redisClient *redis.Client
}

func NewIMServiceImpl(redisClient *redis.Client) *IMServiceImpl {
    return &IMServiceImpl{
        redisClient: redisClient,
    }
}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {

    message := &Message{
        Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
    }

    roomID := getRoomID(req.GetChat())
    if roomID == "" {
        return nil, fmt.Errorf("invalid Chat ID '%s', should be in the format of user1:user2", req.GetChat())
    }

    err := s.SaveMessage(ctx, roomID, message)
    if err != nil {
        return nil, err
    }

    resp := &rpc.SendResponse{
        Code: 0, 
        Msg:  "success",
    }

    return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
    roomID := getRoomID(req.GetChat())
    if roomID == "" {
        return nil, fmt.Errorf("invalid Chat ID '%s', should be in the format of user1:user2", req.GetChat())
    }

    start := req.GetCursor()
    end := start + int64(req.GetLimit())

    messages, err := s.GetMessagesByRoomID(ctx, roomID, start, end, req.GetReverse())
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
            nextCursor = end
            break
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

    resp := &rpc.PullResponse{
        Messages:    respMessages,
        Code:        0,
        Msg:         "success",
        HasMore:     &hasMore,
        NextCursor:  &nextCursor,
    }

    return resp, nil
}

func (s *IMServiceImpl) SaveMessage(ctx context.Context, roomID string, message *Message) error {
    text, err := json.Marshal(message)
    if err != nil {
        return err
    }

    member := &redis.Z{
        Score:  float64(message.Timestamp),
        Member: text,
    }

    _, err = s.redisClient.ZAdd(ctx, roomID, member).Result()
    if err != nil {
        return err
    }

    return nil
}

func (s *IMServiceImpl) GetMessagesByRoomID(ctx context.Context, roomID string, start, end int64, reverse bool) ([]*Message, error) {
    var rawMessages []string

    if reverse {
        rawMessages, err := s.redisClient.ZRevRange(ctx, roomID, start, end).Result()
        if err != nil

