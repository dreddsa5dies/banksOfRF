package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LindsayBradford/go-dbf/godbf"
	"github.com/mholt/archiver"
	"github.com/opesun/goquery"
	"github.com/tealeg/xlsx"
)

const (
	// исходный url получения форм
	url         = `http://cbr.ru/credit/forms.asp`
	urlDownload = `http://cbr.ru/credit/`
)

var (
	dateSave, formName string
	delimiter          = strings.Repeat("=", 40)
)

// определение даты последнего обновления
func init() {
	// какая дата сейчас
	dateNow := time.Now()
	year, mounth, _ := dateNow.Date()
	var strMounth string
	if int(mounth) < 10 {
		strMounth = "0" + strconv.Itoa(int(mounth))
	} else {
		strMounth = strconv.Itoa(int(mounth))
	}
	dateSave = "01." + strMounth + "." + strconv.Itoa(year)
	formName = strconv.Itoa(year) + strMounth + `01.rar`
}

func main() {
	// err := getDataForm()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := unrarForms()
	if err != nil {
		log.Fatal(err)
	}

	// err = readDBF()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// успешное завершение
	os.Exit(0)
}

// открытие и декодирование DBF в exel
func readDBF() error {
	fmt.Println(delimiter)
	// открываем директорию
	fmt.Printf("Чтение DBF: %v\n", dateSave)

	dh, err := os.Open("./" + dateSave)
	if err != nil {
		return fmt.Errorf("Ошибка открытия папки с архивами: %v", err)
	}
	defer dh.Close()

	// считывание списка файлов
	for {
		fis, err := dh.Readdir(10)
		if err == io.EOF {
			break
		}
		// обход всех файлов
		for _, fi := range fis {
			if fi.IsDir() {
				dbfDir, err := os.Open("./" + dateSave + "/" + fi.Name())
				if err != nil {
					return fmt.Errorf("Ошибка открытия папки с архивами: %v", err)
				}
				defer dbfDir.Close()
				fmt.Printf("Чтение: %v\n", dbfDir.Name())
				for {
					f, err := dbfDir.Readdir(10)
					if err == io.EOF {
						break
					}
					// обход всех файлов
					for _, files := range f {
						switch {
						// 3. По форме 101
						case strings.Contains(files.Name(), "B1.DBF"):
							dbfToXLSX(dbfDir.Name(), files.Name())
						// 4. По форме 102
						case strings.Contains(files.Name(), "_P1.DBF"):
							dbfToXLSX(dbfDir.Name(), files.Name())
						// 5. По форме 123
						case strings.Contains(files.Name(), "123D.DBF"):
							dbfToXLSX(dbfDir.Name(), files.Name())
						// 6. По форме 135
						case strings.Contains(files.Name(), "_135_3.dbf"):
							dbfToXLSX(dbfDir.Name(), files.Name())
						case strings.Contains(files.Name(), "_135_4.dbf"):
							dbfToXLSX(dbfDir.Name(), files.Name())
						}
					}
				}
			}
		}
	}

	fmt.Println("Чтение DBF готово")
	fmt.Println(delimiter)

	return nil
}

// DBF to XLS
func dbfToXLSX(dbfDir, files string) error {
	fmt.Printf("Обработка %v\n", files)
	// NOT WORK: csKOI8R windows-1251 windows-1250
	// WORK - DOS 866 - Russian OEM
	dbfTable, err := godbf.NewFromFile(dbfDir+"/"+files, "866")
	if err != nil {
		return fmt.Errorf("Ошибка открытия DBF: %v", err)
	}

	// переменные для сохранения в XLSX
	// файл
	var file *xlsx.File
	// страница
	var sheet *xlsx.Sheet
	// строка
	var row *xlsx.Row

	// создаем новый файл
	file = xlsx.NewFile()
	// добавляем страницу
	sheet, err = file.AddSheet("Sheet")
	if err != nil {
		return fmt.Errorf("Ошибка добавления страницы %v", err)
	}

	switch {
	// 3. По форме 101
	case strings.Contains(files, "B1.DBF"):
		// обход по всей таблице
		for i := 0; i <= dbfTable.NumberOfRecords()-1; i++ {
			// добавление строки в XLS
			row = sheet.AddRow()
			for y := 0; y <= len(dbfTable.FieldNames())-1; y++ {
				// в третьей графе изменения счетов
				if y != 2 {
					// добавление значения в ячейку
					row.AddCell().SetString(dbfTable.FieldValue(i, y))
				} else {
					switch {
					case strings.Compare(dbfTable.FieldValue(i, y), "405") == 0:
						row.AddCell().SetString("40502")
					case strings.Compare(dbfTable.FieldValue(i, y), "406") == 0:
						row.AddCell().SetString("40602")
					case strings.Compare(dbfTable.FieldValue(i, y), "407") == 0:
						row.AddCell().SetString("40702")
					case strings.Compare(dbfTable.FieldValue(i, y), "408.1") == 0:
						row.AddCell().SetString("40817")
					case strings.Compare(dbfTable.FieldValue(i, y), "408.2") == 0:
						row.AddCell().SetString("40820")
					default:
						row.AddCell().SetString(dbfTable.FieldValue(i, y))
					}
				}
			}
		}
	// 4. По форме 102
	case strings.Contains(files, "_P1.DBF"):
		// обход по всей таблице
		for i := 0; i <= dbfTable.NumberOfRecords()-1; i++ {
			// добавление строки в XLS
			row = sheet.AddRow()
			for y := 0; y <= len(dbfTable.FieldNames())-1; y++ {
				// добавление значения в ячейку
				row.AddCell().SetString(dbfTable.FieldValue(i, y))
			}
		}
	// 5. По форме 123
	case strings.Contains(files, "123D.DBF"):
		// обход по всей таблице
		for i := 0; i <= dbfTable.NumberOfRecords()-1; i++ {
			// добавление строки в XLS
			row = sheet.AddRow()
			for y := 0; y <= len(dbfTable.FieldNames())-1; y++ {
				// добавление значения в ячейку
				row.AddCell().SetString(dbfTable.FieldValue(i, y))
			}
			// добавление даты
			row.AddCell().SetString(dateSave)
		}
	// 6. По форме 135
	case strings.Contains(files, "_135_3.dbf"):
		// обход по всей таблице
		for i := 0; i <= dbfTable.NumberOfRecords()-1; i++ {
			// добавление строки в XLS
			row = sheet.AddRow()
			for y := 0; y <= len(dbfTable.FieldNames())-1; y++ {
				// добавление значения в ячейку
				row.AddCell().SetString(dbfTable.FieldValue(i, y))
			}
			// добавление даты
			row.AddCell().SetString(dateSave)
		}
	case strings.Contains(files, "_135_4.dbf"):
		// обход по всей таблице
		for i := 0; i <= dbfTable.NumberOfRecords()-1; i++ {
			// добавление строки в XLS
			row = sheet.AddRow()
			for y := 0; y <= len(dbfTable.FieldNames())-1; y++ {
				// добавление значения в ячейку
				row.AddCell().SetString(dbfTable.FieldValue(i, y))
			}
		}
	}

	// сохранение
	err = file.Save("./" + dateSave + "/" + strings.TrimRight(files, ".DBF") + ".xlsx")
	if err != nil {
		return fmt.Errorf("Ошибка сохранения файла %v", err)
	}
	// сохраню лучше в корень с архивами
	fmt.Printf("Сохранение в %v\n", "./"+dateSave+"/"+strings.TrimRight(files, ".DBF")+".xlsx")
	time.Sleep(5 * time.Second)

	return nil
}

