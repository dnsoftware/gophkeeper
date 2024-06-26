INSERT INTO entity_codes (etype, name)
VALUES ('logopas', 'Логин и пароль'),
       ('card', 'Банковская карта'),
       ('text', 'Текстовые данные'),
       ('binary', 'Бинарные данные');

/* logopas */
INSERT INTO fields (etype, name, ftype, validate_rules, validate_messages)
VALUES ('logopas', 'Логин', 'string', '{"length_gt": "0"}', '{"length_gt": "Логин не может быть пустым"}'),
       ('logopas', 'Пароль', 'string', '{"length_gt": "0"}', '{"length_gt": "Пароль не может быть пустым"}');

/* card */
INSERT INTO fields (etype, name, ftype, validate_rules, validate_messages)
VALUES ('card', 'Номер банковской карты', 'string', '{"length_gte": "16", "length_lte": "19"}', '{"length_gte": "Длина номера должна быть больше или равно %value", "length_lte": "Длина номера должна быть меньше или равно %value"}'),
       ('card', 'Месяц/Год (mm/yy) до которого действует карта', 'string', '{"regex": "^\d\d/\d\d$"}', '{"regex": "Неправильный формат даты окончания действия (mm/dd)"}'),
       ('card', 'Код проверки подлинности', 'string', '{"length_equal": "3", "regex": "^\d{3}$"}', '{"length_equal": "Код должен состоять из трех цифр", "regex": "Код должен состоять из трех цифр"}');

/* text */
INSERT INTO fields (etype, name, ftype, validate_rules, validate_messages)
VALUES ('text', 'Произвольные текстовые данные', 'path', '{"length_gt": "0"}', '{"length_gt": "Укажите путь к текстовому файлу"}');

/* binary */
INSERT INTO fields (etype, name, ftype, validate_rules, validate_messages)
VALUES ('binary', 'Произвольные бинарные данные', 'path', '{"length_gt": "0"}', '{"length_gt": "Укажите путь к бинарному файлу"}');
