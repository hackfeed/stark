package client

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hackfeed/stark/internal/db/cache"
	"github.com/hackfeed/stark/internal/domain"
	"github.com/hackfeed/stark/internal/store/filesrepo"
	"github.com/hackfeed/stark/internal/store/usersrepo"
	"github.com/jroimartin/gocui"
	"github.com/logrusorgru/aurora/v3"
)

func addChatUser(user *domain.User, chat *domain.Chat, repo usersrepo.UsersRepository) (map[string]struct{}, error) {
	users, err := repo.GetUsers(chat.GetName())
	if err != nil {
		return nil, err
	}

	if users == nil {
		users = make(map[string]struct{})
	}
	users[user.GetName()] = struct{}{}

	err = repo.SetUsers(chat.GetName(), users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func removeChatUser(user *domain.User, chat *domain.Chat, repo usersrepo.UsersRepository) (map[string]struct{}, error) {
	users, err := repo.GetUsers(chat.GetName())
	if err != nil {
		return nil, err
	}

	delete(users, user.GetName())

	err = repo.SetUsers(chat.GetName(), users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func addChatFile(file *domain.File, chat *domain.Chat, repo filesrepo.FilesRepository) (map[string][]byte, error) {
	files, err := repo.GetFiles(chat.GetName())
	if err != nil {
		return nil, err
	}

	if files == nil {
		files = make(map[string][]byte)
	}
	files[file.GetName()] = file.GetContent()

	err = repo.SetFiles(chat.GetName(), files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func updateChatsView(g *gocui.Gui, v *gocui.View, user *domain.User, toClear bool) {
	g.Update(func(g *gocui.Gui) error {
		chats := ""
		for chatName, chat := range user.GetChats() {
			counter := ""
			count := len(chat.GetBuffer())
			if count > 0 {
				counter = fmt.Sprintf(" (%d)", count)
			}
			if chat.GetIsActive() {
				chats += aurora.Sprintf(aurora.Green("#%s%s\n"), chatName, aurora.Magenta(counter))
			} else {
				chats += aurora.Sprintf("#%s%s\n", chatName, aurora.Magenta(counter))
			}
		}

		if toClear {
			v.Clear()
		}
		fmt.Fprintln(v, chats)

		return nil
	})
}

func updateMessagesView(g *gocui.Gui, v *gocui.View, message string, toClear bool) {
	g.Update(func(g *gocui.Gui) error {
		if toClear {
			v.Clear()
		}
		fmt.Fprintln(v, message)

		return nil
	})
}

func updateMessagesViewWithBuffer(g *gocui.Gui, v *gocui.View, chat *domain.Chat, toClear bool) {
	g.Update(func(g *gocui.Gui) error {
		if toClear {
			v.Clear()
		}
		for _, bufferMessage := range chat.GetBuffer() {
			fmt.Fprintln(v, bufferMessage)
		}
		chat.SetBuffer(nil)

		return nil
	})
}

func updateUsersView(g *gocui.Gui, v *gocui.View, chat *domain.Chat, repo usersrepo.UsersRepository, toClear bool) {
	g.Update(func(g *gocui.Gui) error {
		users, err := repo.GetUsers(chat.GetName())
		if err != nil {
			log.Fatalln(err)
		}
		chatUsers := ""
		chatUsersCount := len(users)
		for u := range users {
			chatUsers += u + "\n"
		}

		v.Title = fmt.Sprintf("%d users:", chatUsersCount)
		if toClear {
			v.Clear()
		}
		fmt.Fprintln(v, chatUsers)

		return nil
	})
}

func clearView(g *gocui.Gui, v *gocui.View) {
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		return nil
	})
}

func getActiveChat(ctx context.Context, cc *cache.RedisClient, name string, user *domain.User, isNew bool) *domain.Chat {
	var chat *domain.Chat
	if isNew {
		chat = domain.NewChat(name)
		user.AddChat(chat)
	} else {
		user.SetActiveChat(name)
		chat = user.GetActiveChat()
	}
	chat.SetMessages(cc.Subscribe(ctx, name))

	return chat
}

func getUserFromMessage(message string) string {
	return strings.TrimRight(strings.Split(message, " ")[1], ":")
}

func getCommandFromMessage(message string) string {
	return strings.Split(message, " ")[2]
}

func handleJoin(ctx context.Context, cc *cache.RedisClient, chat *domain.Chat, user *domain.User) {
	cc.Publish(ctx, chat.GetName(), "/users>")
	cc.Publish(ctx, chat.GetName(), aurora.Sprintf(aurora.Green("%s just joined!"), aurora.Yellow(user.GetName())))
}

func handleLeave(ctx context.Context, cc *cache.RedisClient, chat *domain.Chat, user *domain.User, message string) {
	cc.Publish(ctx, chat.GetName(), "/users>")
	cc.Publish(ctx, chat.GetName(), aurora.Sprintf(aurora.Red("%s just %s :("), aurora.Yellow(user.GetName()), message))
}

func updateMessagesCounter(g *gocui.Gui, v *gocui.View, user *domain.User) {
	for _, chat := range user.GetChats() {
		go func(c *domain.Chat) {
			if c.GetIsActive() {
				return
			}
			for msg := range c.GetMessages() {
				if c.GetIsActive() {
					return
				}
				msgSplit := strings.Split(msg, " ")
				if strings.HasPrefix(msg, "/") || strings.HasPrefix(msgSplit[2], "/") {
					continue
				}
				c.SetBuffer(append(c.GetBuffer(), msg))
				updateChatsView(g, v, user, true)
			}
		}(chat)
	}
}
