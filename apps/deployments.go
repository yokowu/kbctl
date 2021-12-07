package apps

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewDeploymentCommand(client *kubernetes.Clientset) *cli.Command {
	return &cli.Command{
		Name: "deploy",
		Action: func(c *cli.Context) error {
			return DeployList(c, client)
		},
	}
}

func DeployList(ctx *cli.Context, client *kubernetes.Clientset) error {
	namespace := ctx.String("namespace")
	ds, err := client.AppsV1().Deployments(namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, d := range ds.Items {
		fmt.Println(d.Name)
	}

	return nil
}
