package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
  "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	excelize "github.com/xuri/excelize/v2"
	"cems-dis/model"
	rs "cems-dis/server/response"
)

const (
	contentTypeExcel = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

func (s ApiService) DownloadRawData(c *gin.Context) rs.Response {
	uid := strings.Trim(c.Param("uid"), " ")
	ts1, _ := strconv.Atoi(c.Query("ts1"))
	ts2, _ := strconv.Atoi(c.Query("ts2"))

	if len(uid) == 0 {
		return rs.Error(http.StatusBadRequest, "UID wajib diisi")
	}

	dev, _ := s.model.GetDeviceByUid(uid)
	if dev == nil {
		return rs.Error(http.StatusBadRequest, "UID tidak valid")
	}
	
	if (ts1 <= 0) || (ts2 <= 0) {
		return rs.Error(http.StatusBadRequest, "Tanggal awal dan tanggal akhir wajib diisi")
	} else if (ts1 > ts2) {
		return rs.Error(http.StatusBadRequest, "Tanggal akhir tidak boleh kurang dari tanggal awal")
	}

	rows := []model.RawData{}
	sql := s.model.DB.
		Where("(uid = ?) AND (timestamp BETWEEN ? AND ?)", uid, ts1, ts2).
		Order("timestamp")
	if err := sql.Find(&rows).Error; err != nil {
		log.Warningf("DB error: %s", err.Error())
		return rs.Error(http.StatusInternalServerError, "DB error")
	}
	if len(rows) == 0 {
		return rs.Error(http.StatusBadRequest, "Tidak ada data untuk device dan range waktu tersebut")
	}

	att := fmt.Sprintf("raw-data-%s-%s.xlsx", uid, time.Now().Format("20060102-150405"))
	tmp := fmt.Sprintf("temp/%s", att)
	excel := excelize.NewFile()
	sheet := "Raw Data"
	dftEr := func(msg string, args ...interface{}) rs.Response {
		log.Warningf(msg, args)
		return rs.Error(http.StatusInternalServerError, "Error generating excel file")
	}

	excel.SetDefaultFont("Arial")

	index, err := excel.NewSheet(sheet)
	if err != nil {
		return dftEr("Error creating excel sheet: %s", err.Error())
	}

	sensors := s.model.GetSensorDefinitions()
	excel.SetCellValue(sheet, "A6", "Timestamp")
	for i, s := range sensors {
		excel.SetCellValue(sheet, cellName(6, 2 + i), s.NameAndUnit())
	}

	excel.SetCellValue(sheet, cellName(1, 1), "UID")
	excel.SetCellValue(sheet, cellName(2, 1), "Nama")
	excel.SetCellValue(sheet, cellName(3, 1), "Latitude")
	excel.SetCellValue(sheet, cellName(4, 1), "Longitude")

	excel.SetCellValue(sheet, cellName(1, 2), dev.UID)
	excel.SetCellValue(sheet, cellName(2, 2), dev.Name)
	if dev.Latitude != nil {
		excel.SetCellValue(sheet, cellName(3, 2), *dev.Latitude)
	}
	if dev.Longitude != nil {
		excel.SetCellValue(sheet, cellName(4, 2), *dev.Longitude)
	}

	startRow := 7
	for i, r := range rows {
		row := startRow + i
		excel.SetCellValue(sheet, cellName(row, 1), time.Unix(r.Timestamp, 0))
		for j, v := range r.Values() {
			if v != nil {
				excel.SetCellValue(sheet, cellName(row, 2 + j), *v)
			}
		}
	}

	timestampFormat := "yyyy-mm-dd hh:mm:ss"
	timestampStyle, _ := excel.NewStyle(&excelize.Style{CustomNumFmt: &timestampFormat})
	excel.SetCellStyle(sheet, cellName(startRow, 1), cellName(startRow + len(rows) - 1, 1), timestampStyle)
	excel.SetColWidth(sheet, "A", "A", 20)
	excel.SetActiveSheet(index)

	if err := excel.SaveAs(tmp); err != nil {
		return dftEr("Error saving excel file: %s", err.Error())
	}
	if err := excel.Close(); err != nil {
		return dftEr("Error closing excel file: %s", err.Error())
	}

	b, err := ioutil.ReadFile(tmp)
	if err != nil {
		return dftEr("Error reading excel file: %s", err.Error())
	}
	if err := os.Remove(tmp); err != nil {
		return dftEr("Error deleting temporary excel file: %s", err.Error())
	}

	return rs.Success(func() {
		sendAttachment(c, http.StatusOK, att, contentTypeExcel, b)
	})
}

func cellName(r int, c int) string {
	col, _ := excelize.ColumnNumberToName(c)
	cel, _ := excelize.JoinCellName(col, r)
	return cel
}

func sendAttachment(c *gin.Context, code int, filename string, contentType string, data []byte) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=" + filename)
	// c.Header("Content-Type", "application/octet-stream")
	c.Data(code, contentType, data)
}
