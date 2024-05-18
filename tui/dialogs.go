package tui

import dialog "github.com/ajayd-san/teaDialog"

const (
	dialogRemoveContainer dialog.DialogType = iota
)

func getRemoveContainerDialog(storage map[string]string) dialog.Dialog {
	prompts := []dialog.Prompt{
		dialog.MakeTogglePrompt("remVols", "Remove volumes?"),
		dialog.MakeTogglePrompt("remLinks", "Remove links?"),
		dialog.MakeTogglePrompt("force", "Force?"),
	}

	return dialog.InitDialogue("Remove Container Options:", prompts, dialogRemoveContainer, storage)
}
