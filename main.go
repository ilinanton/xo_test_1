package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html/charset"

	"github.com/ilinanton/xo_test_1/functions"
	"github.com/ilinanton/xo_test_1/structs"
)

var currency string
var value float64
var dValute = make(map[string]*structs.Valute)

func main() {
	defer func() {
		log.Print("test")
	}()

	appendBaseValute()
	initParams()
	getCbrData()

	_, ok := dValute[currency]

	if !ok {
		log.Println("Данные по валюте", currency, "отсутствуют")
		os.Exit(1)
	}

	fmt.Println("Входные параметры: currency -", currency, dValute[currency].Name, " value -", functions.Float64ToString(value))

	var floatValue float64
	var err error

	floatValue, err = functions.StringToFloat64(dValute[currency].Value)
	functions.ChErr(err)
	rub := value * floatValue

	eur := getValue(dValute["EUR"], rub)
	inr := getValue(dValute["INR"], rub)

	fmt.Println("RUB:", functions.Float64ToString(rub), "/EUR:", functions.Float64ToString(eur), "/INR:", functions.Float64ToString(inr))
}

func getValue(valute *structs.Valute, countRub float64) float64 {
	var floatValue float64
	var floatNominal float64
	var err error

	floatValue, err = functions.StringToFloat64(valute.Value)
	functions.ChErr(err)

	floatNominal, err = functions.StringToFloat64(valute.Nominal)
	functions.ChErr(err)

	return countRub / floatValue * floatNominal
}

func appendBaseValute() {
	var rubValute structs.Valute
	rubValute.CharCode = "RUB"
	rubValute.Value = "1.0"
	rubValute.Name = "Русские рубли"
	rubValute.Nominal = "1"

	dValute[rubValute.CharCode] = &rubValute
}

func initParams() {
	if len(os.Args) != 3 {
		log.Printf("Ошибка количества параметров!")
		os.Exit(1)
	}

	flag.StringVar(&currency, "currency", "", "Валюта")
	var tempValue string
	flag.StringVar(&tempValue, "value", "0", "Значение")
	flag.Parse()

	currency = strings.ToUpper(currency)

	var err error
	value, err = functions.StringToFloat64(tempValue)

	if err != nil {
		log.Printf("Ошибка конвертации количества!")
		os.Exit(1)
	}
}

func getCbrData() {
	url := "http://www.cbr.ru/scripts/XML_daily.asp"

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("GET error: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status error: %v", resp.StatusCode)
		os.Exit(1)
	}

	var query structs.Query

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&query)
	if err != nil {
		log.Printf("Read body: %v", err)
		os.Exit(1)
	}

	for i, _ := range query.ValuteList {
		dValute[query.ValuteList[i].CharCode] = &query.ValuteList[i]
	}
}
