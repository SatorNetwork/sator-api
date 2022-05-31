package gapi

import (
	"encoding/json"
)

// Game result type
type GameResult int32

const (
	GameResultLose GameResult = iota
	GameResultWin
)

// GameResult to string
func (g GameResult) String() string {
	switch g {
	case GameResultLose:
		return "lose"
	case GameResultWin:
		return "win"
	default:
		return "undefined"
	}
}

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

// Check if nft type is valid
func (n NFTType) IsValid() bool {
	switch n {
	case NFTTypeCommon, NFTTypeRare, NFTTypeSuperRare, NFTTypeEpic, NFTTypeLegend:
		return true
	}
	return false
}

// NFTType to int
func (n NFTType) ToInt() int {
	switch n {
	case NFTTypeCommon:
		return 0
	case NFTTypeRare:
		return 1
	case NFTTypeSuperRare:
		return 2
	case NFTTypeEpic:
		return 3
	case NFTTypeLegend:
		return 4
	default:
		return -1
	}
}

func getNextNFTType(nftType NFTType) NFTType {
	switch nftType {
	case NFTTypeCommon:
		return NFTTypeRare
	case NFTTypeRare:
		return NFTTypeSuperRare
	case NFTTypeSuperRare:
		return NFTTypeEpic
	case NFTTypeEpic:
		return NFTTypeLegend
	case NFTTypeLegend:
		return NFTTypeLegend
	default:
		return NFTTypeUndefined
	}
}

// Predefined game levels
const (
	GameLevelEasy int32 = iota + 1
	GameLevelNormal
	GameLevelHard
)

func getGameLevelName(level int32) string {
	switch level {
	case GameLevelEasy:
		return "easy"
	case GameLevelNormal:
		return "normal"
	case GameLevelHard:
		return "hard"
	default:
		return "default"
	}
}

// getNFTLevelByType returns max NFT level depending on NFT type
func getNFTLevelByType(t NFTType) int32 {
	switch t {
	case NFTTypeCommon:
		return GameLevelEasy
	case NFTTypeRare:
		return GameLevelNormal
	case NFTTypeSuperRare, NFTTypeEpic, NFTTypeLegend:
		return GameLevelHard
	default:
		return GameLevelEasy
	}
}

// NFTInfo ...
type NFTInfo struct {
	ID       string  `json:"id"`
	MaxLevel int32   `json:"max_level"`
	NftType  NFTType `json:"nft_type"`
}

