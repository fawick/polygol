package polygol

import "testing"

func Test_crossProduct(t *testing.T) {
	pt1 := Vector{x: newBigNumber(1), y: newBigNumber(2)}
	pt2 := Vector{x: newBigNumber(3), y: newBigNumber(4)}
	expect(t, crossProduct(pt1, pt2).equalTo(newBigNumber(-2)))
}

func Test_dotProduct(t *testing.T) {
	pt1 := Vector{x: newBigNumber(1), y: newBigNumber(2)}
	pt2 := Vector{x: newBigNumber(3), y: newBigNumber(4)}
	expect(t, dotProduct(pt1, pt2).equalTo(newBigNumber(11)))
}

func Test_length(t *testing.T) {
	testCases := []struct {
		name string
		v    Vector
		l    BigNumber
	}{
		{
			name: "horizontal",
			v:    Vector{x: newBigNumber(3), y: newBigNumber(0)},
			l:    newBigNumber(3),
		},
		{
			name: "vertical",
			v:    Vector{x: newBigNumber(0), y: newBigNumber(-2)},
			l:    newBigNumber(2),
		},
		{
			name: "3-4-5",
			v:    Vector{x: newBigNumber(3), y: newBigNumber(4)},
			l:    newBigNumber(5),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			expect(t, length(tt.v).equalTo(tt.l))

		})
	}
}

func Test_sineAndCosineOfAngle(t *testing.T) {
	testCases := []struct {
		name    string
		shared  Vector
		base    Vector
		angle   Vector
		sine    BigNumber
		cosine  BigNumber
		closeTo bool
	}{
		{
			name:   "parallel",
			shared: Vector{x: newBigNumber(0), y: newBigNumber(0)},
			base:   Vector{x: newBigNumber(1), y: newBigNumber(0)},
			angle:  Vector{x: newBigNumber(1), y: newBigNumber(0)},
			sine:   newBigNumber(0),
			cosine: newBigNumber(1),
		},
		{
			name:    "45 degrees",
			shared:  Vector{x: newBigNumber(0), y: newBigNumber(0)},
			base:    Vector{x: newBigNumber(1), y: newBigNumber(0)},
			angle:   Vector{x: newBigNumber(1), y: newBigNumber(-1)},
			sine:    newBigNumber(2).sqrt().div(newBigNumber(2)),
			cosine:  newBigNumber(2).sqrt().div(newBigNumber(2)),
			closeTo: true,
		},
		{
			name:   "90 degrees",
			shared: Vector{x: newBigNumber(0), y: newBigNumber(0)},
			base:   Vector{x: newBigNumber(1), y: newBigNumber(0)},
			angle:  Vector{x: newBigNumber(0), y: newBigNumber(-1)},
			sine:   newBigNumber(1),
			cosine: newBigNumber(0),
		},
		{
			name:    "135 degrees",
			shared:  Vector{x: newBigNumber(0), y: newBigNumber(0)},
			base:    Vector{x: newBigNumber(1), y: newBigNumber(0)},
			angle:   Vector{x: newBigNumber(-1), y: newBigNumber(-1)},
			sine:    newBigNumber(2).sqrt().div(newBigNumber(2)),
			cosine:  newBigNumber(2).sqrt().negated().div(newBigNumber(2)),
			closeTo: true,
		},
		{
			name:   "anti-parallel",
			shared: Vector{x: newBigNumber(0), y: newBigNumber(0)},
			base:   Vector{x: newBigNumber(1), y: newBigNumber(0)},
			angle:  Vector{x: newBigNumber(-1), y: newBigNumber(0)},
			sine:   newBigNumber(0),
			cosine: newBigNumber(-1),
		},
		{
			name:    "225 degrees",
			shared:  Vector{x: newBigNumber(0), y: newBigNumber(0)},
			base:    Vector{x: newBigNumber(1), y: newBigNumber(0)},
			angle:   Vector{x: newBigNumber(-1), y: newBigNumber(1)},
			sine:    newBigNumber(2).sqrt().negated().div(newBigNumber(2)),
			cosine:  newBigNumber(2).sqrt().negated().div(newBigNumber(2)),
			closeTo: true,
		},
		{
			name:   "270 degrees",
			shared: Vector{x: newBigNumber(0), y: newBigNumber(0)},
			base:   Vector{x: newBigNumber(1), y: newBigNumber(0)},
			angle:  Vector{x: newBigNumber(0), y: newBigNumber(1)},
			sine:   newBigNumber(-1),
			cosine: newBigNumber(0),
		},
		{
			name:    "315 degrees",
			shared:  Vector{x: newBigNumber(0), y: newBigNumber(0)},
			base:    Vector{x: newBigNumber(1), y: newBigNumber(0)},
			angle:   Vector{x: newBigNumber(1), y: newBigNumber(1)},
			sine:    newBigNumber(2).sqrt().negated().div(newBigNumber(2)),
			cosine:  newBigNumber(2).sqrt().div(newBigNumber(2)),
			closeTo: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("sine", func(t *testing.T) {
				if tt.closeTo {
					expect(t, sineOfAngle(tt.shared, tt.base, tt.angle).closeTo(tt.sine))
				} else {
					expect(t, sineOfAngle(tt.shared, tt.base, tt.angle).equalTo(tt.sine))
				}
			})
			t.Run("cosine", func(t *testing.T) {
				if tt.closeTo {
					expect(t, cosineOfAngle(tt.shared, tt.base, tt.angle).closeTo(tt.cosine))
				} else {
					expect(t, cosineOfAngle(tt.shared, tt.base, tt.angle).equalTo(tt.cosine))
				}
			})
		})
	}
}

