// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/GoogleCloudPlatform/prometheus-engine/pkg/operator/generated/clientset/versioned/typed/monitoring/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeMonitoringV1alpha1 struct {
	*testing.Fake
}

func (c *FakeMonitoringV1alpha1) ClusterRules() v1alpha1.ClusterRulesInterface {
	return &FakeClusterRules{c}
}

func (c *FakeMonitoringV1alpha1) OperatorConfigs(namespace string) v1alpha1.OperatorConfigInterface {
	return &FakeOperatorConfigs{c, namespace}
}

func (c *FakeMonitoringV1alpha1) PodMonitorings(namespace string) v1alpha1.PodMonitoringInterface {
	return &FakePodMonitorings{c, namespace}
}

func (c *FakeMonitoringV1alpha1) Rules(namespace string) v1alpha1.RulesInterface {
	return &FakeRules{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeMonitoringV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
