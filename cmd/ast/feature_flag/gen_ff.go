package feature_flag

type FeatureFlagsKey int

const (
	KeySaveWorld FeatureFlagsKey = iota
	KeyMakeVaccine
	KeyANewField
)

func (a FeatureFlags) IsFeatureEnabled(key FeatureFlagsKey) bool {
	result := false
	switch key {
	case KeySaveWorld:
		result = a.SaveWorld
	case KeyMakeVaccine:
		result = a.MakeVaccine
	case KeyANewField:
		result = a.ANewField
	}
	return result
}
