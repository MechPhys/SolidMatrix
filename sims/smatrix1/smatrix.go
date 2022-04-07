// Copyright (c) 2022, The MechPhys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// SMatrix is the SolidMatrix model
package main

import (
	"math/rand"

	"github.com/emer/emergent/edge"
	"github.com/emer/emergent/egui"
	"github.com/emer/emergent/elog"
	"github.com/emer/emergent/estats"
	"github.com/emer/etable/etensor"
	"github.com/emer/etable/etview"
	_ "github.com/emer/etable/etview" // include to get gui views
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/mat32"
)

func main() {
	TheSim.Config()
	gimain.Main(func() { // this starts gui -- requires valid OpenGL display connection (e.g., X11)
		guirun()
	})
}

func guirun() {
	win := TheSim.ConfigGui()
	win.StartEventLoop()
}

// LogPrec is precision for saving float values in logs
const LogPrec = 4

// Params holds the particle params
type Params struct {
	Dims      int `desc:"number of dimensions"`
	Size      int `desc:"size of world along each dim"`
	Particles int `desc:"n of particles"`
	Mass      float32
	Trace     float32 `desc:"1-rate of decay for display trace"`
}

func (pr *Params) Defaults() {
	pr.Dims = 2
	pr.Size = 100
	pr.Particles = 1
	pr.Mass = 1
	pr.Trace = 0.9
}

// Sim holds the params, table, etc
type Sim struct {
	Params    Params       `view:"inline" desc:"main params"`
	Mom0      mat32.Vec2   `desc:"initial momemntum"`
	TimeSteps int          `desc:"total number of time steps to take"`
	State     estats.Stats `desc:"state and stats"`
	Logs      elog.Logs    `desc:"logs"`
	GUI       egui.GUI     `view:"-" desc:"manages all the gui elements"`
	NoGui     bool         `view:"-" desc:"if true, runing in no GUI mode"`
}

// TheSim is the overall state for this simulation
var TheSim Sim

func (ss *Sim) Defaults() {
	ss.Params.Defaults()
}

// Config configures all the elements using the standard functions
func (ss *Sim) Config() {
	ss.Defaults()
	ss.TimeSteps = 1000
	ss.Update()
	ss.ConfigState()
	ss.ConfigLogs()
}

// Update updates computed values
func (ss *Sim) Update() {
}

// ConfigState configures tables and logs
func (ss *Sim) ConfigState() {
	ss.State.Init()
	pr := &ss.Params
	shp := []int{pr.Size, pr.Size}
	nms := []string{"Y", "X"}
	cur := etensor.NewFloat32(shp, nil, nms)
	ss.State.SetF32Tensor("Cur", cur)
	tr := etensor.NewFloat32(shp, nil, nms)
	ss.State.SetF32Tensor("Tr", tr)
	pos := etensor.NewInt([]int{pr.Particles, pr.Dims}, nil, []string{"Ps", "Dims"})
	ss.State.SetIntTensor("Pos", pos)
	mom := etensor.NewFloat32([]int{pr.Particles, pr.Dims}, nil, []string{"Ps", "Dims"})
	ss.State.SetF32Tensor("Mom", mom)
}

// ConfigLogs configures tables and logs
func (ss *Sim) ConfigLogs() {
	ss.Logs.CreateTables()
}

// Init inits
func (ss *Sim) Init() {
	pr := &ss.Params
	cur := ss.State.F32Tensor("Cur")
	tr := ss.State.F32Tensor("Tr")
	pos := ss.State.IntTensor("Pos")
	mom := ss.State.F32Tensor("Mom")

	cur.SetZeros()
	tr.SetZeros()
	pos.Set([]int{0, 0}, pr.Size/2)
	pos.Set([]int{0, 1}, pr.Size/2)
	mom.Set([]int{0, 0}, ss.Mom0.X)
	mom.Set([]int{0, 1}, ss.Mom0.Y)
}

func Move(v float32) int {
	e := 0.5 * (1.0 + v*v)
	p1 := 0.5 * (e + v)
	p0 := 1 - e
	// pm := 0.5 * (e - v)
	rv := rand.Float32()
	var s int
	if rv < p1 {
		s = 1
	} else if rv < p1+p0 {
		s = 0
	} else {
		s = -1
	}
	// 	fmt.Printf("v: %g  e: %g  p1: %g  p0: %g  pm: %g  s: %d\n", v, e, p1, p0, pm, s)
	return s
}

