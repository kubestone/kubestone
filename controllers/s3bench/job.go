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

package s3bench

import (
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

// NewJob creates a s3bench benchmark job
func NewJob(cr *perfv1alpha1.S3Bench) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	var s3benchCmdLineArgs []string
	s3benchCmdLineArgs = append(s3benchCmdLineArgs, cr.Spec.Mode)
	s3benchCmdLineArgs = append(s3benchCmdLineArgs, ProcessS3BenchArgs(&cr.Spec)...)

	image := cr.Spec.Image
	if (image == perfv1alpha1.ImageSpec{}) {
		image = perfv1alpha1.ImageSpec{
			Name:       "minio/warp:v0.3.5",
			PullPolicy: "IfNotPresent",
		}
	}

	job := k8s.NewPerfJob(objectMeta, "s3bench", image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.Containers[0].Args = s3benchCmdLineArgs
	return job
}

func ProcessS3BenchArgs(spec *perfv1alpha1.S3BenchSpec) []string {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "--host", spec.Host)

	if (spec.S3BenchOptions != perfv1alpha1.S3BenchOptions{}) {
		opts := spec.S3BenchOptions

		checkAndAppendBool(&cmdArgs, opts.NoColor, "--no-color")
		checkAndAppendBool(&cmdArgs, opts.Debug, "--debug")
		checkAndAppendBool(&cmdArgs, opts.Debug, "--insecure")
		checkAndAppendString(&cmdArgs, opts.AccessKey, "--access-key")
		checkAndAppendString(&cmdArgs, opts.SecretKey, "--secret-key")
		checkAndAppendBool(&cmdArgs, opts.Tls, "--tls")
		checkAndAppendString(&cmdArgs, opts.Region, "--region")
		checkAndAppendBool(&cmdArgs, opts.Encrypt, "--encrypt")
		checkAndAppendString(&cmdArgs, opts.Bucket, "--bucket")
		checkAndAppendString(&cmdArgs, opts.HostSelect, "--host-select")
		checkAndAppendInt(&cmdArgs, opts.Concurrent, "--concurrent")
		checkAndAppendBool(&cmdArgs, opts.NoPrefix, "--noprefix")
		checkAndAppendString(&cmdArgs, opts.BenchOutput, "--benchdata")
		checkAndAppendString(&cmdArgs, opts.Duration, "--duration")
		checkAndAppendBool(&cmdArgs, opts.NoClear, "--noclear")

		// TODO: Possibly set this to +1 min from create time or something. May be needed if use clients
		checkAndAppendString(&cmdArgs, opts.SyncStart, "--syncstart")
		checkAndAppendBool(&cmdArgs, opts.Requests, "--requests")
	}

	if (spec.S3ObjectOptions != perfv1alpha1.S3ObjectOptions{}) {
		opts := spec.S3ObjectOptions
		checkAndAppendInt(&cmdArgs, opts.Count, "--objects")
		checkAndAppendString(&cmdArgs, opts.Size, "--obj.size")
		checkAndAppendString(&cmdArgs, opts.Generator, "--obj.generator")
		checkAndAppendBool(&cmdArgs, opts.RandomSize, "--obj.randsize")

	}

	if (spec.S3AutoTermOptions != perfv1alpha1.S3AutoTermOptions{}) {
		opts := spec.S3AutoTermOptions

		if opts.Enabled {
			checkAndAppendBool(&cmdArgs, opts.Enabled, "--autoterm")
			checkAndAppendString(&cmdArgs, opts.Duration, "--autoterm.dur")
			checkAndAppendString(&cmdArgs, opts.Percent, "--autoterm.pct")

		}
	}

	if (spec.S3AnalysisOptions != perfv1alpha1.S3AnalysisOptions{}) {
		opts := spec.S3AnalysisOptions

		checkAndAppendString(&cmdArgs, opts.Duration, "--analyze.dur")
		checkAndAppendString(&cmdArgs, opts.Output, "--analyze.out")
		checkAndAppendString(&cmdArgs, opts.OperationFilter, "--analyze.op")
		checkAndAppendBool(&cmdArgs, opts.PrintErrors, "--analyze.errors")
		checkAndAppendString(&cmdArgs, opts.HostFilter, "--analyze.host")
		checkAndAppendString(&cmdArgs, opts.Skip, "--analyze.skip")
		checkAndAppendBool(&cmdArgs, opts.HostDetails, "--analyze.hostdetails")
	}

	if spec.Mode == "mixed" {
		if (spec.MixedDistributionOptions != perfv1alpha1.MixedDistributionOptions{}) {
			checkAndAppendInt(&cmdArgs, spec.MixedDistributionOptions.GetDist, "--get-distrib")
			checkAndAppendInt(&cmdArgs, spec.MixedDistributionOptions.StatDist, "--stat-distrib")
			checkAndAppendInt(&cmdArgs, spec.MixedDistributionOptions.PutDist, "--put-distrib")
			checkAndAppendInt(&cmdArgs, spec.MixedDistributionOptions.DeleteDist, "--delete-distrib")
		}
	}

	return cmdArgs
}

func checkAndAppendInt(cmdArgs *[]string, dist int32, arg string) {
	if dist > 0 {
		*cmdArgs = append(*cmdArgs, arg, strconv.Itoa(int(dist)))
	}
}

func checkAndAppendBool(cmdArgs *[]string, dist bool, arg string) {
	if dist {
		*cmdArgs = append(*cmdArgs, arg)
	}
}

func checkAndAppendString(cmdArgs *[]string, dist string, arg string) {
	if dist != "" {
		*cmdArgs = append(*cmdArgs, arg, dist)
	}
}
