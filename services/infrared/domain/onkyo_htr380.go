package domain

import (
	"context"
	"math"

	"github.com/jakewright/home-automation/libraries/go/device"
	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	"github.com/jakewright/home-automation/libraries/go/oops"
	deviceregistrydef "github.com/jakewright/home-automation/services/device-registry/def"
	"github.com/jakewright/home-automation/services/infrared/ir"
)

const (
	onkyoHTR380KeyPower      = "KEY_POWER"
	onkyoHTR380KeyBDDVD      = "KEY_DVD"
	onkyoHTR380KeyVCRDVD     = "KEY_VCR"
	onkyoHTR380KeyCBLSAT     = "KEY_SAT"
	onkyoHTR380KeyGAME       = "BTN_GAMEPAD"
	onkyoHTR380KeyAUX        = "KEY_AUX"
	onkyoHTR380KeyTUNER      = "KEY_TUNER"
	onkyoHTR380KeyTVCD       = "KEY_TV"
	onkyoHTR380KeyPORT       = "KEY_TV2"
	onkyoHTR380KeyVolumeUp   = "KEY_VOLUMEUP"
	onkyoHTR380KeyVolumeDown = "KEY_VOLUMEDOWN"
	onkyoHTR380KeyMute       = "KEY_MUTE"

	onkyoHTR380InputBDDVD  = "BD_DVD"
	onkyoHTR380InputVCRDVD = "VCR_DVD"
	onkyoHTR380InputCBLSAT = "CBL_SAT"
	onkyoHTR380InputGAME   = "GAME"
	onkyoHTR380InputAUX    = "AUX"
	onkyoHTR380InputTUNER  = "TUNER"
	onkyoHTR380InputTVCD   = "TV_CD"
	onkyoHTR380InputPORT   = "PORT"
)

var onkyoHTR380InputKeys = map[string]string{
	onkyoHTR380InputBDDVD:  onkyoHTR380KeyBDDVD,
	onkyoHTR380InputVCRDVD: onkyoHTR380KeyVCRDVD,
	onkyoHTR380InputCBLSAT: onkyoHTR380KeyCBLSAT,
	onkyoHTR380InputGAME:   onkyoHTR380KeyGAME,
	onkyoHTR380InputAUX:    onkyoHTR380KeyAUX,
	onkyoHTR380InputTUNER:  onkyoHTR380KeyTUNER,
	onkyoHTR380InputTVCD:   onkyoHTR380KeyTVCD,
	onkyoHTR380InputPORT:   onkyoHTR380KeyPORT,
}

type OnkyoHTR380 struct {
	*deviceregistrydef.DeviceHeader
	power bool
}

func (d *OnkyoHTR380) ID() string {
	return d.Id
}

func (d *OnkyoHTR380) key(key string) ir.Instruction {
	return ir.Key("ONKYO_HT_R380", key)
}

func (d *OnkyoHTR380) LoadState(ctx context.Context) error {
	state, err := device.LoadProvidedState(ctx, d.Id, d.StateProviders)
	if err != nil {
		return oops.WithMessage(err, "failed to load provided state for device %q", d.Id)
	}

	if power, ok := state["power"].(bool); !ok {
		return oops.InternalService("state provider didn't provide power state for device %q", d.Id)
	} else {
		d.power = power
	}

	return nil
}

func (d *OnkyoHTR380) InstructionsFromState(state map[string]interface{}) ([]ir.Instruction, error) {
	if err := device.ValidateState(state, d.ToDef()); err != nil {
		return nil, err
	}

	var instructions []ir.Instruction

	if power, ok := state["power"].(bool); ok && power != d.power {
		instructions = append(instructions, d.key(onkyoHTR380KeyPower))

		// Give the AV receiver plenty of time to turn on
		instructions = append(instructions, ir.Wait(5000))
	}

	return instructions, nil
}

func (d *OnkyoHTR380) InstructionFromCommand(command string, args map[string]interface{}) ([]ir.Instruction, error) {
	if err := device.ValidateCommand(command, args, d.ToDef().Commands); err != nil {
		return nil, err
	}

	var instructions []ir.Instruction

	switch command {
	case "volume":
		delta, ok := args["delta"].(float64)
		if !ok {
			return nil, oops.BadRequest("invalid volume argument 'delta': %v", args["delta"])
		}

		key := onkyoHTR380KeyVolumeUp
		if delta < 0 {
			key = onkyoHTR380KeyVolumeDown
		}

		// Send the key n + 1 times because the key needs to be
		// pressed to activate the volume control before it changes
		for i := 0.0; i <= math.Abs(delta); i++ {
			instructions = append(instructions, d.key(key), ir.Wait(200))
		}

	case "mute":
		instructions = append(instructions, d.key(onkyoHTR380KeyMute), ir.Wait(2000))

	case "input":
		input, ok := args["input"].(string)
		if !ok {
			return nil, oops.BadRequest("invalid input argument 'input': %v", args["input"])
		}

		key, ok := onkyoHTR380InputKeys[input]
		if !ok {
			return nil, oops.InternalService("failed to find key for input option %q", input)
		}

		instructions = append(instructions, d.key(key), ir.Wait(2000))
	}

	return instructions, nil
}

//func (d *OnkyoHTR380) SetProperties(state map[string]interface{}) (bool, error) {
//	if err := device.ValidateState(state, d.ToDef().State); err != nil {
//		return false, err
//	}
//
//	changed := false
//
//	power, ok := state["power"].(bool)
//	if ok {
//		if d.power != power {
//			changed = true
//		}
//
//		d.power = power
//	}
//
//	return changed, nil
//}

func (d *OnkyoHTR380) ToDef() *devicedef.Device {
	return &devicedef.Device{
		Id:             d.ID(),
		Name:           d.Name,
		Type:           d.Type,
		Kind:           d.Kind,
		ControllerName: d.ControllerName,
		Attributes:     d.Attributes,
		StateProviders: d.StateProviders,
		State: map[string]*devicedef.Property{
			"power": device.BoolProperty(d.power),
		},
		Commands: map[string]*devicedef.Command{
			"volume": {
				Args: map[string]*devicedef.Arg{
					"delta": device.IntArg(-10, 10, true),
				},
			},
			"mute": {},
			"input": {
				Args: map[string]*devicedef.Arg{
					"input": device.StringArgWithOptions([]*devicedef.Option{
						{onkyoHTR380InputBDDVD, "Laptop"},
						{onkyoHTR380InputCBLSAT, "Roku"},
						{onkyoHTR380InputGAME, "Game"},
					}, true),
				},
			},
		},
	}
}

func (d *OnkyoHTR380) Copy() Device {
	return &OnkyoHTR380{
		DeviceHeader: d.DeviceHeader,
		power:        d.power,
	}
}
