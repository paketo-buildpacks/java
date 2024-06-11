package integration_test

import (
	"fmt"
	"testing"

	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/occam"
)

func testTomcat(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		pack      occam.Pack
		docker    occam.Docker
		image     occam.Image
		buildLogs fmt.Stringer
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	it.After(func() {
		Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
	})

	context("from source", func() {
		it("uses Tomcat as the app server", func() {
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
				Execute(imageName, "samples/java/war")
			Expect(err).ToNot(HaveOccurred())
			Expect(buildLogs.String()).ToNot(BeEmpty())
			Expect(len(image.Buildpacks)).To(BeNumerically(">", 0))
		})
	})

	context("precompiled", func() {
		it("uses Tomcat as the app server", func() {
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
				Execute(imageName, "samples/java/war/target/demo-0.0.1-SNAPSHOT.war")
			Expect(err).ToNot(HaveOccurred())
			Expect(buildLogs.String()).ToNot(BeEmpty())
			Expect(len(image.Buildpacks)).To(BeNumerically(">", 0))
		})
	})
}
