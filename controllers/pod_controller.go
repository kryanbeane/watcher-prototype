package controllers

import (
	"context"
	"github.com/sirupsen/logrus"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type PodController struct {
	client.Client
	*runtime.Scheme
	dynamic.Interface
}

const (
	serviceAccountNamePrefix = "heimdall-service-account"
	bindingName              = "heimdall-binding"
)

var _ reconcile.Reconciler = &PodController{}

// Add +kubebuilder:rbac:groups=*,resources=*,verbs=get;list;watch
func (p PodController) Add(mgr manager.Manager, selector v1.LabelSelector) error {
	// Create a new Controller
	c, err := controller.New("pod-controller", mgr,
		controller.Options{Reconciler: &PodController{
			Client:    mgr.GetClient(),
			Scheme:    mgr.GetScheme(),
			Interface: dynamic.NewForConfigOrDie(mgr.GetConfig()),
		}})
	if err != nil {
		logrus.Errorf("Failed to create pod controller: %v", err)
		return err
	}

	// Create label selector containing the specified label
	labelSelectorPredicate, err := predicate.LabelSelectorPredicate(selector)
	if err != nil {
		logrus.Errorf("Error creating label selector predicate: %v", err)
		return err
	}

	// Add a watch to objects containing that label
	err = c.Watch(
		&source.Kind{Type: &v12.Pod{}}, &handler.EnqueueRequestForObject{}, labelSelectorPredicate)
	if err != nil {
		logrus.Errorf("Error creating watch for objects: %v", err)
		return err
	}
	return nil
}

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=pods/finalizers,verbs=update

func (p PodController) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logrus.Infof("Reconciling pod %s", request.NamespacedName)

	return reconcile.Result{}, nil
}
