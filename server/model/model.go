package model

import (
	Booking "github.com/iamengg/railyatri/bookingStub"
	booking "github.com/iamengg/railyatri/bookingStub"
)

var (
	//Sections of bogie
	SectionA = 0
	SectionB = 1

	//fix seats per bogie
	SeatsPerBogie = 5

	//Section A , Section B are two availalbe categories of bogies , ex: sleeper, general
	SectionBogies = map[int]int{SectionA: 2, SectionB: 1}

	// Stations registered with railway
	Stations = map[string]struct{}{"London": struct{}{}, "Paris": struct{}{}, "Berlin": struct{}{}, "St.Petersberg": struct{}{}}

	// Fare between source & destination stations
	Fare = map[string]int{"London_Paris": 20, "Paris_Berlin": 30, "London_Berlin": 40}

	// Registered Users, Non registered users first need to register to access RailYatri services
	Users = map[int]booking.User{
		1:  booking.User{FirstName: "Pratik", LastName: "p", Email: "p.b@gmail.com", Id: 1},
		2:  booking.User{FirstName: "Sagar", LastName: "s", Email: "ss@gmail.com", Id: 2},
		3:  booking.User{FirstName: "Vijay", LastName: "v", Email: "sv@gmail.com", Id: 4},
		4:  booking.User{FirstName: "Pratik", LastName: "p", Email: "p.b@gmail.com", Id: 1},
		5:  booking.User{FirstName: "Sagar", LastName: "s", Email: "ss@gmail.com", Id: 2},
		6:  booking.User{FirstName: "Vijay", LastName: "v", Email: "sv@gmail.com", Id: 4},
		7:  booking.User{FirstName: "Pratik", LastName: "p", Email: "p.b@gmail.com", Id: 1},
		8:  booking.User{FirstName: "Sagar", LastName: "s", Email: "ss@gmail.com", Id: 2},
		9:  booking.User{FirstName: "Vijay", LastName: "v", Email: "sv@gmail.com", Id: 4},
		10: booking.User{FirstName: "Pratik", LastName: "p", Email: "p.b@gmail.com", Id: 1},
		11: booking.User{FirstName: "Pratik", LastName: "p", Email: "p.b@gmail.com", Id: 1},
		12: booking.User{FirstName: "Sagar", LastName: "s", Email: "ss@gmail.com", Id: 2},
		13: booking.User{FirstName: "Vijay", LastName: "v", Email: "sv@gmail.com", Id: 4},
		14: booking.User{FirstName: "Pratik", LastName: "p", Email: "p.b@gmail.com", Id: 1},
		15: booking.User{FirstName: "Sagar", LastName: "s", Email: "ss@gmail.com", Id: 2},
		16: booking.User{FirstName: "Vijay", LastName: "v", Email: "sv@gmail.com", Id: 4},
		17: booking.User{FirstName: "Pratik", LastName: "p", Email: "p.b@gmail.com", Id: 1},
		18: booking.User{FirstName: "Sagar", LastName: "s", Email: "ss@gmail.com", Id: 2},
		19: booking.User{FirstName: "Vijay", LastName: "v", Email: "sv@gmail.com", Id: 4},
		20: booking.User{FirstName: "Vijay", LastName: "v", Email: "sv@gmail.com", Id: 4},
	}

	// Contains key, value pair of {TrainNumber, map{section:bogies,}}
	TrainObj = map[int]Train{
		1234: Train{
			Number: 1234,
			Bogies: map[int]int{0: SectionBogies[SectionA], 1: SectionBogies[SectionB]},
		},
		2345: Train{
			Number: 2345,
			Bogies: map[int]int{0: SectionBogies[SectionA], 1: SectionBogies[SectionB]},
		},
	}

	// TrainsBetweenStations is contains train numbers which are
	// actively running from source to destination station 	(sourceStation_destinationStation)
	TrainsBetweenStations = map[string][]int64{
		"London_Paris": {1234, 2345},
	}

	// scheduledTrains map of trainNumbers & there schedule (schedueled for day & time)
	// trains schedueld for dates, {trainNumber, set{weekDays}, time}}
	scheduledTrains = map[int]DaysForRunningTrain{
		1234: DaysForRunningTrain{
			Day:  map[weekday]struct{}{Mon: struct{}{}, Tue: struct{}{}, Wed: struct{}{}},
			Time: "6PM",
		},
		2345: DaysForRunningTrain{
			Day:  map[weekday]struct{}{Mon: struct{}{}},
			Time: "6PM",
		},
	}

	/*
		trains which are starting from specific station to end location
		map[st_end_stations] trains
		then should filter trains for days Train.Day == day of selected date then take this train & so on take all such trains at that day & return array
		what time you want to select from this as well there is vacant seat at preferrdSection
		run trains for weekday, trians
		now return ticket from first available train
		Now also mark the train booking for day+trian = booking id for user
		return all bookingIDs of user
		for modifying take onlyu one bookingId of user
		modification not allowed on specific day from train journey for booking or based on time
	*/
)

