package exceptions

type Exception struct {
	Message         string
	ResponseMessage string
	Code            string
	StatusCode      int
}

type AggregatedRoutineException struct {
	Message         string
	ResponseMessage string
	Code            string
	StatusCode      int
	Errors          []RoutineException
}

type RoutineException struct {
	Message         string
	ResponseMessage string
	RoutineName     string
	Critical        bool
}

type ValidationException struct {
	ResponseMessage string
	Data            interface{}
	Code            string
}
