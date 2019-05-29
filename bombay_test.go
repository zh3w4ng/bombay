package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books Suite")
}

var _ = Describe("Bombay", func() {
	BeforeEach(func() {
		setupDB()
		loadDummyData()
	})
	Context("loadDummyData", func() {
		It("should load 3 buildings", func() {
			Expect(len(m)).To(Equal(3))
		})
	})

	Context("Controllers", func() {
		BeforeEach(func() {
			setupRouter()

		})
		Context("GetBuilding", func() {
			It("should get building #1", func() {
				req, _ := http.NewRequest("GET", "/buildings/1", nil)
				respRec := httptest.NewRecorder()
				r.ServeHTTP(respRec, req)
				Expect(respRec.Code).To(Equal(http.StatusOK))
				b, _ := json.Marshal(m[1])
				Expect(respRec.Body.String()).To(Equal(string(b) + "\n"))
			})
		})
		Context("ListBuildings", func() {
			It("should list all buildings", func() {
				req, _ := http.NewRequest("GET", "/buildings", nil)
				respRec := httptest.NewRecorder()
				r.ServeHTTP(respRec, req)
				Expect(respRec.Code).To(Equal(http.StatusOK))
				b, _ := json.Marshal(buildings)
				Expect(respRec.Body.String()).To(Equal(string(b) + "\n"))
			})
		})
		Context("CreateBuilding", func() {
			It("should create building if ID is new", func() {
				var payload = []byte(`{"id":4, "address":"City Hall", "floors": ["3", "4"]}`)
				req, _ := http.NewRequest("POST", "/buildings", bytes.NewBuffer(payload))
				respRec := httptest.NewRecorder()
				r.ServeHTTP(respRec, req)
				Expect(respRec.Code).To(Equal(http.StatusOK))
				Expect(respRec.Body.String()).To(Equal(`{"created":true}` + "\n"))
				Expect(len(m)).To(Equal(4))
			})
			It("should NOT create building if ID exists", func() {
				var payload = []byte(`{"id":1, "address":"City Hall", "floors": ["3", "4"]}`)
				req, _ := http.NewRequest("POST", "/buildings", bytes.NewBuffer(payload))
				respRec := httptest.NewRecorder()
				r.ServeHTTP(respRec, req)
				Expect(respRec.Code).To(Equal(http.StatusOK))
				Expect(respRec.Body.String()).To(Equal(`{"created":false}` + "\n"))
				Expect(len(m)).To(Equal(3))
			})
		})
		Context("UpdateBuilding", func() {
			It("should update building if ID exists", func() {
				var payload = []byte(`{"address":"City Hall", "floors": ["3", "4"]}`)
				req, _ := http.NewRequest("PUT", "/buildings/3", bytes.NewBuffer(payload))
				respRec := httptest.NewRecorder()
				r.ServeHTTP(respRec, req)
				Expect(respRec.Code).To(Equal(http.StatusOK))
				Expect(respRec.Body.String()).To(Equal(`{"updated":true}` + "\n"))
			})
			It("should update building if ID is not found", func() {
				var payload = []byte(`{"address":"City Hall", "floors": ["3", "4"]}`)
				req, _ := http.NewRequest("PUT", "/buildings/5", bytes.NewBuffer(payload))
				respRec := httptest.NewRecorder()
				r.ServeHTTP(respRec, req)
				Expect(respRec.Code).To(Equal(http.StatusOK))
				Expect(respRec.Body.String()).To(Equal(`{"updated":false}` + "\n"))
			})
		})
		Context("DeleteBuilding", func() {
			It("should delete building if ID exists", func() {
				req, _ := http.NewRequest("DELETE", "/buildings/3", nil)
				respRec := httptest.NewRecorder()
				r.ServeHTTP(respRec, req)
				Expect(respRec.Code).To(Equal(http.StatusOK))
				Expect(respRec.Body.String()).To(Equal(`{"deleted":true}` + "\n"))
				Expect(len(m)).To(Equal(2))
			})
			It("should delete building if ID is not found", func() {
				req, _ := http.NewRequest("DELETE", "/buildings/5", nil)
				respRec := httptest.NewRecorder()
				r.ServeHTTP(respRec, req)
				Expect(respRec.Code).To(Equal(http.StatusOK))
				Expect(respRec.Body.String()).To(Equal(`{"deleted":false}` + "\n"))
				Expect(len(m)).To(Equal(3))
			})
		})
	})

})
