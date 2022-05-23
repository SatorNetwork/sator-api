package gapi

var defaultGameConfig = `{
  "_factoryBlockInfo": {
    "_baseStep": 5,
    "_baseScale": 5,
    "_minRandomPos": -25,
    "_maxRandomPos": 25
  },
  "_factoryMoveInfo": {
    "_movementCurve": [
      "0,0",
      "1,1"
    ]
  },
  "_factoryWaterInfo": {
    "_movementCurve": [
      "0,0",
      "0.06904668,0.2065322",
      "0.2607135,0.3969892",
      "0.8250648,0.855166",
      "1.005554,0.9608002"
    ],
    "_minDistanceToSnap": 3
  },
  "_factoryWaterStaticInfo": {
    "_baseStep": 5,
    "_baseScale": 5,
    "_minRandomPos": -15,
    "_maxRandomPos": 15,
    "_minDistanceToSnap": 3
  },
  "_generalData": {
    "_forwardStep": 10,
    "_resourceFinishRoadPath": "RoadCross/RoadFinishPrefab"
  },
  "_blocksCreateInfo": {
    "_blocksAmount": 75,
    "_blocksInfo": [
      {
        "_id": "start",
        "_type": "Simple",
        "_resourceRoadPath": "RoadCross/RoadStartZonePrefab",
        "_from": 0,
        "_to": 1,
        "_amountMin": 1,
        "_amountMax": 1,
        "_chance": 1,
        "_onExit": []
      },
      {
        "_id": "start_small_car_chain",
        "_type": "Move",
        "_resourceRoadPath": "RoadCross/RoadSimplePrefab",
        "_from": 1,
        "_to": 15,
        "_amountMin": 2,
        "_amountMax": 3,
        "_chance": 0.600000023841858,
        "_onExit": [
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 4,
            "_spawnDelayRandomMax": 5,
            "_moveTimeMin": 10,
            "_moveTimeMax": 11,
            "_reverseMovementChance": 0.5,
            "_resourceMovementPath": "RoadCross/MovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 5,
            "_spawnDelayRandomMax": 6,
            "_moveTimeMin": 9,
            "_moveTimeMax": 10,
            "_reverseMovementChance": 0.5,
            "_resourceMovementPath": "RoadCross/MovementPrefab"
          }
        ]
      },
      {
        "_id": "start_big_car_chain",
        "_type": "Move",
        "_resourceRoadPath": "RoadCross/RoadSimplePrefab",
        "_from": 1,
        "_to": 15,
        "_amountMin": 1,
        "_amountMax": 3,
        "_chance": 0.300000011920929,
        "_onExit": [
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 5,
            "_spawnDelayRandomMax": 6,
            "_moveTimeMin": 10,
            "_moveTimeMax": 12,
            "_reverseMovementChance": 0.5,
            "_resourceMovementPath": "RoadCross/MovementBigPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 5,
            "_spawnDelayRandomMax": 7,
            "_moveTimeMin": 9,
            "_moveTimeMax": 10,
            "_reverseMovementChance": 0.5,
            "_resourceMovementPath": "RoadCross/MovementBigPrefab"
          }
        ]
      },
      {
        "_id": "start_big_car_non_chain",
        "_type": "Move",
        "_resourceRoadPath": "RoadCross/RoadSimplePrefab",
        "_from": 1,
        "_to": 15,
        "_amountMin": 1,
        "_amountMax": 1,
        "_chance": 0.100000001490116,
        "_onExit": [],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 4,
            "_spawnDelayRandomMax": 5,
            "_moveTimeMin": 11,
            "_moveTimeMax": 11,
            "_reverseMovementChance": 0.5,
            "_resourceMovementPath": "RoadCross/MovementBigPrefab"
          }
        ]
      },
      {
        "_id": "water_chain",
        "_type": "Water",
        "_resourceRoadPath": "RoadCross/RoadWaterPrefab",
        "_from": 15,
        "_to": 45,
        "_amountMin": 2,
        "_amountMax": 5,
        "_chance": 0.300000011920929,
        "_onExit": [
          "WaterStaticZone",
          "WaterZone_1",
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 10,
            "_moveTimeMax": 11,
            "_reverseMovementChance": 0,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 3,
            "_spawnDelayRandomMax": 4,
            "_moveTimeMin": 11,
            "_moveTimeMax": 12,
            "_reverseMovementChance": 1,
            "_resourceMovementPath": "RoadCross/WaterMovementBigPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 10,
            "_moveTimeMax": 11,
            "_reverseMovementChance": 0,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 11,
            "_moveTimeMax": 12,
            "_reverseMovementChance": 1,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          }
        ]
      },
      {
        "_id": "water_non_chain",
        "_type": "Water",
        "_resourceRoadPath": "RoadCross/RoadWaterPrefab",
        "_from": 15,
        "_to": 45,
        "_amountMin": 1,
        "_amountMax": 1,
        "_chance": 0.100000001490116,
        "_onExit": [
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 10,
            "_moveTimeMax": 12,
            "_reverseMovementChance": 0.5,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          }
        ]
      },
      {
        "_id": "small_car_chain",
        "_type": "Move",
        "_resourceRoadPath": "RoadCross/RoadSimplePrefab",
        "_from": 15,
        "_to": 45,
        "_amountMin": 3,
        "_amountMax": 6,
        "_chance": 0.300000011920929,
        "_onExit": [
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 4,
            "_spawnDelayRandomMax": 5,
            "_moveTimeMin": 8,
            "_moveTimeMax": 8,
            "_reverseMovementChance": 0.300000011920929,
            "_resourceMovementPath": "RoadCross/MovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 4,
            "_spawnDelayRandomMax": 5,
            "_moveTimeMin": 9,
            "_moveTimeMax": 9,
            "_reverseMovementChance": 0.699999988079071,
            "_resourceMovementPath": "RoadCross/MovementPrefab"
          }
        ]
      },
      {
        "_id": "big_car_chain",
        "_type": "Move",
        "_resourceRoadPath": "RoadCross/RoadSimplePrefab",
        "_from": 15,
        "_to": 45,
        "_amountMin": 2,
        "_amountMax": 5,
        "_chance": 0.200000002980232,
        "_onExit": [
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 5,
            "_spawnDelayRandomMax": 5,
            "_moveTimeMin": 8,
            "_moveTimeMax": 8,
            "_reverseMovementChance": 0.300000011920929,
            "_resourceMovementPath": "RoadCross/MovementBigPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 5,
            "_spawnDelayRandomMax": 6,
            "_moveTimeMin": 10,
            "_moveTimeMax": 10,
            "_reverseMovementChance": 0.699999988079071,
            "_resourceMovementPath": "RoadCross/MovementBigPrefab"
          }
        ]
      },
      {
        "_id": "big_car_non_chain",
        "_type": "Move",
        "_resourceRoadPath": "RoadCross/RoadSimplePrefab",
        "_from": 15,
        "_to": 45,
        "_amountMin": 1,
        "_amountMax": 1,
        "_chance": 0.100000001490116,
        "_onExit": [
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 5,
            "_spawnDelayRandomMax": 5,
            "_moveTimeMin": 8,
            "_moveTimeMax": 10,
            "_reverseMovementChance": 0.5,
            "_resourceMovementPath": "RoadCross/MovementBigPrefab"
          }
        ]
      },
      {
        "_id": "late_water_chain",
        "_type": "Water",
        "_resourceRoadPath": "RoadCross/RoadWaterPrefab",
        "_from": 45,
        "_to": 75,
        "_amountMin": 2,
        "_amountMax": 4,
        "_chance": 0.300000011920929,
        "_onExit": [
          "WaterStaticZone",
          "WaterZoneLate",
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 9,
            "_moveTimeMax": 10,
            "_reverseMovementChance": 0,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 3,
            "_spawnDelayRandomMax": 4,
            "_moveTimeMin": 11,
            "_moveTimeMax": 12,
            "_reverseMovementChance": 1,
            "_resourceMovementPath": "RoadCross/WaterMovementBigPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 9,
            "_moveTimeMax": 10,
            "_reverseMovementChance": 0,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 11,
            "_moveTimeMax": 12,
            "_reverseMovementChance": 1,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          }
        ]
      },
      {
        "_id": "late_water_non_chain",
        "_type": "Water",
        "_resourceRoadPath": "RoadCross/RoadWaterPrefab",
        "_from": 45,
        "_to": 75,
        "_amountMin": 1,
        "_amountMax": 1,
        "_chance": 0.100000001490116,
        "_onExit": [
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 10,
            "_moveTimeMax": 10,
            "_reverseMovementChance": 0.5,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          }
        ]
      },
      {
        "_id": "late_small_car_chain",
        "_type": "Move",
        "_resourceRoadPath": "RoadCross/RoadSimplePrefab",
        "_from": 45,
        "_to": 75,
        "_amountMin": 3,
        "_amountMax": 6,
        "_chance": 0.300000011920929,
        "_onExit": [
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 4,
            "_spawnDelayRandomMax": 5,
            "_moveTimeMin": 7,
            "_moveTimeMax": 7,
            "_reverseMovementChance": 0.300000011920929,
            "_resourceMovementPath": "RoadCross/MovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 4,
            "_spawnDelayRandomMax": 5,
            "_moveTimeMin": 8,
            "_moveTimeMax": 9,
            "_reverseMovementChance": 0.699999988079071,
            "_resourceMovementPath": "RoadCross/MovementPrefab"
          }
        ]
      },
      {
        "_id": "late_big_car_chain",
        "_type": "Move",
        "_resourceRoadPath": "RoadCross/RoadSimplePrefab",
        "_from": 45,
        "_to": 75,
        "_amountMin": 2,
        "_amountMax": 5,
        "_chance": 0.200000002980232,
        "_onExit": [
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 5,
            "_spawnDelayRandomMax": 5,
            "_moveTimeMin": 7,
            "_moveTimeMax": 8,
            "_reverseMovementChance": 0.300000011920929,
            "_resourceMovementPath": "RoadCross/MovementBigPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 5,
            "_spawnDelayRandomMax": 6,
            "_moveTimeMin": 9,
            "_moveTimeMax": 10,
            "_reverseMovementChance": 0.699999988079071,
            "_resourceMovementPath": "RoadCross/MovementBigPrefab"
          }
        ]
      },
      {
        "_id": "late_big_car_non_chain",
        "_type": "Move",
        "_resourceRoadPath": "RoadCross/RoadSimplePrefab",
        "_from": 45,
        "_to": 75,
        "_amountMin": 1,
        "_amountMax": 1,
        "_chance": 0.100000001490116,
        "_onExit": [
          "SafeZone_1"
        ],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 5,
            "_spawnDelayRandomMax": 5,
            "_moveTimeMin": 8,
            "_moveTimeMax": 8,
            "_reverseMovementChance": 0.5,
            "_resourceMovementPath": "RoadCross/MovementBigPrefab"
          }
        ]
      }
    ],
    "_blocksExitInfo": [
      {
        "_id": "SafeZone_1",
        "_type": "Block",
        "_resourceRoadPath": "RoadCross/RoadSafeZonePrefab",
        "_from": 0,
        "_to": 0,
        "_amountMin": 1,
        "_amountMax": 1,
        "_chance": 1,
        "_onExit": [],
        "_sizeRandomMin": 1,
        "_sizeRandomMax": 2,
        "_blockAmountRandomMin": 1,
        "_blockAmountRandomMax": 3,
        "_resourceBlockPath": "RoadCross/BlockPrefab"
      },
      {
        "_id": "CarZone_2",
        "_type": "Move",
        "_resourceRoadPath": "RoadCross/RoadSimplePrefab",
        "_from": 0,
        "_to": 0,
        "_amountMin": 1,
        "_amountMax": 1,
        "_chance": 0,
        "_onExit": [],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 2,
            "_moveTimeMin": 1,
            "_moveTimeMax": 1,
            "_reverseMovementChance": 1,
            "_resourceMovementPath": "RoadCross/MovementBigPrefab"
          }
        ]
      },
      {
        "_id": "WaterStaticZone",
        "_type": "WaterStatic",
        "_resourceRoadPath": "RoadCross/RoadWaterPrefab",
        "_from": 0,
        "_to": 0,
        "_amountMin": 1,
        "_amountMax": 1,
        "_chance": 1,
        "_onExit": [],
        "_sizeRandomMin": 1,
        "_sizeRandomMax": 1,
        "_blockAmountRandomMin": 1,
        "_blockAmountRandomMax": 1,
        "_resourceBlockPath": "RoadCross/WaterStaticPrefab"
      },
      {
        "_id": "WaterZone_1",
        "_type": "Water",
        "_resourceRoadPath": "RoadCross/RoadWaterPrefab",
        "_from": 1,
        "_to": 15,
        "_amountMin": 2,
        "_amountMax": 4,
        "_chance": 1,
        "_onExit": [],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 10,
            "_moveTimeMax": 11,
            "_reverseMovementChance": 0,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 3,
            "_spawnDelayRandomMax": 4,
            "_moveTimeMin": 11,
            "_moveTimeMax": 12,
            "_reverseMovementChance": 1,
            "_resourceMovementPath": "RoadCross/WaterMovementBigPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 10,
            "_moveTimeMax": 11,
            "_reverseMovementChance": 0,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 11,
            "_moveTimeMax": 12,
            "_reverseMovementChance": 1,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          }
        ]
      },
      {
        "_id": "WaterZoneLate",
        "_type": "Water",
        "_resourceRoadPath": "RoadCross/RoadWaterPrefab",
        "_from": 0,
        "_to": 0,
        "_amountMin": 2,
        "_amountMax": 4,
        "_chance": 1,
        "_onExit": [],
        "_possibleMoveItems": [
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 9,
            "_moveTimeMax": 10,
            "_reverseMovementChance": 0,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 3,
            "_spawnDelayRandomMax": 4,
            "_moveTimeMin": 11,
            "_moveTimeMax": 12,
            "_reverseMovementChance": 1,
            "_resourceMovementPath": "RoadCross/WaterMovementBigPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 9,
            "_moveTimeMax": 10,
            "_reverseMovementChance": 0,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          },
          {
            "_id": "",
            "_spawnDelayRandomMin": 2,
            "_spawnDelayRandomMax": 3,
            "_moveTimeMin": 11,
            "_moveTimeMax": 12,
            "_reverseMovementChance": 1,
            "_resourceMovementPath": "RoadCross/WaterMovementPrefab"
          }
        ]
      }
    ]
  },
  "_lazyInfo": {
    "_maxBlocksAllowed": 40,
    "_distanceToRemoveItem": 50
  },
  "_cameraInfo": {
    "_camMoveCurve": [
      "0,0",
      "1,1"
    ],
    "_camLoseTime": 10,
    "_camLoseDistance": 30
  }
}`
