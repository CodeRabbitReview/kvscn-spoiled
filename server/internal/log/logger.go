package zlog

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var Log logr.Logger

// Init inits Log instance
// takes paths where logs will be sent
func Init(paths ...string) {
	conf := zap.NewProductionConfig()
	conf.OutputPaths = paths
	zapLog, err := conf.Build()
	if err != nil {
		panic(fmt.Sprintf("init zap log err: %s;", err))
	}

	Log = zapr.NewLogger(zapLog)
}
