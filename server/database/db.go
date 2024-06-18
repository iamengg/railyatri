package database

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"log"
	"sync"
	"time"

	_ "github.com/iamengg/railyatri/bookingStub"

	//"github.com/iamengg/railyatri/server/model"
	model "github.com/iamengg/railyatri/server/model"
)

// Data structure to store Booking information
var BookingData model.Bookings
var UserBookingsDB model.UserBookings
var BookingIdBookingDetail map[model.BookingId]model.UserBookingDetails

// mutex to secure access to resource
var mt sync.Mutex

// This method called first before any other function in this file by go runtime, So used to initialize
func init() {
	BookingData.BookingsData = make(map[string]map[string][]model.UserBookingDetails)
	UserBookingsDB.UserBookingsData = make(map[int]map[model.BookingId]struct{})
	BookingIdBookingDetail = make(map[model.BookingId]model.UserBookingDetails)
}

// Add userBookings to track all bookings of usre with quick access
func AddUserBooking(userID int, bookingId model.BookingId) {
	bookings, exist := UserBookingsDB.UserBookingsData[userID]
	if !exist {
		bookings = make(map[model.BookingId]struct{})
	}
	bookings[bookingId] = struct{}{}
	UserBookingsDB.UserBookingsData[userID] = bookings
}

// Returns db of all users bookings w.r.t {source_destination_date_of_journey, {section, []booked seats}}
func GetBookingData() *map[string]map[string][]model.UserBookingDetails {
	return &BookingData.BookingsData
}

// Todo: add functionality to check if trains are running at Day corresponding to date
func GetAvailableTrains(srcStation, destinationStation string, date string) []int64 {
	trains, ok := model.TrainsBetweenStations[srcStation+"_"+destinationStation]
	if !ok {
		return []int64{}
	}
	return trains
}

// Create new booking based on availability
// Here we are using only confirmed or notConfirmed status , we are not using waitlist or RAC
func CreateBooking(UserId int64, TrainNum int32, SourceStation string, DestinationStation string,
	Section int, Date string) (model.BookingId, model.SeatNumber, error) {

	//check if user is existing
	log.Println("Section is ", Section)

	if !IsUserExist(UserId) {
		err := errors.New("user is not existing, Create it first before booking")
		log.Println(err, UserId)
		return -1, -1, err
	}

	//check if train is available for Date, Source to Destination station
	trainsRunningOnDate := GetAvailableTrains(SourceStation, DestinationStation, Date)
	if len(trainsRunningOnDate) == 0 {
		err := fmt.Errorf("no trains are availalbe at %s, dest %s, date %s ", SourceStation, DestinationStation, Date)
		log.Println(err, UserId)
		return -1, -1, err
	}

	var NxtAvailableSeat int
	var totalSeats int
	var bogiLen int

	//check if seat with expected section is available at train
	//if not return relavant error message
	for _, TrainNum := range trainsRunningOnDate {

		//for trian number & date key check how many seats allocated in specific section of booking
		//get total seats present at that train for requested section
		//compare with vacant with present seats
		//if yes then get next vacant seat
		//update booked seats count
		sections := BookingData.BookingsData[Date+"_"+fmt.Sprintf("%v", TrainNum)]
		actualSeatsInSection, _ := model.TrainObj[int(TrainNum)]
		totalSeats = 0
		for section, bogies := range actualSeatsInSection.Bogies {
			if section == Section {
				totalSeats = bogies * model.SeatsPerBogie
				break
			}
		}

		for section, bogieSeats := range sections {
			log.Println("Total bookings in train ", TrainNum, " at section ", section, " are ", len(bogieSeats))
			// Split string & get section
			converted, err := strconv.Atoi(strings.Split(section, "_")[0])
			if err != nil {
				err := fmt.Errorf("error in string to int conversion for section index %s", err.Error())
				log.Fatal(err, UserId)
			}

			if converted != Section {
				continue
			}

			//Get next available seatNum
			bogiLen = len(bogieSeats)
			if bogiLen < totalSeats {
				log.Println("Total seats are ", totalSeats, " currently allocated are ", bogiLen)
				NxtAvailableSeat = bogiLen + 1
				break
			} else {
				return -1, -1, errors.New("no seats availalbe to book")
			}
		}
	}

	dateTrainNum, exist := BookingData.BookingsData[Date+"_"+fmt.Sprintf("%d", TrainNum)]
	bookingKey := fmt.Sprintf("%v", Section)
	bookingId := GetBookingId()
	receipt := model.UserBookingDetails{
		BookingId:        bookingId,
		UserId:           int(UserId),
		SeatNum:          NxtAvailableSeat,
		CoachType:        int(Section),
		Status:           model.CONFIRMED,
		BookingDateTime:  time.Now().String(),
		ModifiedDateTime: "",
	}
	if !exist {
		BookingData.BookingsData[Date+"_"+fmt.Sprintf("%d", TrainNum)] = map[string][]model.UserBookingDetails{
			bookingKey: []model.UserBookingDetails{receipt},
		}
	} else {
		allBookingsAtSection := dateTrainNum[bookingKey]
		allBookingsAtSection = append(allBookingsAtSection, receipt)
		dateTrainNum[bookingKey] = allBookingsAtSection

	}
	BookingIdBookingDetail[model.BookingId(bookingId)] = receipt
	AddUserBooking(int(UserId), model.BookingId(bookingId))
	log.Println("Created booking at Database")
	return model.BookingId(bookingId), model.SeatNumber(NxtAvailableSeat), nil
}

