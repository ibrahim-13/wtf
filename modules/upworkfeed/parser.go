package upworkfeed

import (
	"encoding/xml"
	"fmt"
	"html"
	"strings"
	"time"
)

const (
	__parse_pattern_hourly_rate string = "<b>Hourly Range</b>:"
	__parse_pattern_budget      string = "<b>Budget</b>:"
	__parse_pattern_category    string = "<b>Category</b>:"
	__parse_pattern_skills      string = "<b>Skills</b>:"
	__parse_pattern_country     string = "<b>Country</b>:"
)

func ParseXml(src []byte) (*UpworkRss, error) {
	feedData := UpworkRss{}
	err := xml.Unmarshal(src, &feedData)
	if err == nil {
		ct, err := time.Parse(time.RFC1123Z, feedData.Channel.PublishDate)
		if err == nil {
			feedData.Channel.PublishDateTime = ct
			feedData.Channel.PublishDate = parse_format_date_time(ct)
		}
		for i := range feedData.Channel.Items {
			(&feedData.Channel.Items[i]).parseItem()
		}
	}
	return &feedData, err
}

func (item *UpworkItem) parseItem() {
	if len(item.Title) > 8 {
		hasPostfix := strings.TrimSpace(strings.ToLower(item.Title[len(item.Title)-8:])) == "- upwork"
		if hasPostfix {
			item.Title = strings.TrimSpace(item.Title[:len(item.Title)-8])
		}
	}
	ct, err := time.Parse(time.RFC1123Z, item.PublishDate)
	if err == nil {
		item.PublishDateTime = ct
		item.PublishDate = parse_format_date_time(ct)
	}
	parts := strings.Split(item.Description, "<br />")
	dl := 0
	for i, p := range parts {
		dl = i
		if strings.HasPrefix(p, "<b>") {
			break
		}
	}
	if dl > 0 {
		item.ShortDescription = html.UnescapeString(strings.TrimSpace(strings.Join(parts[0:dl], " ")))
	}
	for _, part := range parts {
		section := ""
		if section = parse_extract_section(part, __parse_pattern_hourly_rate); section != "" {
			item.Rate = fmt.Sprintf("Hourly: %s", section)
		} else if section = parse_extract_section(part, __parse_pattern_budget); section != "" {
			item.Rate = fmt.Sprintf("Budget: %s", section)
		} else if section = parse_extract_section(part, __parse_pattern_category); section != "" {
			item.Category = section
		} else if section = parse_extract_section(part, __parse_pattern_skills); section != "" {
			skills := strings.Split(section, ",")
			for i := range skills {
				skills[i] = strings.TrimSpace(skills[i])
			}
			item.SkillsArr = skills
			item.Skills = strings.Join(skills, " ")
		} else if section = parse_extract_section(part, __parse_pattern_country); section != "" {
			item.Country = section
		}
	}
}

func parse_extract_section(target, contain string) string {
	index := strings.Index(target, contain)
	if index > -1 {
		return strings.TrimSpace(string(target[index+len(contain):]))
	}
	return ""
}

func parse_format_date_time(t time.Time) string {
	return t.Local().Format("03:04PM 02/01/06")
}
