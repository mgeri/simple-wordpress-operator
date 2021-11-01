package controllers

import (
	"context"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"

	"github.com/mgeri/simple-wordpress-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *SimpleWordpressReconciler) ensureDeployment(logger logr.Logger,
	request reconcile.Request,
	instance *v1alpha1.SimpleWordpress,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	found := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)

	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		logger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			logger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		}
		// Deployment was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the deployment not existing
		logger.Error(err, "Failed to get Deployment")
		return &ctrl.Result{}, err
	}

	return nil, nil

}

func (r *SimpleWordpressReconciler) ensureService(logger logr.Logger,
	request reconcile.Request,
	instance *v1alpha1.SimpleWordpress,
	s *corev1.Service,
) (*reconcile.Result, error) {
	found := &corev1.Service{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		logger.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.Client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			logger.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &ctrl.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the service not existing
		logger.Error(err, "Failed to get Service")
		return &ctrl.Result{}, err
	}

	return nil, nil
}
