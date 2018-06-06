package main

import (
	"GoTrigger/scheduler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type absoluteRequest struct {
	S       Schedule   `json:"schedule"`
	RestReq RestClient `json:"restRequest"`
}

type relativeRequest struct {
	S struct {
		Hour   int `json:"hour"`
		Minute int `json:"minute"`
		Second int `json:"second"`
	} `json:"schedule"`
	RestReq RestClient `json:"restRequest"`
}

// Response : response for absolute/relative schedule request.
type Response struct {
	ScheduleID string `json:"scheduleID"`
	Msg        string `json:"msg"`
}

// AbsoluteSchedule : controller for absolute schedule request.
func AbsoluteSchedule(c *gin.Context) {
	var ar absoluteRequest
	err := c.BindJSON(&ar)
	if err != nil {
		setResponse(c, http.StatusBadRequest, "", err.Error())
		return
	}
	ID := GetNewUniqueID()
	var rcs RestSchedule
	rcs.ID = ID
	rcs.RestReq = &ar.RestReq
	rcs.Schedule = &ar.S
	err = scheduler.Schedule(&rcs)
	if err != nil {
		setResponse(c, http.StatusInternalServerError, "", err.Error())
	} else {
		setResponse(c, http.StatusOK, ID, "Success")
	}
}

// RelativeSchedule : controller for relative schedule request.
func RelativeSchedule(c *gin.Context) {
	var rr relativeRequest
	err := c.BindJSON(&rr)
	if err != nil {
		setResponse(c, http.StatusBadRequest, "", err.Error())
		return
	}

	ID := GetNewUniqueID()
	var rcs RestSchedule
	rcs.ID = ID
	rcs.RestReq = &rr.RestReq
	rcs.Schedule = GetSchedule(rr.S.Hour, rr.S.Minute, rr.S.Second)
	err = scheduler.Schedule(&rcs)
	if err != nil {
		setResponse(c, http.StatusInternalServerError, "", err.Error())
	} else {
		setResponse(c, http.StatusOK, ID, "Success")
	}
}

// AbortSchedule : controller to abort the already scheduled request.
func AbortSchedule(c *gin.Context) {
	ID := c.Param("id")
	if len(ID) == 0 {
		setResponse(c, http.StatusBadRequest, ID, "Invalid schedule ID")
		return
	}
	err := scheduler.Abort(ID)
	if err != nil {
		setResponse(c, http.StatusInternalServerError, ID, err.Error())
	} else {
		setResponse(c, http.StatusOK, ID, "Success")
	}
}

func setResponse(c *gin.Context, statusCode int, ID, msg string) {
	var resp Response
	resp.ScheduleID = ID
	resp.Msg = msg
	c.JSON(statusCode, resp)
}
