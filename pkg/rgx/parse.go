package rgx

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jdetok/gopuml/pkg/dir"
)

// same idea as one with RgxReady passed, but capture package
func (r *Rgx) Parse(dm *dir.DirMap) error {
	for d, fm := range *dm {
		r.RDirMap[d] = RgxFileMap{}
		// fmt.Println("dir in parse: ", d)x
		for n := range fm {
			f, err := dm.OpenFile(d, n)

			if err != nil {
				return err
			}
			defer f.Close()
			rf, err := r.RgxParseFile(d, f)
			if err != nil {
				return err
			}
			// map *RgxFile to file name, which is mapped to dir name
			if r.RDirMap[d][f.Name()] == nil {
				r.RDirMap[d][f.Name()] = &RgxFile{}
			}
			r.RDirMap[d][f.Name()] = rf
			if r.PkgMap[rf.RgxPkg] == nil {
				r.PkgMap[rf.RgxPkg] = RgxFileMap{}
			}
			r.PkgMap[rf.RgxPkg][f.Name()] = rf

		}
	}
	return nil
}

// use bufio scanner to iterate through each line in passed file
func (r *Rgx) RgxParseFile(dir string, f *os.File) (*RgxFile, error) {
	defer f.Close()
	// fmt.Printf("parsing %s...\n", filepath.Base(f.Name()))

	if r.RDirMap[dir][f.Name()] == nil {
		r.RDirMap[dir][f.Name()] = &RgxFile{}
	}

	rf := RgxFile{
		Pkg:  dir,
		Name: f.Name(),
	}

	lineCount := 0
	// insideStruct := false // used to handle struct fields
	// scan each line of the file
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineCount++
		line := scanner.Text() // get string of current line
		rm := r.RgxParseLine(line)
		if rm == nil {
			continue
		}

		switch rm.FindType {
		case PKG:
			rf.RgxPkg = rm.Groups[0]
			fmt.Println("package", rm.Groups[0], "for file", rf.Name)
		case FUNC:
			rf.Funcs = append(rf.Funcs,
				&RgxFunc{
					Name:   rm.Groups[0],
					Params: rm.Groups[1],
					Rtn:    rm.Groups[2],
				},
			)
			// fmt.Printf("found func %s in RgxFile Pkg %s | File %s\n", rm.Groups[0], rf.Pkg, rf.Name)
		case STRUCT:
			var s = &RgxStruct{
				Name: rm.Groups[0],
			} // get fields from struct (scans lines until '}' is reached)
			matches := r.RgxParseStruct(scanner, &lineCount)
			for _, m := range matches {
				// fmt.Printf("%v\n", *m)
				s.Fields = append(s.Fields,
					RgxStructFld{
						Name:  m.Groups[0],
						DType: m.Groups[1],
					},
				)
			}
			rf.Structs = append(rf.Structs, s)
			// fmt.Printf("found struct %s in RgxFile Pkg %s | File %s\n", rm.Groups[0], rf.Pkg, rf.Name)
		case MTHD:
			rf.Methods = append(rf.Methods,
				&RgxFunc{
					IsMthd:    true,
					BelongsTo: rm.Groups[0],
					Name:      rm.Groups[1],
					Params:    rm.Groups[2],
					Rtn:       rm.Groups[3],
				},
			)
			// fmt.Printf("found method %s in %s in RgxFile Pkg %s | File %s\n", rm.Groups[1], rm.Groups[0], rf.Pkg, rf.Name)
		}

	}
	if r.PkgMap[rf.RgxPkg] == nil {
		r.PkgMap[rf.RgxPkg] = RgxFileMap{}
	}
	if r.PkgMap[rf.RgxPkg][f.Name()] == nil {
		r.PkgMap[rf.RgxPkg][f.Name()] = &RgxFile{}
	}
	r.PkgMap[rf.RgxPkg][f.Name()] = &rf
	return &rf, nil
}

// pass line bytes and linenum, check for regex matches
// return matches and a bool signaling whether next line is within a struct
func (r *Rgx) RgxParseLine(line string) *RgxMatch {
	var m RgxMatch
	pkgRgx := r.ready[PKG]
	if pkgMatch := pkgRgx.FindStringSubmatch(line); pkgMatch != nil {
		m.FindType = PKG
		m.MatchStr = pkgMatch[0]
		m.Groups = pkgMatch[1:]
	}
	for _, key := range RGX_CHECK_ORDER {
		rgx, ok := r.ready[key]
		if !ok {
			continue
		}
		if matches := rgx.FindStringSubmatch(line); matches != nil {
			m.FindType = key
			m.MatchStr = matches[0]
			m.Groups = matches[1:]
			// groupStr := str.AngleWrap(m.Groups)
			// fmt.Printf("MATCH: %s\nKEY: %s | GROUPS: %s\n", m.FindType, key, groupStr)
			// if m.FindType != STRUCT {
			// 	fmt.Println()
			// }
			return &m
		}
	}
	return nil
}

func (r *Rgx) RgxParseStruct(s *bufio.Scanner, lineCount *int) []*RgxMatch {
	insideStruct := true
	matches := []*RgxMatch{}
	for s.Scan() && insideStruct {
		var rm *RgxMatch
		*lineCount++
		structLine := s.Text()

		rm, insideStruct = r.RgxParseStructFld(structLine)
		if rm == nil {
			if insideStruct {
				continue
			} else {
				break
			}
		}
		matches = append(matches, rm)
	}

	// fmt.Printf("finished scanning struct from lines %d - %d\n", startLine, *lineCount)
	return matches
}

// match fields inside a struct, only called from RgxParseStruct
// the bool returned determines whether the struct has more fields or not
func (r *Rgx) RgxParseStructFld(line string) (*RgxMatch, bool) {
	// check for ; at end of struct def - send signal to end RgxParseStruct
	if r.ready[ENDSTMNT].MatchString(line) {
		return nil, false
	}
	// if not a ; check for fields | if nil return nil but true continue signal
	matches := r.ready[STRUCTFLD].FindStringSubmatch(line)
	if matches == nil {
		return nil, true
	}
	// if field found build match struct
	return &RgxMatch{
		FindType: STRUCTFLD,
		MatchStr: matches[0],
		Groups:   matches[1:],
	}, true
}
