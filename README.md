# Basic photo management (very WIP)

[![Build Status](https://travis-ci.org/internetimagery/photos.svg?branch=master)](https://travis-ci.org/internetimagery/photos) [![Coverage Status](https://coveralls.io/repos/github/internetimagery/photos/badge.svg?branch=master)](https://coveralls.io/github/internetimagery/photos?branch=master)

# TODO:

- add --copy flag to "sort", to leave media in its original place and only copy across info
- check and see if os.Rename works across drives. Create a test if it doesn't.
- add phash check before and after compression as a failsafe to ensure the copy is faithful to the original (and a test)
- add more tests for things like bad data
- stop all commands from being run in root
- make a tag searching function to collect all tags. Provide prompt for spelling errors against similar tags.
- make tagging use a web serice, with a basic website for interractive tagging
- add image duplication checker like phash
- add video phash check, if possible

#### Nice to have
- Test for, and create functionality for the situation where premature shutdown happens (ctrl+c / power outage)
- manage all actions with an interface (ie file movements, renames etc)
- track those actions and allow an undo system to exist
- autocomplete actions
- run compressions in a few threads

### Intended use:

#### (1) Init

```
photos init "some name"
```

Before doing anything we need to initialize a directory on our computer that will contain our photos etc. To do this we use the above command, which creates a file (photos-config.yaml) in the current shell working directory.

This file serves to both let the "photos" program know where the root of this project lies, and also contains settings and commands to be run for compressing / backing up.

#### (2) Adding media, sorting

```
photos sort [--copy] <filename> <filename...>
```

The next step is to pull in some media. The sort command above will do the trick. You must use it (like all commands) from within your project. It will grab all the files from the specified directory and put them into a "sorted" directory with the project itself. Each file within a sub-directory sorted according to date.

The intention then is to manually go through the images and put them into a folder structure that makes sense. Also naming the created folders as events is a nice way to go. A useful format can be:

```
project-root / year (2018) / event (18-10-10 eventname) / media
```
If you wish to keep the original media in the directory it was found, add the "--copy" flag to copy the files instead of moving them.


#### (3) Format / compress media, rename

```
photos rename
```

By design files follow a strict naming scheme. They take an element from the directory they reside in, are given an id, and can have tags. Files that do not follow this scheme are assumed to have not yet been added/compressed.

Running the above command will format the names of all files not already formatted in the working directory of the shell. It will also run any compression commands specified in the photos-config.yaml file, based on pattern matching of the filename against the name of the command. For instance, a useful setup using mozjpeg and ffmpeg to compress jpeg and mp4 media could look like this:

```
compress:
  -
    name: "*.jpg *.jpeg"
    command: "mozjpeg -quality 80 '$SOURCEPATH' > '$DESTPATH'"
  -
    name: "*.mp4"
    command: "ffmpeg -i '$SOURCEPATH' -crf 23 '$DESTPATH'"
```

There are environment variables set for use in these commands as they are run. Typically all you'd want is SOURCEPATH and DESTPATH here.

All original media (regardless of if compression happens or not) will be moved into a temporary folder. If you see anything wrong with your renamed and perhaps compressed files, you can easily bring back the original. Once you're happy with the changes however, feel free to delete the originals folder.

#### (4) Tag media

```
photos tag [--remove] <filename> <filename...> -- <tag> <tag...>
```

Now your media is formatted and compressed (optionally). It's time to tag it. This is an important step, because as we are using the filename to store the tags, it makes it quite fragile if we want to edit the tags later thus renaming files. So we want to be as thoughtful as we can about what we want to use in the tags to make these files searchable up front.

An idea is to tag the people seen in the media.
By format convention the filenames inherit the name of their parent directory. So any searchable terms / event names you want to apply to the entire collection can reside there. Leaving media specific info for the tags.

Tagging has a useful shorthand too. You can use the full filename to refer to the file or just its index (assuming the file is within the same working directory of the shell.) ie the following two commands are identical:

```
photos tag "18-01-10 Event_004.jpg" person
photos tag 4 person
```

And now the file will look like this:

```
18-01-10 Event_004[person].jpg
```

#### (4.5) Locking

```
photos lock [--force]
```

As the goal of this system is to protect the data stored within. It makes sense once everything is compressed, named and tagged to lock it all down. The above command does just that.

Upon running the command you'll get a new file "locked.yaml" and all the formatted files within the folder will have become read only. The "locked.yaml" file contains a snapshot of information about the contents of the files themselves at the time of locking.
Running the command on an already locked directory will perform a cross check between the files current state and that of the snapshot to ensure no changes have occurred since locking. If changes have been made, a warning will pop up.

Using the "--force" flag will suppress any warning about changes and update the snapshot data with the current state of the file. Only use this if you know the files current state is what you want to keep. Generally speaking if something changed, you might want to look at a backup of the file to see what the difference is.

To "unlock" the files, just delete the "locked.yaml" file.

Locking in this way is somewhat of an optional step as performing a backup will first run this locking system ahead of the backup process. Thus bailing out if files have changed without you allowing it specifically (--force).

#### (5) Backup

```
photos backup name
```

Runs the provided backup command on the working directory of the shell. The command to run is provided in the "photos-config.yaml" file and is selected via the "name" argument. The name selector can include * characters to match more than one command. For instance "backup-s3" and "backup-b2" would both be run if "name" argument is "backup*".

An example backup solution is the amazing rclone. It can be set up independently of this system, and a command can be added like this:

```
backup:
  -
    name: "b2"
    command: "rclone copy \"$SOURCEPATH\" \"backup:my-bucket/$RELPATH\" -v"
```

Prior to the backup taking place, a lock command is run on the files. This both locks files (see section above) and also checks if they have changed since they were last locked. If any files are found to have been changed, the backup will abort as a safety measure. If the files changing was an intentional situation, you will need to run the lock command above with the "--force" flag to update the lock, then re-run the backup.

This may seem convoluted, but it's a step towards ensuring a corrupt or otherwise unintended change is not backed up and overrides a "real" copy of a file.
