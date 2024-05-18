package tui

import dialog "github.com/ajayd-san/teaDialog"

const (
	dialogRemoveContainer dialog.DialogType = iota
	dialogRemoveImage
)

func getRemoveContainerDialog(storage map[string]string) dialog.Dialog {
	prompts := []dialog.Prompt{
		dialog.MakeTogglePrompt("remVols", "Remove volumes?"),
		dialog.MakeTogglePrompt("remLinks", "Remove links?"),
		dialog.MakeTogglePrompt("force", "Force?"),
	}

	return dialog.InitDialogue("Remove Container Options:", prompts, dialogRemoveContainer, storage)
}

func getRemoveImageDialog(storage map[string]string) dialog.Dialog {
	prompts := []dialog.Prompt{
		dialog.MakeTogglePrompt("force", "Force"),
		dialog.MakeTogglePrompt("pruneChildren", "Prune Children"),
	}

	return dialog.InitDialogue("Remove Image Options:", prompts, dialogRemoveImage, storage)
}
