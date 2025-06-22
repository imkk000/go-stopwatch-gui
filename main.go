package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"time"

	_ "embed"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
)

//go:embed icon.png
var iconBytes []byte

//go:embed digital-7.ttf
var fontBytes []byte

func main() {
	go newTicker()

	icon, err := loadIcon()
	if err != nil {
		log.Fatal("cannot load icon:", err)
		return
	}

	flags := g.MasterWindowFlagsFrameless | g.MasterWindowFlagsNotResizable | g.MasterWindowFlagsFloating | g.MasterWindowFlagsTransparent
	wnd = g.NewMasterWindow(title, winWidth, winHeight, flags)
	wnd.SetIcon(icon)
	wnd.SetBgColor(color.RGBA{0, 0, 0, 0})
	wnd.SetStyle(g.Style().SetColor(g.StyleColorText, color.RGBA{169, 169, 169, 255}))
	fontInfo = g.Context.FontAtlas.AddFontFromBytes(fontName, fontBytes, 12)

	wnd.Run(mainLoop)
}

func mainLoop() {
	g.SingleWindow().BringToFront()

	imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
	g.PushColorWindowBg(color.RGBA{50, 50, 70, 130})
	g.PushColorFrameBg(color.RGBA{30, 30, 60, 110})

	g.SingleWindow().Layout(
		g.Align(g.AlignCenter).To(
			g.Custom(func() {
				g.InvisibleButton().Size(winWidth, 10).Build()

				if imgui.BeginDragDropSource() {
					windowPos := imgui.WindowPos()
					wnd.SetPos(int(windowPos.X-winWidth/2), int(windowPos.Y-5))
					imgui.EndDragDropSource()
				}
			}),
			g.Row(
				g.Button("30s").OnClick(reset(30*time.Second)),
				g.Button("15m").OnClick(reset(15*time.Minute)),
				g.Button("30m").OnClick(reset(30*time.Minute)),
				g.Button("1h").OnClick(reset(time.Hour)),
			),
			g.Row(g.Style().
				SetFont(fontInfo).
				SetFontSize(24).To(g.Label(text)),
			),
		),
	)
	g.PopStyleColor()
	g.PopStyleColor()
	imgui.PopStyleVar()
}

func newTicker() {
	ticker = time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	text = fmt.Sprint("00:00")

	for range ticker.C {
		if duration == 0 {
			continue
		}
		duration -= 1 * time.Second
		min := duration / time.Minute
		sec := (duration % time.Minute) / time.Second
		text = fmt.Sprintf("%02d:%02d", min, sec)

		g.Update()
	}
}

func reset(d time.Duration) func() {
	return func() {
		duration = d
		ticker.Reset(1 * time.Second)
		wnd.SetTitle(fmt.Sprintf("%s (%s)", title, d.String()))
	}
}

func loadIcon() (image.Image, error) {
	img, err := png.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		return nil, err
	}

	return img, nil
}

const (
	winWidth  = 260
	winHeight = 140
	title     = "Stopwatch"
	fontName  = "digital-7.ttf"
)

var (
	duration time.Duration
	ticker   *time.Ticker
	text     string
	fontInfo *g.FontInfo
	wnd      *g.MasterWindow
)
