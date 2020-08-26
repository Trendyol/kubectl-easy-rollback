package client

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	appsV1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
		panic(err.Error())
	}
	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return &K8SClient{client}
}

type CmdFunc func(cmd *cobra.Command, args []string)

func findCurrentReplicaSetOfDeployment(replicaSets *appsV1.ReplicaSetList) appsV1.ReplicaSet {
	for _, replicaSet := range replicaSets.Items {
		if *replicaSet.Spec.Replicas != 0 {
			return replicaSet
		}
	}
	panic("No active replicaSet found for given deployment")
}

type Image struct {
	name     string
	creation time.Time
}

func findOtherDeployedImages(replicaSets *appsV1.ReplicaSetList) *hashmap.Map {
	images := hashmap.New()
	for _, replicaSet := range replicaSets.Items {
		images.Put(replicaSet.Spec.Template.Spec.Containers[0].Image,
			replicaSet.CreationTimestamp.Format("02 January 2006 15:04:05"))
	}
	return images
}

func (k *K8SClient) ListPreviousDeployedImages(deployment, namespace string) {
	deploymentsClient := k.AppsV1().Deployments(namespace)

	deploy, err := deploymentsClient.Get(context.Background(), deployment, metaV1.GetOptions{})

	if err != nil {
		panic(err)
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

	currentReplicaSet := findCurrentReplicaSetOfDeployment(replicaSets)

	deployedImages := findOtherDeployedImages(replicaSets)

	for _, image := range deployedImages.Keys() {
		creationTime, _ := deployedImages.Get(image)
		if strings.Compare(currentReplicaSet.Spec.Template.Spec.Containers[0].Image, image.(string)) == 0 {
			fmt.Println(fmt.Sprintf("Image version: %s , Creation creationTime: %s %s", image, creationTime.(string),
				chalk.Green.Color("*")))
		} else {
			fmt.Println(fmt.Sprintf("Image version: %s , Creation creationTime: %s ", image, creationTime.(string)))
		}
	}
}

func (k *K8SClient) RollbackDeployment(namespace, deployment, toImage string) {

	deploymentsClient := k.AppsV1().Deployments(namespace)

	deploy, err := deploymentsClient.Get(context.Background(), deployment, metaV1.GetOptions{})

	if err != nil {
		panic(err)
	}

	deploy.Spec.Template.Spec.Containers[0].Image = toImage
	_, deploymentUpdateStatus := deploymentsClient.Update(context.Background(), deploy, metaV1.UpdateOptions{})

	if deploymentUpdateStatus != nil {
		panic(err)
	} else {
		fmt.Println("Successfully rollbacked to image", toImage)
	}
}
