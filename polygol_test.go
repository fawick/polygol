package polygol

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/engelsjk/polygol/geojson"
)

func expect(t testing.TB, what bool) {
	t.Helper()
	if !what {
		t.Error("expectation failure")
	}
}

func terr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func loadGeoms(filepath string) ([]Geom, error) {

	// fmt.Println(filepath)
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	newFeatures := unmarshalFeatureOrFeatureCollection(b)

	geoms := make([]Geom, len(newFeatures))
	for i := range newFeatures {
		p := newFeatures[i].Properties
		opts := p["options"]
		if opts != nil {
			prec := opts.(map[string]any)["precision"]
			if prec != nil {
				setPrecision(prec.(float64))
			}
		}
		fg := newFeatures[i].Geometry
		switch fg.Type {
		case "Polygon":
			geoms[i] = Geom{fg.Polygon}
		case "MultiPolygon":
			geoms[i] = fg.MultiPolygon
		default:
			return nil, fmt.Errorf("only polygon or multipolygon geometry types supported")
		}
	}

	return geoms, nil
}

func unmarshalFeatureOrFeatureCollection(b []byte) []*geojson.Feature {
	feature, err := geojson.UnmarshalFeature(b)
	if err != nil {
		return nil
	}
	if feature.Type != "FeatureCollection" {
		return []*geojson.Feature{feature}
	}
	fc, err := geojson.UnmarshalFeatureCollection(b)
	if err != nil {
		return nil
	}
	return fc.Features
}
