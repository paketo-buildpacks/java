package integrationtests_test

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// 	"github.com/onsi/gomega/types"
// 	"github.com/paketo-buildpacks/occam"
// 	"github.com/testcontainers/testcontainers-go"
// 	"github.com/testcontainers/testcontainers-go/wait"
// )

// var _ = Describe("Tomee", Label("integration"), func() {
// 	var (
// 		container      testcontainers.Container
// 		pack           occam.Pack
// 		docker         occam.Docker
// 		testContainers occam.TestContainers
// 		image          occam.Image
// 		buildLogs      fmt.Stringer
// 	)

// 	AfterEach(func() {
// 		currentTest := CurrentSpecReport()
// 		if currentTest.Failed() {
// 			_, _ = fmt.Fprintf(os.Stderr, "Skipping cleanup for test \"%s\" because of failure (container: %s, image: %s)", currentTest.LeafNodeText, container.GetContainerID(), image.ID)
// 			return
// 		}

// 		err := container.Terminate(context.Background())
// 		Expect(err).ToNot(HaveOccurred())

// 		err = docker.Image.Remove.Execute(image.ID)
// 		Expect(err).ToNot(HaveOccurred())
// 	})

// 	BeforeEach(func() {
// 		pack = occam.NewPack()
// 		docker = occam.NewDocker()
// 		testContainers = occam.NewTestContainers()
// 	})

// 	DescribeTable("TomEE Tests",
// 		func(appPath string, envVars map[string]string, runtimeMatchers []types.GomegaMatcher, tomcatMatcher types.GomegaMatcher) {
// 			imageName, err := occam.RandomName()
// 			Expect(err).ToNot(HaveOccurred())

// 			By("builds the image", func() {
// 				image, buildLogs, err = pack.WithNoColor().Build.
// 					WithBuildpacks(buildPack).
// 					WithEnv(envVars).
// 					WithBuilder(builder).
// 					WithPullPolicy("if-not-present").
// 					Execute(imageName, appPath)
// 				Expect(err).ToNot(HaveOccurred())
// 				Expect(buildLogs.String()).ToNot(BeEmpty())
// 				Expect(len(image.Buildpacks)).To(BeNumerically(">", 0))
// 			})

// 			By("runs test containers", func() {
// 				container, err = testContainers.WithExposedPorts("8080/tcp").
// 					WithWaitingFor(wait.ForLog("Server startup in")).
// 					Execute(imageName)
// 				Expect(err).ShouldNot(HaveOccurred())
// 			})

// 			By("runs http tests", func() {
// 				mappedPort, err := container.MappedPort(context.Background(), "8080/tcp")
// 				Expect(err).ShouldNot(HaveOccurred())

// 				runtimeResponse := makeRequest("/runtime", mappedPort.Port())
// 				for _, matcher := range runtimeMatchers {
// 					Expect(runtimeResponse).To(matcher)
// 				}

// 				Expect(makeRequest("/tomcat", mappedPort.Port())).To(tomcatMatcher)
// 			})
// 		},
// 		Entry("Tomee 7 Java 8", tomcat8App, map[string]string{
// 			"BP_JVM_VERSION":     "8",
// 			"BP_JAVA_APP_SERVER": "tomee",
// 			"BP_TOMEE_VERSION":   "7",
// 		}, []types.GomegaMatcher{ContainSubstring("java.version=1.8.0")}, MatchRegexp(`Apache Tomcat \(TomEE\)/8\.[\d]+\.[\d]+ \(7\.[\d]+\.[\d]+\)`)),
// 		Entry("Tomee 7 Java 11", tomcat8App, map[string]string{
// 			"BP_JVM_VERSION":     "11",
// 			"BP_JAVA_APP_SERVER": "tomee",
// 			"BP_TOMEE_VERSION":   "7",
// 		}, []types.GomegaMatcher{ContainSubstring("java.version=11")}, MatchRegexp(`Apache Tomcat \(TomEE\)/8\.[\d]+\.[\d]+ \(7\.[\d]+\.[\d]+\)`)),
// 		Entry("Tomee 7 Java 17", tomcat8App, map[string]string{
// 			"BP_JVM_VERSION":     "17",
// 			"BP_JAVA_APP_SERVER": "tomee",
// 			"BP_TOMEE_VERSION":   "7",
// 		}, []types.GomegaMatcher{ContainSubstring("java.version=17")}, ContainSubstring("Apache Tomcat/8")), // When running with the JVM 17, org.apache.catalina.util.ServerInfo omits the Tomee information
// 		Entry("Tomee 9 Java 11", tomcat10App, map[string]string{
// 			"BP_JVM_VERSION":     "11",
// 			"BP_JAVA_APP_SERVER": "tomee",
// 			"BP_TOMEE_VERSION":   "9.0.0-M8",
// 		}, []types.GomegaMatcher{ContainSubstring("java.version=11")}, MatchRegexp(`Apache Tomcat \(TomEE\)/10\.[\d]+\.[\d]+ \(9\.0\.0-M8\)`)),
// 		Entry("Tomee 9 Java 17", tomcat10App, map[string]string{
// 			"BP_JVM_VERSION":     "17",
// 			"BP_JAVA_APP_SERVER": "tomee",
// 			"BP_TOMEE_VERSION":   "9.0.0-M8",
// 		}, []types.GomegaMatcher{ContainSubstring("java.version=17")}, ContainSubstring("Apache Tomcat/10")), // When running with the JVM 17, org.apache.catalina.util.ServerInfo omits the Tomee information
// 	)
// })
