# YAMM - Yet Another Media Manager

YAMM is a tool to help you organize your media library.

Some of its features:
 - Import any media to its library
 - Deploy media to your media provider (for example Jellyfin)
 - Transcode media
 - Move media between volumes
 - Automatically organize media with ability to fix it by hand
 - TUI for browsing your library

YAMM is lightweight and tries to be as simple as possible. It doesn't have any
daemon, background server nor tasks. Everything is contained within
one simple utility - `yamm` executable.

For now it is targeting `linux`.

## How does it work?

YAMM organize your media and stores it inside its database. It can have multiple
databases - f.e. one for each disk. Media lives inside that databases with
metadata stored inside YAMM master database. From there you can deploy media to
your media provider and create links in the desired folder structure. Also you
can transcode media, rename it and store multiple versions of the same video!

To import media to the YAMM you just invoke `yamm import` on the directory and
it tries to gather metadata from folder structure and filenames. If it can't you
can easily change it by hand with TUI.

## How to use it?

You have two ways of using YAMM:

First one is TUI. Just invoke `yamm` command and you will fall into TUI where
you can manage your library.

Second one is CLI interface. Use `yamm --help` for more info. 