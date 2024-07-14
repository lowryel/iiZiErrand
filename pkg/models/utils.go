package models

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"

	// "strings"
	"time"

	"github.com/gofiber/fiber"
	"github.com/golang-jwt/jwt"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
)

var (
	infoLogger = log.New(os.Stdout, "[iizi]	[INFO]: \t", log.Ldate | log.Ltime | log.Lshortfile)
	errorLogger = log.New(os.Stdout, "[iizi] [ERROR]: \t", log.Ldate | log.Ltime | log.Lshortfile)
)

var (
	secretKey = os.Getenv("SECRET_FLAVOUR")
)



// JWTClaims struct represents the claims for JWT
type JWTClaims struct {
  UserId      string      `json:"user_id"`
  Email       string    `json:"email"`
  UserType    string 		`json:"user_type"`
  jwt.StandardClaims
}


func ValidateEmail(email string) bool {
	// verify email format
    _, err := mail.ParseAddress(email)
	if err != nil{
		errorLogger.Fatalf("%s", err)
	}
    return err == nil
}

func ValidateUser(user *UserModel) error {
	if !ValidateEmail(user.Email) {
		return errors.New("incorrect email format")
	}

	if user.UserType != "USER" && user.UserType != "ERRAND" {
		return errors.New("wrong user type")
	}

	if user.UserType == "" {
		return errors.New("specify user type")
	}

	return nil
}


func ValidateChangePassData(user *ChangePass) error {
	if !ValidateEmail(user.Email) {
		return errors.New("incorrect email format")
	}
	if user.OldPass == "" && user.NewPass == "" {
		return errors.New("wrong inputs")
	}

	if user.OldPass == user.NewPass {
		return errors.New("passwords can't be similar")
	}
	return nil
}



func HashPass(pass string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err!= nil {
		errorLogger.Println("password hashing failed")
	}
	infoLogger.Println(string(hash))
	return string(hash), nil
}

func CompareHashAndPass(hashed, pass string) error{
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	return err
}


// generate jwt token
func GenerateToken(email, user_type string, userId string) (string, error) {
    claims := JWTClaims{
		UserId: userId,
		Email: email,
		UserType: user_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
    if err != nil {
        return "", err
    }
    return tokenString, nil
}


// DecodeToken decodes and validates the JWT token and extracts the claims.
func DecodeToken(tokenString string) (*JWTClaims, error) {
    // Parse and validate the JWT token
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })
    if err != nil || !token.Valid {
        return nil, fmt.Errorf("invalid or expired token: %v", err)
    }

    // Extract the claims from the token
    claims, ok := token.Claims.(*JWTClaims)
    if !ok {
        return nil, fmt.Errorf("failed to extract claims from token")
    }

    return claims, err
}


func GetIdFromToken(tokenString string) (*JWTClaims, error){
	if tokenString == ""{
		errorLogger.Println("token not provided")
		return nil, fmt.Errorf("token missing")
	}
	tokenStr := strings.Split(tokenString, " ")[1]
	claims, err := DecodeToken(tokenStr)
	if err != nil{
		log.Println(err)
	}
	return claims, err
}



// JWTMiddleware checks the JWT token in the request headers
// JWTMiddleware caches the parsed and validated token 
// for performance
func JWTMiddleware() fiber.Handler {
	tokenCache := cache.New(5*time.Minute, 10*time.Minute)
	return func(c *fiber.Ctx) {
		tokenString := c.Get("Authorization")
		if tokenString == ""{
			log.Println("token required")
			c.Status(http.StatusForbidden).JSON(&fiber.Map{
				"msg": "token unavailable",
			})
		}
		tokenString = tokenString[7:]
		//  check cache
		if token, ok := tokenCache.Get(tokenString); ok {
			c.Locals("user", token)
			c.Status(http.StatusForbidden).JSON(&fiber.Map{
				"msg": token,
			})
		}
		token, err := parseAndValidateToken(tokenString)
		if err != nil {
			log.Printf("Token invalid error: %s\n", err)
			c.Status(http.StatusForbidden).JSON(&fiber.Map{
				"msg": err,
			})
		}
	// Cache valid token
	tokenCache.Set(string(tokenString), token, 0)
		// Extract and store user identity
		user := token.Claims.(jwt.MapClaims) 
		c.Locals("user", user)
		c.Next()
	}
}


// Centralized parsing and validation
func parseAndValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["Token"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err 
	}

	// Additional validation checks
	if !token.Valid {
		return nil, ErrorToken
	}
	return token, nil
}
// Custom error types
var (
  	ErrorToken= errors.New("token expired")
)

/* errandProfile.Skills [This implementation will: Keep all existing skills. Add any new skills from the update request that aren't already in the list. Avoid duplicates.] */

// append more skills to the array
func AppendArrayToArray(existingSkills, newSkills []string) ([]string) {
// func appendUniqueSkills(existingSkills, newSkills []string) []string {
    skillMap := make(map[string]bool)
    for _, skill := range existingSkills {
        skillMap[skill] = true
    }
    
    for _, skill := range newSkills {
        if !skillMap[skill] {
            existingSkills = append(existingSkills, skill)
            skillMap[skill] = true
        }
    }
    
    return existingSkills
}



