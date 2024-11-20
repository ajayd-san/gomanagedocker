package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/knadh/koanf/v2"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		UserConfig string
		Want       map[string]any
	}{
// 		{
// 			UserConfig: "",
// 			Want: map[string]any{
// 				"config.Polling-Time":             500,
// 				"config.Tab-Order":                []any{"images", "containers", "volumes"},
// 				"config.Notification-Timeout":     2000,
// 				"keybindings.navigation.Enter":    			   []any{"enter"},
// 				"keybindings.navigation.Back":     			   []any{"esc"},
// 				"keybindings.navigation.Quit":     			   []any{"q", "ctrl+c"},
// 				"keybindings.navigation.Select":   			   []any{" "},
// 				"keybindings.navigation.NextTab":  			   []any{"right", "tab", "l"},
// 				"keybindings.navigation.PrevTab":  			   []any{"left", "shift+tab", "h"},
// 				"keybindings.navigation.NextItem": 			   []any{"down", "j"},
// 				"keybindings.navigation.PrevItem": 			   []any{"up", "k"},
// 				"keybindings.navigation.PrevPage": 			   []any{"["},
// 				"keybindings.navigation.NextPage": 			   []any{"]"},
//     			"keybindings.image.Build":                     []any{"b"},
//     			"keybindings.image.CopyId":                    []any{"c"},
//     			"keybindings.image.Delete":                    []any{"d"},
//     			"keybindings.image.DeleteForce":               []any{"D"},
//     			"keybindings.image.DeleteForceBulk":           []any{"D"},
//     			"keybindings.image.ExitSelectionMode":         []any{"esc"},
//     			"keybindings.image.Prune":                     []any{"p"},
//     			"keybindings.image.Rename":                    []any{"R"},
//     			"keybindings.image.Run":                       []any{"r"},
//     			"keybindings.image.RunAndExec":                []any{"x"},
//     			"keybindings.image.Scout":                     []any{"s"},
// 				"keybindings.container.CopyId":                []any{"c"},
//     			"keybindings.container.Delete":                []any{"d"},
//     			"keybindings.container.DeleteForce":           []any{"D"},
//     			"keybindings.container.Exec":                  []any{"x"},
//     			"keybindings.container.Prune":                 []any{"p"},
//     			"keybindings.container.Restart":               []any{"r"},
//     			"keybindings.container.ShowLogs":              []any{"L"},
//     			"keybindings.container.ToggleListAll":         []any{"a"},
//     			"keybindings.container.TogglePause":           []any{"t"},
//     			"keybindings.container.ToggleStartStop":       []any{"s"},
//     			"keybindings.containerBulk.DeleteForce":       []any{"D"},
//     			"keybindings.containerBulk.ExitSelectionMode": []any{"esc"},
//     			"keybindings.containerBulk.Restart":           []any{"r"},
//     			"keybindings.containerBulk.ToggleListAll":     []any{"a"},
//     			"keybindings.containerBulk.TogglePause":       []any{"t"},
//     			"keybindings.containerBulk.ToggleStartStop":   []any{"s"},
//     			"keybindings.volume.Delete":   				   []any{"d"},
//     			"keybindings.volume.DeleteForce":     		   []any{"D"},
//     			"keybindings.volume.Prune":   				   []any{"p"},
//     			"keybindings.volume.CopyId":   				   []any{"c"},
//     			"keybindings.volumeBulk.DeleteForce":   	   []any{"D"},
//     			"keybindings.volumeBulk.ExitSelectionMode":    []any{"esc"},
// 			},
// 		},
// 		{
// 			UserConfig: `config:
//   Polling-Time: 100`,
// 			Want: map[string]any{
// 				"config.Polling-Time": 100,
// 				"config.Tab-Order": []any{"images", "containers", "volumes"},
// 				"config.Notification-Timeout": 2000,
// 				"keybindings.navigation.Enter":    			   []any{"enter"},
// 				"keybindings.navigation.Back":     			   []any{"esc"},
// 				"keybindings.navigation.Quit":     			   []any{"q", "ctrl+c"},
// 				"keybindings.navigation.Select":   			   []any{" "},
// 				"keybindings.navigation.NextTab":  			   []any{"right", "tab", "l"},
// 				"keybindings.navigation.PrevTab":  			   []any{"left", "shift+tab", "h"},
// 				"keybindings.navigation.NextItem": 			   []any{"down", "j"},
// 				"keybindings.navigation.PrevItem": 			   []any{"up", "k"},
// 				"keybindings.navigation.PrevPage": 			   []any{"["},
// 				"keybindings.navigation.NextPage": 			   []any{"]"},
//     			"keybindings.image.Build":                     []any{"b"},
//     			"keybindings.image.CopyId":                    []any{"c"},
//     			"keybindings.image.Delete":                    []any{"d"},
//     			"keybindings.image.DeleteForce":               []any{"D"},
//     			"keybindings.image.DeleteForceBulk":           []any{"D"},
//     			"keybindings.image.ExitSelectionMode":         []any{"esc"},
//     			"keybindings.image.Prune":                     []any{"p"},
//     			"keybindings.image.Rename":                    []any{"R"},
//     			"keybindings.image.Run":                       []any{"r"},
//     			"keybindings.image.RunAndExec":                []any{"x"},
//     			"keybindings.image.Scout":                     []any{"s"},
// 				"keybindings.container.CopyId":                []any{"c"},
//     			"keybindings.container.Delete":                []any{"d"},
//     			"keybindings.container.DeleteForce":           []any{"D"},
//     			"keybindings.container.Exec":                  []any{"x"},
//     			"keybindings.container.Prune":                 []any{"p"},
//     			"keybindings.container.Restart":               []any{"r"},
//     			"keybindings.container.ShowLogs":              []any{"L"},
//     			"keybindings.container.ToggleListAll":         []any{"a"},
//     			"keybindings.container.TogglePause":           []any{"t"},
//     			"keybindings.container.ToggleStartStop":       []any{"s"},
//     			"keybindings.containerBulk.DeleteForce":       []any{"D"},
//     			"keybindings.containerBulk.ExitSelectionMode": []any{"esc"},
//     			"keybindings.containerBulk.Restart":           []any{"r"},
//     			"keybindings.containerBulk.ToggleListAll":     []any{"a"},
//     			"keybindings.containerBulk.TogglePause":       []any{"t"},
//     			"keybindings.containerBulk.ToggleStartStop":   []any{"s"},
//     			"keybindings.volume.Delete":   				   []any{"d"},
//     			"keybindings.volume.DeleteForce":     		   []any{"D"},
//     			"keybindings.volume.Prune":   				   []any{"p"},
//     			"keybindings.volume.CopyId":   				   []any{"c"},
//     			"keybindings.volumeBulk.DeleteForce":   	   []any{"D"},
//     			"keybindings.volumeBulk.ExitSelectionMode":    []any{"esc"},
// 			},
// 		},
// 		{
// 			UserConfig: `config:
//   Polling-Time: 200
//   Tab-Order: [containers, volumes]
//   Notification-Timeout: 10000`,
// 			Want: map[string]any{
// 				"config.Polling-Time": 200,
// 				"config.Tab-Order": []any{"containers", "volumes"},
// 				"config.Notification-Timeout": 10000,
// 				"keybindings.navigation.Enter":    			   []any{"enter"},
// 				"keybindings.navigation.Back":     			   []any{"esc"},
// 				"keybindings.navigation.Quit":     			   []any{"q", "ctrl+c"},
// 				"keybindings.navigation.Select":   			   []any{" "},
// 				"keybindings.navigation.NextTab":  			   []any{"right", "tab", "l"},
// 				"keybindings.navigation.PrevTab":  			   []any{"left", "shift+tab", "h"},
// 				"keybindings.navigation.NextItem": 			   []any{"down", "j"},
// 				"keybindings.navigation.PrevItem": 			   []any{"up", "k"},
// 				"keybindings.navigation.PrevPage": 			   []any{"["},
// 				"keybindings.navigation.NextPage": 			   []any{"]"},
//     			"keybindings.image.Build":                     []any{"b"},
//     			"keybindings.image.CopyId":                    []any{"c"},
//     			"keybindings.image.Delete":                    []any{"d"},
//     			"keybindings.image.DeleteForce":               []any{"D"},
//     			"keybindings.image.DeleteForceBulk":           []any{"D"},
//     			"keybindings.image.ExitSelectionMode":         []any{"esc"},
//     			"keybindings.image.Prune":                     []any{"p"},
//     			"keybindings.image.Rename":                    []any{"R"},
//     			"keybindings.image.Run":                       []any{"r"},
//     			"keybindings.image.RunAndExec":                []any{"x"},
//     			"keybindings.image.Scout":                     []any{"s"},
// 				"keybindings.container.CopyId":                []any{"c"},
//     			"keybindings.container.Delete":                []any{"d"},
//     			"keybindings.container.DeleteForce":           []any{"D"},
//     			"keybindings.container.Exec":                  []any{"x"},
//     			"keybindings.container.Prune":                 []any{"p"},
//     			"keybindings.container.Restart":               []any{"r"},
//     			"keybindings.container.ShowLogs":              []any{"L"},
//     			"keybindings.container.ToggleListAll":         []any{"a"},
//     			"keybindings.container.TogglePause":           []any{"t"},
//     			"keybindings.container.ToggleStartStop":       []any{"s"},
//     			"keybindings.containerBulk.DeleteForce":       []any{"D"},
//     			"keybindings.containerBulk.ExitSelectionMode": []any{"esc"},
//     			"keybindings.containerBulk.Restart":           []any{"r"},
//     			"keybindings.containerBulk.ToggleListAll":     []any{"a"},
//     			"keybindings.containerBulk.TogglePause":       []any{"t"},
//     			"keybindings.containerBulk.ToggleStartStop":   []any{"s"},
//     			"keybindings.volume.Delete":   				   []any{"d"},
//     			"keybindings.volume.DeleteForce":     		   []any{"D"},
//     			"keybindings.volume.Prune":   				   []any{"p"},
//     			"keybindings.volume.CopyId":   				   []any{"c"},
//     			"keybindings.volumeBulk.DeleteForce":   	   []any{"D"},
//     			"keybindings.volumeBulk.ExitSelectionMode":    []any{"esc"},
// 			},
// 		},
		{
			UserConfig: `config:
  Polling-Time: 200
  Tab-Order: [containers, volumes]
  Notification-Timeout: 10000
keybindings:
  navigation:
    Enter: [e]
    Quit: [q, ctrl+a]
  image:
    Build: [B]
    CopyId: [C]
  container:
    CopyId: [C]
    Delete: [D]
  `,
			Want: map[string]any{
				"config.Polling-Time": 200,
				"config.Tab-Order": []any{"containers", "volumes"},
				"config.Notification-Timeout": 10000,
				"keybindings.navigation.Enter":    			   []any{"e"},
				"keybindings.navigation.Back":     			   []any{"esc"},
				"keybindings.navigation.Quit":     			   []any{"q", "ctrl+a"},
				"keybindings.navigation.Select":   			   []any{" "},
				"keybindings.navigation.NextTab":  			   []any{"right", "tab", "l"},
				"keybindings.navigation.PrevTab":  			   []any{"left", "shift+tab", "h"},
				"keybindings.navigation.NextItem": 			   []any{"down", "j"},
				"keybindings.navigation.PrevItem": 			   []any{"up", "k"},
				"keybindings.navigation.PrevPage": 			   []any{"["},
				"keybindings.navigation.NextPage": 			   []any{"]"},
    			"keybindings.image.Build":                     []any{"B"},
    			"keybindings.image.CopyId":                    []any{"C"},
    			"keybindings.image.Delete":                    []any{"d"},
    			"keybindings.image.DeleteForce":               []any{"D"},
    			"keybindings.image.DeleteForceBulk":           []any{"D"},
    			"keybindings.image.ExitSelectionMode":         []any{"esc"},
    			"keybindings.image.Prune":                     []any{"p"},
    			"keybindings.image.Rename":                    []any{"R"},
    			"keybindings.image.Run":                       []any{"r"},
    			"keybindings.image.RunAndExec":                []any{"x"},
    			"keybindings.image.Scout":                     []any{"s"},
				"keybindings.container.CopyId":                []any{"C"},
    			"keybindings.container.Delete":                []any{"D"},
    			"keybindings.container.DeleteForce":           []any{"D"},
    			"keybindings.container.Exec":                  []any{"x"},
    			"keybindings.container.Prune":                 []any{"p"},
    			"keybindings.container.Restart":               []any{"r"},
    			"keybindings.container.ShowLogs":              []any{"L"},
    			"keybindings.container.ToggleListAll":         []any{"a"},
    			"keybindings.container.TogglePause":           []any{"t"},
    			"keybindings.container.ToggleStartStop":       []any{"s"},
    			"keybindings.containerBulk.DeleteForce":       []any{"D"},
    			"keybindings.containerBulk.ExitSelectionMode": []any{"esc"},
    			"keybindings.containerBulk.Restart":           []any{"r"},
    			"keybindings.containerBulk.ToggleListAll":     []any{"a"},
    			"keybindings.containerBulk.TogglePause":       []any{"t"},
    			"keybindings.containerBulk.ToggleStartStop":   []any{"s"},
    			"keybindings.volume.Delete":   				   []any{"d"},
    			"keybindings.volume.DeleteForce":     		   []any{"D"},
    			"keybindings.volume.Prune":   				   []any{"p"},
    			"keybindings.volume.CopyId":   				   []any{"c"},
    			"keybindings.volumeBulk.DeleteForce":   	   []any{"D"},
    			"keybindings.volumeBulk.ExitSelectionMode":    []any{"esc"},
			},
		},
	}

	for id, test := range tests {
		tempFile, _ := os.CreateTemp("", "")
		tempFile.WriteString(test.UserConfig)
		defer os.Remove(tempFile.Name())

		got := koanf.New(".")
		filePath := tempFile.Name()
		ReadConfig(got, filePath)

		if !cmp.Equal(got.All(), test.Want) {
			t.Errorf("Fail %d: %s", id, cmp.Diff(got.All(), test.Want))
		}

	}

}
