package styles

type Style struct {

	//misc
	TransparentColor [4]float32

	// widget space
	Margin                float32
	Padding               float32
	TopMargin, BotMargin  float32
	LeftMargin            float32
	WidgSpaceDividerColor [4]float32

	//btn
	BtnColor        [4]float32
	BtnHoveredColor [4]float32
	BtnActiveColor  [4]float32

	//tab
	TabBtnColor       [4]float32
	TabBtnActiveColor [4]float32

	//text
	TextColor   [4]float32
	TextSize    int
	TextPadding int
	FontScale   float32

	//slider
	SliderColor    [4]float32
	SliderBtnColor [4]float32
	SliderHeight   float32
	SliderBtnWidth float32
}

var (
	DefaultStyle = Style{
		TransparentColor: [4]float32{0, 0, 0, 0},
		Margin:           10,
		Padding:          10,
		TopMargin:        10,
		BotMargin:        10,
		LeftMargin:       10,
		BtnColor:         [4]float32{124, 90, 156, 1},
		BtnHoveredColor:  [4]float32{200, 270, 30, 1},
		BtnActiveColor:   [4]float32{255, 0, 0, 1},
		TextColor:        [4]float32{5, 134, 222, 1},
		TextSize:         20,
		//TextPadding:      10,
		SliderColor:           [4]float32{231, 240, 162, 0.8},
		SliderBtnColor:        [4]float32{0, 0, 0, 1},
		SliderHeight:          10,
		SliderBtnWidth:        20,
		FontScale:             1,
		TabBtnColor:           [4]float32{44, 81, 132, 1},
		TabBtnActiveColor:     [4]float32{51, 105, 173, 1},
		WidgSpaceDividerColor: [4]float32{70, 70, 79, 1},
	}
)
