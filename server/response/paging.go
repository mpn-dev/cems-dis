package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const DEFAULT_PAGE_SIZE = 20

type Paging struct {
	sql					*gorm.DB
	rows				int64
	size				int
	page				int
	pages				int
}

func (p *Paging) Sql() *gorm.DB {
	return p.sql
}

func (p *Paging) H(r Response, d gin.H) gin.H {
	data := addMetaField(d)
	data["meta"].(gin.H)["pagination"] = gin.H{
		"rows": 		p.rows, 
		"size": 		p.size, 
		"page": 		p.page, 
		"pages": 		p.pages, 
	}

	return data
}

func NewPaging(dbSql *gorm.DB, pageNum int, pageSize int) (*Paging, error) {
	size := pageSize
	page := pageNum
	rows := int64(0)

	if size <= 0 {
		size = DEFAULT_PAGE_SIZE
	}
	if page <= 0 {
		page = 1
	}

	err := dbSql.Count(&rows).Error
	if err != nil {
		log.Warningf("DB error: %s", err.Error())
		return nil, errors.New("DB error")
	}

	extraPage := 0
	if rows % int64(size) > 0 {
		extraPage = 1
	}
	pages := int((rows / int64(size)) + int64(extraPage))
	if page > pages {
		page = pages
	}

	formatter := &Paging{
		sql:				dbSql.Offset((page - 1) * size).Limit(size), 
		rows: 			rows, 
		size: 			size, 
		page: 			page, 
		pages: 			pages, 
	}

	return formatter, nil
}

func (r Response) WithPaging(paging *Paging) Response {
	r.Formats = append(r.Formats, paging)
	return r
}

func (r Response) DefaultWithPaging(paging *Paging) Response {
	r.Formats = []Formatter{
		NewDefaultFormatter(), 
		paging, 
	}
	return r
}
