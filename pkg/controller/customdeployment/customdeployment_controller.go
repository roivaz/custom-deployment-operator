package customdeployment

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	customdeploymentv1alpha1 "github.com/lominorama/custom-deployment-operator/pkg/apis/customdeployment/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CustomDeployment Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCustomDeployment{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("customdeployment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CustomDeployment
	err = c.Watch(&source.Kind{Type: &customdeploymentv1alpha1.CustomDeployment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner CustomDeployment
	// This is required so the controller can detect changes on the created deployments
	// and enforce the desired state on them. For example if someone manually changes something in
	// the deployment, this will notice it and trigger the reconciliation loop
	// err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &customdeploymentv1alpha1.CustomDeployment{},
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

var _ reconcile.Reconciler = &ReconcileCustomDeployment{}

// ReconcileCustomDeployment reconciles a CustomDeployment object
type ReconcileCustomDeployment struct {
	// TODO: Clarify the split client
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CustomDeployment object and makes changes based on the state read
// and what is in the CustomDeployment.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCustomDeployment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("Reconciling CustomDeployment %s/%s\n", request.Namespace, request.Name)

	// Fetch the CustomDeployment instance
	instance := &customdeploymentv1alpha1.CustomDeployment{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Deployment object
	deployment := newDeploymentForCR(instance)

	// Set CustomDeployment instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Deployment already exists

	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Printf("Creating a new Deployment %s/%s\n", deployment.Namespace, deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			log.Print(err)
			return reconcile.Result{}, err
		}

		// Deployment created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		log.Print(err)
		return reconcile.Result{}, err
	}

	// Deployment already exists - update
	log.Printf("Updating deployment %s/%s", deployment.Namespace, deployment.Name)
	err = r.client.Update(context.TODO(), deployment)
	if err != nil {
		log.Print(err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newDeploymentForCR(cr *customdeploymentv1alpha1.CustomDeployment) *appsv1.Deployment {

	var yaml = `
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: placeholder
      namespace: placeholder
      labels:
        k8s-app: placeholder
    spec:
      selector:
        matchLabels:
          k8s-app: placeholder
      replicas: 4
      template:
        metadata:
          labels:
            k8s-app: placeholder
        spec:
          containers:
          - name: placeholder
            image: placeholder
            ports:
            - containerPort: 80
            resources:
              requests:
                cpu: 100m
    `

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(yaml), nil, nil)
	if err != nil {
		log.Printf("%#v", err)
	}
	deployment := obj.(*appsv1.Deployment)
	spew.Dump(deployment)

	labels := map[string]string{
		"k8s-app": cr.Name,
	}

	// Set customizable fields
	deployment.ObjectMeta.Name = cr.Name + "-deployment"
	deployment.ObjectMeta.Namespace = cr.Namespace
	deployment.ObjectMeta.Labels = labels
	deployment.Spec.Selector.MatchLabels = labels
	deployment.Spec.Template.ObjectMeta.Labels = labels

	// TODO: check that deployment definition only contains one container
	deployment.Spec.Template.Spec.Containers[0].Name = cr.Name + "-pod"
	deployment.Spec.Template.Spec.Containers[0].Image = cr.Spec.Image + ":" + cr.Spec.Version

	return deployment
}
