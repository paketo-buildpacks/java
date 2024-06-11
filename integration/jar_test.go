package integration_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
)

func testExecutableJar(t *testing.T, context spec.G, it spec.S) {
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

	it("uses precompiled executable jar", func() {
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
			Execute(imageName, "samples/java/jar")
		Expect(err).ToNot(HaveOccurred())
		Expect(buildLogs.String()).ToNot(BeEmpty())
		Expect(len(image.Buildpacks)).To(BeNumerically(">", 0))
	})
}
