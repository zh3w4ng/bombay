package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Building struct {
	ID      uint32   `json:"id"`
	Address string   `json:"address"`
	Floors  []string `json:"floors"`
}

var m map[uint32]Building // in-memory db
var r *chi.Mux

func setupRouter() {
	r = chi.NewRouter()
	r.Use(middleware.Logger,
		render.SetContentType(render.ContentTypeJSON),
	)
	r.Get("/ping", ping)
	r.Route("/buildings", func(r chi.Router) {
		r.Get("/", ListBuildings)
		r.Post("/", CreateBuiling) // POST /buildings
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", GetBuilding)       // GET /buildings/123
			r.Put("/", UpdateBuilding)    // PUT /buildings/123
			r.Delete("/", DeleteBuilding) // DELETE /buildings/123
		})
	})
}

func setupDB() {
	m = make(map[uint32]Building)
}

var buildings = []*Building{
	&Building{1, "City House", []string{"B1", "2"}},
	&Building{2, "Suntec City", []string{"1", "2"}},
	&Building{3, "Funan Mall", []string{"2"}},
}

func loadDummyData() {
	for _, b := range buildings {
		m[b.ID] = *b
	}
}

func main() {
	setupRouter()
	setupDB()
	loadDummyData()
	http.ListenAndServe(":3333", r)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// ➜  bombay curl http://localhost:3333/buildings
// [{"id":2,"address":"Suntec City","floors":["1","2"]},{"id":3,"address":"Funan Mall","floors":["2"]}]
func ListBuildings(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, buildings)
}

// ➜  bombay curl -X POST http://localhost:3333/buildings/ -d '{"id":4, "address":"City Hall", "floors": ["3", "4"]}'
// {"created":true}
// ➜  bombay curl -X GET http://localhost:3333/buildings/
// [{"id":1,"address":"City House","floors":["B1","2"]},{"id":2,"address":"Suntec City","floors":["1","2"]},{"id":3,"address":"Funan Mall","floors":["2"]},{"id":4,"address":"City Hall","floors":["3","4"]}]
func CreateBuiling(w http.ResponseWriter, r *http.Request) {
	var b Building
	json.NewDecoder(r.Body).Decode(&b)
	_, existed := m[b.ID]
	if existed {
		render.JSON(w, r, map[string]bool{"created": false})
	} else {
		m[b.ID] = b
		render.JSON(w, r, map[string]bool{"created": true})
	}
}

// ➜  bombay curl http://localhost:3333/buildings/1
// {"id":1,"address":"City House","floors":["B1","2"]}
func GetBuilding(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	b, okay := m[uint32(id)]
	if okay {
		render.JSON(w, r, b)
	}
}

// ➜  bombay curl -X PUT http://localhost:3333/buildings/4 -d '{"id":4, "address":"City Hall", "floors": ["3", "4", "5"]}'
// {"updated":true}
// ➜  bombay curl -X GET http://localhost:3333/buildings/
// [{"id":1,"address":"City House","floors":["B1","2"]},{"id":2,"address":"Suntec City","floors":["1","2"]},{"id":3,"address":"Funan Mall","floors":["2"]},{"id":4,"address":"City Hall","floors":["3","4","5"]}]
func UpdateBuilding(w http.ResponseWriter, r *http.Request) {
	var b Building
	json.NewDecoder(r.Body).Decode(&b)
	_, existed := m[getID(r)]

	if existed {
		b.ID = getID(r)
		m[getID(r)] = b
		render.JSON(w, r, map[string]bool{"updated": true})
	} else {
		render.JSON(w, r, map[string]bool{"updated": false})
	}
}

// ➜  bombay curl -X DELETE http://localhost:3333/buildings/1
// {"deleted":true}
func DeleteBuilding(w http.ResponseWriter, r *http.Request) {
	for _, b := range m {
		fmt.Println(b.ID)
	}
	_, existed := m[getID(r)]
	if existed {
		delete(m, getID(r))
		render.JSON(w, r, map[string]bool{"deleted": true})
	} else {
		render.JSON(w, r, map[string]bool{"deleted": false})

	}

}

// helper method
func getID(r *http.Request) uint32 {
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	return uint32(id)
}
