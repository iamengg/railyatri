package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"

	booking "github.com/iamengg/railyatri/bookingStub"
	model "github.com/iamengg/railyatri/server/model"
	"github.com/iamengg/railyatri/util"
)

var (
	users        = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	trainNumbers = []int{1234, 2344}
	sources      = []string{"London", "Paris", "moscow"}
	destinations = []string{"Paris", "London", "London"}
)

func GetBookingReqObj(userId int, source, destination, date string, section booking.Sections) *booking.BookingRequest {
	return &booking.BookingRequest{
		UserId:             int64(userId),
		TrainNum:           int32(trainNumbers[0]),
		SourceStation:      source,
		DestinationStation: destination,
		Date:               date,
		Section:            &section,
	}

}

func TestCreateBooking(t *testing.T) {
	bookingIds := []int64{}
	bkgHandler := NewBookingServerHandler()

	actualBookings := 2
	for i := 0; i < actualBookings; i++ {
		tmpUserId := rand.Intn(len(users))

		//source & dest are keeping at 0th index for which stations are added in routes (sourceStation_destinationStation)
		r := GetBookingReqObj(users[tmpUserId], sources[0], destinations[0], util.GetDate(i+1), booking.Sections{Section: booking.Section_A})
		resp, err := bkgHandler.CreateBooking(context.TODO(), r)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
		bookingIds = append(bookingIds, int64(resp.BookingId))
	}

	if actualBookings != len(bookingIds) {
		t.Fail()
	}
}

// create booking & retrieve same
func TestGetUserBookings(t *testing.T) {
	//bookingIds := []int64{}
	bkgHandler := NewBookingServerHandler()
	totalPossibleBookings := model.GetTotalSeats()
	log.Println("Total seats are ", totalPossibleBookings)
	actualBookings := 0
	for i := 0; i < totalPossibleBookings; i++ {

		//source & dest are keeping at 0th index for which stations are added in routes (sourceStation_destinationStation)
		r := GetBookingReqObj(users[0], sources[0], destinations[0], util.GetDate(i+1), booking.Sections{Section: booking.Section_A})
		resp, err := bkgHandler.CreateBooking(context.TODO(), r)
		if err != nil {
			continue
		}
		actualBookings++
		fmt.Println(resp)
		//bookingIds = append(bookingIds, int64(resp.BookingId))
	}

	if actualBookings != totalPossibleBookings {
		t.Fail()
	}
}

// create few bookings at section A, B & retrieve them & check count
func TestGetSectionBookings(t *testing.T) {
	seatsAtSectionA := model.GetSeatsAtSection(0)

	for i := 0; i < seatsAtSectionA; i++ {
		r := GetBookingReqObj(users[i], sources[0], destinations[0], util.GetDate(i+1), booking.Sections{Section: booking.Section_A})
		_, err := NewBookingServerHandler().CreateBooking(context.TODO(), r)
		if err != nil {
			fmt.Println(err)
		}

		//We shouldn't get any booking here
		r = GetBookingReqObj(users[i], sources[0], destinations[0], util.GetDate(i+1), booking.Sections{Section: booking.Section_A})
		sectionBookings, err := NewBookingServerHandler().GetSectionBookings(context.TODO(), r)
		if err != nil {
			fmt.Println(err)
		}
		//section bookings should be equall to created bookings at that section
		if sectionBookings == nil || len(sectionBookings.Bookings) != i+1 {
			t.Fail()
		}
	}
}
