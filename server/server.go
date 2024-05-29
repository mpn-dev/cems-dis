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

  d := engine.Group("pengiriman-das")
  d.GET("", j(api.ApiService.ListRawData))
  d.POST("", middleware.TokenAuthMiddleware, j(api.ApiService.DasReceiveData))
  d.POST("/login", j(api.ApiService.DasLoginByUid))
}

func registerWebRoutes(engine *gin.Engine, model *model.Model) {
  g := engine.Group("web")
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
  render := NewCustomRender()
  render.Debug = gin.IsDebugging()

  funMap := template.FuncMap{
    // "dropdownMenuItems": menuitem.DropdownMenuItems, 
  }

  tmplMap := map[string][]string{
    "dashboard.html":       []string{"views/content/dashboard.html", "views/layout/map.html"}, 
    "raw_data.html":        []string{"views/content/raw_data.html", "views/layout/admin.html"}, 
    "emission_data.html":   []string{"views/content/emission_data.html", "views/layout/admin.html"}, 
    "percentage_data.html": []string{"views/content/percentage_data.html", "views/layout/admin.html"}, 
    "transmissions.html":   []string{"views/content/transmissions.html", "views/layout/admin.html"}, 
    "push_requests.html":   []string{"views/content/push_requests.html", "views/layout/admin.html"}, 
    "relay_stations.html":  []string{"views/content/relay_stations.html", "views/layout/admin.html"}, 
    "devices.html":         []string{"views/content/devices.html", "views/layout/admin.html"}, 
  }

  // todo: reload templates if env == development
  for k, v := range tmplMap {
    tpl := template.New(k).Funcs(funMap).Delims("<%", "%>")
    tpl, err := tpl.ParseFiles(v...)
    if err != nil {
      panic(err)
    }
    render.Add(k, tpl)
  }

  // todo: panic when template loading failed
  engine.HTMLRender = render
}

func (s Server) Start() {
  log.Infof("Serving on port %d\n", s.port)
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
