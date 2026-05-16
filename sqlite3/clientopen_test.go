package gluasql_sqlite3_test

import (
	"testing"

	"github.com/bbernhard/gluasql"
	util "github.com/bbernhard/gluasql/util"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestClientOpen(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := getLuaDbConnection() + `
		return ok, err;
	`
	assert.NoError(L.DoString(script))

	ok := util.GetValue(L, 1)
	err := util.GetValue(L, 2)
	assert.True(ok.(bool))
	assert.Nil(err)
}
