package main

import (
	"errors"

	"net/http"

	"time"

	rank "github.com/eugene/iizi_errand"
	"github.com/eugene/iizi_errand/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"xorm.io/xorm"
)

type Repository struct {
	DBConn *xorm.Engine
}

var (
	// api_key = os.Getenv("GEO_API_KEY")
	// url = "https://api.ipgeolocation.io/ipgeo?apiKey=%v&ip=%v"
)

func (r *Repository) CreateUser(ctx *fiber.Ctx) error {
	user := &models.UserModel{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if err := models.ValidateUser(user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	hashedPass, err := models.HashPass(user.Password)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
	}

	user.Password = hashedPass
	user.UserId = uuid.NewString()
	user.CreatedAt = time.Now()

	session := r.DBConn.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to begin transaction"})
	}
	
	/* INSERT DATA INTO DB */
	if err := r.insertUser(user); err != nil {
		session.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if err := session.Commit(); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to commit transaction"})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{"message": "User created successfully", "data": user})
}


func (repo *Repository) GetUser(ctx *fiber.Ctx) error{
	tokenString := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenString)
	if err != nil{
		errorLogger.Println("session error", err)
		return err
	}

	session := repo.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("session error", err)
		return nil
	}
	
	userProfile := &models.UserProfile{}
	user_id := claims.UserId
	infoLogger.Println(user_id)

	if claims.UserType == "USER"{
		has, err := repo.DBConn.Where("user_id = ?", user_id).Get(userProfile)
		if err != nil{
			errorLogger.Println("session error", err)
			session.Rollback()
			return err
		}
		if !has {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User profile not found"})
		}
		if err := session.Commit(); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to commit transaction"})
		}

		return ctx.Status(http.StatusCreated).JSON(fiber.Map{"message": "User created successfully", "data": userProfile})
	}

	errandRunner := &models.UserProfile{}
	has, err := repo.DBConn.Where("user_id = ?", user_id).Get(errandRunner)
	if err != nil{
		errorLogger.Println("session error", err)
		session.Rollback()
		return err
	}
	if !has {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User profile not found"})
	}
	if err := session.Commit(); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to commit transaction"})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{"message": "User created successfully", "data": errandRunner})
}


// 
func (repo *Repository) ChangePasswordHandler(ctx *fiber.Ctx) error {
    changePassObj := &models.ChangePass{}
    if err := ctx.BodyParser(changePassObj); err != nil {
        return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
    }

    if err := models.ValidateChangePassData(changePassObj); err != nil {
        return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    session := repo.DBConn.NewSession()
    defer session.Close()

    if err := session.Begin(); err != nil {
        return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
    }

	loginObj := &models.Login{
		Email: changePassObj.Email,
		Password: changePassObj.OldPass,
	}

    user, err := repo.AuthenticateUser(loginObj)
    if err != nil {
        session.Rollback()
        return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
    }

    if err := repo.UpdateUserPassword(user.Email, changePassObj.NewPass); err != nil {
        session.Rollback()
        return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
    }

    if err := session.Commit(); err != nil {
        return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
    }

    return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Password changed successfully"})
}




func (r *Repository) UpdateUserProfile(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenString)
	if err != nil{
		errorLogger.Println("session error", err)
		return err
	}
	user_id := claims.UserId
	infoLogger.Println(user_id)
	if claims.UserType != "USER"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"can't access this resource",
		})
		return err
	}

	userProfile := &models.UserProfile{}
	err = ctx.BodyParser(userProfile)
	if err != nil {
		errorLogger.Printf("failed to parse data: %v", err)
		return nil
	}
	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("session error", err)
		return nil
	}
	
	update := &models.UserProfile{}

	has, err := r.DBConn.Where("user_id = ?", user_id).Get(update)
	if err != nil{
		errorLogger.Println("session error", err)
		session.Rollback()
		return err
	}
    if !has {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User profile not found"})
    }

	location, err := rank.GetLocation()
	if err != nil {
		errorLogger.Println(err)
		session.Rollback()
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get location"})
	}

	update.NationalId=userProfile.NationalId
	update.Phone=userProfile.Phone
	update.UpdatedAt=time.Now()
	update.Latitude=location.Latitude
	update.Longitude=location.Longitude
    // Update specific fields
	_, err = r.DBConn.ID(user_id).Update(update)
	if err != nil{
		errorLogger.Println("session error", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed transaction",
		})
		return err
	}

	// Commit transaction
	err = session.Commit()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed transaction",
		})
	}
	infoLogger.Println("Update successful")
	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{"msg": "user profile update successful"})
}



