package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	api "github.com/vv-sam/otus-project/proto/grpc/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	login    = flag.String("login", "admin", "login")
	password = flag.String("password", "1234", "password")

	token string
	au    api.AuthServiceClient
	ac    api.AgentServiceClient
	cc    api.ConfigurationServiceClient
	tc    api.TaskServiceClient
)

func main() {
	flag.Parse()

	if *password == "" {
		log.Fatal("password is required")
	}

	log.Println("Agent is running...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Agent is stopping...")
		stop()
	}()

	s, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	au = api.NewAuthServiceClient(s)
	ac = api.NewAgentServiceClient(s)
	cc = api.NewConfigurationServiceClient(s)
	tc = api.NewTaskServiceClient(s)

	loginResponse, err := au.Login(ctx, &api.LoginRequest{
		Username: *login,
		Password: *password,
	})
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	log.Printf("token: %v", loginResponse.Token)

	token = loginResponse.Token

	agents, err := ac.GetAll(ctx, &api.GetAllAgentsRequest{})
	if err != nil {
		log.Fatalf("Failed to get agents: %v", err)
	}
	log.Printf("current agents: %v", agents)

	configs, err := cc.GetAll(ctx, &api.GetAllConfigurationsRequest{})
	if err != nil {
		log.Fatalf("Failed to get configs: %v", err)
	}
	log.Printf("current configs: %v", configs)

	tasks, err := tc.GetAll(ctx, &api.GetAllTasksRequest{})
	if err != nil {
		log.Fatalf("Failed to get tasks: %v", err)
	}
	log.Printf("current tasks: %v", tasks)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+token))
	generateStructs(ctx)
}

func generateStructs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("generateStructs is stopping...")
			return
		case <-time.After(time.Second * 5):
			r := rand.Int31n(4)
			switch r {
			case 0:
				s := api.AgentInfo{
					AgentId: uuid.New().String(),
					Status:  0,
				}
				_, err := ac.Post(ctx, &api.PostAgentRequest{
					Agent: &s,
				})
				if err != nil {
					log.Printf("Failed to post agent: %v\n", err)
					continue
				}
				log.Printf("posted agent with id: %v\n", s.AgentId)
				r, err := ac.GetById(ctx, &api.GetAgentByIdRequest{
					Id: s.AgentId,
				})
				if err != nil {
					log.Printf("Failed to get agent: %v\n", err)
					continue
				}
				log.Printf("got agent with id: %v\n", r.Agent.AgentId)
			case 1:
				s := api.FactorioConfig{
					Base: &api.BaseConfig{
						Id:      uuid.New().String(),
						AgentId: uuid.New().String(),
					},
				}
				_, err := cc.Post(ctx, &api.PostConfigurationRequest{
					Configuration: &s,
				})
				if err != nil {
					log.Printf("Failed to post configuration: %v\n", err)
					continue
				}
				log.Printf("posted configuration with id: %v\n", s.Base.Id)
				r, err := cc.GetById(ctx, &api.GetConfigByIdRequest{
					Id: s.Base.Id,
				})
				if err != nil {
					log.Printf("Failed to get configuration: %v\n", err)
					continue
				}
				log.Printf("got configuration with id: %v\n", r.Configuration.Base.Id)
			case 2:
				s := api.Task{
					Id:     uuid.New().String(),
					Status: api.TaskStatus_TASK_STATUS_QUEUED,
					Type:   "factorio",
				}
				_, err := tc.Post(ctx, &api.PostTaskRequest{
					Task: &s,
				})
				if err != nil {
					log.Printf("Failed to post task: %v\n", err)
					continue
				}
				log.Printf("posted task with id: %v\n", s.Id)
				r, err := tc.GetById(ctx, &api.GetTaskByIdRequest{
					Id: s.Id,
				})
				if err != nil {
					log.Printf("Failed to get task: %v\n", err)
					continue
				}
				log.Printf("got task with id: %v\n", r.Task.Id)
			}
		}
	}
}
