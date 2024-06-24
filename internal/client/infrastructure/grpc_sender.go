package infrastructure

import (
	"fmt"
)

type GRPCSender struct {
}

func NewGRPCSender() (*GRPCSender, error) {

	sender := &GRPCSender{}

	return sender, nil
}

// Registration регистрация пользователя
// На входе: логин, пароль, повторный пароль
// Возвращает токен авторизации в случае успеха и ошибку
func (s *GRPCSender) Registration(login string, password string, password2 string) (string, error) {

	if password != password2 {
		return "", fmt.Errorf("пароли не совпадают")
	}

	//resp := pb.RegisterRequest{
	//	Login:          login,
	//	Password:       password,
	//	RepeatPassword: password2,
	//}

	return "", nil
}
