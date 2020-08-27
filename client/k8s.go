package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mitchellh/go-homedir"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

type K8SClient struct {
	*kubernetes.Clientset
}

func NewK8sClient(kubeconfig string, context string) *K8SClient {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err.Error())
	}

	if kubeconfig == "" {
		kubeconfig = filepath.Join(homeDir, ".kube", "config")
	}

	// TODO: Support all kubectl global flags
	overrides := clientcmd.ConfigOverrides{}

	if context != "" {
		overrides.CurrentContext = context
	}

	deferredLoadingClientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&overrides)

	config, err := deferredLoadingClientConfig.ClientConfig()
	if err != nil {
		log.Fatalf("detail: %+v", err)
	}
	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("detail: %+v", err)
	}
	return &K8SClient{client}
}

func (k *K8SClient) ListPreviousDeployedImages(deployment, namespace string) {
	deploymentsClient := k.AppsV1().Deployments(namespace)

	deploy, err := deploymentsClient.Get(context.Background(), deployment, metaV1.GetOptions{})

	if err != nil {
		log.Fatalf("detail: %+v", err)
	}

	matchLabels := deploy.Spec.Selector.MatchLabels
	var labels bytes.Buffer

	l := len(matchLabels)
	c := 0

	for key, value := range matchLabels {
		c += 1
		labels.WriteString(key + "=" + value)
		if c != l {
			labels.WriteString(",")
		}
	}

	replicaSets, _ := k.AppsV1().ReplicaSets(namespace).List(context.Background(), metaV1.ListOptions{
		LabelSelector: labels.String(),
	})

	gtp, err := printers.NewGoTemplatePrinter([]byte("{{range .items}}{{range .spec.template.spec.containers}}{{printf \"Image: %s name: %s\\n\"  .image .name}}{{end}}{{end}}"))

	if err != nil {
		log.Fatalf("detail: %+v", err)
	}

	err = gtp.PrintObj(replicaSets, os.Stdout)

	if err != nil {
		log.Fatalf("detail: %+v", err)
	}
}

func (k *K8SClient) RollbackDeployment(namespace, deployment, toImage string) {

	deploymentsClient := k.AppsV1().Deployments(namespace)

	deploy, err := deploymentsClient.Get(context.Background(), deployment, metaV1.GetOptions{})

	if err != nil {
		log.Fatalf("detail: %+v", err)
	}

	deploy.Spec.Template.Spec.Containers[0].Image = toImage
	_, deploymentUpdateStatus := deploymentsClient.Update(context.Background(), deploy, metaV1.UpdateOptions{})

	if deploymentUpdateStatus != nil {
		log.Fatalf("detail: %+v", err)
	} else {
		fmt.Println("Successfully rollbacked to image", toImage)
	}
}
