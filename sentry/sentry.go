package sentry

import (
	"fmt"
	"errors"
	"runtime"
	"strconv"
	"github.com/kataras/iris/context"
	"github.com/getsentry/raven-go"
	"github.com/kataras/iris"
	"kamva.ir/libraries/exceptions"
)

func New() context.Handler {
	return func(context context.Context) {
		writer, request := context.ResponseWriter(), context.Request()

		defer func() {
			if err := recover(); err != nil {
				var errString string
				var aggregateRoutineException exceptions.AggregatedRoutineException
				var shouldReport bool

				if validation, ok := err.(exceptions.ValidationException); ok {
					shouldReport = false
					context.StatusCode(iris.StatusNotAcceptable)
					context.JSON(iris.Map{
						"code":    validation.Code,
						"message": validation.ResponseMessage,
						"data":    validation.Data,
					})
					return
				} else if exception, ok := err.(exceptions.Exception); ok {
					errString = exception.Message
					context.Values().Set("message", exception.ResponseMessage)
					context.Values().Set("code", exception.Code)
					context.StatusCode(exception.StatusCode)
					shouldReport = exception.StatusCode >= 500
				} else if routineException, ok := err.(exceptions.AggregatedRoutineException); ok {
					errString = routineException.Message
					aggregateRoutineException = routineException
					context.Values().Set("message", routineException.ResponseMessage)
					context.Values().Set("code", routineException.Code)
					context.Values().Set("data", getCriticalErrorArray(routineException.Errors))
					context.StatusCode(routineException.StatusCode)
					shouldReport = true
				} else {
					errString = fmt.Sprint(err)
					writer.WriteHeader(iris.StatusInternalServerError)
					shouldReport = true
				}

				if shouldReport {
					packet := raven.NewPacket(
						errString,
						raven.NewException(errors.New(errString), raven.NewStacktrace(2, 3, nil)),
						raven.NewHttp(request),
					)

					raven.Capture(packet, getCaptureTags(aggregateRoutineException.Errors))
				}

				logWarning(context, errString)
			}
		}()

		context.Next()
	}
}

func getCaptureTags(exceptions []exceptions.RoutineException) map[string]string {
	var tag = make(map[string]string)

	for _, value := range exceptions {
		tag[value.RoutineName] = value.Message
	}

	return tag
}

func getCriticalErrorArray(exceptions []exceptions.RoutineException) map[string]string {
	var errorArray = make(map[string]string)

	for _, value := range exceptions {
		errorArray[value.RoutineName] = value.ResponseMessage
	}

	return errorArray
}

func logWarning(context iris.Context, err string) {
	logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", context.HandlerName())
	logMessage += fmt.Sprintf("At Request: %s\n", getRequestLogs(context))
	logMessage += fmt.Sprintf("Trace: %s\n", err)
	logMessage += fmt.Sprintf("\n%s", getStacktrace())
	context.Application().Logger().Warn(logMessage)
}

func getStacktrace() string {
	var stacktrace string
	for i := 1; ; i++ {
		_, f, l, got := runtime.Caller(i)
		if !got {
			break

		}

		stacktrace += fmt.Sprintf("%s:%d\n", f, l)
	}

	return stacktrace
}

func getRequestLogs(context context.Context) string {
	var status, ip, method, path string
	status = strconv.Itoa(context.GetStatusCode())
	path = context.Path()
	method = context.Method()
	ip = context.RemoteAddr()
	// the date should be logged by iris' Logger, so we skip them
	return fmt.Sprintf("%v %s %s %s", status, path, method, ip)
}

func CaptureRoutineException(exception []exceptions.RoutineException) {
	errString := "routine exception"
	packet := raven.NewPacket(
		errString,
		raven.NewException(errors.New(errString), raven.NewStacktrace(2, 3, nil)),
	)

	raven.Capture(packet, getCaptureTags(exception))
}
