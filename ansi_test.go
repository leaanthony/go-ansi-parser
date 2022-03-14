package ansi

import (
	"testing"

	is "github.com/matryer/is"
)

func TestParseAnsi16Styles(t *testing.T) {
	is2 := is.New(t)
	var got []*StyledText
	var err error

	// Bold
	got, err = Parse("\u001b[1;30mHello World\033[0m")
	is2.NoErr(err)
	is2.Equal(len(got), 1)
	is2.True(got[0].Bold())
	// Faint
	got, err = Parse("\u001b[2;30mHello World\033[0m")
	is2.NoErr(err)
	is2.Equal(len(got), 1)
	is2.True(got[0].Faint())
	// Italic
	got, err = Parse("\u001b[3;30mHello World\033[0m")
	is2.NoErr(err)
	is2.Equal(len(got), 1)
	is2.True(got[0].Italic())
	// Underlined
	got, err = Parse("\u001b[4;30mHello World\033[0m")
	is2.NoErr(err)
	is2.Equal(len(got), 1)
	is2.True(got[0].Underlined())
	// Blinking
	got, err = Parse("\u001b[5;30mHello World\033[0m")
	is2.NoErr(err)
	is2.Equal(len(got), 1)
	is2.True(got[0].Blinking())
	// Inversed
	got, err = Parse("\u001b[7;30mHello World\033[0m")
	is2.NoErr(err)
	is2.Equal(len(got), 1)
	is2.True(got[0].Inversed())
	// Invisible
	got, err = Parse("\u001b[8;30mHello World\033[0m")
	is2.NoErr(err)
	is2.Equal(len(got), 1)
	is2.True(got[0].Invisible())
	// Strikethrough
	got, err = Parse("\u001b[9;30mHello World\033[0m")
	is2.NoErr(err)
	is2.Equal(len(got), 1)
	is2.True(got[0].Strikethrough())
}

func TestParseAnsi16Swap(t *testing.T) {
	is2 := is.New(t)
	var got []*StyledText
	var err error

	// Swap single
	c0ffee := &Col{
		Id:   0,
		Hex:  "c0ffee",
		Name: "Coffee",
	}
	original := ColourMap["Regular"]["30"]
	ColourMap["Regular"]["30"] = c0ffee
	got, err = Parse("\u001b[0;30mHello World\033[0m")
	is2.NoErr(err)
	is2.Equal(len(got), 1)
	is2.Equal(got[0].FgCol.Name, "Coffee")

	// Restore
	ColourMap["Regular"]["30"] = original
	got, err = Parse("\u001b[0;30mHello World\033[0m")
	is2.NoErr(err)
	is2.Equal(len(got), 1)
	is2.Equal(got[0].FgCol.Name, "Black")
}

func TestParseAnsi16SingleColour(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name      string
		input     string
		wantText  string
		wantColor string
		wantErr   bool
	}{
		{"No formatting", "Hello World", "Hello World", "", false},
		{"Black", "\u001b[0;30mHello World\033[0m", "Hello World", "Black", false},
		{"Red", "\u001b[0;31mHello World\033[0m", "Hello World", "Maroon", false},
		{"Green", "\u001b[0;32mGreen\033[0m", "Green", "Green", false},
		{"Yellow", "\u001b[0;33m😀\033[0m", "😀", "Olive", false},
		{"Blue", "\u001b[0;34m123\033[0m", "123", "Navy", false},
		{"Purple", "\u001b[0;35m👩🏽‍🔧\u001B[0m", "👩🏽‍🔧", "Purple", false},
		{"Cyan", "\033[0;36m😀\033[0m", "😀", "Teal", false},
		{"White", "\u001b[0;37m[0;37m\033[0m", "[0;37m", "Silver", false},
		{"Black Bold", "\u001b[1;30mHello World\033[0m", "Hello World", "Grey", false},
		{"Red Bold", "\u001b[1;31mHello World\033[0m", "Hello World", "Red", false},
		{"Green Bold", "\u001b[1;32mGreen\033[0m", "Green", "Lime", false},
		{"Yellow Bold", "\u001b[1;33m😀\033[0m", "😀", "Yellow", false},
		{"Blue Bold", "\u001b[1;34m123\033[0m", "123", "Blue", false},
		{"Purple Bold", "\u001b[1;35m👩🏽‍🔧\u001B[0m", "👩🏽‍🔧", "Fuchsia", false},
		{"Cyan Bold", "\033[1;36m😀\033[0m", "😀", "Aqua", false},
		{"White Bold", "\u001b[1;37m[0;37m\033[0m", "[0;37m", "White", false},
		{"Blank", "", "", "", true},
		{"Emoji", "😀👩🏽‍🔧", "😀👩🏽‍🔧", "", false},
		{"Spaces", "  ", "  ", "", false},
		{"Bad code", "\u001b[1  ", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			is2.Equal(err != nil, tt.wantErr)
			expectedLength := 1
			if tt.wantErr {
				expectedLength = 0
			}
			is2.Equal(len(got), expectedLength)
			if expectedLength == 1 {
				if len(tt.wantColor) > 0 {
					is2.True(got[0].FgCol != nil)
					is2.Equal(got[0].FgCol.Name, tt.wantColor)
				}
			}
		})
	}
}

