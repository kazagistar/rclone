// Package dropbox provides an interface to Dropbox object storage
package dropbox

/*
Limitations of dropbox

File system is case insensitive
*/

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/ncw/rclone/fs"
	"github.com/ncw/rclone/oauthutil"
	"github.com/spf13/pflag"
	"github.com/stacktic/dropbox"
)

// Constants
const (
	rcloneAppKey    = "5jcck7diasz0rqy"
	rcloneAppSecret = "m8WRxJ6b1Z/Y25fDwJWS"
	metadataLimit   = dropbox.MetadataLimitDefault // max items to fetch at once
)

var (
	// A regexp matching path names for files Dropbox ignores
	// See https://www.dropbox.com/en/help/145 - Ignored files
	ignoredFiles = regexp.MustCompile(`(?i)(^|/)(desktop\.ini|thumbs\.db|\.ds_store|icon\r|\.dropbox|\.dropbox.attr)$`)
	// Upload chunk size - setting too small makes uploads slow.
	// Chunks aren't buffered into memory though so can set large.
	uploadChunkSize    = fs.SizeSuffix(128 * 1024 * 1024)
	maxUploadChunkSize = fs.SizeSuffix(150 * 1024 * 1024)
)

// Register with Fs
func init() {
	fs.Register(&fs.Info{
		Name:   "dropbox",
		NewFs:  NewFs,
		Config: configHelper,
		Options: []fs.Option{{
			Name: "app_key",
			Help: "Dropbox App Key - leave blank normally.",
		}, {
			Name: "app_secret",
			Help: "Dropbox App Secret - leave blank normally.",
		}},
	})
	pflag.VarP(&uploadChunkSize, "dropbox-chunk-size", "", fmt.Sprintf("Upload chunk size. Max %v.", maxUploadChunkSize))
}

// Configuration helper - called after the user has put in the defaults
func configHelper(name string) {
	// See if already have a token
	token := fs.ConfigFile.MustValue(name, "token")
	if token != "" {
		fmt.Printf("Already have a dropbox token - refresh?\n")
		if !fs.Confirm() {
			return
		}
	}

	// Get a dropbox
	db, err := newDropbox(name)
	if err != nil {
		log.Fatalf("Failed to create dropbox client: %v", err)
	}

	// This method will ask the user to visit an URL and paste the generated code.
	if err := db.Auth(); err != nil {
		log.Fatalf("Failed to authorize: %v", err)
	}

	// Get the token
	token = db.AccessToken()

	// Stuff it in the config file if it has changed
	old := fs.ConfigFile.MustValue(name, "token")
	if token != old {
		fs.ConfigFile.SetValue(name, "token", token)
		fs.SaveConfig()
	}
}

// Fs represents a remote dropbox server
type Fs struct {
	name           string           // name of this remote
	db             *dropbox.Dropbox // the connection to the dropbox server
	root           string           // the path we are working on
	slashRoot      string           // root with "/" prefix, lowercase
	slashRootSlash string           // root with "/" prefix and postfix, lowercase
}

// Object describes a dropbox object
type Object struct {
	fs          *Fs       // what this object is part of
	remote      string    // The remote path
	bytes       int64     // size of the object
	modTime     time.Time // time it was last modified
	hasMetadata bool      // metadata is valid
}

// ------------------------------------------------------------

// Name of the remote (as passed into NewFs)
func (f *Fs) Name() string {
	return f.name
}

// Root of the remote (as passed into NewFs)
func (f *Fs) Root() string {
	return f.root
}

// String converts this Fs to a string
func (f *Fs) String() string {
	return fmt.Sprintf("Dropbox root '%s'", f.root)
}

// Makes a new dropbox from the config
func newDropbox(name string) (*dropbox.Dropbox, error) {
	db := dropbox.NewDropbox()

	appKey := fs.ConfigFile.MustValue(name, "app_key")
	if appKey == "" {
		appKey = rcloneAppKey
	}
	appSecret := fs.ConfigFile.MustValue(name, "app_secret")
	if appSecret == "" {
		appSecret = fs.Reveal(rcloneAppSecret)
	}

	err := db.SetAppInfo(appKey, appSecret)
	return db, err
}

