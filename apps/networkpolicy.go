package apps

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
	netv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func NewNetworkPolicyCommand(client *kubernetes.Clientset) *cli.Command {
	return &cli.Command{
		Name:    "networkpolicy",
		Aliases: []string{"np"},
		Action: func(c *cli.Context) error {
			return NetworkPolicyList(c, client)
		},
		Subcommands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create new networkpolicy",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "ports",
						Usage:   "ingress ports format 80,8080",
						Aliases: []string{"p"},
					},
				},
				Action: func(c *cli.Context) error {
					return CreateNetworkPolicy(c, client)
				},
			},
			{
				Name:  "delete",
				Usage: "delete networkpolicy by name",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Usage:   "networkpolicy name",
						Aliases: []string{"n"},
					},
				},
				Action: func(c *cli.Context) error {
					return DeleteNetworkPolicy(c, client)
				},
			},
		},
	}
}

func NetworkPolicyList(ctx *cli.Context, client *kubernetes.Clientset) error {
	namespace := ctx.String("namespace")
	np, err := client.NetworkingV1().NetworkPolicies(namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("NAMESPACE\tNAME\n")
	for _, n := range np.Items {
		fmt.Printf("%s \t%s\n", n.Namespace, n.Name)
	}

	return nil
}

func CreateNetworkPolicy(ctx *cli.Context, client *kubernetes.Clientset) error {
	namespace := ctx.String("namespace")
	npi := client.NetworkingV1().NetworkPolicies(namespace)
	np := &netv1.NetworkPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name:      "net-test",
			Namespace: namespace,
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: v1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "gin-test",
				},
			},
		},
	}

	ports := ctx.String("ports")
	var netPorts []netv1.NetworkPolicyPort
	if ports != "" {
		portList := strings.Split(ports, ",")
		for _, port := range portList {
			if port == "" {
				continue
			}

			portInt, err := strconv.Atoi(port)
			if err != nil {
				return err
			}
			netPorts = append(netPorts, netv1.NetworkPolicyPort{
				Port: &intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: int32(portInt),
				},
			})
		}
	}

	rule := netv1.NetworkPolicyIngressRule{}

	if len(netPorts) > 0 {
		rule.Ports = netPorts
	}

	np.Spec.Ingress = append(np.Spec.Ingress, rule)

	np, err := npi.Create(context.Background(), np, v1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Printf("create network success\n\n %+v", np)
	return nil
}

func DeleteNetworkPolicy(ctx *cli.Context, client *kubernetes.Clientset) error {
	namespace := ctx.String("namespace")
	npi := client.NetworkingV1().NetworkPolicies(namespace)

	for _, n := range ctx.Args().Slice() {
		if err := npi.Delete(context.Background(), n, v1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}
