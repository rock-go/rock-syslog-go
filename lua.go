package syslog

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
	"reflect"
)

const (
	RFC3164 int = iota + 1
	RFC5424
	RFC6587
	Automatic
)

var (
	SYSLOGS = reflect.TypeOf((*server)(nil)).String()
)

func newLuaSyslogS(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name, SYSLOGS)
	if proc.IsNil() {
		proc.Set(newSyslogS(cfg))
		goto done
	}
	proc.Value.(*server).cfg = cfg

done:
	L.Push(proc)
	return 1
}

func LuaInjectApi(env xcall.Env) {
	uv := lua.NewUserKV()
	uv.Set("RFC3164", lua.LNumber(RFC3164))
	uv.Set("RFC5424", lua.LNumber(RFC5424))
	uv.Set("RFC6587", lua.LNumber(RFC6587))
	uv.Set("AUTO", lua.LNumber(Automatic))

	uv.Set("server", lua.NewFunction(newLuaSyslogS))
	env.SetGlobal("syslog", uv)
}
