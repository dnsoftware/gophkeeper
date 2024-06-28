package constants

import (
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	LogFile       string = "log.log"
	LogLevel             = zapcore.InfoLevel
	LogReadline   string = "readline.log"
	PW_SALT_BYTES        = 16            // длина "соли" для пароля
	PW_HASH_BYTES        = 24            // длина хеша пароля
	JWTTokenExp          = time.Hour * 1 // время истечения токена авторизации
	JWTSecretKey         = "jwtstrong"   // секретный ключ для подписи токена авторизации
	TokenKey      string = "token"       // ключ JWT токена к передаваемах метаданных (контексте)
	FileBankDir   string = "filebank"    // папка в которой хранятся файлы пользователей на сервере
	ChunkSize     int    = 10240         // chunk size для потоковой передачи бинарных данных
	FileStorage   string = "filestorage" // папка куда скачиваются файлы пользователя на клиента
	UserUD        string = "userID"      // идентификатор кода пользователя в GRPC контексте сервера
)

// типы сущностей
const (
	LogopasEntity string = "logopas" // логин-пароль
	CardEntity    string = "card"    // банковская карта
	TextEntity    string = "text"    // произвольные текстовые данные
	BinaryEntity  string = "binary"  // произвольные бинарные данные
)

// типы полей свойств сущности
const (
	FieldTypeString string = "string" // строка
	FieldTypePath   string = "path"   // путь к файлу
)

// Названия методов для которых применяется симметричное шифрования
// шифровка отправляемых данных
const (
	MethodAddEntity string = "AddEntity"
)

// Названия методов для которых применяется симметричное шифрования
// расшифровка полученных данных
const (
	MethodEntity string = "Entity"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

const (
	DBContextTimeout time.Duration = time.Duration(5) * time.Second // длительность запроса в контексте работы с БД
)

// сообщения об ошибках
const (
	ErrPasswordsNotMatch string = "пароли не совпадают"
	ErrBadPassword       string = "bad password"
	ErrNoSuchUser        string = "no such user"
)

// Методы для которых не проверяем токен авторизации
const (
	ExcludeMethodPing         string = "/proto.Keeper/Ping"
	ExcludeMethodRegistration string = "/proto.Keeper/Registration"
	ExcludeMethodLogin        string = "/proto.Keeper/Login"
)
