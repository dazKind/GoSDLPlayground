package lmath

import "math"

type Quat struct { X,Y,Z,W float32 }

//
func (_q0 Quat) Cross(_q1 Quat) Quat {
	return Quat{
		_q0.Y * _q1.Z - _q0.Z * _q1.Y + _q1.X * _q0.W + _q0.X * _q1.W,
        _q0.Z * _q1.X - _q0.X * _q1.Z + _q1.Y * _q0.W + _q0.Y * _q1.W,
        _q0.X * _q1.Y - _q0.Y * _q1.X + _q1.Z * _q0.W + _q0.Z * _q1.W,
        _q0.W * _q1.W - (_q0.X * _q1.X + _q0.Y * _q1.Y + _q0.Z * _q1.Z),
	}
}

// 
func (_v0 Quat) Dot(_v1 Quat) float32 {
	return _v0.X * _v1.X + _v0.Y * _v1.Y + _v0.Z * _v1.Z + _v0.W * _v1.W
}

//
func (_v0 Quat) LenSq() float32 {
	return _v0.Dot(_v0)
}

//
func (_v0 Quat) Len() float32 {
	return float32(math.Sqrt(float64(_v0.LenSq())))
}

//
func (_v0 Quat) Norm() float32 {
	len := _v0.Len()
	_v0.X /= len; _v0.Y /= len; _v0.Z /= len; _v0.W /= len
	return len
}

func (_v0 *Quat) RotateInX(_a float32) {
	tmp := _v0.Cross(Quat{float32(math.Sin(float64(_a))), 0, 0, float32(math.Cos(float64(_a)))})
	_v0.X = tmp.X; _v0.Y = tmp.Y; _v0.Z = tmp.Z; _v0.W = tmp.W
}

func (_v0 *Quat) RotateInY(_a float32) {
	tmp := _v0.Cross(Quat{0, float32(math.Sin(float64(_a))), 0, float32(math.Cos(float64(_a)))})
	_v0.X = tmp.X; _v0.Y = tmp.Y; _v0.Z = tmp.Z; _v0.W = tmp.W
}

func (_v0 *Quat) RotateInZ(_a float32) {
	tmp := _v0.Cross(Quat{0, 0, float32(math.Sin(float64(_a))), float32(math.Cos(float64(_a)))})
	_v0.X = tmp.X; _v0.Y = tmp.Y; _v0.Z = tmp.Z; _v0.W = tmp.W
}
