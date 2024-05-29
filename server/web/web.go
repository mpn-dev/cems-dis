package web

import (
	"net/http"
  "github.com/gin-gonic/gin"
	"cems-dis/server/session_info"
	"cems-dis/server/menuitems"
)


func Ping(c *gin.Context) {
  c.Data(http.StatusOK, "text/plain", []byte("PONG"))
}

func Devices(c *gin.Context) {
	handlePageByMenu(c, "devices", "devices.html", nil)
}

func RelayStations(c *gin.Context) {
	handlePageByMenu(c, "relay-stations", "relay_stations.html", nil)
}

func RawData(c *gin.Context) {
	handlePageByMenu(c, "raw-data", "raw_data.html", nil)
}

func EmissionData(c *gin.Context) {
	handlePageByMenu(c, "emission-data", "emission_data.html", nil)
}

func PercentageData(c *gin.Context) {
	handlePageByMenu(c, "percentage-data", "percentage_data.html", nil)
}

func Transmissions(c *gin.Context) {
	handlePageByMenu(c, "transmissions", "transmissions.html", nil)
}

func PushRequests(c *gin.Context) {
	handlePageByMenu(c, "push-requests", "push_requests.html", nil)
}

func Dashboard(c *gin.Context) {
	handlePageByMenu(c, "dashboard", "dashboard.html", nil)
}

func Overview(c *gin.Context) {
	data := gin.H{"uid": c.Query("uid")}
	handlePageByMenu(c, "overview", "overview.html", data)
}

func selectMenu(c *gin.Context, menu string) {
  si := session_info.NewSessionInfo(c)
  si.SetSelectedMenu(menu)
}

func webData(ctx *gin.Context, custom interface{}) gin.H {
	si := session_info.NewSessionInfo(ctx)
	menuItems := menuitems.New(si)
	authToken, _ := si.GetAuthToken()
	userId, _ := si.GetUserId()
	userName, _ := si.GetUserName()

	data := gin.H{
		"session_token":	authToken, 
		"csrf_token": 		authToken, 
		"app_title":			"CEMS Management", 
		"logo_left":			"/libs/logo/logo-mpn.png", 
		"user_id":				userId, 
		"user_name":			userName, 
		"menu_items":			menuItems.Html(10), 
	}

	if custom != nil {
		for k,v := range custom.(gin.H) {
			data[k] = v
		}
	}

	return data
}

func handlePageByMenu(c *gin.Context, menu string, template string, data interface{}) {
  selectMenu(c, menu)
  c.HTML(http.StatusOK, template, webData(c, data))
}
