
## Easy Errand (aka 'iizi errand') 
- Congratulations on your startup idea!
  
    Your Local Errand Service app can be a game-changer for people who need help with errands, tasks, and chores.  Here's a description of the app and a potential revenue model to help you earn $100 weekly:
    App Name: ErrandEase

- Tagline: "Get your errands done with ease!"
    
-   Description: ErrandEase is a mobile app that connects users with local errand runners who can help with         various tasks, such as:
    Grocery shopping
    Dog walking
    Household chores
    Delivery services
    And more!

        -  TaskModel {
            -   Location
            -   Budget
            -   Category
            -   TimeReq
            -   Description
            -   TaskRequirements
            -   UserID string
            -   CreatedAt time.Time
        -  }

        -  RatingModel {
            -   UserId string
            -   Rating  float
            -   Review  string
            -   CreatedAt  time.Time
        -  }

### Endpoints
* USER and ERRAND RUNNER [  USER  ]
- Create User ("/user")
- Get User with UserProfile ("/user/:id")
- Get Errand UserProfile with ErrandRunnerProfile ("/user/:id")
- Update Errand UserProfile with ErrandRunnerProfile ("/user/:id")
- Delete Errand UserProfile with ErrandRunnerProfile ("/user/:id")

- Login ("/login") // All users will have the same user data BUT different Profile

* RATING
- Create Rating ("/rating")
- Get Rating ("/rating/:id")
- Delete Rating ("/rating/:id")
- Update Rating ("/rating/:id")


* TASK
- Create Task ("/task")
- Get Task ("/task/:id")
- Get All Tasks ("/tasks")
- Get Most Matched Tasks ("/tasks/matched") [for this task, compare task attributes to errand runner attributes and rank it higher to the best fit errand runner based on the benchmarks like experience, ratings, availability, etc]
- Update Task ("/task/:id")
- Delete Task ("/task/:id")


- 
* 


* How do I determine when an errand runner has got the job?

### Features:
- User registration and profile creation
        -   UserModel {
          <!-- -   UserID -->
          -   FirstName
          -   LastName
          -   Email
          -   UserType  [USER, ERRAND]
          -   Password
          -   JWTToken
          -   CreatedAt
        -   }

        -  UserProfile {
           -  UserID
           -  FirstName
           -  LastName
           -  Phone
           -  Email
           -  Rating *[]RatingModel // Errand Ruuner will update this
           -  Location  // [Name, GPS Address]
           -  UserType  
           -  Tasks *[]TaskModel
           -  NationalID [Ghana Card{preferred}, VotersID, Drivers License]
           -  CreatedAt    time.Time
           -  UpdatedAt
        -  }
- Task posting with details and budget
- Errand runner registration and profile creation
        -   ErrandRunnerModel {
          -   FirstName
          -   LastName
          -   Email
          -   UserType  [USER, ERRAND]
          -   Password
          -   JWTToken
          -   CreatedAt    time.Time
        -   }

        -  ErrandRunnerProfile {
           -  UserID
           -  FirstName
           -  LastName
           -  Phone
           -  Email
           -  UserType  
           -  Tasks *[]TaskModel
           -  Location  // Location [Name, GPS Address]
           -  NationalID  // update profile with ...
           -  Guarantor // update profile with ...
           -  GuarantorPhone // update profile with ...
           -  AvailableTime  // update profile with ...
           -  Rating  *[]RatingModel // update profile with ...
           -  Skills []string // update profile with ...
           -  Photo [Passport Size] // update profile with ...
           -  CreatedAt    time.Time
           -  UpdatedAt     time.Time
        -  }
- Matching algorithm for task assignment
- In-app messaging and task management 
- Payment processing and transaction management
- Rating and review system


    `func (eas *errandApplicationService) CreateApplication(ctx *fiber.Ctx) error {
        // Insert application into database
    }`

    `func (eas *errandApplicationService) GetApplicationsByErrandID(ctx *fiber.Ctx) error {
        // Retrieve applications for errand from database
    }`

    `func (eas *errandApplicationService) UpdateApplicationStatus(ctx *fiber.Ctx) error {
        // Update application status in database
    }`


### Revenue Model:
* Commission-based: Charge a commission (e.g., 15%) on each task completed through the app.
* Service fees: Offer additional services like priority task assignment, extra support, or task insurance for a flat fee.
* Advertising: Display local business ads within the app, targeting users based on their task preferences.

### Potential Earnings:
- Assuming an average task value of $20 and a 15% commission:
    10 tasks completed weekly = $200 in total task value
    15% commission = $30 (your earnings)

### To reach this goal, focus on:
* Building a strong user base through marketing and promotions
* Partnering with local businesses for advertising and task opportunities
* Ensuring excellent user experience and high-quality errand runners


# By scaling your app and expanding services, you can increase earnings beyond $100 weekly. Best of luck with ErrandEase!

-----------------------------------------------------------------------------------------------------------
### Describe an app that can be used as a startup to be making me $100 every week. App Concept: Local Errand Service.

