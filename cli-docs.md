# fontctl

# NAME

fontctl - Install or uninstall a font on MS Windows

# SYNOPSIS

fontctl

```
[--debug|-d]
[--help|-h]
```

# DESCRIPTION

Copyright (C) 2025 Christian Korneck <christian@korneck.de>

**Usage**:

```
fontctl [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--debug, -d**: enable verbose debug logging

**--help, -h**: show help


# COMMANDS

## install

Install a font

>fontctl install [--systemwide] <Font File>

**--help, -h**: show help

**--systemwide, -s**: Install in the system font dir (default: install in the current user's userprofile). Requires Admin privileges.

### help, h

Shows a list of commands or help for one command

## uninstall

Uninstall a font

>fontctl uninstall [--systemwide] <Font File>

**--help, -h**: show help

**--systemwide, -s**: Uninstall from the system font dir (default: uninstall from the user font dir). Requires Admin privileges.

### help, h

Shows a list of commands or help for one command

## getname

Get the font name from a file

>fontctl getname <Font File>

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## load

Load a font into memory

>fontctl load <Font File>

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## unload

Unload a font from memory

>fontctl unload <Font File>

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## refresh

Refresh known fonts for current user session

>fontctl refresh

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## preview

Preview a font

**--help, -h**: show help

### file

Preview a font file using Windows Font Viewer

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### font

Preview a loaded font using Windows GDI

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## help, h

Shows a list of commands or help for one command