func (r *Repository) UpdateErrandRunnerProfile(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenString)
	if err != nil{
		errorLogger.Println("session error: ", err)
		return err
	}
	user_id := claims.UserId
	infoLogger.Println(user_id)
	if claims.UserType != "ERRAND"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"can't access this resource",
		})
		return err
	}

	errandProfile := &models.ErrandRunnerProfile{}
	err = ctx.BodyParser(errandProfile)
	if err != nil {
		errorLogger.Println("failed to parse data")
		return nil
	}
	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("db session error", err)
		return nil
	}
	
	updateErrandModel := &models.ErrandRunnerProfile{}

	// retrieve user data before updating
	has, err := r.DBConn.Where("user_id = ?", user_id).Get(updateErrandModel)
	if err != nil{
		errorLogger.Println("session error", err)
		session.Rollback()
		return err
	}
    if !has {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Errand runner profile not found"})
    }

	location, err := rank.GetLocation()
	if err != nil {
		errorLogger.Println(err)
	}

	updateErrandModel.Latitude = location.Latitude
	updateErrandModel.Longitude = location.Longitude
	updateErrandModel.NationalId = errandProfile.NationalId
	updateErrandModel.Phone = errandProfile.Phone
	updateErrandModel.AvailableTime = errandProfile.AvailableTime
	updateErrandModel.Guarantor = errandProfile.Guarantor

	// errandProfile.Skills [This implementation will: Keep all existing skills. Add any new skills from the update request that aren't already in the list. Avoid duplicates.]
	updateErrandModel.Skills = models.AppendArrayToArray(updateErrandModel.Skills, errandProfile.Skills)
	updateErrandModel.GuarantorPhone = errandProfile.GuarantorPhone
	updateErrandModel.Photo = errandProfile.Photo
	updateErrandModel.UpdatedAt = time.Now()

	_, err = r.DBConn.Where("user_id = ?", user_id).Update(updateErrandModel)
	if err != nil{
		errorLogger.Println("session error", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed transaction",
		})
		return err
	}

	// Commit transaction
	err = session.Commit()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed transaction",
		})
	}
	infoLogger.Println("Update successful")
	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{"msg": "errand runner profile update successful"})
}


func (r *Repository) DeleteErrandRunnerProfile(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenString)
	if err != nil{
		errorLogger.Println("invalid user id", err)
		return err
	}
	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("DB session closed: ", err)
		return err
	}
	user_id := claims.UserId
	infoLogger.Println(user_id)
	if claims.UserType != "ERRAND"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"can't access this resource",
		})
		return err
	}
	errandRunner := &models.ErrandRunnerProfile{}
	_, err = r.DBConn.Where(" user_id = ? ", user_id).Delete(errandRunner)
	if err != nil{
		errorLogger.Println("failed to delete errand runner", err)
		session.Rollback()
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "failed to delete errand runner",
		})
		return err
	}
	err = session.Commit()
	if err != nil{
		errorLogger.Println("transaction commit failed", err)
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "DB session failed",
		})
		return err
	}
	infoLogger.Println("errand profile removed successfully")
	return ctx.Status(204).JSON(&fiber.Map{"msg": "errand runner profile removed"})
}



func (r *Repository) DeleteUserProfile(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenString)
	if err != nil{
		errorLogger.Println("invalid user id", err)
		return err
	}
	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("DB session closed: ", err)
		return err
	}
	user_id := claims.UserId
	infoLogger.Println(user_id)
	if claims.UserType != "USER"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"can't access this resource",
		})
		return err
	}
	user := &models.UserProfile{}
	_, err = r.DBConn.Where(" user_id = ? ", user_id).Delete(user)
	if err != nil{
		errorLogger.Println("failed to delete errand runner", err)
		session.Rollback()
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "failed to delete errand runner",
		})
		return err
	}
	err = session.Commit()
	if err != nil{
		errorLogger.Println("transaction commit failed", err)
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "DB session failed",
		})
		return err
	}
	infoLogger.Println("errand profile removed successfully")
	return ctx.Status(204).JSON(&fiber.Map{"msg": "errand runner profile removed"})
}


