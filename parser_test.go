package querystring

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/go-test/deep"
)

func TestParser(t *testing.T) {
	t.Run("WithoutOperator", func(t *testing.T) {
		cond, err := Parse("a b c")
		require.NoError(t, err)

		expected := Condition(&OrCondition{
			Left: &MatchCondition{Value: "a"},
			Right: &OrCondition{
				Left: &MatchCondition{
					Value: "b",
				},
				Right: &MatchCondition{
					Value: "c",
				},
			},
		})
		spew.Dump(cond)
		diff := deep.Equal(cond, expected)
		require.Nil(t, diff)
	})

	t.Run("WithNested", func(t *testing.T) {
		cond, err := Parse("(a NOT (b OR c)) OR (d)")
		require.NoError(t, err)

		expected := Condition(&OrCondition{
			Left: &OrCondition{
				Left: &MatchCondition{Value: "a"},
				Right: &NotCondition{Condition: &OrCondition{
					Left:  &MatchCondition{Value: "b"},
					Right: &MatchCondition{Value: "c"},
				}},
			},
			Right: &MatchCondition{
				Value: "d",
			},
		})
		spew.Dump(cond)
		diff := deep.Equal(cond, expected)
		require.Nil(t, diff)
	})

	t.Run("WithQuote", func(t *testing.T) {
		cond, err := Parse("\"a b c\" d")
		require.NoError(t, err)

		expected := Condition(&OrCondition{
			Left: &MatchCondition{Value: "a b c"},
			Right: &MatchCondition{
				Value: "d",
			},
		})
		spew.Dump(cond)
		diff := deep.Equal(cond, expected)
		require.Nil(t, diff)
	})

	t.Run("WithNot", func(t *testing.T) {
		cond, err := Parse("a NOT b")
		require.NoError(t, err)

		expected := Condition(&OrCondition{
			Left: &MatchCondition{
				Value: "a",
			},
			Right: &NotCondition{
				Condition: &MatchCondition{
					Value: "b",
				},
			},
		})
		spew.Dump(cond)
		diff := deep.Equal(cond, expected)
		require.Nil(t, diff)
	})

	t.Run("WithNotAndParentheses", func(t *testing.T) {
		cond, err := Parse("(a NOT b)")
		require.NoError(t, err)

		expected := Condition(&OrCondition{
			Left: &MatchCondition{
				Value: "a",
			},
			Right: &NotCondition{
				Condition: &MatchCondition{
					Value: "b",
				},
			},
		})
		spew.Dump(cond)
		diff := deep.Equal(cond, expected)
		require.Nil(t, diff)
	})

	t.Run("WithOr", func(t *testing.T) {
		cond, err := Parse("a OR b OR c")
		require.NoError(t, err)

		expected := Condition(&OrCondition{
			Left: &MatchCondition{
				Value: "a",
			},
			Right: &OrCondition{
				Left:  &MatchCondition{Value: "b"},
				Right: &MatchCondition{Value: "c"},
			},
		})
		spew.Dump(cond)
		diff := deep.Equal(cond, expected)
		require.Nil(t, diff)
	})

	t.Run("WithAnd", func(t *testing.T) {
		cond, err := Parse("a AND b AND c")
		require.NoError(t, err)

		expected := Condition(&AndCondition{
			Left: &MatchCondition{
				Value: "a",
			},
			Right: &AndCondition{
				Left:  &MatchCondition{Value: "b"},
				Right: &MatchCondition{Value: "c"},
			},
		})
		spew.Dump(cond)
		diff := deep.Equal(cond, expected)
		require.Nil(t, diff)
	})

	t.Run("WithParentheses", func(t *testing.T) {
		cond, err := Parse("a NOT (b OR c)")
		require.NoError(t, err)

		expected := Condition(&OrCondition{
			Left: &MatchCondition{
				Value: "a",
			},
			Right: &NotCondition{
				Condition: &OrCondition{
					Left:  &MatchCondition{Value: "b"},
					Right: &MatchCondition{Value: "c"},
				},
			},
		})
		spew.Dump(cond)
		diff := deep.Equal(cond, expected)
		require.Nil(t, diff)
	})

	t.Run("WithPrecedence", func(t *testing.T) {
		cond, err := Parse("a OR b AND NOT c")
		require.NoError(t, err)

		expected := Condition(&OrCondition{
			Left: &MatchCondition{
				Value: "a",
			},
			Right: &AndCondition{
				Left: &MatchCondition{
					Value: "b",
				},
				Right: &NotCondition{
					Condition: &MatchCondition{
						Value: "c",
					},
				},
			},
		})
		spew.Dump(cond)
		diff := deep.Equal(cond, expected)
		require.Nil(t, diff)
	})

	t.Run("Orig", func(t *testing.T) {
		cond, err := Parse("message: test\\ value AND datetime: [\"2020-01-01T00:00:00\" TO \"2020-12-31T00:00:00\"]")
		if err != nil {
			t.Errorf("parse return error, %s", err)
			return
		}

		expected := Condition(&AndCondition{
			Left: &MatchCondition{
				Field: "message",
				Value: "test value",
			},
			Right: &TimeRangeCondition{
				Field:        "datetime",
				Start:        pointer("2020-01-01T00:00:00"),
				End:          pointer("2020-12-31T00:00:00"),
				IncludeStart: true,
				IncludeEnd:   true,
			},
		})

		if diff := deep.Equal(cond, expected); diff != nil {
			t.Errorf("returned condition unexpected: diff= %s", diff)
			return
		}
	})

}

func pointer(s string) *string {
	return &s
}