// NFTPackInfo ...
type NFTPackInfo struct {
	ID          string      `json:"pack_id"`
	Name        string      `json:"name"`
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

// DropChances to map
func (d DropChances) ToMap() map[string]float64 {
	return map[string]float64{
		NFTTypeCommon.String():    d.Common,
		NFTTypeRare.String():      d.Rare,
		NFTTypeSuperRare.String(): d.SuperRare,
		NFTTypeEpic.String():      d.Epic,
		NFTTypeLegend.String():    d.Legend,
	}
}

// DropChances to bytes
func (d DropChances) Bytes() []byte {
	bytes, _ := json.Marshal(d)
	return bytes
}

// GameConfig struct
type GameConfig struct {
	FactoryBlockInfo struct {
		BaseStep     float64 `json:"_baseStep"`
		BaseScale    float64 `json:"_baseScale"`
		MinRandomPos float64 `json:"_minRandomPos"`
		MaxRandomPos float64 `json:"_maxRandomPos"`
	} `json:"_factoryBlockInfo"`
	FactoryMoveInfo struct {
		MovementCurve []string `json:"_movementCurve"`
	} `json:"_factoryMoveInfo"`
	FactoryWaterInfo struct {
		MovementCurve     []string `json:"_movementCurve"`
		MinDistanceToSnap float64  `json:"_minDistanceToSnap"`
	} `json:"_factoryWaterInfo"`
	FactoryWaterStaticInfo struct {
		BaseStep          float64 `json:"_baseStep"`
		BaseScale         float64 `json:"_baseScale"`
		MinRandomPos      float64 `json:"_minRandomPos"`
		MaxRandomPos      float64 `json:"_maxRandomPos"`
		MinDistanceToSnap float64 `json:"_minDistanceToSnap"`
	} `json:"_factoryWaterStaticInfo"`
	GeneralData struct {
		ForwardStep            float64 `json:"_forwardStep"`
		ResourceFinishRoadPath string  `json:"_resourceFinishRoadPath"`
	} `json:"_generalData"`
	BlocksCreateInfo struct {
		BlocksAmount int `json:"_blocksAmount"`
		BlocksInfo   []struct {
			ID                string   `json:"_id"`
			Type              string   `json:"_type"`
			ResourceRoadPath  string   `json:"_resourceRoadPath"`
			From              int      `json:"_from"`
			To                int      `json:"_to"`
			AmountMin         int      `json:"_amountMin"`
			AmountMax         int      `json:"_amountMax"`
			Chance            float64  `json:"_chance"`
			OnExit            []string `json:"_onExit"`
			PossibleMoveItems []struct {
				ID                    string  `json:"_id"`
				SpawnDelayRandomMin   float64 `json:"_spawnDelayRandomMin"`
				SpawnDelayRandomMax   float64 `json:"_spawnDelayRandomMax"`
				MoveTimeMin           float64 `json:"_moveTimeMin"`
				MoveTimeMax           float64 `json:"_moveTimeMax"`
				ReverseMovementChance float64 `json:"_reverseMovementChance"`
				ResourceMovementPath  string  `json:"_resourceMovementPath"`
			} `json:"_possibleMoveItems,omitempty"`
		} `json:"_blocksInfo"`
		BlocksExitInfo []struct {
			ID                   string   `json:"_id"`
			Type                 string   `json:"_type"`
			ResourceRoadPath     string   `json:"_resourceRoadPath"`
			From                 int      `json:"_from"`
			To                   int      `json:"_to"`
			AmountMin            int      `json:"_amountMin"`
			AmountMax            int      `json:"_amountMax"`
			Chance               float64  `json:"_chance"`
			OnExit               []string `json:"_onExit"`
			SizeRandomMin        float64  `json:"_sizeRandomMin,omitempty"`
			SizeRandomMax        float64  `json:"_sizeRandomMax,omitempty"`
			BlockAmountRandomMin int      `json:"_blockAmountRandomMin,omitempty"`
			BlockAmountRandomMax int      `json:"_blockAmountRandomMax,omitempty"`
			ResourceBlockPath    string   `json:"_resourceBlockPath,omitempty"`
			PossibleMoveItems    []struct {
				ID                    string  `json:"_id"`
				SpawnDelayRandomMin   float64 `json:"_spawnDelayRandomMin"`
				SpawnDelayRandomMax   float64 `json:"_spawnDelayRandomMax"`
				MoveTimeMin           float64 `json:"_moveTimeMin"`
				MoveTimeMax           float64 `json:"_moveTimeMax"`
				ReverseMovementChance float64 `json:"_reverseMovementChance"`
				ResourceMovementPath  string  `json:"_resourceMovementPath"`
			} `json:"_possibleMoveItems,omitempty"`
		} `json:"_blocksExitInfo"`
	} `json:"_blocksCreateInfo"`
	LazyInfo struct {
		MaxBlocksAllowed     int     `json:"_maxBlocksAllowed"`
		DistanceToRemoveItem float64 `json:"_distanceToRemoveItem"`
	} `json:"_lazyInfo"`
	CameraInfo struct {
		CamMoveCurve    []string `json:"_camMoveCurve"`
		CamLoseTime     float64  `json:"_camLoseTime"`
		CamLoseDistance float64  `json:"_camLoseDistance"`
	} `json:"_cameraInfo"`
}

func (g GameConfig) String() string {
	json, _ := json.Marshal(g)
	return string(json)
}