/* ERRAND RUNNER RATING USER */
func (r *Repository) RateUser(ctx *fiber.Ctx) error {
	tasker_id := ctx.Params("user_id")

	tokenString := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenString)
	if err != nil{
		errorLogger.Println("session error: ", err)
		return err
	}
	errandRunnerId := claims.UserId
	infoLogger.Println(errandRunnerId)
	if claims.UserType != "ERRAND"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"can't access this resource",
		})
		return err
	}

	rating := &models.RatingModel{}
	err = ctx.BodyParser(rating)
	if err != nil {
		errorLogger.Println("failed to parse rating data")
		return err
	}

	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("session error", err)
		return nil
	}
	
	// retrieve user data before updating
	userProfile, err := r.GetUserProfile(tasker_id)
	if err != nil{
		errorLogger.Println("failed to retrieve user profile", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to retrieve user profile",
		})
		return err
	}

	if userProfile.UserId == ""{
		errorLogger.Println("incorrecr user id")
		return errors.New("incorrect user id")
	}
	
	rating.RatingId = uuid.NewString()
	rating.EmployerId = userProfile.UserId
	rating.RunnerId = errandRunnerId
	rating.CreatedAt = time.Now()
	_, err = r.DBConn.Insert(rating)
	if err != nil{
		errorLogger.Println("failed to create rating", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to create rating",
		})
		return err
	}

	userProfile.Rating = append(userProfile.Rating, rating)
	_, err = r.DBConn.Where(" user_id = ? ", userProfile.UserId).Update(userProfile)
	if err != nil{
		errorLogger.Println("failed to create rating", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to create rating",
		})
		return err
	}

	// Commit transaction
	err = session.Commit()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed transaction",
		})
	}
	infoLogger.Println("user rating created")
	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{"msg": "user rating created"})
}


/* USER RATING ERRAND RUNNER */
func (r *Repository) RateErrandRunner(ctx *fiber.Ctx) error {
	errand_runner_id := ctx.Params("errand_runner_id")

	tokenString := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenString)
	if err != nil{
		errorLogger.Println("session error: ", err)
		return err
	}
	userId := claims.UserId
	infoLogger.Println(userId)
	if claims.UserType != "USER"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"can't access this resource",
		})
		return err
	}

	rating := &models.RatingModel{}
	err = ctx.BodyParser(rating)
	if err != nil {
		errorLogger.Println("failed to parse rating data")
		return err
	}

	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("session error", err)
		return nil
	}
	
	errandRunner := &models.ErrandRunnerProfile{}
	// retrieve user data before updating
	has, err := r.DBConn.Where("user_id = ?", errand_runner_id).Get(errandRunner)
	if err != nil{
		errorLogger.Println("failed to retrieve errand runner profile", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to retrieve errand runner profile",
		})
		return err
	}
    if !has {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Errand runner profile not found"})
    }

	if errandRunner.UserId == ""{
		errorLogger.Println("incorrecr errand runner id")
		return errors.New("incorrect errand runner id")
	}
	
	rating.RatingId = uuid.NewString()
	rating.EmployerId = userId
	rating.RunnerId = errandRunner.UserId
	rating.CreatedAt = time.Now()
	_, err = r.DBConn.Insert(rating)
	if err != nil{
		errorLogger.Println("failed to create rating", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to create rating",
		})
		return err
	}

	errandRunner.Ratings = append(errandRunner.Ratings, rating)
	_, err = r.DBConn.Where(" user_id = ? ", errandRunner.UserId).Update(errandRunner)
	if err != nil{
		errorLogger.Println("failed to create rating", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to create errand runner rating",
		})
		return err
	}

	// Commit transaction
	err = session.Commit()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed transaction",
		})
	}
	infoLogger.Println("errand runner rating created")
	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{"msg": "errand runner rating created"})
}


