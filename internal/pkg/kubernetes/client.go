package kubernetes

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GenerateClient(kubeconfig string) (*kubernetes.Clientset, error) {
	var cfg *rest.Config
	var err error

	cfg, err = rest.InClusterConfig()
	if err != nil && err != rest.ErrNotInCluster {
		return nil, fmt.Errorf("failed to load Kubernetes configuration from cluster: %v", err)
	}

	if cfg == nil {
		if kubeconfig == "" {
			kubeconfig = clientcmd.RecommendedHomeFile
		}

		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load Kubernetes configuration: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	return clientset, nil
}

func RestartDeployment(clientset *kubernetes.Clientset, namespace string, deploymentName string) error {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %v", err)
	}

	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}

	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	_, err = clientset.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %v", err)
	}

	return nil
}
