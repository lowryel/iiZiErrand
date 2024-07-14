package models

import "time"

type UserModel struct {
	UserId    string    `json:"user_id"`
	FirstName string    `xorm:"not null" json:"first_name"`
	LastName  string    `xorm:"not null" json:"last_name"`
	Email     string    `xorm:"unique not null" json:"email"`
	UserType  string    `xorm:"not null" json:"user_type" enum:"USER, ERRAND"` // [USER, ERRAND]
	Password  string    `xorm:"not null" json:"password"`
	JwtToken  string    `json:"jwt_token"`
	CreatedAt time.Time `json:"created_at"`
}



// Define the TaskStatus enum-like type
type TaskStatus string

// Define constants for the different ticket types
const (
    Created TaskStatus = "CREATED"
    Assigned        TaskStatus = "ASSIGNED"
    Completed          TaskStatus = "COMPLETED"
)

type TaskModel struct{
	TaskId	string `json:"task_id" xorm:"pk"`
	Longitude	string	`xorm:"not null" json:"longitude"`
	Latitude	string	`xorm:"not null" json:"latitude"`
	Budget		string	`xorm:"not null" json:"budget"`
	Category	string	`xorm:"not null" json:"category"`
	TimeReq		string	`json:"time_req"`
	Status 		TaskStatus		`xorm:"'status'" default:"CREATED"`
	Description		string	`json:"description"`
	TaskRequirements	[]string	`json:"task_requirements"`
	UserId string	`json:"user_id"`
	ErrandRunnerId	string  `json:"errand_runner_id"`
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}


type Location struct {
    Latitude  string `json:"latitude"`
    Longitude string `json:"longitude"`
}



type RatingModel struct {	
    RatingId    string    `json:"rating_id"`
    EmployerId  string    `json:"employer_id"` // id from url id = ctx.Params(user_id)
    RunnerId    string    `json:"runner_id"`  // must be logged in to get the id
    TaskId      string    `json:"task_id"`	// get task is by the employer id
    Rating      float64   `json:"rating"`	
    Review      string    `json:"review"`
    CreatedAt   time.Time `json:"created_at"`
}



type UserProfile struct{
	UserId	string	`json:"user_id" xorm:"pk"`
	FirstName	string `json:"first_name"`
	LastName	string `json:"last_name"`
	Phone	string `json:"phone"`
	Email	string `json:"email"`
	Rating []*RatingModel  `json:"rating"`// Errand Ruuner will update this
	UserType  string `json:"user_type"`
	Longitude	string	`xorm:"not null" json:"longitude"`
	Latitude	string	`xorm:"not null" json:"latitude"`
	Tasks []*TaskModel
	NationalId string `json:"national_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`
}


type ErrandRunnerProfile struct{
	UserId	string				`json:"user_id" xorm:"pk"`
	FirstName	string 			`json:"first_name"`
	LastName	string 			`json:"last_name"`
	Phone	string 				`json:"phone"`
	Email	string 				`json:"email"`
	UserType  string 			`json:"user_type"`
	Tasks []*TaskModel 			`json:"tasks"`
	NationalId  string 			`json:"national_id"`// update profile with ...		    /
	Guarantor string 			`json:"guarantor"`// update profile with ...			/
	GuarantorPhone  string 		`json:"guarantor_phone"`// update profile with ...		/
	AvailableTime   string 		`json:"available_time"`// update profile with ...		/
	Longitude	string	`xorm:"not null" json:"longitude"`
	Latitude	string	`xorm:"not null" json:"latitude"`
	Ratings  []*RatingModel 	`json:"ratings"`// update profile with ...
	Skills 	[]string 			`json:"skills"`// update profile with ...				/
	Photo 	string 			`json:"photo"`// update profile with ...
	CreatedAt    time.Time	`json:"created_at"`
	UpdatedAt     time.Time	`json:"updated_at"`
}


type ErrandApplication struct {
	AppId        string   `json:"app_id"`
	TaskId  string   `json:"errand_id"`
	UserId    string   `json:"user_id"`
	Status    string `json:"status"` // "pending", "accepted", "rejected"
	Skills  []string	`json:"skills"`
	Email	string	`json:"email"`
	Description  string	`json:"description"`
	CreatedAt    time.Time	`json:"created_at"`
	UpdatedAt     time.Time	`json:"updated_at"`
}



type Login struct {
	Email string `json:"email"`
	Password string `json:"password"`
}


type LoginData struct {
	FirstName string 
	LastName string
	Email string	`json:"email"`
	UserType string
	Password string	`json:"password"`
}


type ChangePass struct {
	Email string `json:"email"`
	OldPass string `json:"old_pass"`
	NewPass string `json:"new_pass"`
}

