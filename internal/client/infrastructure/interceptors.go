// Перехватчики исходящих и входящих grpc запросов (добавление токена авторизации, шифровка/дешифровка данных)
package infrastructure

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/dnsoftware/gophkeeper/internal/constants"
	"github.com/dnsoftware/gophkeeper/internal/proto"
	"github.com/dnsoftware/gophkeeper/internal/utils"
)

//// AuthDataSet интерфейс работы с данными авторизации (установка текущего токена авторизации и пароля)
//type AuthDataSet interface {
//	SetToken(token string)
//}

/********************************* Добавление токена в запрос ***********************************/

// ActualTokenGet интерфейс работы с данными авторизации (получение текущего токена авторизации и пароля)
type ActualTokenGet interface {
	GetToken() string
}

type AuthInterceptor struct {
	actualToken    ActualTokenGet
	excludeMethods map[string]bool // методы для которых не применяется перезватчик
}

// NewAuthInterceptor перехватчик, добавляющий токен авторизации в контекст
// excludeMethods - карта методов для которых он не применяется
func NewAuthInterceptor(actualToken ActualTokenGet, excludeMethods map[string]bool) *AuthInterceptor {
	a := &AuthInterceptor{
		actualToken:    actualToken,
		excludeMethods: excludeMethods,
	}

	return a
}

// TokenInterceptor добавление токена авторизации в исходящий запрос
func (i *AuthInterceptor) TokenInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) error {

		token := i.actualToken.GetToken()
		if ok := i.excludeMethods[method]; !ok {
			ctx = metadata.AppendToOutgoingContext(ctx, constants.TokenKey, token)
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

/******************************** Шифровка исходящих данных *******************************/

// UserPasswordGet интерфейс получения пароля пользователя, чтобы использовать его как часть ключа шифрования
type UserPasswordGet interface {
	GetPassword() string
}

// DataOutInterceptor перехватчик исходящих данных
type DataOutInterceptor struct {
	userPassword UserPasswordGet // для получения пароль пользователя
	validMethods map[string]bool // методы, нуждающиеся в перехвате
	secretKey    string          // секретный ключ, для шифрования
}

// NewDataOutInterceptor шифрование значимых полей в исходящих запросах
// userPassword - интерфейс получения пароля пользователя, который он вводил при входе
// secretKey - секретный ключ, хранящийся на стороне клиента
// methods - методы к которым применяется перехватчик
func NewDataOutInterceptor(userPassword UserPasswordGet, secretKey string, methods map[string]bool) *DataOutInterceptor {
	a := &DataOutInterceptor{
		userPassword: userPassword,
		validMethods: methods,
		secretKey:    secretKey,
	}

	return a
}

// DataOutputInterceptor перехватчик исходящих данных (шифровка/дешифровка)
func (d *DataOutInterceptor) DataOutputInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) error {

		parts := strings.Split(method, "/")
		methodName := parts[len(parts)-1]

		if d.validMethods[methodName] {
			// логика шифрования
			cryptoKey := utils.SymmPassCreate(d.userPassword.GetPassword(), d.secretKey)
			switch methodName {
			case constants.MethodAddEntity:
				entity := req.(*proto.AddEntityRequest)
				for key, prop := range entity.Props {
					if entity.Etype == constants.BinaryEntity || entity.Etype == constants.TextEntity { // название загружаемого файла не шифруем
						entity.Props[key].Value = prop.Value
					} else {
						entity.Props[key].Value = utils.Encrypt(prop.Value, cryptoKey)
					}
				}
				for key, meta := range entity.Metainfo {
					entity.Metainfo[key].Title = utils.Encrypt(meta.Title, cryptoKey)
					entity.Metainfo[key].Value = utils.Encrypt(meta.Value, cryptoKey)
				}
			case constants.MethodSaveEditEntity:
				entity := req.(*proto.SaveEntityRequest)
				for key, prop := range entity.Props {
					if entity.Etype == constants.BinaryEntity || entity.Etype == constants.TextEntity { // название загружаемого файла не шифруем
						entity.Props[key].Value = prop.Value
					} else {
						entity.Props[key].Value = utils.Encrypt(prop.Value, cryptoKey)
					}
				}
				for key, meta := range entity.Metainfo {
					entity.Metainfo[key].Title = utils.Encrypt(meta.Title, cryptoKey)
					entity.Metainfo[key].Value = utils.Encrypt(meta.Value, cryptoKey)
				}

			}
		}

		f := invoker(ctx, method, req, reply, cc, opts...)

		if d.validMethods[methodName] {
			// логика расшифровки
			cryptoKey := utils.SymmPassCreate(d.userPassword.GetPassword(), d.secretKey)
			switch methodName {
			case constants.MethodEntity:
				entity := reply.(*proto.EntityResponse)
				for key, prop := range entity.Props {
					if entity.Etype == constants.BinaryEntity || entity.Etype == constants.TextEntity { // название файла не нуждается в расшифровке
						entity.Props[key].Value = prop.Value
					} else {
						entity.Props[key].Value = utils.Decrypt(prop.Value, cryptoKey)
					}
				}
				for key, meta := range entity.Metainfo {
					entity.Metainfo[key].Title = utils.Decrypt(meta.Title, cryptoKey)
					entity.Metainfo[key].Value = utils.Decrypt(meta.Value, cryptoKey)
				}

			}

		}

		return f
	}
}
