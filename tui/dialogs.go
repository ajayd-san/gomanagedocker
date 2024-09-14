package tui

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/ajayd-san/gomanagedocker/service/dockercmd"
	"github.com/ajayd-san/gomanagedocker/tui/components"
	teadialog "github.com/ajayd-san/teaDialog"
)

const (
	// containers
	dialogRemoveContainer teadialog.DialogType = iota
	dialogPruneContainers

	// images
	dialogRemoveImage
	dialogPruneImages
	dialogRunImage
	dialogImageScout
	dialogImageBuild
	dialogImageBuildProgress

	// volumes
	dialogPruneVolumes
	dialogRemoveVolumes

	// pods
	dialogPrunePods
	dialogDeletePod
)

func getRunImageDialog(storage map[string]string) teadialog.Dialog {
	// MGS nerds will prolly like this
	prompt := []teadialog.Prompt{
		teadialog.MakeTextInputPrompt(
			"port",
			"Port mappings",
			teadialog.WithPlaceHolder("Ex: 1011:2016,226:1984/udp"),
			teadialog.WithTextWidth(30),
		),
		teadialog.MakeTextInputPrompt(
			"name",
			"Name",
			teadialog.WithPlaceHolder("prologueAwakening"),
			teadialog.WithTextWidth(30),
		),
		teadialog.MakeTextInputPrompt(
			"env",
			"Environment variables",
			teadialog.WithPlaceHolder("VENOM=AHAB,DD=goodDoggo"),
			teadialog.WithTextWidth(30),
		),
	}

	title := "Run Image\n(Leave inputs blank for defaults)"
	return teadialog.InitDialogWithPrompt(title, prompt, dialogRunImage, storage, teadialog.WithShowFullHelp(true))
}

func getImageScoutDialog(f func() (*dockercmd.ScoutData, error)) DockerScoutInfoCard {
	infoCard := teadialog.InitInfoCard(
		"Image Scout",
		"",
		dialogImageScout,
		teadialog.WithMinHeight(13),
		teadialog.WithMinWidth(130),
	)
	return DockerScoutInfoCard{
		tableChan: make(chan *TableModel),
		inner:     &infoCard,
		f:         f,
		spinner:   components.InitialModel(),
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

// Gets the build progress bar info card/dialog
func getBuildProgress(progressBar components.ProgressBar) buildProgressModel {

	infoCard := teadialog.InitInfoCard(
		"Image Build",
		"",
		dialogImageBuildProgress,
		teadialog.WithMinHeight(8),
		teadialog.WithMinWidth(100),
	)

	reg := regexp.MustCompile(`(?i)Step\s(\d+)\/(\d+)\s?:\s(.*)`)

	return buildProgressModel{
		progressChan: make(chan string, 10),
		regex:        reg,
		progressBar:  progressBar,
		inner:        &infoCard,
	}
}

// PODS
func getPrunePodsDialog(storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeOptionPrompt("confirmPrunePods", "This will remove all stopped pods, are your sure?", []string{"Yes", "No"}),
	}

	return teadialog.InitDialogWithPrompt("Prune Pods: ", prompts, dialogPrunePods, storage)
}

func getRemovePodDialog(running int, storage map[string]string) teadialog.Dialog {
	prompts := []teadialog.Prompt{
		teadialog.MakeTogglePrompt("force", "Force?"),
	}

	if running > 0 {
		runningContainersString := containerCountForeground.Render(fmt.Sprintf("%d running", running))
		confirmPrompt := teadialog.MakeOptionPrompt(
			"confirm",
			fmt.Sprintf(
				"Are you sure? This pod has %s containers.",
				runningContainersString,
			),
			[]string{"Yes", "No"})
		prompts = slices.Insert(prompts, 0, confirmPrompt)
	}
	return teadialog.InitDialogWithPrompt("Remove Pod Options:", prompts, dialogDeletePod, storage)
}
