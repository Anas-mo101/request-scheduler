// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type BodyType string

const (
	BodyTypeTEXT BodyType = "TEXT"
	BodyTypeJSON BodyType = "JSON"
)

func (e *BodyType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = BodyType(s)
	case string:
		*e = BodyType(s)
	default:
		return fmt.Errorf("unsupported scan type for BodyType: %T", src)
	}
	return nil
}

type NullBodyType struct {
	BodyType BodyType
	Valid    bool // Valid is true if BodyType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBodyType) Scan(value interface{}) error {
	if value == nil {
		ns.BodyType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.BodyType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBodyType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.BodyType), nil
}

type Method string

const (
	MethodGET  Method = "GET"
	MethodPOST Method = "POST"
)

func (e *Method) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Method(s)
	case string:
		*e = Method(s)
	default:
		return fmt.Errorf("unsupported scan type for Method: %T", src)
	}
	return nil
}

type NullMethod struct {
	Method Method
	Valid  bool // Valid is true if Method is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMethod) Scan(value interface{}) error {
	if value == nil {
		ns.Method, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Method.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMethod) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Method), nil
}

type Status string

const (
	StatusScheduled Status = "Scheduled"
	StatusInvoked   Status = "Invoked"
	StatusFailed    Status = "Failed"
)

func (e *Status) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Status(s)
	case string:
		*e = Status(s)
	default:
		return fmt.Errorf("unsupported scan type for Status: %T", src)
	}
	return nil
}

type NullStatus struct {
	Status Status
	Valid  bool // Valid is true if Status is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatus) Scan(value interface{}) error {
	if value == nil {
		ns.Status, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Status.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Status), nil
}

type Schedule struct {
	ID                  int32
	InvocationTimestamp pgtype.Timestamptz
	CreatedAt           pgtype.Timestamptz
	RequestMethod       Method
	RequestBodyType     NullBodyType
	RequestBody         pgtype.Text
	RequestUrl          string
	RequestHeader       []byte
	RequestQuery        []byte
	Status              Status
	RetriesNo           pgtype.Int4
	MaxRetries          pgtype.Int4
	FailureReason       pgtype.Text
}
