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
  "cems-dis/server/web"
)

func registerMiscRoutes(engine *gin.Engine, model *model.Model) {
  engine.GET("/", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "/web/device")
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
  g.GET("/devices/:id", j(api.ApiService.GetDevice))
  g.PATCH("/devices/:id", j(api.ApiService.UpdateDevice))
  g.DELETE("/devices/:id", j(api.ApiService.DeleteDevice))
  g.GET("/devices/:id/raw-data", j(api.ApiService.CemsRawData))
  g.GET("/devices/:id/emission-data", j(api.ApiService.CemsEmissionData))
  g.GET("/devices/:id/percentage-data", j(api.ApiService.CemsPercentageData))
  
  g.POST("/cems/login", j(api.ApiService.DasLogin))
  g.POST("/cems/refresh-token", j(api.ApiService.CemsRefreshToken))
  g.POST("/cems/push-data", j(api.ApiService.CemsPushData))
  g.GET("/cems/records", j(api.ApiService.ListRawData))
  g.GET("/cems/records/:id", j(api.ApiService.GetRawDataById))

  // compatibility support for cems das-data
  g.POST("/pengiriman-das", j(api.ApiService.CemsPushData))
  g.POST("/pengiriman-das/login", j(api.ApiService.DasLoginByUid))
  g.POST("/pengiriman-das/refresh-token", j(api.ApiService.CemsRefreshToken))
}

func registerWebRoutes(engine *gin.Engine, model *model.Model) {
  g := engine.Group("web")
  
  g.GET("/device", web.Device)
  g.GET("/raw-data", web.RawData)
  g.GET("/emission-data", web.EmissionData)
  g.GET("/percentage-data", web.PercentageData)
  g.GET("/push-request", web.PushRequest)
}

func registerTemplates(engine *gin.Engine) {
  render := NewCustomRender()
  render.Debug = gin.IsDebugging()

  funMap := template.FuncMap{
    // "dropdownMenuItems": menuitem.DropdownMenuItems, 
  }

  tmplMap := map[string][]string{
    "device.html":          []string{"views/content/device.html", "views/layout/admin.html"}, 
    "raw_data.html":        []string{"views/content/raw_data.html", "views/layout/admin.html"}, 
    "emission_data.html":   []string{"views/content/emission_data.html", "views/layout/admin.html"}, 
    "percentage_data.html": []string{"views/content/percentage_data.html", "views/layout/admin.html"}, 
    "push_request.html":    []string{"views/content/push_request.html", "views/layout/admin.html"}, 
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

func Start(model *model.Model) {
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

  serverPort := config.ServerPort()
  log.Infof("Serving on port %d\n", serverPort)
  engine.Run(fmt.Sprintf("0.0.0.0:%d", serverPort))
}
