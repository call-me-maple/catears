package main

import (
	"context"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v7"
	"github.com/oklog/run"
	"github.com/thoas/bokchoy"
	"github.com/thoas/bokchoy/logging"
)

var (
	messageCreate  *bokchoy.Queue
	reactionCreate *bokchoy.Queue
	messageSend    *bokchoy.Queue
	messageReact   *bokchoy.Queue
	followUp       *bokchoy.Queue
	dg             *discordgo.Session
	client         *redis.Client
)

func main() {
	var g run.Group
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			dg, _ = discordgo.New(os.Getenv("BOT_TOKEN"))
			dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
			dg.AddHandler(messageCreated)
			dg.AddHandler(reactionCreated)
			err := dg.Open()
			if err != nil {
				log.Println("error opening connection,", err)
				return err
			}
			log.Println("catears! wakey wakey")
			<-cancel
			log.Printf("nyaa getting sleepy\n")
			return dg.Close()
		}, func(err error) {
			close(cancel)
		})
	}
	{
		ctx, cancel := context.WithCancel(context.Background())
		client = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		})
		engine, err := bokchoy.New(ctx, bokchoy.Config{},
			bokchoy.WithBroker(bokchoy.NewRedisBroker(client, "client", "", logging.NewNopLogger())))
		if err != nil {
			log.Println("error opening redis connection,", err)
			return
		}
		messageCreate = engine.Queue(os.Getenv("MSG_CREATE_QUEUE"))
		reactionCreate = engine.Queue(os.Getenv("REACT_CREATE_QUEUE"))
		messageSend = engine.Queue(os.Getenv("MSG_SEND_QUEUE"))
		messageReact = engine.Queue(os.Getenv("MSG_REACT_QUEUE"))
		followUp = engine.Queue(os.Getenv("FOLLOW_UP_QUEUE"))

		messageCreate.HandleFunc(processMessage)
		reactionCreate.HandleFunc(processReaction)
		messageSend.HandleFunc(sendMessage)
		messageReact.HandleFunc(sendReaction)
		followUp.HandleFunc(checkin)

		g.Add(func() error {
			return engine.Run(ctx)
		}, func(error) {
			log.Print("Received signal, gracefully stopping")
			engine.Stop(ctx)
			cancel()
		})
	}
	log.Printf("The group was terminated with: %v\n", g.Run())
}
