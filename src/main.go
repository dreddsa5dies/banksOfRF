package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/LindsayBradford/go-dbf/godbf"
	"github.com/opesun/goquery"
)

const (
	// исходный url получения форм
	url         = `http://cbr.ru/credit/forms.asp`
	urlDownload = `http://cbr.ru/credit/`
)

func main() {
	getDataForm()
}

// загрузка последних данных
func getDataForm() {
	// запрос по url
	resp, err := http.Get(url)
	log.Printf("Загружается страница:	%v", url)
	if err != nil {
		log.Fatalf("Ошибка загрузки %v", err)
	}
	// отложенное закрытие коннекта
	defer resp.Body.Close()

	// парсинг ответа
	x, err := goquery.Parse(resp.Body)

	// какая дата сейчас
	dateNow := time.Now()
	year, mounth, _ := dateNow.Date()
	var strMounth string
	if int(mounth) < 10 {
		strMounth = "0" + strconv.Itoa(int(mounth))
	} else {
		strMounth = strconv.Itoa(int(mounth))
	}

	// ищу ссылочки на формы
	var urls []string
	log.Println(`forms\/1\d\d-` + strconv.Itoa(year) + strMounth + `01.rar`)
	// совпадения на все формы по текущему месяцу
	regLLink, _ := regexp.Compile(`forms\/1\d\d-` + strconv.Itoa(year) + strMounth + "01.rar")

	for _, i := range x.Find("a").Attrs("href") {
		if regLLink.MatchString(i) {
			urls = append(urls, i)
		}
	}
	fmt.Println(urls)

}

// открытие и декодирование DBF
func readDBF() {
	// NOT WORK
	// csKOI8R windows-1251 windows-1250
	// WORK - DOS 866 - Russian OEM
	dbfTable, err := godbf.NewFromFile("./101-20170701/NAMES.DBF", "866")
	if err != nil {
		log.Fatalf("Ошибка открытия %v", err)
	}

	fmt.Printf("%v\n", len(dbfTable.FieldNames()))

	for i := 0; i <= dbfTable.NumberOfRecords()-1; i++ {
		for y := 0; y <= len(dbfTable.FieldNames())-1; y++ {
			fmt.Printf("%v;", dbfTable.FieldValue(i, y))
		}
		fmt.Println()
	}
}
