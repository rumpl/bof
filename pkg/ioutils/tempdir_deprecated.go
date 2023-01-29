package ioutils

import "github.com/rumpl/bof/pkg/longpath"

// TempDir is the equivalent of [os.MkdirTemp], except that on Windows
// the result is in Windows longpath format. On Unix systems it is
// equivalent to [os.MkdirTemp].
//
// Deprecated: use [longpath.MkdirTemp].
var TempDir = longpath.MkdirTemp
