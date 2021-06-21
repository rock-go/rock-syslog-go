package syslog

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
)

const (
	RFC3164 int = iota + 1
	RFC5424
	RFC6587
	Automatic

)

func newLuaSyslogS(L *lua.LState) int {
	cfg := newConfig(L)
	if e := cfg.verify(); e != nil {
		L.RaiseError("%v" , e)
		return 0
	}

	var obj *server
	var ok bool

	proc := L.NewProc(cfg.name)
	if proc.Value == nil {
		proc.Value = newSyslogS(cfg)
		goto done
	}

	obj , ok = proc.Value.(*server)
	if !ok {
		L.RaiseError("want to reset not syslog")
		return 0
	}
	obj.cfg = cfg

done:
	L.Push(proc)
	return 1
}

func LuaInjectApi(env xcall.Env ) {
	uv := lua.NewUserKV()
	uv.Set("RFC3164" , lua.LNumber(RFC3164))
	uv.Set("RFC5424" , lua.LNumber(RFC5424))
	uv.Set("RFC6587" , lua.LNumber(RFC6587))
	uv.Set("AUTO" , lua.LNumber(Automatic))

	uv.Set("server" , lua.NewFunction(newLuaSyslogS))
	env.SetGlobal("syslog", uv)
}