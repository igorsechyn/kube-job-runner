package reporter

// Sink interface for different reporter sinks (stdout, datadog)
type Sink interface {
	Info(message string, fields map[string]interface{})
	Error(err error, fields map[string]interface{})
}

func New(sinks ...Sink) *Reporter {
	return &Reporter{
		sinks: sinks,
	}
}

type Reporter struct {
	sinks []Sink
}

func (reporter *Reporter) Info(code string, message string, fields map[string]interface{}) {
	fields["code"] = code
	for _, sink := range reporter.sinks {
		sink.Info(message, fields)
	}
}

func (reporter *Reporter) Error(code string, err error, fields map[string]interface{}) {
	fields["code"] = code
	for _, sink := range reporter.sinks {
		sink.Error(err, fields)
	}
}

//func (reporter *Reporter) ErrorNew(err error, fields map[string]interface{}) {
//	appError, ok := errors.Cause(err).(AppError)
//	if !ok {
//		reporter.Error("unknown.error", appError, fields)
//	} else {
//		reporter.Error(appError.Code, err, fields)
//	}
//}
