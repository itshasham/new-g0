package postgres

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBuildPagesWherePostgres_EqualityMap(t *testing.T) {
	where, args, err := buildPagesWherePostgres(123, nil, []map[string]any{
		{"response_code": 200},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "crawling_session_id = $1 AND response_code = $2"; where != want {
		t.Fatalf("expected where %q got %q", want, where)
	}
	if !reflect.DeepEqual(args, []any{int64(123), 200}) {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildPagesWherePostgres_FilterGroup(t *testing.T) {
	where, args, err := buildPagesWherePostgres(1, nil, []map[string]any{
		{"filters": []any{map[string]any{"name": "response_code", "operator": "gte", "value": 400}}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "crawling_session_id = $1 AND (response_code >= $2)"; where != want {
		t.Fatalf("expected where %q got %q", want, where)
	}
	if !reflect.DeepEqual(args, []any{int64(1), 400}) {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildProblematicClausePostgres(t *testing.T) {
	raw, _ := json.Marshal(map[string]any{
		"filter_groups": []any{
			map[string]any{
				"filters": []any{
					map[string]any{"name": "response_code", "operator": "gte", "value": 400},
				},
			},
		},
	})
	clause, args, err := buildProblematicClausePostgres([][]byte{raw})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "(response_code >= $1)"; clause != want {
		t.Fatalf("expected clause %q got %q", want, clause)
	}
	if !reflect.DeepEqual(args, []any{float64(400)}) {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestRenumberPostgresPlaceholders(t *testing.T) {
	got := renumberPostgresPlaceholders("a = $1 AND b = $2", 5)
	if want := "a = $5 AND b = $6"; got != want {
		t.Fatalf("expected %q got %q", want, got)
	}
}
