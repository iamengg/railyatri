package main

import (
	"log"
	"net"
	"sync"

	booking "github.com/iamengg/railyatri/bookingStub/Booking"
	handler "github.com/iamengg/railyatri/server/api"

	"google.golang.org/grpc"
)

const (
	grpcPort = ":50051"
)

func setupGRPC(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Info: ", "GRPC started at ", grpcPort)

	// Initialise TCP listener.
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer lis.Close()

	// Bootstrap server.
	uplSrv := booking.NewServer(handler.BookingHandler{})

	// Bootstrap gRPC server.
	rpcSrv := grpc.NewServer()

	log.Print("Registering grpc server")
	// Register and start gRPC server.
	booking.RegisterBookingServiceServer(rpcSrv, uplSrv)
	log.Fatal("Error: ", rpcSrv.Serve(lis))
}

func main() {
	wg := sync.WaitGroup{}
	setupGRPC(&wg)
	log.Println("Welcome to RailYatri! How may i help you?")

	wg.Wait()
}