// user ratings
func (r *Repository) GetUserRatings(ctx *fiber.Ctx) error {
	user_id := ctx.Params("user_id")
	_, err := models.GetIdFromToken(ctx.Get("Authorization"))
	if err != nil{
		errorLogger.Println("session error: ", err)
		return err
	}
	// retrieve user ratings
	user := &[]models.RatingModel{}
	err = r.DBConn.Where("employer_id = ?", user_id).Find(user)
	if err != nil{
		errorLogger.Println("failed to retrieve user profile", err)
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to retrieve user profile",
		})
		return err
	}

	// calculate the average of user.rating
	averageRating := 0.0
	for _, rating := range *user {
		averageRating += rating.Rating
	}
	averageRating = averageRating / float64(len(*user))
	infoLogger.Printf("user rating retrieved %v stars: ", averageRating)
	return ctx.Status(http.StatusOK).JSON(user)
}


/* CREATE TASK */
func (r *Repository) CreateTask(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenString)
	if err != nil{
		errorLogger.Println("session error: ", err)
		return err
	}
	user_id := claims.UserId
	infoLogger.Println(user_id)
	if claims.UserType != "USER"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"can't access this resource",
		})
		return err
	}

	task := &models.TaskModel{}
	err = ctx.BodyParser(task)
	if err != nil {
		errorLogger.Println("failed to parse data")
		return err
	}

	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("session error", err)
		return err
	}

	// get location from a service
	location, err := rank.GetLocation()
	if err != nil{
		errorLogger.Println("session error", err)
		return err
	}

	task.TaskId = uuid.NewString()
	task.Status = models.Created
	task.Latitude = location.Latitude
	task.Longitude = location.Longitude
	task.UserId = user_id
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	// retrieve user data before updating
	_, err = r.DBConn.Insert(task)
	if err != nil{
		errorLogger.Println("session error", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed transaction",
		})
		return err
	}

	// update user profile with task
	updateUserProfileTasks := &models.UserProfile{}
	has, err := r.DBConn.Where(" user_id = ? ", user_id).Get(updateUserProfileTasks)
	if err != nil{
		errorLogger.Println("failed to retrieve user profile", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to retrieve user profile",
		})
		return err
	}

    if !has {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User profile not found"})
    }

	updateUserProfileTasks.Tasks = append(updateUserProfileTasks.Tasks, task)
	_, err = r.DBConn.Where(" user_id = ? ", user_id).Update(updateUserProfileTasks)
	if err != nil{
		errorLogger.Println("failed to retrieve user profile", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to update user profile with task",
		})
		return err
	}

	// Commit transaction
	err = session.Commit()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed transaction",
		})
	}
	infoLogger.Println("task creation successful")
	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{"data": task})
}


// get errand runner id from url param NOT DONE YET
func (r *Repository) TasksNearYou(ctx *fiber.Ctx) error {
    tokenString := ctx.Get("Authorization")
    claims, err := models.GetIdFromToken(tokenString)
    if err != nil {
        errorLogger.Println("session error: ", err)
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
    }

    user_id := claims.UserId
    infoLogger.Println(user_id)

    if claims.UserType != "ERRAND" {
        errorLogger.Println("unauthorized access")
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Can't access this resource"})
    }

    tasks := &[]models.TaskModel{}
    user := &models.ErrandRunnerProfile{}

    session := r.DBConn.NewSession()
    defer session.Close()

    if err := session.Begin(); err != nil {
        errorLogger.Println("session error", err)
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
    }

    if err := r.DBConn.Where(" status=? ", "CREATED").Find(tasks); err != nil {
        errorLogger.Println("session error", err)
        session.Rollback()
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tasks"})
    }

    if _, err := r.DBConn.Where("user_id = ?", user_id).Get(user); err != nil {
        errorLogger.Println("session error", err)
        session.Rollback()
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user data"})
    }

    nearbyTasks := []models.TaskModel{}
    userLocation := &models.Location{
        Latitude:  user.Latitude,
        Longitude: user.Longitude,
    }

    for _, task := range *tasks {
        taskLocation := &models.Location{
            Latitude:  task.Latitude,
            Longitude: task.Longitude,
        }
        score := CalculateDistanceScore(*taskLocation, *userLocation)
        if score <= 10 {
            nearbyTasks = append(nearbyTasks, task)
        }
    }

    if err := session.Commit(); err != nil {
        errorLogger.Println("commit error", err)
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
    }

    infoLogger.Println("tasks near you")
    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": nearbyTasks, "msg": "Tasks near you"})
}


