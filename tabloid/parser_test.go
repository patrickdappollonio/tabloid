package tabloid

import (
	"reflect"
	"testing"
)

func TestTabloid_ParseHeading(t *testing.T) {
	tests := []struct {
		name    string
		heading string
		want    []Column
		wantErr bool
	}{
		{
			name:    "basic",
			heading: "NAME   READY   STATUS    %RESTART   AGE GAP",
			want: []Column{
				{
					Title:      "NAME",
					StartIndex: 0,
					EndIndex:   7,
				},
				{
					Title:      "READY",
					StartIndex: 7,
					EndIndex:   15,
				},
				{
					Title:      "STATUS",
					StartIndex: 15,
					EndIndex:   25,
				},
				{
					Title:      "%RESTART",
					StartIndex: 25,
					EndIndex:   36,
				},
				{
					Title:      "AGE GAP",
					StartIndex: 36,
					EndIndex:   -1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tabloid{}
			got, err := tr.ParseHeading(tt.heading)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("mismatched number of returned values, got %v, want %v", got, tt.want)
			}

			for i, c := range got {
				assertEqual(t, c.Title, tt.want[i].Title, "item %d title = %q, want %q", i+1, c.Title, tt.want[i].Title)
				assertEqual(t, c.StartIndex, tt.want[i].StartIndex, "item %d start index = %d, want %d", i+1, c.StartIndex, tt.want[i].StartIndex)
				assertEqual(t, c.EndIndex, tt.want[i].EndIndex, "item %d end index = %d, want %d", i+1, c.EndIndex, tt.want[i].EndIndex)
			}
		})
	}
}

func assertEqual(t *testing.T, got, want interface{}, msg string, args ...interface{}) {
	if !reflect.DeepEqual(got, want) {
		if len(args) == 0 && msg != "" {
			t.Errorf(msg)
		}

		if len(args) > 0 && msg != "" {
			t.Errorf(msg, args...)
		}

		t.Errorf("got: %v, want: %v", got, want)
	}
}
