// Copyright (c) 2022, The MechPhys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// SMatrix is the SolidMatrix model
package main

import (
	"fmt"
	"math"
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
	"github.com/goki/ki/ints"
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
	Dims      int     `desc:"number of dimensions"`
	Size      int     `desc:"size of world along each dim"`
	StatWin   int     `desc:"time window in which to track stats"`
	Particles int     `desc:"n of particles"`
	Mass      float32 `desc:"particle mass"`
	Trace     float32 `desc:"1-rate of decay for display trace"`
}

func (pr *Params) Defaults() {
	pr.Dims = 2
	pr.Size = 2000
	pr.StatWin = 100
	pr.Particles = 1
	pr.Mass = 1
	pr.Trace = 0.9
}

// Sim holds the params, table, etc
type Sim struct {
	Params   Params       `view:"no-inline" desc:"main params"`
	Mom0     mat32.Vec2   `desc:"initial momemntum"`
	NRuns    int          `desc:"number of runs resetting to start between"`
	RunSteps int          `desc:"total number of time steps to take per Run"`
	State    estats.Stats `desc:"state and stats"`
	Logs     elog.Logs    `desc:"logs"`
	GUI      egui.GUI     `view:"-" desc:"manages all the gui elements"`
	NoGui    bool         `view:"-" desc:"if true, runing in no GUI mode"`
}

// TheSim is the overall state for this simulation
var TheSim Sim

func (ss *Sim) Defaults() {
	ss.Params.Defaults()
	ss.Mom0.Set(0.5, 0)
}

// Config configures all the elements using the standard functions
func (ss *Sim) Config() {
	ss.Defaults()
	ss.NRuns = 20
	ss.RunSteps = 1000
	ss.Update()
	ss.ConfigState()
	ss.ConfigLogs()
	ss.Init()
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

	tr.SetMetaData("fix-max", "false")
	tr.SetMetaData("grid-min", "1")
	pos.SetMetaData("grid-min", "1")

	ss.State.SetInt("Run", 0)
	ss.State.SetInt("Time", 0)
	ss.State.SetInt("PosX", 0)
	ss.State.SetInt("PosY", 0)
}

// ConfigLogs configures tables and logs
func (ss *Sim) ConfigLogs() {
	ss.ConfigLogItems()
	ss.Logs.CreateTables()
	ss.Logs.SetContext(&ss.State, nil)
	ss.Logs.SetMeta(elog.Train, elog.Cycle, "XAxisCol", "Time")
}

func (ss *Sim) ConfigLogItems() {
	ss.Logs.AddItem(&elog.Item{
		Name: "Run",
		Type: etensor.INT64,
		Plot: elog.DFalse,
		Write: elog.WriteMap{
			elog.Scope(elog.Train, elog.Cycle): func(ctx *elog.Context) {
				ctx.SetStatInt("Run")
			}}})
	ss.Logs.AddItem(&elog.Item{
		Name: "Time",
		Type: etensor.INT64,
		Plot: elog.DFalse,
		Write: elog.WriteMap{
			elog.Scope(elog.Train, elog.Cycle): func(ctx *elog.Context) {
				ctx.SetStatInt("Time")
			}}})
	ss.Logs.AddItem(&elog.Item{
		Name: "PosX",
		Type: etensor.INT64,
		Plot: elog.DFalse,
		Write: elog.WriteMap{
			elog.Scope(elog.Train, elog.Cycle): func(ctx *elog.Context) {
				ctx.SetStatInt("PosX")
			}}})
	ss.Logs.AddItem(&elog.Item{
		Name: "PosY",
		Type: etensor.INT64,
		Plot: elog.DFalse,
		Write: elog.WriteMap{
			elog.Scope(elog.Train, elog.Cycle): func(ctx *elog.Context) {
				ctx.SetStatInt("PosY")
			}}})
}

// Init inits
func (ss *Sim) Init() {
	ss.InitState()
	ss.State.SetInt("Run", 0)
	ss.State.SetInt("Time", 0)
	dt := ss.Logs.Table(elog.Train, elog.Cycle)
	dt.SetNumRows(ss.NRuns * ss.RunSteps)
}

// InitState inits state between runs
func (ss *Sim) InitState() {
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
	ss.UpdtStats()
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

// UpdtStats
func (ss *Sim) UpdtStats() {
	pos := ss.State.IntTensor("Pos")
	px := pos.Value([]int{0, 0})
	py := pos.Value([]int{0, 1})
	ss.State.SetInt("PosX", px)
	ss.State.SetInt("PosY", py)
}

// Log is the main logging function, handles special things for different scopes
func (ss *Sim) Log(mode elog.EvalModes, time elog.Times, row int) {
	// dt := ss.Logs.Table(mode, time)
	//	row := dt.Rows
	ss.Logs.LogRow(mode, time, row) // also logs to file, etc
	if time == elog.Cycle {
		ss.GUI.UpdateCyclePlot(elog.Train, row)
	} else {
		ss.GUI.UpdatePlot(mode, time)
	}
}

// DoStats
func (ss *Sim) DoStats() {
	twin := ss.Params.StatWin
	dt := ss.Logs.Table(elog.Train, elog.Cycle)
	mmax := mat32.Max(ss.Mom0.X, ss.Mom0.Y)
	mmax = mat32.Max(0.2, mmax)
	maxt := int(mat32.Ceil(float32(twin) * mmax))
	nt := 2*maxt + 1

	st := etensor.NewFloat32([]int{2, nt, nt}, nil, []string{"XY", "D", "T"})

	var max float32
	for ri := 0; ri < ss.NRuns; ri++ {
		ro := ri * ss.RunSteps
		for j := twin; j < ss.RunSteps-twin; j++ {
			cx := dt.CellFloat("PosX", ro+j)
			cy := dt.CellFloat("PosY", ro+j)
			for i := -twin; i <= twin; i++ {
				d := ro + j + i
				ti := ints.ClipInt(int(mat32.Round(float32(i)*mmax))+maxt, 0, nt)
				tx := dt.CellFloat("PosX", d)
				ty := dt.CellFloat("PosY", d)
				dx := ints.ClipInt(maxt+int(math.Round(tx-cx)), 0, nt)
				mx := st.Add([]int{0, dx, ti}, 1)
				dy := ints.ClipInt(maxt+int(math.Round(ty-cy)), 0, nt)
				my := st.Add([]int{1, dy, ti}, 1)
				max = mat32.Max(max, mx)
				max = mat32.Max(max, my)
			}
		}
		fmt.Printf("Stat run: %d\n", ri)
		if ss.GUI.StopNow {
			break
		}
	}
	st.SetMetaData("max", fmt.Sprintf("%g", .8*max))
	st.SetMetaData("grid-min", "1")
	st.SetMetaData("grid-fill", "1")
	ss.State.SetF32Tensor("Auto", st)
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
	tr.Add([]int{py, px}, 1)

	ss.State.SetInt("Time", ss.State.Int("Time")+1)
	ss.UpdtStats()

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
	row := 0
	for ri := 0; ri < ss.NRuns; ri++ {
		ss.InitState()
		ss.State.SetInt("Run", ri)
		for ti := 0; ti < ss.RunSteps; ti++ {
			ss.Step()
			ss.Log(elog.Train, elog.Cycle, row)
			row++
		}
		fmt.Printf("Run run: %d\n", ri)
		if ss.GUI.StopNow {
			break
		}
	}
	ss.DoStats()
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
	ss.GUI.CycleUpdateInterval = 100

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

	ss.GUI.AddPlots(title, &ss.Logs)

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
