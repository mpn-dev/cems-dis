package server

import (
  "fmt"
  "html/template"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
  log "github.com/sirupsen/logrus"

  "cems-dis/config"
  "cems-dis/model"
  "cems-dis/server/api"
  "cems-dis/server/response"
  "cems-dis/server/middleware"
  "cems-dis/server/web"
)


type Server struct {
  engine  *gin.Engine
  port    int
}


func registerMiscRoutes(engine *gin.Engine, model *model.Model) {
  engine.GET("/", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "/web/devices")
  })

  engine.GET("/ping", web.Ping)
}

func registerApiRoutes(engine *gin.Engine, model *model.Model) {
  s := api.New(model)
  j := func(fn func(api.ApiService, *gin.Context) response.Response) func(*gin.Context) {
    return func(c *gin.Context) {
      fn(s, c).Json(c)
    }
  }

  g := engine.Group("api/v1")
  g.GET("/sensors", j(api.ApiService.ListSensors))
  g.POST("/sensors", j(api.ApiService.UpdateSensor))

  g.GET("/devices", j(api.ApiService.ListDevices))
  g.POST("/devices", j(api.ApiService.InsertDevice))
  g.GET("/devices/new-secret", j(api.ApiService.GenerateDeviceSecret))
  g.GET("/devices/:uid", j(api.ApiService.GetDevice))
  g.PATCH("/devices/:uid", j(api.ApiService.UpdateDevice))
  g.DELETE("/devices/:uid", j(api.ApiService.DeleteDevice))
  g.GET("/devices/:uid/raw-data", j(api.ApiService.ListRawData))
  g.GET("/devices/:uid/latest-data", j(api.ApiService.GetLatestData))
  g.GET("/devices/:uid/emission-data", j(api.ApiService.ListEmissionData))
  g.GET("/devices/:uid/percentage-data", j(api.ApiService.ListPercentageData))

  g.GET("/relay-stations", j(api.ApiService.ListRelayStation))
  g.POST("/relay-stations", j(api.ApiService.InsertRelayStation))
  g.GET("/relay-stations/supported-protocols", j(api.ApiService.RelayStationProtocols))
  g.GET("/relay-stations/:id", j(api.ApiService.GetRelayStation))
  g.PATCH("/relay-stations/:id", j(api.ApiService.UpdateRelayStation))
  g.DELETE("/relay-stations/:id", j(api.ApiService.DeleteRelayStation))

  g.GET("/push-requests", j(api.ApiService.ListPushRequests))

  g.GET("/transmissions", j(api.ApiService.ListTransmissions))

  p := engine.Group("pengiriman-das")
  p.GET("", j(api.ApiService.ListRawData))
  p.POST("", middleware.TokenAuthMiddleware, j(api.ApiService.DasReceiveData))
  p.POST("/login", j(api.ApiService.DasLoginByUid))


  d := func(fn func(api.ApiService, *gin.Context) response.Response) func(*gin.Context) {
    return func(c *gin.Context) {
      res := fn(s, c)
      if res.IsError() {
        res.Json(c)
        return
      }

      res.Data.(func())()
    }
  }
  r := engine.Group("res")
  r.GET("/raw-data/:uid/download", d(api.ApiService.DownloadRawData))
}

func registerWebRoutes(engine *gin.Engine, model *model.Model) {
  g := engine.Group("web")
  g.GET("/sensors", web.Sensors)
  g.GET("/devices", web.Devices)
  g.GET("/relay-stations", web.RelayStations)
  g.GET("/raw-data", web.RawData)
  g.GET("/emission-data", web.EmissionData)
  g.GET("/percentage-data", web.PercentageData)
  g.GET("/transmissions", web.Transmissions)
  g.GET("/push-requests", web.PushRequests)
  g.GET("/dashboard", web.Dashboard)
}

func registerTemplates(engine *gin.Engine) {
  funcMap := func(tpl *template.Template) *template.Template{
    tpl = tpl.Funcs(template.FuncMap{
      // "dropdownMenuItems": menuitem.DropdownMenuItems, 
    })
    tpl = tpl.Delims("<%", "%>")
    return tpl
  }
  render := NewCustomRender(funcMap, gin.IsDebugging())

  render.AddFromFiles("dashboard.html", "views/content/dashboard.html", "views/layout/map.html")
  render.AddFromFiles("raw_data.html", "views/content/raw_data.html", "views/layout/admin.html")
  render.AddFromFiles("emission_data.html", "views/content/emission_data.html", "views/layout/admin.html")
  render.AddFromFiles("percentage_data.html", "views/content/percentage_data.html", "views/layout/admin.html")
  render.AddFromFiles("transmissions.html", "views/content/transmissions.html", "views/layout/admin.html")
  render.AddFromFiles("push_requests.html", "views/content/push_requests.html", "views/layout/admin.html")
  render.AddFromFiles("relay_stations.html", "views/content/relay_stations.html", "views/layout/admin.html")
  render.AddFromFiles("devices.html", "views/content/devices.html", "views/layout/admin.html")
  render.AddFromFiles("sensors.html", "views/content/sensors.html", "views/layout/admin.html")

  engine.HTMLRender = render
}

func (s Server) Start() {
  log.Info("***** START CEMS DIS *****")
  log.Infof("Serving on port %d", s.port)
  s.engine.Run(fmt.Sprintf("0.0.0.0:%d", s.port))
}

func New(model *model.Model) Server {
  secret := cookie.NewStore([]byte("secret"))
  engine := gin.Default()
  engine.ForwardedByClientIP = true
  engine.Use(sessions.Sessions("cems-session", secret))
  engine.SetTrustedProxies([]string{"127.0.0.1"})
  engine.StaticFile("/favicon.ico", "./assets/favicon.ico")
  engine.Static("/libs", "./assets/libs")
  registerTemplates(engine)
  registerApiRoutes(engine, model)
  registerWebRoutes(engine, model)
  registerMiscRoutes(engine, model)

  return Server{
    engine: engine, 
    port:   config.ServerPort(), 
  }
}
