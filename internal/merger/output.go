package merger

import (
	"encoding/json"
	"fmt"
	"os"
)

func (h *Handler) printJson(orgInfo OrganisationInformation) {
	h.log.Debugf("Printing JSON to screen")
	// convert to json
	j, err := json.Marshal(orgInfo)
	if err != nil {
		panic(err)
	}
	// save to disk using org name
	h.log.Debugf("Saving JSON to disk")
	fmt.Println(orgInfo.Organisation.Name)
	err = h.saveJson(j, orgInfo.Organisation.Name)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) saveJson(json []byte, orgName string) error {
	h.log.Debugf("Saving to file as %s", orgName)
	return os.WriteFile(fmt.Sprintf("%s.json", orgName), json, 0644)
}
