package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/chzyer/readline"
	"github.com/go-playground/validator/v10"
)

type CLIReader struct {
	*readline.Instance
	validator   *validator.Validate
	passwordCfg *readline.Config
	etypes      map[string]string   // справочник типов сущностей (код_сущности: наименование)
	fieldsByID  map[int32]*Field    // карта описаний полей сущности с ключом по ID поля из таблицы fields
	fieldsGroup map[string][]*Field // карта описаний полей сущности,сгруппированных по типу сущности (card, logopas, text, binary и т.д.)
}

const (
	loopBreak    string = "break"
	loopContinue string = "continue"
	loopNone     string = "none" // текущий цикл будет продолжен в штатном режиме
)

//type validateMessages map[string]string `json:""`

func NewCLIReadline(cfg *readline.Config) (*CLIReader, error) {
	rl, err := readline.NewEx(cfg)
	if err != nil {
		return nil, err
	}
	//rl.CaptureExitSignal()

	passwordCfg := rl.GenPasswordConfig()
	passwordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		rl.SetPrompt(fmt.Sprintf("Введите пароль(%v): ", len(line)))
		rl.Refresh()
		return nil, 0, false
	})

	vdr := validator.New(validator.WithRequiredStructEnabled())

	cli := &CLIReader{
		rl,
		vdr,
		passwordCfg,
		make(map[string]string),
		make(map[int32]*Field),
		make(map[string][]*Field),
	}

	return cli, nil
}

// MakeFieldsDescription Формирование карт описаний полей сущностей
func (r *CLIReader) MakeFieldsDescription(fields []*Field) {
	for _, val := range fields {
		r.fieldsByID[val.Id] = val
		r.fieldsGroup[val.Etype] = append(r.fieldsGroup[val.Etype], val)
	}
}

// Registration Ввод регистрационных данных
// возвращает введенные логин и пароль
func (r *CLIReader) Registration() (string, string, error) {
	var pswd, pswd2 []byte

	login, err := r.input("Регистрационный логин:", "required", `{"required": "Логин не может быть пустым"}`)
	if err != nil {
		if err.Error() == "interrupt" {
			r.Writeln(err.Error())
			r.Close()
		}
	}

	for {
		pswd, err = r.ReadPasswordWithConfig(r.passwordCfg)

		if r.interrupt(string(pswd), err) == loopBreak {
			return "", "", fmt.Errorf("interrupt")
		}

		if len(pswd) == 0 {
			r.Writeln("Пароль не может быть пустым!")
			continue
		}

		r.Writeln("Повторите ввод пароля")
		pswd2, err = r.ReadPasswordWithConfig(r.passwordCfg)
		if len(pswd2) == 0 {
			r.Writeln("Пароль не может быть пустым!")
			continue
		}

		if string(pswd) != string(pswd2) {
			r.Writeln("Пароли должны совпадать!")
			continue
		}

		break
	}

	return login, string(pswd), nil
}

// login Ввод логина для аутентификации
func (r *CLIReader) Login() (string, string, error) {
	defer r.SetPrompt("")

	var login, password string

	for {
		r.SetPrompt("login:")
		var err error
		login, err = r.Readline()

		if r.interrupt(login, err) == loopBreak {
			return "", "", fmt.Errorf("interrupt")
		}

		login = strings.TrimSpace(login)
		if err != nil {
			r.Writeln(err.Error())
			continue
		}
		if len(login) == 0 {
			r.Writeln("Логин не может быть пустым!")
			continue
		}

		break
	}

	for {
		r.SetPrompt("password:")
		pswd, err := r.ReadPasswordWithConfig(r.passwordCfg)
		password = string(pswd)

		if r.interrupt(password, err) == loopBreak {
			return "", "", fmt.Errorf("interrupt")
		}

		password = strings.TrimSpace(password)
		if err != nil {
			r.Writeln(err.Error())
			continue
		}
		if len(password) == 0 {
			r.Writeln("Пароль не может быть пустым!")
			continue
		}

		break
	}

	return login, password, nil
}

// GetEtypeName получение названия типа сущности по коду
func (r *CLIReader) GetEtypeName(etype string) string {
	return r.etypes[etype]
}
func (r *CLIReader) SetEtypeName(etype string, name string) {
	r.etypes[etype] = name
}

func (r *CLIReader) GetField(fieldID int32) *Field {
	return r.fieldsByID[fieldID]
}

