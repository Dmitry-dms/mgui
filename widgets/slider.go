package widgets

import (
	"github.com/Dmitry-dms/mgui/styles"
	"github.com/Dmitry-dms/mgui/utils"
)

type Slider struct {
	base                   baseWidget
	min, max               float32
	sliderHeight, btnWidth float32
	mainSliderPos          [4]float32
	btnSliderPos           [4]float32
	CurrentPos             float32
	btnCaptured            bool
}

func NewSlider(id string, x, y, w, h, min, max float32, style *styles.Style) *Slider {
	s := Slider{
		base: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h},
			BackgroundColor: [4]float32{123, 32, 12, 0},
		},
		min:          min,
		max:          max,
		sliderHeight: style.SliderHeight,
		btnWidth:     style.SliderBtnWidth,
	}
	return &s
}

var dragSpeed float32 = 3

func (s *Slider) HandleMouseDrag(delta float32, i *float32, f utils.CapturedDragBehavior) {
	f(utils.NewRectS(s.btnSliderPos), &s.btnCaptured)
	if s.btnCaptured {
		delta *= dragSpeed
		if *i <= s.min {
			*i = s.min
			if delta > 0 {
				*i += delta
			}
		} else if *i > s.max {
			*i = s.max
			if delta < 0 {
				*i -= delta
			}
		} else {
			*i += delta
		}
	}
}

func (s *Slider) CalculateNumber(i *float32) {
	baseLength := s.Width() - s.btnSliderPos[2]
	if *i <= s.min {
		s.CurrentPos = 0
	} else if *i >= s.max {
		s.CurrentPos = baseLength
	} else {
		s.CurrentPos = (*i - s.min) / (s.max - s.min) * baseLength
	}
}

func (s *Slider) calculateSliderPos() {
	v := (s.Height() - s.sliderHeight) / 2
	s.mainSliderPos = [4]float32{s.base.boundingBox[0], s.base.boundingBox[1] + v, s.Width(), s.sliderHeight}
	s.btnSliderPos = [4]float32{s.base.boundingBox[0] + s.CurrentPos, s.base.boundingBox[1], s.btnWidth, s.Height()}
}

func (s *Slider) MainSliderPos() [4]float32 {
	return s.mainSliderPos
}
func (s *Slider) BtnSliderPos() [4]float32 {
	return s.btnSliderPos
}
func (s *Slider) BoundingBox() [4]float32 {
	return s.base.boundingBox
}

func (s *Slider) UpdatePosition(pos [4]float32) {
	s.base.boundingBox = pos
	s.calculateSliderPos()
}

func (s *Slider) WidgetId() string {
	return s.base.id
}

func (s *Slider) Height() float32 {
	return s.base.height()
}
func (s *Slider) Width() float32 {
	return s.base.width()
}
