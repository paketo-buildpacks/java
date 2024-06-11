package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var (
	buildPack = os.Getenv("BP_UNDER_TEST")
	builder   = "paketobuildpacks/builder-jammy-buildpackless-base"
	err       error
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

	cmd := exec.Command("./mvnw", "-DskipTests=true", "clean", "package")
	cmd.Dir, err = filepath.Abs("./samples/java/war/")
	Expect(err).To(Succeed())
	out, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "failed to precompile war package, output:\n%s", out)

	suite.Run(t)
}
