package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DPUSetSpec defines the desired state of DPUSet
type DPUSetSpec struct {
	// DpuNodeSelector selects the host nodes running DPUs
	DpuNodeSelector map[string]string `json:"dpuNodeSelector"`

	// BFB defines the bootstream image version or URL
	BFB string `json:"bfb"`

	// Flavor specifies configuration presets for DPUs
	Flavor string `json:"flavor,omitempty"`
}

// DPUSetStatus defines the observed state of DPUSet
type DPUSetStatus struct {
	// Ready indicates if the DPU provisioning has succeeded
	Ready bool `json:"ready"`

	// Phase represents the current stage of provisioning (e.g., Flashing, Ready, Failed)
	Phase string `json:"phase,omitempty"`

	// Message describes the status details
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DPUSet is the Schema for the dpuseets API (NVIDIA DPF CRD)
type DPUSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DPUSetSpec   `json:"spec,omitempty"`
	Status DPUSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DPUSetList contains a list of DPUSet
type DPUSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DPUSet `json:"items"`
}

// DPUServiceSpec defines the desired state of DPUService
type DPUServiceSpec struct {
	// ServiceType defines what kind of service is running (e.g. "ovn", "storage")
	ServiceType string `json:"serviceType"`

	// Config contains the service-specific configuration string or JSON
	Config string `json:"config,omitempty"`

	// DPUSetName associates this service with a DPUSet
	DPUSetName string `json:"dpuSetName"`
}

// DPUServiceStatus defines the observed state of DPUService
type DPUServiceStatus struct {
	// Ready indicates if the service is deployed and active on DPUs
	Ready bool `json:"ready"`

	// ActiveServices tracks active service daemon instances
	ActiveServices int `json:"activeServices"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DPUService is the Schema for the dpuservices API (NVIDIA DPF CRD)
type DPUService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DPUServiceSpec   `json:"spec,omitempty"`
	Status DPUServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DPUServiceList contains a list of DPUService
type DPUServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DPUService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DPUSet{}, &DPUSetList{}, &DPUService{}, &DPUServiceList{})
}

// DeepCopyInto copies the receiver, writing into out. in must be non-nil.
func (in *DPUSet) DeepCopyInto(out *DPUSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy creates a new DPUSet.
func (in *DPUSet) DeepCopy() *DPUSet {
	if in == nil {
		return nil
	}
	out := new(DPUSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject creates a new runtime.Object.
func (in *DPUSet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto copies the receiver, writing into out. in must be non-nil.
func (in *DPUSetSpec) DeepCopyInto(out *DPUSetSpec) {
	*out = *in
	if in.DpuNodeSelector != nil {
		in, out := &in.DpuNodeSelector, &out.DpuNodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopyInto copies the receiver, writing into out. in must be non-nil.
func (in *DPUSetList) DeepCopyInto(out *DPUSetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DPUSet, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy creates a new DPUSetList.
func (in *DPUSetList) DeepCopy() *DPUSetList {
	if in == nil {
		return nil
	}
	out := new(DPUSetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject creates a new runtime.Object.
func (in *DPUSetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto copies the receiver, writing into out. in must be non-nil.
func (in *DPUService) DeepCopyInto(out *DPUService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy creates a new DPUService.
func (in *DPUService) DeepCopy() *DPUService {
	if in == nil {
		return nil
	}
	out := new(DPUService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject creates a new runtime.Object.
func (in *DPUService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto copies the receiver, writing into out. in must be non-nil.
func (in *DPUServiceList) DeepCopyInto(out *DPUServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DPUService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy creates a new DPUServiceList.
func (in *DPUServiceList) DeepCopy() *DPUServiceList {
	if in == nil {
		return nil
	}
	out := new(DPUServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject creates a new runtime.Object.
func (in *DPUServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
