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

func testAkka(t *testing.T, context spec.G, it spec.S) {
	context("from source", func() {
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

		it("builds an SBT app and runs it from the distribution zip", func() {
			var err error
			source, err = occam.Source(filepath.Join("samples", "java", "akka"))
			Expect(err).NotTo(HaveOccurred())

			var buildLogs fmt.Stringer
			image, buildLogs, err = pack.Build.
				WithBuildpacks(javaBuildpack).
				WithGID("123").
				WithTrustBuilder().
				WithPullPolicy("never").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), buildLogs.String())

			Expect(buildLogs).To(ContainLines(ContainSubstring("Paketo Buildpack for BellSoft Liberica")))
			Expect(buildLogs).To(ContainLines(ContainSubstring("Paketo Buildpack for SBT")))
			Expect(buildLogs).To(ContainLines(ContainSubstring("Paketo Buildpack for DistZip")))

			container, err = docker.Container.Run.
				WithPublish("8080").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable())
			Eventually(container).Should(Serve(ContainSubstring("users")).OnPort(8080).WithEndpoint("/users"))
		})
	})
}
