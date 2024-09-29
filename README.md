 # np-run

A CLI tool that lets you browse and execute npm, pnpm, yarn, or bun scripts from a local project.

https://github.com/user-attachments/assets/348e5b90-ffee-4087-8c29-25f082eaa201

## Todo

- [x] Check for `package.json` in current directory / project
- [x] Read file and parse key and values in `scripts` object in `package.json`
- [x] Print scripts to stdout as a list - use BubbleTea
- [x] Detect the package manager - using the lockfile?
- [x] When enter is pressed for a chosen script, that script is run
- [ ] Look at how to properly structure a Go program
- [x] Sort scripts in alphabetical order

### Nice-To-Have Features
- [x] Install dependencies if `node_modules` isn't found
  - [x] Show spinner or progress bar while dependencies are installing
- [x] Print something if multiple lockfiles are found
