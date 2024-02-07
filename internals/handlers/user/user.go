package userHandler

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	// "github.com/google/uuid"
	"github.com/iamatila/branch_ass_backend/database"
	model "github.com/iamatila/branch_ass_backend/internals/models/user"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser func create a user
// @Description Create a User
// @Tags User
// @Accept json
// @Produce json
// @Param firstname body string true "FirstName"
// @Param lastname body string true "LastName"
// @Param email body string true "Email"
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Param phone body string true "Phone"
// @Success 200 {object} model.User
// @router /api/v1/users/signup [post]
func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	user := new(model.User)

	// Store the body in the user and return error if encountered
	err := c.BodyParser(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input"})
	}

	// Add a uuid to the user
	// user.ID = uuid.New()
	rand.Seed(time.Now().UnixNano())

	num := rand.Intn(100) * 100
	num += rand.Intn(100) * 10
	// num += rand.Intn(100) * 1
	// var id = ""

	var NewUserID = fmt.Sprintf("%s%d", "BRANCH", num)
	user.Userid = NewUserID
	// id = NewUserID
	// Adding a user type
	user.Usertype = fmt.Sprintf("%s", "$BranchUser")

	// Hash User Password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	user.Password = string(hash)

	// Create the User and return error if encountered
	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create user",
			"data":    err,
		})
	}

	// Return the created user
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Account created successfully",
		"data":    user,
	})

}

// GetAllUsers func gets all existing users
// @Description Get all existing users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} model.User
// @router /api/v1/user/all [get]
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []model.User

	// find all users in the database
	db.Find(&users)

	// If no user is present return an error
	if len(users) == 0 {
		empty := []int{}
		return c.JSON(fiber.Map{"status": "error", "message": "No users present", "data": empty})
		// return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No users present", "data": nil})
	}

	// Else return users
	return c.JSON(fiber.Map{"status": "success", "message": "Users Found", "data": users})
}

// GetOneUser func one user by Userid
// @Description Get one user by Userid
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @router /api/v1/user/one/:userid [get]
func GetOneUser(c *fiber.Ctx) error {
	db := database.DB
	var user model.User

	// Read the param Userid
	id := c.Params("userid")

	// Find the user with the given Userid
	db.Find(&user, "userid = ?", id)

	// If no such user present return an error
	if user.Userid == "" {
		empty := []int{}
		return c.JSON(fiber.Map{"status": "error", "message": "No user present", "data": empty})
		// return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user present", "data": nil})
	}

	// Return the user with the Userid
	return c.JSON(fiber.Map{"status": "success", "message": "User Found", "data": user})
}

// LOGIN func login a user
// @Description Login a User
// @Tags User
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 200 {object} model.User
// @router /api/v1/users/login [post]
func Login(c *fiber.Ctx) error {
	req := new(model.User)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Input fields can not be empty")
	}

	// user := new(model.User)
	db := database.DB
	var user model.User

	// Find the user with the given Email
	db.Find(&user, "email = ?", req.Email)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// return err
		// return c.Status(fiber.StatusBadRequest).SendString("Password is incorrect")
		return c.JSON(fiber.Map{
			"status":  "fail",
			"message": "Password is incorrect",
		})
	}

	token, exp, err := createdLoginJWTToken(user)
	if err != nil {
		// return err
		return c.JSON(fiber.Map{
			"status":  "fail",
			"message": "Unable to generate token",
		})
	}

	// userDetails, userExp, err := createdUserJWTToken(user)
	// if err != nil {
	// 	// return err
	// 	return c.JSON(fiber.Map{
	// 		"status":  "fail",
	// 		"message": "Unable to generate user token",
	// 	})
	// }

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged In",
		"token":   token,
		"exp":     exp,
		// "user_details": userDetails,
		// "user_exp":     userExp,
		// 	"firstname":   user.FirstName,
		// 	"lastname":    user.LastName,
		// 	"email":        user.Email,
		// 	"username":     user.Username,
		// 	"usertype":    user.UserType,
		// 	"phone": user.PhoneNumber,
		// 	"gender":       user.Gender,
		// 	"dob":          user.DOB,
		// 	"address":      user.Address,
		// 	"city":        user.City,
		// 	"country":      user.Country,
	})
}

// Login JWT
func createdLoginJWTToken(user model.User) (string, int64, error) {
	err := godotenv.Load("env")
	// err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	secret := os.Getenv("BRANCH_SECRET")
	exp := time.Now().Add(time.Hour * 12).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["usertype"] = user.Usertype
	claims["userid"] = user.Userid
	claims["firstname"] = user.Firstname
	claims["lastname"] = user.Lastname
	claims["phone"] = user.Phone
	claims["username"] = user.Username
	claims["exp"] = exp
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}

// // User info JWT
// // func createdUserJWTToken(user string) (string, int64, error) {
// func createdUserJWTToken(user model.User) (string, int64, error) {
// 	err := godotenv.Load("env")
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	secret := os.Getenv("BRANCH_SECRET_TOO")
// 	exp := time.Now().Add(time.Minute * 60).Unix()

// 	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(jus))
// 	// tokenString, err := token.SignedString([]byte(secret))
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	token := jwt.New(jwt.SigningMethodHS256)
// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["firstname"] = user.Firstname
// 	claims["lastname"] = user.Lastname
// 	claims["email"] = user.Email
// 	claims["username"] = user.Username
// 	claims["usertype"] = user.Usertype
// 	claims["userid"] = user.Userid
// 	claims["phone"] = user.Phone
// 	claims["city"] = user.City
// 	claims["country"] = user.Country
// 	claims["exp"] = exp
// 	t, err := token.SignedString([]byte(secret))
// 	if err != nil {
// 		return "", 0, err
// 	}

// 	return t, exp, nil
// 	// return tokenString, exp, nil
// }
