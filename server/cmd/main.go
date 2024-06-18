package main

import (
	"log"
	"net"

	booking "github.com/iamengg/railyatri/bookingStub"
	handler "github.com/iamengg/railyatri/server/api"

	"google.golang.org/grpc"
)

const (
	//GRPC port where grpc server listenes
	GRPCPORT = ":50051"
)

// setupGRPC register handler for grpc apis &
// starts listning at GRPC port
func setupGRPC() {
	log.Println("Info: ", "GRPC started at ", GRPCPORT)

	// Initialise TCP listener.
	lis, err := net.Listen("tcp", GRPCPORT)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer lis.Close()

	// grpc API handler server
	grpcAPIHandlerServer := handler.NewBookingServerHandler()

	// gRPC server.
	rpcServer := grpc.NewServer()

	log.Print("Registering grpc server")

	// Register and start gRPC server.
	booking.RegisterBookingServiceServer(rpcServer, grpcAPIHandlerServer)
	log.Fatal("Error: ", rpcServer.Serve(lis))
}

func main() {
	log.Println("Welcome to RailYatri! Setting up server ...")

	setupGRPC()
	log.Println(" RailYatri! Server is running ...")
}