// 
func (r *Repository) UpdateTask(ctx *fiber.Ctx) error {
	task_id := ctx.Params("task_id")
	tokenString := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenString)
	if err != nil{
		errorLogger.Println("session error: ", err)
		return err
	}

	user_id := claims.UserId
	infoLogger.Println(user_id)
	if claims.UserType != "USER"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"can't access this resource",
		})
		return err
	}

	task := &models.TaskModel{}
	err = ctx.BodyParser(task)
	if err != nil {
		errorLogger.Println("failed to parse data")
		return err
	}

	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("session error", err)
		return err
	}

	// retrieve user data before updating
	has, err := r.DBConn.Where(" task_id = ? ", task_id).Get(task)
	if err != nil{
		errorLogger.Println("failed to retrieve task", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to retrieve task",
		})
		return err
	}

    if !has {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "TaskModel not found"})
    }

	// update task
	updateTask := &models.TaskModel{
		Latitude: task.Latitude,
		Longitude: task.Longitude,
		Budget: task.Budget,
		TimeReq: task.TimeReq,
		Category: task.Category,
		Description: task.Description,
		UpdatedAt: task.UpdatedAt,
		
	}
	updateTask.TaskRequirements = models.AppendArrayToArray(updateTask.TaskRequirements, task.TaskRequirements)

	_, err = r.DBConn.Where(" task_id = ? ", task.TaskId).Update(updateTask)
	if err != nil{
		errorLogger.Println("failed to update task", err)
		session.Rollback()
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to update task",
		})
		return err
	}

	// Commit transaction
	err = session.Commit()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to commit transaction",
		})
	}
	infoLogger.Println("task update successful")
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"msg": task})
}


func (r *Repository) GetAllTasks(ctx *fiber.Ctx) error {
	ctx.Response().Header.Add("Cache-Time", "6000")
	tasks := &[]models.TaskModel{}
	session := r.DBConn.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil{
		errorLogger.Println("DB session closed: ", err)
		return err
	}

	err = r.DBConn.Find(tasks)
	if err != nil{
		errorLogger.Println("failed to delete task", err)
		session.Rollback()
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "failed to delete task",
		})
		return err
	}

	// Commit transaction
	err = session.Commit()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to commit transaction",
		})
	}

	infoLogger.Println("tasks")
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"msg": tasks})
}



func (r *Repository) GetAllUserTasks(ctx *fiber.Ctx) error {
	ctx.Response().Header.Add("Cache-Time", "6000") // cache response
	tokenStr := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenStr)
	if err != nil{
		errorLogger.Println("session error: ", err)
		return err
	}

	user_id := claims.UserId
	if claims.UserType != "USER" {
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"can't access this resource",
		})
		return err
	}

	tasks := &[]models.TaskModel{}
	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("DB session closed: ", err)
		return err
	}

	err = r.DBConn.Where(" user_id = ? ", user_id).Find(tasks)
	if err != nil{
		errorLogger.Println("user tasks not found", err)
		session.Rollback()
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "user tasks not found",
		})
		return err
	}
	// Commit transaction
	err = session.Commit()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"error": "failed to commit transaction",
		})
	}
	infoLogger.Println("user tasks")
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{"data": tasks})
}





func (r *Repository) DeleteTask(ctx *fiber.Ctx) error {
	task_id := ctx.Params("task_id")
	tokenString := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenString)
	if err != nil{
		errorLogger.Println("invalid user id", err)
		return err
	}
	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil{
		errorLogger.Println("DB session closed: ", err)
		return err
	}

	user_id := claims.UserId
	infoLogger.Println(user_id)
	if claims.UserType != "USER"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"can't access this resource",
		})
		return err
	}
	task := &models.TaskModel{}
	_, err = r.DBConn.Where(" user_id = ? AND task_id = ? ", user_id, task_id).Delete(task)
	if err != nil{
		errorLogger.Println("failed to delete task", err)
		session.Rollback()
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "failed to delete task",
		})
		return err
	}
	err = session.Commit()
	if err != nil{
		errorLogger.Println("transaction commit failed", err)
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "DB session failed",
		})
		return err
	}
	infoLogger.Println("task removed successfully")
	return ctx.Status(http.StatusNoContent).JSON(&fiber.Map{"msg": "task removed"})
}


