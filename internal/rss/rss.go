package rss

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Rss struct {
	Channel struct {
		Title string    `xml:"title"`
		Item  []RssItem `xml:"item"`
	} `xml:"channel"`
}

type Source struct {
	URL               string   `yaml:"url"`
	Pattern           []string `yaml:"pattern"`
	DownloadDirectory string   `yaml:"downloadDirectory"`
	DownloadPaused    bool     `yaml:"downloadPaused"`
}

type RssItem struct {
	Title     string `xml:"title"`
	MagnetURI string `xml:"magnetURI"`
}

type DownloadItem struct {
	Item   RssItem
	Source Source
}

func (s *Source) FindMatchingItems() []DownloadItem {
	var items []DownloadItem
	for _, item := range s.FetchRssItems() {
		for _, pattern := range s.Pattern {
			if match, _ := regexp.MatchString(pattern, item.Title); match {
				log.Println(item.Title)
				log.Println(pattern)
				log.Println(match)
				items = append(items, DownloadItem{
					Item:   item,
					Source: *s,
				})
			}
		}
	}
	return items
}

func (s *Source) FetchRssItems() []RssItem {
	resp, err := http.Get(s.URL)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error parsing body: %s", err)
		return nil
	}
	rss := decodeContents(body)
	return rss.Channel.Item
}

func decodeContents(contents []byte) Rss {
	r := Rss{}
	xml.Unmarshal(contents, &r)
	return r
}
