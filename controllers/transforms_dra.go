package controllers

import (
	"strings"

	gpuv1 "github.com/NVIDIA/gpu-operator/api/nvidia/v1"
)

const (
	DRADriverAdditionalNamespacesEnvName = "ADDITIONAL_NAMESPACES"
)

// getDRADriverAdditionalNamespaces returns a list of additional namespaces where
// the DRA driver can manage resources in
func getDRADriverAdditionalNamespaces(config gpuv1.DRADriverSpec) []string {
	envvars := config.ComputeDomains.Controller.Env
	for _, envvar := range envvars {
		if envvar.Name == DRADriverAdditionalNamespacesEnvName && envvar.Value != "" {
			return strings.Split(strings.TrimSpace(envvar.Value), ",")
		}
	}
	return nil
}

func DRADriverServiceAccounts(n ClusterPolicyController) (gpuv1.State, error) {
	status := gpuv1.Ready
	state := n.idx

	draSpec := n.singleton.Spec.DRADriver

	for i, obj := range n.resources[state].ServiceAccounts {
		namespaces := []string{n.operatorNamespace}
		if obj.Name == "compute-domain-daemon" {
			// If the DRA driver is configured to manage multiple namespaces, then
			// we need to create one compute-domain-daemon service account per namespace.
			namespaces = append(namespaces, getDRADriverAdditionalNamespaces(draSpec)...)
		}
		for _, namespace := range namespaces {
			stat, err := createServiceAccount(n, i, namespace)
			if err != nil {
				return stat, err
			}
			if stat == gpuv1.NotReady {
				status = gpuv1.NotReady
			}
		}
	}
	return status, nil
}

func DRADriverRoles(n ClusterPolicyController) (gpuv1.State, error) {
	status := gpuv1.Ready
	state := n.idx

	draSpec := n.singleton.Spec.DRADriver
	namespaces := append([]string{n.operatorNamespace}, getDRADriverAdditionalNamespaces(draSpec)...)

	for i := range n.resources[state].Roles {
		for _, namespace := range namespaces {
			stat, err := createRole(n, i, namespace)
			if err != nil {
				return stat, err
			}
			if stat == gpuv1.NotReady {
				status = gpuv1.NotReady
			}
		}
	}
	return status, nil
}

func DRADriverRoleBindings(n ClusterPolicyController) (gpuv1.State, error) {
	status := gpuv1.Ready
	state := n.idx

	draSpec := n.singleton.Spec.DRADriver
	namespaces := append([]string{n.operatorNamespace}, getDRADriverAdditionalNamespaces(draSpec)...)

	for i, obj := range n.resources[state].RoleBindings {
		for _, namespace := range namespaces {
			serviceAccountNamespace := namespace
			if obj.Name == "nvidia-dra-driver" {
				serviceAccountNamespace = n.operatorNamespace
			}
			stat, err := createRoleBinding(n, i, namespace, serviceAccountNamespace)
			if err != nil {
				return stat, err
			}
			if stat == gpuv1.NotReady {
				status = gpuv1.NotReady
			}
		}
	}
	return status, nil
}

func DRADriverClusterRoleBindings(n ClusterPolicyController) (gpuv1.State, error) {
	status := gpuv1.Ready
	state := n.idx

	draSpec := n.singleton.Spec.DRADriver

	for i, obj := range n.resources[state].ClusterRoleBindings {
		namespaces := []string{n.operatorNamespace}
		appendNamespaceToName := false
		if obj.Name == "compute-domain-daemon" {
			// Multiple service accounts exist for the compute-domain-daemon if the DRA driver
			// is configured to manage multiple namespaces. In this case, we create one ClusterRoleBinding
			// per service account / namespace. Since ClusterRoleBindings are cluster-scoped, we
			// append the namespace to the name of the ClusterRoleBinding resource to avoid naming
			// conflicts.
			namespaces = append(namespaces, getDRADriverAdditionalNamespaces(draSpec)...)
			appendNamespaceToName = true
		}

		for _, namespace := range namespaces {
			stat, err := createClusterRoleBinding(n, i, namespace, appendNamespaceToName)
			if err != nil {
				return stat, err
			}
			if stat == gpuv1.NotReady {
				status = gpuv1.NotReady
			}
		}
	}
	return status, nil
}
