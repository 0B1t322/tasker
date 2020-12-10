package taskerserver

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	pb "github.com/0B1t322/tasker/tasker"
)

// Types
const (
	ALL 			= iota 	// 0
	DONE					// 1
	PERIOD					// 2
)
// Patterns for types
const (
	ALLPATTERN 		= "^all:.*"
	DONEPATTERN 	= "^done:.*"
	PERIODPATTERN 	= "^period:.*"
)


/*
Filter has next filters:
			all:
			done: false
			done: true
			period: [Months] [day] - [Moths] [day]
*/
type Filter struct {
	Data 	map[string]string
	Type   	int
}

// NewFilter make a new Filter
func NewFilter(str string) *Filter {
	filter := &Filter{}
	filter.Data = make(map[string]string)

	filter.Type = checkType(str)
	switch filter.Type {
	case ALL:
		// nothing
	case DONE:
		filter.Data["done"] 	= strings.Replace(str, "done: ", "", 1)
	case PERIOD:
		period := strings.Split(strings.Replace(str, "period: ", "", 1), " - ")
		filter.Data["from"] 	= period[0]
		filter.Data["until"]	= period[1]
	}

	return filter
}

// FilterTasks return a slice of filtered tasks
func (f Filter) FilterTasks(tasks []*pb.Task) []*pb.Task {
	newTasks := []*pb.Task{}
	switch f.Type {
	case ALL:
		newTasks = tasks
	case DONE:
		done, err := strconv.ParseBool( f.Data["done"] )
		if err != nil {
			panic(err)
		}
		for _, task := range tasks {
			if task.Done == done {
				newTasks = append(newTasks, task)
			}
		}
	case PERIOD:
		from 			:= f.Data["from"]
		until 			:= f.Data["until"] 
		timeFrom,  err 	:= time.Parse(time.Stamp , from)
		if err != nil {
			panic(err)
		}
		timeUntil, err 	:= time.Parse(time.Stamp, until)
		if err != nil {
			panic(err)
		}

		for _, task := range tasks {
			curTime, err := time.Parse(time.Stamp ,task.CreatesTime)
			if err != nil {
				panic(err)
			}
			// так как After учитывает не включая то мы выравним это значит чтобы он включал
			if curTime.After(timeFrom.Add(-time.Minute)) && curTime.Before(timeUntil.Add(time.Minute)) {
				newTasks = append(newTasks, task)
			}
		}
	}
	
	return newTasks
}

func checkType(str string) int {
	var Type = make(chan int, 1)

	go func(str string) {
		match, err := regexp.MatchString(ALLPATTERN, str)
		if err != nil {
			panic(err)
		}
		
		if match {
			Type <- ALL
		}
	}(str)

	go func(str string) {
		match, err := regexp.MatchString(DONEPATTERN, str)
		if err != nil {
			panic(err)
		}

		if match {
			Type <- DONE
		}
	}(str)

	go func(str string) {
		match, err := regexp.MatchString(PERIODPATTERN, str)
		if err != nil {
			panic(err)
		}

		if match {
			Type <- PERIOD
		}
	}(str)

	return <- Type
}