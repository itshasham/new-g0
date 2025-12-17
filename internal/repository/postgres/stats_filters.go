package postgres

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var pageIdentRE = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

type jsonFilterCondition struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// buildPagesWherePostgres builds a parameterized WHERE clause fragment for the pages table.
// It supports two filter encodings:
//  1) equality maps: [{"response_code":200,"depth":1}]
//  2) filter groups: [{"filters":[{"name":"response_code","operator":"gte","value":400}]}]
func buildPagesWherePostgres(sessionID int64, prefilters, filters []map[string]any) (string, []any, error) {
	where := "crawling_session_id = $1"
	args := []any{sessionID}
	next := 2

	var err error
	where, args, next, err = appendFilterMapsPostgres(where, args, next, prefilters)
	if err != nil {
		return "", nil, err
	}
	where, args, next, err = appendFilterMapsPostgres(where, args, next, filters)
	if err != nil {
		return "", nil, err
	}
	return where, args, nil
}

func buildProblematicClausePostgres(rawFilterConfigs [][]byte) (string, []any, error) {
	var orParts []string
	var args []any
	next := 1

	for _, raw := range rawFilterConfigs {
		var cfg map[string]any
		if err := json.Unmarshal(raw, &cfg); err != nil {
			continue
		}
		filterGroupsRaw, ok := cfg["filter_groups"]
		if !ok {
			continue
		}
		groups := normalizeFilterGroups(filterGroupsRaw)
		for _, g := range groups {
			var groupClause string
			var groupArgs []any
			var err error
			if _, ok := g["filters"]; ok {
				groupClause, groupArgs, err = buildGroupClausePostgres(g, next)
			} else {
				groupClause, groupArgs, err = buildEqualityClausePostgres(g, next)
			}
			if err != nil || groupClause == "" {
				continue
			}
			next += len(groupArgs)
			args = append(args, groupArgs...)
			orParts = append(orParts, "("+groupClause+")")
		}
	}

	if len(orParts) == 0 {
		return "", nil, nil
	}
	return strings.Join(orParts, " OR "), args, nil
}

func buildEqualityClausePostgres(m map[string]any, start int) (string, []any, error) {
	var parts []string
	var args []any
	next := start
	for key, value := range m {
		col, ok := allowedPageColumn(key)
		if !ok {
			return "", nil, fmt.Errorf("invalid filter column: %s", key)
		}
		if value == nil {
			parts = append(parts, fmt.Sprintf("%s IS NULL", col))
			continue
		}
		parts = append(parts, fmt.Sprintf("%s = $%d", col, next))
		args = append(args, value)
		next++
	}
	return strings.Join(parts, " AND "), args, nil
}

func normalizeFilterGroups(raw any) []map[string]any {
	rawSlice, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(rawSlice))
	for _, item := range rawSlice {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
			continue
		}
	}
	return out
}

func appendFilterMapsPostgres(where string, args []any, next int, filters []map[string]any) (string, []any, int, error) {
	for _, filter := range filters {
		if filter == nil {
			continue
		}

		if _, ok := filter["filters"]; ok {
			groupClause, groupArgs, err := buildGroupClausePostgres(filter, next)
			if err != nil {
				return "", nil, 0, err
			}
			if groupClause != "" {
				where += " AND (" + groupClause + ")"
				args = append(args, groupArgs...)
				next += len(groupArgs)
			}
			continue
		}

		for key, value := range filter {
			col, ok := allowedPageColumn(key)
			if !ok {
				return "", nil, 0, fmt.Errorf("invalid filter column: %s", key)
			}

			if value == nil {
				where += fmt.Sprintf(" AND %s IS NULL", col)
				continue
			}

			where += fmt.Sprintf(" AND %s = $%d", col, next)
			args = append(args, value)
			next++
		}
	}
	return where, args, next, nil
}

func buildGroupClausePostgres(group map[string]any, start int) (string, []any, error) {
	raw, ok := group["filters"]
	if !ok {
		return "", nil, nil
	}
	rawSlice, ok := raw.([]any)
	if !ok {
		return "", nil, errors.New("invalid filters format")
	}

	var parts []string
	var args []any
	next := start

	for _, item := range rawSlice {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		cond := jsonFilterCondition{}
		if v, ok := m["name"].(string); ok {
			cond.Name = v
		}
		if v, ok := m["operator"].(string); ok {
			cond.Operator = strings.ToLower(strings.TrimSpace(v))
		}
		if cond.Operator == "" {
			cond.Operator = "eq"
		}
		cond.Value = m["value"]

		col, ok := allowedPageColumn(cond.Name)
		if !ok {
			return "", nil, fmt.Errorf("invalid filter column: %s", cond.Name)
		}

		clause, clauseArgs, err := buildConditionPostgres(col, cond.Operator, cond.Value, next)
		if err != nil {
			return "", nil, err
		}
		if clause == "" {
			continue
		}
		parts = append(parts, clause)
		args = append(args, clauseArgs...)
		next += len(clauseArgs)
	}

	return strings.Join(parts, " AND "), args, nil
}

func buildConditionPostgres(col string, op string, value any, start int) (string, []any, error) {
	switch op {
	case "eq":
		if value == nil {
			return col + " IS NULL", nil, nil
		}
		return fmt.Sprintf("%s = $%d", col, start), []any{value}, nil
	case "neq":
		if value == nil {
			return col + " IS NOT NULL", nil, nil
		}
		return fmt.Sprintf("%s <> $%d", col, start), []any{value}, nil
	case "gt":
		return fmt.Sprintf("%s > $%d", col, start), []any{value}, nil
	case "gte":
		return fmt.Sprintf("%s >= $%d", col, start), []any{value}, nil
	case "lt":
		return fmt.Sprintf("%s < $%d", col, start), []any{value}, nil
	case "lte":
		return fmt.Sprintf("%s <= $%d", col, start), []any{value}, nil
	case "isnull":
		return col + " IS NULL", nil, nil
	case "notnull":
		return col + " IS NOT NULL", nil, nil
	case "contains":
		return fmt.Sprintf("%s ILIKE '%%' || $%d || '%%'", col, start), []any{value}, nil
	case "in":
		vals, ok := toAnySlice(value)
		if !ok || len(vals) == 0 {
			return "", nil, nil
		}
		ph := make([]string, len(vals))
		args := make([]any, 0, len(vals))
		for i, v := range vals {
			ph[i] = fmt.Sprintf("$%d", start+i)
			args = append(args, v)
		}
		return fmt.Sprintf("%s IN (%s)", col, strings.Join(ph, ",")), args, nil
	default:
		return "", nil, fmt.Errorf("unsupported operator: %s", op)
	}
}

func toAnySlice(v any) ([]any, bool) {
	switch t := v.(type) {
	case []any:
		return t, true
	default:
		return nil, false
	}
}

func allowedPageColumn(name string) (string, bool) {
	if !pageIdentRE.MatchString(name) {
		return "", false
	}
	return name, true
}
