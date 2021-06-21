package syslog

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"github.com/rock-go/rock/lua"
	"gopkg.in/mcuadros/go-syslog.v2"
	"github.com/rock-go/rock/utils"
	"github.com/rock-go/rock/logger"
)

type server struct {
	lua.Super

	uptime time.Time
	state  lua.LightUserDataStatus

	cfg  *config
	serv  *syslog.Server
}

func newSyslogS(cfg *config) *server {
	return &server{cfg:cfg , state:lua.INIT}
}

func (s *server) Name() string {
	return s.cfg.name
}

func (s *server) Type() string {
	return "syslog.server"
}

func (s *server) State() lua.LightUserDataStatus {
	return s.state
}

func (s *server) Start() error {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	serv := syslog.NewServer()
	switch s.cfg.format {
	case RFC3164:
		serv.SetFormat(syslog.RFC3164)
	case RFC6587:
		serv.SetFormat(syslog.RFC6587)
	case RFC5424:
		serv.SetFormat(syslog.RFC5424)
	default:
		serv.SetFormat(syslog.Automatic)

	}

	serv.SetHandler(handler)

	var err error
	switch s.cfg.protocol {
	case "tcp":
		err = serv.ListenTCP(s.cfg.listen)
	case "udp":
		err = serv.ListenUDP(s.cfg.listen)
	case "tcp/udp":
		err = serv.ListenUDP(s.cfg.listen)
		err = serv.ListenTCP(s.cfg.listen)
	default:
		err = errors.New("invalid protocol , must be tcp , udp , tcp/udp; got " + s.cfg.protocol)
	}

	if err != nil {
		return err
	}

	serv.Boot()
	go func(channel syslog.LogPartsChannel){
		for item := range channel {
			switch s.cfg.encode {
			case "json":
				if v, e := json.Marshal( item ); e == nil {
					s.Push( v )
				} else {
					logger.Errorf("syslog-go err: %v" , e)
				}
			case "line":
				s.Push( fmt.Sprintf("%v" , item ))
			}
		}
	}(channel)

	s.serv = serv
	s.state = lua.RUNNING
	s.uptime = time.Now()
	return nil
}

func (s *server) Push( v interface{} ) {
	n := len(s.cfg.output)
	for i := 0; i< n ;i++ {
		w := s.cfg.output[i]
		_ , err := utils.Push(w , v)
		if err != nil {
			logger.Errorf("%s Push io fail , err: %v" , s.Name() , err)
			continue
		}
	}
}

func (s *server) Close() error {
	logger.Errorf("%s stop succeed" , s.Name())
	return s.serv.Kill()
}