package upworkfeed

import "time"

type UpworkItem struct {
	Title            string `xml:"title"`
	Link             string `xml:"link"`
	Description      string `xml:"description"`
	ShortDescription string
	PublishDate      string `xml:"pubDate"`
	Rate             string //Hourly/Budget
	Category         string
	Skills           string
	SkillsArr        []string
	Country          string
	PublishDateTime  time.Time
}

type UpworkChannel struct {
	Title           string `xml:"title"`
	PublishDate     string `xml:"pubDate"`
	PublishDateTime time.Time
	Items           []UpworkItem `xml:"item"`
}

type UpworkRss struct {
	Channel UpworkChannel `xml:"channel"`
}

// Sort items
func (c *UpworkChannel) Len() int {
	return len(c.Items)
}

func (c *UpworkChannel) Less(i, j int) bool {
	return c.Items[i].PublishDateTime.After(c.Items[j].PublishDateTime)
}

func (c *UpworkChannel) Swap(i, j int) {
	c.Items[i], c.Items[j] = c.Items[j], c.Items[i]
}
