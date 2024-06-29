package domain

import (
	"io"
	"testing"
	"time"

	"github.com/chzyer/readline"
	"github.com/stretchr/testify/require"
)

func TestRegistration(t *testing.T) {

	r, w := io.Pipe()

	rl, err := NewCLIReadline(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		Stdin:           r,
		//Stdout:          w,

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
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
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		Stdin:           r,
		//Stdout:          w,

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
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

func sleep() {
	time.Sleep(100 * time.Millisecond)
}
