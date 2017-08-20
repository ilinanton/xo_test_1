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

	"t1/functions"
	"t1/structs"
)

var currency string
var value float64
var dValute = make(map[string]*structs.Valute)

func main() {
	appendBaseValute()
	initParams()
	getCbrData()

	_, ok := dValute[currency]

	if !ok {
		log.Fatalln("Данные по валюте", currency, "отсутствуют")
	}

	fmt.Println("Входные параметры: currency -", currency, dValute[currency].Name, " value -", value)

	var floatValue float64
	var floatNominal float64
	var err error

	floatValue, err = functions.StringToFloat64(dValute[currency].Value)
	functions.ChErr(err)
	rub := value * floatValue

	floatValue, err = functions.StringToFloat64(dValute["EUR"].Value)
	functions.ChErr(err)

	floatNominal, err = functions.StringToFloat64(dValute["EUR"].Nominal)
	functions.ChErr(err)

	eur := rub / floatValue * floatNominal

	floatValue, err = functions.StringToFloat64(dValute["INR"].Value)
	functions.ChErr(err)

	floatNominal, err = functions.StringToFloat64(dValute["INR"].Nominal)
	functions.ChErr(err)
	inr := rub / floatValue * floatNominal

	fmt.Println("RUB:", functions.Float64ToString(rub), "/EUR:", functions.Float64ToString(eur), "/INR:", functions.Float64ToString(inr))
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
		log.Fatalln("Ошибка количества параметров!")
	}

	flag.StringVar(&currency, "currency", "", "Валюта")
	var tempValue string
	flag.StringVar(&tempValue, "value", "0", "Значение")
	flag.Parse()

	currency = strings.ToUpper(currency)

	var err error
	value, err = functions.StringToFloat64(tempValue)

	if err != nil {
		log.Fatalln("Ошибка конвертации количества!")
	}
}

func getCbrData() {
	url := "http://www.cbr.ru/scripts/XML_daily.asp"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Status error: %v", resp.StatusCode)
	}

	var query structs.Query

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&query)
	if err != nil {
		log.Fatalf("Read body: %v", err)
	}

	for i, valute := range query.ValuteList {
		dValute[valute.CharCode] = &query.ValuteList[i]
	}
}