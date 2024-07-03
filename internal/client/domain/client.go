package domain

import (
	"fmt"
	"io"
	"log"
)

// Работа с командной строкой
type Readline interface {
	input(prompt string, validateRules string, validateMessages string) (string, error)
	edit(prompt string, what string, validateRules string, validateMessages string) (string, error)
	Writeln(str string)
	Registration() (string, string, error)
	Login() (string, string, error)
	Stderr() io.Writer
	Close() error
	MakeFieldsDescription(fields []*Field)
	GetEtypeName(etype string) string
	SetEtypeName(etype string, name string)
	GetField(fieldID int32) *Field
	GetFieldsGroup(etype string) []*Field
	interrupt(line string, err error) string
}

// Sender Интерфейс отправки/приема данных с сервера
type Sender interface {
	Registration(login string, password string, password2 string) (string, error)
	Login(login string, password string) (string, error)
	EntityCodes() ([]*EntityCode, error)
	Fields(etype string) ([]*Field, error)
	AddEntity(ae Entity) (int32, error)
	SaveEntity(ae Entity) (int32, error)
	DeleteEntity(id int32) error
	UploadBinary(entityId int32, file string) (int32, error)
	DownloadBinary(entityId int32, fileName string) (string, error)
	UploadCryptoBinary(entityId int32, file string) (int32, error)
	DownloadCryptoBinary(entityId int32, fileName string) (string, error)
	EntityList(etype string) (map[int32]string, error)
	Entity(id int32) (*Entity, error)
}

type Entity struct {
	Id       int32       // ID сущности
	UserID   int32       // ID пользователя
	Etype    string      // тип сущности: card, text, logopas, binary и т.д.
	Props    []*Property // массив значений свойств
	Metainfo []*Metainfo // массив значений метаинформации
}

type Property struct {
	EntityId int32  // код сущности
	FieldId  int32  // код описания поля свйоства
	Value    string // значение свойства
}

type Metainfo struct {
	EntityId int32  // код сущности
	Title    string // наименование метаинформации
	Value    string // значение метаинформации
}

type EntityCode struct {
	Etype string
	Name  string
}

type Field struct {
	Id               int32
	Name             string
	Etype            string
	Ftype            string
	ValidateRules    string
	ValidateMessages string
}

type GophKeepClient struct {
	//rl     *CLIReader
	rl     Readline
	Sender Sender
}

// BinaryFileProperty Данные в поле свойства бинарной сущности содержат JSON в формате:
// {"servername": "имя файла на сервере (полный путь), "clientname": "только имя файла, под которым его грузили с клиента", "chunkcount": "кол-во фрагментов на которые разбит файл"}
type BinaryFileProperty struct {
	Servername string `json:"servername"`
	Clientname string `json:"clientname"`
	Chunkcount int32  `json:"chunkcount"`
}

const (
	WorkAgain string = "again"
	WorkStop  string = "stop"
)

func NewGophKeepClient(readline Readline, sender Sender) (*GophKeepClient, error) {

	client := &GophKeepClient{
		rl:     readline,
		Sender: sender,
	}

	return client, nil

}

func (c *GophKeepClient) Start() {
	defer c.rl.Close()

	log.SetOutput(c.rl.Stderr())

	var token string // токен авторизации

	// Логин или регистрация
	for {
		fmt.Println("Нажмите [Enter] для входа или \"r\" для регистрации  ")
		line, err := c.rl.input(">>", "", "{}")
		if err != nil {
			c.rl.Writeln(err.Error())
			return
		}

		// если не регистрация - переходим к вводу логина и пароля для входа
		if line != "r" {
			break
		}

		// Регистрация
		login, password, err := c.rl.Registration()
		if err != nil {
			return
		}

		token, err = c.Sender.Registration(login, password, password)
		if err != nil {
			c.rl.Writeln("Sender.Registration: " + err.Error())
			continue
		}

		break
	}

	// Если уже ранее регистрировались - запрашиваем логин-пароль
	// без аутентификации дальнейшая работа невозможна
	if token == "" {
		for {
			login, password, err := c.rl.Login()
			if err != nil {
				return
			}

			token, err = c.Sender.Login(login, password)
			if err != nil {
				c.rl.Writeln(err.Error())
				continue
			}

			break
		}
	}

	// Инициализация списка сущностей, с которыми можно работать
	entCodes, err := c.Sender.EntityCodes()

	if err != nil {
		c.rl.Writeln(fmt.Sprintf("Ошибка загрузки сущностей: %v", err))
	}
	for _, val := range entCodes {
		c.rl.SetEtypeName(val.Etype, val.Name)
	}

	// Инициализация описаний полей сущностей
	for _, val := range entCodes {
		fields, err := c.Sender.Fields(val.Etype)
		if err != nil {
			c.rl.Writeln(fmt.Sprintf("Ошибка загрузки полей с описаниями: %v", err))
		}
		c.rl.MakeFieldsDescription(fields)
	}

	/************** Основная логика ************/

	for {
		status, err := c.Base(entCodes)
		if err != nil {
			fmt.Println(err.Error())
		}

		if status == WorkAgain {
			continue
		}

		if status == WorkStop {
			break
		}
	}

	//exit:
}

// DisplayEntity отобразить сущность в консоли
func (c *GophKeepClient) DisplayEntity(ent Entity) {
	fmt.Println("------------------------")
	fmt.Println(" " + c.rl.GetEtypeName(ent.Etype))
	for _, val := range ent.Props {
		fmt.Println("      " + c.rl.GetField(val.FieldId).Name + ": " + val.Value)
	}
	for _, val := range ent.Metainfo {
		fmt.Println("      " + val.Title + ": " + val.Value)
	}
	fmt.Println("------------------------")
}

// DisplayEntityBinary отобразить сущность в консоли и показать путь к загруженному файлу
func (c *GophKeepClient) DisplayEntityBinary(ent Entity, filePath string) {
	fmt.Println("------------------------")
	fmt.Println(" " + c.rl.GetEtypeName(ent.Etype))
	fmt.Println("Путь к загруженному файлу: " + filePath)
	for _, val := range ent.Metainfo {
		fmt.Println("      " + val.Title + ": " + val.Value)
	}
	fmt.Println("------------------------")
}
