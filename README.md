# PACE HOMETEST
Create backend APIs that manage merchant accounts.

## Requirements
1. CRUD for merchant accounts.<br/>
    Merchant code should be unique.
2. CRUD for a team member on a merchant account.<br/>
    Team members should have unique email addresses.
3. Get the list of team members per merchant.
4. Get list should implement pagination.

## System Design
```
						   /-----------> merchant_services_(CRUD)-----\	
	  user --------> router											---> Database
						   \-----------> member_services_(CRUD)-------/
```
### APIs
The HTTP server employs the ```gorilla/mux``` package to create a request router. Here are the list of handler methods supported by the server:
```
	router := mux.NewRouter()

	// CRUD for merchant accounts
	router.HandleFunc("/api/newmerchant", middleware.CreateMerchant).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/merchant/{id}", middleware.GetMerchant).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/merchant/{id}", middleware.UpdateMerchant).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletemerchant/{id}", middleware.DeleteMerchant).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/merchant/{id}/allmembers", middleware.GetMembersWithPagination).Methods("GET", "OPTIONS")

	// CRUD for team members of a merchant account
	router.HandleFunc("/api/newmember", middleware.CreateMember).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/member/{memberid}", middleware.GetMember).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/member/{memberid}", middleware.UpdateMember).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletemember/{memberid}", middleware.DeleteMember).Methods("DELETE", "OPTIONS")
```
The statements above clearly define CRUD methods for merchant accounts and their team members with appropriate HTTP type of requests. Note that for getting members with pagination, two options should be set, namely, ```pageid``` and ```pagesize```. For example, getting a list of team members of merchant 1 with ```pageid=1``` and ```pagesize=10``` (will return the first 10 members of merchant 1):
```
curl -X GET "http://localhost:8080/api/merchant/1/allmembers?pageid=1&pagesize=10"
```
### Data models
Two tables including ```merchants``` and ```members``` are created as follows:
```
merchants (
    merchantID SERIAL,
    name VARCHAR(255) NOT NULL,
    age INT,
    location TEXT,
    PRIMARY KEY (merchantID)
);
members (
    memberID SERIAL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(320) NOT NULL,
    merchantID INT,
    PRIMARY KEY (memberID),
    FOREIGN KEY (merchantID) REFERENCES merchants(merchantID) ON UPDATE CASCADE ON DELETE CASCADE
);
```
The cascade policy is defined as a foreign-key constrain between the two tables. Any update/delete on one table will affect the other.

## Installation

