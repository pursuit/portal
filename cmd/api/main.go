package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pursuit/portal/internal/proto/out"
	"github.com/pursuit/portal/internal/proto/server"
	"github.com/pursuit/portal/internal/repo"
	"github.com/pursuit/portal/internal/service/user"

	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/stdlib"
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

	grpcServer := grpc.NewServer()
	proto.RegisterUserServer(grpcServer, server.UserServer{
		UserService: user.Svc{
			DB:       db,
			UserRepo: repo.UserRepo{},
		},
	})

	go func() {
		fmt.Println("listen to 5001")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	<-sigint

	fmt.Println("Shutting down the server")

	grpcServer.GracefulStop()
}
