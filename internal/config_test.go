package internal

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parseConfig(t *testing.T) {
	cases := map[string]Config{
		"testdata/test.yaml": Config{
			Interval: "5s",
			Buckets:  []float64{0, 1, 2, 5},
			Tasks: []TaskConfig{
				TaskConfig{
					Name: "a", Type: "process",
					Parameters: map[string]interface{}{
						"cmd":  "echo",
						"args": []interface{}{"hi", 1},
					},
				},
				TaskConfig{
					Name: "b", Type: "web",
					Parameters: map[string]interface{}{
						"method":      `GET`,
						"url":         "http://localhost:3000",
						"status_code": 200,
					},
				},
			},
		},
	}
	for path := range cases {
		want := cases[path]
		t.Run(path, func(t *testing.T) {
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			got, err := parseConfig(b)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