// NewFs contstructs an Fs from the path, container:path
func NewFs(name, root string) (fs.Fs, error) {
	if uploadChunkSize > maxUploadChunkSize {
		return nil, fmt.Errorf("Chunk size too big, must be < %v", maxUploadChunkSize)
	}
	db, err := newDropbox(name)
	if err != nil {
		return nil, err
	}
	f := &Fs{
		name: name,
		db:   db,
	}
	f.setRoot(root)

	// Read the token from the config file
	token := fs.ConfigFile.MustValue(name, "token")

	// Set our custom context which enables our custom transport for timeouts etc
	db.SetContext(oauthutil.Context())

	// Authorize the client
	db.SetAccessToken(token)

	// See if the root is actually an object
	entry, err := f.db.Metadata(f.slashRoot, false, false, "", "", metadataLimit)
	if err == nil && !entry.IsDir {
		remote := path.Base(f.root)
		newRoot := path.Dir(f.root)
		if newRoot == "." {
			newRoot = ""
		}
		f.setRoot(newRoot)
		obj := f.NewFsObject(remote)
		// return a Fs Limited to this object
		return fs.NewLimited(f, obj), nil
	}

	return f, nil
}

// Sets root in f
func (f *Fs) setRoot(root string) {
	f.root = strings.Trim(root, "/")
	lowerCaseRoot := strings.ToLower(f.root)

	f.slashRoot = "/" + lowerCaseRoot
	f.slashRootSlash = f.slashRoot
	if lowerCaseRoot != "" {
		f.slashRootSlash += "/"
	}
}

// Return an FsObject from a path
//
// May return nil if an error occurred
func (f *Fs) newFsObjectWithInfo(remote string, info *dropbox.Entry) fs.Object {
	o := &Object{
		fs:     f,
		remote: remote,
	}
	if info != nil {
		o.setMetadataFromEntry(info)
	} else {
		err := o.readEntryAndSetMetadata()
		if err != nil {
			// logged already fs.Debug("Failed to read info: %s", err)
			return nil
		}
	}
	return o
}

// NewFsObject returns an FsObject from a path
//
// May return nil if an error occurred
func (f *Fs) NewFsObject(remote string) fs.Object {
	return f.newFsObjectWithInfo(remote, nil)
}

// Strips the root off path and returns it
func (f *Fs) stripRoot(path string) *string {
	lowercase := strings.ToLower(path)

	if !strings.HasPrefix(lowercase, f.slashRootSlash) {
		fs.Stats.Error()
		fs.ErrorLog(f, "Path '%s' is not under root '%s'", path, f.slashRootSlash)
		return nil
	}

	stripped := path[len(f.slashRootSlash):]
	return &stripped
}

