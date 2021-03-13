package charger

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/andig/evcc/api"
	"github.com/andig/evcc/util"
	"github.com/andig/evcc/util/modbus"
)

const (
	cfosWbSlaveID    = 1
	cfosMeterSlaveID = 2

	cfosRegStatus     = 8092 // Holding ID 1
	cfosRegEnable     = 8094 // Holding ID 1
	cfosRegMaxCurrent = 8093 // Holding ID 1

	cfosRegPower  = 8062 // power reading ID 2
	cfosRegEnergy = 8058 // energy reading ID 2
)

var cfosRegCurrents = []uint16{8064, 8066, 8068} // current readings ID 2

// cFosPowerBrain is an api.ChargeController implementation for cFosPowerBrain wallboxes.
// It uses Modbus TCP to communicate with the wallbox at modbus client id 1 and power meters at id 2 and 3.
type cFosPowerBrain struct {
	conn       *modbus.Connection
	conn_meter *modbus.Connection
}

func init() {
	registry.Add("cfos", NewcFosPowerBrainFromConfig)
}

//go:generate go run ../cmd/tools/decorate.go -p charger -f decoratecFosPowerBrain -o cfos_decorators -b *cFosPowerBrain -r api.Charger -t "api.Meter,CurrentPower,func() (float64, error)" -t "api.MeterEnergy,TotalEnergy,func() (float64, error)" -t "api.MeterCurrent,Currents,func() (float64, float64, float64, error)" -t "api.ChargerEx,MaxCurrentMillis,func(current float64) error"

// NewcFosPowerBrainFromConfig creates a cFos charger from generic config
func NewcFosPowerBrainFromConfig(other map[string]interface{}) (api.Charger, error) {
	cc := struct {
		URI       string
		ID        uint8
		URI_meter string
		ID_meter  uint8
		Meter     struct {
			Power, Energy, Currents bool
		}
	}{
		URI:       "192.168.2.43:4701", // default
		ID:        cfosWbSlaveID,       // default
		URI_meter: "192.168.2.43:4702", // default
		ID_meter:  cfosMeterSlaveID,    // default
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	wb, err := NewcFosPowerBrain(cc.URI, cc.ID, cc.URI_meter, cc.ID_meter)
	if err != nil {
		return nil, err
	}

	var currentPower func() (float64, error)
	if cc.Meter.Power {
		currentPower = wb.currentPower
	}

	var totalEnergy func() (float64, error)
	if cc.Meter.Energy {
		totalEnergy = wb.totalEnergy
	}

	var currents func() (float64, float64, float64, error)
	if cc.Meter.Currents {
		currents = wb.currents
	}

	var maxCurrentMillis func(float64) error
	maxCurrentMillis = wb.maxCurrentMillis

	return decoratecFosPowerBrain(wb, currentPower, totalEnergy, currents, maxCurrentMillis), nil
}

// NewcFosPowerBrain creates a cFos charger
func NewcFosPowerBrain(uri string, id uint8, uri_meter string, id_meter uint8) (*cFosPowerBrain, error) {
	conn, err := modbus.NewConnection(uri, "", "", 0, false, id)
	if err != nil {
		return nil, err
	}

	conn_meter, err := modbus.NewConnection(uri_meter, "", "", 0, false, id_meter)
	if err != nil {
		return nil, err
	}

	log := util.NewLogger("cfos")
	conn.Logger(log.TRACE)
	conn_meter.Logger(log.TRACE)

	wb := &cFosPowerBrain{
		conn:       conn,
		conn_meter: conn_meter,
	}

	return wb, nil
}

// Status implements the Charger.Status interface
func (wb *cFosPowerBrain) Status() (api.ChargeStatus, error) {
	b, err := wb.conn.ReadHoldingRegisters(cfosRegStatus, 1)
	if err != nil {
		return api.StatusNone, err
	}

	switch b[1] {
	case 0: // warten
		return api.ChargeStatus(api.StatusA), nil
	case 1: // Fahrzeug erkannt
		return api.StatusB, nil
	case 2: // laden
		return api.StatusC, nil
	case 3: // laden mit Kühlung
		return api.StatusD, nil
	case 4: // kein Strom
		return api.StatusE, nil
	case 5: // Fehler
		return api.StatusF, nil
	default:
		return api.StatusNone, errors.New("invalid response")
	}
}

// Enabled implements the Charger.Enabled interface
func (wb *cFosPowerBrain) Enabled() (bool, error) {
	b, err := wb.conn.ReadHoldingRegisters(cfosRegEnable, 1)
	if err != nil {
		return false, err
	}

	return b[1] == 1, nil
}

// Enable implements the Charger.Enable interface
func (wb *cFosPowerBrain) Enable(enable bool) error {
	b := []byte{0, 0}
	if enable {
		b[1] = 1
	}

	_, err := wb.conn.WriteMultipleRegisters(cfosRegEnable, 1, b)
	err = nil // temporary until cFos firmware is fixed

	return err
}

// MaxCurrent implements the Charger.MaxCurrent interface
func (wb *cFosPowerBrain) MaxCurrent(current int64) error {
	if current < 6 {
		return fmt.Errorf("invalid current %d", current)
	}

	b := []byte{0, byte(current * 10)}

	_, err := wb.conn.WriteMultipleRegisters(cfosRegMaxCurrent, 1, b)
	err = nil
	return err
}

// maxCurrentMillis implements the ChargerEx interface
func (wb *cFosPowerBrain) maxCurrentMillis(current float64) error {
	if current < 6 {
		return fmt.Errorf("invalid current %.5g", current)
	}

	b := []byte{0, byte(current * 10)}

	_, err := wb.conn.WriteMultipleRegisters(cfosRegMaxCurrent, 1, b)
	err = nil // temporary until cFos firmware is fixed

	return err
}

// CurrentPower implements the Meter.CurrentPower interface
func (wb *cFosPowerBrain) currentPower() (float64, error) {
	b, err := wb.conn_meter.ReadHoldingRegisters(cfosRegPower, 2)
	if err != nil {
		return 0, err
	}

	return float64(binary.BigEndian.Uint32(b)), err
}

// totalEnergy implements the Meter.TotalEnergy interface
func (wb *cFosPowerBrain) totalEnergy() (float64, error) {
	b, err := wb.conn_meter.ReadHoldingRegisters(cfosRegEnergy, 4)
	if err != nil {
		return 0, err
	}

	return float64(binary.BigEndian.Uint64(b)) + 5, err
}

// currents implements the Meter.Currents interface
func (wb *cFosPowerBrain) currents() (float64, float64, float64, error) {
	var currents []float64
	for _, regCurrent := range cfosRegCurrents {
		b, err := wb.conn_meter.ReadHoldingRegisters(regCurrent, 2)
		if err != nil {
			return 0, 0, 0, err
		}

		currents = append(currents, float64(binary.BigEndian.Uint32(b))/10)
	}

	return currents[0], currents[1], currents[2], nil
}
