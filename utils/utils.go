package utils

import (
	"math/rand"
	"strings"
	"time"
)

type Stack[T any] struct {
	Push   func(T)
	Pop    func() T
	Peek   func() T
	Length func() int
}

func NewStack[T any]() Stack[T] {
	slice := make([]T, 0)
	return Stack[T]{
		Push: func(i T) {
			slice = append(slice, i)
		},
		Pop: func() T {
			res := slice[len(slice)-1]
			slice = slice[:len(slice)-1]
			return res
		},
		Peek: func() T {
			return slice[len(slice)-1]
		},
		Length: func() int {
			return len(slice)
		},
	}
}

type Vec2 struct {
	X, Y float32
}

func (v Vec2) Sub(x Vec2) Vec2 {
	return Vec2{v.X - x.X, v.Y - x.Y}
}
func (v Vec2) Add(x Vec2) Vec2 {
	return Vec2{v.X + x.X, v.Y + x.Y}
}

func (v Vec2) LengthSqr() float32 {
	return v.X*v.X + v.Y*v.Y
}

type Vec4 struct {
	X, Y, W, H float32
}

func Max[T int | float32](x1, x2 T) T {
	if x1 > x2 {
		return x1
	} else {
		return x2
	}
}

type Rect struct {
	Min, Max Vec2
}

func NewRect(x, y, w, h float32) Rect {
	min := Vec2{x, y}
	max := Vec2{x + w, y + h}
	return Rect{min, max}
}

func NewRectS(r [4]float32) Rect {
	min := Vec2{r[0], r[1]}
	max := Vec2{r[0] + r[2], r[1] + r[3]}
	return Rect{min, max}
}

func (r Rect) Width() float32 {
	return r.Max.X - r.Min.X
}
func (r *Rect) AddWidth(w float32) {
	prev := r.Max
	r.Max = Vec2{prev.X + w, prev.Y}
}

func (r *Rect) AddHeight(h float32) {
	r.Max = Vec2{r.Max.X, r.Max.Y + h}
}

func (r Rect) Height() float32 {
	return r.Max.Y - r.Min.Y
}

func (r Rect) Contains(p Vec2) bool {
	return p.X >= r.Min.X && p.Y >= r.Min.Y && p.X < r.Max.X && p.Y < r.Max.Y
}

func PointOutsideRect(point Vec2, box Rect) bool {
	if point.X == 0 && point.Y == 0 {
		return false
	}
	return !PointInRect(point, box)
}

func PointInRect(point Vec2, box Rect) bool {
	min := box.Min
	max := box.Max

	return point.X <= max.X && min.X <= point.X &&
		point.Y <= max.Y && min.Y <= point.Y
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type CapturedDragBehavior func(r Rect, captured *bool)

var src = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
