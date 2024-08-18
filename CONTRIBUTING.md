## Contributing Guidelines

Hey 👋, thanks for your interest in contributing to goManageDocker. 

You can contribute in 3 ways: 

1. **Report a Bug:** goManageDocker is still new and quite a few bugs might be snooping around! If you find a bug, you can [open an issue](https://github.com/ajayd-san/gomanagedocker/issues/new?assignees=&labels=&projects=&template=bug_report.md&title=BUG+%F0%9F%90%9E%3A) and report the bug, try to be as descriptive about the bug as possible, it genuinely helps.
2. **Request a feature:** Do you have an idea for a feature that could be a great addition to goManageDocker's blazing-fast (🦀) arsenal? [Open a feature request!](https://github.com/ajayd-san/gomanagedocker/issues/new?assignees=&labels=&projects=&template=feature_request.md&title=)
3. **Work on a feature/issue:** Do you want to pick up a feature/issue? Then keep reading!

## Setting up the development environment: 

### Requirements: 
- Go 1.22.2
- Docker (not really needed, but helpful to debug and final testing)
- just:  [Install Here](https://github.com/casey/just) (command runner, is nice to have, since this lets you run commands quickly) 
- dlv (debugger, optional)

### Setting up the Development environment
1. Fork this repository
2. Clone the forked repo to your local machine
3. Hack away!

### Debugging 
I've added a debug flag `--debug` that writes logs to `./gmd_debug.log`. It is helpful while debugging the TUI quickly (for instance, making sure the control flow works as intended) without running a whole debugger (delve), to write a log just put `log.Println({LOG})` where you find fit. 

Make sure to run using `go run main.go --debug` to enable logging. 

### Create a PR. 
You have now made your changes, congrats! Open a PR and try to be as descriptive as possible, include pertinent information such as a detailed description, issue ID (if it fixes an issue), etc. Adding a test case (if applicable) will reinforce confidence and make the PR move faster.


### Justfile

There are multiple recipes in the included [justfile](./justfile): 

1. `just run`: compiles and runs the project with `--debug` flag
2. `just test`: runs all tests, across all packages
3. `just build`: builds the binary(Not used much)
4. `just race`: runs the project with the race detector. Pretty useful to detect race conditions. 
5. `just debug-server`: starts a `dlv` debug server at `localhost:43000`
6. `just debug-connect`: connects to an existing debug session running at `localhost:43000`. Do `just debug-server` before running this. 
