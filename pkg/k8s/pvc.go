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

package k8s

import (
	"github.com/xridge/kubestone/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewPersistentVolumeClaim creates a PVC based on the provided volumeSpec, name and namespace
func NewPersistentVolumeClaim(volumeSpec *v1alpha1.VolumeSpec, name, namespace string) (*corev1.PersistentVolumeClaim, error) {
	accessModes := make([]corev1.PersistentVolumeAccessMode, len(volumeSpec.PersistentVolumeClaim.AccessModes))
	for i, accessMode := range volumeSpec.PersistentVolumeClaim.AccessModes {
		accessModes[i] = corev1.PersistentVolumeAccessMode(accessMode)
	}

	requests := make(corev1.ResourceList, 1)
	quantity, err := resource.ParseQuantity(string(volumeSpec.PersistentVolumeClaim.Size))
	if err != nil {
		return nil, err
	}
	requests[corev1.ResourceStorage] = quantity

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			Selector:         volumeSpec.PersistentVolumeClaim.Selector,
			VolumeName:       volumeSpec.PersistentVolumeClaim.VolumeName,
			StorageClassName: volumeSpec.PersistentVolumeClaim.StorageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: requests,
			},
		},
	}

	if volumeSpec.PersistentVolumeClaim.VolumeMode != nil {
		volumeMode := corev1.PersistentVolumeMode(*volumeSpec.PersistentVolumeClaim.VolumeMode)
		pvc.Spec.VolumeMode = &volumeMode
	}

	return &pvc, nil
}
