package syslog

import (
	"github.com/rock-go/rock/lua"
)

func checkOutputSDK(L *lua.LState, v lua.LValue) []lua.Writer {
	sdk := make([]lua.Writer, 0)
	if v.Type() != lua.LTTable {
		L.RaiseError("invalid writer , must be table , got %s", v.Type().String())
		return nil
	}

	v.(*lua.LTable).ForEach(func(key lua.LValue, val lua.LValue) {
		if key.Type() != lua.LTNumber {
			L.RaiseError("invalid writer table , got arr")
			return
		}

		w := lua.CheckWriter(val.(*lua.LightUserData))
		if w == nil {
			L.RaiseError("invalid Push userdata")
			return
		}
		sdk = append(sdk, w)
	})

	return sdk

}
