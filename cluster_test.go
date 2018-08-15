package gflowparser

import "testing"

func TestGrowClusters(t *testing.T) {
	specs := []struct {
		givenClusters    clusters
		givenMin         int
		givenMax         int
		expectedClusters clusters
	}{
		{
			givenClusters:    clusters([]int{}),
			givenMin:         2,
			givenMax:         3,
			expectedClusters: clusters([]int{2, 3}),
		}, {
			givenClusters:    clusters(nil),
			givenMin:         2,
			givenMax:         3,
			expectedClusters: clusters([]int{2, 3}),
		}, {
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         2,
			givenMax:         3,
			expectedClusters: clusters([]int{2, 3, 7, 9}),
		}, {
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         7,
			givenMax:         9,
			expectedClusters: clusters([]int{2, 3, 7, 9}),
		}, {
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         10,
			givenMax:         13,
			expectedClusters: clusters([]int{2, 3, 7, 9, 10, 13}),
		}, {
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         0,
			givenMax:         1,
			expectedClusters: clusters([]int{0, 1, 2, 3, 7, 9}),
		}, {
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         4,
			givenMax:         5,
			expectedClusters: clusters([]int{2, 3, 4, 5, 7, 9}),
		}, {
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         1,
			givenMax:         2,
			expectedClusters: clusters([]int{1, 3, 7, 9}),
		}, {
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         1,
			givenMax:         11,
			expectedClusters: clusters([]int{1, 11}),
		}, {
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         3,
			givenMax:         5,
			expectedClusters: clusters([]int{2, 5, 7, 9}),
		}, {
			givenClusters:    clusters([]int{2, 3, 7, 9}),
			givenMin:         3,
			givenMax:         7,
			expectedClusters: clusters([]int{2, 9}),
		},
	}

	for i, spec := range specs {
		t.Logf("TestGrowClusters[%d]:", i)
		act := spec.givenClusters.addCluster(spec.givenMin, spec.givenMax)
		checkClusters(t, spec.expectedClusters, act)
	}
}

func checkClusters(t *testing.T, exp, act clusters) {
	if len(exp) != len(act) {
		t.Errorf("Expected clusters length of %d, got %d.", len(exp), len(act))
		return
	}
	for i, e := range exp {
		if e != act[i] {
			t.Errorf("Expected value %d at index %d, got %d.", 3, i, act[i])
			t.Logf("Expected clusters: %v", exp)
			t.Logf("Actual   clusters: %v", act)
			return
		}
	}
}
