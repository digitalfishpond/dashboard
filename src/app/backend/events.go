// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"k8s.io/kubernetes/pkg/api/unversioned"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
)

// Events response structure.
type Events struct {
	// Namespace.
	Namespace string `json:"namespace"`

	// List of events from given namespace.
	Events []Event `json:"events"`
}

// Single event representation.
type Event struct {
	// Event message.
	Message string `json:"message"`

	// Component from which the event is generated.
	SourceComponent string `json:"sourceComponent"`

	// Host name on which the event is generated.
	SourceHost string `json:"sourceHost"`

	// Reference to a piece of an object, which triggered an event. For example
	// "spec.containers{name}" refers to container within pod with given name, if no container
	// name is specified, for example "spec.containers[2]", then it refers to container with
	// index 2 in this pod.
	SubObject string `json:"object"`

	// The number of times this event has occurred.
	Count int `json:"count"`

	// The time at which the event was first recorded.
	FirstSeen unversioned.Time `json:"firstSeen"`

	// The time at which the most recent occurrence of this event was recorded.
	LastSeen unversioned.Time `json:"lastSeen"`

	// Reason why this event was generated.
	Reason string `json:"reason"`
}

// Return events for particular namespace or error if occurred.
func GetEvents(client *client.Client, namespace string) (*Events, error) {
	list, err := client.Events(namespace).List(unversioned.ListOptions{
		LabelSelector: unversioned.LabelSelector{labels.Everything()},
		FieldSelector: unversioned.FieldSelector{fields.Everything()},
	})

	if err != nil {
		return nil, err
	}

	events := &Events{
		Namespace: namespace,
	}

	for _, element := range list.Items {
		events.Events = append(events.Events, Event{
			Message:         element.Message,
			SourceComponent: element.Source.Component,
			SourceHost:      element.Source.Host,
			SubObject:       element.InvolvedObject.FieldPath,
			Count:           element.Count,
			FirstSeen:       element.FirstTimestamp,
			LastSeen:        element.LastTimestamp,
			Reason:          element.Reason,
		})
	}

	return events, nil
}