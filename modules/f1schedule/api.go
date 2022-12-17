package f1schedule

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	__duration_max = time.Duration(500000) * time.Hour
)

// API: https://ergast.com/mrd/methods/schedule/

type DateTime struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type Location struct {
	Lat      string `json:"lat"`
	Lon      string `json:"lon"`
	Locality string `json:"locality"`
	Country  string `json:"country"`
}

type CircuitData struct {
	Id       string `json:"circuitId"`
	Name     string `json:"circuitName"`
	Location Location
}

type RaceData struct {
	DateTime
	Season  string `json:"season"`
	Round   string `json:"round"`
	Name    string `json:"raceName"`
	Circuit CircuitData
	// Date           string `json:"date"`
	// Time           string `json:"time"`
	FirstPractice  DateTime
	SecondPractice DateTime
	ThirdPractice  DateTime
	Qualifying     DateTime
	Sprint         DateTime
}

type RaceTable struct {
	Season string `json:"season"`
	Races  []RaceData
}

type ScheduleData struct {
	Series    string `json:"series"`
	Limit     string `json:"limit"`
	Offset    string `json:"offset"`
	Total     string `json:"total"`
	RaceTable RaceTable
}

type RaceDataTemplateContext struct {
	Race           *RaceData
	IsSprintFormat bool
	IsFP1          bool
	IsFP2          bool
	IsFP3          bool
	IsSprint       bool
	IsQualifying   bool
	IsRace         bool
}

type RaceSchedule struct {
	Data ScheduleData `json:"MRData"`
}

func (d *DateTime) GetTime() (time.Time, error) {
	return time.Parse(time.RFC3339, d.Date+"T"+d.Time)
}

func (d *DateTime) GetFormattedTime() string {
	t, err := d.GetTime()
	if err != nil {
		return err.Error()
	}
	return t.Local().Format("03:04PM 02/01/06")
}

func (t *RaceTable) GetRaceDataForDisplay() (*RaceDataTemplateContext, error) {
	now := time.Now()
	var last, next *RaceData
	d_next, d_last := __duration_max, __duration_max
	var diff time.Duration
	for _, race := range t.Races {
		rt, err := race.GetTime()
		if err != nil {
			return nil, err
		}
		if rt.After(now) {
			diff = rt.Sub(now)
			if diff < d_next {
				d_next = diff
				next = &race
			}
		} else {
			diff = now.Sub(rt)
			if diff < d_last {
				d_last = diff
				last = &race
			}
		}
	}
	ctx := RaceDataTemplateContext{}
	if next == nil && last == nil {
		return nil, errors.New("ERR: RACE TIME")
	}
	if next == nil {
		ctx.Race = last
	} else {
		ctx.Race = next
	}
	ctx.IsSprintFormat = ctx.Race.Sprint.Date != ""
	return &ctx, nil
}

func F1RaceCurrentSchedule() (*RaceSchedule, error) {
	response, err := http.Get("https://ergast.com/api/f1/current.json")
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var schedule RaceSchedule
	err = json.Unmarshal(body, &schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}
