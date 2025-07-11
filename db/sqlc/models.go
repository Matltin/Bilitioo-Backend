// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type ActivityStatus string

const (
	ActivityStatusPENDING      ActivityStatus = "PENDING"
	ActivityStatusREMINDERSENT ActivityStatus = "REMINDER-SENT"
	ActivityStatusPURCHASED    ActivityStatus = "PURCHASED"
)

func (e *ActivityStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ActivityStatus(s)
	case string:
		*e = ActivityStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ActivityStatus: %T", src)
	}
	return nil
}

type NullActivityStatus struct {
	ActivityStatus ActivityStatus `json:"activity_status"`
	Valid          bool           `json:"valid"` // Valid is true if ActivityStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullActivityStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ActivityStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ActivityStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullActivityStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ActivityStatus), nil
}

type CheckReservationTicketStatus string

const (
	CheckReservationTicketStatusRESERVED    CheckReservationTicketStatus = "RESERVED"
	CheckReservationTicketStatusNOTRESERVED CheckReservationTicketStatus = "NOT_RESERVED"
)

func (e *CheckReservationTicketStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CheckReservationTicketStatus(s)
	case string:
		*e = CheckReservationTicketStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for CheckReservationTicketStatus: %T", src)
	}
	return nil
}

type NullCheckReservationTicketStatus struct {
	CheckReservationTicketStatus CheckReservationTicketStatus `json:"check_reservation_ticket_status"`
	Valid                        bool                         `json:"valid"` // Valid is true if CheckReservationTicketStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCheckReservationTicketStatus) Scan(value interface{}) error {
	if value == nil {
		ns.CheckReservationTicketStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CheckReservationTicketStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCheckReservationTicketStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.CheckReservationTicketStatus), nil
}

type FlightClass string

const (
	FlightClassECONOMY        FlightClass = "ECONOMY"
	FlightClassPREMIUMECONOMY FlightClass = "PREMIUM-ECONOMY"
	FlightClassBUSINESS       FlightClass = "BUSINESS"
	FlightClassFIRST          FlightClass = "FIRST"
)

func (e *FlightClass) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FlightClass(s)
	case string:
		*e = FlightClass(s)
	default:
		return fmt.Errorf("unsupported scan type for FlightClass: %T", src)
	}
	return nil
}

type NullFlightClass struct {
	FlightClass FlightClass `json:"flight_class"`
	Valid       bool        `json:"valid"` // Valid is true if FlightClass is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullFlightClass) Scan(value interface{}) error {
	if value == nil {
		ns.FlightClass, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.FlightClass.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullFlightClass) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.FlightClass), nil
}

type NotificationLogStatus string

const (
	NotificationLogStatusSENT   NotificationLogStatus = "SENT"
	NotificationLogStatusFAILED NotificationLogStatus = "FAILED"
)

func (e *NotificationLogStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = NotificationLogStatus(s)
	case string:
		*e = NotificationLogStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for NotificationLogStatus: %T", src)
	}
	return nil
}

type NullNotificationLogStatus struct {
	NotificationLogStatus NotificationLogStatus `json:"notification_log_status"`
	Valid                 bool                  `json:"valid"` // Valid is true if NotificationLogStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullNotificationLogStatus) Scan(value interface{}) error {
	if value == nil {
		ns.NotificationLogStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.NotificationLogStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullNotificationLogStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.NotificationLogStatus), nil
}

type PaymentStatus string

const (
	PaymentStatusPENDING   PaymentStatus = "PENDING"
	PaymentStatusCOMPLETED PaymentStatus = "COMPLETED"
	PaymentStatusFAILED    PaymentStatus = "FAILED"
	PaymentStatusREFUNDED  PaymentStatus = "REFUNDED"
)

func (e *PaymentStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentStatus(s)
	case string:
		*e = PaymentStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentStatus: %T", src)
	}
	return nil
}

type NullPaymentStatus struct {
	PaymentStatus PaymentStatus `json:"payment_status"`
	Valid         bool          `json:"valid"` // Valid is true if PaymentStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentStatus), nil
}

type PaymentType string

const (
	PaymentTypeCASH         PaymentType = "CASH"
	PaymentTypeCREDITCARD   PaymentType = "CREDIT_CARD"
	PaymentTypeWALLET       PaymentType = "WALLET"
	PaymentTypeBANKTRANSFER PaymentType = "BANK_TRANSFER"
	PaymentTypeCRYPTO       PaymentType = "CRYPTO"
)

func (e *PaymentType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentType(s)
	case string:
		*e = PaymentType(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentType: %T", src)
	}
	return nil
}

type NullPaymentType struct {
	PaymentType PaymentType `json:"payment_type"`
	Valid       bool        `json:"valid"` // Valid is true if PaymentType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentType) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentType), nil
}

type RequestType string

