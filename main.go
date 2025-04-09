package main

import (
	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		firewall, err := civo.NewFirewall(ctx, "civo-firewall", &civo.FirewallArgs{
			Name:   pulumi.String("CivoFirewall"),
			Region: pulumi.StringPtr("FRA1"),
			// CreateDefaultRules: pulumi.BoolPtr(true),
			IngressRules: civo.FirewallIngressRuleArray{
				&civo.FirewallIngressRuleArgs{
					Action: pulumi.String("allow"),
					Cidrs: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
					PortRange: pulumi.String("80"),
					Label:     pulumi.String("allow-ingress-http"),
				},
				&civo.FirewallIngressRuleArgs{
					Action: pulumi.String("allow"),
					Cidrs: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
					PortRange: pulumi.String("443"),
				},
				&civo.FirewallIngressRuleArgs{
					Action: pulumi.String("allow"),
					Cidrs: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
					PortRange: pulumi.String("6443"),
					Label:     pulumi.String("allow-ingress-kube-api"),
				},
			},
		})
		if err != nil {
			return err
		}
		cluster, err := civo.NewKubernetesCluster(ctx, "civo-k3s-cluster", &civo.KubernetesClusterArgs{
			Name:              pulumi.StringPtr("CivoK3sCluster"),
			KubernetesVersion: pulumi.String("1.30.5-k3s1"),
			Cni:               pulumi.String("cilium"),
			Pools: civo.KubernetesClusterPoolsArgs{
				Size:      pulumi.String("g4s.kube.xsmall"),
				NodeCount: pulumi.Int(1),
			},
			Region:       pulumi.StringPtr("FRA1"),
			FirewallId:   firewall.ID(),
			Applications: pulumi.StringPtr("traefik2-nodeport,metrics-server,gitlab"),
		})
		if err != nil {
			return err
		}

		ctx.Export("name", cluster.Name)
		ctx.Export("kubeconfig", cluster.Kubeconfig)
		return nil
	})
}
