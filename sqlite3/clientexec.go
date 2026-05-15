package gluasql_sqlite3

import (
	lua "github.com/yuin/gopher-lua"
)

func luaToGoValue(v lua.LValue) interface{} {
	switch val := v.(type) {
	case lua.LString:
		return string(val)
	case lua.LNumber:
		return float64(val)
	case lua.LBool:
		return bool(val)
	case *lua.LNilType:
		return nil
	default:
		return val.String()
	}
}

func clientExecMethod(L *lua.LState) int {
	client := checkClient(L)

	if client.DB == nil {
		return 0
	}

	query := L.CheckString(2)

	var args []interface{}
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, luaToGoValue(L.Get(i)))
	}

	result, err := client.DB.Exec(query, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	lastID, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()

	tbl := L.NewTable()
	tbl.RawSetString("last_id", lua.LNumber(lastID))
	tbl.RawSetString("rows_affected", lua.LNumber(rowsAffected))

	L.Push(tbl)
	L.Push(lua.LNil)
	return 2
}