// Walk the root returning a channel of FsObjects
func (f *Fs) list(out fs.ListOpts) {
	// Track path component case, it could be different for entries coming from DropBox API
	// See https://www.dropboxforum.com/hc/communities/public/questions/201665409-Wrong-character-case-of-folder-name-when-calling-listFolder-using-Sync-API?locale=en-us
	// and https://github.com/ncw/rclone/issues/53
	nameTree := newNameTree()
	cursor := ""
	for {
		deltaPage, err := f.db.Delta(cursor, f.slashRoot)
		if err != nil {
			fs.Stats.Error()
			fs.ErrorLog(f, "Couldn't list: %s", err)
			break
		} else {
			if deltaPage.Reset && cursor != "" {
				fs.ErrorLog(f, "Unexpected reset during listing - try again")
				fs.Stats.Error()
				break
			}
			fs.Debug(f, "%d delta entries received", len(deltaPage.Entries))
			for i := range deltaPage.Entries {
				deltaEntry := &deltaPage.Entries[i]
				entry := deltaEntry.Entry
				if entry == nil {
					// This notifies of a deleted object
				} else {
					if len(entry.Path) <= 1 || entry.Path[0] != '/' {
						fs.Stats.Error()
						fs.ErrorLog(f, "dropbox API inconsistency: a path should always start with a slash and be at least 2 characters: %s", entry.Path)
						continue
					}

					lastSlashIndex := strings.LastIndex(entry.Path, "/")

					var parentPath string
					if lastSlashIndex == 0 {
						parentPath = ""
					} else {
						parentPath = entry.Path[1:lastSlashIndex]
					}
					lastComponent := entry.Path[lastSlashIndex+1:]

					if entry.IsDir {
						nameTree.PutCaseCorrectDirectoryName(parentPath, lastComponent)
					} else {
						parentPathCorrectCase := nameTree.GetPathWithCorrectCase(parentPath)
						if parentPathCorrectCase != nil {
							path := f.stripRoot(*parentPathCorrectCase + "/" + lastComponent)
							if path == nil {
								// an error occurred and logged by stripRoot
								continue
							}

							out <- f.newFsObjectWithInfo(*path, entry)
						} else {
							nameTree.PutFile(parentPath, lastComponent, entry)
						}
					}
				}
			}
			if !deltaPage.HasMore {
				break
			}
			cursor = deltaPage.Cursor.Cursor
		}
	}

	walkFunc := func(caseCorrectFilePath string, entry *dropbox.Entry) {
		path := f.stripRoot("/" + caseCorrectFilePath)
		if path == nil {
			// an error occurred and logged by stripRoot
			return
		}

		out <- f.newFsObjectWithInfo(*path, entry)
	}
	nameTree.WalkFiles(f.root, walkFunc)
}

// List walks the path returning a channel of FsObjects
func (f *Fs) List() fs.ListOpts {
	out := make(fs.ListOpts, fs.Config.Checkers)
	go func() {
		defer close(out)
		f.list(out)
	}()
	return out
}

// ListDir walks the path returning a channel of FsObjects
func (f *Fs) ListDir() fs.DirChan {
	out := make(fs.DirChan, fs.Config.Checkers)
	go func() {
		defer close(out)
		entry, err := f.db.Metadata(f.root, true, false, "", "", metadataLimit)
		if err != nil {
			fs.Stats.Error()
			fs.ErrorLog(f, "Couldn't list directories in root: %s", err)
		} else {
			for i := range entry.Contents {
				entry := &entry.Contents[i]
				if entry.IsDir {
					name := f.stripRoot(entry.Path)
					if name == nil {
						// an error occurred and logged by stripRoot
						continue
					}

					out <- &fs.Dir{
						Name:  *name,
						When:  time.Time(entry.ClientMtime),
						Bytes: entry.Bytes,
						Count: -1,
					}
				}
			}
		}
	}()
	return out
}

// A read closer which doesn't close the input
type readCloser struct {
	in io.Reader
}

// Read bytes from the object - see io.Reader
func (rc *readCloser) Read(p []byte) (n int, err error) {
	return rc.in.Read(p)
}

// Dummy close function
func (rc *readCloser) Close() error {
	return nil
}

// Put the object
//
// Copy the reader in to the new object which is returned
//
// The new object may have been created if an error is returned
func (f *Fs) Put(in io.Reader, remote string, modTime time.Time, size int64) (fs.Object, error) {
	// Temporary Object under construction
	o := &Object{
		fs:     f,
		remote: remote,
	}
	return o, o.Update(in, modTime, size)
}

// Mkdir creates the container if it doesn't exist
func (f *Fs) Mkdir() error {
	entry, err := f.db.Metadata(f.slashRoot, false, false, "", "", metadataLimit)
	if err == nil {
		if entry.IsDir {
			return nil
		}
		return fmt.Errorf("%q already exists as file", f.root)
	}
	_, err = f.db.CreateFolder(f.slashRoot)
	return err
}

