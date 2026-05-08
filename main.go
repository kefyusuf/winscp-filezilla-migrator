package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"

	"github.com/kefyusuf/winscp-filezilla-migrator/domain/exporter"
	"github.com/kefyusuf/winscp-filezilla-migrator/domain/models"
	"github.com/kefyusuf/winscp-filezilla-migrator/domain/parser"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("WinSCP to FileZilla Migrator"))
		w.Option(app.Size(unit.Dp(900), unit.Dp(650)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type GUI struct {
	win      *app.Window
	explorer *explorer.Explorer

	page     int
	sessions []models.Session
	selected int
	status   string

	openBtn    widget.Clickable
	exportBtn  widget.Clickable
	backBtn    widget.Clickable
	serverList widget.List
	clickables []widget.Clickable
	listSelect []widget.Selectable

	detailSel [6]widget.Selectable
}

const (
	pagePick = 0
	pageList = 1
)

func loop(w *app.Window) error {
	gui := &GUI{
		win:        w,
		explorer:   explorer.NewExplorer(w),
		serverList: widget.List{List: layout.List{Axis: layout.Vertical}},
		selected:   -1,
	}

	th := material.NewTheme()
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			gui.explorer.ListenEvents(e)
			gui.handleClicks(gtx)
			gui.layout(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}

func (g *GUI) handleClicks(gtx layout.Context) {
	for g.openBtn.Clicked(gtx) {
		go func() {
			f, err := g.explorer.ChooseFile(".ini")
			if err != nil {
				return
			}
			defer f.Close()

			data, err := io.ReadAll(f)
			if err != nil {
				g.status = "Error reading file: " + err.Error()
				g.win.Invalidate()
				return
			}

			tmpDir := os.TempDir()
			tmpPath := filepath.Join(tmpDir, "winscp_migrator_temp.ini")
			if err := os.WriteFile(tmpPath, data, 0644); err != nil {
				g.status = "Error saving temp file: " + err.Error()
				g.win.Invalidate()
				return
			}
			defer os.Remove(tmpPath)

			sessions, err := parser.ParseWinSCPIni(tmpPath)
			if err != nil {
				g.status = "Parse error: " + err.Error()
				g.win.Invalidate()
				return
			}
			if len(sessions) == 0 {
				g.status = "No sessions found in INI file"
				g.win.Invalidate()
				return
			}

			sort.Slice(sessions, func(i, j int) bool {
				return sessions[i].Name < sessions[j].Name
			})

			g.sessions = sessions
			g.clickables = make([]widget.Clickable, len(sessions))
			g.listSelect = make([]widget.Selectable, len(sessions))
			g.selected = 0
			g.page = pageList
			g.status = fmt.Sprintf("Loaded %d sessions", len(sessions))
			g.win.Invalidate()
		}()
	}

	for g.exportBtn.Clicked(gtx) {
		go func() {
			w, err := g.explorer.CreateFile("sites.xml")
			if err != nil {
				return
			}
			defer w.Close()

			tmpPath := filepath.Join(os.TempDir(), "winscp_migrator_export.xml")
			if err := exporter.ExportToFileZilla(g.sessions, tmpPath); err != nil {
				g.status = "Export error: " + err.Error()
				g.win.Invalidate()
				return
			}
			defer os.Remove(tmpPath)

			data, err := os.ReadFile(tmpPath)
			if err != nil {
				g.status = "Read error: " + err.Error()
				g.win.Invalidate()
				return
			}
			if _, err := w.Write(data); err != nil {
				g.status = "Write error: " + err.Error()
				g.win.Invalidate()
				return
			}
			g.status = "Exported successfully!"
			g.win.Invalidate()
		}()
	}

	for g.backBtn.Clicked(gtx) {
		g.page = pagePick
		g.status = ""
		g.sessions = nil
		g.clickables = nil
		g.listSelect = nil
		g.selected = -1
	}

	for i := range g.clickables {
		if g.clickables[i].Clicked(gtx) {
			g.selected = i
		}
	}
}

func (g *GUI) layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return g.layoutHeader(gtx, th)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if g.page == pagePick {
				return g.layoutPickPage(gtx, th)
			}
			return g.layoutListPage(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return g.layoutStatus(gtx, th)
		}),
	)
}

func (g *GUI) layoutHeader(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(4), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return material.H5(th, "WinSCP to FileZilla Migrator").Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if g.page == pageList {
					return material.Button(th, &g.backBtn, "Back").Layout(gtx)
				}
				return layout.Dimensions{}
			}),
		)
	})
}

func (g *GUI) layoutStatus(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if g.status == "" {
		return layout.Dimensions{}
	}
	return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(8), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Body1(th, g.status).Layout(gtx)
	})
}

func (g *GUI) layoutPickPage(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.H4(th, "Select your WinSCP.ini file").Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &g.openBtn, "Open WinSCP.ini")
				btn.Inset = layout.UniformInset(unit.Dp(16))
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				appdata := os.Getenv("APPDATA")
				if appdata == "" {
					appdata = "~/.winscp"
				}
				return material.Body2(th, "Typically found at: "+appdata+"\\WinSCP.ini").Layout(gtx)
			}),
		)
	})
}

func (g *GUI) layoutListPage(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
				return g.layoutServerList(gtx, th)
			}),
			layout.Flexed(0.6, func(gtx layout.Context) layout.Dimensions {
				return g.layoutDetail(gtx, th)
			}),
		)
	})
}

func (g *GUI) layoutServerList(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return material.List(th, &g.serverList).Layout(gtx, len(g.sessions), func(gtx layout.Context, i int) layout.Dimensions {
		s := g.sessions[i]
		selected := g.selected == i

		return g.clickables[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, s.Name)
				if selected {
					lbl.Font.Weight = font.Bold
				}
				lbl.State = &g.listSelect[i]
				return lbl.Layout(gtx)
			})
		})
	})
}

func (g *GUI) layoutDetail(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if g.selected < 0 || g.selected >= len(g.sessions) {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, "Select a session to view details").Layout(gtx)
		})
	}

	s := g.sessions[g.selected]

		return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.H6(th, s.Name).Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return g.layoutDetailRowSel(gtx, th, "Host", s.HostName, &g.detailSel[0])
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return g.layoutDetailRowSel(gtx, th, "User", s.UserName, &g.detailSel[1])
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return g.layoutDetailRowSel(gtx, th, "Port", s.PortNumber, &g.detailSel[2])
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return g.layoutDetailRowSel(gtx, th, "Protocol", g.mapProtocol(s.FSProtocol), &g.detailSel[3])
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return g.layoutDetailRowSel(gtx, th, "Remote Dir", s.RemoteDir, &g.detailSel[4])
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return g.layoutDetailRowSel(gtx, th, "Local Dir", s.LocalDir, &g.detailSel[5])
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &g.exportBtn, "Export to FileZilla XML").Layout(gtx)
			}),
		)
	})
}

func (g *GUI) layoutDetailRowSel(gtx layout.Context, th *material.Theme, label, value string, sel *widget.Selectable) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, label+":")
				lbl.Font.Weight = font.Bold
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, value)
				lbl.State = sel
				return lbl.Layout(gtx)
			}),
		)
	})
}

func (g *GUI) mapProtocol(fsProtocol string) string {
	if fsProtocol == "2" {
		return "SFTP"
	}
	return "FTP"
}
