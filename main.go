package main

import (
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget    *widget.Entry    // window pane where we would right text
	PreviewWidget *widget.RichText // window pane where the formatted markdown text will be visible. md => mark down
	CurrentFile   fyne.URI         // To know what current file you have opened. As the application can save contents into a file.
	SaveMenuItem  *fyne.MenuItem   // To save menu.
}

var cfg config

func main() {
	// create a fyne app
	a := app.New()

	// create a window for the fyne app
	win := a.NewWindow("Markdown")

	// get the user interface for the fyne app
	edit, preview := cfg.makeUI()
	// create menu items for the window
	cfg.createMenuItems(win)

	// set the content of the window
	win.SetContent(container.NewHSplit(edit, preview))

	// show window and run the app
	win.Resize(fyne.Size{Width: 800, Height: 500}) // size in pixels.
	win.CenterOnScreen()
	win.ShowAndRun()
}

// makeUI creates two widgets, assignsthem to the app config, and
// adds a listener to the edit widget that updates the preview widget
// with parsed markdown whenever the user types something
func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")
	app.EditWidget = edit
	app.PreviewWidget = preview

	// add a event listener
	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (app *config) createMenuItems(win fyne.Window) {

	// create three menu items
	openMenuItems := fyne.NewMenuItem("Open...", app.openFunc(win))
	saveMenuItems := fyne.NewMenuItem("Save", app.saveFunc(win))
	app.SaveMenuItem = saveMenuItems
	app.SaveMenuItem.Disabled = true
	saveAsMenuItems := fyne.NewMenuItem("Save as...", app.saveAsFunc(win))

	// create a menu(file menu) and add three menu items to it.
	fileMenu := fyne.NewMenu("File", openMenuItems, saveMenuItems, saveAsMenuItems)

	// create the main menu and add the avaialable menu/menus (file menu) to it.
	menu := fyne.NewMainMenu(fileMenu)

	// set the main menu for the application
	win.SetMainMenu(menu)
}

// create a filter for ".md" files, So that we can only open/write files with .md extensions.
// We're going to use this filter with openFunc and saveAsFunc
var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *config) saveFunc(win fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			write, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			write.Write([]byte(app.EditWidget.Text))
			defer write.Close()
		}
	}
}

func (app *config) openFunc(win fyne.Window) func() {
	return func() {
		// create a dialog box of to open file of type .md
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				// if there is an error. It means we couldn't read the file/filesystem. Show the error and return.
				dialog.ShowError(err, win)
				return
			}

			if read == nil {
				// user cancelled (clicked on cancel button). User chose not to open the file. So return
				return
			}

			defer read.Close()

			// read the file.
			data, err := ioutil.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// populate the edit widget with the data read from the file.
			app.EditWidget.SetText(string(data))

			// keep record of the current file.
			app.CurrentFile = read.URI()

			// update/reset the title of the window to "current window title - The new file name".
			win.SetTitle(win.Title() + " - " + read.URI().Name())

			// make sure save menu item is enabled
			app.SaveMenuItem.Disabled = false

		}, win)

		// set the filters so that we are able to open only .md files.
		openDialog.SetFilter(filter)
		// now show the dialog
		openDialog.Show()

	}
}

func (app *config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if write == nil {
				// user cancelled. They chose not to save the file. So return.
				return
			}

			// check if the file has is .md file or not. If it's not .md then show info to the user
			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Please name your file with a .md extension!", win)
				return
			}

			// save the file
			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()

			defer write.Close()

			// reset the title of the window to the "current window title - The new file name".
			win.SetTitle(win.Title() + " - " + write.URI().Name())

			// enable saveAS menu item
			app.SaveMenuItem.Disabled = false
		}, win)
		// giving default name to the file
		saveDialog.SetFileName("Untitled.md")
		// setting filter so that we can write to/save to .md files only.
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}
