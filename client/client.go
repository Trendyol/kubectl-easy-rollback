package client

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewK8sClient() *kubernetes.Clientset {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err.Error())
	}
	var kubeConfigPath = filepath.Join(homeDir, ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		panic(err.Error())
	}

	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return client
}
