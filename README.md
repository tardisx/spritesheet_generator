# spritesheet_generator

This is a simple program to generate a template for isometric sprite 
tilesheets, in png format. For use with software such as [Tiled](https://www.mapeditor.org).

## Installation

`go install github.com/tardisx/spritesheet_generator@latest`

If you really want binaries, pester me :-)

## Usage

`Usage of ./spritesheet_generator:
  -height int
    	base tile height in pixels (default 128)
  -multiplier int
    	tile height multiplier (default 2)
  -output string
    	output filename
  -width int
    	tile width in pixels (default 128)
  -x int
    	number of tiles across (default 8)
  -y int
    	number of tiles down (default 8)`

Hopefully these options are mostly self-explanatory.

The `multiplier` option describes how 'tall' the tiles are. Normally you
want some height to tiles to give them the illusion of depth and the ability
to hide things behind them. If you are unsure, start with the default of 2.

## Example

![Screenshot][screenshot]


[screenshot]: https://raw.githubusercontent.com/tardisx/spritesheet_generator/main/example.png
