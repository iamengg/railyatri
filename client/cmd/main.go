package main

import (
	"context"
	"log"

	booking "github.com/iamengg/railyatri/bookingStub"
	config "github.com/iamengg/railyatri/configs"
	_ "github.com/iamengg/railyatri/server/model"
	util "github.com/iamengg/railyatri/util"

	"google.golang.org/grpc"
)

const (
	GRPCPORT = config.GRPCPORT
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
	log.Printf("Info: CreateBooking request userId %d, src %s , dest %s, date %s section %d", UserId, SourceStation,
		DestinationStation, Date, section)
	client := booking.NewBookingServiceClient(conn)
	response, err := client.CreateBooking(context.Background(),
		&booking.BookingRequest{UserId: int64(UserId),
			SourceStation:      SourceStation,
			DestinationStation: DestinationStation,
			Date:               Date,
			Section:            &booking.Sections{Section: section},
		})

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Info: CreateBooking response ", response)
}

// GetUserBooking returns booking done by User for  given data & from source to destination station
func GetUserBookingReceipts(UserId int, SourceStation string, DestinationStation string, Date string, section booking.Section) {
	client := booking.NewBookingServiceClient(conn)
	receipts, err := client.GetUserBookingReceipts(context.Background(),
		&booking.BookingRequest{UserId: int64(UserId),
			SourceStation:      SourceStation,
			DestinationStation: DestinationStation,
			Date:               Date,
			Section:            &booking.Sections{Section: section},
		})

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Info: GetUserBooking receipt ", receipts)
}

// GetSectionBooking returns all bookings at input train section for given data
func GetSectionBooking(UserId int, trainNum int, SourceStation string, DestinationStation string, Date string, section booking.Section) {
	client := booking.NewBookingServiceClient(conn)
	name, err := client.GetSectionBookings(context.Background(),
		&booking.BookingRequest{UserId: int64(UserId),
			TrainNum:           int32(trainNum),
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
func DeleteUserBooking(UserId int, bookingId int64) {

	client := booking.NewBookingServiceClient(conn)
	name, err := client.DeleteBookings(context.Background(),
		&booking.DeleteBookingRequest{
			UesrId:    int64(UserId),
			BookingId: bookingId,
		})

	if err != nil {
		log.Println(err)
	}

	log.Println("Info: DeleteUserBooking response ", name)
}

func main() {
	defer conn.Close()
	date := util.GetDate(5)

	for i := 1; i < 5; i++ {
		if i%2 == 0 {
			CreateBooking(1, "London", "Paris", date, booking.Section(booking.Section_A))
		} else {
			CreateBooking(i, "London", "Paris", date, booking.Section(booking.Section_A))
		}
	}

	// Get user bookings
	GetUserBookingReceipts(1, "London", "Paris", date, booking.Section(booking.Section_A))
	//GetUserBookingReceipts(3, "London", "Paris", date, booking.Section(booking.Section_B))

	// GetSectoin bookings
	//GetSectionBooking(1, 1234, "London", "Paris", date, booking.Section(model.SectionB))
	GetSectionBooking(1, 1234, "London", "Paris", date, booking.Section(booking.Section_A))

	// update
	//UpdateUserBooking(1, "London", "Paris", date, booking.Section_A)

	// Delete
	bkId := 1718783300964514040
	DeleteUserBooking(1, int64(bkId))
	log.Printf("Bookings aftr removing booking id bkId %d", bkId)
	GetSectionBooking(1, 1234, "London", "Paris", date, booking.Section(booking.Section_A))
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
