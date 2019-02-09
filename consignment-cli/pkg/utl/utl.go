package utl

import (
	"encoding/json"
	"io/ioutil"

	pb "github.com/alamin-mahamud/container-management-system/consignment-service/proto/consignment"
)

// ParseFile - ...
func ParseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, err
	}
	return consignment, err
}
