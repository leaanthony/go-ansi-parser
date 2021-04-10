package ansi

import (
	"fmt"
	"strconv"
	"strings"
)

// TextStyle is a type representing the
// ansi text styles
type TextStyle int

const (
	// Bold Style
	Bold TextStyle = 1 << 0
	// Faint Style
	Faint TextStyle = 1 << 1
	// Italic Style
	Italic TextStyle = 1 << 2
	// Blinking Style
	Blinking TextStyle = 1 << 3
	// Inversed Style
	Inversed TextStyle = 1 << 4
	// Invisible Style
	Invisible TextStyle = 1 << 5
	// Underlined Style
	Underlined TextStyle = 1 << 6
	// Strikethrough Style
	Strikethrough TextStyle = 1 << 7
)

// StyledText represents a single formatted string
type StyledText struct {
	Label string
	FgCol *Col
	BgCol *Col
	Style TextStyle
}

// Bold will return true if the text has a Bold style
func (s *StyledText) Bold() bool {
	return s.Style&Bold == Bold
}

// Faint will return true if the text has a Faint style
func (s *StyledText) Faint() bool {
	return s.Style&Faint == Faint
}

// Italic will return true if the text has an Italic style
func (s *StyledText) Italic() bool {
	return s.Style&Italic == Italic
}

// Blinking will return true if the text has a Blinking style
func (s *StyledText) Blinking() bool {
	return s.Style&Blinking == Blinking
}

// Inversed will return true if the text has an Inversed style
func (s *StyledText) Inversed() bool {
	return s.Style&Inversed == Inversed
}

// Invisible will return true if the text has an Invisible style
func (s *StyledText) Invisible() bool {
	return s.Style&Invisible == Invisible
}

// Underlined will return true if the text has an Underlined style
func (s *StyledText) Underlined() bool {
	return s.Style&Underlined == Underlined
}

// Strikethrough will return true if the text has a Strikethrough style
func (s *StyledText) Strikethrough() bool {
	return s.Style&Strikethrough == Strikethrough
}

// ColourMap maps ansi identifiers to a colour
var ColourMap = map[string]map[string]*Col{
	"Regular": {
		"30": Cols[0],
		"31": Cols[1],
		"32": Cols[2],
		"33": Cols[3],
		"34": Cols[4],
		"35": Cols[5],
		"36": Cols[6],
		"37": Cols[7],
	},
	"Bold": {
		"30": Cols[8],
		"31": Cols[9],
		"32": Cols[10],
		"33": Cols[11],
		"34": Cols[12],
		"35": Cols[13],
		"36": Cols[14],
		"37": Cols[15],
	},
	"Faint": {
		"30": Cols[0],
		"31": Cols[1],
		"32": Cols[2],
		"33": Cols[3],
		"34": Cols[4],
		"35": Cols[5],
		"36": Cols[6],
		"37": Cols[7],
	},
}

