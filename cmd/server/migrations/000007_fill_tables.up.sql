INSERT INTO entity_codes (etype, name)
VALUES ('logopas', 'Логин и пароль'),
       ('card', 'Банковская карта'),
       ('text', 'Текстовые данные'),
       ('binary', 'Бинарные данные');

/* logopas */
INSERT INTO fields (etype, name, ftype, validate_rules, validate_messages)
VALUES ('logopas', 'Логин', 'string', 'required', '{"required": "Логин не может быть пустым"}'),
       ('logopas', 'Пароль', 'string', 'required', '{"required": "Пароль не может быть пустым"}');

/* card */
INSERT INTO fields (etype, name, ftype, validate_rules, validate_messages)
VALUES ('card', 'Номер банковской карты', 'string', 'credit_card', '{"credit_card": "Неправильный формат номера карты"}'),
       ('card', 'Месяц/Год (mm/yy) до которого действует карта', 'string', 'len=5', '{"len": "Месяц/год должны быть в формате mm/dd"}'),
       ('card', 'Код проверки подлинности', 'string', 'len=3,number', '{"len": "Код должен состоять из трех цифр", "number": "Только число"}');

/* text */
INSERT INTO fields (etype, name, ftype, validate_rules, validate_messages)
VALUES ('text', 'Произвольные текстовые данные', 'path', 'required,file', '{"requred": "Путь к файлу не может быть пустым", "file": "Файла не существует"}');

/* binary */
INSERT INTO fields (etype, name, ftype, validate_rules, validate_messages)
VALUES ('binary', 'Произвольные бинарные данные', 'path', 'required,file', '{"required": "Путь к файлу не может быть пустым", "file": "Файла не существует"}');
