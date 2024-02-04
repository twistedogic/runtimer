package internal

import (
	"os"
	"reflect"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	cases := map[string]Config{
		"testdata/test.yaml": Config{
			Interval: "5s",
			Buckets:  []float64{0, 1, 2, 5},
			Tasks: []TaskConfig{
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
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("want: %v, got: %v", want, got)
			}
		})
	}
}