func (r *CLIReader) GetFieldsGroup(etype string) []*Field {
	return r.fieldsGroup[etype]
}

// Логика ввода строки значения поля
// Будет запрашивать ввод до тех пор пока не пройдет валидацию
// validateMessages - правила валидации в формате JSON.
// Например: `{"required": "Логин не может быть пустым"}`
// или `{"gt=1": "Значение поля должно быть больше <param>"}` - вместо <param> будет подставлено значение 1
func (r *CLIReader) input(prompt string, validateRules string, validateMessages string) (string, error) {
	var value string
	var err error

	var vm map[string]string
	err = json.Unmarshal([]byte(validateMessages), &vm)
	if err != nil {
		return "", err
	}

	for {
		r.SetPrompt(prompt)
		value, err = r.Readline()

		if r.interrupt(value, err) == loopBreak {
			return "", readline.ErrInterrupt
		}

		if err != nil {
			r.Writeln(err.Error())
			continue
		}

		if validateRules == "" {
			break
		}

		errs := r.validator.Var(value, validateRules)
		errors, okAssert := errs.(validator.ValidationErrors)
		if okAssert {
			for _, err := range errors {
				message := err.Error()
				if val, ok := vm[err.Tag()]; ok {
					message = val
					message = strings.Replace(message, "<param>", err.Param(), -1)
					r.Writeln(message)
				}
			}
			if len(errors) > 0 {
				continue
			}
		}

		break
	}

	return value, nil
}

// Редактирование
func (r *CLIReader) edit(prompt string, what string, validateRules string, validateMessages string) (string, error) {
	var value string
	var err error

	var vm map[string]string
	err = json.Unmarshal([]byte(validateMessages), &vm)
	if err != nil {
		return "", err
	}

	for {
		r.SetPrompt(prompt)
		value, err = r.ReadlineWithDefault(what)

		if r.interrupt(value, err) == loopBreak {
			return "", readline.ErrInterrupt
		}

		if err != nil {
			r.Writeln(err.Error())
			continue
		}

		if validateRules == "" {
			break
		}

		errs := r.validator.Var(value, validateRules)
		errors, okAssert := errs.(validator.ValidationErrors)
		if okAssert {
			for _, err := range errors {
				message := err.Error()
				if val, ok := vm[err.Tag()]; ok {
					message = val
					message = strings.Replace(message, "<param>", err.Param(), -1)
					r.Writeln(message)
				}
			}
			if len(errors) > 0 {
				continue
			}
		}

		break
	}

	return value, nil
}

func (r *CLIReader) Writeln(str string) {
	r.Write([]byte(str + "\n"))
}

// interrupt прерывание текущего ввода
// нужно для прерывания или продолжения текущего цикла ввода
func (r *CLIReader) interrupt(line string, err error) string {
	if err == readline.ErrInterrupt {
		if len(line) == 0 {
			return loopBreak
		} else {
			return loopContinue
		}
	} else if err == io.EOF {
		return loopBreak
	}

	return loopNone
}

func (r *CLIReader) Close() error {
	return r.Instance.Close()
}

func (r *CLIReader) Stderr() io.Writer {
	return r.Instance.Stderr()
}

//var Completer = readline.NewPrefixCompleter(
//	readline.PcItem("mode",
//		readline.PcItem("vi"),
//		readline.PcItem("emacs"),
//	),
//	readline.PcItem("login"),
//	readline.PcItem("say",
//		readline.PcItemDynamic(ListFiles("./"),
//			readline.PcItem("with",
//				readline.PcItem("following"),
//				readline.PcItem("items"),
//			),
//		),
//		readline.PcItem("hello"),
//		readline.PcItem("bye"),
//	),
//	readline.PcItem("setprompt"),
//	readline.PcItem("setpassword"),
//	readline.PcItem("bye"),
//	readline.PcItem("help"),
//	readline.PcItem("go",
//		readline.PcItem("build", readline.PcItem("-o"), readline.PcItem("-v")),
//		readline.PcItem("install",
//			readline.PcItem("-v"),
//			readline.PcItem("-vv"),
//			readline.PcItem("-vvv"),
//		),
//		readline.PcItem("test"),
//	),
//	readline.PcItem("sleep"),
//)

func FilterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

//func Usage(w io.Writer) {
//	io.WriteString(w, "commands:\n")
//	io.WriteString(w, Completer.Tree("    "))
//}

// Function constructor - constructs new function for listing given directory
func ListFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}
