package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hackfeed/stark/internal/cache/redis"
	"github.com/hackfeed/stark/internal/domain/chatmessage"
	"github.com/hackfeed/stark/internal/store/chatmessagerepo"
	"github.com/hackfeed/stark/internal/store/identifier"
)

func main() {
	ctx := context.Background()

	redisClient, err := redis.NewRedisClient(ctx, &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
	}
	chatMessagesRepo := chatmessagerepo.NewRedisRepo(redisClient, 1*time.Hour)

	_ = redisClient.Publish(ctx, "global", "Noname joined")
	channel := redisClient.Subscribe(ctx, "global")
	go func(ch <-chan string) {
		for msg := range ch {
			id := identifier.New().AddKeyValue("room", "global")
			messages, err := chatMessagesRepo.GetMessages(id)
			if err != nil {
				fmt.Println(err)
			}
			messages = append(messages, chatmessage.New("Serghik", msg))
			err = chatMessagesRepo.SetMessages(id, messages)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(msg)
		}
	}(channel)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		switch {
		case strings.HasPrefix(text, "/join"):
			fmt.Println("Joined")
			// chat := strings.Split(text, " ")[1]
			// fmt.Printf("Joined to %s channel\n", chat)
			// err = redisClient.Publish(ctx, chat, "Noname joined")
			// newChan := redisClient.Subscribe(ctx, chat)
		}
	}
}
