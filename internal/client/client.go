package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hackfeed/stark/internal/db/cache"
	"github.com/hackfeed/stark/internal/domain"
	"github.com/hackfeed/stark/internal/store/filesrepo"
	"github.com/hackfeed/stark/internal/store/usersrepo"
	"github.com/joho/godotenv"
	"github.com/jroimartin/gocui"
	"github.com/logrusorgru/aurora/v3"
)

var (
	ctx         context.Context
	cacheClient *cache.RedisClient
	usersRepo   usersrepo.UsersRepository
	filesRepo   filesrepo.FilesRepository
	user        *domain.User
	activeChat  *domain.Chat
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	ctx = context.Background()
	cc, err := cache.NewRedisClient(ctx, &cache.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
	if err != nil {
		log.Fatalln(err)
	}
	cacheClient = cc

	usersRepo = usersrepo.NewRedisRepo(*cacheClient, 1*time.Hour)
	filesRepo = filesrepo.NewRedisRepo(*cacheClient, 1*time.Hour)
}

func Disconnect(_ *gocui.Gui, _ *gocui.View) error {
	for _, chat := range user.GetChats() {
		_, err := removeChatUser(user, chat, usersRepo)
		if err != nil {
			log.Fatalln(err)
		}
		handleLeave(ctx, cacheClient, chat, user, "disconnected")
	}

	return gocui.ErrQuit
}

func Send(g *gocui.Gui, v *gocui.View) error {
	msg := strings.TrimSpace(v.Buffer())
	if msg == "" {
		return nil
	}

	message := domain.NewMessage(user.GetName(), msg)
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
	activeChat = getActiveChat(ctx, cacheClient, "global", user, true)
	_, err := addChatUser(user, activeChat, usersRepo)
	if err != nil {
		log.Fatalln(err)
	}
	handleJoin(ctx, cacheClient, activeChat, user)

	g.SetViewOnTop("messages")
	g.SetViewOnTop("users")
	g.SetViewOnTop("input")
	g.SetViewOnTop("chats")
	g.SetCurrentView("input")

	messagesView, _ := g.View("messages")
	usersView, _ := g.View("users")
	chatsView, _ := g.View("chats")

	updateChatsView(g, chatsView, user, true)

	go func() {
		for {
		chatLoop:
			for msg := range activeChat.GetMessages() {
				switch {
				case strings.HasPrefix(msg, "/users>"):
					updateUsersView(g, usersView, activeChat, usersRepo, true)

				case getCommandFromMessage(msg) == "/join":
					if getUserFromMessage(msg) != user.GetName() {
						break
					}

					newChatName := strings.TrimSpace(strings.SplitAfter(msg, "/join")[1])
					if _, ok := user.GetChats()[newChatName]; ok {
						updateMessagesView(
							g,
							messagesView,
							aurora.Sprintf(
								aurora.Yellow("You are already in %s. Use %s to move there"),
								aurora.Green(newChatName),
								aurora.Green(fmt.Sprintf("/switch %s", newChatName)),
							),
							false,
						)
						break
					}

					activeChat = getActiveChat(ctx, cacheClient, newChatName, user, true)
					_, err = addChatUser(user, activeChat, usersRepo)
					if err != nil {
						log.Fatalln(err)
					}
					handleJoin(ctx, cacheClient, activeChat, user)

					updateUsersView(g, usersView, activeChat, usersRepo, true)
					updateChatsView(g, chatsView, user, true)
					clearView(g, messagesView)

					updateMessagesCounter(g, chatsView, user)
					break chatLoop

				case getCommandFromMessage(msg) == "/switch":
					if getUserFromMessage(msg) != user.GetName() {
						break
					}

					newChatName := strings.TrimSpace(strings.SplitAfter(msg, "/switch")[1])
					if _, ok := user.GetChats()[newChatName]; !ok {
						updateMessagesView(
							g,
							messagesView,
							aurora.Sprintf(
								aurora.Yellow("There is no %s chat, Use %s to create it"),
								aurora.Green(newChatName),
								aurora.Green(fmt.Sprintf("/join %s", newChatName)),
							),
							false,
						)
						break
					}
					if newChatName == activeChat.GetName() {
						updateMessagesView(g, messagesView, aurora.Yellow("Can't switch to current active chat").String(), false)
						break
					}

					activeChat = getActiveChat(ctx, cacheClient, newChatName, user, false)
					_, err = addChatUser(user, activeChat, usersRepo)
					if err != nil {
						log.Fatalln(err)
					}

					updateMessagesViewWithBuffer(g, messagesView, activeChat, true)
					updateUsersView(g, usersView, activeChat, usersRepo, true)
					updateChatsView(g, chatsView, user, true)

					updateMessagesCounter(g, chatsView, user)
					break chatLoop

				case getCommandFromMessage(msg) == "/leave":
					if getUserFromMessage(msg) != user.GetName() {
						break
					}

					if activeChat.GetName() == "global" {
						updateMessagesView(g, messagesView, aurora.Red("Can't leave global").String(), false)
						break
					}

					_, err = removeChatUser(user, activeChat, usersRepo)
					if err != nil {
						log.Fatalln(err)
					}
					handleLeave(ctx, cacheClient, activeChat, user, "left the chat")
					user.RemoveChat(activeChat.GetName())
					activeChat = getActiveChat(ctx, cacheClient, user.GetActiveChat().GetName(), user, false)

					updateMessagesViewWithBuffer(g, messagesView, activeChat, true)
					updateUsersView(g, usersView, activeChat, usersRepo, true)
					updateChatsView(g, chatsView, user, true)

					updateMessagesCounter(g, chatsView, user)
					break chatLoop

				case getCommandFromMessage(msg) == "/upload":
					if getUserFromMessage(msg) != user.GetName() {
						break
					}

					fileName := strings.Split(msg, " ")[3]
					data, err := os.ReadFile(fileName)
					if err != nil {
						updateMessagesView(g, messagesView, aurora.Red("Can't open file. Check whether path is correct").String(), false)
						break
					}
					saveFileName := strings.Split(msg, " ")[4]
					_, err = addChatFile(domain.NewFile(saveFileName, data), activeChat, filesRepo)
					if err != nil {
						updateMessagesView(g, messagesView, aurora.Red("Error while uploading file. Try again later").String(), false)
						break
					}

					updateMessagesView(g, messagesView, aurora.Green("File is uploaded").String(), false)

				case getCommandFromMessage(msg) == "/download":
					if getUserFromMessage(msg) != user.GetName() {
						break
					}

					files, err := filesRepo.GetFiles(activeChat.GetName())
					if err != nil {
						updateMessagesView(g, messagesView, aurora.Red("Error while fetching chat files").String(), false)
						break
					}
					fileName := strings.Split(msg, " ")[3]

					data, ok := files[fileName]
					if !ok {
						updateMessagesView(g, messagesView, aurora.Red("There is no file with given name for this chat").String(), false)
						break
					}

					saveFileName := strings.Split(msg, " ")[4]
					if err := os.WriteFile(saveFileName, data, 0644); err != nil {
						updateMessagesView(g, messagesView, aurora.Red("Error while saving data").String(), false)
						break
					}

					updateMessagesView(g, messagesView, aurora.Green("File is downloaded").String(), false)

				case getCommandFromMessage(msg) == "/files":
					if getUserFromMessage(msg) != user.GetName() {
						break
					}

					files, err := filesRepo.GetFiles(activeChat.GetName())
					if err != nil {
						updateMessagesView(g, messagesView, aurora.Red("Error while fetching chat files").String(), false)
						break
					}
					if len(files) == 0 {
						updateMessagesView(g, messagesView, aurora.Yellow("No files available for this chat").String(), false)
						break
					}

					header := "Chat files:\n"
					res := ""
					for k := range files {
						res += k + "\n"
					}
					res = strings.TrimSpace(res)

					updateMessagesView(g, messagesView, aurora.Sprintf("%s%s", header, aurora.Blue(res)), false)

				default:
					updateMessagesView(g, messagesView, msg, false)
				}
			}
		}
	}()

	return nil
}
