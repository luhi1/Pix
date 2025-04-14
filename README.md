# The (*Pix*)el sorter!

Insert Gif Here

## A simple image manipulation tool written in GO to help graphic designers and photographers see compression and composition of their photos.

This project aims to solve the issue of figuring out the defining color/features in large and complex images. Using a custom algorithm (basically an insertion sort on a doubly-linked, doubly-circular linked list of pixels), this program sorts images into blocks based on their RGB color groups. This algorithm is optimized to be quickest with images that use a lot of the same colors, and is not best for use in sorting images with many, **highly different** color combinations.

## Installation
### Option 1: Precompiled Binary (MOST USERS)
Currently we only have a Windows binary available. Go over to the releases tab on this repository and download the executable. Run it when finished :)

### Option 2: Compile from Source (Advanced Users)
* Ensure you have golang >= 1.19
* Download the source code as a zip file from the releases tab on this repository.
* Unzip the archive.
* Build the project
  * In a terminal instance run:
      ```bash
      cd path/to/extracted/folder
      go build main.go
  * If you are running a non X11 distro, see [Known Issues](#known-issues) for needed dependencies.
* Finally, execute your binary.
  * Either a)
    * Find your binary and double click it.
  * or b)
    * In a terminal run:
        ```bash
        cd /path/to/folder/where/binary/is
        ./main.go  
## Usage
1. Select the file you want to sort. !!!!!!!!ONLY WORKS WITH .PNG FILES FOR NOW!!!!!!!!!!!!!!!!!!
2. Select your error range (see below)
3. Sort!

## Error Range in Depth
As a user, you can decide how similar colors can be in order for them to be grouped together using the "error range" selector. Larger error ranges are useful to show compression as only exact or close to exact duplicate pixels are put into the same category.   

For example:
* Starting Image:
  [Insert Image Here]
* Error Range 254 Output:
  [Insert Image Here]
* Error Range 1 Output:
  [Insert Image Here]
## Example speeds on my 8 Core, AMD Ryzen 5 3400G (GPU Rendering isn't supported in this release)

## Known Issues

### "Fatal error: X11/X[package]/X[package].h: No Such File Found" when compiling from source on non-X11 linux distros.
Please use your package manager of choice to install the X11 package and then recompile. The Fyne framework currently only supports X11 for rendering.

For debian-based OS's these command should work:

```bash
sudo apt install libxcursor-dev
sudo apt install libxrandr-dev
sudo apt install libxrinerama-dev
sudo apt install libxi-dev
sudo apt install libxxf86vm-dev
