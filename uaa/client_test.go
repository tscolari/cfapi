package uaa_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"code.google.com/p/go-uuid/uuid"
	"github.com/tscolari/cfapi/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		subject     uaa.Client
		username    string
		password    string
		httpHandler http.Handler
		server      *httptest.Server
	)

	JustBeforeEach(func() {
		server = httptest.NewServer(httpHandler)
		subject = uaa.NewClient(server.URL)
	})

	AfterEach(func() {
		server.Close()
	})

	BeforeEach(func() {
		httpHandler = nil
	})

	Describe("Authenticate", func() {
		Context("successfully request", func() {
			BeforeEach(func() {
				username = uuid.New()
				password = uuid.New()

				httpHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					query, _ := ioutil.ReadAll(r.Body)
					values, err := url.ParseQuery(string(query))
					Expect(err).ToNot(HaveOccurred())

					Expect(values["grant_type"][0]).To(Equal("password"))
					Expect(values["username"][0]).To(Equal(username))
					Expect(values["password"][0]).To(Equal(password))
					response := `{"access_token":"1234","refresh_token":"5678","token_type":"bearer"}`
					w.Write([]byte(response))
				})
			})

			It("sends the correct information and returns the correct tokens", func() {
				tokens, err := subject.Authenticate(username, password)
				Expect(err).ToNot(HaveOccurred())
				Expect(tokens.AccessToken).To(Equal("1234"))
				Expect(tokens.RefreshToken).To(Equal("5678"))
				Expect(tokens.TokenType).To(Equal("bearer"))
			})
		})

		Context("when the UAA returns an error", func() {
			BeforeEach(func() {
				httpHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					err := `{"error":"invalid_something", "error_description":"something failed here"}`
					http.Error(w, err, 500)
				})
			})

			It("returns an error message", func() {
				_, err := subject.Authenticate("3", "4")
				Expect(err).To(MatchError("UAA Error: something failed here (invalid_something)"))
			})
		})

		Context("when there's an error parsing the response", func() {
			BeforeEach(func() {
				httpHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := `{"access_token","refresh_token": "token_type":"bearer"}`
					w.Write([]byte(response))
				})
			})

			It("returns an error message", func() {
				_, err := subject.Authenticate("5", "6")
				Expect(err.Error()).To(MatchRegexp("Failed to parse response"))
			})
		})
	})

	Describe("RefreshToken", func() {
		Context("successfully request", func() {
			var refreshToken string

			BeforeEach(func() {
				refreshToken = uuid.New()

				httpHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					Expect(r.FormValue("grant_type")).To(Equal("refresh_token"))
					Expect(r.FormValue("refresh_token")).To(Equal(refreshToken))
					response := `{"access_token":"1232","refresh_token":"5678","token_type":"bearer"}`
					w.Write([]byte(response))
				})
			})

			It("sends the correct information and returns the correct tokens", func() {
				tokens, err := subject.RefreshToken(refreshToken)
				Expect(err).ToNot(HaveOccurred())
				Expect(tokens.AccessToken).To(Equal("1232"))
				Expect(tokens.RefreshToken).To(Equal("5678"))
				Expect(tokens.TokenType).To(Equal("bearer"))
			})
		})
	})
})
