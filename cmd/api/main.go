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
	log.Println("Start the server...")

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	db, err := sql.Open("pgx", "postgres://postgres:password@postgres:5432/portal_development")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ctxDB, cancelDB := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelDB()
	if err := db.PingContext(ctxDB); err != nil {
		panic(err)
	}

	kafkaProducer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
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

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(Interceptor),
	)
	portal_proto.RegisterUserServer(grpcServer, server.UserServer{
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
		Ready:       make(chan struct{}),
		MutationSvc: mutationSvc,
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	kafkaConsumer, err := sarama.NewConsumerGroup([]string{"kafka:9092"}, "free-coin-for-register", sarama.NewConfig())
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	defer cancel()
	go func(c sarama.ConsumerGroup) {
		defer log.Println("finish consum")
		defer wg.Done()
		for {
			if err := c.Consume(ctx, []string{"portal.user.created.x2"}, &freeCoinAfterRegister); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}

			freeCoinAfterRegister.Ready = make(chan struct{})
		}
	}(kafkaConsumer)

	<-freeCoinAfterRegister.Ready

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	log.Println("Server is ready")
	<-sigint

	log.Println("Shutting down the server")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelShutdown()

	gracefulChan := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		gracefulChan <- struct{}{}
	}()

	select {
	case <-gracefulChan:
		break
	case <-ctxShutdown.Done():
		log.Println("Forcing shut down")
		grpcServer.Stop()
	}
}

func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("a mid")

	resp, err := handler(ctx, req)

	log.Println("b mid")

	return resp, err
}
