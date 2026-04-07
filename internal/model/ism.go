package model

// ISM represents the core Information Security Marking object per the json-ism spec.
// Fields map to IC-ISM XML attributes for REST API consumption.
type ISM struct {
	// Core Classification
	Version       string         `json:"version,omitempty"`
	Classification Classification `json:"classification"`
	OwnerProducer []string       `json:"ownerProducer,omitempty"`
	Joint         bool           `json:"joint,omitempty"`

	// CUI Fields
	CUIBasic           string   `json:"cuiBasic,omitempty"`
	CUISpecified       []string `json:"cuiSpecified,omitempty"`
	CategoryMarkings   []string `json:"categoryMarkings,omitempty"`
	ControlledByName   string   `json:"controlledByName,omitempty"`
	ControlledByOffice string   `json:"controlledByOffice,omitempty"`
	POC                string   `json:"poc,omitempty"`

	// Dissemination Controls
	DisseminationControls []string `json:"disseminationControls,omitempty"`
	ReleasableTo          []string `json:"releasableTo,omitempty"`
	DisplayOnlyTo         []string `json:"displayOnlyTo,omitempty"`

	// Distribution Statements (DoDI 5230.24)
	DistributionStatement              string                          `json:"distributionStatement,omitempty"`
	ThirdPartyDistributionStatement    string                          `json:"3rdPartyDistributionStatement,omitempty"`
	ThirdPartyDistributionWarning      string                          `json:"3rdPartyDistributionWarning,omitempty"`
	ThirdPartyDistributionContract     *ThirdPartyDistributionContract `json:"3rdPartyDistributionContract,omitempty"`
	Copyright                          string                          `json:"copyright,omitempty"`

	// Classification Authority
	ClassifiedBy              string `json:"classifiedBy,omitempty"`
	ClassificationReason      string `json:"classificationReason,omitempty"`
	DerivativelyClassifiedBy  string `json:"derivativelyClassifiedBy,omitempty"`
	DerivedFrom               string `json:"derivedFrom,omitempty"`
	CompilationReason         string `json:"compilationReason,omitempty"`

	// Declassification
	DeclassDate      string `json:"declassDate,omitempty"`
	DeclassEvent     string `json:"declassEvent,omitempty"`
	DeclassException string `json:"declassException,omitempty"`

	// Foreign Government Information
	FGISourceOpen      []string `json:"fgiSourceOpen,omitempty"`
	FGISourceProtected []string `json:"fgiSourceProtected,omitempty"`

	// Exemptions
	ExemptFrom []string `json:"exemptFrom,omitempty"`

	// Compliance
	CompliesWith []string `json:"compliesWith,omitempty"`

	// Atomic Energy
	AtomicEnergyMarkings []string `json:"atomicEnergyMarkings,omitempty"`

	// Notices
	NoticeType    []string `json:"noticeType,omitempty"`
	NoticeProseID string   `json:"noticeProseID,omitempty"`

	// Special Access Required
	SARIdentifier []string `json:"sarIdentifier,omitempty"`

	// NATO
	HighWaterNATO string `json:"highWaterNATO,omitempty"`

	// Other
	NonICMarkings          []string `json:"nonICMarkings,omitempty"`
	NonUSControls          []string `json:"nonUSControls,omitempty"`
	BannerLine             string   `json:"bannerLine,omitempty"`
	SecondBannerLine       string   `json:"secondBannerLine,omitempty"`
	CreateDate             string   `json:"createDate,omitempty"`
	HandleViaChannels      string   `json:"handleViaChannels,omitempty"`
	HasApproximateMarkings bool     `json:"hasApproximateMarkings,omitempty"`
	NoAggregation          bool     `json:"noAggregation,omitempty"`
	POCType                string   `json:"pocType,omitempty"`
}
