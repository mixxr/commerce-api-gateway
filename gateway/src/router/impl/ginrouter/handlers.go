package ginrouter

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"main/dataaccess/models"
	sanitizer "main/router/utils/sanitizer"

	"main/logger"

	"github.com/gin-gonic/gin"
)

func formatAndReturn(tables []models.ITable, err error, c *gin.Context, format string) {
	if err == nil {
		switch format {
		case "json":
			c.JSON(http.StatusOK, gin.H{
				"message": "TBD",
			})
		case "csv":
			var buffer bytes.Buffer

			for _, t := range tables {
				fmt.Fprintf(&buffer, "%s\n", t.String())
			}
			c.String(http.StatusOK, buffer.String())
		}
	} else {
		// ERROR
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
}

func (o *GinRouter) GetHandler() http.Handler {
	return o.httpengine
}

func (o *GinRouter) searchHandler(c *gin.Context) {
	owner := sanitizer.GetSearchToken(c.Param("owner"))     // can ends with -
	service := sanitizer.GetSearchToken(c.Param("service")) // can ends with -
	descr := sanitizer.GetDescr(c.Param("descr"))
	tags, ext := getExt(c, "tags", "csv")
	tags = sanitizer.GetTags(tags)
	//if owner != "" || service != "" || owner != "" || service != "" {
	tin := models.Table{Name: service, Owner: owner, Descr: descr, Tags: tags, Status: getStatus(c, owner)}
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET search", tin))
	if tin.IsValid() {
		tables, err := o.dal.ReadTables(&tin)
		//logger.CommonLog.Println("Result n.tables=", len(tables))
		formatAndReturn(models.ConvertToITables(tables), err, c, ext)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "owner and service have to be alphanumeric, use _ in place of spaces, use - in place of * for searches"})
		return
	}
}

func (o *GinRouter) getServiceHandler(c *gin.Context) {
	owner := c.Param("owner")
	service, ext := getExt(c, "service", "csv")
	if owner, service = sanitizer.CheckTokens(owner, service); owner != "" && service != "" {
		tin := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
		logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET service", tin))
		table, err := o.dal.ReadTable(&tin)
		formatAndReturn([]models.ITable{table}, err, c, ext)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "owner and service have to be alphanumeric, use _ in place of spaces"})
		return
	}
}

func (o *GinRouter) getColnamesHandler(c *gin.Context) {
	owner := c.Param("owner")
	service := c.Param("service")
	lang, ext := getExt(c, "lang", "csv")
	t := models.Table{Name: service, Owner: owner, DefLang: lang, Status: getStatus(c, owner)}
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET servicecolnames", t))
	table, err := o.dal.ReadTableColnames(&t, lang)
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET servicecolnames RESULT", table))
	formatAndReturn([]models.ITable{table}, err, c, ext)
}

func (o *GinRouter) getValuesHandler(c *gin.Context) {
	var startNum int
	var countNum int64
	owner := c.Param("owner")
	service := c.Param("service")
	startNum = getInt(c, "start", 0)
	count, ext := getExt(c, "count", "csv")
	countNum = toInt64(count, int64(dataAccessValueMaxCount))
	t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET servicevalues", t, startNum, countNum))
	table, err := o.dal.ReadTableValues(&t, startNum, countNum)
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET servicevalues RESULT", table))
	formatAndReturn([]models.ITable{table}, err, c, ext)
}

func (o *GinRouter) postServiceHandler(c *gin.Context) {
	owner := c.Param("owner")
	service, ext := getExt(c, "service", "csv")
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "CREATE service", owner, service))
	if !isOwner(c, owner) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
		return
	}
	var tableJson models.Table
	if err := c.ShouldBindJSON(&tableJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if tableJson.Name != service || tableJson.Owner != owner {
		c.JSON(http.StatusBadRequest, gin.H{"status": "service or owner does not match the url"})
		return
	}
	err := o.dal.StoreTable(&tableJson)
	formatAndReturn([]models.ITable{&tableJson}, err, c, ext)
}

func (o *GinRouter) putServiceHandler(c *gin.Context) {
	owner := c.Param("owner")
	service, ext := getExt(c, "service", "csv")
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "UPDATE service", owner, service))
	if !isOwner(c, owner) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
		return
	}
	var tableJson models.Table
	if err := c.ShouldBindJSON(&tableJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if tableJson.Name != service || tableJson.Owner != owner {
		c.JSON(http.StatusBadRequest, gin.H{"status": "service or owner does not match the url"})
		return
	}
	err := o.dal.UpdateTable(&tableJson)
	formatAndReturn([]models.ITable{&tableJson}, err, c, ext)
}

func (o *GinRouter) postColnamesHandler(c *gin.Context) {
	owner := c.Param("owner")
	service := c.Param("service")
	lang, ext := getExt(c, "lang", "csv")
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "CREATE servicecolnames", owner, service, lang))
	if !isOwner(c, owner) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "you are not allowed"})
		return
	}
	var tableJson models.TableColnames
	if err := c.ShouldBindJSON(&tableJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if tableJson.Lang != lang {
		c.JSON(http.StatusBadRequest, gin.H{"message": "lang does not match the url"})
		return
	}
	tableJson.SetParent(&models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)})

	err := o.dal.StoreTableColnames(&tableJson)
	formatAndReturn([]models.ITable{&tableJson}, err, c, ext)
}

func (o *GinRouter) postValuesHandler(c *gin.Context) {
	owner := c.Param("owner")
	service := c.Param("service")
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "CREATE servicevalues", owner, service))
	if !isOwner(c, owner) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
		return
	}
	var tableJson models.TableValues
	if err := c.ShouldBindJSON(&tableJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tableJson.SetParent(&models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)})

	err := o.dal.StoreTableValues(&tableJson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": tableJson.Count})
	return
}

func (o *GinRouter) deleteServiceHandler(c *gin.Context) {
	owner := c.Param("owner")
	service, ext := getExt(c, "service", "csv")
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "DELETE service", owner, service))
	if !isOwner(c, owner) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
		return
	}
	t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
	err := o.dal.DeleteTable(&t)
	formatAndReturn([]models.ITable{&t}, err, c, ext) // TODO: StatusAccepted 202
}

func (o *GinRouter) deleteColnamesHandler(c *gin.Context) {
	owner := c.Param("owner")
	service := c.Param("service")
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "DELETE servicecolnames", owner, service))
	if !isOwner(c, owner) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
		return
	}
	langs, ext := getExt(c, "langs", "csv")
	langs = strings.Trim(langs, "/")
	t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
	err := o.dal.DeleteTableColnames(&t, strings.Split(langs, "/"))
	formatAndReturn([]models.ITable{&t}, err, c, ext) // TODO: StatusAccepted 202
	return
}

func (o *GinRouter) deleteValuesHandler(c *gin.Context) {
	owner := c.Param("owner")
	service := c.Param("service")
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "DELETE servicevalues", owner, service))
	if !isOwner(c, owner) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
		return
	}
	//startNum := getInt(c, "start", 0)
	count, ext := getExt(c, "count", "csv")
	countNum := toInt64(count, 0)
	t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
	err := o.dal.DeleteTableValues(&t, countNum)
	formatAndReturn([]models.ITable{&t}, err, c, ext)

	return
}
