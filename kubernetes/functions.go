package kubernetes

import (
	"bytes"
	"fmt"
	"github.com/developer-guy/easy-rollback/client"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

var k8sClient *kubernetes.Clientset

func init() {
	k8sClient = client.NewK8sClient()
}

type CommandFunction func(cmd *cobra.Command, args []string)

func findCurrentReplicaSetOfDeployment(replicaSets *v1.ReplicaSetList) v1.ReplicaSet {
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

func findOtherDeployedImages(replicaSets *v1.ReplicaSetList) *hashmap.Map {
	images := hashmap.New()
	for _, replicaSet := range replicaSets.Items {
		images.Put(replicaSet.Spec.Template.Spec.Containers[0].Image,
			replicaSet.CreationTimestamp.Format("02 January 2006 15:04:05"))
	}
	return images
}

func ListPreviousDeployedImages() CommandFunction {
	return func(cmd *cobra.Command, args []string) {
		deploymentFlag := cmd.Flag("deployment").Value.String()
		namespaceFlag := cmd.Flag("namespace").Value.String()

		deploymentsClient := k8sClient.AppsV1().Deployments(namespaceFlag)

		deployment, err := deploymentsClient.Get(deploymentFlag, metav1.GetOptions{})

		if err != nil {
			panic(err)
		}

		matchLabels := deployment.Spec.Selector.MatchLabels
		var labels bytes.Buffer

		for key, value := range matchLabels {
			labels.WriteString(key + "=" + value)
		}

		replicaSets, _ := k8sClient.AppsV1().ReplicaSets(namespaceFlag).List(metav1.ListOptions{
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
				fmt.Println(fmt.Sprintf("Image version: %s , Creation creationTime: %s ", image, creationTime.(string), ))
			}
		}
	}
}

func RollbackDeployment() CommandFunction {
	return func(cmd *cobra.Command, args []string) {
		deploymentFlag := cmd.Flag("deployment").Value.String()
		namespaceFlag := cmd.Flag("namespace").Value.String()
		toImageFlag := cmd.Flag("to-image").Value.String()

		deploymentsClient := k8sClient.AppsV1().Deployments(namespaceFlag)

		deployment, err := deploymentsClient.Get(deploymentFlag, metav1.GetOptions{})

		if err != nil {
			panic(err)
		}

		deployment.Spec.Template.Spec.Containers[0].Image = toImageFlag
		_, deploymentUpdateStatus := deploymentsClient.Update(deployment)

		if deploymentUpdateStatus != nil {
			panic(err)
		} else {
			fmt.Println("Successfully rollbacked to image", toImageFlag)
		}

	}
}
