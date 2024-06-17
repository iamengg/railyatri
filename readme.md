2nd Technical Assessment Instructions:
1. Please complete the following 2nd Technical Assessment within 3 calendar days.
2. Once completed, please reply to this email with the GitHub link to your results. This link will be received by your Recruiter and added to the calendar invite for your upcoming Hiring Team interview (coming soon).
3. IMPORTANT: If this Technical Assessment is not completed or received within 3 days, your interview will be cancelled.
4. For Junior roles, all effort applied to this 2nd Technical Assessment will be considered!  We look forward to reviewing your results!

Requirements: 

Code must be published in Github with a link we can access (use public repo).
Code must compile with some effort on unit tests, doesn't have to be 100%, but it shouldn't be 0%.
Please code this with Golang and gRPC
Adding a persistence layer will be cumbersome, so just store the data in your current session/in memory.
The results can be in the console output from your grpc-server and grpc-client
App to be coded:

 

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