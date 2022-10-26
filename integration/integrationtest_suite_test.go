package integrationtests_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

const (
	tomcat8App    = "tomcat-8"
	tomcat10App   = "tomcat-10"
	springBootApp = "spring-boot"
)

var (
	buildPack string
	builder   = "paketobuildpacks/builder-jammy-buildpackless-base"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Java Buildpack Suite")
}

var _ = BeforeSuite(func() {
	format.MaxLength = 0

	var ok bool
	buildPack, ok = os.LookupEnv("BP_UNDER_TEST")
	if !ok {
		panic("BP_UNDER_TEST is missing")
	}

	SetDefaultEventuallyTimeout(10 * time.Second)
})

func makeRequest(path, port string) string {
	fullURL, err := url.Parse(fmt.Sprintf("http://localhost:%s", port))
	Expect(err).ToNot(HaveOccurred())
	fullURL.Path = path

	resp, err := http.Get(fullURL.String())
	Expect(err).ToNot(HaveOccurred())

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	return string(body)
}
