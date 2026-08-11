package service

import (
	"autocft/internal/model"
	"reflect"
	"testing"
)

func TestSortConfigsWildcardLast(t *testing.T) {
	in := []*model.IngressConfig{
		{Hostname: "api.example.com", Path: "/"},
		{Hostname: "*.example.com", Path: "/"},
		{Hostname: "app.example.com", Path: "/"},
		{Hostname: "*.other.com", Path: "/"},
	}
	sortConfigs(in)
	want := []string{"api.example.com", "app.example.com", "*.example.com", "*.other.com"}
	got := make([]string, 0, len(in))
	for _, c := range in {
		got = append(got, c.Hostname)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sortConfigs = %v, want %v", got, want)
	}
}
