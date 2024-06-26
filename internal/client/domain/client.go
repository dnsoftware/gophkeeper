package domain

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

// Sender Интерфейс отправки/приема данных с сервера
type Sender interface {
	Registration(login string, password string, password2 string) (string, error)
	Login(login string, password string) (string, error)
	EntityCodes() ([]*EntityCode, error)
	Fields(etype string) ([]*Field, error)
	AddEntity(ae Entity) (int32, error)
	UploadBinary(entityId int32, file string) (int32, error)
	DownloadBinary(entityId int32, fileName string) (string, error)
	UploadCryptoBinary(entityId int32, file string) (int32, error)
	DownloadCryptoBinary(entityId int32, fileName string) (string, error)
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
	rl     *CLIReader
	Sender Sender
}

func NewGophKeepClient(sender Sender) (*GophKeepClient, error) {

	rl, err := NewCLIReadline(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		return nil, err
	}

	client := &GophKeepClient{
		rl:     rl,
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
		c.rl.Writeln("Нажмите [Enter] для входа или \"r\" для регистрации  ")
		line, err := c.rl.Readline()
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
			c.rl.Writeln(err.Error())
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
		c.rl.etypes[val.Etype] = val.Name
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

	c.rl.Writeln("")
	c.rl.Writeln("Доступна работа со следующими объектами:")
	for i, val := range entCodes {
		c.rl.Writeln(fmt.Sprintf("[%v] %v", i+1, val.Name))
	}

	var objStr string
	for {
		objStr, err = c.rl.input("Выберите номер объекта:", "required,number", `{"required": "Не может быть пустым", "number": "Только число"}`)
		if err != nil {
			c.rl.Writeln(err.Error())
			continue
		}
		break
	}
	objIndex, _ := strconv.Atoi(objStr)
	entCode := entCodes[objIndex-1]
	c.rl.Writeln("")
	c.rl.Writeln(fmt.Sprintf(`Для объекта "%v" доступны следующие действия:`, entCode.Name))
	c.rl.Writeln("[1] Добавить новый")
	c.rl.Writeln("[2] Получить сохраненный")
	var doStr string
	for {
		for {
			doStr, err = c.rl.input(">>", "required,number", `{"required": "Не может быть пустым", "number": "Только число"}`)
			if err != nil {
				c.rl.Writeln(err.Error())
				continue
			}
			break
		}

		switch doStr {
		// Добавление сущности, поочередно вводим данные в поля
		case "1":
			var props []*Property
			var metas []*Metainfo
			entity := Entity{
				Id:       0,
				UserID:   0,
				Etype:    entCode.Etype,
				Props:    nil,
				Metainfo: nil,
			}

			// Заполняем обязательные поля
			for _, val := range c.rl.fieldsGroup[entCode.Etype] {
				fieldData, err := c.rl.input(val.Name+":", val.ValidateRules, val.ValidateMessages)
				if err != nil {
					c.rl.Writeln(err.Error())
					continue
				}
				props = append(props, &Property{
					EntityId: 0,
					FieldId:  val.Id,
					Value:    fieldData,
				})
			}
			// Заполняем поля метаданных
			// Добавить или перейти дальше
			nextTag := false
			for {
				c.rl.Writeln("")
				c.rl.Writeln("Выберите дальнейшее действие:")
				c.rl.Writeln("[1] Добавить метаданные")
				c.rl.Writeln("[2] Перейти к сохранению")
				addOrNext, err := c.rl.input(">>", "required,number", `{"required": "Неверный выбор", "number": "Только число"}`)
				if err != nil {
					c.rl.Writeln(err.Error())
					continue
				}

				switch addOrNext {
				case "1":
					for {
						metaName, err := c.rl.input("Название поля метаданных:", "required", `{"required": "Укажите название поля метаданных"}`)
						if err != nil {
							c.rl.Writeln(err.Error())
							continue
						}

						metaValue, err := c.rl.input("Значение поля метаданных:", "required", `{"required": "Укажите значение поля метаданных"}`)
						if err != nil {
							c.rl.Writeln(err.Error())
							continue
						}

						metas = append(metas, &Metainfo{
							EntityId: 0,
							Title:    metaName,
							Value:    metaValue,
						})

						break
					}
				case "2":
					nextTag = true
					break
				default:
					continue
				}

				if nextTag {
					entity.Props = props
					entity.Metainfo = metas

					// Просмотр и сохранение
					c.DisplayEntity(entity)

					break
				}
			}

			break

		// Просмотр сохраненной сущности, получаем список для дальнейшего выбора
		case "2":

		default:
			c.rl.Writeln("Неверный выбор!")
			continue
		}

		break
	}

	for {
		line, err := c.rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "mode "):
			switch line[5:] {
			case "vi":
				c.rl.SetVimMode(true)
			case "emacs":
				c.rl.SetVimMode(false)
			default:
				println("invalid mode:", line[5:])
			}
		case line == "mode":
			if c.rl.IsVimMode() {
				println("current mode: vim")
			} else {
				println("current mode: emacs")
			}
		case line == "login":
			pswd, err := c.rl.ReadPassword("please enter your password: ")
			if err != nil {
				break
			}
			println("you enter:", strconv.Quote(string(pswd)))
		case line == "help":
			usage(c.rl.Stderr())
		case line == "setpassword":
			pswd, err := c.rl.ReadPasswordWithConfig(c.rl.passwordCfg)
			if err == nil {
				println("you set:", strconv.Quote(string(pswd)))
			}
		case strings.HasPrefix(line, "setprompt"):
			if len(line) <= 10 {
				log.Println("setprompt <prompt>")
				break
			}
			c.rl.SetPrompt(line[10:])
		case strings.HasPrefix(line, "say"):
			line := strings.TrimSpace(line[3:])
			if len(line) == 0 {
				log.Println("say what?")
				break
			}
			go func() {
				for range time.Tick(time.Second) {
					log.Println(line)
				}
			}()
		case line == "bye":
			goto exit
		case line == "sleep":
			log.Println("sleep 4 second")
			time.Sleep(4 * time.Second)
		case line == "":
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
exit:
}

// DisplayEntity отобразить сужность в консоли
func (c *GophKeepClient) DisplayEntity(ent Entity) {
	c.rl.Writeln("------------------------")
	c.rl.Writeln(" " + c.rl.etypes[ent.Etype])
	for _, val := range ent.Props {
		c.rl.Writeln("      " + c.rl.fieldsByID[val.FieldId].Name + ": " + val.Value)
	}
	for _, val := range ent.Metainfo {
		c.rl.Writeln("      " + val.Title + ": " + val.Value)
	}
	c.rl.Writeln("------------------------")
}

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

// Function constructor - constructs new function for listing given directory
func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("mode",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),
	readline.PcItem("login"),
	readline.PcItem("say",
		readline.PcItemDynamic(listFiles("./"),
			readline.PcItem("with",
				readline.PcItem("following"),
				readline.PcItem("items"),
			),
		),
		readline.PcItem("hello"),
		readline.PcItem("bye"),
	),
	readline.PcItem("setprompt"),
	readline.PcItem("setpassword"),
	readline.PcItem("bye"),
	readline.PcItem("help"),
	readline.PcItem("go",
		readline.PcItem("build", readline.PcItem("-o"), readline.PcItem("-v")),
		readline.PcItem("install",
			readline.PcItem("-v"),
			readline.PcItem("-vv"),
			readline.PcItem("-vvv"),
		),
		readline.PcItem("test"),
	),
	readline.PcItem("sleep"),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