func TestParseAnsi16SingleBGColour(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name      string
		input     string
		wantText  string
		wantColor string
		wantErr   bool
	}{
		{"No formatting", "Hello World", "Hello World", "", false},
		{"Black", "\u001b[0;40mHello World\033[0m", "Hello World", "Black", false},
		{"Red", "\u001b[0;41mHello World\033[0m", "Hello World", "Maroon", false},
		{"Green", "\u001b[0;42mGreen\033[0m", "Green", "Green", false},
		{"Yellow", "\u001b[0;43m😀\033[0m", "😀", "Olive", false},
		{"Blue", "\u001b[0;44m123\033[0m", "123", "Navy", false},
		{"Purple", "\u001b[0;45m👩🏽‍🔧\u001B[0m", "👩🏽‍🔧", "Purple", false},
		{"Cyan", "\033[0;46m😀\033[0m", "😀", "Teal", false},
		{"White", "\u001b[0;47m[0;47m\033[0m", "[0;47m", "Silver", false},
		{"Black Bold", "\u001b[1;40mHello World\033[0m", "Hello World", "Grey", false},
		{"Red Bold", "\u001b[1;41mHello World\033[0m", "Hello World", "Red", false},
		{"Green Bold", "\u001b[1;42mGreen\033[0m", "Green", "Lime", false},
		{"Yellow Bold", "\u001b[1;43m😀\033[0m", "😀", "Yellow", false},
		{"Blue Bold", "\u001b[1;44m123\033[0m", "123", "Blue", false},
		{"Purple Bold", "\u001b[1;45m👩🏽‍🔧\u001B[0m", "👩🏽‍🔧", "Fuchsia", false},
		{"Cyan Bold", "\033[1;46m😀\033[0m", "😀", "Aqua", false},
		{"White Bold", "\u001b[1;47m[0;47m\033[0m", "[0;47m", "White", false},
		{"Pre text", "Hello\u001b[0m", "Hello", "", false},
		{"Blank", "", "", "", true},
		{"Emoji", "😀👩🏽‍🔧", "😀👩🏽‍🔧", "", false},
		{"Spaces", "  ", "  ", "", false},
		{"Bad code", "\u001b[1  ", "", "", true},
		{"No colour", "\033[m\033[40m    \033[0m", "    ", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			is2.Equal(err != nil, tt.wantErr)
			expectedLength := 1
			if tt.wantErr {
				expectedLength = 0
			}
			is2.Equal(len(got), expectedLength)
			if expectedLength == 1 {
				if len(tt.wantColor) > 0 {
					is2.True(got[0].BgCol != nil)
					is2.Equal(got[0].BgCol.Name, tt.wantColor)
				}
			}
		})
	}
}