const (
	RequestTypePAYMENTISSUE       RequestType = "PAYMENT-ISSUE"
	RequestTypeTRAVELDELAY        RequestType = "TRAVEL-DELAY"
	RequestTypeUNEXPECTEDRESERVED RequestType = "UNEXPECTED-RESERVED"
	RequestTypeETC                RequestType = "ETC."
)

func (e *RequestType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RequestType(s)
	case string:
		*e = RequestType(s)
	default:
		return fmt.Errorf("unsupported scan type for RequestType: %T", src)
	}
	return nil
}

type NullRequestType struct {
	RequestType RequestType `json:"request_type"`
	Valid       bool        `json:"valid"` // Valid is true if RequestType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRequestType) Scan(value interface{}) error {
	if value == nil {
		ns.RequestType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RequestType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRequestType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RequestType), nil
}

type Role string

const (
	RoleADMIN Role = "ADMIN"
	RoleUSER  Role = "USER"
)

func (e *Role) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Role(s)
	case string:
		*e = Role(s)
	default:
		return fmt.Errorf("unsupported scan type for Role: %T", src)
	}
	return nil
}

type NullRole struct {
	Role  Role `json:"role"`
	Valid bool `json:"valid"` // Valid is true if Role is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRole) Scan(value interface{}) error {
	if value == nil {
		ns.Role, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Role.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Role), nil
}

type TicketStatus string

const (
	TicketStatusRESERVED       TicketStatus = "RESERVED"
	TicketStatusRESERVING      TicketStatus = "RESERVING"
	TicketStatusCANCELED       TicketStatus = "CANCELED"
	TicketStatusCANCELEDBYTIME TicketStatus = "CANCELED-BY-TIME"
)

func (e *TicketStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TicketStatus(s)
	case string:
		*e = TicketStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TicketStatus: %T", src)
	}
	return nil
}

type NullTicketStatus struct {
	TicketStatus TicketStatus `json:"ticket_status"`
	Valid        bool         `json:"valid"` // Valid is true if TicketStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTicketStatus) Scan(value interface{}) error {
	if value == nil {
		ns.TicketStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TicketStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTicketStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TicketStatus), nil
}

type UserStatus string

const (
	UserStatusACTIVE    UserStatus = "ACTIVE"
	UserStatusNONACTIVE UserStatus = "NON-ACTIVE"
)

func (e *UserStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserStatus(s)
	case string:
		*e = UserStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for UserStatus: %T", src)
	}
	return nil
}

type NullUserStatus struct {
	UserStatus UserStatus `json:"user_status"`
	Valid      bool       `json:"valid"` // Valid is true if UserStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserStatus) Scan(value interface{}) error {
	if value == nil {
		ns.UserStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserStatus), nil
}

type VehicleType string

const (
	VehicleTypeBUS      VehicleType = "BUS"
	VehicleTypeTRAIN    VehicleType = "TRAIN"
	VehicleTypeAIRPLANE VehicleType = "AIRPLANE"
)

func (e *VehicleType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = VehicleType(s)
	case string:
		*e = VehicleType(s)
	default:
		return fmt.Errorf("unsupported scan type for VehicleType: %T", src)
	}
	return nil
}

type NullVehicleType struct {
	VehicleType VehicleType `json:"vehicle_type"`
	Valid       bool        `json:"valid"` // Valid is true if VehicleType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullVehicleType) Scan(value interface{}) error {
	if value == nil {
		ns.VehicleType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.VehicleType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullVehicleType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.VehicleType), nil
}

type Airplane struct {
	VehicleID   int64       `json:"vehicle_id"`
	FlightClass FlightClass `json:"flight_class"`
	Name        string      `json:"name"`
}

type AirplaneSeat struct {
	SeatID int64 `json:"seat_id"`
}

type Bus struct {
	VehicleID int64 `json:"vehicle_id"`
	VIP       bool  `json:"VIP"`
	BedChair  bool  `json:"bed_chair"`
}

type BusSeat struct {
	SeatID int64 `json:"seat_id"`
}

type ChangeReservation struct {
	ID            int64        `json:"id"`
	ReservationID int64        `json:"reservation_id"`
	AdminID       int64        `json:"admin_id"`
	UserID        int64        `json:"user_id"`
	FromStatus    TicketStatus `json:"from_status"`
	ToStatus      TicketStatus `json:"to_status"`
}

type City struct {
	ID       int64  `json:"id"`
	Province string `json:"province"`
	County   string `json:"county"`
}

