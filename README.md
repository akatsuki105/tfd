# tfd

<img src="https://travis-ci.com/Akatsuki-py/tfd.svg?branch=master" /> <a href="https://godoc.org/github.com/Akatsuki-py/tfd"><img src="https://godoc.org/github.com/Akatsuki-py/tfd?status.svg" alt="GoDoc"></a>

Golang API for text-base file dialog

Demo

<img src="https://imgur.com/4DoSCUR.gif" />



Multi-platform (Win10 and Xubuntu18.04 is checked)

<img src="https://imgur.com/EZWy99a.png" />

## Usage

install

```
go get github.com/akatsuki-py/tfd
```

This package has only one API.

```go
extensions := []string{"md", "txt"}                             // select *.md or *.txt file
hiddenDisplay := false                                          // if hidden-file displays 
path, err := tfd.CreateSelectDialog(extensions, hiddenDisplay)  // e.g. /usr/local/go/README.md
```

## Command 

You can use the <kbd>Arrow Keys</kbd> to navigate, <kbd>Enter</kbd> to select a file or enter directory ,<kbd>Esc</kbd> to quit.