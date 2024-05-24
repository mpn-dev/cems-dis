package menuitems

import (
	"fmt"
  "html/template"
  "cems-dis/server/session_info"
)

type menuItem struct {
  name  string
  text  string
  path  string
  exts  string
}

type menuItems struct {
	items		[]menuItem
	sInfo		session_info.SessionInfo
}


func allMenuItems() []menuItem {
  return []menuItem{
    // menuItem{name: "overview",          text: "Overview",               path: "/web/overview",          exts: ""}, 
    // menuItem{name: "map",               text: "Map",                    path: "/web/map",               exts: ""}, 
    // menuItem{}, 
    menuItem{name: "raw-data",          text: "Raw Data",               path: "/web/raw-data",          exts: ""}, 
    menuItem{name: "emission-data",     text: "Emission Data",          path: "/web/emission-data",     exts: ""}, 
    menuItem{name: "percentage-data",   text: "Percentage Data",        path: "/web/percentage-data",   exts: ""}, 
    menuItem{}, 
    menuItem{name: "transmissions",     text: "Transmissions",          path: "/web/transmissions",     exts: ""}, 
    menuItem{name: "push-requests",     text: "Push Requests",          path: "/web/push-requests",     exts: ""}, 
    menuItem{}, 
    menuItem{name: "relay-stations",    text: "Daftar Relay Station",   path: "/web/relay-stations",    exts: ""}, 
    menuItem{name: "devices",           text: "Daftar Device",          path: "/web/devices",           exts: ""}, 
  }
}

func (m menuItems) Html(indent int) template.HTML {
  padding := ""
  for p := 0; p < indent; p++ {
    padding = padding + " "
  }

	selectedMenu, _ := m.sInfo.GetSelectedMenu()
  html := padding + "<div class=\"dropdown-menu\">\n"
  for _, i := range m.items {
		line := "  <div class=\"dropdown-divider\"></div>\n"
		if len(i.text) > 0 {
			active := ""
			if i.name == selectedMenu {
				active = " active"
			}
			line = fmt.Sprintf("  <a class=\"dropdown-item%s\" href=\"%s\"%s>%s</a>\n", active, i.path, i.exts, i.text)
		}

    html = html + padding + line
  }
  html = html + padding + "</div>\n"

  return template.HTML(html)
}

func New(i session_info.SessionInfo) menuItems {
	return menuItems{
		items: 	allMenuItems(), 
		sInfo:	i, 
	}
}
