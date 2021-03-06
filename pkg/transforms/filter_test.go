//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package transforms

import (
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/dtos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	deviceName1 = "device1"
	deviceName2 = "device2"

	profileName1 = "profile1"
	profileName2 = "profile2"

	resource1 = "resource1"
	resource2 = "resource2"
	resource3 = "resource3"
)

func TestFilter_FilterByProfileName(t *testing.T) {
	profile1Event := dtos.NewEvent(profileName1, deviceName1)

	tests := []struct {
		Name              string
		Filters           []string
		FilterOut         bool
		EventIn           *dtos.Event
		ExpectedNilResult bool
		ExtraParam        bool
	}{
		{"filter for - no event", []string{profileName1}, true, nil, true, false},
		{"filter for - no filter values", []string{}, false, &profile1Event, false, false},
		{"filter for with extra params - found", []string{profileName1}, false, &profile1Event, false, true},
		{"filter for - found", []string{profileName1}, false, &profile1Event, false, false},
		{"filter for - not found", []string{profileName2}, false, &profile1Event, true, false},

		{"filter out - no event", []string{profileName1}, true, nil, true, false},
		{"filter out - no filter values", []string{}, true, &profile1Event, false, false},
		{"filter out extra param - found", []string{profileName1}, true, &profile1Event, true, true},
		{"filter out - found", []string{profileName1}, true, &profile1Event, true, false},
		{"filter out - not found", []string{profileName2}, true, &profile1Event, false, false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var filter Filter
			if test.FilterOut {
				filter = NewFilterOut(test.Filters)
			} else {
				filter = NewFilterFor(test.Filters)
			}

			expectedContinue := !test.ExpectedNilResult

			if test.EventIn == nil {
				continuePipeline, result := filter.FilterByProfileName(context)
				assert.EqualError(t, result.(error), "no Event Received")
				assert.False(t, continuePipeline)
			} else {
				var continuePipeline bool
				var result interface{}
				if test.ExtraParam {
					continuePipeline, result = filter.FilterByProfileName(context, *test.EventIn, "application/event")
				} else {
					continuePipeline, result = filter.FilterByProfileName(context, *test.EventIn)
				}
				assert.Equal(t, expectedContinue, continuePipeline)
				assert.Equal(t, test.ExpectedNilResult, result == nil)
				if result != nil && test.EventIn != nil {
					assert.Equal(t, *test.EventIn, result)
				}
			}
		})
	}
}

func TestFilter_FilterByDeviceName(t *testing.T) {
	device1Event := dtos.NewEvent(profileName1, deviceName1)

	tests := []struct {
		Name              string
		Filters           []string
		FilterOut         bool
		EventIn           *dtos.Event
		ExpectedNilResult bool
		ExtraParams       bool
	}{
		{"filter for - no event", []string{deviceName1}, false, nil, true, false},
		{"filter for - no filter values", []string{}, false, &device1Event, false, false},
		{"filter for with extra params - found", []string{deviceName1}, false, &device1Event, false, true},
		{"filter for - found", []string{deviceName1}, false, &device1Event, false, false},
		{"filter for - not found", []string{deviceName2}, false, &device1Event, true, false},

		{"filter out - no event", []string{deviceName1}, true, nil, true, false},
		{"filter out - no filter values", []string{}, true, &device1Event, false, false},
		{"filter out extra param - found", []string{deviceName1}, true, &device1Event, true, true},
		{"filter out - found", []string{deviceName1}, true, &device1Event, true, false},
		{"filter out - not found", []string{deviceName2}, true, &device1Event, false, false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var filter Filter
			if test.FilterOut {
				filter = NewFilterOut(test.Filters)
			} else {
				filter = NewFilterFor(test.Filters)
			}

			expectedContinue := !test.ExpectedNilResult

			if test.EventIn == nil {
				continuePipeline, result := filter.FilterByDeviceName(context)
				assert.EqualError(t, result.(error), "no Event Received")
				assert.False(t, continuePipeline)
			} else {
				var continuePipeline bool
				var result interface{}
				if test.ExtraParams {
					continuePipeline, result = filter.FilterByDeviceName(context, *test.EventIn, "application/event")
				} else {
					continuePipeline, result = filter.FilterByDeviceName(context, *test.EventIn)
				}
				assert.Equal(t, expectedContinue, continuePipeline)
				assert.Equal(t, test.ExpectedNilResult, result == nil)
				if result != nil && test.EventIn != nil {
					assert.Equal(t, *test.EventIn, result)
				}
			}
		})
	}
}

