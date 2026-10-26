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

func testTomcat(t *testing.T, context spec.G, it spec.S) {
	var (
		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack().WithNoColor()
		docker = occam.NewDocker()
	})

	context("from source", func() {
		var (
			Expect     = NewWithT(t).Expect
			Eventually = NewWithT(t).Eventually

			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("uses Tomcat as the app server", func() {
			var err error
			source, err = occam.Source(filepath.Join("samples", "java", "war"))
			Expect(err).NotTo(HaveOccurred())

			var buildLogs fmt.Stringer
			image, buildLogs, err = pack.Build.
				WithBuildpacks(javaBuildpack).
				WithEnv(map[string]string{
					"BP_TOMCAT_VERSION": "10.1",
				}).
				WithTrustBuilder().
				WithPullPolicy("never").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), buildLogs.String())

			Expect(buildLogs).To(ContainLines(ContainSubstring("Paketo Buildpack for BellSoft Liberica")))
			Expect(buildLogs).To(ContainLines(ContainSubstring("Paketo Buildpack for Gradle")))
			Expect(buildLogs).To(ContainLines(ContainSubstring("Paketo Buildpack for Apache Tomcat")))

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable())
			Eventually(container).Should(Serve(ContainSubstring("UP")).OnPort(8080).WithEndpoint("/actuator/health"))
		})
	})

	context("precompiled", func() {
		var (
			Expect     = NewWithT(t).Expect
			Eventually = NewWithT(t).Eventually

			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("uses Tomcat as the app server", func() {
			var err error
			source, err = occam.Source(filepath.Join("samples", "java", "war-spring", "target"))
			Expect(err).NotTo(HaveOccurred())

			var buildLogs fmt.Stringer
			image, buildLogs, err = pack.Build.
				WithBuildpacks(javaBuildpack).
				WithEnv(map[string]string{
					"BP_TOMCAT_VERSION": "10.1",
				}).
				WithTrustBuilder().
				WithPullPolicy("never").
				Execute(name, filepath.Join(source, "demo-0.0.1-SNAPSHOT.war"))
			Expect(err).NotTo(HaveOccurred(), buildLogs.String())

			Expect(buildLogs).To(ContainLines(ContainSubstring("Paketo Buildpack for BellSoft Liberica")))
			Expect(buildLogs).To(ContainLines(ContainSubstring("Paketo Buildpack for Apache Tomcat")))

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable())
			Eventually(container).Should(Serve(ContainSubstring("UP")).OnPort(8080).WithEndpoint("/actuator/health"))
		})
	})
}
