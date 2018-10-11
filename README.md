# photos

[![Build Status](https://travis-ci.org/internetimagery/photos.svg?branch=master)](https://travis-ci.org/internetimagery/photos)

photos management wip

rebuild.

[ ] Rebuild config to create empty configs, to load configs, and to hunt for configs. (should hunting for a file be a utility?)

Plan:

two main commands to start. simplify.

"rename", "backup"

## rename:
* take folder in working directory
* build regex around folder name (event name)
* apply to all entries in folder. assume anything not matching is a new item.
* compress and rename files to fit format
* offer flag to prevent compression
* offer flag to set compression level (maybe for later, keep it simple)

## backup:
* run backup command on working directory folder relative to root
* require specification on which backup to use

## global:
* create config file that houses information for commands, like the compression and backup options
* force it to be manually created? offer command to create it automatically? prepopulate?
