package cfapi_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/tscolari/cfapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CfClient", func() {
	var server *httptest.Server
	var handlerFunc http.Handler
	var client *cfapi.CfClient
	var response resources.ApplicationResource

	JustBeforeEach(func() {
		response = resources.ApplicationResource{}
		server = httptest.NewServer(handlerFunc)
		client = cfapi.NewCfClient(server.URL, "my-access-token", "my-refresh-token")
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Get", func() {
		BeforeEach(func() {
			handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).To(Equal("GET"))
				Expect(r.Header.Get("Authorization")).To(Equal("bearer my-access-token"))

				w.Header().Set("Content-Type", "application/json")
				w.Write(readResponseJSON("app-response.json"))
			})
		})

		It("sends the correct request, and parses the response", func() {
			err := client.Get("/app/123", &response)
			Expect(err).ToNot(HaveOccurred())

			Expect(response.Metadata.Guid).To(Equal("49934910-756a-46c5-bae1-b82540e28937"))
			Expect(*response.Entity.Name).To(Equal("name-475"))
			Expect(*response.Entity.Memory).To(BeEquivalentTo(1024))
		})

		Context("when something goes wrong", func() {
			Context("when cloudcontroller returns an error", func() {
				BeforeEach(func() {
					handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						http.Error(w, `{"description": "Ups"}`, http.StatusInternalServerError)
					})
				})

				It("returns the correct error message", func() {
					err := client.Get("/app/123", &response)
					Expect(err).To(MatchError("Ups"))
				})
			})

			Context("when cloud controller can't be reached", func() {
				JustBeforeEach(func() {
					client = cfapi.NewCfClient("http://invalid.example.com", "my-access-token", "my-refresh-token")
				})

				It("returns the correct error message", func() {
					err := client.Get("/app/123", &response)
					Expect(err.Error()).To(ContainSubstring("Failed to connect"))
				})
			})

			Context("when cloud controller returns an invalid json", func() {
				BeforeEach(func() {
					handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						w.Write([]byte(`{"invalid_json": } true}`))
					})
				})

				It("returns the correct error message", func() {
					err := client.Get("/app/123", &response)
					Expect(err.Error()).To(ContainSubstring("Failed to parse response"))
				})
			})
		})
	})

	Describe("Post", func() {
		BeforeEach(func() {
			handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).To(Equal("POST"))
				Expect(r.Header.Get("Authorization")).To(Equal("bearer my-access-token"))

				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				Expect(err).ToNot(HaveOccurred())
				var options map[string]string

				err = json.Unmarshal(body, &options)
				Expect(err).ToNot(HaveOccurred())
				Expect(options["key1"]).To(Equal("value1"))
				Expect(options["key2"]).To(Equal("value2"))

				w.Header().Set("Content-Type", "application/json")
				w.Write(readResponseJSON("app-response.json"))
			})
		})

		It("sends the correct request, and parses the response", func() {
			options := map[string]string{
				"key1": "value1",
				"key2": "value2",
			}

			err := client.Post("/app/123", options, &response)
			Expect(err).ToNot(HaveOccurred())

			Expect(response.Metadata.Guid).To(Equal("49934910-756a-46c5-bae1-b82540e28937"))
			Expect(*response.Entity.Name).To(Equal("name-475"))
			Expect(*response.Entity.Memory).To(BeEquivalentTo(1024))
		})
	})

	Describe("Put", func() {
		BeforeEach(func() {
			handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).To(Equal("PUT"))
				Expect(r.Header.Get("Authorization")).To(Equal("bearer my-access-token"))

				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				Expect(err).ToNot(HaveOccurred())
				var options map[string]string

				err = json.Unmarshal(body, &options)
				Expect(err).ToNot(HaveOccurred())
				Expect(options["key1"]).To(Equal("value1"))
				Expect(options["key2"]).To(Equal("value2"))

				w.Header().Set("Content-Type", "application/json")
				w.Write(readResponseJSON("app-response.json"))
			})
		})

		It("sends the correct request, and parses the response", func() {
			options := map[string]string{
				"key1": "value1",
				"key2": "value2",
			}

			err := client.Put("/app/123", options, &response)
			Expect(err).ToNot(HaveOccurred())

			Expect(response.Metadata.Guid).To(Equal("49934910-756a-46c5-bae1-b82540e28937"))
			Expect(*response.Entity.Name).To(Equal("name-475"))
			Expect(*response.Entity.Memory).To(BeEquivalentTo(1024))
		})
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Method).To(Equal("DELETE"))
				Expect(r.Header.Get("Authorization")).To(Equal("bearer my-access-token"))

				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				Expect(err).ToNot(HaveOccurred())
				var options map[string]string

				err = json.Unmarshal(body, &options)
				Expect(err).ToNot(HaveOccurred())
				Expect(options["key1"]).To(Equal("value1"))
				Expect(options["key2"]).To(Equal("value2"))

				w.Header().Set("Content-Type", "application/json")
				w.Write(readResponseJSON("app-response.json"))
			})
		})

		It("sends the correct request, and parses the response", func() {
			options := map[string]string{
				"key1": "value1",
				"key2": "value2",
			}

			err := client.Delete("/app/123", options)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
