package libre

import (
	"fmt"
	"librelink-up-tg/internal/utils"
)

type CountryCode int32

const (
	AE CountryCode = iota
	AP
	AU
	CA
	DE
	EU
	EU2
	FR
	JP
	US
	LA
	RU
)

var (
	LLU_API_ENDPOINTS = map[CountryCode]string{
		AE:  "api-ae.libreview.io",
		AP:  "api-ap.libreview.io",
		AU:  "api-au.libreview.io",
		CA:  "api-ca.libreview.io",
		DE:  "api-de.libreview.io",
		EU:  "api-eu.libreview.io",
		EU2: "api-eu2.libreview.io",
		FR:  "api-fr.libreview.io",
		JP:  "api-jp.libreview.io",
		US:  "api-us.libreview.io",
		LA:  "api-la.libreview.io",
		RU:  "api.libreview.ru",
	}
	trends = map[int]string{
		1: "Резко падает ↓↓",
		2: "Падает ↓",
		3: "Стабильно →",
		4: "Возрастает ↗",
		5: "Резко возрастает ↗↗",
	}
)

func IsValidRegion(region int32) bool {
	if _, ok := LLU_API_ENDPOINTS[CountryCode(region)]; !ok {
		return false
	}
	return true
}

type H struct {
	Th   float64 `json:"th"`
	Thmm float64 `json:"thmm"`
	D    float64 `json:"d"`
	F    float64 `json:"f"`
}

type F struct {
	Th   float64 `json:"th"`
	Thmm float64 `json:"thmm"`
	D    float64 `json:"d"`
	Tl   float64 `json:"tl"`
	Tlmm float64 `json:"tlmm"`
}

type L struct {
	Th   float64 `json:"th"`
	Thmm float64 `json:"thmm"`
	D    float64 `json:"d"`
	Tl   float64 `json:"tl"`
	Tlmm float64 `json:"tlmm"`
}

type Nd struct {
	I float64 `json:"i"`
	R float64 `json:"r"`
	L float64 `json:"l"`
}

type Std struct {
}

type AlarmRules struct {
	C   bool    `json:"c"`
	H   H       `json:"h"`
	F   F       `json:"f"`
	L   L       `json:"l"`
	Nd  Nd      `json:"nd"`
	P   float64 `json:"p"`
	R   float64 `json:"r"`
	Std Std     `json:"std"`
}

type GlucoseItem struct {
	FactoryTimestamp string       `json:"FactoryTimestamp"`
	Timestamp        string       `json:"Timestamp"`
	Type             int          `json:"type"`
	ValueInMgPerDl   float64      `json:"ValueInMgPerDl"`
	TrendArrow       *int         `json:"TrendArrow,omitempty"`
	TrendMessage     *interface{} `json:"TrendMessage,omitempty"`
	MeasurementColor int          `json:"MeasurementColor"`
	GlucoseUnits     int          `json:"GlucoseUnits"`
	Value            float64      `json:"Value"`
	IsHigh           bool         `json:"isHigh"`
	IsLow            bool         `json:"isLow"`
}

type GlucoseMeasurement struct {
	GlucoseItem
	TrendArrow int `json:"TrendArrow"`
}

func (m *GlucoseMeasurement) GetMmolDivideLiter() float64 {
	return 0.0555 * m.ValueInMgPerDl
}

type FixedLowAlarmValues struct {
	Mgdl  float64 `json:"mgdl"`
	Mmoll float64 `json:"mmoll"`
}

type PatientDevice struct {
	DID                 string              `json:"did"`
	Dtid                int                 `json:"dtid"`
	V                   string              `json:"v"`
	Ll                  float64             `json:"ll"`
	Hl                  float64             `json:"hl"`
	U                   float64             `json:"u"`
	FixedLowAlarmValues FixedLowAlarmValues `json:"fixedLowAlarmValues"`
	Alarms              bool                `json:"alarms"`
	FixedLowThreshold   float64             `json:"fixedLowThreshold"`
	L                   *bool               `json:"l,omitempty"`
	H                   *bool               `json:"h,omitempty"`
}

type Sensor struct {
	DeviceID string  `json:"deviceId"`
	Sn       string  `json:"sn"`
	A        float64 `json:"a"`
	W        float64 `json:"w"`
	Pt       float64 `json:"pt"`
	S        bool    `json:"s"`
	Lj       bool    `json:"lj"`
}

type AuthTicket struct {
	Token    string  `json:"token"`
	Expires  float64 `json:"expires"`
	Duration float64 `json:"duration"`
}

type Connection struct {
	ID                 string             `json:"id"`
	PatientID          string             `json:"patientId"`
	Country            string             `json:"country"`
	Status             int                `json:"status"`
	FirstName          string             `json:"firstName"`
	LastName           string             `json:"lastName"`
	TargetLow          float64            `json:"targetLow"`
	TargetHigh         float64            `json:"targetHigh"`
	Uom                float64            `json:"uom"`
	Sensor             Sensor             `json:"sensor"`
	AlarmRules         AlarmRules         `json:"alarmRules"`
	GlucoseMeasurement GlucoseMeasurement `json:"glucoseMeasurement"`
	GlucoseItem        GlucoseItem        `json:"glucoseItem"`
	GlucoseAlarm       *interface{}       `json:"glucoseAlarm,omitempty"`
	PatientDevice      PatientDevice      `json:"patientDevice"`
	Created            float64            `json:"created"`
}

