# goManageDocker

Do Docker commands slip your mind because you don't use Docker often enough? Sick of googling commands for everyday tasks? GoManageDocker is designed to NUKE this annoyance. 

Introducing **goManageDocker** (get it?)! This blazing fast TUI, made using Go and BubbleTea, will make managing your Docker objects a breeze. 

## Contents

1. [Install Instructions](#install-instructions)
2. [Features](#features)
3. [Keybinds](#keybinds)
4. [Configuration](#configuration)
5. [Roadmap](#roadmap)
6. [Found an issue?](#found-an-issue-)
7. [Contributing](#contributing)

## Install Instructions

### Unix

You can install the latest release of goManageDocker on UNIX systems with a simple bash script:

```
bash -c "$(curl -sLo- https://raw.githubusercontent.com/ajayd-san/gomanagedocker/main/install.sh)"
```

Start the program with `gmd`. 

### Build from source

Just build like any other Go binary: 

```
go install github.com/ajayd-san/gomanagedocker@main
```

Start the program with `gomanagedocker` (Rename it to `gmd` if you'd like, the binary will be installed at your `$GOPATH`).

### Windows

You can get the latest precompiled binary from releases or you may build from source. 

Now, **goManageDocker 😏!!**

> [!NOTE]
> goManageDocker runs best on terminals that support ANSI 256 colors and designed to run while the **terminal is maximized**.

### Docker
Want to try this without installing a binary? I gotchu!

```
docker run -it -v /var/run/docker.sock:/var/run/docker.sock kakshipth/gomanagedocker:latest
```

Alias it to something quicker (unless you like typing a lot 🙄)
## Features

### **New in v1.4:**\

1. Global notification system
	![notificationSystem](vhs/gifs/notifications.gif)
	
2. Bulk operation mode: select multiple objects before performing an operations (saves so much time!!)
	![bulkDelete](vhs/gifs/bulkDelete.gif)
	
3. Build image from Dockerfile using `b`
	![build](vhs/gifs/build.gif)
	
4. View live logs from a container using `L`
 	![runImage](vhs/gifs/logs.gif)

5. Run image now takes arguments for port, name and env vars. 
 	![runImage](vhs/gifs/runImage.gif)

### **Previous release features:**


1. Easy navigation with vim keybinds and arrow keys.
   ![intro](vhs/gifs/intro.gif)

2. Exec into selected container with A SINGLE KEYSTROKE: `x`...How cool is that?
   ![exec](vhs/gifs/exec.gif)

3. Delete objects using `d` (You can force delete with `D`, you won't have to answer a prompt this way)
   ![delete](vhs/gifs/delete.gif)

4. Prune objects using `p`
   ![prune](vhs/gifs/prune.gif)

5. start/stop/pause/restart containers with `s`, `t` and `r`
   ![startstop](vhs/gifs/startstop.gif)

6. Filter objects with `/`
   ![search](vhs/gifs/search.gif)
   
7. Perfrom docker scout with `s`
   ![scout](vhs/gifs/scout.gif)
   
8. Run an image directly from the image tab by pressing `r`.
   ![runImage](vhs/gifs/runImage.gif) 
   
9. You can directly copy the ID to your clipboard of an object by pressing `c`.
   ![copyId](vhs/gifs/copyId.gif)
10. You can now run and exec into an image directly from the images tab with  `x`
    ![runAndExec](vhs/gifs/execFromImgs.gif)

## Keybinds

### Navigation
| Operation 		   | Key                                                                 |
|------------------|---------------------------------------------------------------------|
| Back 			   | <kbd>Esc</kbd>                                                      |
| Quit      		   | <kbd>Ctrl</kbd> + <kbd>c</kbd> / <kbd>q</kbd>                       |
| Next Tab  		   | <kbd>→</kbd> / <kbd>l</kbd> / <kbd>Tab</kbd>                        |
| Prev Tab  		   | <kbd>←</kbd> / <kbd>h</kbd> / <kbd>Shift</kbd> + <kbd>Tab</kbd>     |
| Next Item 		   | <kbd>↓</kbd> / <kbd>j</kbd>                                         |
| Prev Item 		   | <kbd>↑</kbd> / <kbd>k</kbd>                                         |
| Next Page 		   | <kbd>[</kbd>                                                        |
| Prev Page 		   | <kbd>]</kbd>                                                        |
| Enter bulk mode  | <kbd>Space</kbd>                                                    |

### Image
| Operation         | Key                                                           |
|-------------------|---------------------------------------------------------------|
| Run               | <kbd>r</kbd>                                                  |
| Build Image       | <kbd>b</kbd>                                                  |
| Scout             | <kbd>s</kbd>                                                  |
| Prune             | <kbd>p</kbd>                                                  |
| Delete            | <kbd>d</kbd>                                                  |
| Delete (Force)    | <kbd>D</kbd>                                                  |
| Copy ID           | <kbd>c</kbd>                                                  |
| Run and Exec      | <kbd>x</kbd>                                                  |

### Container
| Operation         | Key                                                           |
|-------------------|---------------------------------------------------------------|
| Toggle List All   | <kbd>a</kbd>                                                  |
| Toggle Start/Stop | <kbd>s</kbd>                                                  |
| Toggle Pause      | <kbd>t</kbd>                                                  |
| Restart           | <kbd>r</kbd>                                                  |
| Delete            | <kbd>d</kbd>                                                  |
| Delete (Force)    | <kbd>D</kbd>                                                  |
| Exec              | <kbd>x</kbd>                                                  |
| Prune             | <kbd>p</kbd>                                                  |
| Copy ID           | <kbd>c</kbd>                                                  |
| Show Logs         | <kbd>L</kbd>                                                  |

### Volume
| Operation         | Key                                                           |
|-------------------|---------------------------------------------------------------|
| Delete            | <kbd>d</kbd>                                                  |
| Prune             | <kbd>p</kbd>                                                  |
| Copy Volume Name  | <kbd>c</kbd>                                                  |

## Configuration

I've added support for config files from V1.2.

Place `gomanagedocker/gomanagedocker.yaml` in your XDG config folder and configure to your heart's content!

Default Configuration:  

```
config:
  Polling-Time: 500
  Tab-Order: [images, containers, volumes]
  Notification-timeout: 2000
```

- Polling-Time: Set how frequently the program calls the docker API (measured in milliseconds, default: 500ms)
- Tab-Order: Set the order of tabs displayed, the keys must be `images`, `containers` and `volumes`. You can omit the names of the tabs you do not wish to see as well. Say I want to see `containers` tab first and do not want to see the `volumes` tab, I can set `Tab-Order: [containers, images]`
- Notification-Timeout: Set how long a status message sticks around for (measured in milliseconds, default: 2000ms)

## Roadmap

- Add a networks tab
- Make compatible with podman 👀

## Found an issue ?

Feel free to open a new issue, I will take a look ASAP.

## Contributing

Please refer [CONTRIBUTING.md](./CONTRIBUTING.md) for more info. 

## Thanks!!

![image](https://github.com/ajayd-san/gomanagedocker/assets/54715852/61be1ce3-c176-4392-820d-d0e94650ef01)
