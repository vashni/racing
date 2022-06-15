package main

import (
	"database/sql"
	"flag"
	"log"
	"net"

	"github.com/vashni/racing/sports/db"
	"github.com/vashni/racing/sports/proto/sports"
	"github.com/vashni/racing/sports/service"
	"google.golang.org/grpc"
)

var (
	//sports grpc end point listening in 9001.
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9001", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s\n", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9001")
	if err != nil {
		return err
	}

	sportingDB, err := sql.Open("sqlite3", "./db/sports.db")
	if err != nil {
		return err
	}

	sportsRepo := db.NewSportsRepo(sportingDB)
	if err := sportsRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	//Registering new sports server.
	sports.RegisterSportsServer(
		grpcServer,
		service.NewSportsService(
			sportsRepo,
		),
	)

	log.Printf("gRPC server listening on: %s\n", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
