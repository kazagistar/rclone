% rclone(1) User Manual
% Nick Craig-Wood
% Jan 31, 2016

Rclone
======

[![Logo](http://rclone.org/img/rclone-120x120.png)](http://rclone.org/)

Rclone is a command line program to sync files and directories to and from

  * Google Drive
  * Amazon S3
  * Openstack Swift / Rackspace cloud files / Memset Memstore
  * Dropbox
  * Google Cloud Storage
  * Amazon Cloud Drive
  * Microsoft One Drive
  * Hubic
  * Backblaze B2
  * Yandex Disk
  * The local filesystem

Features

  * MD5/SHA1 hashes checked at all times for file integrity
  * Timestamps preserved on files
  * Partial syncs supported on a whole file basis
  * Copy mode to just copy new/changed files
  * Sync (one way) mode to make a directory identical
  * Check mode to check for file hash equality
  * Can sync to and from network, eg two different cloud accounts

Links

  * [Home page](http://rclone.org/)
  * [Github project page for source and bug tracker](http://github.com/ncw/rclone)
  * <a href="https://google.com/+RcloneOrg" rel="publisher">Google+ page</a></li>
  * [Downloads](http://rclone.org/downloads/)

Install
-------

Rclone is a Go program and comes as a single binary file.

[Download](http://rclone.org/downloads/) the relevant binary.

Or alternatively if you have Go installed use

    go get github.com/ncw/rclone

and this will build the binary in `$GOPATH/bin`.  If you have built
rclone before then you will want to update its dependencies first with
this (remove `-f` if using go < 1.4)

    go get -u -v -f github.com/ncw/rclone/...

See the [Usage section](http://rclone.org/docs/) of the docs for how to use rclone, or
run `rclone -h`.

linux binary downloaded files install example
-------

    unzip rclone-v1.17-linux-amd64.zip
    cd rclone-v1.17-linux-amd64
    #copy binary file
    sudo cp rclone /usr/sbin/
    sudo chown root:root /usr/sbin/rclone
    sudo chmod 755 /usr/sbin/rclone
    #install manpage
    sudo mkdir -p /usr/local/share/man/man1
    sudo cp rclone.1 /usr/local/share/man/man1/
    sudo mandb

Configure
---------

First you'll need to configure rclone.  As the object storage systems
have quite complicated authentication these are kept in a config file
`.rclone.conf` in your home directory by default.  (You can use the
`--config` option to choose a different config file.)

The easiest way to make the config is to run rclone with the config
option:

    rclone config

See the following for detailed instructions for

  * [Google drive](http://rclone.org/drive/)
  * [Amazon S3](http://rclone.org/s3/)
  * [Swift / Rackspace Cloudfiles / Memset Memstore](http://rclone.org/swift/)
  * [Dropbox](http://rclone.org/dropbox/)
  * [Google Cloud Storage](http://rclone.org/googlecloudstorage/)
  * [Local filesystem](http://rclone.org/local/)
  * [Amazon Cloud Drive](http://rclone.org/amazonclouddrive/)
  * [Backblaze B2](http://rclone.org/b2/)
  * [Hubic](http://rclone.org/hubic/)
  * [Microsoft One Drive](http://rclone.org/onedrive/)
  * [Yandex Disk](http://rclone.org/yandex/)

Usage
-----

Rclone syncs a directory tree from one storage system to another.

Its syntax is like this

    Syntax: [options] subcommand <parameters> <parameters...>

Source and destination paths are specified by the name you gave the
storage system in the config file then the sub path, eg
"drive:myfolder" to look at "myfolder" in Google drive.

You can define as many storage paths as you like in the config file.

Subcommands
-----------

### rclone copy source:path dest:path ###

Copy the source to the destination.  Doesn't transfer
unchanged files, testing by size and modification time or
MD5SUM.  Doesn't delete files from the destination.

Note that it is always the contents of the directory that is synced,
not the directory so when source:path is a directory, it's the
contents of source:path that are copied, not the directory name and
contents.

If dest:path doesn't exist, it is created and the source:path contents
go there.

For example

    rclone copy source:sourcepath dest:destpath

Let's say there are two files in sourcepath

    sourcepath/one.txt
    sourcepath/two.txt

This copies them to

    destpath/one.txt
    destpath/two.txt

Not to

    destpath/sourcepath/one.txt
    destpath/sourcepath/two.txt

If you are familiar with `rsync`, rclone always works as if you had
written a trailing / - meaning "copy the contents of this directory".
This applies to all commands and whether you are talking about the
source or destination.

### rclone sync source:path dest:path ###

Sync the source to the destination, changing the destination
only.  Doesn't transfer unchanged files, testing by size and
modification time or MD5SUM.  Destination is updated to match
source, including deleting files if necessary.

**Important**: Since this can cause data loss, test first with the
`--dry-run` flag to see exactly what would be copied and deleted.

Note that files in the destination won't be deleted if there were any
errors at any point.

It is always the contents of the directory that is synced, not the
directory so when source:path is a directory, it's the contents of
source:path that are copied, not the directory name and contents.  See
extended explanation in the `copy` command above if unsure.

If dest:path doesn't exist, it is created and the source:path contents
go there.

### rclone ls remote:path ###

List all the objects in the the path with size and path.

### rclone lsd remote:path ###

List all directories/containers/buckets in the the path.

### rclone lsl remote:path ###

List all the objects in the the path with modification time,
size and path.

### rclone md5sum remote:path ###

Produces an md5sum file for all the objects in the path.  This
is in the same format as the standard md5sum tool produces.

### rclone sha1sum remote:path ###

Produces an sha1sum file for all the objects in the path.  This
is in the same format as the standard sha1sum tool produces.

### rclone size remote:path ###

Prints the total size of objects in remote:path and the number of
objects.

### rclone mkdir remote:path ###

Make the path if it doesn't already exist

### rclone rmdir remote:path ###

Remove the path.  Note that you can't remove a path with
objects in it, use purge for that.

### rclone purge remote:path ###

Remove the path and all of its contents.  Note that this does not obey
include/exclude filters - everything will be removed.  Use `delete` if
you want to selectively delete files.

### rclone delete remote:path ###

Remove the contents of path.  Unlike `purge` it obeys include/exclude
filters so can be used to selectively delete files.

Eg delete all files bigger than 100MBytes

Check what would be deleted first (use either)

    rclone --min-size 100M lsl remote:path
    rclone --dry-run --min-size 100M delete remote:path

Then delete

    rclone --min-size 100M delete remote:path

That reads "delete everything with a minimum size of 100 MB", hence
delete all files bigger than 100MBytes.

### rclone check source:path dest:path ###

Checks the files in the source and destination match.  It
compares sizes and MD5SUMs and prints a report of files which
don't match.  It doesn't alter the source or destination.

### rclone dedupe remote:path ###

Interactively find duplicate files and offer to delete all but one or
rename them to be different. Only useful with Google Drive which can
have duplicate file names.

```
$ rclone dedupe drive:dupes
2016/01/31 14:13:11 Google drive root 'dupes': Looking for duplicates
two.txt: Found 3 duplicates
  1:       564374 bytes, 2016-01-31 14:07:22.159000000, md5sum 7594e7dc9fc28f727c42ee3e0749de81
  2:      1744073 bytes, 2016-01-31 14:07:12.490000000, md5sum 851957f7fb6f0bc4ce76be966d336802
  3:      6048320 bytes, 2016-01-31 14:07:02.111000000, md5sum 1eedaa9fe86fd4b8632e2ac549403b36
s) Skip and do nothing
k) Keep just one (choose which in next step)
r) Rename all to be different (by changing file.jpg to file-1.jpg)
s/k/r> r
two-1.txt: renamed from: two.txt
two-2.txt: renamed from: two.txt
two-3.txt: renamed from: two.txt
one.txt: Found 2 duplicates
  1:         6579 bytes, 2016-01-31 14:05:01.235000000, md5sum 2b76c776249409d925ae7ccd49aea59b
  2:         6579 bytes, 2016-01-31 12:50:30.318000000, md5sum 2b76c776249409d925ae7ccd49aea59b
s) Skip and do nothing
k) Keep just one (choose which in next step)
r) Rename all to be different (by changing file.jpg to file-1.jpg)
s/k/r> k
Enter the number of the file to keep> 2
one.txt: Deleted 1 extra copies
```

The result being

```
$ rclone lsl drive:dupes
   564374 2016-01-31 14:07:22.159000000 two-1.txt
  1744073 2016-01-31 14:07:12.490000000 two-2.txt
  6048320 2016-01-31 14:07:02.111000000 two-3.txt
     6579 2016-01-31 12:50:30.318000000 one.txt
```

### rclone config ###

Enter an interactive configuration session.

### rclone help ###

Prints help on rclone commands and options.

Server Side Copy
----------------

Drive, S3, Dropbox, Swift and Google Cloud Storage support server side
copy.

This means if you want to copy one folder to another then rclone won't
download all the files and re-upload them; it will instruct the server
to copy them in place.

Eg

    rclone copy s3:oldbucket s3:newbucket

Will copy the contents of `oldbucket` to `newbucket` without
downloading and re-uploading.

Remotes which don't support server side copy (eg local) **will**
download and re-upload in this case.

Server side copies are used with `sync` and `copy` and will be
identified in the log when using the `-v` flag.

Server side copies will only be attempted if the remote names are the
same.

This can be used when scripting to make aged backups efficiently, eg

    rclone sync remote:current-backup remote:previous-backup
    rclone sync /path/to/files remote:current-backup

Options
-------

Rclone has a number of options to control its behaviour.

Options which use TIME use the go time parser.  A duration string is a
possibly signed sequence of decimal numbers, each with optional
fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid
time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

Options which use SIZE use kByte by default.  However a suffix of `k`
for kBytes, `M` for MBytes and `G` for GBytes may be used.  These are
the binary units, eg 2\*\*10, 2\*\*20, 2\*\*30 respectively.

### --bwlimit=SIZE ###

Bandwidth limit in kBytes/s, or use suffix k|M|G.  The default is `0`
which means to not limit bandwidth.

For example to limit bandwidth usage to 10 MBytes/s use `--bwlimit 10M`

This only limits the bandwidth of the data transfer, it doesn't limit
the bandwith of the directory listings etc.

### --checkers=N ###

The number of checkers to run in parallel.  Checkers do the equality
checking of files during a sync.  For some storage systems (eg s3,
swift, dropbox) this can take a significant amount of time so they are
run in parallel.

The default is to run 8 checkers in parallel.

### -c, --checksum ###

Normally rclone will look at modification time and size of files to
see if they are equal.  If you set this flag then rclone will check
the file hash and size to determine if files are equal.

This is useful when the remote doesn't support setting modified time
and a more accurate sync is desired than just checking the file size.

This is very useful when transferring between remotes which store the
same hash type on the object, eg Drive and Swift. For details of which
remotes support which hash type see the table in the [overview
section](http://rclone.org/overview/).

Eg `rclone --checksum sync s3:/bucket swift:/bucket` would run much
quicker than without the `--checksum` flag.

When using this flag, rclone won't update mtimes of remote files if
they are incorrect as it would normally.

### --config=CONFIG_FILE ###

Specify the location of the rclone config file.  Normally this is in
your home directory as a file called `.rclone.conf`.  If you run
`rclone -h` and look at the help for the `--config` option you will
see where the default location is for you.  Use this flag to override
the config location, eg `rclone --config=".myconfig" .config`.

### --contimeout=TIME ###

Set the connection timeout. This should be in go time format which
looks like `5s` for 5 seconds, `10m` for 10 minutes, or `3h30m`.

The connection timeout is the amount of time rclone will wait for a
connection to go through to a remote object storage system.  It is
`1m` by default.

### -n, --dry-run ###

Do a trial run with no permanent changes.  Use this to see what rclone
would do without actually doing it.  Useful when setting up the `sync`
command which deletes files in the destination.

### --ignore-existing ###

Using this option will make rclone unconditionally skip all files
that exist on the destination, no matter the content of these files.

While this isn't a generally recommended option, it can be useful
in cases where your files change due to encryption. However, it cannot
correct partial transfers in case a transfer was interrupted.

### --log-file=FILE ###

Log all of rclone's output to FILE.  This is not active by default.
This can be useful for tracking down problems with syncs in
combination with the `-v` flag.

### --modify-window=TIME ###

When checking whether a file has been modified, this is the maximum
allowed time difference that a file can have and still be considered
equivalent.

The default is `1ns` unless this is overridden by a remote.  For
example OS X only stores modification times to the nearest second so
if you are reading and writing to an OS X filing system this will be
`1s` by default.

This command line flag allows you to override that computed default.

### -q, --quiet ###

Normally rclone outputs stats and a completion message.  If you set
this flag it will make as little output as possible.

### --retries int ###

Retry the entire sync if it fails this many times it fails (default 3).

Some remotes can be unreliable and a few retries helps pick up the
files which didn't get transferred because of errors.

Disable retries with `--retries 1`.

### --size-only ###

Normally rclone will look at modification time and size of files to
see if they are equal.  If you set this flag then rclone will check
only the size.

This can be useful transferring files from dropbox which have been
modified by the desktop sync client which doesn't set checksums of
modification times in the same way as rclone.

When using this flag, rclone won't update mtimes of remote files if
they are incorrect as it would normally.

### --stats=TIME ###

Rclone will print stats at regular intervals to show its progress.

This sets the interval.

The default is `1m`. Use 0 to disable.

### --delete-(before,during,after) ###

This option allows you to specify when files on your destination are
deleted when you sync folders.

Specifying the value `--delete-before` will delete all files present on the
destination, but not on the source *before* starting the transfer
of any new or updated files.

Specifying `--delete-during` (default value) will delete files while checking
and uploading files. This is usually the fastest option.

Specifying `--delete-after` will delay deletion of files until all new/updated
files have been successfully transfered.

### --timeout=TIME ###

This sets the IO idle timeout.  If a transfer has started but then
becomes idle for this long it is considered broken and disconnected.

The default is `5m`.  Set to 0 to disable.

### --transfers=N ###

The number of file transfers to run in parallel.  It can sometimes be
useful to set this to a smaller number if the remote is giving a lot
of timeouts or bigger if you have lots of bandwidth and a fast remote.

The default is to run 4 file transfers in parallel.

### -v, --verbose ###

If you set this flag, rclone will become very verbose telling you
about every file it considers and transfers.

Very useful for debugging.

### -V, --version ###

Prints the version number

Developer options
-----------------

These options are useful when developing or debugging rclone.  There
are also some more remote specific options which aren't documented
here which are used for testing.  These start with remote name eg
`--drive-test-option` - see the docs for the remote in question.

### --cpuprofile=FILE ###

Write CPU profile to file.  This can be analysed with `go tool pprof`.

### --dump-bodies ###

Dump HTTP headers and bodies - may contain sensitive info.  Can be
very verbose.  Useful for debugging only.

### --dump-filters ###

Dump the filters to the output.  Useful to see exactly what include
and exclude options are filtering on.

### --dump-headers ###

Dump HTTP headers - may contain sensitive info.  Can be very verbose.
Useful for debugging only.

### --memprofile=FILE ###

Write memory profile to file. This can be analysed with `go tool pprof`.

### --no-check-certificate=true/false ###

`--no-check-certificate` controls whether a client verifies the
server's certificate chain and host name.
If `--no-check-certificate` is true, TLS accepts any certificate
presented by the server and any host name in that certificate.
In this mode, TLS is susceptible to man-in-the-middle attacks.

This option defaults to `false`.

**This should be used only for testing.**

Filtering
---------

For the filtering options

  * `--delete-excluded`
  * `--filter`
  * `--filter-from`
  * `--exclude`
  * `--exclude-from`
  * `--include`
  * `--include-from`
  * `--files-from`
  * `--min-size`
  * `--max-size`
  * `--min-age`
  * `--max-age`
  * `--dump-filters`

See the [filtering section](http://rclone.org/filtering/).

Exit Code
---------

If any errors occurred during the command, rclone will set a non zero
exit code.  This allows scripts to detect when rclone operations have
failed.

# Configuring rclone on a remote / headless machine #

Some of the configurations (those involving oauth2) require an
Internet connected web browser.

If you are trying to set rclone up on a remote or headless box with no
browser available on it (eg a NAS or a server in a datacenter) then
you will need to use an alternative means of configuration.  There are
two ways of doing it, described below.

## Configuring using rclone authorize ##

On the headless box

```
...
Remote config
Use auto config?
 * Say Y if not sure
 * Say N if you are working on a remote or headless machine
y) Yes
n) No
y/n> n
For this to work, you will need rclone available on a machine that has a web browser available.
Execute the following on your machine:
	rclone authorize "amazon cloud drive"
Then paste the result below:
result>
```

Then on your main desktop machine

```
rclone authorize "amazon cloud drive"
If your browser doesn't open automatically go to the following link: http://127.0.0.1:53682/auth
Log in and authorize rclone for access
Waiting for code...
Got code
Paste the following into your remote machine --->
SECRET_TOKEN
<---End paste
```

Then back to the headless box, paste in the code

```
result> SECRET_TOKEN
--------------------
[acd12]
client_id = 
client_secret = 
token = SECRET_TOKEN
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d>
```

## Configuring by copying the config file ##

Rclone stores all of its config in a single configuration file.  This
can easily be copied to configure a remote rclone.

So first configure rclone on your desktop machine

    rclone config

to set up the config file.

Find the config file by running `rclone -h` and looking for the help for the `--config` option

```
$ rclone -h
[snip]
      --config="/home/user/.rclone.conf": Config file.
[snip]
```

Now transfer it to the remote box (scp, cut paste, ftp, sftp etc) and
place it in the correct place (use `rclone -h` on the remote box to
find out where).

# Filtering, includes and excludes #

Rclone has a sophisticated set of include and exclude rules. Some of
these are based on patterns and some on other things like file size.

The filters are applied for the `copy`, `sync`, `move`, `ls`, `lsl`,
`md5sum`, `sha1sum`, `size`, `delete` and `check` operations.
Note that `purge` does not obey the filters.

Each path as it passes through rclone is matched against the include
and exclude rules.  The paths are matched without a leading `/`.

For example the files might be passed to the matching engine like this

  * `file1.jpg`
  * `file2.jpg`
  * `directory/file3.jpg`

## Patterns ##

The patterns used to match files for inclusion or exclusion are based
on "file globs" as used by the unix shell.

If the pattern starts with a `/` then it only matches at the top level
of the directory tree.  If it doesn't start with `/` then it is
matched starting at the end of the path, but it will only match a
complete path element.

    file.jpg  - matches "file.jpg"
              - matches "directory/file.jpg"
              - doesn't match "afile.jpg"
              - doesn't match "directory/afile.jpg"
    /file.jpg - matches "file.jpg"
              - doesn't match "afile.jpg"
              - doesn't match "directory/file.jpg"

A `*` matches anything but not a `/`.

    *.jpg  - matches "file.jpg"
           - matches "directory/file.jpg"
           - doesn't match "file.jpg/something"

Use `**` to match anything, including slashes (`/`).

    dir/** - matches "dir/file.jpg"
           - matches "dir/dir1/dir2/file.jpg"
           - doesn't match "directory/file.jpg"
           - doesn't match "adir/file.jpg"

A `?` matches any character except a slash `/`.

    l?ss  - matches "less"
          - matches "lass"
          - doesn't match "floss"

A `[` and `]` together make a a character class, such as `[a-z]` or
`[aeiou]` or `[[:alpha:]]`.  See the [go regexp
docs](https://golang.org/pkg/regexp/syntax/) for more info on these.

    h[ae]llo - matches "hello"
             - matches "hallo"
             - doesn't match "hullo"

A `{` and `}` define a choice between elements.  It should contain a
comma seperated list of patterns, any of which might match.  These
patterns can contain wildcards.

    {one,two}_potato - matches "one_potato"
                     - matches "two_potato"
                     - doesn't match "three_potato"
                     - doesn't match "_potato"

Special characters can be escaped with a `\` before them.

    \*.jpg       - matches "*.jpg"
    \\.jpg       - matches "\.jpg"
    \[one\].jpg  - matches "[one].jpg"
  
### Differences between rsync and rclone patterns ###

Rclone implements bash style `{a,b,c}` glob matching which rsync doesn't.

Rclone ignores `/` at the end of a pattern.

Rclone always does a wildcard match so `\` must always escape a `\`.

## How the rules are used ##

Rclone maintains a list of include rules and exclude rules.

Each file is matched in order against the list until it finds a match.
The file is then included or excluded according to the rule type.

If the matcher falls off the bottom of the list then the path is
included.

For example given the following rules, `+` being include, `-` being
exclude,

    - secret*.jpg
    + *.jpg
    + *.png
    + file2.avi
    - *

This would include

  * `file1.jpg`
  * `file3.png`
  * `file2.avi`

This would exclude

  * `secret17.jpg`
  * non `*.jpg` and `*.png`

## Adding filtering rules ##

Filtering rules are added with the following command line flags.

### `--exclude` - Exclude files matching pattern ###

Add a single exclude rule with `--exclude`.

Eg `--exclude *.bak` to exclude all bak files from the sync.

### `--exclude-from` - Read exclude patterns from file ###

Add exclude rules from a file.

Prepare a file like this `exclude-file.txt`

    # a sample exclude rule file
    *.bak
    file2.jpg

Then use as `--exclude-from exclude-file.txt`.  This will sync all
files except those ending in `bak` and `file2.jpg`.

This is useful if you have a lot of rules.

### `--include` - Include files matching pattern ###

Add a single include rule with `--include`.

Eg `--include *.{png,jpg}` to include all `png` and `jpg` files in the
backup and no others.

This adds an implicit `--exclude *` at the very end of the filter
list. This means you can mix `--include` and `--include-from` with the
other filters (eg `--exclude`) but you must include all the files you
want in the include statement.  If this doesn't provide enough
flexibility then you must use `--filter-from`.

### `--include-from` - Read include patterns from file ###

Add include rules from a file.

Prepare a file like this `include-file.txt`

    # a sample include rule file
    *.jpg
    *.png
    file2.avi

Then use as `--include-from include-file.txt`.  This will sync all
`jpg`, `png` files and `file2.avi`.

This is useful if you have a lot of rules.

This adds an implicit `--exclude *` at the very end of the filter
list. This means you can mix `--include` and `--include-from` with the
other filters (eg `--exclude`) but you must include all the files you
want in the include statement.  If this doesn't provide enough
flexibility then you must use `--filter-from`.

### `--filter` - Add a file-filtering rule ###

This can be used to add a single include or exclude rule.  Include
rules start with `+ ` and exclude rules start with `- `.  A special
rule called `!` can be used to clear the existing rules.

Eg `--filter "- *.bak"` to exclude all bak files from the sync.

### `--filter-from` - Read filtering patterns from a file ###

Add include/exclude rules from a file.

Prepare a file like this `filter-file.txt`

    # a sample exclude rule file
    - secret*.jpg
    + *.jpg
    + *.png
    + file2.avi
    # exclude everything else
    - *

Then use as `--filter-from filter-file.txt`.  The rules are processed
in the order that they are defined.

This example will include all `jpg` and `png` files, exclude any files
matching `secret*.jpg` and include `file2.avi`.  Everything else will
be excluded from the sync.

### `--files-from` - Read list of source-file names ###

This reads a list of file names from the file passed in and **only**
these files are transferred.  The filtering rules are ignored
completely if you use this option.

Prepare a file like this `files-from.txt`

    # comment
    file1.jpg
    file2.jpg

Then use as `--files-from files-from.txt`.  This will only transfer
`file1.jpg` and `file2.jpg` providing they exist.

### `--min-size` - Don't transfer any file smaller than this ###

This option controls the minimum size file which will be transferred.
This defaults to `kBytes` but a suffix of `k`, `M`, or `G` can be
used.

For example `--min-size 50k` means no files smaller than 50kByte will be
transferred.

### `--max-size` - Don't transfer any file larger than this ###

This option controls the maximum size file which will be transferred.
This defaults to `kBytes` but a suffix of `k`, `M`, or `G` can be
used.

For example `--max-size 1G` means no files larger than 1GByte will be
transferred.

### `--max-age` - Don't transfer any file older than this ###

This option controls the maximum age of files to transfer.  Give in
seconds or with a suffix of:

  * `ms` - Milliseconds
  * `s` - Seconds
  * `m` - Minutes
  * `h` - Hours
  * `d` - Days
  * `w` - Weeks
  * `M` - Months
  * `y` - Years

For example `--max-age 2d` means no files older than 2 days will be
transferred.

### `--min-age` - Don't transfer any file younger than this ###

This option controls the minimum age of files to transfer.  Give in
seconds or with a suffix (see `--max-age` for list of suffixes)

For example `--min-age 2d` means no files younger than 2 days will be
transferred.

### `--delete-excluded` - Delete files on dest excluded from sync ###

**Important** this flag is dangerous - use with `--dry-run` and `-v` first.

When doing `rclone sync` this will delete any files which are excluded
from the sync on the destination.

If for example you did a sync from `A` to `B` without the `--min-size 50k` flag

    rclone sync A: B:

Then you repeated it like this with the `--delete-excluded`

    rclone --min-size 50k --delete-excluded sync A: B:

This would delete all files on `B` which are less than 50 kBytes as
these are now excluded from the sync.

Always test first with `--dry-run` and `-v` before using this flag.

### `--dump-filters` - dump the filters to the output ###

This dumps the defined filters to the output as regular expressions.

Useful for debugging.

## Quoting shell metacharacters ##

The examples above may not work verbatim in your shell as they have
shell metacharacters in them (eg `*`), and may require quoting.

Eg linux, OSX

  * `--include \*.jpg`
  * `--include '*.jpg'`
  * `--include='*.jpg'`

In Windows the expansion is done by the command not the shell so this
should work fine

  * `--include *.jpg`

# Overview of cloud storage systems #

Each cloud storage system is slighly different.  Rclone attempts to
provide a unified interface to them, but some underlying differences
show through.

## Features ##

Here is an overview of the major features of each cloud storage system.

| Name                   | Hash    | ModTime | Case Insensitive | Duplicate Files |
| ---------------------- |:-------:|:-------:|:----------------:|:---------------:|
| Google Drive           | MD5     | Yes     | No               | Yes             |
| Amazon S3              | MD5     | Yes     | No               | No              |
| Openstack Swift        | MD5     | Yes     | No               | No              |
| Dropbox                | -       | No      | Yes              | No              |
| Google Cloud Storage   | MD5     | Yes     | No               | No              |
| Amazon Cloud Drive     | MD5     | No      | Yes              | No              |
| Microsoft One Drive    | SHA1    | Yes     | Yes              | No              |
| Hubic                  | MD5     | Yes     | No               | No              |
| Backblaze B2           | SHA1    | Partial | No               | No              |
| Yandex Disk            | MD5     | Yes     | No               | No              |
| The local filesystem   | All     | Yes     | Depends          | No              |

### Hash ###

The cloud storage system supports various hash types of the objects.  
The hashes are used when transferring data as an integrity check and
can be specifically used with the `--checksum` flag in syncs and in
the `check` command.

To use the checksum checks between filesystems they must support a 
common hash type.

### ModTime ###

The cloud storage system supports setting modification times on
objects.  If it does then this enables a using the modification times
as part of the sync.  If not then only the size will be checked by
default, though the MD5SUM can be checked with the `--checksum` flag.

All cloud storage systems support some kind of date on the object and
these will be set when transferring from the cloud storage system.

Backblaze B2 preserves file modification times on files uploaded and
downloaded, but doesn't use them to decide which objects to sync.

### Case Insensitive ###

If a cloud storage systems is case sensitive then it is possible to
have two files which differ only in case, eg `file.txt` and
`FILE.txt`.  If a cloud storage system is case insensitive then that
isn't possible.

This can cause problems when syncing between a case insensitive
system and a case sensitive system.  The symptom of this is that no
matter how many times you run the sync it never completes fully.

The local filesystem may or may not be case sensitive depending on OS.

  * Windows - usually case insensitive, though case is preserved
  * OSX - usually case insensitive, though it is possible to format case sensitive
  * Linux - usually case sensitive, but there are case insensitive file systems (eg FAT formatted USB keys)

Most of the time this doesn't cause any problems as people tend to
avoid files whose name differs only by case even on case sensitive
systems.

### Duplicate files ###

If a cloud storage system allows duplicate files then it can have two
objects with the same name.

This confuses rclone greatly when syncing.

Google Drive
-----------------------------------------

Paths are specified as `drive:path`

Drive paths may be as deep as required, eg `drive:directory/subdirectory`.

The initial setup for drive involves getting a token from Google drive
which you need to do in your browser.  `rclone config` walks you
through it.

Here is an example of how to make a remote called `remote`.  First run:

     rclone config

This will guide you through an interactive setup process:

```
n) New remote
d) Delete remote
q) Quit config
e/n/d/q> n
name> remote
What type of source is it?
Choose a number from below
 1) swift
 2) s3
 3) local
 4) drive
type> 4
Google Application Client Id - leave blank normally.
client_id> 
Google Application Client Secret - leave blank normally.
client_secret> 
Remote config
Use auto config?
 * Say Y if not sure
 * Say N if you are working on a remote or headless machine or Y didn't work
y) Yes
n) No
y/n> y
If your browser doesn't open automatically go to the following link: http://127.0.0.1:53682/auth
Log in and authorize rclone for access
Waiting for code...
Got code
--------------------
[remote]
client_id = 
client_secret = 
token = {"AccessToken":"xxxx.x.xxxxx_xxxxxxxxxxx_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx","RefreshToken":"1/xxxxxxxxxxxxxxxx_xxxxxxxxxxxxxxxxxxxxxxxxxx","Expiry":"2014-03-16T13:57:58.955387075Z","Extra":null}
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

Note that rclone runs a webserver on your local machine to collect the
token as returned from Google if you use auto config mode. This only
runs from the moment it opens your browser to the moment you get back
the verification code.  This is on `http://127.0.0.1:53682/` and this
it may require you to unblock it temporarily if you are running a host
firewall, or use manual mode.

You can then use it like this,

List directories in top level of your drive

    rclone lsd remote:

List all the files in your drive

    rclone ls remote:

To copy a local directory to a drive directory called backup

    rclone copy /home/source remote:backup

### Modified time ###

Google drive stores modification times accurate to 1 ms.

### Revisions ###

Google drive stores revisions of files.  When you upload a change to
an existing file to google drive using rclone it will create a new
revision of that file.

Revisions follow the standard google policy which at time of writing
was

  * They are deleted after 30 days or 100 revisions (whatever comes first).
  * They do not count towards a user storage quota.

### Deleting files ###

By default rclone will delete files permanently when requested.  If
sending them to the trash is required instead then use the
`--drive-use-trash` flag.

### Specific options ###

Here are the command line options specific to this cloud storage
system.

#### --drive-chunk-size=SIZE ####

Upload chunk size. Must a power of 2 >= 256k. Default value is 256kB.

#### --drive-full-list ####

Use a full listing for directory list. More data but usually
quicker. On by default, disable with `--full-drive-list=false`.

#### --drive-upload-cutoff=SIZE ####

File size cutoff for switching to chunked upload.  Default is 256kB.

#### --drive-use-trash ####

Send files to the trash instead of deleting permanently. Defaults to
off, namely deleting files permanently.

#### --drive-auth-owner-only ####

Only consider files owned by the authenticated user. Requires
that --drive-full-list=true (default).

#### --drive-formats ####

Google documents can only be exported from Google drive.  When rclone
downloads a Google doc it chooses a format to download depending upon
this setting.

By default the formats are `docx,xlsx,pptx,svg` which are a sensible
default for an editable document.

When choosing a format, rclone runs down the list provided in order
and chooses the first file format the doc can be exported as from the
list. If the file can't be exported to a format on the formats list,
then rclone will choose a format from the default list.

If you prefer an archive copy then you might use `--drive-formats
pdf`, or if you prefer openoffice/libreoffice formats you might use
`--drive-formats ods,odt`.

Note that rclone adds the extension to the google doc, so if it is
calles `My Spreadsheet` on google docs, it will be exported as `My
Spreadsheet.xlsx` or `My Spreadsheet.pdf` etc.

Here are the possible extensions with their corresponding mime types.

| Extension | Mime Type | Description |
| --------- |-----------| ------------|
| csv  | text/csv | Standard CSV format for Spreadsheets |
| doc  | application/msword | Micosoft Office Document |
| docx | application/vnd.openxmlformats-officedocument.wordprocessingml.document | Microsoft Office Document |
| html | text/html | An HTML Document |
| jpg  | image/jpeg | A JPEG Image File |
| ods  | application/vnd.oasis.opendocument.spreadsheet | Openoffice Spreadsheet |
| ods  | application/x-vnd.oasis.opendocument.spreadsheet | Openoffice Spreadsheet |
| odt  | application/vnd.oasis.opendocument.text | Openoffice Document |
| pdf  | application/pdf | Adobe PDF Format |
| png  | image/png | PNG Image Format|
| pptx | application/vnd.openxmlformats-officedocument.presentationml.presentation | Microsoft Office Powerpoint |
| rtf  | application/rtf | Rich Text Format |
| svg  | image/svg+xml | Scalable Vector Graphics Format |
| txt  | text/plain | Plain Text |
| xls  | application/vnd.ms-excel | Microsoft Office Spreadsheet |
| xlsx | application/vnd.openxmlformats-officedocument.spreadsheetml.sheet | Microsoft Office Spreadsheet |
| zip  | application/zip | A ZIP file of HTML, Images CSS |

### Limitations ###

Drive has quite a lot of rate limiting.  This causes rclone to be
limited to transferring about 2 files per second only.  Individual
files may be transferred much faster at 100s of MBytes/s but lots of
small files can take a long time.

Amazon S3
---------------------------------------

Paths are specified as `remote:bucket` (or `remote:` for the `lsd`
command.)  You may put subdirectories in too, eg `remote:bucket/path/to/dir`.

Here is an example of making an s3 configuration.  First run

    rclone config

This will guide you through an interactive setup process.

```
No remotes found - make a new one
n) New remote
q) Quit config
n/q> n
name> remote
What type of source is it?
Choose a number from below
 1) swift
 2) s3
 3) local
 4) google cloud storage
 5) dropbox
 6) drive
type> 2
AWS Access Key ID.
access_key_id> accesskey
AWS Secret Access Key (password). 
secret_access_key> secretaccesskey
Region to connect to.
Choose a number from below, or type in your own value
 * The default endpoint - a good choice if you are unsure.
 * US Region, Northern Virginia or Pacific Northwest.
 * Leave location constraint empty.
 1) us-east-1
 * US West (Oregon) Region
 * Needs location constraint us-west-2.
 2) us-west-2
[snip]
 * South America (Sao Paulo) Region
 * Needs location constraint sa-east-1.
 9) sa-east-1
 * If using an S3 clone that only understands v2 signatures - eg Ceph - set this and make sure you set the endpoint.
10) other-v2-signature
 * If using an S3 clone that understands v4 signatures set this and make sure you set the endpoint.
11) other-v4-signature
region> 1
Endpoint for S3 API.
Leave blank if using AWS to use the default endpoint for the region.
Specify if using an S3 clone such as Ceph.
endpoint> 
Location constraint - must be set to match the Region. Used when creating buckets only.
Choose a number from below, or type in your own value
 * Empty for US Region, Northern Virginia or Pacific Northwest.
 1) 
 * US West (Oregon) Region.
 2) us-west-2
 * US West (Northern California) Region.
 3) us-west-1
 * EU (Ireland) Region.
 4) eu-west-1
[snip]
location_constraint> 1
Remote config
--------------------
[remote]
access_key_id = accesskey
secret_access_key = secretaccesskey
region = us-east-1
endpoint = 
location_constraint = 
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
Current remotes:

Name                 Type
====                 ====
remote               s3

e) Edit existing remote
n) New remote
d) Delete remote
q) Quit config
e/n/d/q> q
```

This remote is called `remote` and can now be used like this

See all buckets

    rclone lsd remote:

Make a new bucket

    rclone mkdir remote:bucket

List the contents of a bucket

    rclone ls remote:bucket

Sync `/home/local/directory` to the remote bucket, deleting any excess
files in the bucket.

    rclone sync /home/local/directory remote:bucket

### Modified time ###

The modified time is stored as metadata on the object as
`X-Amz-Meta-Mtime` as floating point since the epoch accurate to 1 ns.

### Multipart uploads ###

rclone supports multipart uploads with S3 which means that it can
upload files bigger than 5GB. Note that files uploaded with multipart
upload don't have an MD5SUM.

### Buckets and Regions ###

With Amazon S3 you can list buckets (`rclone lsd`) using any region,
but you can only access the content of a bucket from the region it was
created in.  If you attempt to access a bucket from the wrong region,
you will get an error, `incorrect region, the bucket is not in 'XXX'
region`.

### Anonymous access to public buckets ###

If you want to use rclone to access a public bucket, configure with a
blank `access_key_id` and `secret_access_key`.  Eg

```
e) Edit existing remote
n) New remote
d) Delete remote
q) Quit config
e/n/d/q> n
name> anons3
What type of source is it?
Choose a number from below
 1) amazon cloud drive
 2) drive
 3) dropbox
 4) google cloud storage
 5) local
 6) s3
 7) swift
type> 6
AWS Access Key ID - leave blank for anonymous access.
access_key_id> 
AWS Secret Access Key (password) - leave blank for anonymous access.
secret_access_key> 
Region to connect to.
region> 1
endpoint> 
location_constraint> 
```

Then use it as normal with the name of the public bucket, eg

    rclone lsd anons3:1000genomes

You will be able to list and copy data but not upload it.

### Ceph ###

Ceph is an object storage system which presents an Amazon S3 interface.

To use rclone with ceph, you need to set the following parameters in
the config.

```
access_key_id = Whatever
secret_access_key = Whatever
endpoint = https://ceph.endpoint.goes.here/
region = other-v2-signature
```

Note also that Ceph sometimes puts `/` in the passwords it gives
users.  If you read the secret access key using the command line tools
you will get a JSON blob with the `/` escaped as `\/`.  Make sure you
only write `/` in the secret access key.

Eg the dump from Ceph looks something like this (irrelevant keys
removed).

```
{
    "user_id": "xxx",
    "display_name": "xxxx",
    "keys": [
        {
            "user": "xxx",
            "access_key": "xxxxxx",
            "secret_key": "xxxxxx\/xxxx"
        }
    ],
}
```

Because this is a json dump, it is encoding the `/` as `\/`, so if you
use the secret key as `xxxxxx/xxxx`  it will work fine.

Swift
----------------------------------------

Swift refers to [Openstack Object Storage](http://www.openstack.org/software/openstack-storage/).
Commercial implementations of that being:

  * [Rackspace Cloud Files](http://www.rackspace.com/cloud/files/)
  * [Memset Memstore](http://www.memset.com/cloud/storage/)

Paths are specified as `remote:container` (or `remote:` for the `lsd`
command.)  You may put subdirectories in too, eg `remote:container/path/to/dir`.

Here is an example of making a swift configuration.  First run

    rclone config

This will guide you through an interactive setup process.

```
No remotes found - make a new one
n) New remote
q) Quit config
n/q> n
name> remote
What type of source is it?
Choose a number from below
 1) swift
 2) s3
 3) local
 4) drive
type> 1
User name to log in.
user> user_name
API key or password.
key> password_or_api_key
Authentication URL for server.
Choose a number from below, or type in your own value
 * Rackspace US
 1) https://auth.api.rackspacecloud.com/v1.0
 * Rackspace UK
 2) https://lon.auth.api.rackspacecloud.com/v1.0
 * Rackspace v2
 3) https://identity.api.rackspacecloud.com/v2.0
 * Memset Memstore UK
 4) https://auth.storage.memset.com/v1.0
 * Memset Memstore UK v2
 5) https://auth.storage.memset.com/v2.0
 * OVH
 6) https://auth.cloud.ovh.net/v2.0
auth> 1
Tenant name - optional
tenant>
Remote config
--------------------
[remote]
user = user_name
key = password_or_api_key
auth = https://auth.api.rackspacecloud.com/v1.0
tenant =
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

This remote is called `remote` and can now be used like this

See all containers

    rclone lsd remote:

Make a new container

    rclone mkdir remote:container

List the contents of a container

    rclone ls remote:container

Sync `/home/local/directory` to the remote container, deleting any
excess files in the container.

    rclone sync /home/local/directory remote:container

### Specific options ###

Here are the command line options specific to this cloud storage
system.

#### --swift-chunk-size=SIZE ####

Above this size files will be chunked into a _segments container.  The
default for this is 5GB which is its maximum value.
      
### Modified time ###

The modified time is stored as metadata on the object as
`X-Object-Meta-Mtime` as floating point since the epoch accurate to 1
ns.

This is a defacto standard (used in the official python-swiftclient
amongst others) for storing the modification time for an object.

Dropbox
---------------------------------

Paths are specified as `remote:path`

Dropbox paths may be as deep as required, eg
`remote:directory/subdirectory`.

The initial setup for dropbox involves getting a token from Dropbox
which you need to do in your browser.  `rclone config` walks you
through it.

Here is an example of how to make a remote called `remote`.  First run:

     rclone config

This will guide you through an interactive setup process:

```
n) New remote
d) Delete remote
q) Quit config
e/n/d/q> n
name> remote
What type of source is it?
Choose a number from below
 1) swift
 2) s3
 3) local
 4) google cloud storage
 5) dropbox
 6) drive
type> 5
Dropbox App Key - leave blank normally.
app_key> 
Dropbox App Secret - leave blank normally.
app_secret> 
Remote config
Please visit:
https://www.dropbox.com/1/oauth2/authorize?client_id=XXXXXXXXXXXXXXX&response_type=code
Enter the code: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX_XXXXXXXXXX
--------------------
[remote]
app_key = 
app_secret = 
token = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX_XXXX_XXXXXXXXXXXXXXXXXXXXXXXXXXXXX
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

You can then use it like this,

List directories in top level of your dropbox

    rclone lsd remote:

List all the files in your dropbox

    rclone ls remote:

To copy a local directory to a dropbox directory called backup

    rclone copy /home/source remote:backup

### Modified time and MD5SUMs ###

Dropbox doesn't have the capability of storing modification times or
MD5SUMs so syncs will effectively have the `--size-only` flag set.

### Specific options ###

Here are the command line options specific to this cloud storage
system.

#### --dropbox-chunk-size=SIZE ####

Upload chunk size. Max 150M. The default is 128MB.  Note that this
isn't buffered into memory.

### Limitations ###

Note that Dropbox is case insensitive so you can't have a file called
"Hello.doc" and one called "hello.doc".

There are some file names such as `thumbs.db` which Dropbox can't
store.  There is a full list of them in the ["Ignored Files" section
of this document](https://www.dropbox.com/en/help/145).  Rclone will
issue an error message `File name disallowed - not uploading` if it
attempt to upload one of those file names, but the sync won't fail.

Google Cloud Storage
-------------------------------------------------

Paths are specified as `remote:bucket` (or `remote:` for the `lsd`
command.)  You may put subdirectories in too, eg `remote:bucket/path/to/dir`.

The initial setup for google cloud storage involves getting a token from Google Cloud Storage
which you need to do in your browser.  `rclone config` walks you
through it.

Here is an example of how to make a remote called `remote`.  First run:

     rclone config

This will guide you through an interactive setup process:

```
n) New remote
d) Delete remote
q) Quit config
e/n/d/q> n
name> remote
What type of source is it?
Choose a number from below
 1) swift
 2) s3
 3) local
 4) google cloud storage
 5) dropbox
 6) drive
type> 4
Google Application Client Id - leave blank normally.
client_id> 
Google Application Client Secret - leave blank normally.
client_secret> 
Project number optional - needed only for list/create/delete buckets - see your developer console.
project_number> 12345678
Access Control List for new objects.
Choose a number from below, or type in your own value
 * Object owner gets OWNER access, and all Authenticated Users get READER access.
 1) authenticatedRead
 * Object owner gets OWNER access, and project team owners get OWNER access.
 2) bucketOwnerFullControl
 * Object owner gets OWNER access, and project team owners get READER access.
 3) bucketOwnerRead
 * Object owner gets OWNER access [default if left blank].
 4) private
 * Object owner gets OWNER access, and project team members get access according to their roles.
 5) projectPrivate
 * Object owner gets OWNER access, and all Users get READER access.
 6) publicRead
object_acl> 4
Access Control List for new buckets.
Choose a number from below, or type in your own value
 * Project team owners get OWNER access, and all Authenticated Users get READER access.
 1) authenticatedRead
 * Project team owners get OWNER access [default if left blank].
 2) private
 * Project team members get access according to their roles.
 3) projectPrivate
 * Project team owners get OWNER access, and all Users get READER access.
 4) publicRead
 * Project team owners get OWNER access, and all Users get WRITER access.
 5) publicReadWrite
bucket_acl> 2
Remote config
Remote config
Use auto config?
 * Say Y if not sure
 * Say N if you are working on a remote or headless machine or Y didn't work
y) Yes
n) No
y/n> y
If your browser doesn't open automatically go to the following link: http://127.0.0.1:53682/auth
Log in and authorize rclone for access
Waiting for code...
Got code
--------------------
[remote]
type = google cloud storage
client_id = 
client_secret = 
token = {"AccessToken":"xxxx.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx","RefreshToken":"x/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx_xxxxxxxxx","Expiry":"2014-07-17T20:49:14.929208288+01:00","Extra":null}
project_number = 12345678
object_acl = private
bucket_acl = private
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

Note that rclone runs a webserver on your local machine to collect the
token as returned from Google if you use auto config mode. This only
runs from the moment it opens your browser to the moment you get back
the verification code.  This is on `http://127.0.0.1:53682/` and this
it may require you to unblock it temporarily if you are running a host
firewall, or use manual mode.

This remote is called `remote` and can now be used like this

See all the buckets in your project

    rclone lsd remote:

Make a new bucket

    rclone mkdir remote:bucket

List the contents of a bucket

    rclone ls remote:bucket

Sync `/home/local/directory` to the remote bucket, deleting any excess
files in the bucket.

    rclone sync /home/local/directory remote:bucket

### Modified time ###

Google google cloud storage stores md5sums natively and rclone stores
modification times as metadata on the object, under the "mtime" key in
RFC3339 format accurate to 1ns.

Amazon Cloud Drive
-----------------------------------------

Paths are specified as `remote:path`

Paths may be as deep as required, eg `remote:directory/subdirectory`.

The initial setup for Amazon cloud drive involves getting a token from
Amazon which you need to do in your browser.  `rclone config` walks
you through it.

Here is an example of how to make a remote called `remote`.  First run:

     rclone config

This will guide you through an interactive setup process:

```
n) New remote
d) Delete remote
q) Quit config
e/n/d/q> n
name> remote
What type of source is it?
Choose a number from below
 1) amazon cloud drive
 2) drive
 3) dropbox
 4) google cloud storage
 5) local
 6) s3
 7) swift
type> 1
Amazon Application Client Id - leave blank normally.
client_id> 
Amazon Application Client Secret - leave blank normally.
client_secret> 
Remote config
If your browser doesn't open automatically go to the following link: http://127.0.0.1:53682/auth
Log in and authorize rclone for access
Waiting for code...
Got code
--------------------
[remote]
client_id = 
client_secret = 
token = {"access_token":"xxxxxxxxxxxxxxxxxxxxxxx","token_type":"bearer","refresh_token":"xxxxxxxxxxxxxxxxxx","expiry":"2015-09-06T16:07:39.658438471+01:00"}
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

See the [remote setup docs](http://rclone.org/remote_setup/) for how to set it up on a
machine with no Internet browser available.

Note that rclone runs a webserver on your local machine to collect the
token as returned from Amazon. This only runs from the moment it
opens your browser to the moment you get back the verification
code.  This is on `http://127.0.0.1:53682/` and this it may require
you to unblock it temporarily if you are running a host firewall.

Once configured you can then use `rclone` like this,

List directories in top level of your Amazon cloud drive

    rclone lsd remote:

List all the files in your Amazon cloud drive

    rclone ls remote:

To copy a local directory to an Amazon cloud drive directory called backup

    rclone copy /home/source remote:backup

### Modified time and MD5SUMs ###

Amazon cloud drive doesn't allow modification times to be changed via
the API so these won't be accurate or used for syncing.

It does store MD5SUMs so for a more accurate sync, you can use the
`--checksum` flag.

### Deleting files ###

Any files you delete with rclone will end up in the trash.  Amazon
don't provide an API to permanently delete files, nor to empty the
trash, so you will have to do that with one of Amazon's apps or via
the Amazon cloud drive website.

### Specific options ###

Here are the command line options specific to this cloud storage
system.

#### --acd-templink-threshold=SIZE ####

Files this size or more will be downloaded via their `tempLink`. This
is to work around a problem with Amazon Cloud Drive which blocks
downloads of files bigger than about 10GB.  The default for this is
9GB which shouldn't need to be changed.

To download files above this threshold, rclone requests a `tempLink`
which downloads the file through a temporary URL directly from the
underlying S3 storage.

### Limitations ###

Note that Amazon cloud drive is case insensitive so you can't have a
file called "Hello.doc" and one called "hello.doc".

Amazon cloud drive has rate limiting so you may notice errors in the
sync (429 errors).  rclone will automatically retry the sync up to 3
times by default (see `--retries` flag) which should hopefully work
around this problem.

Amazon cloud drive has an internal limit of file sizes that can be
uploaded to the service. This limit is not officially published,
but all files larger than this will fail.

At the time of writing (Jan 2016) is in the area of 50GB per file.
This means that larger files are likely to fail.

Unfortunatly there is no way for rclone to see that this failure is 
because of file size, so it will retry the operation, as any other
failure. To avoid this problem, use `--max-size=50GB` option to limit
the maximum size of uploaded files.

Microsoft One Drive
-----------------------------------------

Paths are specified as `remote:path`

Paths may be as deep as required, eg `remote:directory/subdirectory`.

The initial setup for One Drive involves getting a token from
Microsoft which you need to do in your browser.  `rclone config` walks
you through it.

Here is an example of how to make a remote called `remote`.  First run:

     rclone config

This will guide you through an interactive setup process:

```
n) New remote
d) Delete remote
q) Quit config
e/n/d/q> n
name> remote
What type of source is it?
Choose a number from below
 1) amazon cloud drive
 2) drive
 3) dropbox
 4) google cloud storage
 5) local
 6) onedrive
 7) s3
 8) swift
type> 6
Microsoft App Client Id - leave blank normally.
client_id> 
Microsoft App Client Secret - leave blank normally.
client_secret> 
Remote config
If your browser doesn't open automatically go to the following link: http://127.0.0.1:53682/auth
Log in and authorize rclone for access
Waiting for code...
Got code
--------------------
[remote]
client_id = 
client_secret = 
token = {"access_token":"XXXXXX"}
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

See the [remote setup docs](http://rclone.org/remote_setup/) for how to set it up on a
machine with no Internet browser available.

Note that rclone runs a webserver on your local machine to collect the
token as returned from Microsoft. This only runs from the moment it
opens your browser to the moment you get back the verification
code.  This is on `http://127.0.0.1:53682/` and this it may require
you to unblock it temporarily if you are running a host firewall.

Once configured you can then use `rclone` like this,

List directories in top level of your One Drive

    rclone lsd remote:

List all the files in your One Drive

    rclone ls remote:

To copy a local directory to an One Drive directory called backup

    rclone copy /home/source remote:backup

### Modified time and hashes ###

One Drive allows modification times to be set on objects accurate to 1
second.  These will be used to detect whether objects need syncing or
not.

One drive supports SHA1 type hashes, so you can use `--checksum` flag.


### Deleting files ###

Any files you delete with rclone will end up in the trash.  Microsoft
doesn't provide an API to permanently delete files, nor to empty the
trash, so you will have to do that with one of Microsoft's apps or via
the One Drive website.

### Specific options ###

Here are the command line options specific to this cloud storage
system.

#### --onedrive-chunk-size=SIZE ####

Above this size files will be chunked - must be multiple of 320k. The
default is 10MB.  Note that the chunks will be buffered into memory.

#### --onedrive-upload-cutoff=SIZE ####

Cutoff for switching to chunked upload - must be <= 100MB. The default
is 10MB.

### Limitations ###

Note that One Drive is case insensitive so you can't have a
file called "Hello.doc" and one called "hello.doc".

Rclone only supports your default One Drive, and doesn't work with One
Drive for business.  Both these issues may be fixed at some point
depending on user demand!

There are quite a few characters that can't be in One Drive file
names.  These can't occur on Windows platforms, but on non-Windows
platforms they are common.  Rclone will map these names to and from an
identical looking unicode equivalent.  For example if a file has a `?`
in it will be mapped to `？` instead.

Hubic
-----------------------------------------

Paths are specified as `remote:path`

Paths are specified as `remote:container` (or `remote:` for the `lsd`
command.)  You may put subdirectories in too, eg `remote:container/path/to/dir`.

The initial setup for Hubic involves getting a token from Hubic which
you need to do in your browser.  `rclone config` walks you through it.

Here is an example of how to make a remote called `remote`.  First run:

     rclone config

This will guide you through an interactive setup process:

```
n) New remote
d) Delete remote
q) Quit config
e/n/d/q> n
name> remote
What type of source is it?
Choose a number from below
 1) amazon cloud drive
 2) drive
 3) dropbox
 4) google cloud storage
 5) local
 6) onedrive
 7) hubic
 8) s3
 9) swift
type> 7
Hubic App Client Id - leave blank normally.
client_id> 
Hubic App Client Secret - leave blank normally.
client_secret> 
Remote config
If your browser doesn't open automatically go to the following link: http://localhost:53682/auth
Log in and authorize rclone for access
Waiting for code...
Got code
--------------------
[remote]
client_id = 
client_secret = 
token = {"access_token":"XXXXXX"}
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

See the [remote setup docs](http://rclone.org/remote_setup/) for how to set it up on a
machine with no Internet browser available.

Note that rclone runs a webserver on your local machine to collect the
token as returned from Hubic. This only runs from the moment it opens
your browser to the moment you get back the verification code.  This
is on `http://127.0.0.1:53682/` and this it may require you to unblock
it temporarily if you are running a host firewall.

Once configured you can then use `rclone` like this,

List containers in the top level of your Hubic

    rclone lsd remote:

List all the files in your Hubic

    rclone ls remote:

To copy a local directory to an Hubic directory called backup

    rclone copy /home/source remote:backup

### Modified time ###

The modified time is stored as metadata on the object as
`X-Object-Meta-Mtime` as floating point since the epoch accurate to 1
ns.

This is a defacto standard (used in the official python-swiftclient
amongst others) for storing the modification time for an object.

Note that Hubic wraps the Swift backend, so most of the properties of
are the same.

### Limitations ###

Code to refresh the OpenStack token isn't done yet which may cause
problems with very long transfers.

Backblaze B2
----------------------------------------

B2 is [Backblaze's cloud storage system](https://www.backblaze.com/b2/).

Paths are specified as `remote:bucket` (or `remote:` for the `lsd`
command.)  You may put subdirectories in too, eg `remote:bucket/path/to/dir`.

Here is an example of making a b2 configuration.  First run

    rclone config

This will guide you through an interactive setup process.  You will
need your account number (a short hex number) and key (a long hex
number) which you can get from the b2 control panel.

```
No remotes found - make a new one
n) New remote
q) Quit config
n/q> n
name> remote
What type of source is it?
Choose a number from below
 1) amazon cloud drive
 2) b2
 3) drive
 4) dropbox
 5) google cloud storage
 6) swift
 7) hubic
 8) local
 9) onedrive
10) s3
type> 2
Account ID
account> 123456789abc
Application Key
key> 0123456789abcdef0123456789abcdef0123456789
Endpoint for the service - leave blank normally.
endpoint> 
Remote config
--------------------
[remote]
account = 123456789abc
key = 0123456789abcdef0123456789abcdef0123456789
endpoint = 
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

This remote is called `remote` and can now be used like this

See all buckets

    rclone lsd remote:

Make a new bucket

    rclone mkdir remote:bucket

List the contents of a bucket

    rclone ls remote:bucket

Sync `/home/local/directory` to the remote bucket, deleting any
excess files in the bucket.

    rclone sync /home/local/directory remote:bucket

### Modified time ###

The modified time is stored as metadata on the object as
`X-Bz-Info-src_last_modified_millis` as milliseconds since 1970-01-01
in the Backblaze standard.  Other tools should be able to use this as
a modified time.

Modified times are set on upload, read on download and shown in
listings.  They are not used in syncing as unfortunately B2 doesn't
have an API method to set them independently of doing an upload.

### SHA1 checksums ###

The SHA1 checksums of the files are checked on upload and download and
will be used in the syncing process. You can use the `--checksum` flag.

### Versions ###

When rclone uploads a new version of a file it creates a [new version
of it](https://www.backblaze.com/b2/docs/file_versions.html).
Likewise when you delete a file, the old version will still be
available.

The old versions of files are visible in the B2 web interface, but not
via rclone yet.

Rclone doesn't provide any way of managing old versions (downloading
them or deleting them) at the moment.  When you `purge` a bucket, all
the old versions will be deleted.

### Bugs ###

Note that when uploading a file, rclone has to make a temporary copy
of it in your temp filing system.  This is due to a weakness in the B2
API which I'm hoping will be addressed soon.

### API ###

Here are [some notes I made on the backblaze
API](https://gist.github.com/ncw/166dabf352b399f1cc1c) while
integrating it with rclone which detail the changes I'd like to see.
With a couple of small tweaks Backblaze could enable rclone to not
make a temporary copy of all files and fully support modification
times.

Yandex Disk
----------------------------------------

[Yandex Disk](https://disk.yandex.com) is a cloud storage solution created by [Yandex](http://yandex.com).

Yandex paths may be as deep as required, eg `remote:directory/subdirectory`.

Here is an example of making a yandex configuration.  First run

    rclone config

This will guide you through an interactive setup process:

```
No remotes found - make a new one
n) New remote
q) Quit config
n/q> n
name> remote
What type of source is it?
Choose a number from below
 1) amazon cloud drive
 2) b2
 3) drive
 4) dropbox
 5) google cloud storage
 6) swift
 7) hubic
 8) local
 9) onedrive
10) s3
11) yandex
type> 11
Yandex Client Id - leave blank normally.
client_id> 
Yandex Client Secret - leave blank normally.
client_secret> 
Remote config
If your browser doesn't open automatically go to the following link: http://127.0.0.1:53682/auth
Log in and authorize rclone for access
Waiting for code...
Got code
--------------------
[remote]
client_id = 
client_secret = 
token = {"access_token":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx","token_type":"bearer","expiry":"2016-12-29T12:27:11.362788025Z"}
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

See the [remote setup docs](http://rclone.org/remote_setup/) for how to set it up on a
machine with no Internet browser available.

Note that rclone runs a webserver on your local machine to collect the
token as returned from Yandex Disk. This only runs from the moment it
opens your browser to the moment you get back the verification code.
This is on `http://127.0.0.1:53682/` and this it may require you to
unblock it temporarily if you are running a host firewall.

Once configured you can then use `rclone` like this,

See top level directories

    rclone lsd remote:

Make a new directory

    rclone mkdir remote:directory

List the contents of a directory

    rclone ls remote:directory

Sync `/home/local/directory` to the remote path, deleting any
excess files in the path.

    rclone sync /home/local/directory remote:directory

### Modified time ###

Modified times are supported and are stored accurate to 1 ns in custom
metadata called `rclone_modified` in RFC3339 with nanoseconds format.

### MD5 checksums ###

MD5 checksums are natively supported by Yandex Disk.

Local Filesystem
-------------------------------------------

Local paths are specified as normal filesystem paths, eg `/path/to/wherever`, so

    rclone sync /home/source /tmp/destination

Will sync `/home/source` to `/tmp/destination`

These can be configured into the config file for consistencies sake,
but it is probably easier not to.

### Modified time ###

Rclone reads and writes the modified time using an accuracy determined by
the OS.  Typically this is 1ns on Linux, 10 ns on Windows and 1 Second
on OS X.

### Filenames ###

Filenames are expected to be encoded in UTF-8 on disk.  This is the
normal case for Windows and OS X.  There is a bit more uncertainty in
the Linux world, but new distributions will have UTF-8 encoded files
names.

If an invalid (non-UTF8) filename is read, the invalid caracters will
be replaced with the unicode replacement character, '�'.  `rclone`
will emit a debug message in this case (use `-v` to see), eg

```
Local file system at .: Replacing invalid UTF-8 characters in "gro\xdf"
```

### Long paths on Windows ###

Rclone handles long paths automatically, by converting all paths to long
[UNC paths](https://msdn.microsoft.com/en-us/library/windows/desktop/aa365247(v=vs.85).aspx#maxpath)
which allows paths up to 32,767 characters.

This is why you will see that your paths, for instance `c:\files` is
converted to the UNC path `\\?\c:\files` in the output,
and `\\server\share` is converted to `\\?\UNC\server\share`.

However, in rare cases this may cause problems with buggy file
system drivers like [EncFS](https://github.com/ncw/rclone/issues/261).
To disable UNC conversion globally, add this to your `.rclone.conf` file:

```
[local]
nounc = true
```

If you want to selectively disable UNC, you can add it to a separate entry like this:

```
[nounc]
type = local
nounc = true
```
And use rclone like this:

`rclone copy c:\src nounc:z:\dst`

This will use UNC paths on `c:\src` but not on `z:\dst`.
Of course this will cause problems if the absolute path length of a
file exceeds 258 characters on z, so only use this option if you have to.

Changelog
---------

  * v1.27 - 2016-01-31
    * New Features
      * Easier headless configuration with `rclone authorize`
      * Add support for multiple hash types - we now check SHA1 as well as MD5 hashes.
      * `delete` command which does obey the filters (unlike `purge`)
      * `dedupe` command to deduplicate a remote.  Useful with Google Drive.
      * Add `--ignore-existing` flag to skip all files that exist on destination.
      * Add `--delete-before`, `--delete-during`, `--delete-after` flags.
      * Add `--memprofile` flag to debug memory use.
      * Warn the user about files with same name but different case
      * Make `--include` rules add their implict exclude * at the end of the filter list
      * Deprecate compiling with go1.3
    * Amazon Cloud Drive
      * Fix download of files > 10 GB
      * Fix directory traversal ("Next token is expired") for large directory listings
      * Remove 409 conflict from error codes we will retry - stops very long pauses
    * Backblaze B2
      * SHA1 hashes now checked by rclone core
    * Drive
      * Add `--drive-auth-owner-only` to only consider files owned by the user - thanks Björn Harrtell
      * Export Google documents
    * Dropbox
      * Make file exclusion error controllable with -q
    * Swift
      * Fix upload from unprivileged user.
    * S3
      * Fix updating of mod times of files with `+` in.
    * Local
      * Add local file system option to disable UNC on Windows.
  * v1.26 - 2016-01-02
    * New Features
      * Yandex storage backend - thank you Dmitry Burdeev ("dibu")
      * Implement Backblaze B2 storage backend
      * Add --min-age and --max-age flags - thank you Adriano Aurélio Meirelles
      * Make ls/lsl/md5sum/size/check obey includes and excludes
    * Fixes
      * Fix crash in http logging
      * Upload releases to github too
    * Swift
      * Fix sync for chunked files
    * One Drive
      * Re-enable server side copy
      * Don't mask HTTP error codes with JSON decode error
    * S3
      * Fix corrupting Content-Type on mod time update (thanks Joseph Spurrier)
  * v1.25 - 2015-11-14
    * New features
      * Implement Hubic storage system
    * Fixes
      * Fix deletion of some excluded files without --delete-excluded
        * This could have deleted files unexpectedly on sync
        * Always check first with `--dry-run`!
    * Swift
      * Stop SetModTime losing metadata (eg X-Object-Manifest)
        * This could have caused data loss for files > 5GB in size
      * Use ContentType from Object to avoid lookups in listings
    * One Drive
      * disable server side copy as it seems to be broken at Microsoft
  * v1.24 - 2015-11-07
    * New features
      * Add support for Microsoft One Drive
      * Add `--no-check-certificate` option to disable server certificate verification
      * Add async readahead buffer for faster transfer of big files
    * Fixes
      * Allow spaces in remotes and check remote names for validity at creation time
      * Allow '&' and disallow ':' in Windows filenames.
    * Swift
      * Ignore directory marker objects where appropriate - allows working with Hubic
      * Don't delete the container if fs wasn't at root
    * S3
      * Don't delete the bucket if fs wasn't at root
    * Google Cloud Storage
      * Don't delete the bucket if fs wasn't at root
  * v1.23 - 2015-10-03
    * New features
      * Implement `rclone size` for measuring remotes
    * Fixes
      * Fix headless config for drive and gcs
      * Tell the user they should try again if the webserver method failed
      * Improve output of `--dump-headers`
    * S3
      * Allow anonymous access to public buckets
    * Swift
      * Stop chunked operations logging "Failed to read info: Object Not Found"
      * Use Content-Length on uploads for extra reliability
  * v1.22 - 2015-09-28
    * Implement rsync like include and exclude flags
    * swift
      * Support files > 5GB - thanks Sergey Tolmachev
  * v1.21 - 2015-09-22
    * New features
      * Display individual transfer progress
      * Make lsl output times in localtime
    * Fixes
      * Fix allowing user to override credentials again in Drive, GCS and ACD
    * Amazon Cloud Drive
      * Implement compliant pacing scheme
    * Google Drive
      * Make directory reads concurrent for increased speed.
  * v1.20 - 2015-09-15
    * New features
      * Amazon Cloud Drive support
      * Oauth support redone - fix many bugs and improve usability
        * Use "golang.org/x/oauth2" as oauth libary of choice
        * Improve oauth usability for smoother initial signup
        * drive, googlecloudstorage: optionally use auto config for the oauth token
      * Implement --dump-headers and --dump-bodies debug flags
      * Show multiple matched commands if abbreviation too short
      * Implement server side move where possible
    * local
      * Always use UNC paths internally on Windows - fixes a lot of bugs
    * dropbox
      * force use of our custom transport which makes timeouts work
    * Thanks to Klaus Post for lots of help with this release
  * v1.19 - 2015-08-28
    * New features
      * Server side copies for s3/swift/drive/dropbox/gcs
      * Move command - uses server side copies if it can
      * Implement --retries flag - tries 3 times by default
      * Build for plan9/amd64 and solaris/amd64 too
    * Fixes
      * Make a current version download with a fixed URL for scripting
      * Ignore rmdir in limited fs rather than throwing error
    * dropbox
      * Increase chunk size to improve upload speeds massively
      * Issue an error message when trying to upload bad file name
  * v1.18 - 2015-08-17
    * drive
      * Add `--drive-use-trash` flag so rclone trashes instead of deletes
      * Add "Forbidden to download" message for files with no downloadURL
    * dropbox
      * Remove datastore
        * This was deprecated and it caused a lot of problems
        * Modification times and MD5SUMs no longer stored
      * Fix uploading files > 2GB
    * s3
      * use official AWS SDK from github.com/aws/aws-sdk-go
      * **NB** will most likely require you to delete and recreate remote
      * enable multipart upload which enables files > 5GB
      * tested with Ceph / RadosGW / S3 emulation
      * many thanks to Sam Liston and Brian Haymore at the [Utah
        Center for High Performance Computing](https://www.chpc.utah.edu/) for a Ceph test account
    * misc
      * Show errors when reading the config file
      * Do not print stats in quiet mode - thanks Leonid Shalupov
      * Add FAQ
      * Fix created directories not obeying umask
      * Linux installation instructions - thanks Shimon Doodkin
  * v1.17 - 2015-06-14
    * dropbox: fix case insensitivity issues - thanks Leonid Shalupov
  * v1.16 - 2015-06-09
    * Fix uploading big files which was causing timeouts or panics
    * Don't check md5sum after download with --size-only
  * v1.15 - 2015-06-06
    * Add --checksum flag to only discard transfers by MD5SUM - thanks Alex Couper
    * Implement --size-only flag to sync on size not checksum & modtime
    * Expand docs and remove duplicated information
    * Document rclone's limitations with directories
    * dropbox: update docs about case insensitivity
  * v1.14 - 2015-05-21
    * local: fix encoding of non utf-8 file names - fixes a duplicate file problem
    * drive: docs about rate limiting
    * google cloud storage: Fix compile after API change in "google.golang.org/api/storage/v1"
  * v1.13 - 2015-05-10
    * Revise documentation (especially sync)
    * Implement --timeout and --conntimeout
    * s3: ignore etags from multipart uploads which aren't md5sums
  * v1.12 - 2015-03-15
    * drive: Use chunked upload for files above a certain size
    * drive: add --drive-chunk-size and --drive-upload-cutoff parameters
    * drive: switch to insert from update when a failed copy deletes the upload
    * core: Log duplicate files if they are detected
  * v1.11 - 2015-03-04
    * swift: add region parameter
    * drive: fix crash on failed to update remote mtime
    * In remote paths, change native directory separators to /
    * Add synchronization to ls/lsl/lsd output to stop corruptions
    * Ensure all stats/log messages to go stderr
    * Add --log-file flag to log everything (including panics) to file
    * Make it possible to disable stats printing with --stats=0
    * Implement --bwlimit to limit data transfer bandwidth
  * v1.10 - 2015-02-12
    * s3: list an unlimited number of items
    * Fix getting stuck in the configurator
  * v1.09 - 2015-02-07
    * windows: Stop drive letters (eg C:) getting mixed up with remotes (eg drive:)
    * local: Fix directory separators on Windows
    * drive: fix rate limit exceeded errors
  * v1.08 - 2015-02-04
    * drive: fix subdirectory listing to not list entire drive
    * drive: Fix SetModTime
    * dropbox: adapt code to recent library changes
  * v1.07 - 2014-12-23
    * google cloud storage: fix memory leak
  * v1.06 - 2014-12-12
    * Fix "Couldn't find home directory" on OSX
    * swift: Add tenant parameter
    * Use new location of Google API packages
  * v1.05 - 2014-08-09
    * Improved tests and consequently lots of minor fixes
    * core: Fix race detected by go race detector
    * core: Fixes after running errcheck
    * drive: reset root directory on Rmdir and Purge
    * fs: Document that Purger returns error on empty directory, test and fix
    * google cloud storage: fix ListDir on subdirectory
    * google cloud storage: re-read metadata in SetModTime
    * s3: make reading metadata more reliable to work around eventual consistency problems
    * s3: strip trailing / from ListDir()
    * swift: return directories without / in ListDir
  * v1.04 - 2014-07-21
    * google cloud storage: Fix crash on Update
  * v1.03 - 2014-07-20
    * swift, s3, dropbox: fix updated files being marked as corrupted
    * Make compile with go 1.1 again
  * v1.02 - 2014-07-19
    * Implement Dropbox remote
    * Implement Google Cloud Storage remote
    * Verify Md5sums and Sizes after copies
    * Remove times from "ls" command - lists sizes only
    * Add add "lsl" - lists times and sizes
    * Add "md5sum" command
  * v1.01 - 2014-07-04
    * drive: fix transfer of big files using up lots of memory
  * v1.00 - 2014-07-03
    * drive: fix whole second dates
  * v0.99 - 2014-06-26
    * Fix --dry-run not working
    * Make compatible with go 1.1
  * v0.98 - 2014-05-30
    * s3: Treat missing Content-Length as 0 for some ceph installations
    * rclonetest: add file with a space in
  * v0.97 - 2014-05-05
    * Implement copying of single files
    * s3 & swift: support paths inside containers/buckets
  * v0.96 - 2014-04-24
    * drive: Fix multiple files of same name being created
    * drive: Use o.Update and fs.Put to optimise transfers
    * Add version number, -V and --version
  * v0.95 - 2014-03-28
    * rclone.org: website, docs and graphics
    * drive: fix path parsing
  * v0.94 - 2014-03-27
    * Change remote format one last time
    * GNU style flags
  * v0.93 - 2014-03-16
    * drive: store token in config file
    * cross compile other versions
    * set strict permissions on config file
  * v0.92 - 2014-03-15
    * Config fixes and --config option
  * v0.91 - 2014-03-15
    * Make config file
  * v0.90 - 2013-06-27
    * Project named rclone
  * v0.00 - 2012-11-18
    * Project started

Bugs and Limitations
--------------------

### Empty directories are left behind / not created ##

With remotes that have a concept of directory, eg Local and Drive,
empty directories may be left behind, or not created when one was
expected.

This is because rclone doesn't have a concept of a directory - it only
works on objects.  Most of the object storage systems can't actually
store a directory so there is nowhere for rclone to store anything
about directories.

You can work round this to some extent with the`purge` command which
will delete everything under the path, **inluding** empty directories.

This may be fixed at some point in
[Issue #100](https://github.com/ncw/rclone/issues/100)

### Directory timestamps aren't preserved ##

For the same reason as the above, rclone doesn't have a concept of a
directory - it only works on objects, therefore it can't preserve the
timestamps of directories.

Frequently Asked Questions
--------------------------

### Do all cloud storage systems support all rclone commands ###

Yes they do.  All the rclone commands (eg `sync`, `copy` etc) will
work on all the remote storage systems.

### Can I copy the config from one machine to another ###

Sure!  Rclone stores all of its config in a single file.  If you want
to find this file, the simplest way is to run `rclone -h` and look at
the help for the `--config` flag which will tell you where it is.

See the [remote setup docs](http://rclone.org/remote_setup/) for more info.

### How do I configure rclone on a remote / headless box with no browser? ###

This has now been documented in its own [remote setup page](http://rclone.org/remote_setup/).

### Can rclone sync directly from drive to s3 ###

Rclone can sync between two remote cloud storage systems just fine.

Note that it effectively downloads the file and uploads it again, so
the node running rclone would need to have lots of bandwidth.

The syncs would be incremental (on a file by file basis).

Eg

    rclone sync drive:Folder s3:bucket


### Using rclone from multiple locations at the same time ###

You can use rclone from multiple places at the same time if you choose
different subdirectory for the output, eg

```
Server A> rclone sync /tmp/whatever remote:ServerA
Server B> rclone sync /tmp/whatever remote:ServerB
```

If you sync to the same directory then you should use rclone copy
otherwise the two rclones may delete each others files, eg

```
Server A> rclone copy /tmp/whatever remote:Backup
Server B> rclone copy /tmp/whatever remote:Backup
```

The file names you upload from Server A and Server B should be
different in this case, otherwise some file systems (eg Drive) may
make duplicates.

### Why doesn't rclone support partial transfers / binary diffs like rsync? ###

Rclone stores each file you transfer as a native object on the remote
cloud storage system.  This means that you can see the files you
upload as expected using alternative access methods (eg using the
Google Drive web interface).  There is a 1:1 mapping between files on
your hard disk and objects created in the cloud storage system.

Cloud storage systems (at least none I've come across yet) don't
support partially uploading an object. You can't take an existing
object, and change some bytes in the middle of it.

It would be possible to make a sync system which stored binary diffs
instead of whole objects like rclone does, but that would break the
1:1 mapping of files on your hard disk to objects in the remote cloud
storage system.

All the cloud storage systems support partial downloads of content, so
it would be possible to make partial downloads work.  However to make
this work efficiently this would require storing a significant amount
of metadata, which breaks the desired 1:1 mapping of files to objects.

### Can rclone do bi-directional sync? ###

No, not at present.  rclone only does uni-directional sync from A ->
B. It may do in the future though since it has all the primitives - it
just requires writing the algorithm to do it.

### Can I use rclone with an HTTP proxy? ###

Yes. rclone will use the environment variables `HTTP_PROXY`,
`HTTPS_PROXY` and `NO_PROXY`, similar to cURL and other programs.

`HTTPS_PROXY` takes precedence over `HTTP_PROXY` for https requests.

The environment values may be either a complete URL or a "host[:port]",
in which case the "http" scheme is assumed.

The `NO_PROXY` allows you to disable the proxy for specific hosts.
Hosts must be comma separated, and can contain domains or parts.
For instance "foo.com" also matches "bar.foo.com".

### Rclone gives x509: failed to load system roots and no roots provided error ###

This means that `rclone` can't file the SSL root certificates.  Likely
you are running `rclone` on a NAS with a cut-down Linux OS.

Rclone (via the Go runtime) tries to load the root certificates from
these places on Linux.

    "/etc/ssl/certs/ca-certificates.crt", // Debian/Ubuntu/Gentoo etc.
    "/etc/pki/tls/certs/ca-bundle.crt",   // Fedora/RHEL
    "/etc/ssl/ca-bundle.pem",             // OpenSUSE
    "/etc/pki/tls/cacert.pem",            // OpenELEC

So doing something like this should fix the problem.  It also sets the
time which is important for SSL to work properly.

```
mkdir -p /etc/ssl/certs/
curl -o /etc/ssl/certs/ca-certificates.crt https://raw.githubusercontent.com/bagder/ca-bundle/master/ca-bundle.crt
ntpclient -s -h pool.ntp.org
```

License
-------

This is free software under the terms of MIT the license (check the
COPYING file included with the source code).

```
Copyright (C) 2012 by Nick Craig-Wood http://www.craig-wood.com/nick/

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```

Authors
-------

  * Nick Craig-Wood <nick@craig-wood.com>

Contributors
------------

  * Alex Couper <amcouper@gmail.com>
  * Leonid Shalupov <leonid@shalupov.com>
  * Shimon Doodkin <helpmepro1@gmail.com>
  * Colin Nicholson <colin@colinn.com>
  * Klaus Post <klauspost@gmail.com>
  * Sergey Tolmachev <tolsi.ru@gmail.com>
  * Adriano Aurélio Meirelles <adriano@atinge.com>
  * C. Bess <cbess@users.noreply.github.com>
  * Dmitry Burdeev <dibu28@gmail.com>
  * Joseph Spurrier <github@josephspurrier.com>
  * Björn Harrtell <bjorn@wololo.org>
  * Xavier Lucas <xavier.lucas@corp.ovh.com>
  * Werner Beroux <werner@beroux.com>

Contact the rclone project
--------------------------

The project website is at:

  * https://github.com/ncw/rclone

There you can file bug reports, ask for help or contribute pull
requests.

See also

  * <a href="https://google.com/+RcloneOrg" rel="publisher">Google+ page for general comments</a></li>

Or email [Nick Craig-Wood](mailto:nick@craig-wood.com)

