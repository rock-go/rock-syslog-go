package syslog

import (
	"errors"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/utils"
)

type config struct {
	name     string
	protocol string
	listen   string
	format   int
	encode   string
	output   []lua.Writer
}



func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{format: Automatic}
	tab.ForEach(func(key lua.LValue, val lua.LValue) {
		switch key.String() {

		case "name":
			cfg.name = val.String()

		case "protocol":
			cfg.protocol = val.String()

		case "listen":
			cfg.listen = val.String()

		case "format":
			cfg.format = lua.CheckInt(L , val)

		case "encode":
			cfg.encode = val.String()

		case "output":
			cfg.output = checkOutputSDK(L , val)

		default:
			L.RaiseError("invalid syslog config field")
			return
		}
	})

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v" , e)
		return nil
	}

	return cfg
}

func (cfg *config) verify() error {
	if e := utils.Name(cfg.name);e != nil {
		return e
	}


	switch cfg.protocol {
	case "tcp","udp", "tcp/udp":
		//ok
	default:
		return errors.New("invalid network protocol")
	}

	switch cfg.format {
	case RFC3164 , RFC5424 , RFC6587 , Automatic:
		//ok
	default:
		return errors.New("invalid syslog format")
	}

	switch cfg.encode {
	case "json","raw":
		//ok

	default:
		return errors.New("invalid syslog encode")
	}



	return nil
}
