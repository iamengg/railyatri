package api

import (
	"context"
	"fmt"
	"log"

	. "github.com/iamengg/railyatri/bookingStub"
	Booking "github.com/iamengg/railyatri/bookingStub"
	db "github.com/iamengg/railyatri/server/database"
)

type BookingHandler struct{}

func NewBookingServerHandler() *BookingHandler {
	return &BookingHandler{}
}

func (b *BookingHandler) CreateBooking(c context.Context, r *BookingRequest) (*BookingResponse, error) {
	log.Printf("CreateBooking Request is %v\n", r)
	bookingId, SeatNum, err := db.CreateBooking(r.UserId, r.SourceStation, r.DestinationStation, int(r.Section.Section), r.Date)

	return &BookingResponse{BookingId: int64(bookingId), SeatNumber: int32(SeatNum)}, err
}

// return []{bookingIds, seatNumber for bookingId}
// func (b *BookingHandler) GetUserBookings(c context.Context, r *BookingRequest) (*Booking.BookingsResponse, error) {
// 	bookingIdSeatNumbers, err := db.GetUserBookings(r.UserId, r.TrainNum, r.SourceStation, r.DestinationStation, int(r.Section.Section), r.Date)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return nil, err
// 	}
// 	log.Println(bookingIdSeatNumbers)

// 	//Get all receipts of this users bookings
// 	userReceipts := make([]model.UserBookingDetails, 0)
// 	for _, bookingId := range bookingIdSeatNumbers.Bookings {
// 		receipt := db.GetBookingReceipt(bookingId.BookingId)
// 		userReceipts = append(userReceipts, receipt)
// 	}
// 	fmt.Println(userReceipts)
// 	return bookingIdSeatNumbers, nil
// }

// return []{bookingIds, seatNumber for bookingId}
func (b *BookingHandler) GetUserBookingReceipts(c context.Context, r *BookingRequest) (*Booking.BookingReceipts, error) {
	userReceipts, err := db.GetUserBookingReceipts(r.UserId, r.TrainNum, r.SourceStation, r.DestinationStation, int(r.Section.Section), r.Date)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return userReceipts, nil
}

// Returns bookings for both sections
func (b *BookingHandler) GetSectionBookings(c context.Context, r *BookingRequest) (*Booking.BookingsResponse, error) {
	log.Printf("Handler GetSectionBookings for %d-%s-%v", r.TrainNum, r.Date, fmt.Sprintf("%v", r.Section))

	bookingIdSeatNumbers, err := db.GetSectionBookings(r.UserId, r.TrainNum, r.SourceStation, r.DestinationStation, int(r.Section.Section), r.Date)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	log.Println(bookingIdSeatNumbers)

	return bookingIdSeatNumbers, nil
}

// Delete users booking for Date, from given source to destination
func (b *BookingHandler) DeleteBookings(c context.Context, r *DeleteBookingRequest) (*DeleteBookingResponse, error) {

	err := db.DeleteUserBookings(r.UesrId, r.BookingId)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	return &DeleteBookingResponse{Success: true}, nil
}

// can update user details only
func (b *BookingHandler) UpdateBooking(c context.Context, r *BookingRequest) (*BookingResponse, error) {
	err := db.UpdateUserBooking(r.UserId, r.TrainNum, r.SourceStation, r.DestinationStation, int(r.Section.Section), r.Date)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	return nil, nil
}