type Company struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type NotificationLog struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"user_id"`
	SendEmailSmsID int64  `json:"send_email_sms_id"`
	MessageText    string `json:"message_text"`
	UserActivityID int64  `json:"user_activity_id"`
}

type Payment struct {
	ID          int64         `json:"id"`
	FromAccount int64         `json:"from_account"`
	ToAccount   string        `json:"to_account"`
	Amount      int64         `json:"amount"`
	Type        PaymentType   `json:"type"`
	Status      PaymentStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
}

type Penalty struct {
	ID          int64  `json:"id"`
	VehicleID   int64  `json:"vehicle_id"`
	PenaltyText string `json:"penalty_text"`
	BeforDay    int32  `json:"befor_day"`
	AfterDay    int32  `json:"after_day"`
}

type Profile struct {
	UserID       int64  `json:"user_id"`
	PicDir       string `json:"pic_dir"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	CityID       int64  `json:"city_id"`
	Wallet       int64  `json:"wallet"`
	NationalCode string `json:"national_code"`
}

type Report struct {
	ID            int64       `json:"id"`
	ReservationID int64       `json:"reservation_id"`
	UserID        int64       `json:"user_id"`
	AdminID       int64       `json:"admin_id"`
	RequestType   RequestType `json:"request_type"`
	RequestText   string      `json:"request_text"`
	ResponseText  string      `json:"response_text"`
}

type Reservation struct {
	ID           int64         `json:"id"`
	UserID       int64         `json:"user_id"`
	TicketID     int64         `json:"ticket_id"`
	PaymentID    int64         `json:"payment_id"`
	Status       TicketStatus  `json:"status"`
	DurationTime time.Duration `json:"duration_time"`
	CreatedAt    time.Time     `json:"created_at"`
}

type Route struct {
	ID                    int64         `json:"id"`
	OriginCityID          int64         `json:"origin_city_id"`
	DestinationCityID     int64         `json:"destination_city_id"`
	OriginTerminalID      sql.NullInt64 `json:"origin_terminal_id"`
	DestinationTerminalID sql.NullInt64 `json:"destination_terminal_id"`
	Distance              int32         `json:"distance"`
}

type Seat struct {
	ID          int64       `json:"id"`
	VehicleID   int64       `json:"vehicle_id"`
	VehicleType VehicleType `json:"vehicle_type"`
	SeatNumber  int32       `json:"seat_number"`
	IsAvailable bool        `json:"is_available"`
}

type SendEmailSm struct {
	ID          int64                 `json:"id"`
	MessageType string                `json:"message_type"`
	SentAt      time.Time             `json:"sent_at"`
	Status      NotificationLogStatus `json:"status"`
}

type SendVerificationCode struct {
	ID             int64         `json:"id"`
	UserID         int64         `json:"user_id"`
	SendEmailSmsID int64         `json:"send_email_sms_id"`
	Token          string        `json:"token"`
	DurationTime   time.Duration `json:"duration_time"`
	CreatedAt      time.Time     `json:"created_at"`
}

type Terminal struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Ticket struct {
	ID            int64                        `json:"id"`
	VehicleID     int64                        `json:"vehicle_id"`
	SeatID        int64                        `json:"seat_id"`
	VehicleType   VehicleType                  `json:"vehicle_type"`
	RouteID       int64                        `json:"route_id"`
	Amount        int64                        `json:"amount"`
	DepartureTime time.Time                    `json:"departure_time"`
	ArrivalTime   time.Time                    `json:"arrival_time"`
	CountStand    int32                        `json:"count_stand"`
	Status        CheckReservationTicketStatus `json:"status"`
	CreatedAt     time.Time                    `json:"created_at"`
}

type Train struct {
	VehicleID       int64 `json:"vehicle_id"`
	Rank            int32 `json:"rank"`
	HaveCompartment bool  `json:"have_compartment"`
}

type TrainSeat struct {
	SeatID          int64         `json:"seat_id"`
	Salon           int32         `json:"salon"`
	HaveCompartment bool          `json:"have_compartment"`
	CuopeNumber     sql.NullInt32 `json:"cuope_number"`
}

type User struct {
	ID               int64      `json:"id"`
	Email            string     `json:"email"`
	PhoneNumber      string     `json:"phone_number"`
	HashedPassword   string     `json:"hashed_password"`
	PasswordChangeAt time.Time  `json:"password_change_at"`
	Role             Role       `json:"role"`
	Status           UserStatus `json:"status"`
	PhoneVerified    bool       `json:"phone_verified"`
	EmailVerified    bool       `json:"email_verified"`
	CreatedAt        time.Time  `json:"created_at"`
}

type UserActivity struct {
	ID           int64          `json:"id"`
	UserID       int64          `json:"user_id"`
	RouteID      int64          `json:"route_id"`
	VehicleType  VehicleType    `json:"vehicle_type"`
	Status       ActivityStatus `json:"status"`
	DurationTime time.Duration  `json:"duration_time"`
	CreatedAt    time.Time      `json:"created_at"`
}

type Vehicle struct {
	ID          int64           `json:"id"`
	CompanyID   int64           `json:"company_id"`
	Capacity    int32           `json:"capacity"`
	VehicleType VehicleType     `json:"vehicle_type"`
	Feature     json.RawMessage `json:"feature"`
}

type VerifyEmail struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}
