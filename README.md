# Remora 

<p align="center">
  <img src="./frontend/src/assets/images/remora.png" alt="The Remora Logo; an Orca whale with a fish swimming below it" width="40%"/>
</p>

Remora swims alongside the Orca screen reader and uses IPC to:
- Provide OCR
- Interact with Ollama Vision Models
- Provide quick hotkeys for common settings

## Installation

Download a precompiled binary from Github or [build from source](#building)

## Limitations: 

- You must use a recent version of Orca (at least v49.0 beta or above); you may need to build from source
- This app only works on X11; Wayland does not support global hotkeys or coordinates and thus is impossible to support across compositors
- To use hotkeys you cannot have numlock or another modifier key enabled

## Building

You must have [wails](https://wails.io/docs/gettingstarted/installation) installed to build this program. Then follow the commands in the [makefile](./makefile)