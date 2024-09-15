# jsdr
A Software Defined Radio written in Go. The intention is to provide an application that
will run on Windows, Linux, and macOS.

## Project Status

This project is very early in its development, and therefore contains little or no useful
functionality. I have created this project to:

1. Learn programming in Go.
2. Learn Digital Signal Processing (DSP).
3. Use my SDR dongles.

Given my track record with large projects, this project will probably never be finished. I
have started this project using a number of different programming languages and tools and
have never found the right combination for my liking, so be forewarned!

## About the [License](LICENSE)

I would have preferred that this project be under a very permissive license, but one or more of the libraries used indirectly in this project is licensed under the GPL version 2
or later.
From my understanding of that license, any project that includes code licensed under one
or more versions of the GPL license must also be licensed under GPL.

## Building jsdr

As mentioned earlier, jsdr is being developed using the Go programming language and tools.
In addition to the Go toolset, I also use Visual Studio Code, so the build instructions
below assume that.

Development of jsdr is currently being done on Linux (specifically Kubuntu 24.04), so build
instructions will only be given for Unbuntu Linux at this time.

### Needed Development Tools

In order to build jsdr, a number of development tools are required on each build system:

1. System development tools
2. Go
3. Visual Studio Code
4. SoapySDR libraries

For each of the operating systems that can be used to build jsdr, how to install the
various required tools will be listed.

### Building jsdr on Ubuntu

The following instructions assume you are starting from a freshly installed system.

* sudo apt update
* sudo apt install build-essential libgl1-mesa-dev xorg-dev
* sudo apt install libsoapysdr-dev
* sudo snap install go
* sudo snap install code
* cd ~
* mkdir go
* cd go
* git clone https://github.com/jimorc/jsdr.git
  
A number of other libraries will be needed. Installation instructions will be added as
needed.