// Rmdir deletes the container
//
// Returns an error if it isn't empty
func (f *Fs) Rmdir() error {
	entry, err := f.db.Metadata(f.slashRoot, true, false, "", "", 16)
	if err != nil {
		return err
	}
	if len(entry.Contents) != 0 {
		return errors.New("Directory not empty")
	}
	return f.Purge()
}

// Precision returns the precision
func (f *Fs) Precision() time.Duration {
	return fs.ModTimeNotSupported
}

// Copy src to this remote using server side copy operations.
//
// This is stored with the remote path given
//
// It returns the destination Object and a possible error
//
// Will only be called if src.Fs().Name() == f.Name()
//
// If it isn't possible then return fs.ErrorCantCopy
func (f *Fs) Copy(src fs.Object, remote string) (fs.Object, error) {
	srcObj, ok := src.(*Object)
	if !ok {
		fs.Debug(src, "Can't copy - not same remote type")
		return nil, fs.ErrorCantCopy
	}

	// Temporary Object under construction
	dstObj := &Object{
		fs:     f,
		remote: remote,
	}

	srcPath := srcObj.remotePath()
	dstPath := dstObj.remotePath()
	entry, err := f.db.Copy(srcPath, dstPath, false)
	if err != nil {
		return nil, fmt.Errorf("Copy failed: %s", err)
	}
	dstObj.setMetadataFromEntry(entry)
	return dstObj, nil
}

// Purge deletes all the files and the container
//
// Optional interface: Only implement this if you have a way of
// deleting all the files quicker than just running Remove() on the
// result of List()
func (f *Fs) Purge() error {
	// Let dropbox delete the filesystem tree
	_, err := f.db.Delete(f.slashRoot)
	return err
}

// Move src to this remote using server side move operations.
//
// This is stored with the remote path given
//
// It returns the destination Object and a possible error
//
// Will only be called if src.Fs().Name() == f.Name()
//
// If it isn't possible then return fs.ErrorCantMove
func (f *Fs) Move(src fs.Object, remote string) (fs.Object, error) {
	srcObj, ok := src.(*Object)
	if !ok {
		fs.Debug(src, "Can't move - not same remote type")
		return nil, fs.ErrorCantMove
	}

	// Temporary Object under construction
	dstObj := &Object{
		fs:     f,
		remote: remote,
	}

	srcPath := srcObj.remotePath()
	dstPath := dstObj.remotePath()
	entry, err := f.db.Move(srcPath, dstPath)
	if err != nil {
		return nil, fmt.Errorf("Move failed: %s", err)
	}
	dstObj.setMetadataFromEntry(entry)
	return dstObj, nil
}

// DirMove moves src to this remote using server side move operations.
//
// Will only be called if src.Fs().Name() == f.Name()
//
// If it isn't possible then return fs.ErrorCantDirMove
//
// If destination exists then return fs.ErrorDirExists
func (f *Fs) DirMove(src fs.Fs) error {
	srcFs, ok := src.(*Fs)
	if !ok {
		fs.Debug(srcFs, "Can't move directory - not same remote type")
		return fs.ErrorCantDirMove
	}

	// Check if destination exists
	entry, err := f.db.Metadata(f.slashRoot, false, false, "", "", metadataLimit)
	if err == nil && !entry.IsDeleted {
		return fs.ErrorDirExists
	}

	// Do the move
	_, err = f.db.Move(srcFs.slashRoot, f.slashRoot)
	if err != nil {
		return fmt.Errorf("MoveDir failed: %v", err)
	}
	return nil
}

// Hashes returns the supported hash sets.
func (f *Fs) Hashes() fs.HashSet {
	return fs.HashSet(fs.HashNone)
}

// ------------------------------------------------------------

// Fs returns the parent Fs
func (o *Object) Fs() fs.Fs {
	return o.fs
}

// Return a string version
func (o *Object) String() string {
	if o == nil {
		return "<nil>"
	}
	return o.remote
}

