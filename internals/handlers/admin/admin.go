package admin

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	// "github.com/google/uuid"
	"github.com/iamatila/branch_ass_backend/database"
	model "github.com/iamatila/branch_ass_backend/internals/models/admin"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// CreateAdmin func create a Admin
// @Description Create a Admin
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} model.Admin
// @router /api/v1/admins/signup [post]
func CreateAdmin(c *fiber.Ctx) error {
	db := database.DB
	admin := new(model.Admin)

	// Store the body in the admin and return error if encountered
	err := c.BodyParser(admin)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// Add a uuid to the admin
	// admin.ID = uuid.New()
	rand.Seed(time.Now().UnixNano())

	num := rand.Intn(100) * 100
	num += rand.Intn(100) * 10
	// num += rand.Intn(100) * 1

	var NewAdminID = fmt.Sprintf("%s%d", "BRANCH", num)
	admin.Adminid = NewAdminID

	// Adding a admin type
	admin.Admintype = fmt.Sprintf("%s", "$BranchAdmin")

	// Hash Admin Password
	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	admin.Password = string(hash)

	// Create the Admin and return error if encountered
	err = db.Create(&admin).Error
	if err != nil {
		if err.Error() == "email taken" {
			return c.Status(409).JSON(fiber.Map{
				"status":  "error",
				"message": "Email already taken",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create admin",
			"data":    err,
		})
	}

	// Return the created admin
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Created Admin",
		"data":    admin,
	})

}

// GetAllAdmins func gets all existing admins
// @Description Get all existing admins
// @Tags Admins
// @Accept json
// @Produce json
// @Success 200 {array} model.Admin
// @router /api/v1/admin/all [get]
func GetAllAdmins(c *fiber.Ctx) error {
	db := database.DB
	var admins []model.Admin

	// find all admins in the database
	db.Find(&admins)

	// If no user is present return an error
	if len(admins) == 0 {
		empty := []int{}
		return c.JSON(fiber.Map{"status": "error", "message": "No admin present", "data": empty})
		// return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No admins present", "data": nil})
	}

	// Else return admins
	return c.JSON(fiber.Map{"status": "success", "message": "Admins Found", "data": admins})
}

// GetOneAdmin func one admin by Adminid
// @Description Get one admin by Adminid
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} model.Admin
// @router /api/v1/admin/one/:adminid [get]
func GetOneOfTheAdmins(c *fiber.Ctx) error {
	db := database.DB
	var admin model.Admin

	// Read the param Adminid
	id := c.Params("adminid")

	// Find the admin with the given Adminid
	db.Find(&admin, "adminid = ?", id)

	// If no such admin present return an error
	if admin.Adminid == "" {
		empty := []int{}
		return c.JSON(fiber.Map{"status": "error", "message": "No admin present", "data": empty})
		// return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No admin present", "data": nil})
	}

	// Return the admin with the Adminid
	return c.JSON(fiber.Map{"status": "success", "message": "Admin Found", "data": admin})
}

// LOGIN func login a admin
// @Description Login a admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 200 {object} model.Admin
// @router /api/v1/admins/login [post]
func AdminLogin(c *fiber.Ctx) error {
	req := new(model.Admin)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Input fields can not be empty")
	}

	// admin := new(model.Admin)
	db := database.DB
	var admin model.Admin

	// Find the admin with the given Email
	db.Find(&admin, "email = ?", req.Email)

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, exp, err := createAdminLoginJWTToken(admin)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged In",
		"token":   token,
		"exp":     exp,
	})
}

// Login JWT
func createAdminLoginJWTToken(admin model.Admin) (string, int64, error) {
	err := godotenv.Load("env")
	// err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	secret := os.Getenv("BRANCH_SECRET")
	exp := time.Now().Add(time.Hour * 12).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = admin.Email
	claims["admintype"] = admin.Admintype
	claims["adminid"] = admin.Adminid
	claims["firstname"] = admin.Firstname
	claims["lastname"] = admin.Lastname
	claims["phone"] = admin.Phone
	claims["username"] = admin.Username
	claims["exp"] = exp
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}
