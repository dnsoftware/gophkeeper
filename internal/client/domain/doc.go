/*
Package domain содержит бизнес логику работы CLI клиента

### Аутентификация
Аутентификация происходит с помошью JWT токена.

Токен клиент получает сразу после успешной регистрации - можно сразу приступать к работе. И при успешном входе, после ввода логина и пароля.

Во избежание кражи, токен хранится только в оперативной памяти при работе программы. Также он имеет заданное время "жизни".

При отправке запросов на сервер перехватчик автоматически добавляет токен с ID пользователя к контексту запроса.

Пароль в базе хранится в виде хеша. Для дополнительной криптостойкости при формировании хеша используется "соль"-фрагмент случайных данных.

### Передача данных
В силу того, что файлы могут иметь большие размеры - их передача происходит в потоковом режиме gRPC. Потоки однонаправленные - от клиента к серверу при сохранении и от сервера к клиенту при получении. Размер чанков/фрагментов задается константой в коде программы.

Обмен происходит по защищенному TLS протоколу. Используются заранее сгенерированные сертификаты.

### Шифрование
Помимо защищенного канала все критичные данные шифруются симметричным AES алгоритмом.

Шифровка происходит на стороне клиента, поэтому на сервере все данные зашифрованы, а следовательно потенциальная возможность взлома сервера ничего не даст злоумышленнику.

Ключ шифрования формируется из двух частей. Первая часть - пароль пользователя, вводимый при входе в систему. Вторая часть - секретный ключ (SecretKey), который задается в конфиге.

При гипотетическом перехвате пароля на стороне сервера злоумышленник не сможет расшифровать данные из-за отсутствия секретного ключа. Заполучив секретный ключ на стороне клиента, злоумышленник также не сможет ничего сделать из-за отсутствия пароля, который в идеале хранится только в голове пользователя!))

Перед отправкой файлы разбиваются на фрагменты и каждый фрагмент шифруется.

## Получение файлов клиентом
В момент запроса файла клиента последовательно считываются все файлы-фрагменты из соответствующей папки и в потоковом режиме передаются на клиент.

На клиенте каждый фрагмент расшифровывается и добавляется к результирующему файлу. После успешного получения файла пользователю сообщается путь к нему.

### Структура файлов Клиента
cmd/client/ - папка с бинарным исполняемым файлом

cmd/client/cert - самоподписанный корневой сертификат для клиента

cmd/client/filestorage - сюда будут загружаться файлы с сервера

cmd/client/testbinary - пара файлов для тестовых выгрузок-загрузок

internal/client - код клиента

### Валидация вводимых данных
Валидация происходит на стороне клиента.

Правила же валидации хранятся на стороне сервера для каждого поля в таблице fields.
Также в этой таблице хранятся сообщения, выдающиеся при возникновении ошибки валидации.

Таким образом для добавления новой сущности в систему достаточно на стороне сервера указать все ее поля и правила их валидации.

В момент подключения клиент загружает все описания полей доступных сущностей и правил их валидации.

### Ключи запуска клиента
-с - путь к файлу кофигурации
-a - адрес сервера
-k - секретный ключ для AES шифрования

### Интерфейс клиента
На начальном этапе будет предложено ввести логин и пароль для входа в систему. Если же таковые отсутствуют, то можно будет пройти процедуру регистрации.

После успешной авторизации будет предложен список возможных сущностей, с которыми можно работать.

Управление действиями в клиенте производится путем выбора предложенного варианта действий из пронумерованного списка.

Выбор пункта с номером "0" используется для перехода в начало, к выбору сущности из списка.

Комбинация клавиш Ctrl+C приводит к завершению работы клиента с удалением всех скачанных в течение текущего сеанса файлов. Это сделано во избежание хранения расшифрованных данных между сеансами. Предполагается, что в процессе сеанса пользователь просмотрит всю необходимую информацию или же скопирует ее в надежное место.
*/
package domain