type ConnectionsResponse struct {
	Status int          `json:"status"`
	Data   []Connection `json:"data"`
	Ticket AuthTicket   `json:"ticket"`
}

type ActiveSensor struct {
	Sensor Sensor        `json:"sensor"`
	Device PatientDevice `json:"device"`
}

type GraphData struct {
	Connection    Connection     `json:"connection"`
	ActiveSensors []ActiveSensor `json:"activeSensors"`
	GraphData     []GlucoseItem  `json:"graphData"`
}

func (g *GraphData) IsBullshit() bool {
	return g.Connection.GlucoseMeasurement.IsLow
}

// Функция для преобразования TrendArrow в текст
func getTrend(trend int) string {
	return trends[trend]
}

func (data *GraphData) String() string {
	// Определяем статус
	status := "В пределах целевого диапазона"
	if data.Connection.GlucoseMeasurement.IsHigh {
		status = "Выше целевого диапазона!"
	} else if data.Connection.GlucoseMeasurement.IsLow {
		status = "Ниже целевого диапазона!"
	}

	// Форматируем сообщение
	return fmt.Sprintf(
		"Текущий уровень глюкозы: %.1f (%.1f) ммоль/л\n"+
			"Время измерения: %s\n"+
			"Тренд: %s\n"+
			"Состояние: %s (%.1f - %.1f ммоль/л)",
		data.Connection.GlucoseMeasurement.Value,
		data.Connection.GlucoseMeasurement.GetMmolDivideLiter(),
		utils.FormatTime(data.Connection.GlucoseMeasurement.Timestamp),
		getTrend(data.Connection.GlucoseMeasurement.TrendArrow),
		status,
		data.Connection.TargetLow/18,  // Конвертация из mg/dL в ммоль/л
		data.Connection.TargetHigh/18, // Конвертация из mg/dL в ммоль/л
	)
}

type GraphResponse struct {
	Status int        `json:"status"`
	Data   GraphData  `json:"data"`
	Ticket AuthTicket `json:"ticket"`
}

type Messages struct {
	FirstUsePhoenix                  float64 `json:"firstUsePhoenix"`
	FirstUsePhoenixReportsDataMerged float64 `json:"firstUsePhoenixReportsDataMerged"`
	LluAnalyticsNewAccount           float64 `json:"lluAnalyticsNewAccount"`
	LluGettingStartedBanner          float64 `json:"lluGettingStartedBanner"`
	LluNewFeatureModal               float64 `json:"lluNewFeatureModal"`
	LvWebPostRelease                 string  `json:"lvWebPostRelease"`
}

type System struct {
	Messages Messages `json:"messages"`
}

type TwoFactor struct {
	PrimaryMethod   string `json:"primaryMethod"`
	PrimaryValue    string `json:"primaryValue"`
	SecondaryMethod string `json:"secondaryMethod"`
	SecondaryValue  string `json:"secondaryValue"`
}

type History struct {
	PolicyAccept float64 `json:"policyAccept"`
}

type RealWorldEvidence struct {
	PolicyAccept float64   `json:"policyAccept"`
	TouAccept    float64   `json:"touAccept"`
	History      []History `json:"history"`
}

type Consents struct {
	RealWorldEvidence RealWorldEvidence `json:"realWorldEvidence"`
}

type User struct {
	ID                    string      `json:"id"`
	FirstName             string      `json:"firstName"`
	LastName              string      `json:"lastName"`
	Email                 string      `json:"email"`
	Country               string      `json:"country"`
	UiLanguage            string      `json:"uiLanguage"`
	CommunicationLanguage string      `json:"communicationLanguage"`
	AccountType           string      `json:"accountType"`
	Uom                   string      `json:"uom"`
	DateFormat            string      `json:"dateFormat"`
	TimeFormat            string      `json:"timeFormat"`
	EmailDay              []int       `json:"emailDay"`
	System                System      `json:"system"`
	TwoFactor             TwoFactor   `json:"twoFactor"`
	Created               float64     `json:"created"`
	LastLogin             float64     `json:"lastLogin"`
	Programs              interface{} `json:"programs"`
	DateOfBirth           float64     `json:"dateOfBirth"`
	Practices             interface{} `json:"practices"`
	Consents              Consents    `json:"consents"`
}

type MessagesSummary struct {
	Unread int `json:"unread"`
}

type Notifications struct {
	Unresolved int `json:"unresolved"`
}

type LoginData struct {
	User               User            `json:"user"`
	Messages           MessagesSummary `json:"messages"`
	Notifications      Notifications   `json:"notifications"`
	AuthTicket         AuthTicket      `json:"authTicket"`
	Invitations        *interface{}    `json:"invitations,omitempty"`
	TrustedDeviceToken string          `json:"trustedDeviceToken"`
	Redirect           *bool           `json:"redirect,omitempty"`
	Region             *string         `json:"region,omitempty"`
}

type LoginResponse struct {
	Status int       `json:"status"`
	Data   LoginData `json:"data"`
}

type LibreLinkUpHttpHeaders struct {
	Version string `json:"version"`
	Product string `json:"product"`
}