func TestFilter_FilterByResourceName(t *testing.T) {
	// event with a reading for resource 1
	resource1Event := dtos.NewEvent(profileName1, deviceName1)
	err := resource1Event.AddSimpleReading(resource1, v2.ValueTypeInt32, int32(123))
	require.NoError(t, err)

	// event with a reading for resource 2
	resource2Event := dtos.NewEvent(profileName1, deviceName1)
	err = resource2Event.AddSimpleReading(resource2, v2.ValueTypeInt32, int32(123))
	require.NoError(t, err)

	// event with a reading for resource 3
	resource3Event := dtos.NewEvent(profileName1, deviceName1)
	err = resource3Event.AddSimpleReading(resource3, v2.ValueTypeInt32, int32(123))
	require.NoError(t, err)

	// event with readings for resource 1 & 2
	twoResourceEvent := dtos.NewEvent(profileName1, deviceName1)
	err = twoResourceEvent.AddSimpleReading(resource1, v2.ValueTypeInt32, int32(123))
	require.NoError(t, err)
	err = twoResourceEvent.AddSimpleReading(resource2, v2.ValueTypeInt32, int32(123))
	require.NoError(t, err)

	tests := []struct {
		Name                 string
		Filters              []string
		FilterOut            bool
		EventIn              *dtos.Event
		ExpectedNilResult    bool
		ExtraParams          bool
		ExpectedReadingCount int
	}{
		{"filter for - no event", []string{resource1}, false, nil, true, false, 0},
		{"filter for extra param - found", []string{resource1}, false, &resource1Event, false, true, 1},
		{"filter for 0 in R1 - no change", []string{}, false, &resource1Event, false, false, 1},
		{"filter for 1 in R1 - 1 of 1 found", []string{resource1}, false, &resource1Event, false, false, 1},
		{"filter for 1 in 2R - 1 of 2 found", []string{resource1}, false, &twoResourceEvent, false, false, 1},
		{"filter for 2 in R1 - 1 of 1 found", []string{resource1, resource2}, false, &resource1Event, false, false, 1},
		{"filter for 2 in 2R - 2 of 2 found", []string{resource1, resource2}, false, &twoResourceEvent, false, false, 2},
		{"filter for 2 in R2 - 1 of 2 found", []string{resource1, resource2}, false, &resource2Event, false, false, 1},
		{"filter for 1 in R2 - not found", []string{resource1}, false, &resource2Event, true, false, 0},

		{"filter out - no event", []string{resource1}, true, nil, true, false, 0},
		{"filter out extra param - found", []string{resource1}, true, &resource1Event, true, true, 0},
		{"filter out 0 in R1 - no change", []string{}, true, &resource1Event, false, false, 1},
		{"filter out 1 in R1 - 1 of 1 found", []string{resource1}, true, &resource1Event, true, false, 0},
		{"filter out 1 in R2 - not found", []string{resource1}, true, &resource2Event, false, false, 1},
		{"filter out 1 in 2R - 1 of 2 found", []string{resource1}, true, &twoResourceEvent, false, false, 1},
		{"filter out 2 in R1 - 1 of 1 found", []string{resource1, resource2}, true, &resource1Event, true, false, 0},
		{"filter out 2 in R2 - 1 of 1 found", []string{resource1, resource2}, true, &resource2Event, true, false, 0},
		{"filter out 2 in 2R - 2 of 2 found", []string{resource1, resource2}, true, &twoResourceEvent, true, false, 0},
		{"filter out 2 in R3 - not found", []string{resource1, resource2}, true, &resource3Event, false, false, 1},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var filter Filter
			if test.FilterOut {
				filter = NewFilterOut(test.Filters)
			} else {
				filter = NewFilterFor(test.Filters)
			}

			expectedContinue := !test.ExpectedNilResult

			if test.EventIn == nil {
				continuePipeline, result := filter.FilterByResourceName(context)
				assert.EqualError(t, result.(error), "no Event Received")
				assert.False(t, continuePipeline)
			} else {
				var continuePipeline bool
				var result interface{}
				if test.ExtraParams {
					continuePipeline, result = filter.FilterByResourceName(context, *test.EventIn, "application/event")
				} else {
					continuePipeline, result = filter.FilterByResourceName(context, *test.EventIn)
				}
				assert.Equal(t, expectedContinue, continuePipeline)
				assert.Equal(t, test.ExpectedNilResult, result == nil)
				if result != nil {
					actualEvent, ok := result.(dtos.Event)
					require.True(t, ok)
					assert.Equal(t, test.ExpectedReadingCount, len(actualEvent.Readings))
				}
			}
		})
	}
}
