package main

import (
	"crypto/md5"
	"fmt"

	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	Id      int16 `gorm:"primary_key"`
	Name    string
	Url     string
	Url_img string
	Price   string
}

func main() {
	dsn := "badchaos:pe0038900@tcp(127.0.0.1:3306)/cenozavr?charset=utf8&parseTime=True&loc=Local"
	url := "https://www.ozon.ru/category/moloko-9283/"
	browser := rod.New().Timeout(time.Minute).MustConnect()
	defer browser.MustClose()

	fmt.Printf("js: %x\n\n", md5.Sum([]byte(stealth.JS)))

	page := stealth.MustPage(browser)

	page.MustNavigate(url)
	fmt.Println(page)

	dvs, err := page.ElementX("//div")
	if err != nil {
		panic(err)
	}

	res, err := dvs.Text()
	if err != nil {
		panic(err)
	}

	check_table(dsn)

	first_split := strings.Split(res, "Популярные")
	second := strings.Split(first_split[1], "Дальше")
	last := strings.Split(second[0], "Доставит Ozon, продавец Ozon")
	to_clean := last[:4]

	for _, n := range to_clean {

		res_row := clean_data(strings.Split(n, "\n"))
		fmt.Println(res_row)
		insert_row(dsn, res_row)
	}
}

func insert_row(dsn string, data Product) {

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Select("name", "price").Create(&data)
	fmt.Println("Row has been created")

}

func check_table(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if db.Migrator().HasTable(&Product{}) {
		fmt.Println("Table already exists in database")
	} else {
		fmt.Println("Creating new table")
		db.AutoMigrate(&Product{})
	}
}

func clean_data(data []string) Product {
	result := Product{}
	if len(data[2]) > 10 {
		result.Price = data[3]
	} else if len(data[3]) > 10 {
		result.Price = data[4]
	}

	if strings.Contains(data[6], "Молоко") || strings.Contains(data[6], "молоко") {
		result.Name = data[6]
	} else if strings.Contains(data[7], "Молоко") || strings.Contains(data[7], "молоко") {
		result.Name = data[7]
	} else if strings.Contains(data[5], "Молоко") || strings.Contains(data[5], "молоко") {
		result.Name = data[5]
	} else if strings.Contains(data[4], "Молоко") || strings.Contains(data[4], "молоко") {
		result.Name = data[4]
	} else if strings.Contains(data[8], "Молоко") || strings.Contains(data[8], "молоко") {
		result.Name = data[8]
	}
	return result
}
