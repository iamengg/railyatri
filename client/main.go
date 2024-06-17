package main

import (
	"context"
	"log"

	booking "github.com/iamengg/railyatri/bookingStub/Booking"

	//handler "railyatri/server/api/handlers"

	"google.golang.org/grpc"
)

const (
	GRPCPORT = ":50051"
)

func sendOverGrpc() {

	// Initialise gRPC connection.
	conn, err := grpc.Dial(GRPCPORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := booking.NewBookingServiceClient(conn)
	name, err := client.Get(context.Background())
	if err != nil {
		log.Println(err)
	}
	log.Println("Info: ", "Uploaded file is ", name)
}

func main() {
	sendOverGrpc()
}

func NewClient() {
	return booking.NewClient
}

type BookingClient struct {
	obj *booking.BookingServiceClient
}

func (b *BookingClient) CreateBooking(ctx context.Context) error {
	firstName := "Pratik"
	lastName := "Shitole"
	email := "pratikshitole1@gmail.com"
	user := booking.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	req := booking.BookingRequest{
		User:   user,
		Travel: booking.Travel{},
	}
	resp := b.obj.CreateBooking(ctx, req)

	log.Println(resp)
}

func (b *BookingClient) GetUserBookings(ctx context.Context, in *User, opts ...grpc.CallOption) (*BookingsResponse, error) {
	return nil, nil
}

// Returns bookings for both sections
func (b *BookingClient) GetSectionBookings(ctx context.Context, in *Sections, opts ...grpc.CallOption) (*BookingsResponse, error) {
	return nil, nil
}

// can update user details only
func (b *BookingClient) UpdateBooking(ctx context.Context, in *User, opts ...grpc.CallOption) (*BookingResponse, error) {
	return nil, nil
}

func (b *BookingClient) DeleteBookings(ctx context.Context, in *User, opts ...grpc.CallOption) (*BookingResponse, error) {
	return nil, nil
}
