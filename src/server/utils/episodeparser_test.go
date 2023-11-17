package utils

import (
	"testing"
)

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

func TestGetTokens(t *testing.T) {
	s := "[爱恋字幕社][10月新番][星灵感应][Hoshikuzu Telepath][06][720P][MP4][GB][简中]"
	tokens, err := getTokens(s)
	if err != nil {
		t.Error(err)
	}
	if len(tokens) != 9 {
		t.Error("tokens size error")
	}

	s = "[爱恋字幕社][10月新番][星灵感应]Hoshikuzu Telepath[06][720P][MP4][GB][简中]"
	tokens, err = getTokens(s)
	if err != nil {
		t.Error(err)
	}
	if len(tokens) != 10 {
		t.Error("tokens size error")
	}
}

func TestChinese2Num(t *testing.T) {
	s := "零"
	num, _ := chinese2Num(s)
	if num != 0 {
		t.Error("num != 18")
	}
	s = "一"
	num, _ = chinese2Num(s)
	if num != 1 {
		t.Error("num != 18")
	}
	s = "十"
	num, _ = chinese2Num(s)
	if num != 10 {
		t.Error("num != 10")
	}
	s = "十八"
	num, _ = chinese2Num(s)
	if num != 18 {
		t.Error("num != 10")
	}
	s = "十十八"
	num, _ = chinese2Num(s)
	if num != 108 {
		t.Error("num != 18")
	}
}
