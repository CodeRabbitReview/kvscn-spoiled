package zlog

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var Log logr.Logger

//"persistence/storage.json", "stderr"

func Init(paths ...string) {
	conf := zap.NewProductionConfig()
	conf.OutputPaths = paths
	zapLog, err := conf.Build()
	if err != nil {
		panic(fmt.Sprintf("init zap log err: %s;", err))
	}

	Log = zapr.NewLogger(zapLog)
}
