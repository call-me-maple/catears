package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/oklog/run"
	"github.com/thoas/bokchoy"
	"log"
	"os"
)

var (
	messageCreate  *bokchoy.Queue
	reactionCreate *bokchoy.Queue
	messageSend    *bokchoy.Queue
	messageReact   *bokchoy.Queue
	dg             *discordgo.Session
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
			select {
			case <-cancel:
				log.Printf("nyaa getting sleepy\n")
				return dg.Close()
			}
		}, func(err error) {
			close(cancel)
		})
	}
	{
		ctx, cancel := context.WithCancel(context.Background())
		engine, err := bokchoy.New(ctx, bokchoy.Config{
			Broker: bokchoy.BrokerConfig{
				Type: "redis",
				Redis: bokchoy.RedisConfig{
					Type: "client",
					Client: bokchoy.RedisClientConfig{
						Addr:     os.Getenv("REDIS_ADDR"),
						Password: os.Getenv("REDIS_PASSWORD"),
					},
				},
			},
		})
		if err != nil {
			log.Println("error opening redis connection,", err)
			return
		}
		messageCreate = engine.Queue(os.Getenv("MSG_CREATE_QUEUE"))
		messageCreate.HandleFunc(processMessage)
		reactionCreate = engine.Queue(os.Getenv("REACT_CREATE_QUEUE"))
		reactionCreate.HandleFunc(processReaction)
		messageSend = engine.Queue(os.Getenv("MSG_SEND_QUEUE"))
		messageSend.HandleFunc(sendMessage)
		messageReact = engine.Queue(os.Getenv("MSG_REACT_QUEUE"))
		messageReact.HandleFunc(sendReaction)
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
