-------------------------------   Requirements

Background: All API referenced are gRPC APIs, not REST ones. 

I want to board a train from London to France. The train ticket will cost $20.  

Create API where you can submit a purchase for a ticket.  Details included in the receipt are: 
a. From, To, User , price paid.
       i. User should include first and last name, email address
The user is allocated a seat in the train.  Assume the train has only 2 sections, section A and section B.
An API that shows the details of the receipt for the user
An API that lets you view the users and seat they are allocated by the requested section
An API to remove a user from the train
An API to modify a user's seat

     --------------------------------- Test cases
Note -> 'F -> ' indicates remaining extension for existing apis
CreateBooking
	Seatnumber should get assigned incrementally section wise
	Are Sent & Recieved pmtr values are matching (user details , stations etc.)
	When booking is existing for same user or users for x train at y date , then don't allow duplicate booking with same user
	If allready section is full booking should not happen -> Put in waiting list by taking payment, If booking not confirmed on journyey day return full payment , Trigger event for same
	F -> If seat came from deleted section then iterate for seat number & book (Or maintain such seats separety for booking )

GetUserReceipts
	Check for not existing user
	Check for existing user but no booking record
	Check for existing user with booking/s
	F -> Get receipts based on date, source, destination

Get all bookings for section
	Check if bookings not exist for section
	Check all bookings for section where booking exist.

Update booking
	- try updating not existing booking id
	- try updating existing booking id
	F ->User with booked user (who booked same ticket previously) only can trigger update booking
		Can modify booking pmtrs related with User

Delete booking
	- Delete booking & see that number of available seats are increased by 1
	- Next time the deleted seat should be available to book
 

	------------------------ Features to extend Platform
Feature Reservation update:
	If reservation deleted or cancelled then provide this seat to CreateBooking thr. separate
		cache as this will be used only when total seat bookeing counter is crosses total physical seats but misses cancelled seats

Features to extend/add/modify infra
	Modifying Sections
	Add/Remove num of bogies to train
	Add/Remove num of seats per bogie

Feature extension : Notifications
	Notification svc if train got cancelled or is getting delayed
	Mail API to send ticket to users mail
	Message to user on journey date, few hours before starting journye , and half hour before offboardign at destination station

Features for Concurrency :
	When millions of users will try to book seat
		If any user selected DS for booking then lock & release after booking
	ReadOnly mutexes for
		- getSectionWise bookings
		- get users all bookings
		- Locks for DB for write/modify API's
			- CreateBooking, UpdateBooking, DeleteBooking