package domain

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"

	"github.com/dnsoftware/gophkeeper/internal/constants"
)

func (c *GophKeepClient) Base(entCodes []*EntityCode) (string, error) {

	fmt.Println("Доступна работа со следующими объектами:")
	for i, val := range entCodes {
		fmt.Println(fmt.Sprintf("[%v] %v", i+1, val.Name))
	}

	var objStr string
	var err error
	for {
		objStr, err = c.rl.input("Выберите номер объекта:", "required,number", `{"required": "Не может быть пустым", "number": "Только число"}`)
		if c.rl.interrupt(objStr, err) == loopBreak {
			return WorkStop, err
		}
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		break
	}
	objIndex, _ := strconv.Atoi(objStr)
	entCode := entCodes[objIndex-1]
	fmt.Println("")
	fmt.Println(fmt.Sprintf(`Для объекта "%v" доступны следующие действия:`, entCode.Name))
	fmt.Println("[1] Добавить новый")
	fmt.Println("[2] Получить сохраненный")
	fmt.Println("[0] Начать сначала")
	var doStr string
	for {
		for {
			doStr, err = c.rl.input("Действия для объекта>>", "required,number", `{"required": "Не может быть пустым", "number": "Только число"}`)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			break
		}

		switch doStr {
		// Добавление сущности, поочередно вводим данные в поля
		case "1":
			var props []*Property
			var metas []*Metainfo
			entity := Entity{
				Id:       0,
				UserID:   0,
				Etype:    entCode.Etype,
				Props:    nil,
				Metainfo: nil,
			}

			// Заполняем обязательные поля
			for _, val := range c.rl.GetFieldsGroup(entCode.Etype) {
				fieldData, err := c.rl.input(val.Name+":", val.ValidateRules, val.ValidateMessages)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				props = append(props, &Property{
					EntityId: 0,
					FieldId:  val.Id,
					Value:    fieldData,
				})
			}
			// Заполняем поля метаданных
			// Добавить или перейти дальше
			nextTag := false
			prefix := "Метаданные"
			for {
				fmt.Println("")
				fmt.Println("Выберите дальнейшее действие:")
				fmt.Println("[1] Добавить метаданные")
				fmt.Println("[2] Перейти к сохранению")
				fmt.Println("[0] Начать сначала")
				addOrNext, err := c.rl.input(prefix+" или сохранение>>", "required,number", `{"required": "Неверный выбор", "number": "Только число"}`)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				switch addOrNext {
				case "1":
					metas = c.createMetainfo(metas)
					prefix = "Еще метаданные"
				case "2":
					nextTag = true
					break
				case "0":
					return WorkAgain, nil

				default:
					continue
				}

				if nextTag {
					entity.Props = props
					entity.Metainfo = metas

					// Просмотр и сохранение
					c.DisplayEntity(entity)

					for {
						fmt.Println("")
						fmt.Println("Выберите дальнейшее действие:")
						fmt.Println("[1] Сохранить")
						fmt.Println("[0] Начать заново")
						againOrSave, err := c.rl.input("Сохранить или заново>>", "required,number", `{"required": "Неверный выбор", "number": "Только число"}`)
						if err != nil {
							fmt.Println(err.Error())
							continue
						}

						switch againOrSave {
						case "1":
							id, err := c.Sender.AddEntity(entity)
							if err != nil || id <= 0 {
								return WorkAgain, err
							}

							if entCode.Etype == constants.BinaryEntity || entCode.Etype == constants.TextEntity {
								// Если бинарные данные - после заведения записи на сервере загружаем бинарник на сервер
								size, err := c.Sender.UploadCryptoBinary(id, entity.Props[0].Value)
								if err != nil {
									fmt.Println("При сохранении возникли ошибки:" + err.Error())
									return WorkAgain, err
								}
								fmt.Println(fmt.Sprintf("Данные успешно сохранены! Загружен файл размером %v байт", size))
							} else {
								fmt.Println(fmt.Sprintf("Данные успешно сохранены!"))
							}

							return WorkAgain, nil
						case "0":
							return WorkAgain, nil
						case "-1":
							return WorkStop, nil
						default:
							continue
						}
					}

				}
			}

		// Просмотр сохраненной сущности, получаем список для дальнейшего выбора
		case "2":
			list, err := c.Sender.EntityList(entCode.Etype)
			if err != nil {
				fmt.Println(err.Error())
			}

			if len(list) == 0 {
				fmt.Println("Нет ни одного объекта!")
				return WorkAgain, nil
			}

			for {
				fmt.Println("")
				fmt.Println(fmt.Sprintf("%v. Выберите номер объекта, данные которого хотите получить:", c.rl.GetEtypeName(entCode.Etype)))

				// соответствие межну консольными номерами сущностей и реальными идентификаторами
				index := 0
				mapIndexToEntityID := make(map[int]int32, len(list))
				for key, val := range list {
					index++
					fmt.Println(fmt.Sprintf("[%v] %v", index, val))
					mapIndexToEntityID[index] = key
				}

				entityIndex, err := c.rl.input("Просмотр объекта>>", "required,number", `{"required": "Неверный выбор", "number": "Только число"}`)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				entIndex, err := strconv.Atoi(entityIndex)
				if err != nil {
					fmt.Println("Неверный номер!")
					continue
				}

				entityID := mapIndexToEntityID[entIndex]
				ent, err := c.Sender.Entity(entityID)

				// Если бинарные данные или произвольный текст - скачиваем файл
				if entCode.Etype == constants.BinaryEntity || entCode.Etype == constants.TextEntity {
					fd := &BinaryFileProperty{}
					err = json.Unmarshal([]byte(ent.Props[0].Value), fd)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					pathDownload, err := c.Sender.DownloadCryptoBinary(entityID, path.Base(fd.Clientname))
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					ent.Props[0].Value = pathDownload

					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					c.DisplayEntityBinary(*ent, pathDownload)

				} else {

					c.DisplayEntity(*ent)

				}

				for {
					fmt.Println("")
					fmt.Println("Выберите дальнейшее действие:")
					fmt.Println("[1] Изменить")
					fmt.Println("[2] Удалить")
					fmt.Println("[0] Начать все сначала")
					againOrSave, err := c.rl.input("Действия с объектом>>", "required,number", `{"required": "Неверный выбор", "number": "Только число"}`)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					switch againOrSave {
					case "1":

						for propKey, propVal := range ent.Props {
							field := c.rl.GetField(propVal.FieldId)
							if ent.Etype == constants.BinaryEntity || ent.Etype == constants.TextEntity {
								ent.Props[propKey].Value, err = c.rl.edit(field.Name+":", "", field.ValidateRules, field.ValidateMessages)
							} else {
								ent.Props[propKey].Value, err = c.rl.edit(field.Name+":", propVal.Value, field.ValidateRules, field.ValidateMessages)
							}
						}

						for metaKey, metaVal := range ent.Metainfo {
							ent.Metainfo[metaKey].Title, err = c.rl.edit("Название метаданных:", metaVal.Title, "required", `{"required": "Укажите название поля метаданных"}`)
							ent.Metainfo[metaKey].Value, err = c.rl.edit("Значение метаданных:", metaVal.Value, "required", `{"required": "Укажите значение поля метаданных"}`)
						}

						id, err := c.Sender.SaveEntity(*ent)
						if err != nil || id <= 0 {
							return WorkAgain, err
						}

						if entCode.Etype == constants.BinaryEntity || entCode.Etype == constants.TextEntity {
							// Если бинарные данные - после редактирования записи на сервере загружаем бинарник на сервер
							size, err := c.Sender.UploadCryptoBinary(id, ent.Props[0].Value)
							if err != nil {
								fmt.Println("При изменении возникли ошибки:" + err.Error())
								return WorkAgain, err
							}
							fmt.Println(fmt.Sprintf("Данные успешно изменены! Загружен файл размером %v байт", size))
						} else {
							fmt.Println(fmt.Sprintf("Данные успешно изменены!"))
						}

						return WorkAgain, nil

					case "2":
						fmt.Println("Не реализовано")
					case "0":
						return WorkAgain, nil
					default:
						continue
					}
				}

			}
		case "0":
			return WorkAgain, nil

		default:
			fmt.Println("Неверный выбор!")
			continue
		}

		break
	}

	return WorkAgain, nil
}

func (c *GophKeepClient) createMetainfo(metas []*Metainfo) []*Metainfo {
	for {
		metaName, err := c.rl.input("Название поля метаданных:", "required", `{"required": "Укажите название поля метаданных"}`)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		metaValue, err := c.rl.input("Значение поля метаданных:", "required", `{"required": "Укажите значение поля метаданных"}`)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		metas = append(metas, &Metainfo{
			EntityId: 0,
			Title:    metaName,
			Value:    metaValue,
		})

		break
	}

	return metas

}
