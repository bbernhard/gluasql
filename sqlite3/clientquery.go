package gluasql_sqlite3

import (
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func clientQueryMethod(L *lua.LState) int {
	client := checkClient(L)
	query := L.ToString(2)

	if client.DB == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("database not initialized"))
		return 2
	}

	if query == "" {
		L.ArgError(2, "query string required")
		return 0
	}

	rows, err := client.DB.Query(query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	resultTable := L.NewTable()

	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan row values
		if err := rows.Scan(valuePtrs...); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		// Convert row to Lua table
		rowTable := L.NewTable()

		for i, col := range cols {
			val := values[i]

			switch v := val.(type) {
			case nil:
				rowTable.RawSetString(col, lua.LNil)

			case int:
				rowTable.RawSetString(col, lua.LNumber(v))

			case int8:
				rowTable.RawSetString(col, lua.LNumber(v))

			case int16:
				rowTable.RawSetString(col, lua.LNumber(v))

			case int32:
				rowTable.RawSetString(col, lua.LNumber(v))

			case int64:
				rowTable.RawSetString(col, lua.LNumber(v))

			case uint:
				rowTable.RawSetString(col, lua.LNumber(v))

			case uint8:
				rowTable.RawSetString(col, lua.LNumber(v))

			case uint16:
				rowTable.RawSetString(col, lua.LNumber(v))

			case uint32:
				rowTable.RawSetString(col, lua.LNumber(v))

			case uint64:
				rowTable.RawSetString(col, lua.LNumber(v))

			case float32:
				rowTable.RawSetString(col, lua.LNumber(v))

			case float64:
				rowTable.RawSetString(col, lua.LNumber(v))

			case bool:
				rowTable.RawSetString(col, lua.LBool(v))

			case string:
				rowTable.RawSetString(col, lua.LString(v))

			case []byte:
				rowTable.RawSetString(col, lua.LString(string(v)))

			case time.Time:
				rowTable.RawSetString(
					col,
					lua.LString(v.Format(time.RFC3339)),
				)

			default:
				rowTable.RawSetString(
					col,
					lua.LString(fmt.Sprint(v)),
				)
			}
		}

		resultTable.Append(rowTable)
	}

	if err := rows.Err(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(resultTable)
	return 1
}
