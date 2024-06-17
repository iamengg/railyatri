package api

import (
	"context"
	"log"
	. "railyatri/BookingStub/Booking"
)

type BookingHandler struct {
}

type BookingServiceServer interface {
	CreateBooking(context.Context, *BookingRequest) (*BookingResponse, error)
	GetUserBookings(context.Context, *User) (*BookingsResponse, error)
	//Returns bookings for both sections
	GetSectionBookings(context.Context, *Sections) (*BookingsResponse, error)
	//can update user details only
	UpdateBooking(context.Context, *User) (*BookingResponse, error)
	DeleteBookings(context.Context, *User) (*BookingResponse, error)
}

func (b *BookingHandler) CreateBooking(context.Context, *BookingRequest) (*BookingResponse, error) {
	log.Println("Handler CreateBooking")
	return nil, nil
}

func (b *BookingHandler) GetUserBookings(context.Context, *User) (*BookingsResponse, error) {
	log.Println("Handler GetUserBookings")
	return nil, nil
}

// Returns bookings for both sections
func (b *BookingHandler) GetSectionBookings(context.Context, *Sections) (*BookingsResponse, error) {
	log.Println("Handler GetSectionBookings")
	return nil, nil
}

// can update user details only
func (b *BookingHandler) UpdateBooking(context.Context, *User) (*BookingResponse, error) {
	log.Println("Handler UpdateBooking")
	return nil, nil
}

func (b *BookingHandler) DeleteBookings(context.Context, *User) (*BookingResponse, error) {
	log.Println("Handler DeleteBookings")
	return nil, nil
}
