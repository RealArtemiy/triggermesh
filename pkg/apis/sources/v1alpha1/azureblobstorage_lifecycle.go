/*
Copyright 2022 TriggerMesh Inc.

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

package v1alpha1

import (
	"sort"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	"github.com/triggermesh/triggermesh/pkg/apis/common/v1alpha1"
)

// GetGroupVersionKind implements kmeta.OwnerRefable.
func (s *AzureBlobStorageSource) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("AzureBlobStorageSource")
}

// GetConditionSet implements duckv1.KRShaped.
func (s *AzureBlobStorageSource) GetConditionSet() apis.ConditionSet {
	return azureBlobStorageSourceConditionSet
}

// GetStatus implements duckv1.KRShaped.
func (s *AzureBlobStorageSource) GetStatus() *duckv1.Status {
	return &s.Status.Status.Status
}

// GetSink implements EventSender.
func (s *AzureBlobStorageSource) GetSink() *duckv1.Destination {
	return &s.Spec.Sink
}

// GetStatusManager implements Reconcilable.
func (s *AzureBlobStorageSource) GetStatusManager() *v1alpha1.StatusManager {
	return &v1alpha1.StatusManager{
		ConditionSet: s.GetConditionSet(),
		Status:       &s.Status.Status,
	}
}

// AsEventSource implements EventSource.
func (s *AzureBlobStorageSource) AsEventSource() string {
	return s.Spec.StorageAccountID.String()
}

// GetAdapterOverrides implements AdapterConfigurable.
func (s *AzureBlobStorageSource) GetAdapterOverrides() *v1alpha1.AdapterOverrides {
	return s.Spec.AdapterOverrides
}

// Default event types.
// This list is non-exhaustive, see AzureBlobStorageSourceSpec.
const (
	AzureBlobStorageBlobCreatedEventType = "Microsoft.Storage.BlobCreated"
	AzureBlobStorageBlobDeletedEventType = "Microsoft.Storage.BlobDeleted"
)

// GetEventTypes returns the event types generated by the source.
func (s *AzureBlobStorageSource) GetEventTypes() []string {
	if s.Spec.EventTypes == nil {
		return defaultAzureStorageEventTypes()
	}

	selectedTypes := make(map[string]struct{})
	for _, t := range s.Spec.EventTypes {
		if _, alreadySet := selectedTypes[t]; !alreadySet {
			selectedTypes[t] = struct{}{}
		}
	}

	eventTypes := make([]string, 0, len(selectedTypes))

	for t := range selectedTypes {
		eventTypes = append(eventTypes, t)
	}

	sort.Strings(eventTypes)

	return eventTypes
}

// defaultAzureStorageEventTypes returns the list of storage event types which
// are enabled by default by Azure when not explicitly defined by the user.
func defaultAzureStorageEventTypes() []string {
	return []string{
		AzureBlobStorageBlobCreatedEventType,
		AzureBlobStorageBlobDeletedEventType,
	}
}

// Status conditions
const (
	// AzureBlobStorageConditionSubscribed has status True when an event subscription exists for the source.
	AzureBlobStorageConditionSubscribed apis.ConditionType = "Subscribed"
)

// azureBlobStorageSourceConditionSet is a set of conditions for
// AzureBlobStorageSource objects.
var azureBlobStorageSourceConditionSet = v1alpha1.NewConditionSet(
	AzureBlobStorageConditionSubscribed,
)

// MarkSubscribed sets the Subscribed condition to True.
func (s *AzureBlobStorageSourceStatus) MarkSubscribed() {
	azureBlobStorageSourceConditionSet.Manage(s).MarkTrue(AzureBlobStorageConditionSubscribed)
}

// MarkNotSubscribed sets the Subscribed condition to False with the given
// reason and message.
func (s *AzureBlobStorageSourceStatus) MarkNotSubscribed(reason, msg string) {
	azureBlobStorageSourceConditionSet.Manage(s).MarkFalse(AzureBlobStorageConditionSubscribed, reason, msg)
}
