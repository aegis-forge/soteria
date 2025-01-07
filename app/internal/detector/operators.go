package detector

import (
	"errors"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

var operators = []string{"Equals", "Match", "Exists"}
var logicalOperators = []string{"And", "Or"}

type Operator interface {
	Evaluate(yamlContent []byte) error
	GetValue() bool
	GetLines() []int
	GetChildren() []Operator
	ClearResults()
}

type operator struct {
	value bool
	lines []int

	NOT bool
	LHS interface{}
	RHS interface{}
}

// ================
// ==== EQUALS ====
// ================

type Equals operator

func (o *Equals) Evaluate(yamlContent []byte) error {
	if reflect.TypeOf(o.LHS).Kind() != reflect.String && reflect.TypeOf(o.RHS).Kind() != reflect.Slice {
		return errors.New("expected strings as RHS and LHS")
	}

	toCompare, lines, err := Resolve(o.LHS.(string), yamlContent)

	if err != nil {
		return err
	}

	for ind, res := range toCompare {
		if res == o.RHS {
			o.lines = append(o.lines, lines[ind])
		}
	}

	if len(o.lines) > 0 {
		o.value = true
	}

	return nil
}

func (o *Equals) GetValue() bool {
	return o.value
}

func (o *Equals) GetLines() []int {
	return o.lines
}

func (o *Equals) GetChildren() []Operator {
	return nil
}

func (o *Equals) ClearResults() {
	o.lines = []int{}
	o.value = false
}

// ===============
// ==== MATCH ====
// ===============

type Match operator

func (o *Match) Evaluate(yamlContent []byte) error {
	toCompare, lines, err := Resolve(o.LHS.(string), yamlContent)

	if err != nil {
		return err
	}

	regex := regexp.MustCompile(o.RHS.(string))

	for ind, res := range toCompare {
		if regex.MatchString(res) {
			o.lines = append(o.lines, lines[ind])
		}
	}

	if len(o.lines) > 0 {
		o.value = true
	}

	return nil
}

func (o *Match) GetValue() bool {
	return o.value
}

func (o *Match) GetLines() []int {
	return o.lines
}

func (o *Match) GetChildren() []Operator {
	return nil
}

func (o *Match) ClearResults() {
	o.lines = []int{}
	o.value = false
}

// ================
// ==== EXISTS ====
// ================

type Exists operator

func (o *Exists) Evaluate(yamlContent []byte) error {
	exists, err := CheckExistence(o.LHS.(string), o.RHS.(string), yamlContent)

	if err != nil {
		return err
	}

	if o.NOT {
		o.value = !exists
	} else {
		o.value = exists
	}

	if o.value {
		o.lines = []int{1}
	}

	return nil
}

func (o *Exists) GetValue() bool {
	return o.value
}

func (o *Exists) GetLines() []int {
	return o.lines
}

func (o *Exists) GetChildren() []Operator {
	return nil
}

func (o *Exists) ClearResults() {
	o.value = false
	o.lines = []int{}
}

// =============
// ==== AND ====
// =============

type And operator

func (o *And) Evaluate(yamlContent []byte) error {
	lhs := o.LHS.(Operator)
	rhs := o.RHS.(Operator)

	lhsType := strings.Split(reflect.TypeOf(lhs).String(), ".")
	rhsType := strings.Split(reflect.TypeOf(rhs).String(), ".")

	if slices.Contains(operators, lhsType[len(lhsType)-1]) {
		if err := lhs.Evaluate(yamlContent); err != nil {
			return err
		}
	}

	if slices.Contains(operators, rhsType[len(rhsType)-1]) {
		if err := rhs.Evaluate(yamlContent); err != nil {
			return err
		}
	}

	o.value = lhs.GetValue() && rhs.GetValue()

	if o.value {
		for _, line := range lhs.GetLines() {
			o.lines = append(o.lines, line)
		}

		for _, line := range rhs.GetLines() {
			o.lines = append(o.lines, line)
		}
	}

	return nil
}

func (o *And) GetValue() bool {
	return o.value
}

func (o *And) GetLines() []int {
	slices.Sort(o.lines)
	return slices.Compact(o.lines)
}

func (o *And) GetChildren() []Operator {
	return []Operator{o.RHS.(Operator), o.LHS.(Operator)}
}

func (o *And) ClearResults() {
	if len(o.GetChildren()) > 0 {
		for _, child := range o.GetChildren() {
			child.ClearResults()
		}
	}

	o.lines = []int{}
	o.value = false
}

// ============
// ==== OR ====
// ============

type Or operator

func (o *Or) Evaluate(yamlContent []byte) error {
	lhs := o.LHS.(Operator)
	rhs := o.RHS.(Operator)

	lhsType := strings.Split(reflect.TypeOf(lhs).String(), ".")
	rhsType := strings.Split(reflect.TypeOf(rhs).String(), ".")

	if slices.Contains(operators, lhsType[len(lhsType)-1]) {
		if err := lhs.Evaluate(yamlContent); err != nil {
			return err
		}
	}

	if slices.Contains(operators, rhsType[len(rhsType)-1]) {
		if err := rhs.Evaluate(yamlContent); err != nil {
			return err
		}
	}

	o.value = lhs.GetValue() || rhs.GetValue()

	if o.value {
		for _, line := range lhs.GetLines() {
			o.lines = append(o.lines, line)
		}

		for _, line := range rhs.GetLines() {
			o.lines = append(o.lines, line)
		}
	}

	return nil
}

func (o *Or) GetValue() bool {
	return o.value
}

func (o *Or) GetLines() []int {
	slices.Sort(o.lines)
	return slices.Compact(o.lines)
}

func (o *Or) GetChildren() []Operator {
	return []Operator{o.RHS.(Operator), o.LHS.(Operator)}
}

func (o *Or) ClearResults() {
	if len(o.GetChildren()) > 0 {
		for _, child := range o.GetChildren() {
			child.ClearResults()
		}
	}

	o.lines = []int{}
	o.value = false
}
