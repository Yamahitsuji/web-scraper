package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/xerrors"

	"github.com/PuerkitoBio/goquery"
)

const wikiDomain = "ja.wikipedia.org"

func main() {
	run("https://ja.wikipedia.org/wiki/Category:%E6%9D%B1%E4%BA%AC%E9%83%BD%E3%81%AE%E8%A6%B3%E5%85%89%E5%9C%B0")
	//article, err := scrapeDetailPage("https://ja.wikipedia.org/wiki/%E3%81%8A%E5%8F%B0%E5%A0%B4")
	//if err != nil {
	//	fmt.Printf("%+v\n", err)
	//	return
	//}
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
	latitudeStr := ""
	doc.Find(".latitude").Each(func(_ int, s *goquery.Selection) {
		if len(s.Text()) > len(latitudeStr) {
			latitudeStr = s.Text()
		}
	})
	latitude, err := cnvLatitudeToFloat(latitudeStr)
	if err != nil {
		return nil, xerrors.Errorf("[ERROR] :%w", err)
	}
	longitudeStr := ""
	doc.Find(".longitude").Each(func(_ int, s *goquery.Selection) {
		if len(s.Text()) > len(longitudeStr) {
			longitudeStr = s.Text()
		}
	})
	longitude, err := cnvLongitudeToFloat(longitudeStr)
	if err != nil {
		return nil, xerrors.Errorf("[ERROR] :%w", err)
	}

	text := ""
	details := make(map[string]string)
	currentCategory := "リード文"
	doc.Find("#mw-content-text > .mw-parser-output > p, h2, h3").Each(func(_ int, s *goquery.Selection) {
		if strings.HasPrefix(s.Text(), "座標:") {
			return
		}

		nodeName := goquery.NodeName(s)
		switch nodeName {
		case "h2":
			currentCategory = strings.Trim(s.Text(), "[編集]")
		case "p", "h3":
			details[currentCategory] += strings.Trim(s.Text(), "[編集]")
		default:
			fmt.Printf("%+v\n", xerrors.New("[ERROR] got unexpected element "+nodeName))
		}
		text += strings.Trim(s.Text(), "[編集]")
	})

	article := NewArticle(title, url, details["リード文"], text, latitude, longitude, details)
	return article, nil
}

// 北緯35度42分17.66秒
func cnvLatitudeToFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	var sign float64
	switch {
	case strings.HasPrefix(s, "北緯"):
		sign = 1
	case strings.HasPrefix(s, "南緯"):
		sign = -1
	default:
		return 0, nil
	}
	ss := strings.SplitAfterN(s, "緯", 2)
	ss = strings.SplitAfterN(ss[1], "度", 2)
	deg, err := strconv.ParseFloat(strings.TrimRight(ss[0], "度"), 64)
	if err != nil {
		return 0, xerrors.Errorf("[ERROR] :%w", err)
	}
	ss = strings.SplitAfterN(ss[1], "分", 2)
	min, err := strconv.ParseFloat(strings.TrimRight(ss[0], "分"), 64)
	if err != nil {
		return 0, xerrors.Errorf("[ERROR] :%w", err)
	}
	sec, err := strconv.ParseFloat(strings.TrimRight(ss[1], "秒"), 64)
	if err != nil {
		return 0, xerrors.Errorf("[ERROR] :%w", err)
	}
	return sign * (deg + min/60 + sec/3600), nil
}

func cnvLongitudeToFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	var sign float64
	switch {
	case strings.HasPrefix(s, "東経"):
		sign = 1
	case strings.HasPrefix(s, "西経"):
		sign = -1
	default:
		return 0, nil
	}
	ss := strings.SplitAfterN(s, "経", 2)
	ss = strings.SplitAfterN(ss[1], "度", 2)
	deg, err := strconv.ParseFloat(strings.TrimRight(ss[0], "度"), 64)
	if err != nil {
		return 0, xerrors.Errorf("[ERROR] :%w", err)
	}
	ss = strings.SplitAfterN(ss[1], "分", 2)
	min, err := strconv.ParseFloat(strings.TrimRight(ss[0], "分"), 64)
	if err != nil {
		return 0, xerrors.Errorf("[ERROR] :%w", err)
	}
	sec, err := strconv.ParseFloat(strings.TrimRight(ss[1], "秒"), 64)
	if err != nil {
		return 0, xerrors.Errorf("[ERROR] :%w", err)
	}
	return sign * (deg + min/60 + sec/3600), nil
}
