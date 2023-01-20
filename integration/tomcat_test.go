package integrationtests_test

import (
	"context"
	"fmt"
	"os"

	"github.com/testcontainers/testcontainers-go"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/paketo-buildpacks/occam"
	"github.com/testcontainers/testcontainers-go/wait"
)

var _ = Describe("Tomcat", Label("integration"), func() {
	var (
		container      testcontainers.Container
		pack           occam.Pack
		docker         occam.Docker
		testContainers occam.TestContainers
		image          occam.Image
		buildLogs      fmt.Stringer
	)

	AfterEach(func() {
		currentTest := CurrentSpecReport()
		if currentTest.Failed() {
			_, _ = fmt.Fprintf(os.Stderr, "Skipping cleanup for test \"%s\" because of failure (container: %s, image: %s)", currentTest.LeafNodeText, container.GetContainerID(), image.ID)
			return
		}

		err := container.Terminate(context.Background())
		Expect(err).ToNot(HaveOccurred())

		err = docker.Image.Remove.Execute(image.ID)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
		testContainers = occam.NewTestContainers()
	})

	DescribeTable("Tomcat Tests",
		func(appPath string, envVars map[string]string, runtimeMatchers []types.GomegaMatcher, tomcatMatcher types.GomegaMatcher) {
			imageName, err := occam.RandomName()
			Expect(err).ToNot(HaveOccurred())

			By("builds the image", func() {
				image, buildLogs, err = pack.WithNoColor().Build.
					WithBuildpacks(buildPack).
					WithEnv(envVars).
					WithBuilder(builder).
					WithPullPolicy("if-not-present").
					Execute(imageName, appPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(buildLogs.String()).ToNot(BeEmpty())
				Expect(len(image.Buildpacks)).To(BeNumerically(">", 0))
			})

			By("runs test containers", func() {
				container, err = testContainers.WithExposedPorts("8080/tcp").
					WithWaitingFor(wait.ForLog("Server startup in")).
					Execute(imageName)
				Expect(err).ShouldNot(HaveOccurred())
			})

			By("runs http tests", func() {
				mappedPort, err := container.MappedPort(context.Background(), "8080/tcp")
				Expect(err).ShouldNot(HaveOccurred())

				runtimeResponse := makeRequest("/runtime", mappedPort.Port())
				for _, matcher := range runtimeMatchers {
					Expect(runtimeResponse).To(matcher)
				}

				Expect(makeRequest("/tomcat", mappedPort.Port())).To(tomcatMatcher)
			})
		},
		// Entry("Tomcat 8 Java 8", tomcat8App, map[string]string{
		// 	"BP_JVM_VERSION":    "8",
		// 	"BP_TOMCAT_VERSION": "8",
		// }, []types.GomegaMatcher{ContainSubstring("java.version=1.8.0")}, ContainSubstring("Apache Tomcat/8")),
		Entry("Tomcat 8 Java 11", tomcat8App, map[string]string{
			"BP_JVM_VERSION":    "11",
			"BP_TOMCAT_VERSION": "8",
		}, []types.GomegaMatcher{ContainSubstring("java.version=11")}, ContainSubstring("Apache Tomcat/8")),
		// Entry("Tomcat 8 Java 17", tomcat8App, map[string]string{
		// 	"BP_JVM_VERSION":    "17",
		// 	"BP_TOMCAT_VERSION": "8",
		// }, []types.GomegaMatcher{ContainSubstring("java.version=17")}, ContainSubstring("Apache Tomcat/8")),
		// Entry("Tomcat 10 Java 11", tomcat10App, map[string]string{
		// 	"BP_JVM_VERSION":    "11",
		// 	"BP_TOMCAT_VERSION": "10",
		// }, []types.GomegaMatcher{ContainSubstring("java.version=11")}, ContainSubstring("Apache Tomcat/10")),
		// Entry("Tomcat 10 Java 17", tomcat10App, map[string]string{
		// 	"BP_JVM_VERSION":    "17",
		// 	"BP_TOMCAT_VERSION": "10",
		// }, []types.GomegaMatcher{ContainSubstring("java.version=17")}, ContainSubstring("Apache Tomcat/10")),
	)
})
