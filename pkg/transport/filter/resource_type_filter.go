// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0
package filter

import (
	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

type resourceTypeFilter struct {
	includeResourceTypes []string
}

func (f resourceTypeFilter) Matches(cd *cdv2.ComponentDescriptor, r cdv2.Resource) bool {
	for _, resourceType := range f.includeResourceTypes {
		if r.Type == resourceType {
			return true
		}
	}
	return false
}

func NewResourceTypeFilter(includeResourceTypes ...string) (Filter, error) {
	filter := resourceTypeFilter{
		includeResourceTypes: includeResourceTypes,
	}

	return &filter, nil
}
