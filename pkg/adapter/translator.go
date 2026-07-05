package adapter

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	opiv1alpha1 "opi-nvidia-adapter/api/v1alpha1"
)

// TranslateDPUClusterToDPUSet maps an OPI DPUCluster to a DPF DPUSet
func TranslateDPUClusterToDPUSet(cluster *opiv1alpha1.DPUCluster) *opiv1alpha1.DPUSet {
	if cluster == nil {
		return nil
	}

	dpuset := &opiv1alpha1.DPUSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-dpuset", cluster.Name),
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"opi.io/managed-by": "opi-nvidia-adapter",
				"opi.io/cluster":    cluster.Name,
			},
		},
		Spec: opiv1alpha1.DPUSetSpec{
			DpuNodeSelector: cluster.Spec.NodeSelector,
			BFB:             cluster.Spec.BFB,
			Flavor:          cluster.Spec.DpuFlavor,
		},
	}

	// Set default flavor if not specified
	if dpuset.Spec.Flavor == "" {
		dpuset.Spec.Flavor = "default"
	}

	return dpuset
}

// TranslateDPUClusterToDPUService maps an OPI DPUCluster to a DPF DPUService
func TranslateDPUClusterToDPUService(cluster *opiv1alpha1.DPUCluster, dpuSetName string) *opiv1alpha1.DPUService {
	if cluster == nil || cluster.Spec.NetworkOffloadMode == "" {
		return nil
	}

	// We map network services if the offload mode is specified.
	serviceType := "network"
	configStr := fmt.Sprintf(`{"mode": "%s", "vpc": "%s"}`, cluster.Spec.NetworkOffloadMode, cluster.Spec.VpcName)

	dpuservice := &opiv1alpha1.DPUService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-dpuservice", cluster.Name),
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"opi.io/managed-by": "opi-nvidia-adapter",
				"opi.io/cluster":    cluster.Name,
			},
		},
		Spec: opiv1alpha1.DPUServiceSpec{
			ServiceType: serviceType,
			Config:      configStr,
			DPUSetName:  dpuSetName,
		},
	}

	return dpuservice
}
