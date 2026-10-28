package tdp

import "testing"

func TestVersionOrder(t *testing.T) {
	order := []string{
		"v0.9.16", "v0.11.5", "v1.0.2",
		"v1.1.0-alpha", "v1.1.0-alpha.1", "v1.1.0-beta.1", "v1.1.0-beta.2", "v1.1.0-beta.10", "v1.1.0-rc.1",
		"v1.1.0", "v1.1.1",
	}
	for i := range order {
		a, ok := semver(order[i])
		if !ok {
			t.Fatalf("%s is not a version", order[i])
		}
		for j := range order {
			b, _ := semver(order[j])
			if less(a, b) != (i < j) {
				t.Errorf("less(%s, %s) = %v", order[i], order[j], !(i < j))
			}
		}
	}
	for _, tag := range []string{"1.0.2", "v1.1", "v1.1.0-", "v1.1.0-beta..1", "v1.x.0", "main"} {
		if _, ok := semver(tag); ok {
			t.Errorf("%q taken as a version", tag)
		}
	}
}
