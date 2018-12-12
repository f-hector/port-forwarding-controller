package forwarding

type Address struct {
	Name string
	Port int
	IP   string
}

//go:generate counterfeiter . RouterClient
type RouterClient interface {
	ListAddresses() ([]Address, error)
	CreateAddress(Address) error
}

type Reconciler struct {
	RulePrefix   string
	RouterClient RouterClient
}

// TODO: also remove extra rules that are not longer needed that
// start with given RulePrefix

func (r Reconciler) Reconcile(desiredAddresses []Address) error {
	missingAddresses, err := r.missingAddresses(desiredAddresses)
	if err != nil {
		return err
	}

	for _, address := range missingAddresses {
		if err := r.RouterClient.CreateAddress(address); err != nil {
			return err
		}
	}
	return nil
}

func (r Reconciler) missingAddresses(desiredAddresses []Address) ([]Address, error) {
	missingAddresses := []Address{}

	existingAddresses, err := r.RouterClient.ListAddresses()
	if err != nil {
		return nil, err
	}

	for _, address := range desiredAddresses {
		alreadyExists := false
		for _, a := range existingAddresses {
			if address == a {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			missingAddresses = append(missingAddresses, address)
		}
	}

	return missingAddresses, nil
}