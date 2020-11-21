package ginrouter

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func toInt(s string, def int) int {
	intVar, err := strconv.Atoi(s)
	if err == nil {
		return intVar
	}
	return def
}

func toInt64(s string, def int64) int64 {
	intVar, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return intVar
	}
	return def
}

func getExt(c *gin.Context, paramname string, defval string) (string, string) {
	value := c.Param(paramname)
	i := strings.LastIndex(value, ".")
	if i < 0 {
		return value, defval
	} else {
		return value[:i], value[i+1:]
	}
}

func getInt(c *gin.Context, paramname string, defval int) int {
	value := c.Param(paramname)
	return toInt(value, defval)
}
