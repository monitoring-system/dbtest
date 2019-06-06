package util

import (
	"github.com/prometheus/common/log"
	"plugin"
)

func GetPluginSymbol(pluginName string) plugin.Symbol {
	p, err := plugin.Open(pluginName)
	if err != nil {
		log.Warn("Open a plugin failed", err)
		return nil
	}

	f, err := p.Lookup("Filter")
	if err != nil {
		log.Warn("Lookup the Filter function failed", err)
		return nil
	}

	return f
}
