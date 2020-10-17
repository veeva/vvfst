/*
This code serves as an example and is not meant for production use.

Copyright 2020 Veeva Systems Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing permissions
and limitations under the License.
*/
package util

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// Return truncated with ellipses when string is maximum than the given size
func FixedWidth(msg string, size int, ellipses bool) string {
	s := size
	e := ""
	if len(msg) > size && ellipses {
		s = s - 2
		e = ".."
	}

	// -20.20s
	//f := fmt.Sprintf("\\%\\-%d\\.%ds", s, s)
	f := fmt.Sprintf("%%-%d.%ds%s", s, s, e)
	return fmt.Sprintf(f, msg)
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func EndWithFileSeparator(s string) bool {
	return s[len(s)-1] == filepath.Separator
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func SplitParentAndName(path string) (string, string) {
	i := strings.LastIndex(path, "/")
	parent := path[0:i]
	name := path[i+1:]
	if parent == "" {
		parent = "/" // use the root folder when only / is the available
	}
	return parent, name
}

func GetFilename(path string) string {
	i := strings.LastIndex(path, "/")
	return path[i+1:]
}
