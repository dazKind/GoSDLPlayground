package lmath

import (
	"math"
)

type Vec3 struct { X,Y,Z float32 }

//
func (_v0 Vec3) Cross(_v1 Vec3) Vec3 {
	return Vec3{
		_v0.Y * _v1.Z - _v0.Z * _v1.Y,
		_v0.Z * _v1.X - _v0.X * _v1.Z,
        _v0.X * _v1.Y - _v0.Y * _v1.X,
	}
}

// 
func (_v0 Vec3) Dot(_v1 Vec3) float32 {
	return _v0.X*_v1.X + _v0.Y*_v1.Y + _v0.Z*_v1.Z
}

//
func (_v0 Vec3) LenSq() float32 {
	return _v0.Dot(_v0)
}

//
func (_v0 Vec3) Len() float32 {
	return float32(math.Sqrt(float64(_v0.LenSq())))
}

//
func (_v0 Vec3) Norm() float32 {
	len := _v0.Len()
	_v0.X /= len; _v0.Y /= len; _v0.Z /= len
	return len
}

//
func (_a Vec3) Add(_b Vec3) Vec3 {
	return Vec3{_a.X+_b.X, _a.Y+_b.Y, _a.Z+_b.Z}
}

//
func (_a Vec3) Subtract(_b Vec3) Vec3 {
	return Vec3{_a.X-_b.X, _a.Y-_b.Y, _a.Z-_b.Z}
}

//
func (_a Vec3) AlmostEq(_b Vec3, _eps float32) bool {
	return AlmostEqual(_a.X, _a.X, _eps) && AlmostEqual(_a.Y, _a.Y, _eps) && 
		AlmostEqual(_a.Z, _a.Z, _eps)
}

//
func (_a Vec3) Eq(_b Vec3) bool {
	return _a.AlmostEq(_b, EPSILON)
}
