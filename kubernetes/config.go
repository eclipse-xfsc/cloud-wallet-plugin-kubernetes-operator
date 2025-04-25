package kubernetes

import (
	"errors"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	k8rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var client *kubernetes.Clientset

func InitializeKubernetes() error {
	config, err := k8rest.InClusterConfig()
	if viper.GetBool("OVERRIDE_INCLUSTER") || (err != nil && errors.Is(err, k8rest.ErrNotInCluster)) {
		file := viper.GetString("KUBE_FILE")
		masterUrl := viper.GetString("KUBE_CLUSTER_URL")
		config, err = clientcmd.BuildConfigFromFlags(masterUrl, file)
	}

	if err != nil {
		return errors.Join(err, errors.New("error building kubernetes config"))
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return errors.Join(err, errors.New("error creating kubernetes client from config"))
	}

	client = clientset

	return nil
}
