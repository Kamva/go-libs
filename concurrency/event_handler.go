package concurrency

import (
	"sync"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kamva/go-libs/exceptions"
	"github.com/kamva/go-libs/translation"
	"github.com/kamva/go-libs/sentry"
)

type EventListener struct {
	Listener []func(*Event, interface{})
	RollBack func(interface{})
}

type Event struct {
	eventMap  map[string]EventListener
	issueCode string
	waitGroup *sync.WaitGroup
	channel   chan exceptions.RoutineException
}

func (e *Event) Fire(event string, data interface{}) {
	for _, function := range e.eventMap[event].Listener {
		e.waitGroup.Add(1)
		go function(e, data)
	}

	e.waitGroup.Wait()
	close(e.channel)

	var errors []exceptions.RoutineException
	var criticalErrors []exceptions.RoutineException
	for exception := range e.channel {
		errors = append(errors, exception)
		if exception.Critical {
			criticalErrors = append(criticalErrors, exception)
		}
	}

	if len(criticalErrors) > 0 {
		if RollBack := e.eventMap[event].RollBack; RollBack != nil {
			RollBack(data)
		}

		panic(exceptions.AggregatedRoutineException{
			Message:         fmt.Sprintf("Error in event [%s] handlers", event),
			ResponseMessage: translation.Translate("internal_error"),
			Code:            e.issueCode,
			StatusCode:      iris.StatusInternalServerError,
			Errors:          errors,
		})
	}

	if len(errors) > 0 {
		sentry.CaptureRoutineException(errors)
	}
}

func (e *Event) RecoverRoutinePanic(caller string, critical bool) {
	defer e.waitGroup.Done()

	if err := recover(); err != nil {
		if err, ok := err.(exceptions.Exception); ok {
			e.channel <- exceptions.RoutineException{
				Message:         err.Message,
				ResponseMessage: err.ResponseMessage,
				RoutineName:     caller,
				Critical:        critical,
			}
		} else {
			e.channel <- exceptions.RoutineException{
				Message:         fmt.Sprint(err),
				ResponseMessage: fmt.Sprint(err),
				RoutineName:     caller,
				Critical:        critical,
			}
		}
	}
}

func NewEvent(eventMap map[string]EventListener, exceptionCode string) *Event {
	return &Event{
		eventMap:  eventMap,
		issueCode: exceptionCode,
		waitGroup: GetWaitGroup(),
		channel:   make(chan exceptions.RoutineException, 10),
	}
}

func GetWaitGroup() *sync.WaitGroup {
	return &sync.WaitGroup{}
}
