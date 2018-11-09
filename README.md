# Basic photo management (very WIP)

[![Build Status](https://travis-ci.org/internetimagery/photos.svg?branch=master)](https://travis-ci.org/internetimagery/photos)

[![Coverage Status](https://coveralls.io/repos/github/internetimagery/photos/badge.svg?branch=master)](https://coveralls.io/github/internetimagery/photos?branch=master)

# TODO:

- fix up tagging, add proper tests in main_test
- use standard file structure for importing images, and placing them etc
- make tagging use a web serice, with a basic website for interractive tagging
- add image duplication checker like phash
- add video phash check, if possible

#### Nice to have
- manage all actions with an interface (ie file movements, renames etc)
- track those actions and allow an undo system to exist
- autocomplete actions

### Intended use:

```
photos init "some name"
```

Initialize working directory as project root. Create a configuration file at that location that serves to mark the root, and also provides a space to add custom commands for backups / compression.

```
photos sort
```

Take all loose media in working directory, and add them to a folders named based on their date. Use EXIF where available.

```
photos rename
```

Run through all files within the working directory.
Use the parent directory name as the namespace (or event) and checks the filenames against a predetermined format ("event_index[tag tag].ext"). Files that do not match this format are determined to be newly added, and are renamed. If a compression command is provided in the config file, this will be run on the file (ie mozjpeg "$SOURCEPATH" > "$DESTPATH").

```
photos backup name
```

Runs the named backup command (from the config) on files/dirs in the current working directory, relative to the root directory, etc. Allows for quick shortcuts/aliases to otherwise more complicated backup code.

### Environment Vars

Commands run within the config file inherit the parent commands environment. Additional contextural variable names will be added. ie:

$SOURCEPATH = path to source file
$DESTPATH = path to destination file, which should not yet exist

A command such as the following, would perform a basic copy.

```
cp "$SOURCEPATH" "$DESTPATH"
```
