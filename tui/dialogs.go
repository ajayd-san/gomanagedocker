package tui

import (
	"github.com/ajayd-san/gomanagedocker/dockercmd"
	teadialog "github.com/ajayd-san/teaDialog"
)

const (
	dialogRemoveContainer teadialog.DialogType = iota
	dialogPruneContainers
	dialogRemoveImage
	dialogPruneImages
	dialogPruneVolumes
	dialogRemoveVolumes
	dialogImageScout
	dialogImageBuild
	dialogImageBuildProgress
)

func getImageScoutDialog(f func() (*dockercmd.ScoutData, error)) InfoCardWrapperModel {
	infoCard := teadialog.InitInfoCard(
		"Image Scout",
		"",
		dialogImageScout,
		teadialog.WithMinHeight(13),
		teadialog.WithMinWidth(130),
	)
	return InfoCardWrapperModel{
		tableChan: make(chan *TableModel),
		inner:     &infoCard,
		f:         f,
		spinner:   initialModel(),
	}
}

func getRemoveContainerDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeTogglePrompt("remVols", "Remove volumes?"),
		teadialog.MakeTogglePrompt("remLinks", "Remove links?"),
		teadialog.MakeTogglePrompt("force", "Force?"),
	}

	return teadialog.InitDialogWithPrompt("Remove Container Options:", prompts, dialogRemoveContainer, storage)
}

func getRemoveVolumeDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeTogglePrompt("force", "Force?"),
	}

	return teadialog.InitDialogWithPrompt("Remove Volume Options:", prompts, dialogRemoveVolumes, storage)
}

func getPruneContainersDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeOptionPrompt("confirm", "This will remove all stopped containers, are your sure?", []string{"Yes", "No"}),
	}

	return teadialog.InitDialogWithPrompt("Prune Containers: ", prompts, dialogPruneContainers, storage)
}

func getRemoveImageDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeTogglePrompt("force", "Force"),
		teadialog.MakeTogglePrompt("pruneChildren", "Prune Children"),
	}

	return teadialog.InitDialogWithPrompt("Remove Image Options:", prompts, dialogRemoveImage, storage)
}

func getPruneImagesDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeOptionPrompt("confirm", "This will remove all unused images, are your sure?", []string{"Yes", "No"}),
	}

	return teadialog.InitDialogWithPrompt("Prune Containers: ", prompts, dialogPruneImages, storage)
}

func getPruneVolumesDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeTogglePrompt("all", "Removed all unused volumes(not just anonymous ones)"),
		teadialog.MakeOptionPrompt("confirm", "This will remove all unused volumes, are your sure?", []string{"Yes", "No"}),
	}

	return teadialog.InitDialogWithPrompt("Prune Containers: ", prompts, dialogPruneVolumes, storage)
}

func getBuildImageDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		// teadialog.NewFilePicker("browser"),
		// NewFilePicker("filepicker"),
		teadialog.MakeTextInputPrompt("image_tags", "Image Tags:"),
	}

	return teadialog.InitDialogWithPrompt("Build Image: ", prompts, dialogImageBuild, storage)
}

func getBuildProgress(loading loadingModel) buildProgressModel {

	infoCard := teadialog.InitInfoCard(
		"Image Scout",
		"",
		dialogImageBuildProgress,
		teadialog.WithMinHeight(8),
		teadialog.WithMinWidth(100),
	)

	return buildProgressModel{
		loading: loading,
		inner:   &infoCard,
	}
}
