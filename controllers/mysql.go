package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/mgeri/simple-wordpress-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SimpleWordpressReconciler) deploymentForMysql(logger logr.Logger,
	cr *v1alpha1.SimpleWordpress) *appsv1.Deployment {
	labels := map[string]string{
		"app": cr.Name,
	}
	var matchLabels = map[string]string{
		"app":  cr.Name,
		"tier": "mysql",
	}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wordpress-mysql",
			Namespace: cr.Namespace,
			Labels:    labels,
		},

		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: matchLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "mysql:5.6",
						Name:  "mysql",

						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: cr.Spec.SqlRootPassword,
							},
						},

						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
					},
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, dep, r.Scheme)
	return dep
}

func (r *SimpleWordpressReconciler) serviceForMysql(logger logr.Logger,
	cr *v1alpha1.SimpleWordpress) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}
	matchLabels := map[string]string{
		"app":  cr.Name,
		"tier": "mysql",
	}

	ser := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "wordpress-mysql",
			Namespace: cr.Namespace,
			Labels:    labels,
		},

		Spec: corev1.ServiceSpec{
			Selector: matchLabels,

			Ports: []corev1.ServicePort{
				{
					Port: 3306,
					Name: cr.Name,
				},
			},
			ClusterIP: "None",
		},
	}

	controllerutil.SetControllerReference(cr, ser, r.Scheme)
	return ser
}

func (r *SimpleWordpressReconciler) isMysqlUp(logger logr.Logger,
	v *v1alpha1.SimpleWordpress) bool {
	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      "wordpress-mysql",
		Namespace: v.Namespace,
	}, deployment)

	if err != nil {
		logger.Error(err, "Deployment mysql not found")
		return false
	}
	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false

}
