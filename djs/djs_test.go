// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package djs

import "testing"

func TestBasic(t *testing.T) {
	ds := New()
	CheckedFind := func(i, j int, expected bool, format string) {
		if (ds.Find(i) == ds.Find(j)) != expected {
			t.Errorf(format, i, j)
		}
	}

	ds.MakeSet(1)
	ds.MakeSet(2)
	ds.MakeSet(3)

	CheckedFind(1, 2, false, "%v and %v in same set")
	CheckedFind(1, 3, false, "%v and %v in same set")
	CheckedFind(2, 3, false, "%v and %v in same set")

	ds.Union(1, 2)
	CheckedFind(1, 2, true, "%v and %v not in same set")
	CheckedFind(1, 3, false, "%v and %v in same set")
	CheckedFind(2, 3, false, "%v and %v in same set")

	ds.Union(2, 3)
	CheckedFind(1, 2, true, "%v and %v not in same set")
	CheckedFind(1, 3, true, "%v and %v not in same set")
	CheckedFind(2, 3, true, "%v and %v not in same set")
}

func TestDisjointSet(t *testing.T) {
	vec := []int{
		1, 2, 4, 8, 16, 32, 64, 128, 256, 512,
		1024, 2048, 4096, 8192, 16384, 32768,
	}
	ds := New()
	for _, v := range vec {
		ds.MakeSet(v)
	}

	for i := range vec {
		for j := range vec {
			if i != j && ds.Find(vec[i]) == ds.Find(vec[j]) {
				t.Errorf("%d and %d are in same set", vec[i], vec[j])
			}
		}
	}

	for i := range vec {
		if i%2 == 0 {
			ds.Union(vec[i], vec[i+1])
		}
	}

	for i := range vec {
		if i%2 == 0 {
			if ds.Find(vec[i]) != ds.Find(vec[i+1]) {
				t.Errorf("%d and %d should be in same set", vec[i], vec[i+1])
			}
		} else {
			if i+1 < len(vec) && ds.Find(vec[i]) == ds.Find(vec[i+1]) {
				t.Errorf("%d and %d should be in different sets", vec[i], vec[i+1])
			}
		}
	}

	// Unify the upper and lower part of the array
	for i := 1; i < len(vec)/2; i++ {
		ds.Union(vec[i-1], vec[i])
	}
	for i := len(vec)/2 + 2; i < len(vec); i++ {
		ds.Union(vec[i-1], vec[i])
	}

	for i := range vec {
		m := vec[i]
		n := vec[len(vec)-i-1]
		if ds.Find(m) == ds.Find(n) {
			t.Errorf("%d and %d should be in different sets", m, n)
		}
	}
}
