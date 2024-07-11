package main

import (
	"errors"
	"time"

	"fmt"
	"net/smtp"

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
	if user.UserType == "USER" {
		return &models.UserProfile{
			UserId:    user.UserId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			UserType:  user.UserType,
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



type EmailService struct {
    smtpServer string
    smtpPort   int
    username   string
    password   string
}

func NewEmailService(server string, port int, username, password string) *EmailService {
    return &EmailService{
        smtpServer: server,
        smtpPort:   port,
        username:   username,
        password:   password,
    }
}

func (s *EmailService) SendEmail(to, subject, body string) error {
    auth := smtp.PlainAuth("", s.username, s.password, s.smtpServer)

    msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

    err := smtp.SendMail(
        fmt.Sprintf("%s:%d", s.smtpServer, s.smtpPort),
        auth,
        s.username,
        []string{to},
        []byte(msg),
    )

    if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

    return nil
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

