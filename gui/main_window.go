package gui

import (
	"io/ioutil"
	"os"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/conf"
)

// InitMainWindow initialize main window
func InitMainWindow() {
	var mw *walk.MainWindow
	var treeView *walk.TreeView
	var hostConfigText *walk.TextEdit
	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Lily",
		MinSize:  Size{720, 500},
		Layout:   VBox{},
		// MenuItems: newMenuItems(mw),
		ToolBar: newToolBar(treeView),
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					newTreeView(&treeView, &hostConfigText),
					TextEdit{
						AssignTo:      &hostConfigText,
						StretchFactor: 3,
						OnKeyUp: func(key walk.Key) {
							common.Info("============================ Key up =================================")
						},
					},
				},
			},
		},
	}).Create(); err != nil {
		common.Error("Error creating main window: ", err)
		os.Exit(-1)
	}
	setWindowIcon(mw)
	setXY(mw)
	setTreeViewBackground(treeView)
	showCurrentItem(hostConfigText)
	mw.Run()
}

func setWindowIcon(mw *walk.MainWindow) {
	icon, _ := walk.NewIconFromFile("res/lily.ico")
	mw.SetIcon(icon)
}

func setXY(mw *walk.MainWindow) {
	hDC := win.GetDC(0)
	defer win.ReleaseDC(0, hDC)
	mw.SetX((int(win.GetDeviceCaps(hDC, win.HORZRES)) - mw.Width()) / 2)
	mw.SetY((int(win.GetDeviceCaps(hDC, win.VERTRES)) - mw.Height()) / 2)
}

func setTreeViewBackground(treeView *walk.TreeView) {
	bg, err := walk.NewSolidColorBrush(walk.RGB(218, 223, 230))
	if err != nil {
		common.Error("Error new color brush", err)
	} else {
		treeView.SetBackground(bg)
	}
}

func showCurrentItem(hostConfigText *walk.TextEdit) {
	index := conf.Config.CurrentHostIndex
	model := conf.Config.HostConfigModel
	if index < 0 || len(model.Roots) == 0 {

		return
	}
	current := model.RootAt(index)
	if current == nil {
		common.Error("Invalid CurrentHostIndex in config.json, cannot find the specific hosts")
	} else {
		if bytes, err := ioutil.ReadFile("conf/hosts/" + current.Text() + ".hosts"); err != nil {
			common.Error("Error reading host config: ", err)
		} else {
			hostConfigText.SetText(string(bytes))
		}
	}
}
