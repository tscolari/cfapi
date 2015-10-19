package cf_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/tscolari/cfapi/cf"
	"github.com/tscolari/cfapi/uaa"
	uaafakes "github.com/tscolari/cfapi/uaa/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RefresherClient", func() {
	var server *httptest.Server
	var handlerFunc http.Handler
	var tokens uaa.Tokens
	var uaaRefresher *uaafakes.FakeRefresher
	var client *cf.RefresherClient
	var response resources.ApplicationResource

	JustBeforeEach(func() {
		response = resources.ApplicationResource{}
		server = httptest.NewServer(handlerFunc)
		client = cf.NewRefresherClient(server.URL, tokens, uaaRefresher)
	})

	BeforeEach(func() {
		tokens = uaa.Tokens{
			AccessToken:  "12345",
			RefreshToken: "old-refresh-token",
		}
		uaaRefresher = new(uaafakes.FakeRefresher)
	})

	AfterEach(func() {
		server.Close()
	})

	Context("when receives a valid response", func() {
		BeforeEach(func() {
			handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
		})

		It("doesn't refresh the tokens", func() {
			err := client.Get("/app/123", nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(uaaRefresher.RefreshTokenCallCount()).To(Equal(0))
		})
	})

	Context("when the first token is not valid", func() {
		BeforeEach(func() {
			uaaRefresher.RefreshTokenReturns(&uaa.Tokens{
				AccessToken:  "refreshed-access-token",
				RefreshToken: "another-refresh-token",
			}, nil)

			handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") == "bearer refreshed-access-token" {
					w.WriteHeader(200)
					return
				}
				w.WriteHeader(401)
			})
		})

		It("refreshes the tokens", func() {
			err := client.Get("/app/123", nil)
			Expect(err).ToNot(HaveOccurred())

			Expect(uaaRefresher.RefreshTokenCallCount()).To(Equal(1))
			Expect(uaaRefresher.RefreshTokenArgsForCall(0)).To(Equal("old-refresh-token"))
		})

		Context("when `OnTokenRefresh` is given", func() {
			It("calls with the updated tokens", func() {
				client.OnTokenRefresh = func(newTokens uaa.Tokens) {
					Expect(newTokens.AccessToken).To(Equal("refreshed-access-token"))
					Expect(newTokens.RefreshToken).To(Equal("another-refresh-token"))
				}

				err := client.Get("/app/123", nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