// SeatNumber in bogie
type SeatNumber int

// BookingId is id for booking
type BookingId int

// Contains all seats in bogi, and
// NxtAvailable seat shows vacant seat in bogi
type Bogies struct {
	Bogi             []SeatNumber
	NxtAvailableSeat int
}

// Weekday is type to store days when trains are scheduled to run
type weekday int

const (
	Mon weekday = iota
	Tue
	Wed
	Thur
	Fri
	Sat
	Sun
)

// DaysForRunningTrain contains map of days where train is scheduled ,
// and time , train will run at same time for all scheduled days
type DaysForRunningTrain struct {
	Day  map[weekday]struct{}
	Time string
}

// UserBookingDetails is final booking receipt, So contains all details related with booking
type UserBookingDetails struct {
	BookingId        int64           `json:"bookingId"`
	UserId           int             `json:"usrId"`
	SeatNum          int             `json:"seatNum"`
	TrainNumber      int             `json:"TrainNum"`
	Section          Booking.Section `json:"BogiType"`
	Status           status          `json:"status"`
	BookingDateTime  string          `json:"bookingDate"`
	ModifiedDateTime string          `json:"modifiedDate"`
	SrcStation       string          `json:"srcStation"`
	DestStation      string          `json:"destStation"`
}

// DB to store bookings per train for each day
type Bookings struct {
	//date+trainNum, map{section+bogi,userBookingDetails}
	BookingsData map[string]map[string][]UserBookingDetails
}

// Stores all bookings of user
type UserBookings struct {
	//map for all bookings of user {userid, {bookingIds,} }
	UserBookingsData map[int]map[BookingId]struct{}
}

// Train is struct of trainNumber & its { sections,bogiesPerSection}
type Train struct {
	Number int64
	Bogies map[int]int //section, bogies
}

// Journey From to To station
type Journey struct {
	From string
	To   string
}

// booking confirmation status
type status int

const (
	CONFIRMED status = iota + 1 // EnumIndex = 1
	NOTCONFIRMED
	CANCELLED
)

type Response struct {
	BookingId  BookingId
	SeatNumber SeatNumber
}

// -----------TODO ----------------------------------------------------------
// add train to railways
func AddTrains(trainNumber int, SectionABogies int, SectionBBogies int) {}

// GetBogies number for section
func GetBogies(section int) {}

// Add stations
func AddStations() {}

// register user
func AddUser() {}

// Methods where dynamically update bogies, trains, stations which are currently taken as constant values

//method to archival of data db
//methods to data cleanup
//Feature to onboard intermediate stations,
//Feature to select multiple bookings & put details for same
//Feature for waitlist & allocating once someone cancelles booking
//Payment option once seat found at for booking, and blocking on receiving payment

//Concurrency -
//simulate multiiple bookings from client with goroutines