func TestParseAnsi16MultiColour(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name    string
		input   string
		want    []*StyledText
		wantErr bool
	}{
		{"Black & Red", "\u001B[0;30mHello World\u001B[0m\u001B[0;31mHello World\u001B[0m", []*StyledText{
			{Label: "Hello World", FgCol: &Col{Name: "Black"}},
			{Label: "Hello World", FgCol: &Col{Name: "Maroon"}},
		}, false},
		{"Text then Black & Red", "This is great!\u001B[0;30mHello World\u001B[0m\u001B[0;31mHello World\u001B[0m", []*StyledText{
			{Label: "This is great!"},
			{Label: "Hello World", FgCol: &Col{Name: "Black"}},
			{Label: "Hello World", FgCol: &Col{Name: "Maroon"}},
		}, false},
		{"Text Reset then Black & Red", "This is great!\u001B[0m\u001B[0;30mHello World\u001B[0m\u001B[0;31mHello World\u001B[0m", []*StyledText{
			{Label: "This is great!"},
			{Label: "Hello World", FgCol: &Col{Name: "Black"}},
			{Label: "Hello World", FgCol: &Col{Name: "Maroon"}},
		}, false},
		{"Text Reset then Black & Red", "This is great!\u001B[0m", []*StyledText{
			{Label: "This is great!"},
		}, false},
		{"Black & Red no reset", "\u001B[0;30mHello World\u001B[0;31mHello World", []*StyledText{
			{Label: "Hello World", FgCol: &Col{Name: "Black"}},
			{Label: "Hello World", FgCol: &Col{Name: "Maroon"}},
		}, false},
		{"Black,space,Red", "\u001B[0;30mHello World\u001B[0m \u001B[0;31mHello World\u001B[0m", []*StyledText{
			{Label: "Hello World", FgCol: &Col{Name: "Black"}},
			{Label: " "},
			{Label: "Hello World", FgCol: &Col{Name: "Maroon"}},
		}, false},
		{"Black,Red,Blue,Green underlined", "\033[4;30mBlack\u001B[0m\u001B[4;31mRed\u001B[0m\u001B[4;34mBlue\u001B[0m\u001B[4;32mGreen\u001B[0m", []*StyledText{
			{Label: "Black", FgCol: &Col{Name: "Black"}, Style: Underlined},
			{Label: "Red", FgCol: &Col{Name: "Maroon"}, Style: Underlined},
			{Label: "Blue", FgCol: &Col{Name: "Navy"}, Style: Underlined},
			{Label: "Green", FgCol: &Col{Name: "Green"}, Style: Underlined},
		}, false},
		{"Black,Red,Blue,Green bold", "\033[1;30mBlack\u001B[0m\u001B[1;31mRed\u001B[0m\u001B[1;34mBlue\u001B[0m\u001B[1;32mGreen\u001B[0m", []*StyledText{
			{Label: "Black", FgCol: &Col{Name: "Grey"}, Style: Bold},
			{Label: "Red", FgCol: &Col{Name: "Red"}, Style: Bold},
			{Label: "Blue", FgCol: &Col{Name: "Blue"}, Style: Bold},
			{Label: "Green", FgCol: &Col{Name: "Lime"}, Style: Bold},
		}, false},
		{"Green Feint & Yellow Italic", "\u001B[2;32m👩🏽‍🔧\u001B[0m\u001B[0;3;33m👩🏽‍🔧\u001B[0m", []*StyledText{
			{Label: "👩🏽‍🔧", FgCol: &Col{Name: "Green"}, Style: Faint},
			{Label: "👩🏽‍🔧", FgCol: &Col{Name: "Olive"}, Style: Italic},
		}, false},
		{"Green Blinking & Yellow Inversed", "\u001B[5;32m👩🏽‍🔧\u001B[0m\u001B[0;7;33m👩🏽‍🔧\u001B[0m", []*StyledText{
			{Label: "👩🏽‍🔧", FgCol: &Col{Name: "Green"}, Style: Blinking},
			{Label: "👩🏽‍🔧", FgCol: &Col{Name: "Olive"}, Style: Inversed},
		}, false},
		{"Green Invisible & Yellow Invisible & Strikethrough", "\u001B[8;32m👩🏽‍🔧\u001B[9;33m👩🏽‍🔧\u001B[0m", []*StyledText{
			{Label: "👩🏽‍🔧", FgCol: &Col{Name: "Green"}, Style: Invisible},
			{Label: "👩🏽‍🔧", FgCol: &Col{Name: "Olive"}, Style: Invisible | Strikethrough},
		}, false},
		{"Red Foregraound & Black Background", "\u001b[1;31;40mHello World\033[0m", []*StyledText{
			{Label: "Hello World", FgCol: &Col{Name: "Red"}, BgCol: &Col{Name: "Black"}, Style: Bold},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			is2.Equal(err != nil, tt.wantErr)
			for index, w := range tt.want {
				is2.Equal(got[index].Label, w.Label)
				if w.FgCol != nil {
					is2.Equal(got[index].FgCol.Name, w.FgCol.Name)
				}
				is2.Equal(got[index].Style, w.Style)
			}
		})
	}
}

