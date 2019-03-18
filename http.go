package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func start() {
	r := mux.NewRouter()
	r.HandleFunc("/_ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(version))
	})

	r.HandleFunc("/", getIndex)
	r.HandleFunc("/v1/topo", generateTopo)
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})

	logrus.Println(http.ListenAndServe(":3456", handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}

func generateTopo(w http.ResponseWriter, r *http.Request) {
	configrue, err := getNginxConfigure()
	if err != nil {
		fmt.Printf("Receive Error [%s] \n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	nginx, err := extractNginxMeta(configrue)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// fmt.Println(nginx)
	// generateSVG(w, r, nginx)
	level := generateNginxLevel(nginx)

	data, err := json.Marshal(level)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(data)
}

func generateNginxLevel(nginx map[string]map[string]nginxMeta) (nl nignxLevel) {
	nl.Name = "Nginx"
	var children []nignxLevel
	for domain, value := range nginx {
		var child nignxLevel
		child.Name = domain
		var _c []nignxLevel
		for loc, dest := range value {
			_c = append(_c, nignxLevel{
				Name: loc,
				Children: []nignxLevel{
					{
						Name: dest.Dest,
					},
				},
			})
		}
		child.Children = _c
		children = append(children, child)
	}
	nl.Children = children
	return
}
