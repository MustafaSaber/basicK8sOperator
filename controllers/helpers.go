package controllers

import (
	msaberdevv1 "custom-k8s-operator/api/v1"
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"

	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *OnekindReconciler) desiredConfigMap(kind msaberdevv1.Onekind) (*corev1.ConfigMap, error) {

	configMapData := make(map[string]string)

	nginxconfig := fmt.Sprintf(
		`
			events {}
			http {
			server {
				location / {
						add_header Content-Type text/plain;
						return 200 "%s";
					}
				}
			}
		`, kind.Spec.Message)

	configMapData["nginx.conf"] = nginxconfig

	config := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kind.Name,
			Namespace: kind.Namespace,
		},
		Data: configMapData,
	}
	if err := ctrl.SetControllerReference(&kind, config, r.Scheme); err != nil {
		return config, err
	}

	return config, nil
}

func (r *OnekindReconciler) desiredDeployment(kind msaberdevv1.Onekind, cm corev1.ConfigMap) (*apps.Deployment, error) {
	depl := &apps.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: apps.SchemeGroupVersion.String(), Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kind.Name,
			Namespace: kind.Namespace,
		},
		Spec: apps.DeploymentSpec{
			Replicas: &kind.Spec.Replicas, // won't be nil because defaulting
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"onekind": kind.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"onekind": kind.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
							Ports: []corev1.ContainerPort{
								{ContainerPort: 80, Name: "http", Protocol: "TCP"},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: cm.Name, ReadOnly: true, MountPath: "/etc/nginx"},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: cm.Name,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: cm.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(&kind, depl, r.Scheme); err != nil {
		return depl, err
	}

	return depl, nil
}

func (r *OnekindReconciler) desiredService(kind msaberdevv1.Onekind) (*corev1.Service, error) {
	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{APIVersion: corev1.SchemeGroupVersion.String(), Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kind.Name,
			Namespace: kind.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "http", Port: 80, Protocol: "TCP", TargetPort: intstr.FromString("http")},
			},
			Selector: map[string]string{"onekind": kind.Name},
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	// always set the controller reference so that we know which object owns this.
	if err := ctrl.SetControllerReference(&kind, svc, r.Scheme); err != nil {
		return svc, err
	}

	return svc, nil
}
