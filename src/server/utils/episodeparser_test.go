package utils

import "testing"

func TestParseEpisode(t *testing.T) {
	names := []string{"abc 1 ss", "abc 2 ss", "abc 10 ss"}
	eps := ParseEpisode(names)
	if len(eps) != 3 {
		t.Errorf("len(eps)(%d) != 3", len(eps))
		return
	}
	ep1, ok := eps[1]
	if !ok {
		t.Errorf("1 not int eps")
		return
	}
	ep2, ok := eps[2]
	if !ok {
		t.Errorf("2 not int eps")
		return
	}
	ep10, ok := eps[10]
	if !ok {
		t.Errorf("10 not int eps")
		return
	}
	if ep1 != 0 {
		t.Errorf("ep1(%d) != 0", ep1)
	}
	if ep2 != 1 {
		t.Errorf("ep2(%d) != 0", ep2)
	}
	if ep10 != 2 {
		t.Errorf("ep10(%d) != 0", ep10)
	}
}