// Remote returns the remote path
func (o *Object) Remote() string {
	return o.remote
}

// Hash is unsupported on Dropbox
func (o *Object) Hash(t fs.HashType) (string, error) {
	return "", fs.ErrHashUnsupported
}

// Size returns the size of an object in bytes
func (o *Object) Size() int64 {
	return o.bytes
}

// setMetadataFromEntry sets the fs data from a dropbox.Entry
//
// This isn't a complete set of metadata and has an inacurate date
func (o *Object) setMetadataFromEntry(info *dropbox.Entry) {
	o.bytes = info.Bytes
	o.modTime = time.Time(info.ClientMtime)
	o.hasMetadata = true
}

// Reads the entry from dropbox
func (o *Object) readEntry() (*dropbox.Entry, error) {
	entry, err := o.fs.db.Metadata(o.remotePath(), false, false, "", "", metadataLimit)
	if err != nil {
		fs.Debug(o, "Error reading file: %s", err)
		return nil, fmt.Errorf("Error reading file: %s", err)
	}
	return entry, nil
}

// Read entry if not set and set metadata from it
func (o *Object) readEntryAndSetMetadata() error {
	// Last resort set time from client
	if !o.modTime.IsZero() {
		return nil
	}
	entry, err := o.readEntry()
	if err != nil {
		return err
	}
	o.setMetadataFromEntry(entry)
	return nil
}

// Returns the remote path for the object
func (o *Object) remotePath() string {
	return o.fs.slashRootSlash + o.remote
}

// Returns the key for the metadata database for a given path
func metadataKey(path string) string {
	// NB File system is case insensitive
	path = strings.ToLower(path)
	hash := md5.New()
	_, _ = hash.Write([]byte(path))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Returns the key for the metadata database
func (o *Object) metadataKey() string {
	return metadataKey(o.remotePath())
}

// readMetaData gets the info if it hasn't already been fetched
func (o *Object) readMetaData() (err error) {
	if o.hasMetadata {
		return nil
	}
	// Last resort
	return o.readEntryAndSetMetadata()
}

// ModTime returns the modification time of the object
//
// It attempts to read the objects mtime and if that isn't present the
// LastModified returned in the http headers
func (o *Object) ModTime() time.Time {
	err := o.readMetaData()
	if err != nil {
		fs.Log(o, "Failed to read metadata: %s", err)
		return time.Now()
	}
	return o.modTime
}

// SetModTime sets the modification time of the local fs object
//
// Commits the datastore
func (o *Object) SetModTime(modTime time.Time) {
	// FIXME not implemented
	return
}

// Storable returns whether this object is storable
func (o *Object) Storable() bool {
	return true
}

// Open an object for read
func (o *Object) Open() (in io.ReadCloser, err error) {
	in, _, err = o.fs.db.Download(o.remotePath(), "", 0)
	return
}

// Update the already existing object
//
// Copy the reader into the object updating modTime and size
//
// The new object may have been created if an error is returned
func (o *Object) Update(in io.Reader, modTime time.Time, size int64) error {
	remote := o.remotePath()
	if ignoredFiles.MatchString(remote) {
		fs.Log(o, "File name disallowed - not uploading")
		return nil
	}
	entry, err := o.fs.db.UploadByChunk(ioutil.NopCloser(in), int(uploadChunkSize), remote, true, "")
	if err != nil {
		return fmt.Errorf("Upload failed: %s", err)
	}
	o.setMetadataFromEntry(entry)
	return nil
}

// Remove an object
func (o *Object) Remove() error {
	_, err := o.fs.db.Delete(o.remotePath())
	return err
}

// Check the interfaces are satisfied
var (
	_ fs.Fs       = (*Fs)(nil)
	_ fs.Copier   = (*Fs)(nil)
	_ fs.Purger   = (*Fs)(nil)
	_ fs.Mover    = (*Fs)(nil)
	_ fs.DirMover = (*Fs)(nil)
	_ fs.Object   = (*Object)(nil)
)
