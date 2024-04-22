package response

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const DEFAULT_PAGE_SIZE = 20

type PagingFormatter struct {
	sql					*gorm.DB
	rows				int64
	size				int
	page				int
	pages				int
}

func (p *PagingFormatter) Sql() *gorm.DB {
	return p.sql
}

func (p *PagingFormatter) Json(c *gin.Context, r Response) {
	data := gin.H{
		"error": 	r.Error, 
		"data":		r.Data, 
		"rows": 	p.rows, 
		"size": 	p.size, 
		"page": 	p.page, 
		"pages": 	p.pages, 
	}

	c.JSON(r.Code, data)
}

func NewPagingFormatter(dbSql *gorm.DB, pageSize int, pageNum int) (*PagingFormatter, error) {
	size := pageSize
	page := pageNum
	rows := int64(0)

	if size <= 0 {
		size = DEFAULT_PAGE_SIZE
	}
	if page <= 0 {
		page = 1
	}

	sql := dbSql.Offset((page - 1) * size).Limit(size)
	err := sql.Count(&rows).Error
	if err != nil {
		log.Warningf("DB error: %s", err.Error())
		return nil, err
	}

	extraPage := 0
	if rows % int64(size) > 0 {
		extraPage = 1
	}
	pages := int((rows / int64(size)) + int64(extraPage))
	if page > pages {
		page = pages
	}

	formatter := &PagingFormatter{
		sql:				sql, 
		rows: 			rows, 
		size: 			size, 
		page: 			page, 
		pages: 			pages, 
	}

	return formatter, nil
}
