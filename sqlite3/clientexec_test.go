package gluasql_sqlite3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluasql"
	util "github.com/tengattack/gluasql/util"
	lua "github.com/yuin/gopher-lua"
)

func TestClientExec(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := getLuaDbConnection() + `	
		res, err = c:exec("insert into mytable(name, email, phone_numbers, created_at) values(?, ?, ?, ?)", "a", "b", "c", "d")
		return res, err;
	`

	assert.NoError(L.DoString(script))

	res := util.GetValue(L, 1)
	err := util.GetValue(L, 2)
	assert.Nil(err)

	assert.Equal(
		map[string]interface{}{
			"last_id":       5,
			"rows_affected": 1,
		}, res)
}

func TestClientCreateDb(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := `
		c = require 'sqlite3'.new();
		ok, err = c:open("/tmp/my.db", { cache = "shared", mode = "rwc" });
		if ok then
			res, err = c:exec("create table if not exists mytable (name text not null)");
			return res, err
		else
			return ok, err
		end
	`
	assert.NoError(L.DoString(script))
}
