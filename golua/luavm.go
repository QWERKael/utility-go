package golua

import (
	"github.com/QWERKael/utility-go/log"
	lua "github.com/yuin/gopher-lua"
)

type LuaVM struct {
	L *lua.LState
}

func NewLuaVM() *LuaVM {
	return &LuaVM{
		L: lua.NewState(),
	}
}

func (lvm *LuaVM) Destruct() {
	lvm.L.Close()
}

func (lvm *LuaVM) Exec(command string) error {
	err := lvm.L.DoString(command)
	if err != nil {
		return err
	}
	return nil
}

func (lvm *LuaVM) Run(scriptPath string) error {
	err := lvm.L.DoFile(scriptPath)
	if err != nil {
		return err
	}
	return nil
}

func (lvm *LuaVM) ExecuteLuaScriptWithArgsResult(scriptPath string, luaFuncName string, args ...lua.LValue) (lua.LValue, error) {
	//执行脚本
	if err := lvm.L.DoFile(scriptPath); err != nil {
		return nil, err
	}

	//执行指定的函数
	luaFunc := lvm.L.GetGlobal(luaFuncName)
	if luaFunc.Type() != lua.LTFunction {
		log.SugarLogger.Debug("未获取到result函数！")
		return lua.LString("没有返回值！"), nil
	}
	log.SugarLogger.Debugf("执行lua脚本【%s】的【%s】函数，参数为：%#v", scriptPath, luaFuncName, args)
	if err := lvm.L.CallByParam(lua.P{
		Fn:      luaFunc,
		NRet:    1,
		Protect: true,
	}, args...); err != nil {
		return nil, err
	}
	rst := lvm.L.Get(-1) // returned value
	lvm.L.Pop(1)         // remove received value
	return rst, nil
}
