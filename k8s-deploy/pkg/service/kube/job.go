/*
 * Copyright 2019-2022 VMware, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package kube

import (
	"context"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Job interface {
	GetJobByName(namespace, jobName string) (*v1.Job, error)
}

func (e *Kube) GetJobByName(namespace, jobName string) (*v1.Job, error) {
	job, err := e.client.BatchV1().Jobs(namespace).Get(context.Background(), jobName, metav1.GetOptions{})
	return job, err
}
