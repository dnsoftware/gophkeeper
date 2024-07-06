package domain

import (
	"io"
	"testing"
	"time"

	"github.com/chzyer/readline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistration(t *testing.T) {

	r, w := io.Pipe()

	rl, err := NewCLIReadline(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    nil,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		Stdin:           r,
		//Stdout:          w,

		HistorySearchFold:   true,
		FuncFilterInputRune: FilterInput,
	})

	// позитивный тест, логин и пароль проходят валидацию
	login, password := "", ""
	go func() {
		login, password, err = rl.Registration()
		return
	}()
	time.Sleep(1 * time.Second)

	w.Write([]byte("logintest\n"))
	sleep()
	w.Write([]byte("passwordtest\n"))
	sleep()
	w.Write([]byte("passwordtest\n"))
	sleep()

	require.Equal(t, "logintest", login)
	require.Equal(t, "passwordtest", password)

	// негативный тест, логин пустой
	login, password = "", ""
	go func() {
		login, password, err = rl.Registration()
		return
	}()
	time.Sleep(1 * time.Second)
	w.Write([]byte("\n"))
	sleep()
	w.Write([]byte("passwordtest\n"))
	sleep()
	w.Write([]byte("passwordtest\n"))
	sleep()

	require.Equal(t, "", login)
	require.Equal(t, "", password)

	// негативный тест, пароли не совпадают
	login, password = "", ""
	go func() {
		login, password, err = rl.Registration()
		return
	}()
	time.Sleep(1 * time.Second)
	w.Write([]byte("logintest\n"))
	sleep()
	w.Write([]byte("passwordtest\n"))
	sleep()
	w.Write([]byte("passwordbad\n"))
	sleep()

	require.Equal(t, "", login)
	require.Equal(t, "", password)

	// негативный тест, пароль пустой
	login, password = "", ""
	go func() {
		login, password, err = rl.Registration()
		return
	}()
	time.Sleep(1 * time.Second)
	w.Write([]byte("logintest\n"))
	sleep()
	w.Write([]byte("\n"))
	sleep()
	w.Write([]byte("\n"))
	sleep()

	require.Equal(t, "", login)
	require.Equal(t, "", password)

	require.NoError(t, err)

}

func TestLogin(t *testing.T) {

	r, w := io.Pipe()

	rl, err := NewCLIReadline(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    nil,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		Stdin:           r,
		//Stdout:          w,

		HistorySearchFold:   true,
		FuncFilterInputRune: FilterInput,
	})

	// позитивный тест, логин и пароль проходят валидацию
	login, password := "", ""
	go func() {
		login, password, err = rl.Login()
		return
	}()
	time.Sleep(1 * time.Second)

	w.Write([]byte("logintest\n"))
	sleep()
	w.Write([]byte("passwordtest\n"))
	sleep()

	require.Equal(t, "logintest", login)
	require.Equal(t, "passwordtest", password)

	// негативный тест, логин пустой
	login, password = "", ""
	go func() {
		login, password, err = rl.Login()
		return
	}()
	time.Sleep(1 * time.Second)
	w.Write([]byte("\n"))
	sleep()
	w.Write([]byte("passwordtest\n"))
	sleep()

	require.Equal(t, "", login)
	require.Equal(t, "", password)

	// негативный тест, пароль пустой
	login, password = "", ""
	go func() {
		login, password, err = rl.Login()
		return
	}()
	time.Sleep(1 * time.Second)
	w.Write([]byte("\n"))
	sleep()

	require.Equal(t, "", login)
	require.Equal(t, "", password)

	require.NoError(t, err)

}

func TestEdit(t *testing.T) {
	r, w := io.Pipe()

	rl, _ := NewCLIReadline(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    nil,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		Stdin:           r,
		//Stdout:          w,

		HistorySearchFold:   true,
		FuncFilterInputRune: FilterInput,
	})

	var err error
	newval := ""
	go func() {
		newval, err = rl.edit("test", "test", "required", `{"required": "Не может быть пустым"}`)
		assert.NoError(t, err)
	}()
	sleep()
	w.Write([]byte("new\n"))
	sleep()

	assert.NoError(t, err)
	assert.Equal(t, "testnew", newval)

	var newval2 string
	go func() {
		newval2, err = rl.edit("test", "", "required", `{"required": "Не может быть пустым"}`)
		assert.NoError(t, err)
	}()
	sleep()
	w.Write([]byte("\n"))
	sleep()
	require.Equal(t, "", newval2)
	require.NoError(t, err)
	w.Write([]byte("noempty\n"))

}

func TestGet(t *testing.T) {
	r, _ := io.Pipe()

	rl, _ := NewCLIReadline(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    nil,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		Stdin:           r,
		//Stdout:          w,

		HistorySearchFold:   true,
		FuncFilterInputRune: FilterInput,
	})

	rl.SetEtypeName("card", "карта")
	name := rl.GetEtypeName("card")
	assert.Equal(t, "карта", name)

	fields := []*Field{&Field{
		Id:               1,
		Name:             "test",
		Etype:            "test",
		Ftype:            "test",
		ValidateRules:    "",
		ValidateMessages: "",
	}}
	rl.MakeFieldsDescription(fields)
	f := rl.GetField(1)
	assert.Equal(t, fields[0], f)
	f2 := rl.GetFieldsGroup("test")
	assert.Equal(t, fields, f2)
}

func sleep() {
	time.Sleep(100 * time.Millisecond)
}
