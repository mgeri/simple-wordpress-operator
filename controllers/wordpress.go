package controllers

import (
	"github.com/go-logr/logr"
	"github.com/mgeri/simple-wordpress-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SimpleWordpressReconciler) deploymentForWordpress(logger logr.Logger, cr *v1alpha1.SimpleWordpress) *appsv1.Deployment {

	labels := map[string]string{
		"app": cr.Name,
	}
	matchLabels := map[string]string{
		"app":  cr.Name,
		"tier": "frontend",
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wordpress",
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
						Image: "wordpress:5.4-apache",
						Name:  "wordpress",

						Env: []corev1.EnvVar{{
							Name:  "WORDPRESS_DB_HOST",
							Value: "wordpress-mysql",
						},
							{
								Name:  "WORDPRESS_DB_PASSWORD",
								Value: cr.Spec.SqlRootPassword,
							},
						},

						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "wordpress-port",
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

func (r *SimpleWordpressReconciler) serviceForWordpress(logger logr.Logger, cr *v1alpha1.SimpleWordpress) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}
	matchLabels := map[string]string{
		"app":  cr.Name,
		"tier": "frontend",
	}

	ser := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "wordpress",
			Namespace: cr.Namespace,
			Labels:    labels,
		},

		Spec: corev1.ServiceSpec{
			Selector: matchLabels,

			Ports: []corev1.ServicePort{
				{
					Port: 80,
					Name: "port",
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(cr, ser, r.Scheme)
	return ser

}
