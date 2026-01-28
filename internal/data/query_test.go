package data

import (
	"reflect"
	"testing"
)

func TestQuery(t *testing.T) {
	tests := []struct {
		name       string
		data       any
		expression string
		want       any
		wantErr    bool
	}{
		{
			name:       "extract single field",
			data:       map[string]any{"name": "Alice", "age": 30},
			expression: ".name",
			want:       "Alice",
		},
		{
			name: "extract nested field",
			data: map[string]any{
				"user": map[string]any{
					"profile": map[string]any{
						"name": "Bob",
					},
				},
			},
			expression: ".user.profile.name",
			want:       "Bob",
		},
		{
			name: "construct new object",
			data: map[string]any{
				"items": []any{"a", "b", "c"},
			},
			expression: "{count: (.items | length), first: .items[0]}",
			want:       map[string]any{"count": 3, "first": "a"},
		},
		{
			name:       "extract first element of array",
			data:       []any{1, 2, 3},
			expression: ".[0]",
			want:       1,
		},
		{
			name:       "invalid expression",
			data:       map[string]any{},
			expression: ".name[",
			wantErr:    true,
		},
		{
			name: "expression returns no result",
			data: map[string]any{
				"items": []any{},
			},
			expression: ".items[]",
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query(tt.data, tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestQueryAll(t *testing.T) {
	tests := []struct {
		name       string
		data       any
		expression string
		want       []any
		wantErr    bool
	}{
		{
			name: "extract array elements",
			data: map[string]any{
				"users": []any{
					map[string]any{"id": 1, "name": "Alice"},
					map[string]any{"id": 2, "name": "Bob"},
				},
			},
			expression: ".users[]",
			want: []any{
				map[string]any{"id": 1, "name": "Alice"},
				map[string]any{"id": 2, "name": "Bob"},
			},
		},
		{
			name: "filter with select",
			data: map[string]any{
				"items": []any{
					map[string]any{"id": "a", "active": true},
					map[string]any{"id": "b", "active": false},
					map[string]any{"id": "c", "active": true},
				},
			},
			expression: ".items[] | select(.active)",
			want: []any{
				map[string]any{"id": "a", "active": true},
				map[string]any{"id": "c", "active": true},
			},
		},
		{
			name: "deeply nested extraction",
			data: map[string]any{
				"org": map[string]any{
					"dept": map[string]any{
						"team": map[string]any{
							"members": []any{"Alice", "Bob"},
						},
					},
				},
			},
			expression: ".org.dept.team.members[]",
			want:       []any{"Alice", "Bob"},
		},
		{
			name: "empty array returns empty result",
			data: map[string]any{
				"items": []any{},
			},
			expression: ".items[]",
			want:       nil,
		},
		{
			name:       "invalid expression",
			data:       map[string]any{},
			expression: ".items[",
			wantErr:    true,
		},
		{
			name: "transform each element",
			data: map[string]any{
				"numbers": []any{1, 2, 3},
			},
			expression: ".numbers[] | . * 2",
			want:       []any{2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryAll(tt.data, tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
