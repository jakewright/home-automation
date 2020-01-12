package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/rpc"
	"github.com/jakewright/home-automation/libraries/go/util"
	sceneproto "github.com/jakewright/home-automation/service.scene/proto"
)

const (
	funcSleep           = "sleep"
	propertyTypeString  = "string"
	propertyTypeBoolean = "boolean"
	propertyTypeNumber  = "number"
	propertyTypeNull    = "null"
)

// Action is a single step in a scene
type Action struct {
	gorm.Model
	SceneID  int
	Stage    int
	Sequence int

	Func           string
	ControllerName string
	DeviceID       string
	Command        string
	Property       string
	PropertyValue  string
	PropertyType   string
}

// Validate checks that the action makes sense
func (a *Action) Validate() error {
	if !util.ExactlyOne(a.Func != "", a.Command != "", a.Property != "") {
		return errors.BadRequest("exactly one of func, command and property should be set")
	}

	switch {
	case a.Stage == 0:
		return errors.BadRequest("stage should be set to 1 or more")
	case a.Sequence == 0:
		return errors.BadRequest("sequence should be set to 1 or more")
	case a.Func == "" && a.ControllerName == "":
		return errors.BadRequest("controller_name should be set if setting property or calling command")
	case a.Func == "" && a.DeviceID == "":
		return errors.BadRequest("device_id should be set if setting property or calling command")
	case a.Property != "" && a.PropertyType != propertyTypeNull && a.PropertyValue == "":
		return errors.BadRequest("property_value cannot be blank unless property_type is \"null\"")
	case a.Property != "" && a.PropertyType == "":
		return errors.BadRequest("property_type should be set if setting property")

	case a.Func != "":
		if _, err := a.parseFunc(); err != nil {
			return err
		}

	case a.Command != "":
		if _, err := a.parseCommand(); err != nil {
			return err
		}

	case a.Property != "":
		if _, err := a.parseProperty(); err != nil {
			return err
		}
	}

	return nil
}

// Perform does the action
func (a *Action) Perform() error {
	if err := a.Validate(); err != nil {
		return err
	}

	var f func() error

	switch {
	case a.Func != "":
		f, _ = a.parseFunc()
	case a.Command != "":
		f, _ = a.parseCommand()
	case a.Property != "":
		f, _ = a.parseProperty()
	default:
		return nil
	}

	return f()
}

func (a *Action) parseFunc() (func() error, error) {
	parts := strings.Split(a.Func, " ")
	if len(parts) == 0 {
		return nil, errors.BadRequest("failed to extract func name from '%s'", a.Func)
	}

	switch parts[0] {
	case funcSleep:
		if len(parts) != 2 {
			return nil, errors.BadRequest("sleep func should have one argument")
		}

		d, err := time.ParseDuration(parts[1])
		if err != nil {
			return nil, err
		}

		return func() error {
			time.Sleep(d)
			return nil
		}, nil
	}

	return nil, errors.BadRequest("unknown func %s", parts[0])
}

func (a *Action) parseCommand() (func() error, error) {
	// todo
	return func() error {
		return nil
	}, nil
}

func (a *Action) parseProperty() (func() error, error) {
	url := fmt.Sprintf("%s/device/%s", a.ControllerName, a.DeviceID)

	val, err := marshalPropertyValue(a.PropertyType, a.PropertyValue)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal property value %s into type %s", a.PropertyValue, a.PropertyType)
	}

	body := map[string]interface{}{
		a.Property: val,
	}

	return func() error {
		_, err := rpc.Patch(url, body, nil)
		return err
	}, nil
}

func marshalPropertyValue(t, v string) (interface{}, error) {
	switch t {
	case propertyTypeString:
		return v, nil

	case propertyTypeBoolean:
		return strconv.ParseBool(v)

	case propertyTypeNumber:
		return strconv.ParseFloat(v, 64)

	case propertyTypeNull:
		return nil, nil
	}

	return nil, errors.BadRequest("unknown property type %s", t)
}

// ToProto marshals to the proto type
func (a *Action) ToProto() *sceneproto.Action {
	return &sceneproto.Action{
		Id:             uint32(a.ID),
		Stage:          int32(a.Stage),
		Sequence:       int32(a.Sequence),
		Func:           a.Func,
		ControllerName: a.ControllerName,
		Command:        a.Command,
		Property:       a.Property,
		PropertyValue:  a.PropertyValue,
		CreatedAt:      util.TimeToProto(a.CreatedAt),
		UpdatedAt:      util.TimeToProto(a.UpdatedAt),
		DeletedAt:      util.PTimeToProto(a.DeletedAt),
	}
}
