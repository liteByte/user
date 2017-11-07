package main

import (
	"crypto/rand"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Usecase interface {
	Create(email string, password string, name string, age uint, number int, date time.Time) (*Model, error)
	Login(email string, password string) (string, *Model, error)
	Find(filter []map[string]string, order map[string]string, offset, limit int) ([]Model, error)
	Update(updates map[string]interface{}, filter []map[string]string) error
	Delete(filter []map[string]string) error
	UpdateOne(model *Model, updates map[string]interface{}) (*Model, error)
	DeleteOne(model *Model) error
}

type UsecaseHandler struct {
	persistenceHandler Persistence
	config             *Config
}

func (h *UsecaseHandler) Create(email string, password string, name string, age uint, number int, date time.Time) (*Model, error) {

	model := Model{
		Email:  email,
		Name:   name,
		Age:    age,
		Number: number,
		Date:   date,
	}

	if len(name) < 5 {
		return nil, Error{Code: http.StatusBadRequest, Message: "Name should be longer than 5 characters"}
	}

	if age < 5 {
		return nil, Error{Code: http.StatusBadRequest, Message: "Age should be greater than 5"}
	}

	model.ProtectionScheme = "lizard.v1"
	protectedForm, err := h.ProtectedFormFromPassword(password)
	if err != nil {
		panic(err)
	}
	model.Password = protectedForm
	model.Compromised = false

	inUse, err := h.isEmailInUse(model.Email)
	if err != nil {
		panic(err)
	}
	if inUse == false {
		return nil, Error{Code: http.StatusConflict, Message: "Email is already in use"}
	}

	//TODO remove later
	if err := h.persistenceHandler.Migrate(&model); err != nil {
		panic(err)
	}

	if err := h.persistenceHandler.Create(&model); err != nil {
		panic(err)
	}

	return &model, nil
}

func (h *UsecaseHandler) isEmailInUse(email string) (bool, error) {
	const other = false
	_, err := h.FindByEmail(email)
	if err != nil {
		switch err.(type) {
		case Error:
			if err.(Error).Code == http.StatusNotFound {
				return false, nil
			} else {
				return other, err
			}
		default:
			return other, err
		}
	} else {
		return true, nil
	}
}

func (h *UsecaseHandler) ProtectedFormFromPassword(password string) (string, error) {

	salt, err := h.generateSalt()
	if err != nil {
		return "", err
	}

	hash, err := h.generateHash([]byte(password), salt)
	if err != nil {
		return "", err
	}

	protectedForm := append(salt, hash...)

	return fmt.Sprintf("%x", protectedForm), nil
}

func (h *UsecaseHandler) generateSalt() ([]byte, error) {
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func (h *UsecaseHandler) generateHash(password []byte, salt []byte) ([]byte, error) {
	//1<<15 = 2ˆ15 sorry for the bit of magic
	hash, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 64)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (h *UsecaseHandler) ComparePasswordAndProtectedForm(password, protectedForm string) (bool, error) {

	salt := protectedForm[:32]
	oldHash := protectedForm[32:]

	newHash, err := h.generateHash([]byte(password), []byte(salt))
	if err != nil {
		return false, err
	}

	if fmt.Sprintf("%x", newHash) != oldHash {
		return false, nil
	}

	return true, nil
}

func (h *UsecaseHandler) Login(email string, password string) (string, *Model, error) {

	user, err := h.FindByEmail(email)
	if err != nil {
		if v, ok := err.(Error); ok && v.Code == http.StatusNotFound {
			return "", nil, Error{Code: http.StatusUnauthorized, Message: "Email or password incorrect"}
		}
		return "", nil, err
	}

	ok, err := h.ComparePasswordAndProtectedForm(password, user.Password)
	if err != nil {
		return "", nil, err
	}
	if ok == false {
		return "", nil, Error{Code: http.StatusUnauthorized, Message: "Email or password incorrect"}
	}

	token, err := h.CreateToken(h.config, user.ID, user.Email)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

type JWTCustomClaims struct {
	Email string
	jwt.StandardClaims
}

type JWTToken struct {
	Id    uint
	Name  string
	Email string
	Admin bool
}

func (h *UsecaseHandler) CreateToken(config *Config, id uint, email string) (string, error) {
	claims := JWTCustomClaims{
		email,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			Issuer:    config.AppName,
			Subject:   strconv.FormatUint(uint64(id), 10),
			Audience:  config.AppName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *UsecaseHandler) Find(filter []map[string]string, order map[string]string, offset, limit int) ([]Model, error) {

	models, err := h.persistenceHandler.Find(filter, order, offset, limit)
	if err != nil {
		panic(err)
	}

	//TODO change these responses to return something usable to other usecases. The transformation to a body type will happen in the endpoints!
	return models, nil
}

func (h *UsecaseHandler) FindByEmail(email string) (*Model, error) {

	models, err := h.Find([]map[string]string{{"email = ?": email}}, nil, 0, 1)
	if err != nil {
		return nil, err
	}

	if len(models) < 1 {
		return nil, Error{Code: http.StatusNotFound, Message: "Not found"}
	}

	return &models[0], nil
}

func (h *UsecaseHandler) Update(updates map[string]interface{}, filter []map[string]string) error {

	if err := h.persistenceHandler.UpdateMany(updates, filter); err != nil {
		panic(err)
	}

	return nil
}

func (h *UsecaseHandler) Delete(filter []map[string]string) error {

	if err := h.persistenceHandler.DeleteMany(filter); err != nil {
		panic(err)
	}

	return nil
}

func (h *UsecaseHandler) UpdateOne(model *Model, updates map[string]interface{}) (*Model, error) {

	if err := h.persistenceHandler.UpdateFields(model, updates); err != nil {
		panic(err)
	}

	return model, nil
}

func (h *UsecaseHandler) DeleteOne(model *Model) error {

	if err := h.persistenceHandler.Delete(model); err != nil {
		panic(err)
	}

	return nil
}