func TestParseAnsi256(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name    string
		input   string
		want    []*StyledText
		wantErr bool
	}{
		{"Grey93 & DarkViolet", "\u001B[38;5;255mGrey93\u001B[0m\u001B[38;5;128mDarkViolet\u001B[0m", []*StyledText{
			{Label: "Grey93", FgCol: &Col{Name: "Grey93"}},
			{Label: "DarkViolet", FgCol: &Col{Name: "DarkViolet"}},
		}, false},
		{"Grey93 Bold & DarkViolet Italic", "\u001B[0;1;38;5;255mGrey93\u001B[0m\u001B[0;3;38;5;128mDarkViolet\u001B[0m", []*StyledText{
			{Label: "Grey93", FgCol: &Col{Name: "Grey93"}, Style: Bold},
			{Label: "DarkViolet", FgCol: &Col{Name: "DarkViolet"}, Style: Italic},
		}, false},
		{"Grey93 Bold & DarkViolet Italic", "\u001B[0;1;38;5;256mGrey93\u001B[0m", nil, true},
		{"Grey93 Bold & DarkViolet Italic", "\u001B[0;1;38;5;-1mGrey93\u001B[0m", nil, true},
		{"Bad No of Params", "\u001B[0;1;38;5mGrey93\u001B[0m", nil, true},
		{"Bad Params", "\u001B[0;1;38;fivemGrey93\u001B[0m", nil, true},
		{"Bad Params 2", "\u001B[0;1;38;3mGrey93\u001B[0m", nil, true},
		{"Bad Params 3", "\u001B[0;1;38;5;fivemGrey93\u001B[0m", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			is2.Equal(err != nil, tt.wantErr)
			for index, w := range tt.want {
				is2.Equal(got[index].Label, w.Label)
				if w.FgCol != nil {
					is2.Equal(got[index].FgCol.Name, w.FgCol.Name)
				}
				is2.Equal(got[index].Style, w.Style)
			}
		})
	}
}

func TestParseAnsiBG256(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name    string
		input   string
		want    []*StyledText
		wantErr bool
	}{
		{"Grey93 & DarkViolet", "\u001B[48;5;255mGrey93\u001B[0m\u001B[48;5;128mDarkViolet\u001B[0m", []*StyledText{
			{Label: "Grey93", BgCol: &Col{Name: "Grey93"}},
			{Label: "DarkViolet", BgCol: &Col{Name: "DarkViolet"}},
		}, false},
		{"Grey93 Bold & DarkViolet Italic", "\u001B[0;1;48;5;255mGrey93\u001B[0m\u001B[0;3;48;5;128mDarkViolet\u001B[0m", []*StyledText{
			{Label: "Grey93", BgCol: &Col{Name: "Grey93"}, Style: Bold},
			{Label: "DarkViolet", BgCol: &Col{Name: "DarkViolet"}, Style: Italic},
		}, false},
		{"Grey93 Bold & DarkViolet Italic", "\u001B[0;1;48;5;256mGrey93\u001B[0m", nil, true},
		{"Grey93 Bold & DarkViolet Italic", "\u001B[0;1;48;5;-1mGrey93\u001B[0m", nil, true},
		{"Bad No of Params", "\u001B[0;1;48;5mGrey93\u001B[0m", nil, true},
		{"Bad Params", "\u001B[0;1;48;fivemGrey93\u001B[0m", nil, true},
		{"Bad Params 2", "\u001B[0;1;48;3mGrey93\u001B[0m", nil, true},
		{"Bad Params 2", "\u001B[0;1;50;3mGrey93\u001B[0m", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			is2.Equal(err != nil, tt.wantErr)
			for index, w := range tt.want {
				is2.Equal(got[index].Label, w.Label)
				if w.FgCol != nil {
					is2.Equal(got[index].BgCol.Name, w.BgCol.Name)
				}
				is2.Equal(got[index].Style, w.Style)
			}
		})
	}
}

