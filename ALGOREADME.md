- The matching algorithm is a crucial component of your ErrandEase app, responsible for assigning tasks to the most suitable errand runners. Here's a high-level overview of a potential matching algorithm:
#### Inputs:
- Task details (e.g., location, category, budget, time requirements)

- Errand runner profiles (e.g., location, skills, availability, ratings)

#### Algorithm Steps:
- Filtering: Narrow down the pool of errand runners based on task requirements (e.g., location, category)
- Ranking: Score remaining errand runners based on their profiles (e.g., ratings, experience[how long and how many jobs have you so far], availability) [Rank the filtered errand runners based on ratings, experience and availability]
- Assignment: Match the task with the top-ranked errand runner

#### Factors to Consider:
- Distance: Prioritize errand runners nearby the task location
- Skills: Match tasks with errand runners having relevant skills or experience
- Availability: Consider errand runners' schedules and availability
- Ratings: Prioritize errand runners with high ratings and positive reviews
- Task history: Consider errand runners' past performance on similar tasks

#### Algorithm Types:
- Simple ranking: Assign tasks based on a single ranking score
- Multi-factor ranking: Use a weighted scoring system to consider multiple factors
- Machine learning: Train a model to predict the best match based on historical data

#### Benefits:
- Efficient task assignment
- Increased user satisfaction
- Improved errand runner utilization
- Enhanced platform reliability
- By developing an effective matching algorithm, you'll create a seamless experience for both users and errand runners, setting your app up for success!

#### Nuked Codes

`
 func (repo *Repository) LoginHandler(ctx *fiber.Ctx) error {
 	loginObj := &models.Login{}
 	err := ctx.BodyParser(loginObj)
 	if err != nil{
 		errorLogger.Println("invalid data")
 		return err
 	}
 	// Begin transaction
 	session := repo.DBConn.NewSession()
 	defer session.Close()
 	err = session.Begin()
 	if err != nil {
 		errorLogger.Println("session error", err)
 		return err
 	}
 	// retrieve registered user object with email
 	user := &models.LoginData{}
 	_, err = repo.DBConn.SQL("select * from user_model where email = ?", loginObj.Email).Get(user)
 	if err != nil {
 		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"error fetching login data"})
 		session.Rollback()
 		errorLogger.Println(err)
 		return err
 	}
 	infoLogger.Println(user)
 	// match the saved hash in login table to the login input password
 	if err = models.CompareHashAndPass(user.Password, loginObj.Password); err != nil{
 		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"error matching password"})
 		errorLogger.Println(err)
 		session.Rollback()
 		return err
 	}
 	if user.UserType == "ERRAND"{
 		infoLogger.Println(user.Email)
 		infoLogger.Println(user.UserType)
 		errand_user := models.UserModel{}
 		_, err = repo.DBConn.Where("email = ? AND user_type = ?", user.Email, user.UserType).Get(&errand_user)
 		if err != nil{
 			errorLogger.Println("failed to retrieve user")
 			session.Rollback()
 			return err
 		}
 		// Generate jwt token
 		token, err := models.GenerateToken(user.Email, user.UserType, errand_user.UserId)
 		if err != nil {
 			session.Rollback()
 			return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"msg":"error in generating token"})
 		}
 		// Commit transaction
 		err = session.Commit()
 		if err != nil {
 			return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
 				"error": "failed transaction",
 			})
 		}
 		infoLogger.Println("Login successful")
 		return ctx.Status(http.StatusCreated).JSON(&fiber.Map{"token": token})
 	}
 	// 
 	client_user := models.UserModel{}
 	infoLogger.Println(user.Email)
 	infoLogger.Println(user.UserType)
 	_, err = repo.DBConn.Where("email = ? AND user_type = ?", user.Email, user.UserType).Get(&client_user)
 	if err != nil{
 		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"msg":"error retrieving user"})
 		session.Rollback()
 		return err
 	}
 	// Generate jwt token
 	token, err := models.GenerateToken(user.Email, user.UserType, client_user.UserId)
 	if err != nil {
 		session.Rollback()
 		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"msg":"error in generating token"})
 	}
 	// Commit transaction
 	err = session.Commit()
 	if err != nil {
 		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
 			"error": "failed transaction",
 		})
 	}
 	infoLogger.Println("Login successful")
 	ctx.Status(http.StatusCreated).JSON(&fiber.Map{"token": token})
 	return err
 }`
