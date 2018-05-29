# GoTrigger

Go-Trigger is a scheduling server writen in golang. It can schedule a call to a REST endpoint. It provides APIs to schedule and abort the REST call. It internally uses time package of golang for scheduling.

GoTrigger supports two type of scheduling

1. Absolute Schedule
2. Relative Schedule (relative to current time)

## Absolute Schedule

The below example shows an absolute schedule, It will be wake up on 29/05/2018 18:05:00 and makes http GET call to www.google.com.

```json
http://localhost:9999/api/v1/schedule/absolute
{
	"schedule" : {
		"date" : 29,
		"month": 5,
		"year" : 2018,
		"hour" : 18,
		"minute" : 5,
		"second" : 0
	},
	"restRequest" :{
		"url" : "http://www.google.com",
		"httpmethod" : "GET",
		"httpbody" : "NA"
	}
}
```

## Relative Schedule

The below example shows relative schedule, It will be wake up after 6 seconds of current time and makes http GET call to www.google.com.

```json
POST http://localhost:9999/api/v1/schedule/relative
{
	"schedule" : {
		"hour" : 0,
		"minute" : 0,
		"second" : 6
	},
	"restRequest" :{
		"url" : "http://www.google.com/login",
		"httpmethod" : "POST",
		"httpbody" : "this is the post body"
	}
}
```
### Response
Schedule API gives scheduleID as response in case of success.
```json
{
    "scheduleID": "scheduleID",
    "msg": "Success"
}
```

## Abort Schedule

GoTrigger also supports the Abort operation on already scheduled Rest Calls. To About the Schedule, need to provide scheduleID in API.
```json
DELETE http://localhost:9999/api/v1/schedule/<scheduleID>
```