func Test_perpendicular(t *testing.T) {
	testCases := []struct {
		name string
		v    Vector
	}{
		{
			name: "vertical",
			v:    Vector{x: newBigNumber(0), y: newBigNumber(1)},
		},
		{
			name: "horizontal",
			v:    Vector{x: newBigNumber(1), y: newBigNumber(0)},
		},
		{
			name: "45 degrees",
			v:    Vector{x: newBigNumber(1), y: newBigNumber(1)},
		},
		{
			name: "120 degrees",
			v:    Vector{x: newBigNumber(-1), y: newBigNumber(2)},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := perpendicular(tt.v)
			expect(t, dotProduct(tt.v, r).equalTo(newBigNumber(0)))
			expect(t, crossProduct(tt.v, r).notEqualTo(newBigNumber(0)))
		})
	}
}

func Test_verticalIntersection(t *testing.T) {
	t.Run("horizontal", func(t *testing.T) {
		p := Vector{x: newBigNumber(42), y: newBigNumber(3)}
		v := Vector{x: newBigNumber(-2), y: newBigNumber(0)}
		x := newBigNumber(37)
		i := verticalIntersection(p, v, x)
		expect(t, i.x.equalTo(newBigNumber(37)))
		expect(t, i.y.equalTo(newBigNumber(3)))
	})
	t.Run("vertical", func(t *testing.T) {
		p := Vector{x: newBigNumber(42), y: newBigNumber(3)}
		v := Vector{x: newBigNumber(0), y: newBigNumber(4)}
		x := newBigNumber(-2)
		i := verticalIntersection(p, v, x)
		if i != nil {
			t.FailNow()
		}
	})
	t.Run("45 degrees", func(t *testing.T) {
		p := Vector{x: newBigNumber(1), y: newBigNumber(1)}
		v := Vector{x: newBigNumber(1), y: newBigNumber(1)}
		x := newBigNumber(-2)
		i := verticalIntersection(p, v, x)
		expect(t, i.x.equalTo(newBigNumber(-2)))
		expect(t, i.y.equalTo(newBigNumber(-2)))
	})
	t.Run("upper left quadrant", func(t *testing.T) {
		p := Vector{x: newBigNumber(-1), y: newBigNumber(1)}
		v := Vector{x: newBigNumber(-2), y: newBigNumber(1)}
		x := newBigNumber(-3)
		i := verticalIntersection(p, v, x)
		expect(t, i.x.equalTo(newBigNumber(-3)))
		expect(t, i.y.equalTo(newBigNumber(2)))
	})
}

