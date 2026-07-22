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

**Reserved characters:** `(`, `)`, `|`, `*` — use grouping and operators only in their syntactic roles.

Examples:

```
a(b|c)*     matches a, ab, ac, abb, …
(01)*1      matches 1, 01, 001, 0101, …
a*|b*       matches any string of only a's, only b's, or empty
```

## Project layout

```
regex/
├── cmd/matcher/              CLI: two regex args → true/false
├── internal/regexEval/
│   ├── ast.go                AST node types
│   ├── parser.go             Recursive-descent parser
│   ├── nfa.go                Thompson construction, ε-NFA, product search
│   └── main.go               Entry API: HasUnion
└── go.mod
```

## Status

| Stage | Status |
|-------|--------|
| Parser → AST | Done |
| Thompson construction | Done |
| Product automaton search | Done |
| CLI | Done |

## Usage

**Library** (within this module — `internal/regexEval` is not importable outside `regex/`):

```go
package main

import (
    "fmt"
    rx "regex/internal/regexEval"
)

func main() {
    ok, errs := rx.HasUnion(`a(b|c)*`, `ab*`)
    if len(errs) > 0 {
        panic(errs)
    }
    fmt.Println(ok) // true — both match "a", "ab", "abb", …
}
```

**CLI:**

```bash
go run ./cmd/matcher 'a(b|c)*' 'ab*'
# true

go run ./cmd/matcher 'ab' 'c*'
# false
```

`HasUnion` returns `(bool, []string)`: the boolean is whether the languages intersect; the slice holds parse errors (empty on success).

## Development

Requires Go 1.24+.

```bash
go build ./...
go run ./cmd/matcher 'a|b' 'b|c'
```

There are no automated tests yet; add `_test.go` files under `internal/regexEval` as you extend the language.

## How it works

**Thompson construction** (`thompsonConstruction`) compiles each AST node into an ε-NFA fragment with a single start and a single accept state:

- Literal → two states, one labeled transition
- Concatenation → ε-link from the left accept to the right start
- Alternation → new start ε-branches to both fragments; both accepts ε-link to a new accept
- Kleene star → ε-skip to accept, ε-enter fragment, loop back, ε-exit to accept
- Grouping → compile the inner node (no extra states)

States store labeled moves in `Transitions map[rune]*State` and ε-moves in `Epsilons map[*State]struct{}`.

**Product search** (`evaluateNfaProduct`) asks whether some pair `(q, p)` is reachable where `q` is accepting in the first NFA and `p` in the second. At each step it:

1. Takes the ε-closure of the current state pair (`flattenEpsilonTransitions` on each side).
2. If any pair in the product of those closures is `(Accept₁, Accept₂)`, returns `true`.
3. Otherwise, for every symbol that appears on a transition from both sides, recurses on the pair of target states in sync.
4. Uses a visited set on `(stateA, stateB)` to avoid infinite ε-loops.

Because both NFAs may have ε-transitions, reachability is ε-closure-aware rather than plain labeled-edge BFS alone.
