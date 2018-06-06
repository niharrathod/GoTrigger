package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"
)

const contentType = "Content-Type"
const appJSON = "application/json;charset=UTF-8"

// RestClient represents the rest request.
type RestClient struct {
	URL    string `json:"url"`
	Method string `json:"httpmethod"`
	Body   string `json:"httpbody"`
}

// Execute method act as rest client for the rest request
func (r *RestClient) Execute() error {
	req, err := http.NewRequest(r.Method,
		r.URL,
		bytes.NewBufferString(r.Body),
	)

	if err != nil {
		log.Print(err)
		return errors.New("Could not create http.request")
	}

	req.Header.Set(contentType, appJSON)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return errors.New("Could not make http.request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Http Request Failed with Response Status:" + resp.Status)
	}

	return nil
}

// Schedule represents a absolute schedule
type Schedule struct {
	Date   int        `json:"date"`
	Month  time.Month `json:"month"`
	Year   int        `json:"year"`
	Hour   int        `json:"hour"`
	Minute int        `json:"minute"`
	Second int        `json:"second"`
}

// GetTimer .
func (s *Schedule) GetTimer() *time.Timer {
	return time.NewTimer(s.GetDuration())
}

// GetDuration .
func (s *Schedule) GetDuration() time.Duration {
	curTime := time.Now()
	scheduleTime := time.Date(
		s.Year,
		s.Month,
		s.Date,
		s.Hour,
		s.Minute,
		s.Second,
		0, curTime.Location())
	return scheduleTime.Sub(curTime)
}

// RestSchedule .
type RestSchedule struct {
	ID       string
	Schedule *Schedule
	RestReq  *RestClient
}

// GetID .
func (rcs *RestSchedule) GetID() string {
	return rcs.ID
}

// GetDuration .
func (rcs *RestSchedule) GetDuration() time.Duration {
	return rcs.Schedule.GetDuration()
}

// Execute .
func (rcs *RestSchedule) Execute() {
	err := rcs.RestReq.Execute()
	if err != nil {
		log.Println("For ID:", rcs.ID, " ", err.Error())
		return
	}
}

// IsValidSchedule validates the schedule with repect to min and max windows configuration.
func IsValidSchedule(s *Schedule) error {
	minWindowSec := config.Scheduler.MinScheduleWindowSeconds
	maxWindowHour := config.Scheduler.MaxScheduleWindowHours
	maxWindowSec := maxWindowHour * int(time.Hour.Seconds())
	durationSec := int(s.GetDuration().Seconds())

	if durationSec <= minWindowSec {
		return errors.New("Schedule must be minimum " +
			IntToString(minWindowSec) + "secs from current time")
	}

	if durationSec > maxWindowSec {
		return errors.New("Shedule can not be more than " +
			IntToString(maxWindowHour) + "Hours from current time")
	}

	return nil
}

// GetSchedule converts the relative schedule to absolute
func GetSchedule(hour, minute, second int) *Schedule {
	scheduleTime := time.Now()
	scheduleTime = scheduleTime.Add(time.Duration(hour) * time.Hour)
	scheduleTime = scheduleTime.Add(time.Duration(minute) * time.Minute)
	scheduleTime = scheduleTime.Add(time.Duration(second) * time.Second)

	var schedule Schedule
	schedule.Date = scheduleTime.Day()
	schedule.Month = scheduleTime.Month()
	schedule.Year = scheduleTime.Year()
	schedule.Hour = scheduleTime.Hour()
	schedule.Minute = scheduleTime.Minute()
	schedule.Second = scheduleTime.Second()
	return &schedule
}