func TestParseAnsiTrueColor(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name    string
		input   string
		want    []*StyledText
		wantErr bool
	}{
		{"Red", "\u001B[38;2;255;0;0mRed\u001B[0m", []*StyledText{
			{Label: "Red", FgCol: &Col{Rgb: Rgb{255, 0, 0}, Hex: "#ff0000"}},
		}, false},
		{"Red BG", "\u001B[48;2;255;0;0mRed\u001B[0m", []*StyledText{
			{Label: "Red", BgCol: &Col{Rgb: Rgb{255, 0, 0}, Hex: "#ff0000"}},
		}, false},
		{"Red, text, Green", "\u001B[38;2;255;0;0mRed\u001B[0mI am plain text\u001B[38;2;0;255;0mGreen\u001B[0m", []*StyledText{
			{Label: "Red", FgCol: &Col{Rgb: Rgb{255, 0, 0}, Hex: "#ff0000"}},
			{Label: "I am plain text"},
			{Label: "Green", FgCol: &Col{Rgb: Rgb{0, 255, 0}, Hex: "#00ff00"}},
		}, false},
		{"Bad 1", "\u001B[38;2;256;0;0mRed\u001B[0m", nil, true},
		{"Bad 2", "\u001B[38;2;-1;0;0mRed\u001B[0m", nil, true},
		{"Bad no of params", "\u001B[38;2;0;0mRed\u001B[0m", nil, true},
		{"Bad params", "\u001B[38;2;0;onemRed\u001B[0m", nil, true},
		{"Bad params 2", "\u001B[38;2;red;0;0mRed\u001B[0m", nil, true},
		{"Bad params 3", "\u001B[38;2;0;red;0mRed\u001B[0m", nil, true},
		{"Bad params 4", "\u001B[38;2;0;0;redmRed\u001B[0m", nil, true},
		{"Bad params 4", "\u001B[38;3;0;0;redmRed\u001B[0m", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			is2.Equal(err != nil, tt.wantErr)
			for index, w := range tt.want {
				is2.Equal(got[index].Label, w.Label)
				if w.FgCol != nil {
					is2.Equal(got[index].FgCol.Hex, w.FgCol.Hex)
					is2.Equal(got[index].FgCol.Rgb, w.FgCol.Rgb)
				}
				if w.BgCol != nil {
					is2.Equal(got[index].BgCol.Hex, w.BgCol.Hex)
					is2.Equal(got[index].BgCol.Rgb, w.BgCol.Rgb)
				}
				is2.Equal(got[index].Style, w.Style)
			}
		})
	}
}

func TestParseAnsiWithOptions(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name    string
		input   string
		options []ParseOption
		want    []*StyledText
		wantErr bool
	}{
		{
			"Unexpected code errored", "\u001b[0;99mHello World\033[0m",
			nil, nil, true,
		},
		{
			"Unexpected code ignored", "\u001b[0;99mHello World\033[0m",
			[]ParseOption{WithIgnoreInvalidCodes()},
			[]*StyledText{{Label: "Hello World"}}, false,
		},
		{
			"Foreground code default", "\u001b[0;39mHello World\033[0m", nil,
			[]*StyledText{{Label: "Hello World", FgCol: Cols[7]}}, false,
		},
		{
			"Foreground code specified", "\u001b[0;39mHello World\033[0m",
			[]ParseOption{WithDefaultForegroundColor("35")},
			[]*StyledText{{Label: "Hello World", FgCol: Cols[5]}}, false,
		},
		{
			"Background code default", "\u001b[0;49mHello World\033[0m", nil,
			[]*StyledText{{Label: "Hello World", BgCol: Cols[0]}}, false,
		},
		{
			"Background code specified", "\u001b[0;49mHello World\033[0m",
			[]ParseOption{WithDefaultBackgroundColor("36")},
			[]*StyledText{{Label: "Hello World", BgCol: Cols[6]}}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input, tt.options...)
			is2.Equal(err != nil, tt.wantErr)
			for index, w := range tt.want {
				is2.Equal(got[index].Label, w.Label)
				if w.FgCol != nil {
					is2.Equal(got[index].FgCol.Name, w.FgCol.Name)
				}
				if w.BgCol != nil {
					is2.Equal(got[index].BgCol.Name, w.BgCol.Name)
				}
				is2.Equal(got[index].Style, w.Style)
			}
		})
	}
}

