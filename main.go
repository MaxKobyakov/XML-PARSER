package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Items struct {
	XMLName xml.Name `xml:"CATALOG"`
	Items   []Item   `xml:"ITEMS>ITEM"`
}

type Item struct {
	XMLName xml.Name `xml:"ITEM"`
	Name    string   `xml:"NAME"`
	Article string   `xml:"ARTICLE"`
	Weight  string   `xml:"WEIGHT"`
	Props   Props    `xml:"PROPS"`
	Photo []string `xml:"PHOTOS>PHOTO"`
}

type Props struct {
	XMLName xml.Name `xml:"PROPS"`
	NameEN  string   `xml:"NAME_EN"`
	Brand   string   `xml:"TORGOVAYA_MARKA"`
}

func DownloadFile(filepath string, url string) error {

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Сливаю файл, подожди плиз....")
	fileUrl := "http://portal.autofamily.ru/export/catalog_export.xml"
	err := DownloadFile("catalog_export.xml", fileUrl)
	if err != nil {
		panic(err)
	}

	fmt.Println("Файл catalog_export.xml успешно загружен")

	xmlfile, err := os.Open("catalog_export.xml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Файл catalog_export.xml открыт успешно")
	defer xmlfile.Close()

	byteValue, _ := ioutil.ReadAll(xmlfile)

	var items Items

	xml.Unmarshal(byteValue, &items)

	file, err := os.Create("art.csv")
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	for i := 0; i < len(items.Items); i++ {
		fmt.Fprintf(file, "%s\t%s\t%s\t%s\t%s\t%s\t\n",items.Items[i].Props.Brand, items.Items[i].Article, items.Items[i].Name, items.Items[i].Props.NameEN, items.Items[i].Weight,items.Items[i].Photo )

	}

	fmt.Println("Файл catalog_export.xml распарсен в art.csv")
}