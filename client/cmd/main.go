package main

import (
	"context"
	"fmt"
	"log"
	"time"

	booking "github.com/iamengg/railyatri/bookingStub"
	"github.com/iamengg/railyatri/server/model"

	//handler "railyatri/server/api/handlers"

	"google.golang.org/grpc"
)

const (
	GRPCPORT = ":50051"
)

var conn *grpc.ClientConn
var err error

func init() {
	// Initialise gRPC connection.
	conn, err = grpc.NewClient(GRPCPORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err.Error())
	}

}

// CreateBooking takes arguments to create booking for UserId ,for input parameters
func CreateBooking(UserId int, SourceStation string, DestinationStation string, Date string, section booking.Section) {

	client := booking.NewBookingServiceClient(conn)
	name, err := client.CreateBooking(context.Background(),
		&booking.BookingRequest{UserId: int64(UserId),
			TrainNum:           1234,
			SourceStation:      SourceStation,
			DestinationStation: DestinationStation,
			Date:               Date,
			Section:            &booking.Sections{Section: section},
		})

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Info: CreateBooking response ", name)
}

// GetUserBooking returns booking done by User for  given data & from source to destination station
func GetUserBooking(UserId int, SourceStation string, DestinationStation string, Date string, section booking.Section) {
	client := booking.NewBookingServiceClient(conn)
	name, err := client.GetUserBookings(context.Background(),
		&booking.BookingRequest{UserId: int64(UserId),
			TrainNum:           1234,
			SourceStation:      SourceStation,
			DestinationStation: DestinationStation,
			Date:               Date,
			Section:            &booking.Sections{Section: section},
		})

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Info: GetUserBooking response ", name)
}

// GetSectionBooking returns all bookings at input train section for given data
func GetSectionBooking(UserId int, SourceStation string, DestinationStation string, Date string, section booking.Section) {
	client := booking.NewBookingServiceClient(conn)
	name, err := client.GetSectionBookings(context.Background(),
		&booking.BookingRequest{UserId: int64(UserId), TrainNum: 1234,
			SourceStation:      SourceStation,
			DestinationStation: DestinationStation,
			Date:               Date,
			Section:            &booking.Sections{Section: section}})

	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Info: GetSectionBooking response ", name)
}

// ToDo : UpdateUserBooking updates userBooking
func UpdateUserBooking(UserId int, SourceStation string, DestinationStation string, Date string, section booking.Section) {
	
	client := booking.NewBookingServiceClient(conn)
	name, err := client.UpdateBooking(context.Background(),
		&booking.BookingRequest{UserId: int64(UserId),
			TrainNum:           1234,
			SourceStation:      SourceStation,
			DestinationStation: DestinationStation,
			Date:               Date,
			Section:            &booking.Sections{Section: section}})

	if err != nil {
		log.Println(err)
	}
	log.Println("Info: UpdateUserBooking response ", name)
}

// ToDo : DeleteUserBooking deletes userBooking
func DeleteUserBooking(UserId int, SourceStation string, DestinationStation string, Date string, section booking.Section) {
	
	client := booking.NewBookingServiceClient(conn)
	name, err := client.DeleteBookings(context.Background(),
		&booking.BookingRequest{UserId: int64(UserId),
			TrainNum:           1234,
			SourceStation:      SourceStation,
			DestinationStation: DestinationStation,
			Date:               Date,
			Section:            &booking.Sections{Section: section},
		})

	if err != nil {
		log.Println(err)
	}

	log.Println("Info: DeleteUserBooking response ", name)
}

func main() {
	defer conn.Close()
	section := booking.Section_A

	year, mon, day := time.Now().Date()
	date := fmt.Sprintf("%v-%v-%v", year, int(mon), day)

	for i := 1; i < 5; i++ {
		if i%2 == 0 {
			CreateBooking(i, "London", "Paris", date, booking.Section(model.SectionB))
		} else {
			CreateBooking(i, "London", "Paris", date, booking.Section(model.SectionA))
		}
	}

	// Get user bookings
	GetUserBooking(1, "London", "Paris", date, section)
	GetUserBooking(33, "London", "Paris", date, section)

	// GetSectoin bookings
	GetSectionBooking(1, "London", "Paris", date, booking.Section_B)
	GetSectionBooking(1, "London", "Paris", date, booking.Section_A)

	// update
	UpdateUserBooking(1, "London", "Paris", date, booking.Section_A)

	// Delete
	DeleteUserBooking(1, "London", "Paris", date, booking.Section_A)
}

// userID, trainName, src, dest, journeyDate, sectionType
func GetUser() *booking.User {
	return &booking.User{
		FirstName: "Pratik",
		LastName:  "shitole",
		Email:     "pratikshitole1@gmail.com",
		Id:        1,
	}
}

func GetJourney() *booking.Travel {
	return &booking.Travel{
		To:   "mumbai",
		From: "bengaluru",
	}
}

// func (b *BookingClient) CreateBooking(ctx context.Context) error {
//  resp := b.CreateBooking(ctx, &booking.BookingRequest{User: &booking.User{
// func (b *BookingClient) GetUserBookings(ctx context.Context, in *booking.User, opts ...grpc.CallOption) (*booking.BookingsResponse, error) {
// Returns bookings for both sections
// func (b *BookingClient) GetSectionBookings(ctx context.Context, in *booking.Sections, opts ...grpc.CallOption) (*booking.BookingsResponse, error) {
// can update user details only
// func (b *BookingClient) UpdateBooking(ctx context.Context, in *booking.User, opts ...grpc.CallOption) (*booking.BookingResponse, error) {
// func (b *BookingClient) DeleteBookings(ctx context.Context, in *booking.User, opts ...grpc.CallOption) (*booking.BookingResponse, error) {
