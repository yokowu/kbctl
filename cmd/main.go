package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yokowu/kbctl/apps"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// var kubeconfig string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = filepath.Join(home, ".kube", "config")
	// }

	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	app := cli.App{
		Name: "kbctl",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "namespace",
				Usage:   "namespace for deployment",
				Value:   "default",
				Aliases: []string{"n"},
			},
		},
		Commands: []*cli.Command{
			apps.NewDeploymentCommand(client),
			apps.NewStatusfulsetCommand(client),
			apps.NewNetworkPolicyCommand(client),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
