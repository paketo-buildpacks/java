package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/occam"
	. "github.com/paketo-buildpacks/occam/matchers"
	"github.com/sclevine/spec"
)

func testLiberty(t *testing.T, context spec.G, it spec.S) {
	context("precompiled war", func() {
		var (
			Expect     = NewWithT(t).Expect
			Eventually = NewWithT(t).Eventually

			pack   occam.Pack
			docker occam.Docker

			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error

			pack = occam.NewPack().WithNoColor()
			docker = occam.NewDocker()

			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("uses Liberty as the app server", func() {
			var err error
			source, err = occam.Source(filepath.Join("samples", "java", "war", "build", "libs"))
			Expect(err).NotTo(HaveOccurred())

			var buildLogs fmt.Stringer
			image, buildLogs, err = pack.Build.
				WithBuildpacks(javaBuildpack).
				WithEnv(map[string]string{
					"BP_JAVA_APP_SERVER":  "liberty",
					"BP_LIBERTY_PROFILE":  "webProfile11",
					"BP_LIBERTY_FEATURES": "webProfile-11.0",
				}).
				WithTrustBuilder().
				WithPullPolicy("never").
				Execute(name, filepath.Join(source, "war-0.0.1-SNAPSHOT.war"))
			Expect(err).NotTo(HaveOccurred(), buildLogs.String())

			Expect(buildLogs).To(ContainLines(ContainSubstring("Paketo Buildpack for BellSoft Liberica")))
			Expect(buildLogs).To(ContainLines(ContainSubstring("Paketo Buildpack for Liberty")))

			// Liberty's default HTTP port is 9080, not 8080.
			container, err = docker.Container.Run.
				WithPublish("9080").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable())
			Eventually(container).Should(Serve(ContainSubstring("UP")).OnPort(9080).WithEndpoint("/actuator/health"))
		})
	})
}
