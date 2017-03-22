package mongo

import (
	"errors"

	"github.com/pivotal-cf/brokerapi"
)

type InstanceCreator struct {
	adminService *AdminService
	repository   *Repository
}

func NewInstanceCreator(adminService *AdminService, repository *Repository) *InstanceCreator {
	return &InstanceCreator{
		adminService,
		repository,
	}
}

func (instanceCreator *InstanceCreator) Create(instanceID string, details brokerapi.ProvisionDetails) error {
	// TODO: ATOM
	instanceExists, error := instanceCreator.repository.InstanceExists(instanceID)

	if error != nil {
		return error
	}

	if instanceExists {
		return instanceExistsError(instanceID, details)
	}

	databaseExists, error := instanceCreator.adminService.DatabaseExists(instanceID)

	if error != nil {
		return error
	}

	// ensure the instance is empty
	if databaseExists {
		error := instanceCreator.adminService.DeleteDatabase(instanceID)

		if error != nil {
			return error
		}
	}

	database, error := instanceCreator.adminService.CreateDatabase(instanceID)

	if error != nil {
		return error
	}

	if database == nil {
		return errors.New("Failed to create new DB instance: " + instanceID)
	}

	error = instanceCreator.repository.SaveInstance(instanceID, details)

	if error != nil {
		return error
	}

	return nil
}

func (instanceCreator *InstanceCreator) Destroy(instanceID string, details brokerapi.DeprovisionDetails) error {
	// TODO: ATOM
	instanceExists, error := instanceCreator.repository.InstanceExists(instanceID)

	if error != nil {
		return error
	}

	if !instanceExists {
		return instanceDoesNotExistError(instanceID, details)
	}

	error = instanceCreator.adminService.DeleteDatabase(instanceID)

	if error != nil {
		return error
	}

	error = instanceCreator.repository.DeleteInstance(instanceID, details)

	if error != nil {
		return error
	}

	return nil
}

func instanceExistsError(instanceID string, details brokerapi.ProvisionDetails) error {
	return errors.New("Instance exists, incetanceID: " + instanceID + ", serviceID: " + details.ServiceID)
}

func instanceDoesNotExistError(instanceID string, details brokerapi.DeprovisionDetails) error {
	return errors.New("Instance doesn't exist, incetanceID: " + instanceID + ", serviceID: " + details.ServiceID)
}