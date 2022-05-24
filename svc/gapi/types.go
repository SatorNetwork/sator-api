package gapi

import (
	"encoding/json"
)

// NFTType ...
type NFTType string

func (n NFTType) String() string {
	return string(n)
}

// Predefined NFT types
// common, rare, super rare, epic, legend
const (
	NFTTypeUndefined NFTType = "undefined"
	NFTTypeCommon    NFTType = "common"
	NFTTypeRare      NFTType = "rare"
	NFTTypeSuperRare NFTType = "super_rare"
	NFTTypeEpic      NFTType = "epic"
	NFTTypeLegend    NFTType = "legend"
)

// NFT types array
var nftTypesSlice = []NFTType{
	NFTTypeCommon,
	NFTTypeRare,
	NFTTypeSuperRare,
	NFTTypeEpic,
	NFTTypeLegend,
}

// Predefined game levels
const (
	GameLevelEasy = iota + 1
	GameLevelMedium
	GameLevelHard
)

// Game levels array
var gameLevels = []int{
	GameLevelEasy,
	GameLevelMedium,
	GameLevelHard,
}

// NFTInfo ...
type NFTInfo struct {
	ID       string  `json:"id"`
	MaxLevel int     `json:"max_level"`
	NftType  NFTType `json:"nft_type"`
}

// NFTPackInfo ...
type NFTPackInfo struct {
	ID          string      `json:"pack_id"`
	DropChances DropChances `json:"drop_chances"`
	Price       float64     `json:"price"`
}

// DropChances ...
type DropChances struct {
	Common    float64 `json:"common"`
	Rare      float64 `json:"rare"`
	SuperRare float64 `json:"super_rare"`
	Epic      float64 `json:"epic"`
	Legend    float64 `json:"legend"`
}

// GameConfig struct
type GameConfig struct {
	FactoryBlockInfo struct {
		BaseStep     int `json:"_baseStep"`
		BaseScale    int `json:"_baseScale"`
		MinRandomPos int `json:"_minRandomPos"`
		MaxRandomPos int `json:"_maxRandomPos"`
	} `json:"_factoryBlockInfo"`
	FactoryMoveInfo struct {
		MovementCurve []string `json:"_movementCurve"`
	} `json:"_factoryMoveInfo"`
	FactoryWaterInfo struct {
		MovementCurve     []string `json:"_movementCurve"`
		MinDistanceToSnap int      `json:"_minDistanceToSnap"`
	} `json:"_factoryWaterInfo"`
	FactoryWaterStaticInfo struct {
		BaseStep          int `json:"_baseStep"`
		BaseScale         int `json:"_baseScale"`
		MinRandomPos      int `json:"_minRandomPos"`
		MaxRandomPos      int `json:"_maxRandomPos"`
		MinDistanceToSnap int `json:"_minDistanceToSnap"`
	} `json:"_factoryWaterStaticInfo"`
	GeneralData struct {
		ForwardStep            int    `json:"_forwardStep"`
		ResourceFinishRoadPath string `json:"_resourceFinishRoadPath"`
	} `json:"_generalData"`
	BlocksCreateInfo struct {
		BlocksAmount int `json:"_blocksAmount"`
		BlocksInfo   []struct {
			ID                string        `json:"_id"`
			Type              string        `json:"_type"`
			ResourceRoadPath  string        `json:"_resourceRoadPath"`
			From              int           `json:"_from"`
			To                int           `json:"_to"`
			AmountMin         int           `json:"_amountMin"`
			AmountMax         int           `json:"_amountMax"`
			Chance            int           `json:"_chance"`
			OnExit            []interface{} `json:"_onExit"`
			PossibleMoveItems []struct {
				ID                    string  `json:"_id"`
				SpawnDelayRandomMin   int     `json:"_spawnDelayRandomMin"`
				SpawnDelayRandomMax   int     `json:"_spawnDelayRandomMax"`
				MoveTimeMin           int     `json:"_moveTimeMin"`
				MoveTimeMax           int     `json:"_moveTimeMax"`
				ReverseMovementChance float64 `json:"_reverseMovementChance"`
				ResourceMovementPath  string  `json:"_resourceMovementPath"`
			} `json:"_possibleMoveItems,omitempty"`
		} `json:"_blocksInfo"`
		BlocksExitInfo []struct {
			ID                   string        `json:"_id"`
			Type                 string        `json:"_type"`
			ResourceRoadPath     string        `json:"_resourceRoadPath"`
			From                 int           `json:"_from"`
			To                   int           `json:"_to"`
			AmountMin            int           `json:"_amountMin"`
			AmountMax            int           `json:"_amountMax"`
			Chance               int           `json:"_chance"`
			OnExit               []interface{} `json:"_onExit"`
			SizeRandomMin        int           `json:"_sizeRandomMin,omitempty"`
			SizeRandomMax        int           `json:"_sizeRandomMax,omitempty"`
			BlockAmountRandomMin int           `json:"_blockAmountRandomMin,omitempty"`
			BlockAmountRandomMax int           `json:"_blockAmountRandomMax,omitempty"`
			ResourceBlockPath    string        `json:"_resourceBlockPath,omitempty"`
			PossibleMoveItems    []struct {
				ID                    string `json:"_id"`
				SpawnDelayRandomMin   int    `json:"_spawnDelayRandomMin"`
				SpawnDelayRandomMax   int    `json:"_spawnDelayRandomMax"`
				MoveTimeMin           int    `json:"_moveTimeMin"`
				MoveTimeMax           int    `json:"_moveTimeMax"`
				ReverseMovementChance int    `json:"_reverseMovementChance"`
				ResourceMovementPath  string `json:"_resourceMovementPath"`
			} `json:"_possibleMoveItems,omitempty"`
		} `json:"_blocksExitInfo"`
	} `json:"_blocksCreateInfo"`
	LazyInfo struct {
		MaxBlocksAllowed     int `json:"_maxBlocksAllowed"`
		DistanceToRemoveItem int `json:"_distanceToRemoveItem"`
	} `json:"_lazyInfo"`
	CameraInfo struct {
		CamMoveCurve    []string `json:"_camMoveCurve"`
		CamLoseTime     int      `json:"_camLoseTime"`
		CamLoseDistance int      `json:"_camLoseDistance"`
	} `json:"_cameraInfo"`
}

func (g GameConfig) String() string {
	json, _ := json.Marshal(g)
	return string(json)
}
