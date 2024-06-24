package handlers

//// UserService интерфейс для работы с регистрацией и аутентификацией/авторизацией
//type UserService interface {
//	// Registration регистрация нового пользователя. Возвращает токен доступа в случае удачи и ошибку, если что-то пошло не так
//	Registration(ctx context.Context, login string, password string, repeatPassword string) (string, error)
//
//	// Login вход пользователя. Возвращает токен доступа в случае удачи и ошибку, если что-то пошло не так
//	Login(ctx context.Context, login string, password string) (string, error)
//}

//// EntityCodeService интерфейс для работы со справочником сущностей
//type EntityCodeService interface {
//
//	// EntityCodes запрос списка доступных к добавлению типов сущностей (таблица entity_codes)
//	EntityCodes(ctx context.Context) (map[string]string, error)
//}

// запрос данных для добавления новой сущности (набор полей и их характеристик)
//EntityProperties(ctx context.Context, code string)

// сохранение новой сущности в базу

// запрос существующей сущности для просмотра (скачивание текстовых и бинарных файлов)

// запрос существующей сущности для редактирования

// сохранение отредактированной сущности

// запрос на удаление существующей сущности
