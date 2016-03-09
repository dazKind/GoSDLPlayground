package lmath

import (
	"math" 
)

type Mat44 [16] float32

func Mat44PerspLH(_fov float32, _aspect float32, _nz float32, _fz float32) Mat44 {
	ymax := _nz * float32(math.Atan(float64(HALF_DEG2RAD * _fov)))
	xmax := ymax * _aspect
	return Mat44PerspOffCenter(xmax, ymax, _nz, _fz)
}

func Mat44PerspOffCenter(_maxX float32, _maxY float32, _nz float32, _fz float32) Mat44 {
	m := Mat44{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	m[0] = 2 * _nz / (_maxX * 2)
	m[5] = 2 * _nz / (_maxY * 2)
	m[8] = 0
	m[9] = 0
	m[10] = -(_fz + _nz) / (_fz - _nz)
	m[11] = -1
	m[14] = -2 * _fz * _nz / (_fz - _nz)
	m[15] = 0
	return m
}

func (_m0 Mat44) Mult(_m1 Mat44) Mat44 {
	return Mat44 {
		_m0[0] * _m1[0] + _m0[4] * _m1[1] + _m0[8] * _m1[2] + _m0[12] * _m1[3],
        _m0[1] * _m1[0] + _m0[5] * _m1[1] + _m0[9] * _m1[2] + _m0[13] * _m1[3],
        _m0[2] * _m1[0] + _m0[6] * _m1[1] + _m0[10] * _m1[2] + _m0[14] * _m1[3],
        _m0[3] * _m1[0] + _m0[7] * _m1[1] + _m0[11] * _m1[2] + _m0[15] * _m1[3],

        _m0[0] * _m1[4] + _m0[4] * _m1[5] + _m0[8] * _m1[6] + _m0[12] * _m1[7],
        _m0[1] * _m1[4] + _m0[5] * _m1[5] + _m0[9] * _m1[6] + _m0[13] * _m1[7],
        _m0[2] * _m1[4] + _m0[6] * _m1[5] + _m0[10] * _m1[6] + _m0[14] * _m1[7],
        _m0[3] * _m1[4] + _m0[7] * _m1[5] + _m0[11] * _m1[6] + _m0[15] * _m1[7],

        _m0[0] * _m1[8] + _m0[4] * _m1[9] + _m0[8] * _m1[10] + _m0[12] * _m1[11],
        _m0[1] * _m1[8] + _m0[5] * _m1[9] + _m0[9] * _m1[10] + _m0[13] * _m1[11],
        _m0[2] * _m1[8] + _m0[6] * _m1[9] + _m0[10] * _m1[10] + _m0[14] * _m1[11],
        _m0[3] * _m1[8] + _m0[7] * _m1[9] + _m0[11] * _m1[10] + _m0[15] * _m1[11],

        _m0[0] * _m1[12] + _m0[4] * _m1[13] + _m0[8] * _m1[14] + _m0[12] * _m1[15],
        _m0[1] * _m1[12] + _m0[5] * _m1[13] + _m0[9] * _m1[14] + _m0[13] * _m1[15],
        _m0[2] * _m1[12] + _m0[6] * _m1[13] + _m0[10] * _m1[14] + _m0[14] * _m1[15],
        _m0[3] * _m1[12] + _m0[7] * _m1[13] + _m0[11] * _m1[14] + _m0[15] * _m1[15],
	}
}

func (_m *Mat44) Recompose(_o Quat, _s Vec3, _t Vec3) {
	rot := Mat44{ 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1 }
	rot.SetOrientation(_o);
	
	scale := Mat44{ 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1 }
	scale.SetScale(_s.X, _s.Y, _s.Z)
	
	tmp := rot.Mult(scale);

	_m[0] = tmp[0];
	_m[1] = tmp[1];
	_m[2] = tmp[2];
	_m[3] = tmp[3];
	_m[4] = tmp[4];
	_m[5] = tmp[5];
	_m[6] = tmp[6];
	_m[7] = tmp[7];
	_m[8] = tmp[8];
	_m[9] = tmp[9];
	_m[10] = tmp[10];
	_m[11] = tmp[11];
	_m[12] = _t.X;
	_m[13] = _t.Y;
	_m[14] = _t.Z;
	_m[15] = 1;
}

func (_m *Mat44) SetOrientation(_q Quat) {
	Tx := 2 * _q.X;
    Ty := 2 * _q.Y;
    Tz := 2 * _q.Z;
    
    Twx := Tx * _q.W;
    Twy := Ty * _q.W;
    Twz := Tz * _q.W;
    
    Txx := Tx * _q.X;
    Txy := Ty * _q.X;
    Txz := Tz * _q.X;
    
    Tyy := Ty * _q.Y;
    Tyz := Tz * _q.Y;
    Tzz := Tz * _q.Z;
    
    _m[0] = 1 - (Tyy + Tzz);
    _m[1] = Txy - Twz;
    _m[2] = Txz + Twy;
    
    _m[4] = Txy + Twz;
    _m[5] = 1 - (Txx + Tzz);
    _m[6] = Tyz - Twx;
    
    _m[8] = Txz - Twy;
    _m[9] = Tyz + Twx;
    _m[10] = 1 - (Txx + Tyy);
}

func (_m *Mat44) SetScale(_x float32, _y float32, _z float32) {
	_m[0] = _x; _m[5] = _y; _m[10] = _z
}

func (_m *Mat44) SetTranslation(_x float32, _y float32, _z float32) {
	_m[12] = _x; _m[13] = _y; _m[14] = _z
}



