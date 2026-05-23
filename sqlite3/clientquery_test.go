package gluasql_sqlite3_test

import (
	"testing"
	"time"

	"github.com/bbernhard/gluasql"
	util "github.com/bbernhard/gluasql/util"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestClientQuery(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := getLuaDbConnection() + `
		res, err = c:query("SELECT * FROM mytable LIMIT 2");
		return res, err;
	`

	assert.NoError(L.DoString(script))

	res := util.GetValue(L, 1)
	err := util.GetValue(L, 2)
	assert.Nil(err)

	// source:
	// sql.NewRow("John Doe", "john@doe.com", []string{"555-555-555"}, timeNow)
	// sql.NewRow("John Doe", "johnalt@doe.com", []string{}, timeNow)
	timeNowFormat := timeNow.UTC().Format("2006-01-02 15:04:05")
	assert.Equal([]interface{}{
		map[string]interface{}{
			"email":         "john@doe.com",
			"phone_numbers": "[\"555-555-555\"]",
			"name":          "John Doe",
			"created_at":    timeNowFormat,
		},
		map[string]interface{}{
			"email":         "johnalt@doe.com",
			"phone_numbers": "[]",
			"name":          "John Doe",
			"created_at":    timeNowFormat,
		},
	}, res)
}

func TestClientQueryJsonAndTimestamp(t *testing.T) {
	assert := assert.New(t)

	L := lua.NewState()
	defer L.Close()

	gluasql.Preload(L)

	timeNowFormat := timeNow.UTC().Format(time.RFC3339)

	script := getLuaDbConnection() + `
		c:exec([[
			CREATE TABLE IF NOT EXISTS messages (
				id INTEGER PRIMARY KEY,
				data JSON,
				timestamp DATETIME
			);
		]])

		c:exec([[
			INSERT INTO messages(data, timestamp)
			VALUES(
				'{"message":"hello","count":1}',
				'` + timeNowFormat + `'
			);
		]])

		res, err = c:query("SELECT data, timestamp FROM messages LIMIT 1");
		return res, err;
	`

	assert.NoError(L.DoString(script))

	res := util.GetValue(L, 1)
	err := util.GetValue(L, 2)

	assert.Nil(err)

	assert.Equal([]interface{}{
		map[string]interface{}{
			"data":      `{"message":"hello","count":1}`,
			"timestamp": timeNowFormat,
		},
	}, res)
}
