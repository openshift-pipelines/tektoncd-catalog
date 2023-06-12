package controller

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/openshift-pipelines/tektoncd-catalog/acceptance-tests/steps"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func BeforeSuite() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	nodes, _ := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	for _, node := range nodes.Items {
		fmt.Printf("%s\n", node.Name)
		for _, condition := range node.Status.Conditions {
			fmt.Printf("\t%s: %s\n", condition.Type, condition.Status)
		}
	}
}

func AfterSuite() {
	// Add Cleanup logic!
}
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the Kubernetes cluster is available$`, steps.TheKubernetesClusterIsAvailable)
	ctx.Step(`^"([^"]*)" is deployed and "([^"]*)" on "([^"]*)" namespace$`, steps.IsDeployedAndOnNamespace)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(BeforeSuite)
	ctx.AfterSuite(AfterSuite)
}