func TestHasEscapeCodes(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"yes", "\u001B[0;30mHello World\033[0m\u001B[0;31mHello World\u001B[0m", true},
		{"no", "This is great!", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasEscapeCodes(tt.input)
			is2.Equal(got, tt.want)
		})
	}
}

func TestString(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name  string
		input []*StyledText
		want  string
	}{
		{"Blank", []*StyledText{}, ""},
		{"ANSI16 Fg", []*StyledText{{Label: "Red", FgCol: Cols[1]}}, "\033[0;31mRed\033[0m"},
		{"ANSI16 Fg Bold", []*StyledText{{Label: "Red", FgCol: Cols[1], Style: Bold}}, "\033[0;1;31mRed\033[0m"},
		{"ANSI16 Fg Strikethrough", []*StyledText{{Label: "Red", FgCol: Cols[1], Style: Strikethrough}}, "\033[0;9;31mRed\033[0m"},
		{"ANSI16 Fg Bold & Italic", []*StyledText{{Label: "Red", FgCol: Cols[1], Style: Bold | Italic}}, "\033[0;1;3;31mRed\033[0m"},
		{"ANSI16 Bg", []*StyledText{{Label: "Black", BgCol: Cols[0]}}, "\033[0;40mBlack\033[0m"},
		{"ANSI16 Mixed", []*StyledText{{Label: "Mixed", FgCol: Cols[1], BgCol: Cols[0]}}, "\033[0;31;40mMixed\033[0m"},
		{"ANSI256 Fg", []*StyledText{{Label: "Dark Blue", FgCol: Cols[18]}}, "\033[0;38;5;18mDark Blue\033[0m"},
		{"ANSI256 Fg Bold", []*StyledText{{Label: "Dark Blue", FgCol: Cols[18], Style: Bold}}, "\033[0;1;38;5;18mDark Blue\033[0m"},
		{"ANSI256 Bg", []*StyledText{{Label: "Dark Blue", BgCol: Cols[18]}}, "\033[0;48;5;18mDark Blue\033[0m"},
		{"ANSI256 Bg Bold", []*StyledText{{Label: "Dark Blue", BgCol: Cols[18], Style: Bold}}, "\033[0;1;48;5;18mDark Blue\033[0m"},
		{"Truecolor Fg", []*StyledText{{Label: "TrueColor!", FgCol: &Col{Id: 256, Rgb: Rgb{R: 128, G: 127, B: 126}}}}, "\033[0;38;2;128;127;126mTrueColor!\033[0m"},
		{"Truecolor Fg Bold", []*StyledText{{Label: "TrueColor!", FgCol: &Col{Id: 256, Rgb: Rgb{R: 128, G: 127, B: 126}}, Style: Bold}}, "\033[0;1;38;2;128;127;126mTrueColor!\033[0m"},
		{"Truecolor Bg", []*StyledText{{Label: "TrueColor!", BgCol: &Col{Id: 256, Rgb: Rgb{R: 128, G: 127, B: 126}}}}, "\033[0;48;2;128;127;126mTrueColor!\033[0m"},
		{"Truecolor Bg Bold", []*StyledText{{Label: "TrueColor!", BgCol: &Col{Id: 256, Rgb: Rgb{R: 128, G: 127, B: 126}}, Style: Bold}}, "\033[0;1;48;2;128;127;126mTrueColor!\033[0m"},
		{"Truecolor Mixed", []*StyledText{{Label: "TrueColor!", FgCol: &Col{Id: 256, Rgb: Rgb{R: 90, G: 91, B: 92}}, BgCol: &Col{Id: 256, Rgb: Rgb{R: 128, G: 127, B: 126}}, Style: Bold | Faint | Underlined | Strikethrough | Italic | Invisible | Blinking | Inversed}}, "\033[0;1;2;3;4;5;7;8;9;38;2;90;91;92;48;2;128;127;126mTrueColor!\033[0m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := String(tt.input)
			is2.Equal(got, tt.want)
		})
	}
}

