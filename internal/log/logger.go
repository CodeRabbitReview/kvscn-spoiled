package zlog

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var Log logr.Logger

//nolint
func init() {
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{
		"storage.json", "stderr",
	}
	zapLog, err := conf.Build()
	if err != nil {
		panic(fmt.Sprintf("init zap log err: %s;", err))
	}

	Log = zapr.NewLogger(zapLog)
}
