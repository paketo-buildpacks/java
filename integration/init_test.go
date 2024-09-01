package integration_test

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var (
	buildPack      = os.Getenv("BP_UNDER_TEST")
	builder        = "paketobuildpacks/builder-jammy-buildpackless-base"
	testContainers = occam.NewTestContainers()
	err            error
)

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect
	if buildPack == "" {
		t.Error("BP_UNDER_TEST is missing")
	}

	format.MaxLength = 0
	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("java", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("SpringBoot", testSpringBoot)
	suite("ExecutableJar", testExecutableJar)
	suite("Tomcat", testTomcat)
	suite("TomEE", testTomee)

	cmd := exec.Command("./gradlew", "clean", "build", "-x", "test")
	cmd.Dir, err = filepath.Abs("./samples/java/war/")
	Expect(err).To(Succeed())
	out, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "failed to precompile war package, output:\n%s", out)

	cmd = exec.Command("./mvnw", "-DskipTests=true", "clean", "package")
	cmd.Dir, err = filepath.Abs("./samples/java/war-spring/")
	Expect(err).To(Succeed())
	out, err = cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "failed to precompile war-spring package, output:\n%s", out)

	suite.Run(t)
}

func makeRequest(path, port string) (*http.Response, error) {
	fullURL, err := url.Parse(fmt.Sprintf("http://localhost:%s", port))
	if err != nil {
		return nil, err
	}
	fullURL.Path = path

	resp, err := http.Get(fullURL.String())
	if err != nil {
		return nil, err
	}

	return resp, nil
}
