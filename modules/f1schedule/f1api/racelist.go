package f1api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	__url_race_list string = "https://www.formula1.com/en/racing/%s.html"
)

const ()

type RaceLocation struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Race struct {
	Url           string       `json:"@id"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	StartDate     string       `json:"startDate"`
	StartDateTime time.Time    `json:"-"`
	EndDate       string       `json:"endDate"`
	EndDateTime   time.Time    `json:"-"`
	Location      RaceLocation `json:"location"`
}

type RaceEvent struct {
	Url           string    `json:"@id"`
	Name          string    `json:"name"`
	StartDate     string    `json:"startDate"`
	StartDateTime time.Time `json:"-"`
	EndDate       string    `json:"endDate"`
	EndDateTime   time.Time `json:"-"`
}

type RaceEventData struct {
	Race
	SubEvents []RaceEvent `json:"subEvent"`
}

func getLinkedData[K interface{}](ctx context.Context, url string) ([]K, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(response.Body)
	if err != nil {
		return nil, err
	}
	var races []K
	stack := &NodeStack{}
	stack.Push(doc)
	isHeadFound := false
	for {
		if stack.IsEmpty() {
			break
		}
		current, _ := stack.Pop()
		if isLinkedDataTag(current) {
			var race K
			err = json.Unmarshal([]byte(current.FirstChild.Data), &race)
			if err != nil {
				panic(err)
			}
			races = append(races, race)
		}
		if !isHeadFound {
			for c := current.FirstChild; c != nil; c = c.NextSibling {
				stack.Push(c)
			}
		}
		isHeadFound = isHeadTag(current)
	}
	return races, nil
}

func isLinkedDataTag(node *html.Node) bool {
	if node != nil && node.Type == html.ElementNode && node.Data == "script" {
		for i := range node.Attr {
			if node.Attr[i].Key == "type" &&
				node.Attr[i].Val == "application/ld+json" &&
				node.FirstChild != nil &&
				node.FirstChild.Data != "" {
				return true
			}
		}
	}
	return false
}

func isHeadTag(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "head"
}

func (r *Race) UpdateNameAndDescription() {
	p1 := strings.Split(r.Url, "/")
	page := p1[len(p1)-1]
	np := strings.Split(page, ".")
	name := np[0]
	var sb strings.Builder
	var last byte = 0
	for i := range name {
		if name[i] == '_' {
			sb.WriteRune(' ')
			last = ' '
		} else if name[i] >= 'A' && name[i] <= 'Z' && i > 0 && last != ' ' {
			sb.WriteRune(' ')
			sb.WriteByte(name[i])
			last = name[i]
		} else {
			sb.WriteByte(name[i])
			last = name[i]
		}
	}
	name = sb.String()
	r.Name = name
	r.Description = name + " Grand Prix"
}