func UpdtPos(p, s, sz int) int {
	p, _ = edge.Edge(p+s, sz, true) // wrap
	return p
}

// Step runs one update
func (ss *Sim) Step() {
	pr := &ss.Params
	cur := ss.State.F32Tensor("Cur")
	tr := ss.State.F32Tensor("Tr")
	pos := ss.State.IntTensor("Pos")
	mom := ss.State.F32Tensor("Mom")

	vx := mom.Value([]int{0, 0})
	vy := mom.Value([]int{0, 1})

	sx := Move(vx)
	sy := Move(vy)

	px := pos.Value([]int{0, 0})
	py := pos.Value([]int{0, 1})

	// cv := cur.Value([]int{py, px})
	cur.Set([]int{py, px}, 0)

	px = UpdtPos(px, sx, pr.Size)
	py = UpdtPos(py, sy, pr.Size)

	pos.Set([]int{0, 0}, px)
	pos.Set([]int{0, 1}, py)

	cur.Set([]int{py, px}, 1)
	tr.Set([]int{py, px}, tr.Value([]int{py, px})+1)

	if ss.GUI.Active {
		cg := ss.GUI.Grid("Cur")
		cg.UpdateSig()
		tg := ss.GUI.Grid("Tr")
		tg.UpdateSig()
	}
}

// Stop stops running
func (ss *Sim) Run() {
	ss.GUI.StopNow = false
	for t := 0; t < ss.TimeSteps; t++ {
		ss.Step()
		if ss.GUI.StopNow {
			break
		}
	}
	ss.Stopped()
}

// Stop tells the sim to stop running
func (ss *Sim) Stop() {
	ss.GUI.StopNow = true
}

// Stopped is called when a run method stops running -- updates the IsRunning flag and toolbar
func (ss *Sim) Stopped() {
	ss.GUI.Stopped()
}

// ConfigGui configures the GoGi gui interface for this simulation,
func (ss *Sim) ConfigGui() *gi.Window {
	title := "SolidMatrix"
	ss.GUI.MakeWindow(ss, "smatrix", title, `SolidMatrix model. See <a href="https://github.com/MechPhys/SolidMatrix">GitHub</a>.</p>`)
	ss.GUI.CycleUpdateInterval = 10

	cur := ss.State.F32Tensor("Cur")
	tr := ss.State.F32Tensor("Tr")

	tg := ss.GUI.TabView.AddNewTab(etview.KiT_TensorGrid, "Cur").(*etview.TensorGrid)
	tg.SetStretchMax()
	ss.GUI.SetGrid("Cur", tg)
	tg.SetTensor(cur)

	tg = ss.GUI.TabView.AddNewTab(etview.KiT_TensorGrid, "Tr").(*etview.TensorGrid)
	tg.SetStretchMax()
	ss.GUI.SetGrid("Tr", tg)
	tg.SetTensor(tr)

	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Init", Icon: "update",
		Tooltip: "Initialize state, and start over.  Also applies current params.",
		Active:  egui.ActiveStopped,
		Func: func() {
			ss.Init()
			ss.GUI.UpdateWindow()
		},
	})
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Run",
		Icon:    "run",
		Tooltip: "Runs n steps.",
		Active:  egui.ActiveStopped,
		Func: func() {
			if !ss.GUI.IsRunning {
				ss.GUI.IsRunning = true
				ss.GUI.ToolBar.UpdateActions()
				go ss.Run()
			}
		},
	})
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Stop",
		Icon:    "stop",
		Tooltip: "Interrupts running.  Hitting Train again will pick back up where it left off.",
		Active:  egui.ActiveRunning,
		Func: func() {
			ss.Stop()
		},
	})
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "Step",
		Icon:    "step-fwd",
		Tooltip: "Advances one step at a time.",
		Active:  egui.ActiveStopped,
		Func: func() {
			if !ss.GUI.IsRunning {
				ss.GUI.IsRunning = true
				ss.Step()
				ss.GUI.IsRunning = false
				ss.GUI.UpdateWindow()
			}
		},
	})
	ss.GUI.AddToolbarItem(egui.ToolbarItem{Label: "README",
		Icon:    "file-markdown",
		Tooltip: "Opens your browser on the README file that contains instructions for how to run this model.",
		Active:  egui.ActiveAlways,
		Func: func() {
			gi.OpenURL("https://github.com/lvis/blob/main/sims/lvis_cu3d100_te16deg/README.md")
		},
	})
	ss.GUI.FinalizeGUI(false)
	return ss.GUI.Win
}
