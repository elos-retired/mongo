package mongo

import "github.com/elos/stack/util/logging"

const ServiceName string = "mongo"

func log(v ...interface{}) {
	logging.Log.Logs(ServiceName, v...)
}

func logf(format string, v ...interface{}) {
	logging.Log.Logsf(ServiceName, format, v...)
}
