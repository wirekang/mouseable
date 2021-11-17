# Motivation

This program was inspired by
[Ultimate Hacking Keyboard Demo](https://youtu.be/4rjnkHqnA3s?t=20),  
[Windows built-in function](https://support.microsoft.com/en-us/windows/use-mouse-keys-to-move-the-mouse-pointer-9e0c72c8-b882-7918-8e7b-391fd62adf33)
and [NeatMouse](https://github.com/neatdecisions/neatmouse).

## Improvement

Unlike Windows built-in, you can customize every single shortcut. And unlike
NeatMouse, there are no stammering or lag, and there are more customizable
options.  
Also mouseable is written in pure Go, so it's easy to install and maintain.

# Getting Started

## Install

```go install github.com/wirekang/mouseable/cmd/mouseable@latest```

## Register Service

Open terminal as **Admin**. Run command to register Windows Service. This step
may be optional if you don't want to use mouseable in background.

```mouseable -register```

## Edit Config

After register, mouseable was started automatically.  
Edit ```[Home Directory]/mouseable.json``` for
example ```C:/Users/user1/mouseable.json```

## All Flags

```
  -register
        Register and run service
  -reload
        Reload config file at  %USERPROFILE%/mouseable.json
  -run
        This flag run mouseable in foreground that usually NOT NEEDED
  -unregister
        Unregister windows service
```

# Roadmap

* [ ] Support double press
* [ ] UI based config
* [ ] UI based config
* [ ] Fine error handling
