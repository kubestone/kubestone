/*
Copyright 2019 The xridge kubestone contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fio

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// NewPersistentVolumeClaim creates a PVC
func NewPersistentVolumeClaim(cr *perfv1alpha1.Fio) (*corev1.PersistentVolumeClaim, error) {
	if cr.Spec.PersistentVolumeClaim == nil {
		return nil, nil
	}

	accessModes := make([]corev1.PersistentVolumeAccessMode, len(cr.Spec.PersistentVolumeClaim.AccessModes))
	for i, accessMode := range cr.Spec.PersistentVolumeClaim.AccessModes {
		accessModes[i] = corev1.PersistentVolumeAccessMode(accessMode)
	}

	requests := make(corev1.ResourceList, 1)
	quantity, err := resource.ParseQuantity(string(cr.Spec.PersistentVolumeClaim.Size))
	if err != nil {
		return nil, err
	}
	requests[corev1.ResourceStorage] = quantity

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			Selector:         cr.Spec.PersistentVolumeClaim.Selector,
			VolumeName:       cr.Spec.PersistentVolumeClaim.VolumeName,
			StorageClassName: cr.Spec.PersistentVolumeClaim.StorageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: requests,
			},
		},
	}

	if cr.Spec.PersistentVolumeClaim.VolumeMode != nil {
		volumeMode := corev1.PersistentVolumeMode(*cr.Spec.PersistentVolumeClaim.VolumeMode)
		pvc.Spec.VolumeMode = &volumeMode
	}

	return &pvc, nil
}
