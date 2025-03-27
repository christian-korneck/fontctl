# fontctl

fontctl - CLI to install or uninstall a font on MS Windows.

![](logo.png)

## Usage

see auto-generated [cli docs](cli-docs.md) or run `fontctl --help`:

```
NAME:
   fontctl - Install or uninstall a font on MS Windows

USAGE:
   fontctl [global options] [command [command options]]

COMMANDS:
   install    Install a font
   uninstall  Uninstall a font
   getname    Get the font name from a file
   load       Load a font into memory
   unload     Unload a font from memory
   refresh    Refresh known fonts for current user session
   preview    Preview a font
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d  enable verbose debug logging (default: false)
   --help, -h   show help
   ```

## Status

WIP / alpha version.    
    
All current features are working, but needs better handling for edge cases, batch processing, etc and needs overall cleanup/refactoring.    
    
Both CLI and Go API will likely change.

## Concept

### What

`fontctl` implements the win32 gdi font installation ["whitepaper"](https://learn.microsoft.com/en-us/windows/win32/gdi/font-installation-and-deletion) with Go ([cgo-free syscalls](https://www.youtube.com/watch?v=EsPcKkESYPA)) using documented and [undocumented](https://github.com/reactos/reactos/blob/88d9285bd01cc20b742b94e5412d5a7983d71296/sdk/include/reactos/undocgdi.h#L46) windows functions.


### Goals
- automation friendly (make it easy to install fonts from shell scripts, etc)
- renderfarm friendly (temporarily load/unload fonts without installing them)
- primarily manages things that you can do with font *files*

### Non-Goals
- manage things *not* related to font *files* ~~(advanced management of the in-memory windows font table~~, etc)
-  being a ~~full featured font manager~~ (have a look at i.e. [fontbase](https://fontba.se/) instead)

