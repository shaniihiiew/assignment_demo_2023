package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	import (
    client := redis.NewClient(&redis.Options{
        Addr:     "your-redis-address:6379",
        Password: "your-redis-password", 
        DB:       0,                     
    })

    pong, err := client.Ping().Result()
    if err != nil {
        return err
    }
    fmt.Println(pong) 

    err = client.Set("your-key", "your-value", 0).Err()
    if err != nil {
        return err
    }

    // Request sent successfully
    fmt.Println("Send request executed successfully")
    return nil
	}
	resp := rpc.NewSendResponse()
	resp.Code, resp.Messages = 0, "success"
	return resp, nil
}
	
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {

func pullMessages(start, end int64) ([]string, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     "your-redis-address:6379",
        Password: "your-redis-password", 
        DB:       0,                    
    })
    result, err := client.ZRange("your-messages-key", start, end).Result()
    if err != nil {
        return nil, err
    }
    return result, nil
}

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
