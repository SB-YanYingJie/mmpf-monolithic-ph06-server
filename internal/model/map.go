package model

import (
	gi "github.com/machinemapplatform/grpc-interface/golang"
	"github.com/machinemapplatform/library/model"
)

// ToModelImateType defines map from gi.ImageType to model.ImageType
var ToModelImageType = map[gi.NumberOfLenses]model.NumberOfLenses{
	gi.NumberOfLenses_MONO:   model.MONO,
	gi.NumberOfLenses_STEREO: model.STEREO,
}

var (
	// ToGrpcSlamState defines map from model.SlamState to gi.SlamState
	ToGrpcSlamState = map[model.SlamState]gi.SlamState{
		model.IDLE:                    gi.SlamState_IDLE,
		model.TRACKING_INITIALIZING:   gi.SlamState_TRACKING_INITIALIZING,
		model.TRACKING_GOOD:           gi.SlamState_TRACKING_GOOD,
		model.TRACKING_LOST:           gi.SlamState_TRACKING_LOST,
		model.MAP_LOADING_IN_PROGRESS: gi.SlamState_MAP_LOADING_IN_PROGRESS,
	}
	// ToGrpcMoveState defines map from model.MoveState to gi.MoveState
	ToGrpcMoveState = map[model.MoveState]gi.MoveState{
		model.ReachedWaypoint: gi.MoveState_REACHED_WAYPOINT,
		model.ReachedGoal:     gi.MoveState_REACHED_GOAL,
		model.Moving:          gi.MoveState_MOVING,
	}
	// ToGrpcErrorStatus defines map from model.ErrorStatus to uint32(codes.Code)
	ToGrpcErrorStatus = map[model.ErrorStatus]gi.ErrorStatus{
		model.NO_ERROR:         gi.ErrorStatus_NO_ERROR,
		model.INVALID_ARGUMENT: gi.ErrorStatus_INVALID_ARGUMENT,
		model.NOT_FOUND:        gi.ErrorStatus_NOT_FOUND,
		model.INTERNAL:         gi.ErrorStatus_INTERNAL,
	}
)
