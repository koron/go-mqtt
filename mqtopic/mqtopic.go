// Package mqtopic provides MQTT Topic Name and Filter capability.
package mqtopic

import (
	"errors"
	"strings"
)

var (
	// ErrWildcardsInTopicName shows wildcard characters in Topic Name are not
	// allowed.
	ErrWildcardsInTopicName = errors.New("wildcard characters MUST NOT be used within a Topic Name [MQTT-4.7.1-1]")

	// ErrAtLeastOneCharacter shows empty Topic Name or Filter are not allowed.
	ErrAtLeastOneCharacter = errors.New("at least one character long [MQTT-4.7.3-1]")

	// ErrMultiLevelWildcardNotLast shows `#` multi-level wildcard is not
	// allowed at top or middle of Topic Filter. It must be placed at last
	// part.
	ErrMultiLevelWildcardNotLast = errors.New("multi-level wildcard MUST be the last of Topic Filter [MQTT-4.7.1-2]")

	// ErrWildcardsCombinedInLevel shows wildcard characters are combined with
	// other characters in a level of Topic Filter.
	ErrWildcardsCombinedInLevel = errors.New("combined wildcards in a level is not allowed [MQTT-4.7.1-2 MQTT-4.7.1-3]")
)

// Topic is topic name.
type Topic []string

// Parse parses a string as topic name.
func Parse(s string) (Topic, error) {
	if s == "" {
		return nil, ErrAtLeastOneCharacter
	}
	topic := Topic(strings.Split(s, "/"))
	for _, n := range topic {
		m := strings.IndexAny(n, "#+")
		if m >= 0 {
			return nil, ErrWildcardsInTopicName
		}
	}
	return topic, nil
}

// Filter is topic filter.
type Filter []string

// ParseFilter parses a string as topic filter.
func ParseFilter(s string) (Filter, error) {
	if s == "" {
		return nil, ErrAtLeastOneCharacter
	}
	filter := Filter(strings.Split(s, "/"))
	for i, n := range filter {
		if x := strings.IndexRune(n, '#'); x >= 0 {
			if len(n) != 1 {
				return nil, ErrWildcardsCombinedInLevel
			}
			if i != len(filter)-1 {
				return nil, ErrMultiLevelWildcardNotLast
			}
			continue
		}
		if x := strings.IndexRune(n, '+'); x >= 0 {
			if len(n) != 1 {
				return nil, ErrWildcardsCombinedInLevel
			}
			continue
		}
	}
	return filter, nil
}

// Match checks whether a topic name matches filter or not.
func (f Filter) Match(topic Topic) bool {
	if last := len(f) - 1; f[last] == "#" {
		if last == 0 {
			if strings.HasPrefix(topic[0], "$") {
				return false
			}
			return true
		}
		if len(topic) < last {
			return false
		}
		if !f[:last].match(topic[:last]) {
			return false
		}
		return true
	}
	return f.match(topic)
}

func (f Filter) match(topic Topic) bool {
	if len(f) != len(topic) {
		return false
	}
	for i, item := range f {
		n := topic[i]
		if item == "+" {
			if i == 0 && strings.HasPrefix(n, "$") {
				return false
			}
			continue
		}
		if n != item {
			return false
		}
	}
	return true
}
