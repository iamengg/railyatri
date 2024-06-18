package api

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	Booking "github.com/iamengg/railyatri/bookingStub"
	db "github.com/iamengg/railyatri/server/database"
	"github.com/iamengg/railyatri/server/model"
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

func GetBookingReqObj(userId int, source, destination, date string, section Booking.Sections) *Booking.BookingRequest {
	return &Booking.BookingRequest{
		UserId:             int64(userId),
		SourceStation:      source,
		DestinationStation: destination,
		Date:               date,
		Section:            &section,
	}

}

func TestCreateBooking(t *testing.T) {
	bookingIds := []int64{}
	bkgHandler := NewBookingServerHandler()
	//Test createBooking
	for i := 0; i < 20; i++ {
		//log.Printf("Creating booking for %d", i+1)
		tmpUserId := rand.Intn(len(users))

		//source & dest are keeping at 0th index for which stations are added in routes (sourceStation_destinationStation)
		r := GetBookingReqObj(users[tmpUserId], sources[0], destinations[0], util.GetDate(i+1), Booking.Sections{Section: Booking.Section_A})
		resp, err := bkgHandler.CreateBooking(context.TODO(), r)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
		bookingIds = append(bookingIds, int64(resp.BookingId))
	}

	//Verify created booking details
	var receipt model.UserBookingDetails
	prevSeatNum := -1
	for _, bookingId := range bookingIds {
		receipt = db.GetBookingReceipt(bookingId)

		fmt.Printf("%#v => %#v\n", bookingId, receipt)
		if prevSeatNum != -1 && prevSeatNum == receipt.SeatNum {
			fmt.Println("Different seatnumbser are not getting assigned ")
			t.Fail()
		}

		prevSeatNum = receipt.SeatNum
	}
}

// create booking & retrieve same
func TestGetUserBookings(t *testing.T) {
	for i := 0; i < 2; i++ {
		r := GetBookingReqObj(users[i], sources[i], destinations[i], util.GetDate(i+1), Booking.Sections{Section: Booking.Section_A})
		resp, err := NewBookingServerHandler().GetUserBookings(context.TODO(), r)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
	}
}

// create few bookings at section A, B & retrieve them & check count
func TestGetSectionBookings(t *testing.T) {
	for i := 0; i < 1; i++ {
		r := GetBookingReqObj(users[i], sources[i], destinations[i], util.GetDate(i+1), Booking.Sections{Section: Booking.Section_A})
		resp, err := NewBookingServerHandler().GetSectionBookings(context.TODO(), r)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)

		// We shouldn't get any booking here
		r = GetBookingReqObj(users[i], sources[i], destinations[i], util.GetDate(i+1), Booking.Sections{Section: Booking.Section_B})
		resp, err = NewBookingServerHandler().GetSectionBookings(context.TODO(), r)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
	}
}
