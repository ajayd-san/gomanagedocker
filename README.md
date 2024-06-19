# goManageDocker

Do Docker commands slip your mind because you don't use Docker often enough? Sick of googling commands for everyday tasks? GoManageDocker is designed to NUKE this annoyance. 

Introducing **goManageDocker** (get it?)! This blazing fast TUI, made using Go and BubbleTea, will make managing your Docker objects a breeze. 

## Contents
1. [Install Instructions](#install-instructions)
2. [Features](#features)
3. [Roadmap](#roadmap)
4. [Found an issue?](#found-an-issue-)

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
go install github.com/ajayd-san/gomanagedocker@v1.1.1
```
Start the program with `gmd`.

### Windows 

You can get the latest precompiled binary from releases or you may build from source. 

Now, **goManageDocker 😏!!**

>[!NOTE]
>goManageDocker runs best on terminals that support ANSI 256 colors and designed to run while the **terminal is maximized**.

## Features

**New in v1.1.1**: Ability to perform `docker scout quickview` on any image (just press `s`).

![scout](vhs/gifs/scout.gif)

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

## Roadmap
- Make the program work with minimized terminal state
- Add a networks tab

## Found an issue ?

Feel free to open a new issue, I will take a look ASAP.

![image](https://github.com/ajayd-san/gomanagedocker/assets/54715852/61be1ce3-c176-4392-820d-d0e94650ef01)


