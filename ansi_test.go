package ansi

import (
	"testing"

	"github.com/matryer/is"
)

func TestParseAnsi16Styles(t *testing.T) {
	is := is.New(t)
	var got []*StyledText
	var err error

	// Bold
	got, err = Parse("\u001b[1;30mHello World\033[0m")
	is.NoErr(err)
	is.Equal(len(got), 1)
	is.True(got[0].Bold())
	// Faint
	got, err = Parse("\u001b[2;30mHello World\033[0m")
	is.NoErr(err)
	is.Equal(len(got), 1)
	is.True(got[0].Faint())
	// Italic
	got, err = Parse("\u001b[3;30mHello World\033[0m")
	is.NoErr(err)
	is.Equal(len(got), 1)
	is.True(got[0].Italic())
	// Underlined
	got, err = Parse("\u001b[4;30mHello World\033[0m")
	is.NoErr(err)
	is.Equal(len(got), 1)
	is.True(got[0].Underlined())
	// Blinking
	got, err = Parse("\u001b[5;30mHello World\033[0m")
	is.NoErr(err)
	is.Equal(len(got), 1)
	is.True(got[0].Blinking())
	// Inversed
	got, err = Parse("\u001b[7;30mHello World\033[0m")
	is.NoErr(err)
	is.Equal(len(got), 1)
	is.True(got[0].Inversed())
	// Invisible
	got, err = Parse("\u001b[8;30mHello World\033[0m")
	is.NoErr(err)
	is.Equal(len(got), 1)
	is.True(got[0].Invisible())
	// Strikethrough
	got, err = Parse("\u001b[9;30mHello World\033[0m")
	is.NoErr(err)
	is.Equal(len(got), 1)
	is.True(got[0].Strikethrough())
}

func TestParseAnsi16Swap(t *testing.T) {
	is := is.New(t)
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
	is.NoErr(err)
	is.Equal(len(got), 1)
	is.Equal(got[0].FgCol.Name, "Coffee")

	// Restore
	ColourMap["Regular"]["30"] = original
	got, err = Parse("\u001b[0;30mHello World\033[0m")
	is.NoErr(err)
	is.Equal(len(got), 1)
	is.Equal(got[0].FgCol.Name, "Black")
}

func TestParseAnsi16SingleColour(t *testing.T) {
	is := is.New(t)
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
			is.Equal(err != nil, tt.wantErr)
			expectedLength := 1
			if tt.wantErr {
				expectedLength = 0
			}
			is.Equal(len(got), expectedLength)
			if expectedLength == 1 {
				if len(tt.wantColor) > 0 {
					is.True(got[0].FgCol != nil)
					is.Equal(got[0].FgCol.Name, tt.wantColor)
				}
			}
		})
	}
}

func TestParseAnsi16SingleBGColour(t *testing.T) {
	is := is.New(t)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			is.Equal(err != nil, tt.wantErr)
			expectedLength := 1
			if tt.wantErr {
				expectedLength = 0
			}
			is.Equal(len(got), expectedLength)
			if expectedLength == 1 {
				if len(tt.wantColor) > 0 {
					is.True(got[0].BgCol != nil)
					is.Equal(got[0].BgCol.Name, tt.wantColor)
				}
			}
		})
	}
}

func TestParseAnsi16MultiColour(t *testing.T) {
	is := is.New(t)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			is.Equal(err != nil, tt.wantErr)
			for index, w := range tt.want {
				is.Equal(got[index].Label, w.Label)
				if w.FgCol != nil {
					is.Equal(got[index].FgCol.Name, w.FgCol.Name)
				}
				is.Equal(got[index].Style, w.Style)
			}
		})
	}
}

func TestParseAnsi256(t *testing.T) {
	is := is.New(t)
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
			is.Equal(err != nil, tt.wantErr)
			for index, w := range tt.want {
				is.Equal(got[index].Label, w.Label)
				if w.FgCol != nil {
					is.Equal(got[index].FgCol.Name, w.FgCol.Name)
				}
				is.Equal(got[index].Style, w.Style)
			}
		})
	}
}

func TestParseAnsiBG256(t *testing.T) {
	is := is.New(t)
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
		{"Bad Params 2", "\u001B[0;1;49;3mGrey93\u001B[0m", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			is.Equal(err != nil, tt.wantErr)
			for index, w := range tt.want {
				is.Equal(got[index].Label, w.Label)
				if w.FgCol != nil {
					is.Equal(got[index].BgCol.Name, w.BgCol.Name)
				}
				is.Equal(got[index].Style, w.Style)
			}
		})
	}
}

func TestParseAnsiTrueColor(t *testing.T) {
	is := is.New(t)
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
			is.Equal(err != nil, tt.wantErr)
			for index, w := range tt.want {
				is.Equal(got[index].Label, w.Label)
				if w.FgCol != nil {
					is.Equal(got[index].FgCol.Hex, w.FgCol.Hex)
					is.Equal(got[index].FgCol.Rgb, w.FgCol.Rgb)
				}
				if w.BgCol != nil {
					is.Equal(got[index].BgCol.Hex, w.BgCol.Hex)
					is.Equal(got[index].BgCol.Rgb, w.BgCol.Rgb)
				}
				is.Equal(got[index].Style, w.Style)
			}
		})
	}
}
