package asp

type GranularityType string
type Name string

var (
	Marketplace                    GranularityType = "Marketplace"
	ResearchingQuantityInShortTerm Name            = "researchingQuantityInShortTerm"
	ResearchingQuantityInMidTerm   Name            = "researchingQuantityInMidTerm"
	ResearchingQuantityInLongTerm  Name            = "researchingQuantityInLongTerm"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type GetInventorySummariesResponse struct {
	Payload    GetInventorySummariesResult `json:"payload"`
	Pagination Pagination                  `json:"pagination"`
	Errors     []Error                     `json:"errors"`
}

type Pagination struct {
	NextToken string `json:"nextToken"`
}

type GetInventorySummariesResult struct {
	Granularity        Granularity          `json:"granularity"`
	InventorySummaries []InventorySummaries `json:"inventorySummaries"`
}

type Granularity struct {
	GranularityID   string          `json:"granularityId"`
	GranularityType GranularityType `json:"granularityType"`
}

type InventorySummaries struct {
	Asin             string           `json:"asin"`
	FnSku            string           `json:"fnSku"`
	SellerSku        string           `json:"sellerSku"`
	Condition        string           `json:"condition"`
	InventoryDetails InventoryDetails `json:"inventoryDetails"`
	LastUpdatedTime  string           `json:"lastUpdatedTime"`
	ProductName      string           `json:"productName"`
	TotalQuantity    int              `json:"totalQuantity"`
}

type InventoryDetails struct {
	FulfillableQuantity      int                   `json:"fulfillableQuantity"`
	InboundWorkingQuantity   int                   `json:"inboundWorkingQuantity"`
	InboundShippedQuantity   int                   `json:"inboundShippedQuantity"`
	InboundReceivingQuantity int                   `json:"inboundReceivingQuantity"`
	ReservedQuantity         ReservedQuantity      `json:"reservedQuantity"`
	ResearchingQuantity      ResearchingQuantity   `json:"researchingQuantity"`
	UnfulfillableQuantity    UnfulfillableQuantity `json:"unfulfillableQuantity"`
}

type ReservedQuantity struct {
	TotalReservedQuantity        int `json:"totalReservedQuantity"`
	PendingCustomerOrderQuantity int `json:"pendingCustomerOrderQuantity"`
	PendingTransshipmentQuantity int `json:"pendingTransshipmentQuantity"`
	FcProcessingQuantity         int `json:"fcProcessingQuantity"`
}

type ResearchingQuantity struct {
	TotalResearchingQuantity     int                        `json:"totalResearchingQuantity"`
	ResearchingQuantityBreakdown []ResearchingQuantityEntry `json:"researchingQuantityBreakdown"`
}

type ResearchingQuantityEntry struct {
	Name     Name `json:"name"`
	Quantity int  `json:"quantity"`
}

type UnfulfillableQuantity struct {
	TotalUnfulfillableQuantity int `json:"totalUnfulfillableQuantity"`
	CustomerDamagedQuantity    int `json:"customerDamagedQuantity"`
	WarehouseDamagedQuantity   int `json:"warehouseDamagedQuantity"`
	DistributorDamagedQuantity int `json:"distributorDamagedQuantity"`
	CarrierDamagedQuantity     int `json:"carrierDamagedQuantity"`
	DefectiveQuantity          int `json:"defectiveQuantity"`
	ExpiredQuantity            int `json:"expiredQuantity"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}
