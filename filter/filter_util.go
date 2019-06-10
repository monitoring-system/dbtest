package filter

import (
	"github.com/pingcap/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"plugin"
	"reflect"
)

func AddFilter(pluginName string) error {
	p, err := plugin.Open(pluginName)
	if err != nil {
		log.Warn("Open a plugin failed", zap.Error(err))
		return err
	}

	f, err := p.Lookup("Filter")
	if err != nil || f == nil {
		log.Warn("Lookup the Filter function failed", zap.Error(err))
		return err
	}

	if reflect.TypeOf(f).ConvertibleTo(errMsgType) {
		RegisterErrMsgFilter((reflect.ValueOf(f).Convert(errMsgType).Interface()).(ErrMsgFilter))
		log.Info("Register filter success", zap.String("plugin", pluginName))
		return nil
	}

	if reflect.TypeOf(f).ConvertibleTo(diffType) {
		RegisterDiffFilter((reflect.ValueOf(f).Convert(diffType).Interface()).(DiffFilter))
		log.Info("Register filter success", zap.String("plugin", pluginName))
		return nil
	}

	log.Error("Unkown filter type")
	return errors.New("Unkown filter type")
}