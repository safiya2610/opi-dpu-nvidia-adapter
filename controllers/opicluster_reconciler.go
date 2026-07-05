package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	opiv1alpha1 "opi-nvidia-adapter/api/v1alpha1"
	"opi-nvidia-adapter/pkg/adapter"
)

// DPUClusterReconciler reconciles a DPUCluster object
type DPUClusterReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=opi.io,resources=dpuclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=opi.io,resources=dpuclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=opi.io,resources=dpuclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=opi.io,resources=dpuseets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=opi.io,resources=dpuseets/status,verbs=get
// +kubebuilder:rbac:groups=opi.io,resources=dpuservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=opi.io,resources=dpuservices/status,verbs=get

func (r *DPUClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("dpucluster", req.NamespacedName)

	// 1. Fetch the DPUCluster instance
	var cluster opiv1alpha1.DPUCluster
	if err := r.Get(ctx, req.NamespacedName, &cluster); err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get DPUCluster")
		return ctrl.Result{}, err
	}

	// 2. Check vendor - if not nvidia, skip it
	if cluster.Spec.Vendor != "nvidia" {
		log.Info("Skipping DPUCluster: vendor is not nvidia", "vendor", cluster.Spec.Vendor)
		return ctrl.Result{}, nil
	}

	log.Info("Reconciling NVIDIA DPUCluster")

	// 3. Reconcile DPUSet
	desiredDPUSet := adapter.TranslateDPUClusterToDPUSet(&cluster)
	if err := controllerutil.SetControllerReference(&cluster, desiredDPUSet, r.Scheme); err != nil {
		log.Error(err, "Failed to set controller reference on DPUSet")
		return ctrl.Result{}, err
	}

	var existingDPUSet opiv1alpha1.DPUSet
	err := r.Get(ctx, client.ObjectKey{Namespace: desiredDPUSet.Namespace, Name: desiredDPUSet.Name}, &existingDPUSet)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Creating a new DPUSet", "Namespace", desiredDPUSet.Namespace, "Name", desiredDPUSet.Name)
			if err := r.Create(ctx, desiredDPUSet); err != nil {
				log.Error(err, "Failed to create DPUSet")
				r.Recorder.Event(&cluster, "Warning", "CreationFailed", fmt.Sprintf("Failed to create DPUSet %s", desiredDPUSet.Name))
				return ctrl.Result{}, err
			}
			r.Recorder.Event(&cluster, "Normal", "Created", fmt.Sprintf("Created DPUSet %s", desiredDPUSet.Name))
		} else {
			log.Error(err, "Failed to get DPUSet")
			return ctrl.Result{}, err
		}
	} else {
		// Update existing DPUSet if specs mismatch
		if existingDPUSet.Spec.BFB != desiredDPUSet.Spec.BFB || existingDPUSet.Spec.Flavor != desiredDPUSet.Spec.Flavor {
			existingDPUSet.Spec.BFB = desiredDPUSet.Spec.BFB
			existingDPUSet.Spec.Flavor = desiredDPUSet.Spec.Flavor
			existingDPUSet.Spec.DpuNodeSelector = desiredDPUSet.Spec.DpuNodeSelector
			log.Info("Updating existing DPUSet", "Namespace", existingDPUSet.Namespace, "Name", existingDPUSet.Name)
			if err := r.Update(ctx, &existingDPUSet); err != nil {
				log.Error(err, "Failed to update DPUSet")
				return ctrl.Result{}, err
			}
		}
	}

	// 4. Reconcile DPUService
	var existingDPUService opiv1alpha1.DPUService
	serviceExists := false
	desiredDPUService := adapter.TranslateDPUClusterToDPUService(&cluster, desiredDPUSet.Name)
	if desiredDPUService != nil {
		if err := controllerutil.SetControllerReference(&cluster, desiredDPUService, r.Scheme); err != nil {
			log.Error(err, "Failed to set controller reference on DPUService")
			return ctrl.Result{}, err
		}

		err = r.Get(ctx, client.ObjectKey{Namespace: desiredDPUService.Namespace, Name: desiredDPUService.Name}, &existingDPUService)
		if err != nil {
			if apierrors.IsNotFound(err) {
				log.Info("Creating a new DPUService", "Namespace", desiredDPUService.Namespace, "Name", desiredDPUService.Name)
				if err := r.Create(ctx, desiredDPUService); err != nil {
					log.Error(err, "Failed to create DPUService")
					r.Recorder.Event(&cluster, "Warning", "CreationFailed", fmt.Sprintf("Failed to create DPUService %s", desiredDPUService.Name))
					return ctrl.Result{}, err
				}
				r.Recorder.Event(&cluster, "Normal", "Created", fmt.Sprintf("Created DPUService %s", desiredDPUService.Name))
			} else {
				log.Error(err, "Failed to get DPUService")
				return ctrl.Result{}, err
			}
		} else {
			serviceExists = true
			// Update existing DPUService if spec mismatch
			if existingDPUService.Spec.Config != desiredDPUService.Spec.Config {
				existingDPUService.Spec.Config = desiredDPUService.Spec.Config
				log.Info("Updating existing DPUService", "Namespace", existingDPUService.Namespace, "Name", existingDPUService.Name)
				if err := r.Update(ctx, &existingDPUService); err != nil {
					log.Error(err, "Failed to update DPUService")
					return ctrl.Result{}, err
				}
			}
		}
	}

	// 5. Check status of underlying resources and propagate to DPUCluster
	// Fetch latest DPUSet to get the current status
	err = r.Get(ctx, client.ObjectKey{Namespace: desiredDPUSet.Namespace, Name: desiredDPUSet.Name}, &existingDPUSet)
	if err != nil {
		return ctrl.Result{}, err
	}

	dpusetReady := existingDPUSet.Status.Ready
	dpusetPhase := existingDPUSet.Status.Phase

	serviceReady := true
	if desiredDPUService != nil {
		err = r.Get(ctx, client.ObjectKey{Namespace: desiredDPUService.Namespace, Name: desiredDPUService.Name}, &existingDPUService)
		if err != nil {
			return ctrl.Result{}, err
		}
		serviceReady = existingDPUService.Status.Ready
	} else if serviceExists {
		// If desired service is nil but it exists, it means offloading was disabled/removed.
		// Delete the resource
		log.Info("Deleting old DPUService", "Namespace", existingDPUService.Namespace, "Name", existingDPUService.Name)
		if err := r.Delete(ctx, &existingDPUService); err != nil {
			log.Error(err, "Failed to delete DPUService")
			return ctrl.Result{}, err
		}
		serviceReady = true
	}

	// 6. Propagate Phase and Ready conditions
	var newPhase string
	var newReady bool
	var msg string

	if !dpusetReady {
		newReady = false
		if dpusetPhase != "" {
			newPhase = dpusetPhase
		} else {
			newPhase = "DPUProvisioning"
		}
		msg = fmt.Sprintf("Waiting for DPUSet: %s", existingDPUSet.Status.Message)
	} else if !serviceReady {
		newReady = false
		newPhase = "ServiceConfiguring"
		msg = "DPUs provisioned. Waiting for network service config."
	} else {
		newReady = true
		newPhase = "Ready"
		msg = "All DPUs provisioned and services deployed."
	}

	// Apply status updates
	if cluster.Status.Ready != newReady || cluster.Status.Phase != newPhase {
		cluster.Status.Ready = newReady
		cluster.Status.Phase = newPhase

		// Set Condition
		cond := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             newPhase,
			Message:            msg,
		}
		if newReady {
			cond.Status = metav1.ConditionTrue
		}

		// Basic condition helper: replace or append
		found := false
		for i, c := range cluster.Status.Conditions {
			if c.Type == "Ready" {
				cluster.Status.Conditions[i] = cond
				found = true
				break
			}
		}
		if !found {
			cluster.Status.Conditions = append(cluster.Status.Conditions, cond)
		}

		log.Info("Updating DPUCluster status", "Phase", newPhase, "Ready", newReady)
		if err := r.Status().Update(ctx, &cluster); err != nil {
			log.Error(err, "Failed to update DPUCluster status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DPUClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&opiv1alpha1.DPUCluster{}).
		Owns(&opiv1alpha1.DPUSet{}).
		Owns(&opiv1alpha1.DPUService{}).
		Complete(r)
}
