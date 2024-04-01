package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/html"
)

const (
	Url       = "https://pro.imdb.com/title/%s/"
	UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1"
)

func main() {
	hc := newClient()

	if len(os.Args) != 2 {
		check(errors.New("need imdb ID first arg"))
	}
	id := os.Args[1]

	doc, err := parseIMDB(hc, id)
	check(err)

	img, err := findImage(doc)
	check(err)

	fmt.Println(img)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func parseIMDB(hc http.Client, id string) (*html.Node, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(IMDBURL, id), nil)
	check(err)

	req.Header.Set("accept", "text/html")
	req.Header.Set("user-agent", UserAgent)
	resp, err := hc.Do(req)
	check(err)
	defer resp.Body.Close()

	return html.Parse(resp.Body)
}

func findImage(doc *html.Node) (string, error) {
	var src string
	var found bool
	var crawl func(n *html.Node)
	crawl = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "primary_image" {
					found = true
				} else if a.Key == "src" {
					src = a.Val
				}
			}
			if found {
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawl(c)
		}
	}
	crawl(doc)
	if found {
		return src, nil
	}
	return "", errors.New("could not find image")
}

func newClient() http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 20
	t.MaxConnsPerHost = 20
	t.MaxIdleConnsPerHost = 20

	return http.Client{
		Timeout:   5 * time.Second,
		Transport: t,
	}
}
