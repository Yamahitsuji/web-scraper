package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/xerrors"
)

const outDir = "./default"

func main() {
	var url string
	fmt.Print("Enter URL : ")
	_, err := fmt.Scan(&url)
	if err != nil {
		fmt.Printf("%+v\n", xerrors.Errorf("[ERROR] :%w", err))
		return
	}
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
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		switch {
		case strings.HasPrefix(src, "http"):
			if err := saveImageFromURL(src); err != nil {
				fmt.Printf("%+v\n", err)
			}
		case strings.HasPrefix(src, "/"):
			if err := saveImageFromURL(strings.Join(strings.Split(url, "/")[:3], "/") + src); err != nil {
				fmt.Printf("%+v\n", err)
			}
		case strings.HasPrefix(src, "data:image/jpeg;base64") ||
			strings.HasPrefix(src, "data:image/jpg;base64") ||
			strings.HasPrefix(src, "data:image/png;base64"):
			if err := saveBase64Image(src); err != nil {
				fmt.Printf("%+v\n", err)
			}
		default:
			fmt.Printf("[WARN] Unknown img src is detected. src=%s\n", src)
		}
	})
}

func saveImageFromURL(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	defer res.Body.Close()
	img, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return saveImage(img)
}

func saveBase64Image(base64Img string) error {
	img, err := base64.StdEncoding.DecodeString(base64Img)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return saveImage(img)
}

func saveImage(img []byte) error {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.Mkdir(outDir, 0777); err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}
	randStr, err := makeRandomStr(10)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	file, err := os.Create(strings.TrimRight(outDir, "/") + "/" + randStr + ".png")
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	defer file.Close()
	if _, err = file.Write(img); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func makeRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", xerrors.New("unexpected error...")
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}
