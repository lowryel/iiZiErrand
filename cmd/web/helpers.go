package main

import (
	"errors"
	"math"
	"strconv"
	"time"

	"fmt"

	// rank "github.com/eugene/iizi_errand"
	"github.com/eugene/iizi_errand/pkg/models"
	"golang.org/x/crypto/bcrypt"
)



func (r *Repository) insertUser(user *models.UserModel) error {
	if _, err := r.DBConn.Insert(user); err != nil {
		return err
	}

	profile := createUserProfile(user)
	if _, err := r.DBConn.Insert(profile); err != nil {
		return err
	}

	return nil
}


func createUserProfile(user *models.UserModel) interface{} {
	// location, err := rank.GetLocation(fmt.Sprintf(url, api_key))
    // if err != nil {
    //     return err
    // }
	if user.UserType == "USER" {
		return &models.UserProfile{
			UserId:    user.UserId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			UserType:  user.UserType,
			// Location: location,
		}
	}
	return &models.ErrandRunnerProfile{
		UserId:    user.UserId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		UserType:  user.UserType,
		// Location: location,
	}
}


func ValidateLoginData(loginObj *models.Login) error {
    if loginObj.Email == "" || loginObj.Password == "" {
        return errors.New("email and password are required")
    }
    // Add more validations as needed
    return nil
}

func (repo *Repository) AuthenticateUser(loginObj *models.Login) (*models.UserModel, error) {
    user := &models.LoginData{}
    _, err := repo.DBConn.SQL("SELECT * FROM user_model WHERE email = ?", loginObj.Email).Get(user)
    if err != nil {
        return nil, errors.New("error fetching login data")
    }

    if err = models.CompareHashAndPass(user.Password, loginObj.Password); err != nil {
        return nil, errors.New("invalid credentials")
    }

    fullUser := &models.UserModel{}
    _, err = repo.DBConn.Where("email = ? AND user_type = ?", user.Email, user.UserType).Get(fullUser)
    if err != nil {
        return nil, errors.New("error retrieving user data")
    }

    return fullUser, nil
}


func (r *Repository) GetUserProfile(userId string) (*models.UserProfile, error) {
    var userProfile models.UserProfile
    _, err := r.DBConn.Where("user_id = ?", userId).Get(&userProfile)
    if err != nil {
		return nil, fmt.Errorf("failed to retrieve user profile: %v", err)
    }
    return &userProfile, nil
}





func (repo *Repository) UpdateUserPassword(email string, newPassword string) error {
    // Hash the new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash new password: %v", err)
    }

    // Update the user's password in the database
    _, err = repo.DBConn.Exec("UPDATE user_model SET password = ? WHERE email = ?", string(hashedPassword), email)
    if err != nil {
        return fmt.Errorf("failed to update password in database: %v", err)
    }

    return nil
}




func CalculateDistanceScore(taskLoc, runnerLoc models.Location) float64 {
    distance := calculateDistance(taskLoc, runnerLoc)
    return math.Max(1-distance/10, 0) // Assume 10km is the max preferred distance
}


func calculateDistance(loc1, loc2 models.Location) float64 {
    // Haversine formula for calculating distance between two points on a sphere
    const earthRadius = 6371 // km

	lat1, err := strconv.ParseFloat(loc1.Latitude, 64)
    if err != nil {
        return 0.0
    }
    lat1 = lat1 * math.Pi / 180
    lon1, err := strconv.ParseFloat(loc1.Longitude, 64)
    if err != nil {
		return 0.0
    }
	lon1 = lon1 * math.Pi / 180

	lat2, err := strconv.ParseFloat(loc2.Latitude, 64)
    if err != nil {
		return 0.0
    }
    lat2 = lat2 * math.Pi / 180
	lon2, err := strconv.ParseFloat(loc2.Longitude, 64)
    if err != nil {
		return 0.0
    }
	lon2 = lon2 * math.Pi / 180

    dlat := lat2 - lat1
    dlon := lon2 - lon1

    a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    return earthRadius * c
}

