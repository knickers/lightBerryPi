/*
type event struct {
	pins        []gpio.Pin
	state       gpio.State
	nextTime    time.Time
	repeatDays  []bool
	repeatWeeks []bool
}
*/

[
	{
		"Pins": [0,2,...],
		"State": 1, // or 0
		"NextTime": Date(),
		"RepeatDays": [false,true,true,true,true,true,false],
		"RepeatWeeks": [false,true,true,true,true,false,false,true,...] // 52X
	},
]
