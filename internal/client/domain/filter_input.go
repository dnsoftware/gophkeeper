package domain

import (
	"os"

	"github.com/chzyer/readline"

	"github.com/dnsoftware/gophkeeper/internal/constants"
)

// Filter Фильтрация ввода с консоли и завершение работы клиента по нажатию Ctrl+C
type Filter struct {
	filestorage string
	stopChan    chan bool
}

func NewFilter(filestorage string, stopChan chan bool) *Filter {
	f := Filter{
		filestorage: filestorage,
		stopChan:    stopChan,
	}

	return &f
}

func (s *Filter) FilterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false

	case constants.CharCtrlC: // завершение работы
		// удаляем содержимое директории со скачанными файлами
		_ = os.RemoveAll(s.filestorage)
		s.stopChan <- true
	}
	return r, true
}
