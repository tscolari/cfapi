package cf_test

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCfapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cfapi Suite")
}

func readResponseJSON(filename string) []byte {
	response, err := ioutil.ReadFile(filepath.Join("./assets/", filename))
	Expect(err).ToNot(HaveOccurred())
	return response
}
