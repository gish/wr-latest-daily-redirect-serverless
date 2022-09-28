package main

import "testing"

func TestBuildUserAgentReturnsCorrectAgent(t *testing.T) {
	expected := "platform:app-id:1.0.0 (by /u/kevin)"
	actual := buildUserAgent("platform", "app-id", "1.0.0", "kevin")

	if expected != actual {
		t.Errorf("[expected, actual] [%s, %s]", expected, actual)
	}
}
