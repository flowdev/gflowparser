package data2svg

import (
	"reflect"
	"testing"
)

func TestAddCluster(t *testing.T) {
	specs := []struct {
		name             string
		givenClusters    clusters
		givenMin         int
		givenMax         int
		expectedClusters clusters
	}{
		{
			name:             "nil clusters",
			givenClusters:    clusters(nil),
			givenMin:         2,
			givenMax:         3,
			expectedClusters: clusters([]int{2, 3}),
		}, {
			name:             "no clusters",
			givenClusters:    clusters([]int{}),
			givenMin:         2,
			givenMax:         3,
			expectedClusters: clusters([]int{2, 3}),
		}, {
			name:             "double first clusters",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         2,
			givenMax:         3,
			expectedClusters: clusters([]int{2, 3, 7, 9}),
		}, {
			name:             "double second clusters",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         7,
			givenMax:         9,
			expectedClusters: clusters([]int{2, 3, 7, 9}),
		}, {
			name:             "add third cluster",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         10,
			givenMax:         13,
			expectedClusters: clusters([]int{2, 3, 7, 9, 10, 13}),
		}, {
			name:             "add first cluster",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         0,
			givenMax:         1,
			expectedClusters: clusters([]int{0, 1, 2, 3, 7, 9}),
		}, {
			name:             "add middle cluster",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         4,
			givenMax:         5,
			expectedClusters: clusters([]int{2, 3, 4, 5, 7, 9}),
		}, {
			name:             "extend first cluster at start",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         1,
			givenMax:         2,
			expectedClusters: clusters([]int{1, 3, 7, 9}),
		}, {
			name:             "span over all clusters",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         1,
			givenMax:         11,
			expectedClusters: clusters([]int{1, 11}),
		}, {
			name:             "extend first cluster at end",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         3,
			givenMax:         5,
			expectedClusters: clusters([]int{2, 5, 7, 9}),
		}, {
			name:             "merge both clusters",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         3,
			givenMax:         7,
			expectedClusters: clusters([]int{2, 9}),
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got := spec.givenClusters.addCluster(spec.givenMin, spec.givenMax)
			checkClusters(t, spec.expectedClusters, got)
		})
	}
}

func TestGetCluster(t *testing.T) {
	specs := []struct {
		name          string
		givenClusters clusters
		givenIdx      int
		expectedMin   int
		expectedMax   int
	}{
		{
			name:          "nil clusters",
			givenClusters: clusters(nil),
			givenIdx:      2,
			expectedMin:   2,
			expectedMax:   2,
		}, {
			name:          "before first cluster",
			givenClusters: clusters([]int{2, 3, 7, 9}),
			givenIdx:      1,
			expectedMin:   1,
			expectedMax:   1,
		}, {
			name:          "start of first cluster",
			givenClusters: clusters([]int{2, 3, 7, 9}),
			givenIdx:      2,
			expectedMin:   2,
			expectedMax:   3,
		}, {
			name:          "end of first cluster",
			givenClusters: clusters([]int{2, 3, 7, 9}),
			givenIdx:      3,
			expectedMin:   2,
			expectedMax:   3,
		}, {
			name:          "middle of second cluster",
			givenClusters: clusters([]int{2, 3, 7, 9}),
			givenIdx:      8,
			expectedMin:   7,
			expectedMax:   9,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			gotMin, gotMax := spec.givenClusters.getCluster(spec.givenIdx)
			if gotMin != spec.expectedMin {
				t.Errorf("Expected min %d, got %d.", spec.expectedMin, gotMin)
			}
			if gotMax != spec.expectedMax {
				t.Errorf("Expected max %d, got %d.", spec.expectedMax, gotMax)
			}
		})
	}
}

func TestDeleteLine(t *testing.T) {
	specs := []struct {
		name             string
		givenClusters    clusters
		givenIdx         int
		expectedClusters clusters
	}{
		{
			name:             "no clusters",
			givenClusters:    clusters([]int{}),
			givenIdx:         2,
			expectedClusters: clusters([]int{}),
		}, {
			name:             "delete after last cluster",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenIdx:         10,
			expectedClusters: clusters([]int{2, 3, 7, 9}),
		}, {
			name:             "delete last line of last cluster",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenIdx:         9,
			expectedClusters: clusters([]int{2, 3, 7, 8}),
		}, {
			name:             "delete middle line of last cluster",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenIdx:         8,
			expectedClusters: clusters([]int{2, 3, 7, 8}),
		}, {
			name:             "delete first line of last cluster",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenIdx:         7,
			expectedClusters: clusters([]int{2, 3, 7, 8}),
		}, {
			name:             "delete line between clusters",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenIdx:         5,
			expectedClusters: clusters([]int{2, 3, 6, 8}),
		}, {
			name:             "delete first line at all",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenIdx:         0,
			expectedClusters: clusters([]int{1, 2, 6, 8}),
		}, {
			name:             "delete line that removes first cluster",
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenIdx:         3,
			expectedClusters: clusters([]int{6, 8}),
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			got := spec.givenClusters.deleteLine(spec.givenIdx)
			checkClusters(t, spec.expectedClusters, got)
		})
	}
}

func checkClusters(t *testing.T, exp, got clusters) {
	if !reflect.DeepEqual(got, exp) {
		t.Error("Expected and actual clusters differ:")
		t.Logf("Expected clusters: %v", exp)
		t.Logf("Actual   clusters: %v", got)
	}
}
