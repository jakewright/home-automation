package repository

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/oops"
	deviceregistrydef "github.com/jakewright/home-automation/services/device-registry/def"
	"github.com/jakewright/home-automation/services/dmx/domain"
)

// FixtureRepository holds an in-memory collection of fixtures
type FixtureRepository struct {
	fixtures map[string]domain.Fixture
}

// Init loads devices from the device registry and populates a new repository
func Init(
	ctx context.Context,
	serviceName string,
	deviceRegistry deviceregistrydef.DeviceRegistryService,
) (*FixtureRepository, error) {
	// Load devices from the registry
	rsp, err := deviceRegistry.ListDevices(ctx, &deviceregistrydef.ListDevicesRequest{
		ControllerName: &serviceName,
	}).Wait()
	if err != nil {
		return nil, oops.WithMessage(err, "failed to fetch devices")
	}

	headers, _ := rsp.GetDeviceHeaders()

	fixtures := make([]domain.Fixture, len(headers))

	for i, header := range rsp.DeviceHeaders {
		// Be defensive against the device registry returning the wrong devices
		switch {
		case header.GetControllerName() != serviceName:
			return nil, oops.InternalService("device %s is not for this controller", header.GetId())
		}

		fixture, err := domain.NewFixture(header)
		if err != nil {
			return nil, oops.WithMessage(err, "failed to create fixture")
		}

		fixtures[i] = fixture
	}

	if err := validate(fixtures); err != nil {
		return nil, err
	}

	return New(fixtures...), nil
}

// New returns a repository holding the given fixtures
func New(fixtures ...domain.Fixture) *FixtureRepository {
	m := make(map[string]domain.Fixture, len(fixtures))
	for _, f := range fixtures {
		m[f.ID()] = f
	}
	return &FixtureRepository{
		fixtures: m,
	}
}

// Find returns the fixture with the specified ID or nil if it doesn't exist
func (r *FixtureRepository) Find(id string) domain.Fixture {
	if f, ok := r.fixtures[id]; ok {
		return f
	}

	return nil
}

func validate(fs []domain.Fixture) error {
	m := map[domain.UniverseNumber][]domain.Fixture{}

	for _, f := range fs {
		m[f.UniverseNumber()] = append(m[f.UniverseNumber()], f)
	}

	for _, fs := range m {
		if err := domain.ValidateFixtures(fs); err != nil {
			return err
		}
	}

	return nil
}
