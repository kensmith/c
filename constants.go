package main

import (
	"math"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	_ftPerM         = 3.280839895
	_jPerFtLb       = 1.3558179483314004
	_lPerGal        = 3.785411784
	_pPerKg         = 0.45359237
	_wPerHp         = 745.699872
	_defaultMaxRand = math.MaxInt16
)

var (
	_histDirname  = filepath.Join(xdg.StateHome, "github.com", "kensmith", "c")
	_histFilename = filepath.Join(_histDirname, "history")
)
