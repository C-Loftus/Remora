# Remora 

Remora is an application that runs in parallel with the Orca screen reader on Linux and extends Orca's functionality. It provides the ability to:

- Summarize text and visual content with private offline ollama vision models
- OCR the screen
- Provide hotkeys for helpful actions not in Orca by default (for example, a screen curtain)

## Quickstart Overview Video

A video demonstrating how to use Remora can be found [here](https://www.youtube.com/watch?v=Qh0w-DOH12s)

[![The thumbnail of a youtube video demonstrating a quick start of Remora](./frontend/src/assets/thumbnail.png)](https://www.youtube.com/watch?v=Qh0w-DOH12s)

## Installation

- Download [the latest release](https://github.com/C-Loftus/Remora/releases/latest) from Github or [build from source](#building) 
- To install Remora as a desktop application, you can run `make install` in the makefile.

## Limitations: 

- You must use a recent version of Orca (at least v49.0 beta or above)
- This app only works on X11; Wayland does not support global hotkeys or coordinates and thus is impossible to support across compositors
- To use hotkeys you cannot have numlock or another modifier key enabled
- Currently supports only amd64
- You must tab into the GUI window to read its contents with Orca due to [a bug in webkit gtk](https://gitlab.gnome.org/GNOME/orca/-/issues/493)

## Building

You must have [wails](https://wails.io/docs/gettingstarted/installation) installed to build this program. Then follow the commands in the [makefile](./makefile)

## Screenshot

![A screenshot of the remora application, showing the logo, the title page, an enumeration of keyboard shortcuts, ollama output, and ocr](./frontend/src/assets/images/screenshot.png)