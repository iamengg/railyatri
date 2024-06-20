package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"

	booking "github.com/iamengg/railyatri/bookingStub"
	"github.com/iamengg/railyatri/server/model"

	// model "github.com/iamengg/railyatri/server/model"
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

// Test to create bookings
// Test to make sure for same trains different sections will start seatnumbers with 1
func TestCreateBooking(t *testing.T) {
	bookingIds := []int64{}
	bkgHandler := NewBookingServerHandler()

	actualBookings := 5
	var section booking.Section
	SectionWiseBookingIds := make(map[booking.Section][]int)
	//section A should have 1,2 seat booking while , Section B should have 1,2,3 seat bookings
	for i := 0; i < actualBookings; i++ {
		tmpUserId := rand.Intn(len(users))
		if i < 2 {
			//Create 2 bookings in section A
			section = booking.Section_A
		} else {
			//Create 3 bookings in section B
			section = booking.Section_B
		}
		//source & dest are keeping at 0th index for which stations are added in routes (sourceStation_destinationStation)
		r := GetBookingReqObj(users[tmpUserId],
			sources[0],
			destinations[0],
			util.GetDate(5),
			booking.Sections{Section: section},
		)

		resp, err := bkgHandler.CreateBooking(context.TODO(), r)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
		bookingIds = append(bookingIds, int64(resp.BookingId))
		SectionWiseBookingIds[section] = append(SectionWiseBookingIds[section], int(resp.SeatNumber))
	}

	for itSection, itArr := range SectionWiseBookingIds {
		if itSection == booking.Section_A {
			//there are 3 bookings in sectionA
			for _, seatNum := range itArr {
				if seatNum > 2 {
					t.Fail()
				}
			}
		} else if itSection == booking.Section_B {
			//there are 3 bookings in sectionB
			for _, seatNum := range itArr {
				if seatNum > 3 {
					t.Fail()
				}
			}
		}
	}

	//Total bookings done are same as that of bookings done
	if actualBookings != len(bookingIds) {
		t.Fail()
	}
}

// Test for validating after booking each section able to return it's bookings
func TestSectionWiseBookigs(t *testing.T) {
	bookingIds := []int64{}
	bkgHandler := NewBookingServerHandler()

	actualBookings := 5
	var section booking.Section
	SectionWiseBookingIds := make(map[booking.Section][]int)
	//section A should have 1,2 seat booking while , Section B should have 1,2,3 seat bookings
	for i := 0; i < actualBookings; i++ {
		tmpUserId := rand.Intn(len(users))
		if i < 2 {
			//Create 2 bookings in section A
			section = booking.Section_A
		} else {
			//Create 3 bookings in section B
			section = booking.Section_B
		}
		//source & dest are keeping at 0th index for which stations are added in routes (sourceStation_destinationStation)
		r := GetBookingReqObj(users[tmpUserId],
			sources[0],
			destinations[0],
			util.GetDate(5),
			booking.Sections{Section: section},
		)

		resp, err := bkgHandler.CreateBooking(context.TODO(), r)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
		bookingIds = append(bookingIds, int64(resp.BookingId))
		SectionWiseBookingIds[section] = append(SectionWiseBookingIds[section], int(resp.SeatNumber))
	}

	for itSection, itArr := range SectionWiseBookingIds {
		if itSection == booking.Section_A {
			//there are 3 bookings in sectionA
			for _, seatNum := range itArr {
				if seatNum > 2 {
					t.Fail()
				}
			}
		} else if itSection == booking.Section_B {
			//there are 3 bookings in sectionB
			for _, seatNum := range itArr {
				if seatNum > 3 {
					t.Fail()
				}
			}
		}
	}

	if actualBookings != len(bookingIds) {
		t.Fail()
	}

	r := GetBookingReqObj(users[1],
		sources[0],
		destinations[0],
		util.GetDate(5),
		booking.Sections{Section: booking.Section_A},
	)
	respSectionA, err := bkgHandler.GetSectionBookings(context.TODO(), r)
	if err != nil {
		t.Fail()
	}
	r = GetBookingReqObj(users[1],
		sources[0],
		destinations[0],
		util.GetDate(5),
		booking.Sections{Section: booking.Section_B},
	)
	respSectionB, err := bkgHandler.GetSectionBookings(context.TODO(), r)
	if err != nil {
		t.Fail()
	}

	//section A had two bookings while section B had 3, Validate same
	if (respSectionA.Bookings == nil || respSectionB.Bookings == nil) || len(respSectionA.Bookings) != 2 || len(respSectionB.Bookings) != 3 {
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
	r := GetBookingReqObj(users[0], sources[0], destinations[0], util.GetDate(5), booking.Sections{Section: booking.Section_A})
	for i := 0; i < totalPossibleBookings; i++ {

		//source & dest are keeping at 0th index for which stations are added in routes (sourceStation_destinationStation)

		resp, err := bkgHandler.CreateBooking(context.TODO(), r)
		if err != nil {
			continue
		}
		actualBookings++
		fmt.Println(resp)
		//bookingIds = append(bookingIds, int64(resp.BookingId))
	}
	receipts, err := bkgHandler.GetUserBookingReceipts(context.TODO(), r)
	if err != nil || receipts == nil || len(receipts.Receipts) != actualBookings{
		t.Fail()
	}
}
