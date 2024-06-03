# np-run

A CLI tool that lets you browse and execute npm, pnpm, yarn, or bun scripts directly from the command line. No more memorising commands - just pick, run, and go!

## Todo

- [ ] Check for `package.json` in current directory / project
- [ ] Read file and parse key and values in `scripts` object in `package.json`
- [ ] Print scripts to stdout as a list - use BubbleTea
- [ ] Detect the package manager - using the lockfile?
- [ ] When enter is pressed for a chosen script, that script is run

### Nice-To-Have Features
- [ ] Install dependencies if `node_modules` isn't found
