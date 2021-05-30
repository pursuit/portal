package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pursuit/portal/internal/consumer"
	"github.com/pursuit/portal/internal/proto/out"
	"github.com/pursuit/portal/internal/proto/server"
	"github.com/pursuit/portal/internal/repo"
	"github.com/pursuit/portal/internal/rest"
	"github.com/pursuit/portal/internal/service/mutation"
	"github.com/pursuit/portal/internal/service/user"

	"github.com/pursuit/event-go/pkg"

	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/Shopify/sarama"
)

func main() {
	defer fmt.Println("Shutdown the server success")

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	db, err := sql.Open("pgx", "postgres://postgres:password@localhost:5432/portal_development")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	kafkaProducer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}

	kafkaPublishFromSQL := pkg.KafkaPublishFromSQL{
		Batch:     uint(2),
		DB:        db,
		Kafka:     kafkaProducer,
		WorkerNum: uint(2),
	}
	go kafkaPublishFromSQL.Run()
	defer kafkaPublishFromSQL.Shutdown()

	mutationSvc := mutation.Svc{
		DB:           db,
		MutationRepo: repo.MutationRepo{},
	}
	userSvc := user.Svc{
		DB:       db,
		UserRepo: repo.UserRepo{},
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServer(grpcServer, server.UserServer{
		UserService: userSvc,
	})

	go func() {
		fmt.Println("listen to 5001")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	restHandler := rest.Handler{userSvc}
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))
	r.Post("/users", restHandler.CreateUser)

	restServer := http.Server{
		Addr:    ":5002",
		Handler: r,
	}

	go func() {
		fmt.Println("listen to 5002")
		if err := restServer.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			panic(err)
		}
	}()

	freeCoinAfterRegister := consumer.FreeCoinRegisterConsumer{
		Ready:       make(chan bool),
		MutationSvc: mutationSvc,
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	kafkaConsumer, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "free-coin-for-register", sarama.NewConfig())
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	defer cancel()
	go func() {
		defer fmt.Println("finish consum")
		defer wg.Done()
		for {
			fmt.Println("a")
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := kafkaConsumer.Consume(ctx, []string{"portal.user.created.x2"}, &freeCoinAfterRegister); err != nil {
				fmt.Printf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			freeCoinAfterRegister.Ready = make(chan bool)
		}
	}()

	<-freeCoinAfterRegister.Ready

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Server is ready")
	<-sigint

	fmt.Println("Shutting down the server")

	if err := restServer.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	grpcServer.GracefulStop()
}
