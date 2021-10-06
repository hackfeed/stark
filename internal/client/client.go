package client

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hackfeed/stark/internal/db/cache"
	"github.com/hackfeed/stark/internal/domain"
	usersrepo "github.com/hackfeed/stark/internal/store/users_repo"
	"github.com/jroimartin/gocui"
)

var (
	cacheClient *cache.RedisClient
	usersRepo   usersrepo.UsersRepository
	ctx         context.Context
	activeChat  *domain.Chat
	user        *domain.User
)

func init() {
	var err error

	ctx = context.Background()

	cacheClient, err = cache.NewRedisClient(ctx, &cache.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err != nil {
		log.Fatalln(err)
	}

	usersRepo = usersrepo.NewRedisRepo(*cacheClient, 1*time.Hour)

	activeChat = domain.NewChat("global")
}

func Disconnect(g *gocui.Gui, v *gocui.View) error {
	users, err := usersRepo.GetUsers(activeChat.GetName())
	if err != nil {
		log.Fatalln(err)
	}
	delete(users, user.GetName())
	err = usersRepo.SetUsers(activeChat.GetName(), users)
	if err != nil {
		log.Fatalln(err)
	}
	cacheClient.Publish(ctx, activeChat.GetName(), "/users>")
	cacheClient.Publish(ctx, activeChat.GetName(), fmt.Sprintf("%s just disconnected :(", user.GetName()))

	return gocui.ErrQuit
}

func Send(g *gocui.Gui, v *gocui.View) error {
	message := domain.NewMessage(user.GetName(), strings.TrimSpace(v.Buffer()))

	cacheClient.Publish(ctx, activeChat.GetName(), message.String())

	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
		return nil
	})

	return nil
}

func Connect(g *gocui.Gui, v *gocui.View) error {
	user = domain.NewUser(strings.TrimSpace(v.Buffer()))
	user.AddChat(*activeChat)
	activeChat.SetMessages(cacheClient.Subscribe(ctx, activeChat.GetName()))

	users, err := usersRepo.GetUsers(activeChat.GetName())
	if err != nil {
		log.Fatalln(err)
	}
	if users == nil {
		users = make(map[string]struct{})
	}
	users[user.GetName()] = struct{}{}
	err = usersRepo.SetUsers(activeChat.GetName(), users)
	if err != nil {
		log.Fatalln(err)
	}

	cacheClient.Publish(ctx, activeChat.GetName(), "/users>")
	cacheClient.Publish(ctx, activeChat.GetName(), fmt.Sprintf("%s just joined!", user.GetName()))

	g.SetViewOnTop("messages")
	g.SetViewOnTop("users")
	g.SetViewOnTop("input")
	g.SetCurrentView("input")

	messagesView, _ := g.View("messages")
	usersView, _ := g.View("users")

	go func() {
		for msg := range activeChat.GetMessages() {
			switch {
			case strings.HasPrefix(msg, "/users>"):
				users, err := usersRepo.GetUsers(activeChat.GetName())
				if err != nil {
					log.Fatalln(err)
				}
				chatUsers := ""
				chatUsersCount := len(users)
				for user := range users {
					chatUsers += user + "\n"
				}
				g.Update(func(g *gocui.Gui) error {
					usersView.Title = fmt.Sprintf("%d users:", chatUsersCount)
					usersView.Clear()
					fmt.Fprintln(usersView, chatUsers)
					return nil
				})
			case strings.Contains(msg, "/join"):
				newChatName := strings.TrimSpace(strings.SplitAfter(msg, "/join")[1])
				activeChat = domain.NewChat(newChatName)
				activeChat.SetMessages(cacheClient.Subscribe(ctx, activeChat.GetName()))
				user.AddChat(*activeChat)

				users, err := usersRepo.GetUsers(activeChat.GetName())
				if err != nil {
					log.Fatalln(err)
				}
				if users == nil {
					users = make(map[string]struct{})
				}
				users[user.GetName()] = struct{}{}
				err = usersRepo.SetUsers(activeChat.GetName(), users)
				if err != nil {
					log.Fatalln(err)
				}

				chatUsers := ""
				chatUsersCount := len(users)
				for user := range users {
					chatUsers += user + "\n"
				}

				cacheClient.Publish(ctx, activeChat.GetName(), "/users>")
				cacheClient.Publish(ctx, activeChat.GetName(), fmt.Sprintf("%s just joined!", user.GetName()))

				g.Update(func(g *gocui.Gui) error {
					usersView.Title = fmt.Sprintf("%d users:", chatUsersCount)
					usersView.Clear()
					messagesView.Clear()
					fmt.Fprintln(usersView, chatUsers)
					return nil
				})
			default:
				g.Update(func(g *gocui.Gui) error {
					fmt.Fprintln(messagesView, msg)
					return nil
				})
			}
		}
	}()
	return nil
}
