package main

import (
	"embed"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Settings struct {
	TargetFolder  string   `json:"targetFolder"`
	FilesToDelete []string `json:"filesToDelete"`
}

type SettingsNew struct {
	TargetFolder  string               `json:"targetFolder"`
	FilesToDelete []SettingsFilesGroup `json:"filesToDelete"`
}

type SettingsFilesGroup struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

//go:embed *.json
var settingFile embed.FS
var settings Settings
var settingsNew SettingsNew
var filesToDelete []string

func (s *Settings) read() string {
	var logText string

	settingsData, err := settingFile.ReadFile("settings.json")
	if err != nil {
		logText += "❌ Error reading the settings file => " + err.Error() + "\n"
	} else {
		logText += "✅ Reading settings.json successfull\n"

		if err := json.Unmarshal(settingsData, s); err != nil {
			logText += "❌ Error parsing settings from settings.json => " + err.Error() + "\n"
		} else {
			logText += "🔄 The settings have been parsed successfully\n"
		}
	}

	return logText
}

func (s *SettingsNew) read() string {
	var logText string

	settingsData, err := settingFile.ReadFile("settings.new.json")
	if err != nil {
		logText += "❌ Error reading the settings file => " + err.Error() + "\n"
	} else {
		logText += "✅ Reading NEW settings.new.json successfull\n"

		if err := json.Unmarshal(settingsData, s); err != nil {
			logText += "❌ Error parsing settings from settings.json => " + err.Error() + "\n"
		} else {
			logText += "🔄 The NEW settings have been parsed successfully\n"
		}
	}

	return logText
}

func createLogArea() (*widget.Label, *container.Scroll) {
	logsOutput := widget.NewLabel("Welcome to XIV Translation Cleaner\n\n")
	logsOutput.Wrapping = fyne.TextWrapBreak
	logContainer := container.NewScroll(logsOutput)

	return logsOutput, logContainer
}

func createFolderSelectionEntry() *widget.Entry {
	folderEntry := widget.NewEntry()
	folderEntry.SetPlaceHolder("Enter the path")
	folderEntry.SetText(settings.TargetFolder)

	return folderEntry
}

func createCheckboxes(logsOutput *widget.Label) *fyne.Container {
	checkboxesLabel := widget.NewLabel("Select the translation that you want to delete:")
	checkboxesLabel.TextStyle.Bold = true
	checkboxesLayer := container.NewVBox(
		checkboxesLabel,
	)

	for _, data := range settingsNew.FilesToDelete {
		cbData := &SettingsFilesGroup{Name: data.Name, Files: data.Files}
		cb := widget.NewCheck(cbData.Name, func(checked bool) {
			if checked {
				logsOutput.SetText(logsOutput.Text + "Checked Files: " + strings.Join(cbData.Files, ", ") + "\n")
				filesToDelete = append(filesToDelete, cbData.Files...)
			}
		})
		checkboxesLayer.Add(cb)
	}

	return checkboxesLayer
}

// func getOnlySelected() {

// }

func createDeleteButton(logsOutput *widget.Label, logContainer *container.Scroll, folderEntry *widget.Entry, myWindow fyne.Window) *widget.Button {
	deleteButton := widget.NewButton("Delete selected translation", func() {
		targetFolder := folderEntry.Text

		// Check if the target folder exists
		if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
			logsOutput.SetText(logsOutput.Text + "Target folder does not exist: " + targetFolder + "\n")
			logContainer.ScrollToBottom()
			return
		}

		// Confirm with the user before proceeding
		confirm := dialog.NewConfirm("Confirmation", "Do you want to delete selected translation?", func(accept bool) {
			if accept {
				deleteSelectedItems(logsOutput, logContainer, targetFolder)
			}
		}, myWindow)

		confirm.SetDismissText("Cancel")
		confirm.Show()
	})

	return deleteButton
}

func deleteSelectedItems(logsOutput *widget.Label, logContainer *container.Scroll, targetFolder string) {
	// Loop through the list of files to delete and delete them
	for _, fileToDelete := range settings.FilesToDelete {
		filePath := filepath.Join(targetFolder, "exd", fileToDelete)

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			logsOutput.SetText(logsOutput.Text + "❌ " + err.Error() + "\n")
			logContainer.ScrollToBottom()
			continue
		}

		if fileInfo.IsDir() {
			err = os.RemoveAll(filePath)
		} else {
			err = os.Remove(filePath)
		}

		if err != nil {
			logsOutput.SetText(logsOutput.Text + "❌ " + err.Error() + "\n")
			logContainer.ScrollToBottom()
		} else {
			logsOutput.SetText(logsOutput.Text + "✅ Deleted: " + fileToDelete + "\n")
			logContainer.Refresh()
		}
	}
	logsOutput.SetText(logsOutput.Text + "\n---- Selected items were deleted ----\n\n")
	logContainer.Refresh()
}

func main() {
	myApp := app.New()
	// myApp.Settings().SetTheme(theme.DarkTheme())
	myWindow := myApp.NewWindow("XIV Translation Cleaner")

	logsOutput, logContainer := createLogArea()

	// Init Settings
	logsOutput.SetText(logsOutput.Text + settings.read())
	logsOutput.SetText(logsOutput.Text + settingsNew.read())
	logContainer.ScrollToBottom()

	folderEntry := createFolderSelectionEntry()

	checkboxes := createCheckboxes(logsOutput)

	deleteButton := createDeleteButton(logsOutput, logContainer, folderEntry, myWindow)

	// Layout
	topRow := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Path to the trasnlation folder: "),
		folderEntry,
	)
	middleRow := container.NewVBox(
		checkboxes,
		container.New(
			layout.NewCenterLayout(),
			deleteButton,
		),
	)

	// Create main container
	content := container.NewVSplit(
		container.NewVBox(
			topRow,
			middleRow,
		),
		logContainer,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(900, 600))
	myWindow.ShowAndRun()
}
