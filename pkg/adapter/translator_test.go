package adapter

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	opiv1alpha1 "opi-nvidia-adapter/api/v1alpha1"
)

func TestTranslateDPUClusterToDPUSet(t *testing.T) {
	cluster := &opiv1alpha1.DPUCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "default",
		},
		Spec: opiv1alpha1.DPUClusterSpec{
			Vendor:    "nvidia",
			BFB:       "http://example.com/bfb-v3.bfb",
			DpuFlavor: "high-perf",
			NodeSelector: map[string]string{
				"dpu-node": "true",
			},
		},
	}

	dpuset := TranslateDPUClusterToDPUSet(cluster)
	if dpuset == nil {
		t.Fatalf("expected dpuset to be non-nil")
	}

	if dpuset.Name != "test-cluster-dpuset" {
		t.Errorf("expected name to be test-cluster-dpuset, got %s", dpuset.Name)
	}

	if dpuset.Namespace != "default" {
		t.Errorf("expected namespace default, got %s", dpuset.Namespace)
	}

	if dpuset.Spec.BFB != "http://example.com/bfb-v3.bfb" {
		t.Errorf("expected BFB image to match, got %s", dpuset.Spec.BFB)
	}

	if dpuset.Spec.Flavor != "high-perf" {
		t.Errorf("expected flavor high-perf, got %s", dpuset.Spec.Flavor)
	}

	if dpuset.Spec.DpuNodeSelector["dpu-node"] != "true" {
		t.Errorf("expected node selector to map correctly")
	}
}

func TestTranslateDPUClusterToDPUSet_DefaultFlavor(t *testing.T) {
	cluster := &opiv1alpha1.DPUCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		Spec: opiv1alpha1.DPUClusterSpec{
			Vendor: "nvidia",
		},
	}

	dpuset := TranslateDPUClusterToDPUSet(cluster)
	if dpuset.Spec.Flavor != "default" {
		t.Errorf("expected default flavor, got %s", dpuset.Spec.Flavor)
	}
}

func TestTranslateDPUClusterToDPUService(t *testing.T) {
	cluster := &opiv1alpha1.DPUCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "default",
		},
		Spec: opiv1alpha1.DPUClusterSpec{
			Vendor:             "nvidia",
			NetworkOffloadMode: "ovn-kubernetes",
			VpcName:            "prod-vpc",
		},
	}

	service := TranslateDPUClusterToDPUService(cluster, "test-cluster-dpuset")
	if service == nil {
		t.Fatalf("expected service to be non-nil")
	}

	if service.Name != "test-cluster-dpuservice" {
		t.Errorf("expected name to be test-cluster-dpuservice, got %s", service.Name)
	}

	if service.Spec.ServiceType != "network" {
		t.Errorf("expected service type to be network, got %s", service.Spec.ServiceType)
	}

	expectedConfig := `{"mode": "ovn-kubernetes", "vpc": "prod-vpc"}`
	if service.Spec.Config != expectedConfig {
		t.Errorf("expected config to be %s, got %s", expectedConfig, service.Spec.Config)
	}

	if service.Spec.DPUSetName != "test-cluster-dpuset" {
		t.Errorf("expected dpuset association, got %s", service.Spec.DPUSetName)
	}
}
