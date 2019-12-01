package tfd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcusolsson/tui-go"
)

// CreateSelectDialog create file dialog tui box
func CreateSelectDialog(extensions []string, hiddenDisplay bool) (path string, err error) {
	if len(extensions) == 0 {
		errorMsg := "target extensions is empty"
		return path, errors.New(errorMsg)
	}

	now, _ := os.Getwd()
	defer os.Chdir(now)

	t := tui.NewTheme()
	normal := tui.Style{Bg: tui.ColorWhite, Fg: tui.ColorBlack}
	t.SetStyle("normal", normal)
	t.SetStyle("list.item.selected", tui.Style{Bg: tui.ColorBlue, Fg: tui.ColorWhite})
	t.SetStyle("label.cd", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorGreen})

	title := "Please select "
	for _, extension := range extensions {
		title += "*." + extension + " "
	}
	titleLabel := tui.NewLabel(title)
	titleBox := tui.NewHBox(titleLabel)
	titleBox.SetSizePolicy(tui.Expanding, tui.Minimum)

	// 現在のパス
	cdBox := tui.NewHBox()
	setCDBox(cdBox)
	cdBox.SetBorder(true)
	cdBox.SetSizePolicy(tui.Expanding, tui.Minimum)

	// 兄弟ディレクトリの表示
	parentList := tui.NewList()
	makeParentList(parentList, hiddenDisplay)
	parentBox := tui.NewVBox(parentList)
	parentBox.SetTitle(" PARENT DIR ")
	parentBox.SetBorder(true)
	parentBox.SetSizePolicy(tui.Minimum, tui.Maximum)

	// 現在のディレクトリのファイル一覧の表示
	currentList := tui.NewList()
	makeCurrentList(currentList, extensions, hiddenDisplay)
	scrollArea := tui.NewScrollArea(currentList)
	currentBox := tui.NewVBox(scrollArea)
	currentBox.SetTitle(" CURRENT DIR ")
	currentBox.SetBorder(true)
	currentBox.SetSizePolicy(tui.Expanding, tui.Expanding)

	tree := tui.NewHBox(parentBox, currentBox)

	root := tui.NewVBox(titleBox, cdBox, tree)
	ui, err := tui.New(root)
	if err != nil {
		panic(err)
	}

	ui.SetTheme(t)
	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		path = ""
	})
	ui.SetKeybinding("Up", func() {
		currentList.OnKeyEvent(tui.KeyEvent{Key: tui.KeyUp})
		if currentBox.Size().Y >= 4 {
			if currentList.Selected() > currentBox.Size().Y-4 {
				if currentList.Selected() < currentList.Length()-2 {
					scrollArea.Scroll(0, -1)
				}
			} else {
				scrollArea.ScrollToTop()
			}
		}
	})
	ui.SetKeybinding("Down", func() {
		if currentList.Selected() < currentList.Length() {
			currentList.OnKeyEvent(tui.KeyEvent{Key: tui.KeyDown})

			if currentBox.Size().Y >= 4 {
				if currentList.Selected() < currentList.Length()-1 && currentList.Selected() > currentBox.Size().Y-4 {
					scrollArea.Scroll(0, 1)
				}
			}
		}
	})
	ui.SetKeybinding("Left", func() {
		os.Chdir("../")
		setCDBox(cdBox)
		makeCurrentList(currentList, extensions, hiddenDisplay)
		makeParentList(parentList, hiddenDisplay)
	})
	ui.SetKeybinding("Right", func() {
		if currentList.Length() > 0 {
			err := os.Chdir(currentList.SelectedItem())
			if err == nil {
				setCDBox(cdBox)
				makeCurrentList(currentList, extensions, hiddenDisplay)
				makeParentList(parentList, hiddenDisplay)
			}
		}
	})
	ui.SetKeybinding("Enter", func() {
		if currentList.Length() > 0 {
			cd, _ := os.Getwd()
			tmp := filepath.Join(cd, currentList.SelectedItem())
			_, err := ioutil.ReadDir(tmp)
			if err != nil {
				path = tmp
				ui.Quit()
			} else {
				err := os.Chdir(currentList.SelectedItem())
				if err == nil {
					setCDBox(cdBox)
					makeCurrentList(currentList, extensions, hiddenDisplay)
					makeParentList(parentList, hiddenDisplay)
				}
			}
		}
	})

	if err := ui.Run(); err != nil {
		panic(err)
	}
	if path == "" {
		errorMsg := "file is not selected"
		return path, errors.New(errorMsg)
	}
	return path, nil
}

func setCDBox(box *tui.Box) {
	box.Remove(0)
	cd, _ := os.Getwd()
	cdLabel := tui.NewLabel(" " + cd)
	cdLabel.SetStyleName("cd")
	box.Append(cdLabel)
}

func makeParentList(list *tui.List, hiddenDisplay bool) {
	currentPath, _ := os.Getwd()
	currentDir := filepath.Base(currentPath)
	parentDir := filepath.Join(currentDir, "../")

	list.RemoveItems()
	if currentDir != parentDir {
		dirs, _ := ioutil.ReadDir("../")
		for _, dir := range dirs {
			if dir.IsDir() && dir.Name() == currentDir {
				list.AddItems(" " + dir.Name() + "/")
				list.AddItems("                             ")
			}
		}
		for _, dir := range dirs {
			if dir.IsDir() && dir.Name() != currentDir {
				// 隠しファイルを表示するか
				if !hiddenDisplay && strings.HasPrefix(dir.Name(), ".") {
					continue
				}
				list.AddItems(" " + dir.Name() + "/")
			}
		}
		list.SetSelected(0)
	} else {
		list.AddItems("/")
		list.AddItems("                                        ")
		list.SetSelected(0)
	}

	// なんかこれやらんとレイアウト崩れる
	for i := 0; i < 50; i++ {
		list.AddItems("")
	}
}

func makeCurrentList(list *tui.List, extensions []string, hiddenDisplay bool) {
	list.RemoveItems()
	list.SetFocused(true)
	files, _ := ioutil.ReadDir("./")
	for _, file := range files {
		filename := file.Name()

		// 隠しファイルを表示するか
		if !hiddenDisplay && strings.HasPrefix(filename, ".") {
			continue
		}

		isDir := file.IsDir()
		if isDir {
			list.AddItems(filename + "/")
		} else {
			for _, extension := range extensions {
				if filepath.Ext(filename) == "."+extension {
					list.AddItems(filename)
				}
			}
		}
	}
	list.SetSelected(0)
}
