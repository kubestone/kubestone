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

import "k8s.io/apimachinery/pkg/api/errors"

func ignoreByErrorFn(err error, errorFn func(error) bool) error {
	if errorFn(err) {
		return nil
	}
	return err
}

// IgnoreNotFound returns nil on k8s Not Found type of errors,
// but returns the error as-is otherwise
func IgnoreNotFound(err error) error {
	return ignoreByErrorFn(err, errors.IsNotFound)
}

// IgnoreAlreadyExists returns nil on k8s Already Exists type of errors,
// but returns the error as-is otherwise
func IgnoreAlreadyExists(err error) error {
	return ignoreByErrorFn(err, errors.IsAlreadyExists)
}
