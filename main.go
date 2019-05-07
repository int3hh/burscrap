package main // import "github.com/int3hh/burscrap"

import (
	"fmt"
	"html"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/common-nighthawk/go-figure"
	"github.com/gocolly/colly"
)

var date *string
var dirpath string

func init() {
	figure := figure.NewFigure("BurScrap", "", true)
	figure.Print()
	parser := argparse.NewParser("BurScrap", "Scraps Bursa news paper for articles")
	date = parser.String("d", "date", &argparse.Options{Required: true, Help: "Scrap paper from day"})
	err := parser.Parse(os.Args)
	if err != nil {
		log.Panic("Invalid arguments")
	}
	_, err = time.Parse("2006-01-02", *date)
	if err != nil {
		log.Panic("Bad date specified")
	}

	dirpath = path.Join("extracted", *date)
	os.MkdirAll(dirpath, os.ModePerm)

}

func isNumeric(val string) bool {
	_, err := strconv.ParseFloat(val, 32)
	return err == nil
}

func main() {
	fmt.Printf("Collecting data from %s \n", "http://www.bursa.ro/ziar/"+*date)
	collect := colly.NewCollector(colly.AllowedDomains("www.bursa.ro"))
	collect.OnRequest(func(r *colly.Request) {
		fmt.Printf("Browsing to %s \n", r.URL)
	})
	collect.OnHTML("header.caseta-medie-header a[href]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Request.URL.String(), "ziar/") {
			link := e.Attr("href")
			links := strings.Split(link, "-")
			if len(links) > 0 && isNumeric(links[len(links)-1]) {
				collect.Visit(e.Request.AbsoluteURL(link))
			}
		}
	})
	collect.OnHTML("#articol-text", func(e *colly.HTMLElement) {
		title := e.Request.URL.RequestURI()
		title = title[1:strings.LastIndex(title, "-")] + ".txt"
		fil, err := os.Create(path.Join(dirpath, title))
		if err == nil {
			defer fil.Close()
			fmt.Printf("Writing article : %s \n", title)
			e.ForEach("p.par", func(_ int, elem *colly.HTMLElement) {
				text := html.UnescapeString(elem.Text) + "\n"
				fil.Write([]byte(text))
			})
		} else {
			fmt.Printf("Could not write file %s : %s \n", title, err.Error())
		}
	})
	collect.Visit("http://www.bursa.ro/ziar/" + *date)
}