func TestTruncate(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name     string
		input    string
		maxChars int
		want     string
		wantErr  bool
	}{
		{"No formatting", "Hello World", 11, "Hello World", false},
		{"No formatting truncated", "Hello World", 5, "Hello", false},
		{"No formatting many chars", "Hello World", 50, "Hello World", false},
		{"Black", "\u001b[0;30mHello World\033[0m", 5, "\u001B[0;30mHello\u001B[0m", false},
		{"Red Bold", "\u001b[0;1;31mHello World\033[0m", 5, "\033[0;1;31mHello\033[0m", false},
		{"Red Bold & Black", "\u001b[0;1;31mI am Red\033[0m\u001B[0;30mI am Black\u001B[0m", 12, "\u001B[0;1;31mI am Red\u001B[0m\u001B[0;30mI am\u001B[0m", false},
		{"Red Bold & text & Black", "\u001b[0;1;31mI am Red\033[0m and \u001B[0;30mI am Black\u001B[0m", 17, "\u001B[0;1;31mI am Red\u001B[0m and \u001B[0;30mI am\u001B[0m", false},
		{"Emoji", "\u001B[0;1;31m😀👩🏽‍🔧\u001B[0m", 1, "\u001B[0;1;31m😀\u001B[0m", false},
		{"Emoji 2", "\u001B[0;1;31m😀👩🏽‍🔧\u001B[0m", 2, "\u001B[0;1;31m😀👩🏽‍🔧\u001B[0m", false},
		{"Emoji 3", "\u001B[0;1;31m😀👩🏽‍🔧😀\u001B[0m", 2, "\u001B[0;1;31m😀👩🏽‍🔧\u001B[0m", false},
		{"Bad", "\033[44;32;12", 10, "", true},

		//{"Spaces", "  ", "  ", "", false},
		//{"Bad code", "\u001b[1  ", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Truncate(tt.input, tt.maxChars)
			is2.Equal(err != nil, tt.wantErr)
			is2.Equal(got, tt.want)
		})
	}
}

func TestCleanse(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"Blank", "", "", false},
		{"No formatting", "Hello World", "Hello World", false},
		{"ANSI16 Fg", "\033[0;31mRed\033[0m", "Red", false},
		{"Black", "\u001b[0;30mHello World\033[0m", "Hello World", false},
		{"Red Bold & Black", "\u001b[0;1;31mI am Red\033[0m & \u001B[0;30mI am Black\u001B[0m", "I am Red & I am Black", false},
		{"Red All the styles", "\u001b[0;1;2;3;4;5;7;8;9;31mI am Red\033[0m", "I am Red", false},
		{"Emoji", "\u001B[0;1;31m😀\u001B[0m", "😀", false},
		{"Emoji 2", "\u001B[0;1;31m😀👩🏽‍🔧\u001B[0m", "😀👩🏽‍🔧", false},
		{"Bad", "\033[44;32;12", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Cleanse(tt.input)
			is2.Equal(err != nil, tt.wantErr)
			is2.Equal(got, tt.want)
		})
	}
}

func TestLength(t *testing.T) {
	is2 := is.New(t)
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"Blank", "", 0, false},
		{"No formatting", "Hello World", 11, false},
		{"ANSI16 Fg", "\033[0;31mRed\033[0m", 3, false},
		{"Zero chats", "\u001b[0;30m\033[0m", 0, false},
		{"Red Bold & Black", "\u001b[0;1;31mI am Red\033[0m & \u001B[0;30mI am Black\u001B[0m", 21, false},
		{"Red All the styles", "\u001b[0;1;2;3;4;5;7;8;9;31mI am Red\033[0m", 8, false},
		{"Emoji", "\u001B[0;1;31m😀\u001B[0m", 1, false},
		{"Emoji 2", "\u001B[0;1;31m😀👩🏽‍🔧\u001B[0m", 2, false},
		{"Bad", "\033[44;32;12", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Length(tt.input)
			is2.Equal(err != nil, tt.wantErr)
			is2.Equal(got, tt.want)
		})
	}
}
