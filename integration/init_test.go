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
	javaBuildpack  string
	testContainers = occam.NewTestContainers()
)

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	output, err := exec.Command("bash", "-c", "../scripts/package.sh --version 1.2.3").CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	javaBuildpack, err = filepath.Abs("../build/buildpackage.cnb")
	Expect(err).NotTo(HaveOccurred())

	format.MaxLength = 0
	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("java", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("SpringBoot", testSpringBoot)
	suite("ExecutableJar", testExecutableJar)
	suite("Tomcat", testTomcat)
	suite("TomEE", testTomee)

	samplesDir, err := filepath.Abs("./samples")
	Expect(err).To(Succeed())

	if _, err := os.Stat(samplesDir); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/paketo-buildpacks/samples", samplesDir)
		out, err := cmd.CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), "failed to clone samples, output:\n%s", out)
	}

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
