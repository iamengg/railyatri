package database

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"log"
	"sync"
	"time"

	Booking "github.com/iamengg/railyatri/bookingStub"
	_ "github.com/iamengg/railyatri/bookingStub"

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
	if bookingId == -1 {
		return
	}

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
func CreateBooking(UserId int64, SourceStation string, DestinationStation string,
	Section int, Date string) (model.BookingId, model.SeatNumber, error) {

	//check if user is existing
	//log.Println("Section is ", Section)

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
	var TrainNum int64
	var seatFound bool
	//check if seat with expected section is available at train
	//if not return relavant error message
	for _, TrainNum = range trainsRunningOnDate {

		//for trian number & date key check how many seats allocated in specific section of booking
		//get total seats present at that train for requested section
		//compare with vacant with present seats
		//if yes then get next vacant seat
		//update booked seats count
		sections, exist := BookingData.BookingsData[Date+"_"+fmt.Sprintf("%v", TrainNum)]
		if !exist {
			//this train is totally empty for current date
			NxtAvailableSeat = 1
			break
		}
		actualSeatsInSection := model.TrainObj[int(TrainNum)]
		totalSeats = 0
		for section, bogies := range actualSeatsInSection.Bogies {
			if section == Section {
				totalSeats = bogies * model.SeatsPerBogie
				break
			}
		}

		//if section which we looking to book is not available in trains booking means it's empty and available in full capacity
		_, exist = sections[fmt.Sprintf("%d", Section)]
		if !exist {
			NxtAvailableSeat = 1
			break
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
				log.Printf("At train num %d ,total seats at section %s are %d, currently allocated are %d", TrainNum, section, totalSeats, bogiLen)
				NxtAvailableSeat = bogiLen + 1
				seatFound = true
				break
			} else {
				return -1, -1, errors.New("no seats availalbe to book")
			}
		}
		if seatFound {
			break
		}
	}

	dateTrainNum, exist := BookingData.BookingsData[Date+"_"+fmt.Sprintf("%d", TrainNum)]
	bookingKey := fmt.Sprintf("%v", Section)
	bookingId := GetBookingId()
	fare := model.GetFare(SourceStation, DestinationStation)
	//NxtAvailableSeat++
	receipt := model.UserBookingDetails{
		BookingId:        bookingId,
		UserId:           int(UserId),
		SeatNum:          NxtAvailableSeat,
		TrainNumber:      int(TrainNum),
		Section:          Booking.Section(Section),
		Status:           model.CONFIRMED,
		Price:            fare,
		BookingDateTime:  time.Now().String(),
		ModifiedDateTime: "",
		SrcStation:       SourceStation,
		DestStation:      DestinationStation,
		FirstName:        model.Users[int(UserId)].FirstName,
		LastName:         model.Users[int(UserId)].LastName,
		Email:            model.Users[int(UserId)].Email,
	}
	if !exist {
		BookingData.BookingsData[Date+"_"+fmt.Sprintf("%d", TrainNum)] = map[string][]model.UserBookingDetails{
			bookingKey: []model.UserBookingDetails{receipt},
		}
		//dateTrainNum = BookingData.BookingsData[Date+"_"+fmt.Sprintf("%d", TrainNum)]
	} else {
		allBookingsAtSection := dateTrainNum[bookingKey]
		allBookingsAtSection = append(allBookingsAtSection, receipt)
		dateTrainNum[bookingKey] = allBookingsAtSection
	}

	BookingIdBookingDetail[model.BookingId(bookingId)] = receipt
	AddUserBooking(int(UserId), model.BookingId(bookingId))

	return model.BookingId(bookingId), model.SeatNumber(NxtAvailableSeat), nil
}

