package apps

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewStatusfulsetCommand(client *kubernetes.Clientset) *cli.Command {
	return &cli.Command{
		Name:    "statusfulset",
		Aliases: []string{"sts"},
		Subcommands: []*cli.Command{
			{
				Name:  "scale",
				Usage: "to scale statefulset replica",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "namespace",
						Usage:   "namespace for deployment",
						Aliases: []string{"n"},
					},
					&cli.StringFlag{
						Name:     "name",
						Usage:    "statefulset name",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "num",
						Usage:    "replica number",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					return UpdateReplica(ctx, client)
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			return StatusfulsetList(ctx, client)
		},
	}
}

func StatusfulsetList(ctx *cli.Context, client *kubernetes.Clientset) error {
	namespace := ctx.String("namespace")
	sts, err := client.AppsV1().StatefulSets(namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, s := range sts.Items {
		fmt.Println(s.Name)
	}

	return nil
}

func UpdateReplica(ctx *cli.Context, client *kubernetes.Clientset) error {
	namespace := ctx.String("namespace")
	name := ctx.String("name")
	scale, err := client.AppsV1().StatefulSets(namespace).GetScale(context.Background(), name, v1.GetOptions{})
	if err != nil {
		return err
	}

	num := ctx.Int("num")
	if num == int(scale.Spec.Replicas) {
		return nil
	}

	scale.Spec.Replicas = int32(num)
	newScale, err := client.AppsV1().StatefulSets(namespace).UpdateScale(context.Background(), name, scale, v1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("SCALE SUCCESS: %+v\n", newScale)
	return nil
}