### Prerequisites
- Go version 1.16 or later. Please follow this [guide](https://go.dev/doc/install) for the Go installation.
- PostgreSQL with the latest version. Please follow [postgresql.org](https://www.postgresql.org/download) or [digitalocean](https://www.digitalocean.com/community/tutorials/how-to-install-postgresql-on-ubuntu-22-04-quickstart) for the installation guide.

### Clone ```pace-hometest``` project from github
Assume after the installation, ```~go/src``` repo contains projects of Go souce code. Let's clone the ```pace-hometest``` project into it.
```
$ cd ~/go/src
~/go/src$ git clone https://github.com/nguyenvietan94/pace-hometest.git
```
Move into this project repo and download all the required packages:
```
~/go/src$ cd pace-hometest/
~/go/src/pace-hometest$ go get ./...
```

### Create databases and tables with PostgreSQL
Switch to ```postgres``` user:
```
$ sudo -u postgres -i
```
Setup password of PostgreSQL database. In this project, the default password is set to ```postgres```.
```
postgres$ psql
postgres=# \password
Enter new password:
Enter it again:
```
Create ```hometest``` database and switch to this database:
```
postgres=# CREATE DATABASE hometest;
CREATE DATABASE
postgres=# \c hometest
You are now connected to database "hometest" as user "postgres".
hometest=#
```
Create ```merchants``` and ```members``` tables by running ```schemas.sql``` script.
```
hometest=# \i /absolute/path/to/schemas.sql
psql:/home/annv/go/src/pace-hometest/models/schemas.sql:1: NOTICE:  table "merchants" does not exist, skipping
DROP TABLE
psql:/home/annv/go/src/pace-hometest/models/schemas.sql:2: NOTICE:  table "members" does not exist, skipping
DROP TABLE
CREATE TABLE
CREATE TABLE
CREATE INDEX
hometest=# \dt
           List of relations
 Schema |   Name    | Type  |  Owner
--------+-----------+-------+----------
 public | members   | table | postgres
 public | merchants | table | postgres
(2 rows)

```
Or use PostgreSQL command editor to run the script:
```
hometest=# \e # paste the schemas.sql script in the prompted editor
```

Take a look at the configuration file ```.env``` for database connection in Go code. Make sure the password is set correctly:
```
POSTGRES_URL="host=localhost user=postgres password=postgres dbname=hometest sslmode=disable"
```
## Run and Test

As the PostgreSQL database connection is required, it is more convenient to test the program by sending HTTP requests via tools (eg. [Restman](https://chrome.google.com/webstore/detail/restman/ihgpcfpkpmdcghlnaofdmjkoemnlijdi?hl=en) or [cURL](https://curl.se/)), instead of writing unit tests.

These tests are run on a local host.

### Run HTTP server
Let's run the Go code:
```
~/go/src/pace-hometest$ go run main.go
2022/08/08 16:16:06 Successfully connected to database!
2022/08/08 16:16:06 Started server on port 8080
```
**Note**: At this step, I ran into the problem ```pq: password authentication failed for user```. I did try many solutions [here](https://stackoverflow.com/questions/55038942/fatal-password-authentication-failed-for-user-postgres-postgresql-11-with-pg), but only this one worked for me: make sure the PostgreSQL password is explicitly set up (I set ```postgres``` as default):
```
postgres=# \password
Enter new password:
Enter it again:
```
and the password must be correcly declared in the ```.env``` file:
```
~/go/src/pace-hometest$ cat .env
POSTGRES_URL="host=localhost user=postgres password=postgres dbname=hometest sslmode=disable"
```

### CRUD on merchant accounts
HTTP requests will be sent from a seperate terminal to the server by ```cURL``` command.

#### Create
Send ```POST``` requests to create new merchants:
```
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"Alice", "age":18, "location":"US"}' http://localhost:8080/api/newmerchant
{"id":1,"message":"Merchant created successfully"}
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"Bob", "age":20, "location":"Singapore"}' http://localhost:8080/api/newmerchant
{"id":2,"message":"Merchant created successfully"}
```
Logs showing on the server terminal:
```
2022/08/08 16:37:22 Inserted a new merchant: id=1, name=Alice, age=18, location=US
2022/08/08 16:39:09 Inserted a new merchant: id=2, name=Bob, age=20, location=Singapore
```
#### Read
Send ```GET``` requests to receive info of merchants.
```
$ curl -X GET http://localhost:8080/api/merchant/1
{"merchantID":1,"name":"Alice","age":18,"location":"US"}
$ curl -X GET http://localhost:8080/api/merchant/2
{"merchantID":2,"name":"Bob","age":20,"location":"Singapore"}
```
Logs showing on the server terminal:
```
2022/08/08 16:47:31 Get merchant: merchantID=1, name=Alice, age=18, merchantID=US
2022/08/08 16:47:32 Get merchant: merchantID=2, name=Bob, age=20, merchantID=Singapore
```
#### Update
Send ```PUT``` requests to update merchants:
```
$ curl -X PUT -H "Content-Type: application/json" -d '{"name":"Bob", "age":25, "location":"China"}' http://localhost:8080/api/merchant/2
{"id":2,"message":"Updated successfully"}
:~$ curl -X GET http://localhost:8080/api/merchant/2
{"merchantID":2,"name":"Bob","age":25,"location":"China"}
```
Logs showing on the server terminal:
```
2022/08/08 16:58:19 merchantID=2, Updated successfully
2022/08/08 16:58:24 Get merchant: merchantID=2, name=Bob, age=25, merchantID=China
```
#### Delete
Send ```DELETE``` requests to delete merchants.
```
$ curl -X DELETE  http://localhost:8080/api/deletemerchant/2
{"id":2,"message":"Deleted a merchant successfully"}
$ curl -X GET http://localhost:8080/api/merchant/2
{"id":-1,"message":"sql: no rows in result set"}
```
Logs showing on the server terminal:
```
2022/08/08 17:01:06 Deleted a merchant successfully, merchantID 2
2022/08/08 17:01:13 Unable to scan the row. sql: no rows in result set
2022/08/08 17:01:13 Unable to get merchant. sql: no rows in result set
```

### CRUD on team members of a merchant account
HTTP requests will be sent from a seperate terminal to the server by cURL command.

#### Create
Send ```POST``` requests to create new team members of merchant 1.
```
curl -X POST -H "Content-Type: application/json" -d '{"name":"person1", "email":"person1@gmail.com", "merchantID":1}' http://localhost:8080/api/newmember
{"id":1,"message":"Member created successfully."}
```
New members with duplicate emails will be aborted:
```
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"person2", "email":"person1@gmail.com", "merchantID":1}' http://localhost:8080/api/newmember
{"id":-1,"message":"Unable to create a new member. email already exists"}
```
Let's create three more team members of merchant 1:
```
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"person2", "email":"person2@gmail.com", "merchantID":1}' http://localhost:8080/api/newmember
{"id":2,"message":"Member created successfully."}
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"person3", "email":"person3@gmail.com", "merchantID":1}' http://localhost:8080/api/newmember
{"id":3,"message":"Member created successfully."}
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"person4", "email":"person4@gmail.com", "merchantID":1}' http://localhost:8080/api/newmember
{"id":4,"message":"Member created successfully."}
```
Logs showing on the server terminal:
```
2022/08/08 17:17:10 MerchantID 1: Inserted a new team member with memberID 1
2022/08/08 17:17:10 Member created successfully. memberID=1, name=person1, email=person1@gmail.com, merchantID=1
2022/08/08 17:17:34 Unable to create a new member. email already exists
2022/08/08 17:20:33 MerchantID 1: Inserted a new team member with memberID 2
2022/08/08 17:20:33 Member created successfully. memberID=2, name=person2, email=person2@gmail.com, merchantID=1
2022/08/08 17:20:44 MerchantID 1: Inserted a new team member with memberID 3
2022/08/08 17:20:44 Member created successfully. memberID=3, name=person3, email=person3@gmail.com, merchantID=1
2022/08/08 17:20:51 MerchantID 1: Inserted a new team member with memberID 4
2022/08/08 17:20:51 Member created successfully. memberID=4, name=person4, email=person4@gmail.com, merchantID=1
```

#### Read
Send ```GET``` requests to receive members' info.
```
$ curl -X GET http://localhost:8080/api/member/1
{"memberID":1,"name":"person1","email":"person1@gmail.com","merchantID":1}
```
Logs showing on the server terminal:
```
2022/08/08 17:22:54 Get member: memberID=1, name=person1, email=person1@gmail.com, merchantID=1
```
#### Update
Send ```PUT``` requests to update members.
```
$ curl -X PUT -H "Content-Type: application/json" -d '{"name":"Tony", "email":"tony@gmail.com", "merchantID":1}' http://localhost:8080/api/member/1
{"id":1,"message":"Member updated successfully."}
$ curl -X GET http://localhost:8080/api/member/1
{"memberID":1,"name":"Tony","email":"tony@gmail.com","merchantID":1}
```
Logs showing on the server terminal:
```
2022/08/08 17:26:57 memberID=1, Member updated successfully.
2022/08/08 17:27:08 Get member: memberID=1, name=Tony, email=tony@gmail.com, merchantID=1
```
#### Get the list of team members per merchant with pagination
Send a ```get``` request with ```pageid=1``` and ```pagesize=2```. The first two members (out of four) will be transfered:
```
curl -X GET "http://localhost:8080/api/merchant/1/allmembers?pageid=1&pagesize=2"
[{"memberID":2,"name":"person2","email":"person2@gmail.com","merchantID":1},{"memberID":3,"name":"person3","email":"person3@gmail.com","merchantID":1}]
```
Send a ```get``` request with ```pageid=2``` and ```pagesize=2```. The last two members will be transfered:
```
$ curl -X GET "http://localhost:8080/api/merchant/1/allmembers?pageid=2&pagesize=2"
[{"memberID":4,"name":"person4","email":"person4@gmail.com","merchantID":1},{"memberID":1,"name":"Tony","email":"tony@gmail.com","merchantID":1}]
```
Logs showing on the server terminal:
```
2022/08/08 17:36:54 Get team members: merchantID=1, pageID=1, pageSize=2
2022/08/08 17:37:17 Get team members: merchantID=1, pageID=2, pageSize=2
```
#### Delete
Send ```DELETE``` requests to delete members.

```
$ curl -X DELETE  http://localhost:8080/api/deletemember/3
{"id":3,"message":"Deleted a member successfully"}
$ curl -X GET http://localhost:8080/api/member/3
{"id":-1,"message":"sql: no rows in result set"}
```
Logs showing on the server terminal:
```
2022/08/08 17:43:16 Deleted a member successfully. memberID=3
2022/08/08 17:43:37 Unable to scan the row. sql: no rows in result set
2022/08/08 17:43:37 Unable to get member, memberID=3. sql: no rows in result set
```

### Note: The Cascade Policy
When creating ```merchants``` and ```members``` tables, the foreign-key constrant is declared as the cascade policy:
```
CREATE TABLE members (
...
merchantID INT,
FOREIGN KEY (merchantID) REFERENCES merchants(merchantID) ON UPDATE CASCADE ON DELETE CASCADE
);
```
This disallows creating a team member with ```merchantID``` not existing in ```merchants``` table. And when deleting a merchant, all the team member associated with this merchant will be deleted as well.

For example, merchant 2 does not exist in the ```merchants``` table, hence no member with ```merchantID=2``` can be created:
```
$ curl -X GET http://localhost:8080/api/merchant/2
{"id":-1,"message":"sql: no rows in result set"}
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"Alice", "email":"alice@gmail.com", "merchantID":2}' http://localhost:8080/api/newmember
{"id":-1,"message":"Unable to create a new member. pq: insert or update on table \"members\" violates foreign key constraint \"members_merchantid_fkey\""}
```
Removing merchant 1 will result in removing all of its members as well:
```
$ curl -X GET "http://localhost:8080/api/merchant/1/allmembers"
[{"memberID":2,"name":"person2","email":"person2@gmail.com","merchantID":1},{"memberID":4,"name":"person4","email":"person4@gmail.com","merchantID":1},{"memberID":1,"name":"Tony","email":"tony@gmail.com","merchantID":1}]
$ curl -X DELETE  http://localhost:8080/api/deletemerchant/1
{"id":1,"message":"Deleted a merchant successfully"}
$ curl -X GET "http://localhost:8080/api/merchant/1/allmembers"
null
```
