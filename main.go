package main

import (
	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
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
				Size:      pulumi.String("g4s.kube.small"),
				NodeCount: pulumi.Int(1),
			},
			Region:       pulumi.StringPtr("FRA1"),
			FirewallId:   firewall.ID(),
			Applications: pulumi.StringPtr("metrics-server"),
		})
		if err != nil {
			return err
		}

		// Create a Kubernetes provider instance using the cluster's kubeconfig
		k8sProvider, err := kubernetes.NewProvider(ctx, "k8s-provider", &kubernetes.ProviderArgs{
			Kubeconfig: cluster.Kubeconfig,
		})
		if err != nil {
			return err
		}

		// Install Argo CD using Helm
		_, err = helm.NewRelease(ctx, "argocd", &helm.ReleaseArgs{
			Chart:     pulumi.String("argo-cd"),
			Version:   pulumi.String("7.8.24"),
			Namespace: pulumi.String("argocd"),
			RepositoryOpts: &helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://argoproj.github.io/argo-helm"),
			},
			CreateNamespace: pulumi.Bool(true),
			Values: pulumi.Map{
				"server": pulumi.Map{
					"service": pulumi.Map{
						"type": pulumi.String("NodePort"),
					},
				},
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return err
		}

		ctx.Export("name", cluster.Name)
		ctx.Export("kubeconfig", cluster.Kubeconfig)
		return nil
	})
}
