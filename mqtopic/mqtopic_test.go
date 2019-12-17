package mqtopic

import (
	"reflect"
	"testing"
)

func err2str(err error) string {
	if err == nil {
		return "(nil error)"
	}
	return err.Error()
}

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		in  string
		exp Topic
		err error
	}{
		{"", nil, ErrAtLeastOneCharacter},
		{"/", Topic{"", ""}, nil},
		{"aaa/bbb/ccc/ddd", Topic{"aaa", "bbb", "ccc", "ddd"}, nil},
		{"aaa/bbb/ccc/", Topic{"aaa", "bbb", "ccc", ""}, nil},
		{"aaa/bbb//ddd", Topic{"aaa", "bbb", "", "ddd"}, nil},
		{"aaa//ccc/ddd", Topic{"aaa", "", "ccc", "ddd"}, nil},
		{"/bbb/ccc/ddd", Topic{"", "bbb", "ccc", "ddd"}, nil},
		{"#", nil, ErrWildcardsInTopicName},
		{"+", nil, ErrWildcardsInTopicName},
		{"aaa/bbb/#", nil, ErrWildcardsInTopicName},
		{"aaa/#/ccc", nil, ErrWildcardsInTopicName},
		{"#/bbb/ccc", nil, ErrWildcardsInTopicName},
		{"aaa/bbb/+", nil, ErrWildcardsInTopicName},
		{"aaa/+/ccc", nil, ErrWildcardsInTopicName},
		{"+/bbb/ccc", nil, ErrWildcardsInTopicName},
		{"#+", nil, ErrWildcardsInTopicName},
		{"foo#", nil, ErrWildcardsInTopicName},
		{"foo+", nil, ErrWildcardsInTopicName},
		{"aaa/bbb/cc#", nil, ErrWildcardsInTopicName},
		{"aaa/#bb/ccc", nil, ErrWildcardsInTopicName},
		{"a#a/bbb/ccc", nil, ErrWildcardsInTopicName},
		{"aaa/bbb/cc+", nil, ErrWildcardsInTopicName},
		{"aaa/+bb/ccc", nil, ErrWildcardsInTopicName},
		{"a+a/bbb/ccc", nil, ErrWildcardsInTopicName},
	} {
		act, err := Parse(tc.in)
		if act := err2str(err); act != err2str(tc.err) {
			t.Fatalf("unexpecetd error: in=%q\nexpect=%s\nactual=%s", tc.in, tc.err, act)
			continue
		}
		if !reflect.DeepEqual(tc.exp, act) {
			t.Fatalf("unexpected result: in=%q\nexpect=%+v\nactual=%+v", tc.in, tc.exp, act)
		}
	}
}

func TestParseFilter(t *testing.T) {
	for _, tc := range []struct {
		in  string
		exp Filter
		err error
	}{
		{"", nil, ErrAtLeastOneCharacter},
		{"/", Filter{"", ""}, nil},
		{"aaa/bbb/ccc/ddd", Filter{"aaa", "bbb", "ccc", "ddd"}, nil},
		{"aaa/bbb/ccc/", Filter{"aaa", "bbb", "ccc", ""}, nil},
		{"aaa/bbb//ddd", Filter{"aaa", "bbb", "", "ddd"}, nil},
		{"aaa//ccc/ddd", Filter{"aaa", "", "ccc", "ddd"}, nil},
		{"/bbb/ccc/ddd", Filter{"", "bbb", "ccc", "ddd"}, nil},

		{"#", Filter{"#"}, nil},
		{"aaa/bbb/#", Filter{"aaa", "bbb", "#"}, nil},
		{"aaa/#/ccc", nil, ErrMultiLevelWildcardNotLast},
		{"#/bbb/ccc", nil, ErrMultiLevelWildcardNotLast},
		{"foo#", nil, ErrWildcardsCombinedInLevel},
		{"aaa/bbb/cc#", nil, ErrWildcardsCombinedInLevel},
		{"aaa/#bb/ccc", nil, ErrWildcardsCombinedInLevel},
		{"a#a/bbb/ccc", nil, ErrWildcardsCombinedInLevel},

		{"+", Filter{"+"}, nil},
		{"aaa/bbb/+", Filter{"aaa", "bbb", "+"}, nil},
		{"aaa/+/ccc", Filter{"aaa", "+", "ccc"}, nil},
		{"+/bbb/ccc", Filter{"+", "bbb", "ccc"}, nil},
		{"foo+", nil, ErrWildcardsCombinedInLevel},
		{"aaa/bbb/cc+", nil, ErrWildcardsCombinedInLevel},
		{"aaa/+bb/ccc", nil, ErrWildcardsCombinedInLevel},
		{"a+a/bbb/ccc", nil, ErrWildcardsCombinedInLevel},

		{"#+", nil, ErrWildcardsCombinedInLevel},
	} {
		act, err := ParseFilter(tc.in)
		if act := err2str(err); act != err2str(tc.err) {
			t.Fatalf("unexpecetd error: in=%q\nexpect=%s\nactual=%s", tc.in, tc.err, act)
			continue
		}
		if !reflect.DeepEqual(tc.exp, act) {
			t.Fatalf("unexpected result: in=%q\nexpect=%+v\nactual=%+v", tc.in, tc.exp, act)
		}
	}
}

func TestFilter_Match(t *testing.T) {
	for _, tc := range []struct {
		filter    string
		matches   []string
		unmatches []string
	}{
		{"#",
			[]string{"/", "abc", "foo/bar/baz"},
			[]string{"$SYS", "$SYS/control"}},
		{"sport/#",
			[]string{"sport", "sport/", "sport/tennis", "sport/tennis/player1"},
			[]string{"sports", "/", "abc", "foo/bar/baz"}},
		{"sport/tennis/player1/#",
			[]string{
				"sport/tennis/player1",
				"sport/tennis/player1/ranking",
				"sport/tennis/player1/ranking/wimbledon",
			}, nil},
		{"sport/tennis/+",
			[]string{
				"sport/tennis/player1",
				"sport/tennis/player2",
			}, []string{
				"sport/tennis/player1/ranking",
				"sport/tennis",
				"sport/",
				"sport",
				"/",
			}},
		{"sport/+", []string{"sport/"}, []string{"sport"}},
		{"sport/+/player1",
			[]string{
				"sport//player1",
				"sport/tennis/player1",
				"sport/baseball/player1",
				"sport/football/player1",
			}, []string{
				"sport//player2",
				"sport/tennis/player2",
				"sport/baseball/player2",
				"sport/football/player2",
			}},
		{"+/+", []string{"/finance"}, nil},
		{"/+", []string{"/finance"}, nil},
		{"+", nil, []string{"/finance"}},
	} {
		f, err := ParseFilter(tc.filter)
		if err != nil {
			t.Fatalf("parse filter failed %q: %s", tc.filter, err)
		}
		for _, s := range tc.matches {
			topic, err := Parse(s)
			if err != nil {
				t.Fatalf("parse topic failed %q: %s", s, err)
			}
			if !f.Match(topic) {
				t.Fatalf("unexpected unmatch: filter=%s topic=%s", tc.filter, s)
			}
		}
		for _, s := range tc.unmatches {
			topic, err := Parse(s)
			if err != nil {
				t.Fatalf("parse topic failed %q: %s", s, err)
			}
			if f.Match(topic) {
				t.Fatalf("unexpected match: filter=%s topic=%s", tc.filter, s)
			}
		}
	}
}
