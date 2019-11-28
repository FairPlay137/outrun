package muxhandlers

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/fluofoxxo/outrun/analytics"
	"github.com/fluofoxxo/outrun/analytics/factors"
	"github.com/fluofoxxo/outrun/config/infoconf"
	"github.com/fluofoxxo/outrun/db"
	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/logic"
	"github.com/fluofoxxo/outrun/logic/conversion"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/requests"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
)

func Login(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.LoginRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	uid := request.LineAuth.UserID
	password := request.LineAuth.Password

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if uid == "0" && password == "" {
		helper.Out("Entering LoginAlpha")
		newPlayer := db.NewAccount()
		err = db.SavePlayer(newPlayer)
		if err != nil {
			helper.InternalErr("Error saving player", err)
			return
		}
		baseInfo.StatusCode = status.InvalidPassword
		baseInfo.SetErrorMessage(emess.BadPassword)
		response := responses.LoginRegister(
			baseInfo,
			newPlayer.ID,
			newPlayer.Password,
			newPlayer.Key,
		)
		err = helper.SendResponse(response)
		if err != nil {
			helper.InternalErr("Error responding", err)
		}
		return
	} else if uid == "0" && password != "" {
		helper.Out("Entering LoginBravo")
		// invalid request
		helper.InvalidRequest()
		return
	} else if uid != "0" && password == "" {
		helper.Out("Entering LoginCharlie")
		// game wants to log in
		baseInfo.StatusCode = status.InvalidPassword
		baseInfo.SetErrorMessage(emess.BadPassword)
		player, err := db.GetPlayer(uid)
		if err != nil {
			helper.InternalErr("Error getting player", err)
			// likely account that wasn't found, so let's tell them that:
			response := responses.LoginCheckKey(baseInfo, "")
			baseInfo.StatusCode = status.MissingPlayer
			helper.SendResponse(response)
			return
		}
		response := responses.LoginCheckKey(baseInfo, player.Key)
		err = helper.SendResponse(response)
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	} else if uid != "0" && password != "" {
		helper.Out("Entering LoginDelta")
		// game is attempting to log in using key
		// for now, we pretend that it worked no matter what
		// TODO: fix this obvious security flaw
		baseInfo.StatusCode = status.OK
		baseInfo.SetErrorMessage(emess.OK)
		sid, err := db.AssignSessionID(uid)
		if err != nil {
			helper.InternalErr("Error assigning session ID", err)
			return
		}
		player, err := db.GetPlayer(uid)
		if err != nil {
			helper.InternalErr("Error getting player", err)
			return
		}
		player.LastLogin = time.Now().UTC().Unix()
		err = db.SavePlayer(player)
		if err != nil {
			helper.InternalErr("Error saving player", err)
			return
		}
		response := responses.LoginSuccess(baseInfo, sid, player.Username)
		err = helper.SendResponse(response)
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		analytics.Store(player.ID, factors.AnalyticTypeLogins)
		return
	}
}

func GetVariousParameter(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.VariousParameter(baseInfo, player)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
}

func GetInformation(helper *helper.Helper) {
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	infos := []obj.Information{}
	helper.DebugOut("%v", infoconf.CFile.EnableInfos)
	if infoconf.CFile.EnableInfos {
		for _, ci := range infoconf.CFile.Infos {
			newInfo := conversion.ConfiguredInfoToInformation(ci)
			infos = append(infos, newInfo)
			helper.DebugOut(newInfo.Param)
		}
	}
	operatorInfos := []obj.OperatorInformation{}
	numOpUnread := int64(len(operatorInfos))
	response := responses.Information(baseInfo, infos, operatorInfos, numOpUnread)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetTicker(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultTicker(baseInfo, player)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func LoginBonus(helper *helper.Helper) {
	// TODO: Is agnostic, but shouldn't be!
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultLoginBonus(baseInfo)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetCountry(helper *helper.Helper) {
	// TODO: Should get correct country code!
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultGetCountry(baseInfo)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetMigrationPassword(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.GetMigrationPasswordRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	player.UserPassword = request.UserPassword // TODO: Confirm that this is the right behavior
	db.SavePlayer(player)
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.MigrationPassword(baseInfo, player)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func Migration(helper *helper.Helper) {
	randChar := func(charset string, length int64) string {
		runes := []rune(charset)
		final := make([]rune, 10)
		for i := range final {
			final[i] = runes[rand.Intn(len(runes))]
		}
		return string(final)
	}

	recv := helper.GetGameRequest()
	var request requests.LoginRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	password := request.LineAuth.Password
	migrationUserPassword := request.LineAuth.MigrationPassword

	baseInfo := helper.BaseInfo(emess.OK, status.OK)

	foundPlayers, err := logic.FindPlayersByPassword(password, false)
	if err != nil {
		helper.Err("Error finding players by password", err)
		return
	}
	playerIDs := []string{}
	for _, player := range foundPlayers {
		playerIDs = append(playerIDs, player.ID)
	}
	if len(playerIDs) > 0 {
		migratePlayer, err := db.GetPlayer(playerIDs[0])
		if err != nil {
			helper.InternalErr("Error getting player", err)
			return
		}
		if migrationUserPassword == migratePlayer.UserPassword {
			baseInfo.StatusCode = status.OK
			baseInfo.SetErrorMessage(emess.OK)
			migratePlayer.SetPassword(randChar("abcdefghijklmnopqrstuvwxyz1234567890", 10)) //generate a brand new password
			migratePlayer.LastLogin = time.Now().UTC().Unix()
			err = db.SavePlayer(migratePlayer)
			if err != nil {
				helper.InternalErr("Error saving player", err)
				return
			}
			sid, err := db.AssignSessionID(migratePlayer.ID)
			if err != nil {
				helper.InternalErr("Error assigning session ID", err)
				return
			}
			response := responses.MigrationSuccess(baseInfo, sid, migratePlayer.ID, migratePlayer.Username, migratePlayer.Password)
			helper.SendResponse(response)
		} else {
			response := responses.NewBaseResponse(baseInfo)
			baseInfo.StatusCode = status.InvalidPassword
			baseInfo.SetErrorMessage(emess.BadPassword)
			helper.SendResponse(response)
		}
	} else {
		response := responses.NewBaseResponse(baseInfo)
		baseInfo.StatusCode = status.MissingPlayer // TODO: Is this the correct error code?
		helper.SendResponse(response)
	}
}
