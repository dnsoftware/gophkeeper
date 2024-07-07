// Работа в консоли, обмен данными с сервером
package domain

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dnsoftware/gophkeeper/internal/constants"
)

// Readline Работа с командной строкой
type Readline interface {
	// input управляет вводом строки в консоли
	input(prompt string, validateRules string, validateMessages string) (string, error)
	// edit редактирование строки данных в консоли
	edit(prompt string, what string, validateRules string, validateMessages string) (string, error)
	// Writeln вывод в консоль с переводом строки
	Writeln(str string)
	// Registration регистрация
	Registration() (string, string, error)
	// Login логин
	Login() (string, string, error)
	Stderr() io.Writer
	// Close завершение работы в консоли
	Close() error
	// MakeFieldsDescription Формирование карт описаний полей сущностей
	MakeFieldsDescription(fields []*Field)
	// GetEtypeName получение названия типа сущности по коду
	GetEtypeName(etype string) string
	// SetEtypeName установка названия типа сущности
	SetEtypeName(etype string, name string)
	// GetField получение описания поля сущности по ID поля
	GetField(fieldID int32) *Field
	// GetFieldsGroup получение группы полей сущности по коду типа сущности
	GetFieldsGroup(etype string) []*Field
	// interrupt прерывание ввода
	interrupt(line string, err error) string
}

// Sender Интерфейс отправки/приема данных с сервера
type Sender interface {
	// Registration регистрация
	Registration(login string, password string, password2 string) (string, error)
	// Login логин
	Login(login string, password string) (string, error)
	// EntityCodes получение справочника типов сущностей
	EntityCodes() ([]*EntityCode, error)
	// Fields получение описаний полей сущностей
	Fields(etype string) ([]*Field, error)
	// AddEntity добавление сущности
	AddEntity(ae Entity) (int32, error)
	// SaveEntity сохранение сущности
	SaveEntity(ae Entity) (int32, error)
	// DeleteEntity удаление сущности
	DeleteEntity(id int32) error
	// UploadBinary загрузка незашифрованных бинарных данных (клиент -> сервер)
	UploadBinary(entityId int32, file string) (int32, error)
	// DownloadBinary отдача незашифрованных бинарных данных клиенту (сервер -> клиент)
	DownloadBinary(entityId int32, fileName string) (string, error)
	// UploadCryptoBinary получение зашифрованных бинарных данных с клиента (клиент -> сервер)
	UploadCryptoBinary(entityId int32, file string) (int32, error)
	// DownloadCryptoBinary отдача зашифрованных бинарных данных клиенту (сервер -> клиент)
	DownloadCryptoBinary(entityId int32, fileName string) (string, error)
	// EntityList Получение списка сущностей указанного типа для конкретного пользователя
	// Простая карта с кодом сущности и названием(составляется из метаданных)
	EntityList(etype string) (map[int32]string, error)
	// Entity получение сущности
	Entity(id int32) (*Entity, error)
}

// Entity сущность
type Entity struct {
	Id       int32       // ID сущности
	UserID   int32       // ID пользователя
	Etype    string      // тип сущности: card, text, logopas, binary и т.д.
	Props    []*Property // массив значений свойств
	Metainfo []*Metainfo // массив значений метаинформации
}

// Property свойство сущности
type Property struct {
	EntityId int32  // код сущности
	FieldId  int32  // код описания поля свйоства
	Value    string // значение свойства
}

// Metainfo метаинформация сущности
type Metainfo struct {
	EntityId int32  // код сущности
	Title    string // наименование метаинформации
	Value    string // значение метаинформации
}

// EntityCode название типа сущности
type EntityCode struct {
	Etype string // код типа сущности
	Name  string // название типа сущности
}

// Field описание поля сущности
type Field struct {
	Id               int32  // ID поля
	Name             string // название поля
	Etype            string // код типа сущности
	Ftype            string // тип поля
	ValidateRules    string // правила валидации при вводе поля
	ValidateMessages string // сообщения валидации при ее непрохождении
}

// GophKeepClient клиент, управляет вводом данных в консоли и отправкой/получением данных с/на сервер
type GophKeepClient struct {
	rl     Readline // работа в консоли
	Sender Sender   // отправка-получение данных на/с сервера
}

// BinaryFileProperty Данные в поле свойства бинарной сущности содержат JSON в формате:
// {"servername": "имя файла на сервере (полный путь), "clientname": "только имя файла, под которым его грузили с клиента", "chunkcount": "кол-во фрагментов на которые разбит файл"}
type BinaryFileProperty struct {
	Servername string `json:"servername"`
	Clientname string `json:"clientname"`
	Chunkcount int32  `json:"chunkcount"`
}

const (
	WorkAgain string = "again" // повторение цикла ввода в консоли сначала
	WorkStop  string = "stop"  // завершение цикла ввода в консоли
)

// NewGophKeepClient конструктор
func NewGophKeepClient(readline Readline, sender Sender) (*GophKeepClient, error) {

	client := &GophKeepClient{
		rl:     readline,
		Sender: sender,
	}

	return client, nil

}

// Start старт консольного клиента
func (c *GophKeepClient) Start(stopChan chan bool) error {

	var token string // токен авторизации

	// Логин или регистрация
	for {
		line, err := c.rl.input(`Нажмите [Enter] для входа или "r" для регистрации>>`, "", "{}")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		// если не регистрация - переходим к вводу логина и пароля для входа
		if line != "r" {
			break
		}

		// Регистрация
		login, password, err := c.rl.Registration()
		if err != nil {
			return err
		}

		token, err = c.Sender.Registration(login, password, password)
		if err != nil {
			fmt.Println("Sender.Registration: " + err.Error())
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
				return err
			}

			token, err = c.Sender.Login(login, password)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			break
		}
	}

	// Инициализация списка сущностей, с которыми можно работать
	entCodes, err := c.Sender.EntityCodes()

	if err != nil {
		fmt.Printf("Ошибка загрузки сущностей: %v\n", err)
	}
	for _, val := range entCodes {
		c.rl.SetEtypeName(val.Etype, val.Name)
	}

	// Инициализация описаний полей сущностей
	for _, val := range entCodes {
		fields, err := c.Sender.Fields(val.Etype)
		if err != nil {
			fmt.Printf("Ошибка загрузки полей с описаниями: %v\n", err)
		}
		c.rl.MakeFieldsDescription(fields)
	}

	/************** Основная логика ************/
	go func() {
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

		stopChan <- true
	}()

	<-stopChan

	fmt.Println("\nПрограмма завершена!")
	return nil
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

// FilestorageDir получение директории для хранения файлов, полученных с сервера
func FilestorageDir() (string, error) {
	wd, _ := os.Getwd()
	uploadDir := wd + "/" + constants.FileStorage

	if strings.Contains(wd, "/internal/") {
		parts := strings.Split(wd, "internal")
		uploadDir = parts[0] + "cmd/client/" + constants.FileStorage
	}

	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return uploadDir, nil
}
