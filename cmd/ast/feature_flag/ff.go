package feature_flag

type FeatureFlags struct {
	SaveWorld   bool `json:"save_world"`
	MakeVaccine bool `json:"make_vaccine"`
	ANewField   bool
}