func GetUserBookings(UserId int64, TrainNum int32, SourceStation string, DestinationStation string,
	Section int, Date string) ([]model.Response, error) {

	// validate if userId exist
	if !IsUserExist(UserId) {
		err := errors.New("user is not existing, Create it first before booking ")
		log.Println(err, UserId)
		return []model.Response{}, err
	}
	//check for bookings
	bookings := UserBookingsDB.UserBookingsData[int(UserId)]

	userBookigs := make([]model.Response, 0, 5)
	//make []pair{bookigId, seatnumber}
	for bookingIdNumbers, _ := range bookings {
		receipt, exist := BookingIdBookingDetail[model.BookingId(bookingIdNumbers)]
		if !exist {
			continue
		}
		userBookigs = append(userBookigs, model.Response{
			BookingId:  model.BookingId(bookingIdNumbers),
			SeatNumber: model.SeatNumber(receipt.SeatNum)})

	}
	return userBookigs, nil
}

// create timestamp based unique booking id, here we are using lock so
// even at concurrent acces each booking request gets incrementing booking ids
func GetBookingId() int64 {
	mt.Lock()
	defer mt.Unlock()
	return time.Now().UnixNano()
}

func IsUserExist(userId int64) bool {
	if _, exist := model.Users[int(userId)]; !exist {
		return false
	}
	return true
}

func GetSectionBookings(UserId int64, TrainNum int32, SourceStation string, DestinationStation string,
	Section int, Date string) ([]model.Response, error) {

	bookings := BookingData.BookingsData[Date+"_"+fmt.Sprintf("%d", TrainNum)]
	sectionWiseBookings, exist := bookings[fmt.Sprintf("%d", Section)]
	if !exist {
		return []model.Response{}, errors.New("Section not having any bookings")
	}

	sectionBookigs := make([]model.Response, 0, 5)

	for _, receipt := range sectionWiseBookings {
		sectionBookigs = append(sectionBookigs, model.Response{
			BookingId:  model.BookingId(receipt.BookingId),
			SeatNumber: model.SeatNumber(receipt.SeatNum)})

	}
	return sectionBookigs, nil
}

// TODO :  pass bookingId
func DeleteUserBookings(UserId int64, TrainNum int32, SourceStation string, DestinationStation string,
	Section int, Date string) error {
	bookingId := 123 // pass this thr input
	//validate if userId exist
	if !IsUserExist(UserId) {
		err := errors.New("User is not existing, Create it first before booking ")
		log.Println(err, UserId)
		return err
	}

	deleteUserBooking(int(UserId), bookingId)
	deleteFromMapping(model.BookingId(bookingId))
	deleteFromMainDB()
	return nil
}

func deleteUserBooking(userId int, bookingId int) {
	userBookings, exist := UserBookingsDB.UserBookingsData[userId]
	if !exist {
		return
	}
	delete(userBookings, model.BookingId(bookingId))
	UserBookingsDB.UserBookingsData[userId] = userBookings
}

func deleteFromMapping(bookingId model.BookingId) {
	delete(BookingIdBookingDetail, bookingId)
}

func deleteFromMainDB() {
	//BookingData.BookingsData
	log.Fatal("deleteFromMainDB Not implemented")
}

// TODO : pass bookingId
func UpdateUserBooking(UserId int64, TrainNum int32, SourceStation string, DestinationStation string,
	Section int, Date string) error {
	//bookingId := 123 // pass this thr input
	//validate if userId exist
	if !IsUserExist(UserId) {
		err := errors.New("User is not existing, Create it first before booking ")
		log.Println(err, UserId)
		return err
	}

	// UpdateUserBooking(int(UserId), bookingId)
	// deleteFromMapping(model.BookingId(bookingId))
	// deleteFromMainDB()
	return nil
}
