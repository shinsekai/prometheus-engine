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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/GoogleCloudPlatform/prometheus-engine/pkg/operator/apis/monitoring/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ClusterPodMonitoringLister helps list ClusterPodMonitorings.
// All objects returned here must be treated as read-only.
type ClusterPodMonitoringLister interface {
	// List lists all ClusterPodMonitorings in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ClusterPodMonitoring, err error)
	// Get retrieves the ClusterPodMonitoring from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ClusterPodMonitoring, error)
	ClusterPodMonitoringListerExpansion
}

// clusterPodMonitoringLister implements the ClusterPodMonitoringLister interface.
type clusterPodMonitoringLister struct {
	indexer cache.Indexer
}

// NewClusterPodMonitoringLister returns a new ClusterPodMonitoringLister.
func NewClusterPodMonitoringLister(indexer cache.Indexer) ClusterPodMonitoringLister {
	return &clusterPodMonitoringLister{indexer: indexer}
}

// List lists all ClusterPodMonitorings in the indexer.
func (s *clusterPodMonitoringLister) List(selector labels.Selector) (ret []*v1alpha1.ClusterPodMonitoring, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ClusterPodMonitoring))
	})
	return ret, err
}

// Get retrieves the ClusterPodMonitoring from the index for a given name.
func (s *clusterPodMonitoringLister) Get(name string) (*v1alpha1.ClusterPodMonitoring, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("clusterpodmonitoring"), name)
	}
	return obj.(*v1alpha1.ClusterPodMonitoring), nil
}
