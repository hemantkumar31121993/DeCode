package types

type address struct {
	House    string
	Street   string
	Landmark string
	Pincode  int
	_        int8
}

func NewAddress(house, street, landmark string, pincode int) address {
	return address{House: house, Street: street, Landmark: landmark, Pincode: pincode}
}
