package integration_test

import (
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

const (
	adoptiumBuildpack = "docker.io/paketobuildpacks/adoptium:latest"
	correttoBuildpack = "docker.io/paketobuildpacks/amazon-corretto:latest"
)

var javaBuildpack string

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	output, err := exec.Command("bash", "-c", "../scripts/package.sh --version 1.2.3").CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	javaBuildpack, err = filepath.Abs("../build/buildpackage.cnb")
	Expect(err).NotTo(HaveOccurred())

	format.MaxLength = 0
	SetDefaultEventuallyTimeout(60 * time.Second)

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("SpringBoot", testSpringBoot)
	suite("Tomcat", testTomcat)
	suite("TomEE", testTomee)
	suite("Liberty", testLiberty)
	suite("JavaNode", testJavaNode)
	suite("Akka", testAkka)

	samplesDir, err := filepath.Abs("./samples")
	Expect(err).To(Succeed())

	if _, err := os.Stat(samplesDir); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/paketo-buildpacks/samples", samplesDir)
		out, err := cmd.CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), "failed to clone samples, output:\n%s", out)
	}

	// The builder and run image are pulled by scripts/integration.sh. These
	// buildpacks are passed directly to the build (before the java composite)
	// so the tests can select a non-BellSoft JVM vendor until jvm-vendors is
	// integrated into the java buildpack.
	docker := occam.NewDocker()
	Expect(docker.Pull.Execute(adoptiumBuildpack)).To(Succeed())
	Expect(docker.Pull.Execute(correttoBuildpack)).To(Succeed())

	cmd := exec.Command("./mvnw", "-DskipTests=true", "clean", "package")
	cmd.Dir, err = filepath.Abs("./samples/java/maven/")
	Expect(err).To(Succeed())
	out, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "failed to precompile maven package, output:\n%s", out)

	cmd = exec.Command("./gradlew", "clean", "build", "-x", "test")
	cmd.Dir, err = filepath.Abs("./samples/java/war/")
	Expect(err).To(Succeed())
	out, err = cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "failed to precompile war package, output:\n%s", out)

	suite.Run(t)
}
