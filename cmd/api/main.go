package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pursuit/portal/internal/consumer"
	"github.com/pursuit/portal/internal/proto/out/api/portal"
	"github.com/pursuit/portal/internal/proto/server"
	"github.com/pursuit/portal/internal/repo"
	"github.com/pursuit/portal/internal/service/mutation"
	"github.com/pursuit/portal/internal/service/user"

	"github.com/pursuit/event-go/pkg"

	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/Shopify/sarama"
)

func main() {
	defer log.Println("Shutdown the server success")

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
		QInterval: 1 * time.Second,
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
	pursuit_api_portal_proto.RegisterUserServer(grpcServer, server.UserServer{
		UserService:     userSvc,
		MutationService: mutationSvc,
	})

	go func() {
		log.Println("listen to 5001")
		if err := grpcServer.Serve(lis); err != nil {
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
		defer log.Println("finish consum")
		defer wg.Done()
		for {
			if err := kafkaConsumer.Consume(ctx, []string{"portal.user.created.x2"}, &freeCoinAfterRegister); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}

			freeCoinAfterRegister.Ready = make(chan bool)
		}
	}()

	<-freeCoinAfterRegister.Ready

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	log.Println("Server is ready")
	<-sigint

	log.Println("Shutting down the server")

	grpcServer.GracefulStop()
}
