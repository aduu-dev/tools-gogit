package install

import (
	"fmt"
	"regexp"
	"strings"
)

func combinedLine(line string, comment string) string {
	return strings.TrimSpace(fmt.Sprintf("%s # %s", strings.TrimSpace(line), strings.TrimSpace(comment)))
}

var (
	errLineContainsLineSeparators = fmt.Errorf("error: line contains one or more line separators")
	errCommentContainsLineSeparators = fmt.Errorf("error: comment contains line separators")
	errCommentIsEmpty = fmt.Errorf("error: comment is empty")
	errTwoCommentsMatch = fmt.Errorf("error: two in-file comments match the given one")
	errCommentContainsHash = fmt.Errorf("error: comment contains hash")
)

func commentLineIndexes(lines []string, comment string) (indexes []int, err error) {
	containsComment, err := regexp.Compile(`# [\s]*` + regexp.QuoteMeta(comment))
	if err != nil {
		return
	}

	var indexe []int
	for i, l := range lines {
		if containsComment.MatchString(l) {
			indexe = append(indexe, i)
		}
	}

	return indexe, nil
}

func fixLineSeparators(lines []string) (out []string){
	// Remove all spaces until only one is between line and a non-empty line.
	out = make([]string, 0, len(lines))
	for _, l := range lines {
		out = append(out, l)
	}

	lines = out

done:
	for true {
		last := len(lines) - 1

		switch {
		case last == 0:
			break done
		case last == 1:
			if len(lines[0]) == 0 {
				lines = lines[1:]
			} else {
				// Add empty line
				lines = []string{lines[0], "", lines[1]}
			}

			break done
		default:
			// I can use my algorith.
			break
		}

		beforeIndex := last - 1
		beforeBeforeIndex := last - 2

		before := lines[beforeIndex]
		beforeBefore := lines[beforeBeforeIndex]

		if len(before) != 0 {
			// add empty line and we are finished.
			lines = append(lines, "")
			tmp := lines[last]
			lines[last] = ""
			lines[last + 1] = tmp

			break
		}

		if len(before) == 0 && len(beforeBefore) != 0 {
			break
		}

		// before is zero, so remove it.
		lines = append(lines[0:beforeIndex], lines[last])
	}

	return lines
}

// addToShellFile adds the given line with the given comment to the file.
//
// If there is a line with the given comment already it replaces the line's content
// with the given line + comment.
// It ignores further lines with the same comment.
func addToShellFile(content string, line string, comment string) (changedContent string, err error) {
	lines := strings.Split(content, "\n")

	if len(comment) == 0 {
		return "", errCommentIsEmpty
	}

	if strings.Contains(comment, "\n") {
		return "", errCommentContainsLineSeparators
	}

	if strings.Contains(line, "\n") {
		return "", errLineContainsLineSeparators
	}

	if strings.Contains(comment, "#") {
		return "", errCommentContainsHash
	}

	matchedLines, err := commentLineIndexes(lines, comment)
	if err != nil {
		return
	}

	// Add the line at the back.
	if len(matchedLines) == 0 {
		lines = append(lines, combinedLine(line, comment))

		lines = fixLineSeparators(lines)

		return strings.Join(lines, "\n"), nil
	}

	// Add the line where the first matched line is.
	index := matchedLines[0]

	lines[index] = combinedLine(line, comment)

	return strings.Join(lines, "\n"), nil
}

// removeFromShellFile removes one line containing the given comment.
//
// If there are two lines with that comment then reutrn an error.
func removeFromShellFile(content string, comment string) (changedContent string, err error) {
	lines := strings.Split(content, "\n")

	matchedLines, err := commentLineIndexes(lines, comment)
	if err != nil {
		return
	}

	if len(matchedLines) == 0 {
		return content, nil
	}

	if len(matchedLines) != 1 {
		return "", errTwoCommentsMatch
	}

	i := matchedLines[0]

	lines = append(lines[0:i], lines[i+1:]...)

	// Remove empty lines.
	for i := len(lines) - 1; i >= 0; i-- {
		// Only as long as the last line is empty do we continue.
		if len(lines[i]) != 0 {
			break
		}

		lines = lines[0:i]
	}

	return strings.Join(lines,"\n"), nil
}
