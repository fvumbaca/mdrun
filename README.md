# `mdrun`

> `mdrun` is still very much  under development and not quite ready for serious
scrutiny (even before considering it for production use). At this stage, I am
still experimenting with ideas and experiences to achieve my goals with the
project. Nothing here is set in stone and might change overnight.

`mdrun` is a tool for getting the most out of your technical documentation. It
allows you to live evaluate code blocks in your markdown with _no modification_
to your documents.

`mdrun` is not just great as a teaching tool or impressing your friends. It can
also be used as a powerful runbook and incident response tool by granting
responders diagnosis and remediation processes at the click of a button.


## Install

Install `mdrun` from source with the Go tool chain:

```sh
go install github.com/fvumbaca/mdrun
```

## Usage

The simplest and most common way to use `mdrun` is to just run it from the command line:

```sh
mdrun
```

This will startup a server on your localhost allowing you to browse files in
the current directory.

## Still TODO

- Support basic cli flags like `--port`
- Support metadata at the top of the file in front matter
- Support required cli-tools to throw errors early if it can even run on the
system
- Support out-of-the-box containerization??
- CLI mode will print out the runbook to stdout and prompt before each exec
- Script mode that will execute markdown as a script
- Infer language dependencies from code blocks
