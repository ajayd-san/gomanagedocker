# goManageDocker

Do Docker commands slip your mind because you don't use Docker often enough? Sick of googling commands for everyday tasks? Worry no more!

Introducing **goManageDocker** (get it?)! This blazing fast TUI, made using Go and BubbleTea, will make managing your Docker objects a breeze. 

## Contents
1. [Install Instructions](#install-instructions)
2. [Features](#features)
3. [Roadmap](#roadmap)
4. [Found an issue?](#found-an-issue-)

## Install Instructions

Building from source is pretty easy: 

```
go install github.com/ajayd-san/gomanagedocker@v1.0.1
```

Now, **goManageDocker 😏!!**

>[!NOTE]
>goManageDocker runs best on terminals that support ANSI 256 colors and designed to run while the **terminal is maximized**.

## Features

1. Easy navigation with vim keybinds and arrow keys.

  ![intro](https://github.com/ajayd-san/gomanagedocker/assets/54715852/00bf4e8e-44fa-417c-a8cf-7cbccd687ad6)

2. Exec into selected container with A SINGLE KEYSTROKE: `x`...How cool is that?

![exec](https://github.com/ajayd-san/gomanagedocker/assets/54715852/b168b3d7-75f5-4339-884e-573a6e6fb688)


3. Delete objects using `d` (You can force delete with `D`, you won't have to answer a prompt this way)
   
  ![delete](https://github.com/ajayd-san/gomanagedocker/assets/54715852/a4b54c6c-11ad-4ed8-9111-ffad85567188)

4. Prune objects using `p`
   
  ![prune](https://github.com/ajayd-san/gomanagedocker/assets/54715852/1ff3809d-d08e-4200-b00b-aefc7b9f2485)

5. start/stop/pause/restart containers with `s`, `t` and `r`
   
  ![startstop](https://github.com/ajayd-san/gomanagedocker/assets/54715852/3e54bc51-1d7c-4669-8f8e-18eae0ca18bf)

6. Filter objects with `/`

  ![search](https://github.com/ajayd-san/gomanagedocker/assets/54715852/513564e5-dacf-4f8a-8eca-c575dcfe6be2)


## Roadmap
- Make the program work with minimized terminal state
- Add a networks tab

## Found an issue ?

Feel free to open a new issue, I will take a look ASAP.

![image](https://github.com/ajayd-san/gomanagedocker/assets/54715852/61be1ce3-c176-4392-820d-d0e94650ef01)


