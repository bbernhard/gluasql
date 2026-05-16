package gluasql_mysql_test

import (
	"testing"

	"github.com/bbernhard/gluasql"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestClientSetTimeout(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := `
		c=require 'mysql'.new();
		c:set_timeout(1000);
	`

	assert.NoError(L.DoString(script))
}
