package rhsm

type APIType string

func (at APIType) String() string {
	return string(at)
}

const (
	CloudV2APIType APIType = "cloudapi-v2"
	WeldrAPIType   APIType = "weldr"
	TestAPIType    APIType = "test-manifest"
)

// The FactsImageOptions specify things to be stored into the Insights facts
// storage. This mostly relates to how the build of the image was performed.
type FactsImageOptions struct {
	APIType APIType
}
