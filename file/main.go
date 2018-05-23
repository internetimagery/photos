
// safe file writing.
// 1) check file doesn't exist already (watch for race conditions)
// 2) write file to different name in same dir
// 3) move file to different name
// 4) move new file over old file
// 5) remove old file
