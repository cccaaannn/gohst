package url

import "strings"

type segmentType string

const (
	static   segmentType = "static"
	param    segmentType = "param"
	wildcard segmentType = "wildcard"
)

type segment struct {
	value       string
	segmentType segmentType
}

type Path struct {
	pattern  string
	segments []segment
}

func parsePathSegments(pattern string) []segment {
	segments := make([]segment, 0)
	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		segmentType := static
		if strings.HasPrefix(part, ":") {
			segmentType = param
			part = part[1:]
		}
		if part == "*" {
			segmentType = wildcard
		}
		segments = append(segments, segment{value: part, segmentType: segmentType})
	}
	return segments
}

func (path Path) hasWildcard() bool {
	hasWildcard := false
	for _, segment := range path.segments {
		if segment.segmentType == wildcard {
			hasWildcard = true
			break
		}
	}
	return hasWildcard
}

func CreatePath(pattern string) Path {
	segments := parsePathSegments(pattern)
	return Path{pattern: pattern, segments: segments}
}

func (path Path) Match(text string) (map[string]string, bool) {
	textSegments := strings.Split(text, "/")
	params := make(map[string]string)
	hasWildcard := path.hasWildcard()

	if !hasWildcard {
		if len(path.segments) != len(textSegments) {
			return nil, false
		}
	} else {
		if len(path.segments) > len(textSegments) {
			return nil, false
		}
	}

	for i, segment := range path.segments {
		textSegment := textSegments[i]
		if segment.segmentType == param {
			params[segment.value] = textSegment
			continue
		}
		if segment.segmentType == static && segment.value != textSegment {
			return nil, false
		}
		if segment.segmentType == wildcard {
			return params, true
		}
	}
	return params, true
}
