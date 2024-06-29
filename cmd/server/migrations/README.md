# Миграции

Используем пакет https://github.com/golang-migrate/migrate

Для создания новой миграции из командной строки запускаем:

make migration_add name=<migration_file_name>

Создадутся два файла:

    index_migration_file_name.up.sql
    index_migration_file_name.down.sql

В index_migration_file_name.up.sql пишем SQL запрос по модификации БД.
