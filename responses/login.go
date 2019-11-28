package responses

import (
	"strconv"
	"time"

	"github.com/fluofoxxo/outrun/config/infoconf"
	"github.com/fluofoxxo/outrun/logic/conversion"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/obj/constobjs"
	"github.com/fluofoxxo/outrun/responses/responseobjs"
)

type LoginCheckKeyResponse struct {
	BaseResponse
	Key string `json:"key"`
}

func LoginCheckKey(base responseobjs.BaseInfo, key string) LoginCheckKeyResponse {
	baseResponse := NewBaseResponse(base)
	out := LoginCheckKeyResponse{
		baseResponse,
		key,
	}
	return out
}

type LoginRegisterResponse struct {
	BaseResponse
	UserID      string `json:"userId"`
	Password    string `json:"password"`
	Key         string `json:"key"`
	CountryID   int64  `json:"countryId,string"`
	CountryCode string `json:"countryCode"`
}

func LoginRegister(base responseobjs.BaseInfo, uid, password, key string) LoginRegisterResponse {
	// TODO: Fetch correct country code
	baseResponse := NewBaseResponse(base)
	out := LoginRegisterResponse{
		baseResponse,
		uid,
		password,
		key,
		1,
		"US",
	}
	return out
}

type GetCountryResponse struct {
	BaseResponse
	CountryID   int64  `json:"countryId,string"`
	CountryCode string `json:"countryCode"`
}

func GetCountry(base responseobjs.BaseInfo, countryID int64, countryCode string) GetCountryResponse {
	baseResponse := NewBaseResponse(base)
	return GetCountryResponse{
		baseResponse,
		countryID,
		countryCode,
	}
}

func DefaultGetCountry(base responseobjs.BaseInfo) GetCountryResponse {
	return GetCountry(
		base,
		1,
		"US",
	)
}

type LoginSuccessResponse struct {
	BaseResponse
	Username             string   `json:"userName"`
	SessionID            string   `json:"sessionId"`
	SessionTimeLimit     int64    `json:"sessionTimeLimit"`         // game will log in again after this (non-UTC apparently) time
	EnergyRecoveryTime   int64    `json:"energyRecveryTime,string"` // seconds until energy regenerates (misspelling is _actually_ in the game!)
	EnergyRecoveryMax    int64    `json:"energyRecoveryMax,string"` // maximum energy recoverable over time
	InviteBasicIncentive obj.Item `json:"inviteBasicIncentiv"`
}

func LoginSuccess(base responseobjs.BaseInfo, sid, username string) LoginSuccessResponse {
	baseResponse := NewBaseResponse(base)
	out := LoginSuccessResponse{
		baseResponse,
		username,
		sid,
		time.Now().Unix() + 3600, // hour from now  // TODO: does this need to be UTC?
		//time.Now().Unix() + 120, // two minutes from now
		600, // 10 minutes from now, regen energy
		5,
		obj.NewItem("900000", 5),
	}
	return out
}

type VariousParameterResponse struct {
	BaseResponse
	netobj.PlayerVarious
}

func VariousParameter(base responseobjs.BaseInfo, player netobj.Player) VariousParameterResponse {
	baseResponse := NewBaseResponse(base)
	out := VariousParameterResponse{
		baseResponse,
		player.PlayerVarious,
	}
	return out
}

type InformationResponse struct {
	BaseResponse
	Infos             []obj.Information         `json:"informations"`
	OperatorInfos     []obj.OperatorInformation `json:"operatorEachInfos"`
	NumOperatorUnread int64                     `json:"numOperatorInfo"`
}

func Information(base responseobjs.BaseInfo, infos []obj.Information, opinfos []obj.OperatorInformation, numOpUnread int64) InformationResponse {
	baseResponse := NewBaseResponse(base)
	out := InformationResponse{
		baseResponse,
		infos,
		opinfos,
		numOpUnread,
	}
	return out
}

func DefaultInformation(base responseobjs.BaseInfo) InformationResponse {
	infos := constobjs.DefaultInformations
	opinfos := []obj.OperatorInformation{}
	numOpUnread := int64(len(opinfos))
	return Information(
		base,
		infos,
		opinfos,
		numOpUnread,
	)
}

type TickerResponse struct {
	BaseResponse
	TickerList []obj.Ticker `json:"tickerList"`
}

func Ticker(base responseobjs.BaseInfo, tickerList []obj.Ticker) TickerResponse {
	baseResponse := NewBaseResponse(base)
	return TickerResponse{
		baseResponse,
		tickerList,
	}
}