// Parse will convert an ansi encoded string and return
// a slice of StyledText structs that represent the text.
// If parsing is unsuccessful, an error is returned.
func Parse(input string) ([]*StyledText, error) {
	var result []*StyledText
	invalid := fmt.Errorf("invalid ansi string")
	missingTerminator := fmt.Errorf("missing escape terminator 'm'")
	invalidTrueColorSequence := fmt.Errorf("invalid TrueColor sequence")
	invalid256ColSequence := fmt.Errorf("invalid 256 colour sequence")
	index := 0
	var currentStyledText *StyledText = &StyledText{}

	if len(input) == 0 {
		return nil, invalid
	}

	for {
		// Read all chars to next escape code
		esc := strings.Index(input, "\033[")

		// If no more esc chars, save what's left and return
		if esc == -1 {
			text := input[index:]
			if len(text) > 0 {
				currentStyledText.Label = text
				result = append(result, currentStyledText)
			}
			return result, nil
		}
		label := input[:esc]
		if len(label) > 0 {
			currentStyledText.Label = label
			result = append(result, currentStyledText)
			currentStyledText = &StyledText{
				Label: "",
				FgCol: currentStyledText.FgCol,
				BgCol: currentStyledText.BgCol,
				Style: currentStyledText.Style,
			}
		}
		input = input[esc:]
		// skip
		input = input[2:]

		// Read in params
		endesc := strings.Index(input, "m")
		if endesc == -1 {
			return nil, missingTerminator
		}
		paramText := input[:endesc]
		input = input[endesc+1:]
		params := strings.Split(paramText, ";")
		colourMap := ColourMap["Regular"]
		skip := 0
		for index, param := range params {
			if skip > 0 {
				skip--
				continue
			}
			switch param {
			case "0":
				colourMap = ColourMap["Regular"]
				currentStyledText.Style = 0
				currentStyledText.FgCol = nil
				currentStyledText.BgCol = nil
			case "1":
				// Bold
				colourMap = ColourMap["Bold"]
				currentStyledText.Style |= Bold
			case "2":
				// Dim/Feint
				colourMap = ColourMap["Faint"]
				currentStyledText.Style |= Faint
			case "3":
				// Italic
				currentStyledText.Style |= Italic
			case "4":
				// Underlined
				currentStyledText.Style |= Underlined
			case "5":
				// Blinking
				currentStyledText.Style |= Blinking
			case "7":
				// Inversed
				currentStyledText.Style |= Inversed
			case "8":
				// Invisible
				currentStyledText.Style |= Invisible
			case "9":
				// Strikethrough
				currentStyledText.Style |= Strikethrough
			case "30", "31", "32", "33", "34", "35", "36", "37":
				currentStyledText.FgCol = colourMap[param]
			case "40", "41", "42", "43", "44", "45", "46", "47":
				bgcol := "3" + param[1:] // Equivalent of -10
				currentStyledText.BgCol = colourMap[bgcol]
			case "38", "48":
				if len(params)-index < 3 {
					return nil, invalid
				}
				// 256 colours
				if params[index+1] == "5" {
					skip = 2
					colIndexText := params[index+2]
					colIndex, err := strconv.Atoi(colIndexText)
					if err != nil {
						return nil, invalid256ColSequence
					}
					if colIndex < 0 || colIndex > 255 {
						return nil, invalid256ColSequence
					}
					if param == "38" {
						currentStyledText.FgCol = Cols[colIndex]
						continue
					}
					currentStyledText.BgCol = Cols[colIndex]
					continue
				}
				// we must have 4 params left
				if len(params)-index < 5 {
					return nil, invalidTrueColorSequence
				}
				if params[index+1] != "2" {
					return nil, invalidTrueColorSequence
				}
				var r, g, b uint8
				ri, err := strconv.Atoi(params[index+2])
				if err != nil {
					return nil, invalidTrueColorSequence
				}
				gi, err := strconv.Atoi(params[index+3])
				if err != nil {
					return nil, invalidTrueColorSequence
				}
				bi, err := strconv.Atoi(params[index+4])
				if err != nil {
					return nil, invalidTrueColorSequence
				}
				if bi > 255 || gi > 255 || ri > 255 {
					return nil, invalidTrueColorSequence
				}
				if bi < 0 || gi < 0 || ri < 0 {
					return nil, invalidTrueColorSequence
				}
				r = uint8(ri)
				g = uint8(gi)
				b = uint8(bi)
				skip = 4
				colvalue := fmt.Sprintf("#%02x%02x%02x", r, g, b)
				if param == "38" {
					currentStyledText.FgCol = &Col{Hex: colvalue, Rgb: Rgb{r, g, b}}
					continue
				}
				currentStyledText.BgCol = &Col{Hex: colvalue, Rgb: Rgb{r, g, b}}
			default:
				return nil, invalid
			}
		}
	}
}
