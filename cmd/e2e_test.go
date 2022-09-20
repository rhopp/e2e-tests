package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/redhat-appstudio/e2e-tests/pkg/framework"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	"github.com/onsi/gomega"
	_ "github.com/redhat-appstudio/e2e-tests/tests/build"
	_ "github.com/redhat-appstudio/e2e-tests/tests/cluster-registration"
	. "github.com/redhat-appstudio/e2e-tests/tests/e2e-demos"
	_ "github.com/redhat-appstudio/e2e-tests/tests/has"
	_ "github.com/redhat-appstudio/e2e-tests/tests/integration-service"
	_ "github.com/redhat-appstudio/e2e-tests/tests/release"

	"flag"

	"github.com/spf13/viper"

	"k8s.io/klog/v2"
)

const ()

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	return nil
}, func(data []byte) {})

var webhookConfigPath string
var demoSuitesPath string

func init() {
	rootDir, _ := os.Getwd()
	flag.StringVar(&webhookConfigPath, "webhookConfigPath", "", "path to webhook config file")
	flag.StringVar(&demoSuitesPath, "config-suites", fmt.Sprintf(rootDir+"/tests/e2e-demos/config/default.yaml"), "path to e2e demo suites definition")
}

func TestE2E(t *testing.T) {
	klog.Info("Starting Red Hat App Studio e2e tests...")
	// Setting viper configurations in cache
	viper.Set("config-suites", demoSuitesPath)

	GenerateE2EDemoTests()

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Red Hat App Studio E2E tests")
}

var _ = ginkgo.SynchronizedAfterSuite(func() {}, func() {
	//Send webhook only it the parameter configPath is not empty
	if len(webhookConfigPath) > 0 {
		klog.Info("Send webhook")
		framework.SendWebhook(webhookConfigPath)
	}
})

var _ = ginkgo.ReportAfterSuite("myReporter", func(report types.Report) {

	// m := make(map[string][]string)
	// for i := range report.SpecReports {
	// 	specReport := report.SpecReports[i]
	// 	if specReport.LeafNodeType == types.NodeTypeIt {
	// 		m[strings.Join(specReport.ContainerHierarchyTexts, ",")] = append(m[strings.Join(specReport.ContainerHierarchyTexts, ",")], specReport.LeafNodeText)
	// 	}
	// }
	// for testCase, testSteps := range m {
	// 	fmt.Printf("Test case: %s\n", testCase)
	// 	for _, testStep := range testSteps {
	// 		fmt.Printf("\t\tStep: %s\n", testStep)
	// 	}
	// }

	m := make(map[string][]string)
	for i := range report.SpecReports {
		specReport := report.SpecReports[i]
		if specReport.LeafNodeType == types.NodeTypeIt {
			m[specReport.ContainerHierarchyTexts[0]] = append(m[specReport.ContainerHierarchyTexts[0]], strings.Join(append(specReport.ContainerHierarchyTexts[1:], specReport.LeafNodeText), ","))
		}
	}
	for testCase, testSteps := range m {
		fmt.Printf("Test case: %s\n", testCase)
		for _, testStep := range testSteps {
			fmt.Printf("\tStep: %s\n", testStep)
		}
	}
})