func (r *Repository) CreateApplication(ctx *fiber.Ctx) error {
	task_id := ctx.Params("task_id")
	tokenStr := ctx.Get("Authorization")
	infoLogger.Println(task_id)
	claims, err := models.GetIdFromToken(tokenStr)
	errandRunnerId := claims.UserId
	if claims.UserType != "ERRAND"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"User cannot access this resource",
		})
		return err
	}
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"error":"Incorrect token used"})
	}
    app := &models.ErrandApplication{}
    if err := ctx.BodyParser(app); err != nil {
        return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
    }

    app.Status = "pending" // Set default status
	app.AppId = uuid.New().String()
    app.CreatedAt = time.Now()
    app.UpdatedAt = time.Now()
	app.TaskId = task_id
	app.UserId = errandRunnerId
	app.Email = claims.Email

    _, err = r.DBConn.Insert(app)
    if err != nil {
		errorLogger.Println(err)
        return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create application"})
    }

	infoLogger.Println("application successful")
    return ctx.Status(http.StatusCreated).JSON(app)
}


func (r *Repository) GetApplicationsByErrandID(ctx *fiber.Ctx) error {
	ctx.Response().Header.Add("Cache-Time", "6000")
    task_id := ctx.Params("task_id")
    if task_id == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "task ID is required"})
    }

	tokenStr := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenStr)

	if claims.UserType != "USER"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"User cannot access this resource",
		})
		return err
	}
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"error":"Incorrect token used"})
	}

    var applications []models.ErrandApplication
    err = r.DBConn.Where("task_id = ? AND status = ?", task_id, "pending").Find(&applications)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve applications"})
    }

	infoLogger.Println("errand runner application")
    return ctx.JSON(applications)
}

func (r *Repository) UpdateApplicationStatus(ctx *fiber.Ctx) error {
    appID := ctx.Params("appID")
    if appID == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Application ID is required"})
    }

	tokenStr := ctx.Get("Authorization")
	claims, err := models.GetIdFromToken(tokenStr)

	if claims.UserType != "USER"{
		errorLogger.Println("unauthorized access")
		ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"msg":"User cannot access this resource",
		})
		return err
	}

    var updateData struct {
        Status string `json:"status"`
    }
    if err := ctx.BodyParser(&updateData); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
    }

    if updateData.Status == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status is required"})
    }

	session := r.DBConn.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
        return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}
    app := new(models.ErrandApplication)
    has, err := r.DBConn.Where(" app_id = ? AND status = ? ", appID, "pending").Get(app)
    if err != nil {
		session.Rollback()
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve application"})
    }
    if !has {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Application not found"})
    }

    app.Status = updateData.Status
    app.UpdatedAt = time.Now()

    _, err = r.DBConn.Where(" app_id = ? ",appID).Update(app)
    if err != nil {
		session.Rollback()
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update application"})
    }

	if app.Status == "rejected"{
		infoLogger.Printf("Your application status is: %v", app.Status)
	}

    if err := session.Commit(); err != nil {
        return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
    }

	infoLogger.Println("application status updated")
    return ctx.JSON(app)
}




func (repo *Repository) LoginHandler(ctx *fiber.Ctx) error {
    loginObj := &models.Login{}
    if err := ctx.BodyParser(loginObj); err != nil {
        return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
    }

    if err := ValidateLoginData(loginObj); err != nil {
        return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    session := repo.DBConn.NewSession()
    defer session.Close()

    if err := session.Begin(); err != nil {
        return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
    }

    user, err := repo.AuthenticateUser(loginObj)
    if err != nil {
        session.Rollback()
        return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
    }

    token, err := models.GenerateToken(user.Email, user.UserType, user.UserId)
    if err != nil {
        session.Rollback()
        return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
    }

    if err := session.Commit(); err != nil {
        return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
    }

	infoLogger.Println("login successful")
    return ctx.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}
