package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/xerrors"

	"github.com/PuerkitoBio/goquery"
)

const wikiDomain = "ja.wikipedia.org"

func main() {
	run("https://ja.wikipedia.org/wiki/Category:%E6%9D%B1%E4%BA%AC%E9%83%BD%E3%81%AE%E8%A6%B3%E5%85%89%E5%9C%B0")
	//article, _ := scrapeDetailPage("https://ja.wikipedia.org/wiki/%E6%B5%85%E8%8D%89%E5%85%AC%E5%9C%92%E5%85%AD%E5%8C%BA")
	//if err := CreateIfNotExist(article); err != nil {
	//	fmt.Printf("%+v\n", err)
	//}
}

func run(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%+v\n", xerrors.Errorf("[ERROR] :%w", err))
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("%+v\n", xerrors.Errorf("[ERROR] :%w", err))
		return
	}
	doc.Find("#mw-pages > .mw-content-ltr > .mw-category > .mw-category-group > ul > li > a").Each(func(_ int, s *goquery.Selection) {
		absPath, exists := s.Attr("href")
		if !exists {
			return
		}
		article, err := scrapeDetailPage("https://" + wikiDomain + absPath)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		if err := CreateIfNotExist(article); err != nil {
			fmt.Printf("%+v\n", err)
		}
	})
}

func scrapeDetailPage(url string) (*Article, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, xerrors.Errorf("[ERROR] :%w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("[ERROR] :%w", err)
	}
	title := doc.Find("h1").Text()
	latitude := doc.Find(".latitude").Text()
	longitude := doc.Find(".longitude").Text()

	details := make(map[string]string)
	currentCategory := "リード文"
	doc.Find("#mw-content-text > .mw-parser-output > h2, h3, p").Each(func(_ int, s *goquery.Selection) {
		nodeName := goquery.NodeName(s)
		switch nodeName {
		case "h2":
			currentCategory = strings.Trim(s.Text(), "[編集]")
		case "p", "h3":
			details[currentCategory] += strings.Trim(s.Text(), "[編集]")
		default:
			fmt.Printf("%+v\n", xerrors.New("[ERROR] got unexpected element "+nodeName))
		}
	})

	article := NewArticle(title, url, latitude, longitude, details)
	return article, nil
}
