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

package ioping

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// NewPersistentVolumeClaim creates a PVC
// TODO: Factor this into a common place
func NewPersistentVolumeClaim(cr *perfv1alpha1.Ioping) (*corev1.PersistentVolumeClaim, error) {
	accessModes := make([]corev1.PersistentVolumeAccessMode, len(cr.Spec.Volume.PersistentVolumeClaim.AccessModes))
	for i, accessMode := range cr.Spec.Volume.PersistentVolumeClaim.AccessModes {
		accessModes[i] = corev1.PersistentVolumeAccessMode(accessMode)
	}

	requests := make(corev1.ResourceList, 1)
	quantity, err := resource.ParseQuantity(string(cr.Spec.Volume.PersistentVolumeClaim.Size))
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
			Selector:         cr.Spec.Volume.PersistentVolumeClaim.Selector,
			VolumeName:       cr.Spec.Volume.PersistentVolumeClaim.VolumeName,
			StorageClassName: cr.Spec.Volume.PersistentVolumeClaim.StorageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: requests,
			},
		},
	}

	if cr.Spec.Volume.PersistentVolumeClaim.VolumeMode != nil {
		volumeMode := corev1.PersistentVolumeMode(*cr.Spec.Volume.PersistentVolumeClaim.VolumeMode)
		pvc.Spec.VolumeMode = &volumeMode
	}

	return &pvc, nil
}
