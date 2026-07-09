# Regex Intersection Detector

A small Go library that answers one question: **do two regular expressions match any of the same strings?**

Given regexes `R` and `S`, the tool returns `true` if there exists at least one string accepted by both — that is, if `L(R) ∩ L(S) ≠ ∅`. It returns `false` when the languages are disjoint.

The implementation is built on automata theory: parse each pattern into an AST, build ε-NFAs via Thompson construction, then search the product automaton for a reachable accepting configuration.

## Pipeline

```
regex
  ↓
parser
  ↓
AST
  ↓
Thompson construction
  ↓
ε-NFA
  ↓
product automaton search
  ↓
true / false
```

Two patterns intersect if the product of their NFAs has a path from the start pair of states to an accepting pair `(q₁, q₂)` where both `q₁` and `q₂` are accept states.

## Supported language

This project targets **classic regular expressions** — the kind that define regular languages and map cleanly to finite automata. No backreferences, lookaheads, or other features that push beyond regular languages.

| Construct | Syntax | Meaning |
|-----------|--------|---------|
| Literal | `a`, `x`, `9` | Match one character |
| Concatenation | `ab` | Match `a` then `b` |
| Alternation | `a\|b` | Match `a` or `b` |
| Kleene star | `a*` | Match zero or more `a` |
| Grouping | `(…)` | Override precedence |

**Precedence** (lowest to highest): alternation → concatenation → repetition.

**Reserved characters:** `( )
`)`, `|`, `*` — use grouping and operators only in their syntactic roles.

Examples:

```
a(b|c)*     matches a, ab, ac, abb, …
(01)*1      matches 1, 01, 001, 0101, …
a*|b*       matches any string of only a's, only b's, or empty
```

## Project layout

```
regex/
├── cmd/matcher/          CLI entry point (planned)
├── internal/regex/
│   ├── parser.go         Recursive-descent parser
│   └── ast.go            AST node types
└── go.mod
```

Planned additions under `internal/regex/`:

- `nfa.go` — Thompson construction and ε-NFA representation
- `product.go` — product automaton and intersection search
- `intersect.go` — public API: `Intersects(r, s string) (bool, error)`

## Status

| Stage | Status |
|-------|--------|
| Parser → AST | In progress |
| Thompson construction | Planned |
| Product automaton search | Planned |
| CLI / public API | Planned |

## Usage (planned)

```go
package main

import (
    "fmt"
    "regex/internal/regex"
)

func main() {
    ok, err := regex.Intersects(`a(b|c)*`, `ab*`)
    if err != nil {
        panic(err)
    }
    fmt.Println(ok) // true — both match "a", "ab", "abb", …
}
```

```bash
go run ./cmd/matcher 'a(b|c)*' 'ab*'
# true
```

## Development

Requires Go 1.24+.

```bash
go test ./...
go build ./...
```

## How it works

**Thompson construction** compiles each AST into an ε-NFA with a single start state and a single accept state. Each grammar construct maps to a small fragment:

- Literal → two states, one labeled transition
- Concatenation → connect fragments in series
- Alternation → ε-branch between two fragments
- Star → ε-loop back to the fragment

**Product automaton** — given NFAs `M₁ = (Q₁, Σ, δ₁, q₀₁, F₁)` and `M₂ = (Q₂, Σ, δ₂, q₀₂, F₂)`, build `M₁ × M₂` with states `(q, p) ∈ Q₁ × Q₂`. A pair is accepting when `q ∈ F₁` and `p ∈ F₂`. `(q₀₁, q₀₂)` is reachable in the product iff the two languages intersect.

Because both NFAs may have ε-transitions, reachability uses a subset construction or ε-closure-aware search rather than a plain BFS on labeled edges alone.
