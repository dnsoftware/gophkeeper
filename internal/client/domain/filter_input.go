package domain

import (
	"os"

	"github.com/chzyer/readline"

	"github.com/dnsoftware/gophkeeper/internal/constants"
)

// Filter Фильтрация ввода с консоли и завершение работы клиента по нажатию Ctrl+C
type Filter struct {
	filestorage string    // путь к директории хранения полученных файлов (удаляется после завершения сеанса работы)
	stopChan    chan bool // канал, при поступлении в него сигнала завершения работы - удаляется папка с полученными файлами
}

// NewFilter конструктор
func NewFilter(filestorage string, stopChan chan bool) *Filter {
	f := Filter{
		filestorage: filestorage,
		stopChan:    stopChan,
	}

	return &f
}

// FilterInput Фильтрация ввода с консоли и завершение работы клиента по нажатию Ctrl+C
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
