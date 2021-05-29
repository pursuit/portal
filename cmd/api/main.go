package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pursuit/portal/internal/proto/out"
	"github.com/pursuit/portal/internal/proto/server"
	"github.com/pursuit/portal/internal/repo"
	"github.com/pursuit/portal/internal/rest"
	"github.com/pursuit/portal/internal/service/user"

	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	<-sigint

	fmt.Println("Shutting down the server")

	if err := restServer.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	grpcServer.GracefulStop()
}