// разархивирование все форм
func unrarForms() error {
	fmt.Println(delimiter)
	// открываем директорию
	fmt.Printf("Разархивирование: %v\n", dateSave)
	dh, err := os.Open("./" + dateSave)
	if err != nil {
		return fmt.Errorf("Ошибка открытия папки с архивами: %v", err)
	}
	defer dh.Close()

	// считывание списка файлов
	for {
		fis, err := dh.Readdir(10)
		if err == io.EOF {
			break
		}
		// обход всех файлов
		for _, fi := range fis {
			// если это архив
			if strings.Contains(fi.Name(), ".rar") {
				fmt.Printf("Открытие архива: %v\n", fi.Name())

				// куда сохранить
				saveRarForms := "./" + dateSave + "/" + strings.TrimRight(fi.Name(), ".rar")
				fmt.Printf("Создание папки для хранения: %v\n", strings.TrimRight(fi.Name(), ".rar"))
				err = os.Mkdir(saveRarForms, 0775)
				if err != nil {
					// тут не возвращаю, т.к. может быть уже создана папка
					log.Printf("Ошибка создания папки сохранения: %v", err)
				}
				// и сохраняю все
				err := archiver.Rar.Open("./"+dateSave+"/"+fi.Name(), saveRarForms)
				if err != nil {
					return fmt.Errorf("Ошибка разархивирования: %v", err)
				}
			}
		}
	}

	fmt.Println("Разархивирование готово")
	fmt.Println(delimiter)

	return nil
}

// загрузка последних данных
func getDataForm() error {
	fmt.Println(delimiter)
	// запрос по url
	resp, err := http.Get(url)
	fmt.Printf("Загружается страница: %v\n", url)
	if err != nil {
		return fmt.Errorf("Ошибка загрузки %v", err)
	}
	// отложенное закрытие коннекта
	defer resp.Body.Close()

	// парсинг ответа
	x, err := goquery.Parse(resp.Body)

	fmt.Println("Поиск данных за: ", dateSave)

	// ищу ссылочки на формы
	var urls []string

	// совпадения на все формы по текущему месяцу
	regLLink, _ := regexp.Compile(`forms\/1\d\d-` + formName)
	for _, i := range x.Find("a").Attrs("href") {
		if regLLink.MatchString(i) {
			urls = append(urls, i)
		}
	}

	fmt.Printf("Найдены базы: %v\n", len(urls))
	for _, i := range urls {
		fmt.Printf("%v ", i)
	}
	fmt.Println()

	fmt.Printf("Создание папки для хранения: %v\n", dateSave)
	err = os.Mkdir("./"+dateSave, 0775)
	if err != nil {
		// тут не возвращаю, т.к. может быть уже создана папка
		log.Printf("Ошибка создания папки сохранения: %v", err)
	}

	fmt.Println("Скачивание баз...")
	for _, i := range urls {
		bodyForm, err := http.Get(urlDownload + i)
		fmt.Printf("Загружается: %v\n", i)
		if err != nil {
			return fmt.Errorf("Ошибка загрузки одной из форм: %v", err)
		}
		defer bodyForm.Body.Close()

		// запись ответа в переменную
		form, err := ioutil.ReadAll(bodyForm.Body)
		if err != nil {
			return fmt.Errorf("Ошибка записи ответа в переменную: %v", err)
		}

		// создание файла для загрузки
		fileName := strings.TrimLeft(i, "forms/")
		name := "./" + dateSave + "/" + fileName
		fSave, err := os.Create(name)
		if err != nil {
			return fmt.Errorf("Ошибка создания файла для загрузки: %v", err)
		}
		defer fSave.Close()

		// сохранение
		fmt.Printf("Сохрание в:	%v\n", name)
		fSave.Write(form)
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Загрузка завершена")
	fmt.Println(delimiter)

	return nil
}