package domain

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
)

func (c *GophKeepClient) Base(entCodes []*EntityCode) (string, error) {

	//cmd := exec.Command("clear")
	//cmd.Stdout = os.Stdout
	//cmd.Run()

	fmt.Println("Доступна работа со следующими объектами:")
	for i, val := range entCodes {
		fmt.Println(fmt.Sprintf("[%v] %v", i+1, val.Name))
	}

	var objStr string
	var err error
	for {
		objStr, err = c.rl.input("Выберите номер объекта:", "required,number", `{"required": "Не может быть пустым", "number": "Только число"}`)
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
	var doStr string
	for {
		for {
			doStr, err = c.rl.input(">>", "required,number", `{"required": "Не может быть пустым", "number": "Только число"}`)
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
			for _, val := range c.rl.fieldsGroup[entCode.Etype] {
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
			for {
				fmt.Println("")
				fmt.Println("Выберите дальнейшее действие:")
				fmt.Println("[1] Добавить метаданные")
				fmt.Println("[2] Перейти к сохранению")
				addOrNext, err := c.rl.input(">>", "required,number", `{"required": "Неверный выбор", "number": "Только число"}`)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				switch addOrNext {
				case "1":
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
				case "2":
					nextTag = true
					break
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
						fmt.Println("[2] Начать заново")
						againOrSave, err := c.rl.input(">>", "required,number", `{"required": "Неверный выбор", "number": "Только число"}`)
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

							// Если бинарные данные - после заведения записи на сервере загружаем бинарник на сервер
							size, err := c.Sender.UploadCryptoBinary(id, entity.Props[0].Value)
							if err != nil {
								fmt.Println("При сохранении возникли ошибки:" + err.Error())
								return WorkAgain, err
							}

							fmt.Println(fmt.Sprintf("Данные успешно сохранены! Загружен файл размером %v байт", size))
							return WorkAgain, nil
						case "2":
							return WorkAgain, nil
						default:
							continue
						}
					}

					break
				}
			}

			break

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
				fmt.Println(fmt.Sprintf("%v. Выберите номер объекта, данные которого хотите получить:", c.rl.etypes[entCode.Etype]))

				index := 0
				mapIndexToEntityID := make(map[int]int32, len(list))

				for key, val := range list {
					index++
					fmt.Println(fmt.Sprintf("[%v] %v", index, val))
					mapIndexToEntityID[index] = key
				}

				entityID, err := c.rl.input(">>", "required,number", `{"required": "Неверный выбор", "number": "Только число"}`)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				eID, err := strconv.Atoi(entityID)
				if err != nil {
					fmt.Println("Неверный номер!")
					continue
				}

				ent, err := c.Sender.Entity(mapIndexToEntityID[eID])
				// Если бинарные данные - скачиваем файл
				fd := &BinaryFileProperty{}
				err = json.Unmarshal([]byte(ent.Props[0].Value), fd)

				pathDownload, err := c.Sender.DownloadCryptoBinary(int32(eID), path.Base(fd.Clientname))
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

				for {
					fmt.Println("")
					fmt.Println("Выберите дальнейшее действие:")
					fmt.Println("[1] Изменить")
					fmt.Println("[2] Удалить")
					fmt.Println("[0] Начать все сначала")
					againOrSave, err := c.rl.input(">>", "required,number", `{"required": "Неверный выбор", "number": "Только число"}`)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					switch againOrSave {
					case "1":
						fmt.Println("Не реализовано")
					case "2":
						fmt.Println("Не реализовано")
					case "0":
						return WorkAgain, nil
					default:
						continue
					}
				}

			}

		default:
			fmt.Println("Неверный выбор!")
			continue
		}

		break
	}

	return WorkAgain, nil
}
