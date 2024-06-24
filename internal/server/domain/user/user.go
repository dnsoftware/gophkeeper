package user

import (
	"context"
	"errors"
	"time"

	"github.com/dnsoftware/gophkeeper/internal/constants"
	"github.com/dnsoftware/gophkeeper/internal/utils"
)

type UserStorage interface {
	// UserCreate регистрация нового пользователя
	UserCreate(ctx context.Context, login string, password string, salt string) (int, error)

	// GetUser получение данных пользователя по логину (ID, дата добавления)
	GetUser(ctx context.Context, login string) (int, time.Time, error)

	// LoginUser получение ID зарегистрированного пользователя или 0, если нет в базе
	LoginUser(ctx context.Context, login string, password string) (int, string)
}

type User struct {
	storage UserStorage
}

func NewUser(storage UserStorage) (*User, error) {
	user := &User{
		storage: storage,
	}

	return user, nil
}

// Registration регистрация нового пользователя. Возвращает токен доступа в случае удачи и ошибку, если что-то пошло не так.
func (k *User) Registration(ctx context.Context, login string, password string, repeatPassword string) (string, error) {

	// проверка на совпадение паролей
	if password != repeatPassword {
		return "", errors.New(constants.ErrPasswordsNotMatch)
	}

	// проверка на наличие логина в базе (id = 0, если логина нет в базе)
	id, _, err := k.storage.GetUser(ctx, login)
	if err != nil {
		return "", err
	}

	// если такой логин уже есть
	if id > 0 {
		return "", errors.New("такой логин уже занят")
	}

	// генерируем пароль и заносим в базу
	_, saltStr := utils.SaltGenerate()
	passHash := utils.PassGenerate(password, saltStr)
	userId, err := k.storage.UserCreate(ctx, login, passHash, saltStr)
	jwt, err := utils.BuildJWTString(userId)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

// Login Вход пользователя. При успехе возвращает токен доступа, при неудаче - пустую строку и ошибку.
func (k *User) Login(ctx context.Context, login string, password string) (string, error) {
	userID, errorMessage := k.storage.LoginUser(ctx, login, password)
	if userID == 0 {
		return "", errors.New(errorMessage)
	}

	token, err := utils.BuildJWTString(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}
