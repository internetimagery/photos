# Basic photo management (very WIP)

[![Build Status](https://travis-ci.org/internetimagery/photos.svg?branch=master)](https://travis-ci.org/internetimagery/photos)

### Intended useage:

```
photos init "some name"
```

Initialize working directory as project root. Create a configuration file at that location that serves to mark the root, and also provides a space to add custom commands for backups / compression.

```
photos rename
```

Runs through all files within the working directory. Uses the parent directory name as the namespace (or event) and checks the filenames against a predetermined format ("event_index[tag tag].ext"). Files that do not match this format are determined to be new, and are renamed. If a compression command is provided in the config file, this will be run on the file.

```
photos backup name
```

Runs the named backup command (from the config) providing variables for the current working directory, and root directory, etc. Allows for quick shortcuts/aliases to otherwise more complicated code.

### Environment Vars

Commands run within the config file inherit the parent commands environment. However variable names "$var" will be expanded upon in a separate pass with contextural info. ie:

$SOURCEPATH = path to source file
$DESTPATH = path to destination file, which should not yet exist

A command such as the following, would perform a basic copy.

```
cp "$SOURCEPATH" "$DESTPATH"
```
