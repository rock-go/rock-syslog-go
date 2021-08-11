package syslog

import (
	"github.com/rock-go/rock/lua"
)

func (s *server) info(out lua.Printer) {
	out.Printf("type: %s", s.Type())
	out.Printf("uptime: %s", s.uptime.Format("2006-01-02 15:04:06"))
	out.Println("version:  v1.0.0")
	out.Println("")
}

func (s *server) Show(out lua.Printer) {
	s.info(out)

	out.Printf("name = %s", s.Name())
	out.Printf("protocol = %s", s.cfg.protocol)
	out.Printf("listen = %s", s.cfg.listen)
	out.Printf("format = %d", s.cfg.format)

	n := len(s.cfg.output)
	for i := 0; i < n; i++ {
		out.Printf("output.%d = %s", i+1, s.cfg.output[i].Name())
	}
}

func (s *server) Help(out lua.Printer) {
	s.info(out)
}
