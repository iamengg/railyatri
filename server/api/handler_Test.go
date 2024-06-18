package api

import (
	"context"
	"fmt"
	"testing"

	Booking "github.com/iamengg/railyatri/bookingStub"
	"github.com/iamengg/railyatri/util"
	//booking "github.com/iamengg/railyatri/server/api"
)

var (
	dataLength   = 3
	users        = []int{1, 2, 3}
	trainNumbers = []int{123, 234, 345}
	sources      = []string{"London", "Paris", "moscow"}
	destinations = []string{"Paris", "London", "London"}
)

func GetBookingReq(userId int, source, destination, date string, section Booking.Sections) *Booking.BookingRequest {
	return &Booking.BookingRequest{
		UserId:             int64(userId),
		TrainNum:           1234,
		SourceStation:      source,
		DestinationStation: destination,
		Date:               date,
		Section:            &section,
	}

}

func TestCreateBooking(t *testing.T) {
	for i := 0; i < 2; i++ {
		r := GetBookingReq(users[i], sources[i], destinations[i], util.GetDate(i+1), Booking.Sections{Section: Booking.Section_A})
		resp, err := NewBookingServerHandler().CreateBooking(context.TODO(), r)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
	}
	
}