func Test_horizontalIntersection(t *testing.T) {
	t.Run("horizontal", func(t *testing.T) {
		p := Vector{x: newBigNumber(42), y: newBigNumber(3)}
		v := Vector{x: newBigNumber(-2), y: newBigNumber(0)}
		x := newBigNumber(37)
		i := horizontalIntersection(p, v, x)
		if i != nil {
			t.FailNow()
		}
	})
	t.Run("vertical", func(t *testing.T) {
		p := Vector{x: newBigNumber(42), y: newBigNumber(3)}
		v := Vector{x: newBigNumber(0), y: newBigNumber(4)}
		x := newBigNumber(37)
		i := horizontalIntersection(p, v, x)
		expect(t, i.x.equalTo(newBigNumber(42)))
		expect(t, i.y.equalTo(newBigNumber(37)))
	})
	t.Run("45 degrees", func(t *testing.T) {
		p := Vector{x: newBigNumber(1), y: newBigNumber(1)}
		v := Vector{x: newBigNumber(1), y: newBigNumber(1)}
		x := newBigNumber(4)
		i := horizontalIntersection(p, v, x)
		expect(t, i.x.equalTo(newBigNumber(4)))
		expect(t, i.y.equalTo(newBigNumber(4)))
	})
	t.Run("bottom left quadrant", func(t *testing.T) {
		p := Vector{x: newBigNumber(-1), y: newBigNumber(-1)}
		v := Vector{x: newBigNumber(-2), y: newBigNumber(-1)}
		x := newBigNumber(-3)
		i := horizontalIntersection(p, v, x)
		expect(t, i.x.equalTo(newBigNumber(-5)))
		expect(t, i.y.equalTo(newBigNumber(-3)))
	})
}

func Test_intersection(t *testing.T) {
	p1 := Vector{x: newBigNumber(42), y: newBigNumber(42)}
	p2 := Vector{x: newBigNumber(-32), y: newBigNumber(46)}

	testCases := []struct {
		name    string
		v1      Vector
		v2      Vector
		invalid bool
		x       BigNumber
		y       BigNumber
	}{
		{
			name:    "parallel",
			v1:      Vector{x: newBigNumber(1), y: newBigNumber(2)},
			v2:      Vector{x: newBigNumber(-1), y: newBigNumber(-2)},
			invalid: true,
		},
		{
			name: "horizontal and vertical",
			v1:   Vector{x: newBigNumber(0), y: newBigNumber(2)},
			v2:   Vector{x: newBigNumber(-1), y: newBigNumber(0)},
			x:    newBigNumber(42),
			y:    newBigNumber(46),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			v1 := Vector{x: newBigNumber(1), y: newBigNumber(2)}
			v2 := Vector{x: newBigNumber(-1), y: newBigNumber(-2)}
			i := intersection(p1, v1, p2, v2)
			if tt.invalid && i != nil {
				t.FailNow()
				expect(t, i.x.equalTo(tt.x))
				expect(t, i.y.equalTo(tt.y))
			}
		})
	}
	t.Run("consistency", func(t *testing.T) {
		p1 := Vector{x: newBigNumber(0.523787), y: newBigNumber(51.281453)}
		v1 := Vector{x: newBigNumber(0.0002729999999999677), y: newBigNumber(0.0002729999999999677)}
		p2 := Vector{x: newBigNumber(0.523985), y: newBigNumber(51.281651)}
		v2 := Vector{x: newBigNumber(0.000024999999999941735), y: newBigNumber(0.000049000000004184585)}
		i1 := intersection(p1, v1, p2, v2)
		i2 := intersection(p2, v2, p1, v1)
		expect(t, i1.x.equalTo(i2.x))
		expect(t, i1.y.equalTo(i2.y))
	})
}