func GetUserBookingReceipts(UserId int64, TrainNum int32, SourceStation string, DestinationStation string,
	Section int, Date string) (*Booking.BookingReceipts, error) {

	// validate if userId exist
	if !IsUserExist(UserId) {
		err := errors.New("user is not existing, Create it first before booking ")
		log.Println(err, UserId)
		return &Booking.BookingReceipts{}, err
	}
	//check for bookings
	bookings := UserBookingsDB.UserBookingsData[int(UserId)]

	userBookigReceipts := make([]*Booking.BookingReceipt, 0, 5)

	for bookingIdNumbers, _ := range bookings {
		receipt, exist := BookingIdBookingDetail[model.BookingId(bookingIdNumbers)]
		if !exist {
			continue
		}
		if receipt.SrcStation != SourceStation || receipt.DestStation != DestinationStation {
			//this is not booking which user looking for
			continue
		}
		bookingReceipt := &Booking.BookingReceipt{
			BookingId:       receipt.BookingId,
			UserId:          int32(receipt.UserId),
			SrcStation:      receipt.SrcStation,
			DestStation:     receipt.DestStation,
			TrainNumber:     int32(receipt.TrainNumber),
			SeatNum:         int32(receipt.SeatNum),
			Price:           receipt.Price,
			Section:         int32(receipt.Section),
			FirstName:       receipt.FirstName,
			LastName:        receipt.LastName,
			Email:           receipt.Email,
			Status:          GetStatusString(int(receipt.Status)),
			BookingDateTime: receipt.BookingDateTime,
		}
		userBookigReceipts = append(userBookigReceipts, bookingReceipt)
	}

	//return userBookigReceipts, nil
	return &Booking.BookingReceipts{
		Receipts: userBookigReceipts,
	}, nil
}

func GetStatusString(status int) string {
	switch status {
	case 0:
		return "CONFIRMED"
	case 1:
		return "NOTCONFIRMED"
	case 2:
		return "CANCELLED"
	}
	return "PENDING"
}

// Get booking receipt for bookingId
func GetBookingReceipt(bookingId int64) model.UserBookingDetails {
	receipt, exist := BookingIdBookingDetail[model.BookingId(bookingId)]
	if !exist {
		return model.UserBookingDetails{}
	}
	return receipt
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
	Section int, Date string) (*Booking.BookingsResponse, error) {

	bookings := BookingData.BookingsData[Date+"_"+fmt.Sprintf("%d", TrainNum)]
	sectionWiseBookings, exist := bookings[fmt.Sprintf("%d", Section)]
	if !exist {
		return &Booking.BookingsResponse{}, errors.New("section not having any bookings")
	}

	sectionBookigs := make([]*Booking.BookingResponse, 0)

	//	sectionBookigs. = make([]Booking.BookingResponse,0)
	for _, receipt := range sectionWiseBookings {
		sectionBookigs = append(sectionBookigs, &Booking.BookingResponse{
			BookingId:  int64(receipt.BookingId),
			SeatNumber: int32(receipt.SeatNum)})
	}

	return &Booking.BookingsResponse{
		Bookings: sectionBookigs,
	}, nil
}

// TODO :  pass bookingId
func DeleteUserBookings(userId int64, bookingId int64) error {

	//validate if userId exist
	if !IsUserExist(userId) {
		err := errors.New("user is not existing, Create it first before booking ")
		log.Println(err, userId)
		return err
	}
	receipt := GetBookingReceipt(bookingId)
	deleteUserBooking(receipt.UserId, receipt.BookingId)

	deleteFromMapping(model.BookingId(bookingId))
	deleteFromMainDB(receipt.BookingDateTime, receipt.TrainNumber, receipt.Section, receipt.UserId, receipt.BookingId)
	return nil
}

func deleteUserBooking(userId int, bookingId int64) {
	userBookings, exist := UserBookingsDB.UserBookingsData[userId]
	if !exist {
		//bookingId is wrong
		return
	}

	delete(userBookings, model.BookingId(bookingId))
	UserBookingsDB.UserBookingsData[userId] = userBookings
}

func deleteFromMapping(bookingId model.BookingId) {
	delete(BookingIdBookingDetail, bookingId)
}

func deleteElement(data []model.UserBookingDetails, index int) []model.UserBookingDetails {
	return append(data[:index], data[index+1:]...)
}

// TODO : testing delete from db
func deleteFromMainDB(date string, trainNum int, sectionToDel Booking.Section, userId int, bookingId int64) {
	sections, ok := BookingData.BookingsData[date+"_"+fmt.Sprintf("%d", trainNum)]
	if !ok {
		return
	}

	for section, data := range sections {
		if section == fmt.Sprintf("%v", sectionToDel) {
			for index, userBooking := range data {
				if userBooking.BookingId == bookingId {
					data = deleteElement(data, index)
				}
			}
		}
	}
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
