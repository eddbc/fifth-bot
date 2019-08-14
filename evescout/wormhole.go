package evescout

import "time"

//Wormhole Wormhole as reported by eve-scout.com
type Wormhole struct {
	ID                                int32       `json:"id"`
	SignatureID                       string      `json:"signatureId"`
	Type                              string      `json:"type"`
	Status                            string      `json:"status"`
	WormholeMass                      string      `json:"wormholeMass"`
	WormholeEol                       string      `json:"wormholeEol"`
	WormholeEstimatedEol              time.Time   `json:"wormholeEstimatedEol"`
	WormholeDestinationSignatureID    string      `json:"wormholeDestinationSignatureId"`
	CreatedAt                         time.Time   `json:"createdAt"`
	UpdatedAt                         time.Time   `json:"updatedAt"`
	DeletedAt                         interface{} `json:"deletedAt"`
	StatusUpdatedAt                   interface{} `json:"statusUpdatedAt"`
	CreatedBy                         string      `json:"createdBy"`
	CreatedByID                       string      `json:"createdById"`
	DeletedBy                         interface{} `json:"deletedBy"`
	DeletedByID                       interface{} `json:"deletedById"`
	WormholeSourceWormholeTypeID      int32       `json:"wormholeSourceWormholeTypeId"`
	WormholeDestinationWormholeTypeID int32       `json:"wormholeDestinationWormholeTypeId"`
	SolarSystemID                     int32       `json:"solarSystemId"`
	WormholeDestinationSolarSystemID  int32       `json:"wormholeDestinationSolarSystemId"`
	SourceWormholeType                struct {
		ID       int32  `json:"id"`
		Name     string `json:"name"`
		Src      string `json:"src"`
		Dest     string `json:"dest"`
		Lifetime int32  `json:"lifetime"`
		JumpMass int32  `json:"jumpMass"`
		MaxMass  int32  `json:"maxMass"`
	} `json:"sourceWormholeType"`
	DestinationWormholeType struct {
		ID       int32  `json:"id"`
		Name     string `json:"name"`
		Src      string `json:"src"`
		Dest     string `json:"dest"`
		Lifetime int32  `json:"lifetime"`
		JumpMass int32  `json:"jumpMass"`
		MaxMass  int32  `json:"maxMass"`
	} `json:"destinationWormholeType"`
	SourceSolarSystem struct {
		ID              int32   `json:"id"`
		Name            string  `json:"name"`
		ConstellationID int32   `json:"constellationID"`
		Security        float64 `json:"security"`
		RegionID        int32   `json:"regionId"`
		Region          struct {
			ID   int32  `json:"id"`
			Name string `json:"name"`
		} `json:"region"`
	} `json:"sourceSolarSystem"`
	DestinationSolarSystem struct {
		ID              int32   `json:"id"`
		Name            string  `json:"name"`
		ConstellationID int32   `json:"constellationID"`
		Security        float64 `json:"security"`
		RegionID        int32   `json:"regionId"`
		Region          struct {
			ID   int32  `json:"id"`
			Name string `json:"name"`
		} `json:"region"`
	} `json:"destinationSolarSystem"`
}
