package tui

import (
	teadialog "github.com/ajayd-san/teaDialog"
)

const (
	dialogRemoveContainer teadialog.DialogType = iota
	dialogPruneContainers
	dialogRemoveImage
	dialogPruneImages
	dialogPruneVolumes
	dialogRemoveVolumes
)

func getRemoveContainerDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeTogglePrompt("remVols", "Remove volumes?"),
		teadialog.MakeTogglePrompt("remLinks", "Remove links?"),
		teadialog.MakeTogglePrompt("force", "Force?"),
	}

	return teadialog.InitDialogue("Remove Container Options:", prompts, dialogRemoveContainer, storage)
}

func getRemoveVolumeDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeTogglePrompt("force", "Force?"),
	}

	return teadialog.InitDialogue("Remove Volume Options:", prompts, dialogRemoveVolumes, storage)
}

func getPruneContainersDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeOptionPrompt("confirm", "This will remove all stopped containers, are your sure?", []string{"Yes", "No"}),
	}

	return teadialog.InitDialogue("Prune Containers: ", prompts, dialogPruneContainers, storage)
}

func getRemoveImageDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeTogglePrompt("force", "Force"),
		teadialog.MakeTogglePrompt("pruneChildren", "Prune Children"),
	}

	return teadialog.InitDialogue("Remove Image Options:", prompts, dialogRemoveImage, storage)
}

func getPruneImagesDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeOptionPrompt("confirm", "This will remove all unused images, are your sure?", []string{"Yes", "No"}),
	}

	return teadialog.InitDialogue("Prune Containers: ", prompts, dialogPruneImages, storage)
}

func getPruneVolumesDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeTogglePrompt("all", "Removed all unused volumes"),
		teadialog.MakeOptionPrompt("confirm", "This will remove all unused volumes, are your sure?", []string{"Yes", "No"}),
	}

	return teadialog.InitDialogue("Prune Containers: ", prompts, dialogPruneVolumes, storage)
}
