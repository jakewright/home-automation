// Copied from https://github.com/monzo/slog/blob/master/fmt.go
//
// The MIT License (MIT)
//
// Copyright (c) 2016 Mondo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package oops

import (
	"regexp"
	"strconv"
)

var formatterRe = regexp.MustCompile(`%` +
	`[\+\-# 0]*` + // Flags
	`(?:\d*\.|\[(\d+)\]\*\.)?(?:\d+|\[(\d+)\]\*)?` + // Width and precision
	`(?:\[(\d+)\])?` + // Argument index
	`[vTtbcdoqxXUbeEfFgGsqxXpt%]`, // Verb
)

func countFmtOperands(input string) int {
	count, point := 0, 0
	for _, match := range formatterRe.FindAllStringSubmatch(input, -1) {
		if match[0] == "%%" {
			// Deliberately match the regexp on %% (to prevent overlapping matches), but stop them here
			continue
		}

		for _, flag := range match[1:] {
			if flag == "" {
				continue
			} else if i, err := strconv.Atoi(flag); err == nil && i > 0 {
				point = i
				if point > count {
					count = point
				}
			}
		}
		if match[3] == "" {
			point++
		}
		if point > count {
			count = point
		}
	}
	return count
}