func DefaultTicker(base responseobjs.BaseInfo, player netobj.Player) TickerResponse {
	/*
		tickerList := []obj.Ticker{
			obj.NewTicker(
				1,
				time.Now().UTC().Unix()+3600, // one hour later
				"Welcome to [ff0000]OUTRUN!",
			),
			obj.NewTicker(
				2,
				time.Now().UTC().Unix()+7200,
				"ID: [0000ff]"+player.ID,
			),
			obj.NewTicker(
				3,
				time.Now().UTC().Unix()+7200, // two hours later
				"High score (Timed Mode): [0000ff]"+strconv.Itoa(int(player.PlayerState.TimedHighScore)),
			),
			obj.NewTicker(
				4,
				time.Now().UTC().Unix()+7200, // two hours later
				"High score (Story Mode): [0000ff]"+strconv.Itoa(int(player.PlayerState.HighScore)),
			),
			obj.NewTicker(
				5,
				time.Now().UTC().Unix()+7200, // two hours later
				"Total distance ran (Story Mode): [0000ff]"+strconv.Itoa(int(player.PlayerState.TotalDistance)),
			),
		}
	*/
	tickerList := []obj.Ticker{}
	if infoconf.CFile.EnableTickers {
		di := 0
		if !infoconf.CFile.HideWatermarkTicker {
			tickerList = []obj.Ticker{
				obj.NewTicker(
					1,
					time.Now().UTC().Unix()+3600, // one hour later
					"This server is powered by [ff0000]Outrun!",
				),
				obj.NewTicker(
					2,
					time.Now().UTC().Unix()+7200,
					"ID: [0000ff]"+player.ID,
				),
				obj.NewTicker(
					3,
					time.Now().UTC().Unix()+7200, // two hours later
					"High score (Timed Mode): [0000ff]"+strconv.Itoa(int(player.PlayerState.TimedHighScore)),
				),
				obj.NewTicker(
					4,
					time.Now().UTC().Unix()+7200, // two hours later
					"High score (Story Mode): [0000ff]"+strconv.Itoa(int(player.PlayerState.HighScore)),
				),
				obj.NewTicker(
					5,
					time.Now().UTC().Unix()+7200, // two hours later
					"Total distance ran (Story Mode): [0000ff]"+strconv.Itoa(int(player.PlayerState.TotalDistance)),
				)}
			di = 5
		}
		for i, ct := range infoconf.CFile.Tickers {
			newTicker := conversion.ConfiguredTickerToTicker(int64(di+i+1), ct)
			tickerList = append(tickerList, newTicker)
		}
	}
	return Ticker(
		base,
		tickerList,
	)
}

type LoginBonusResponse struct {
	BaseResponse
	LoginBonusStatus          obj.LoginBonusStatus   `json:"loginBonusStatus"`
	LoginBonusRewardList      []obj.LoginBonusReward `json:"loginBonusRewardList"`
	FirstLoginBonusRewardList []obj.LoginBonusReward `json:"firstLoginBonusRewardList"`
	StartTime                 int64                  `json:"startTime"`
	EndTime                   int64                  `json:"endTime"`
	RewardID                  int64                  `json:"rewardId"`
	RewardDays                int64                  `json:"rewardDays"`
	FirstRewardDays           int64                  `json:"firstRewardDays"`
}

func LoginBonus(base responseobjs.BaseInfo, lbs obj.LoginBonusStatus, lbrl, flbrl []obj.LoginBonusReward, st, et, rid, rd, frd int64) LoginBonusResponse {
	baseResponse := NewBaseResponse(base)
	return LoginBonusResponse{
		baseResponse,
		lbs,
		lbrl,
		flbrl,
		st,
		et,
		rid,
		rd,
		frd,
	}
}

func DefaultLoginBonus(base responseobjs.BaseInfo) LoginBonusResponse {
	lbs := obj.NewLoginBonusStatus(2, 2, 1465830000)
	lbrl := constobjs.DefaultLoginBonusRewardList
	flbrl := constobjs.DefaultFirstLoginBonusRewardList
	st := int64(1465743600)
	et := int64(1466348400)
	rid := int64(-1)
	rd := int64(-1)
	frd := int64(-1)
	return LoginBonus(base, lbs, lbrl, flbrl, st, et, rid, rd, frd)
}

type MigrationPasswordResponse struct {
	BaseResponse
	Password string `json:"password"`
}

func MigrationPassword(base responseobjs.BaseInfo, player netobj.Player) MigrationPasswordResponse {
	baseResponse := NewBaseResponse(base)
	return MigrationPasswordResponse{
		baseResponse,
		player.Password,
	}
}

type MigrationSuccessResponse struct {
	BaseResponse
	UserID                   string   `json:"userId"`
	Username                 string   `json:"userName"`
	Password                 string   `json:"password"`
	SessionID                string   `json:"sessionId"`
	SessionTimeLimit         int64    `json:"sessionTimeLimit"`         // game will log in again after this (non-UTC apparently) time
	EnergyRecoveryTime       int64    `json:"energyRecveryTime,string"` // seconds until energy regenerates (misspelling is _actually_ in the game!)
	EnergyRecoveryMax        int64    `json:"energyRecoveryMax,string"` // maximum energy recoverable over time
	InviteBasicIncentive     obj.Item `json:"inviteBasicIncentiv"`
	ChaoRentalBasicIncentive obj.Item `json:"chaoRentalBasicIncentiv"`
}

func MigrationSuccess(base responseobjs.BaseInfo, sid, uid, username, password string) MigrationSuccessResponse {
	baseResponse := NewBaseResponse(base)
	out := MigrationSuccessResponse{
		baseResponse,
		uid,
		username,
		password,
		sid,
		time.Now().Unix() + 3600, // hour from now  // TODO: does this need to be UTC?
		//time.Now().Unix() + 120, // two minutes from now
		600, // 10 minutes from now, regen energy
		5,
		obj.NewItem("900000", 5),
		obj.NewItem("900000", 5),
	}
	return out
}
