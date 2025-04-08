package integration_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/occam"
	. "github.com/paketo-buildpacks/occam/matchers"
	"github.com/sclevine/spec"
)

func testSpringBoot(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		pack      occam.Pack
		docker    occam.Docker
		image     occam.Image
		container occam.Container
		buildLogs fmt.Stringer
		imageName string
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	it.After(func() {
		Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
		Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(imageName))).To(Succeed())
		Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
	})

	context("Maven", func() {
		it("builds SpringBoot app from source with Maven", func() {
			imageName, err := occam.RandomName()
			Expect(err).ToNot(HaveOccurred())

			image, buildLogs, err = pack.WithNoColor().Build.
				WithBuildpacks(buildPack).
				WithEnv(map[string]string{
					"BP_ARCH": "amd64",
				}).
				WithBuilder(builder).
				WithTrustBuilder().
				WithPullPolicy("if-not-present").
				Execute(imageName, "samples/java/maven")
			Expect(err).ToNot(HaveOccurred())
			Expect(buildLogs.String()).ToNot(BeEmpty())
			Expect(len(image.Buildpacks)).To(BeNumerically(">", 0))

			docker := occam.NewDocker()
			container, err = docker.Container.Run.
				WithPublish("8080").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(Serve(ContainSubstring("UP")).OnPort(8080).WithEndpoint("/actuator/health"))
		})
	})

	context("Gradle", func() {
		it("builds SpringBoot app from source with Gradle", func() {
			imageName, err := occam.RandomName()
			Expect(err).ToNot(HaveOccurred())

			image, buildLogs, err = pack.WithNoColor().Build.
				WithBuildpacks(buildPack).
				WithEnv(map[string]string{
					"BP_ARCH": "amd64",
				}).
				WithBuilder(builder).
				WithTrustBuilder().
				WithPullPolicy("if-not-present").
				Execute(imageName, "samples/java/gradle")
			Expect(err).ToNot(HaveOccurred())
			Expect(buildLogs.String()).ToNot(BeEmpty())
			Expect(len(image.Buildpacks)).To(BeNumerically(">", 0))

			docker := occam.NewDocker()
			container, err = docker.Container.Run.
				WithPublish("8080").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(Serve(ContainSubstring("UP")).OnPort(8080).WithEndpoint("/actuator/health"))
		})
	})
}
