package main

import (
	"math"
)

type Point struct {
	X, Y float64
}

func (p Point) XVal() float64            { return p.X }
func (p Point) YVal() float64            { return p.Y }
func (p Point) XErr() (float64, float64) { return math.NaN(), math.NaN() }
func (p Point) YErr() (float64, float64) { return math.NaN(), math.NaN() }

type TableSize struct {
	Name               string
	Total, Index, Data float64
}

func (c TableSize) Category() string { return c.Name }
func (c TableSize) Value() float64   { return c.Total }
func (c TableSize) Flaged() bool     { return false }

type TableSizes []*TableSize

func (s TableSizes) Len() int      { return len(s) }
func (s TableSizes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ByName implements sort.Interface by providing Less and using the Len and
// Swap methods of the embedded TableSizes value.
type ByName struct{ TableSizes }

func (s ByName) Less(i, j int) bool { return s.TableSizes[i].Name < s.TableSizes[j].Name }

// ByTotal implements sort.Interface by providing Less and using the Len and
// Swap methods of the embedded TableSizes value.
type ByTotal struct{ TableSizes }

func (s ByTotal) Less(i, j int) bool { return s.TableSizes[i].Total < s.TableSizes[j].Total }
