package fetcher

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/knqyf263/go-cpe/common"
	"github.com/knqyf263/go-cpe/naming"
	log "github.com/kotakanbe/go-cve-dictionary/log"
	"github.com/kotakanbe/go-cve-dictionary/models"
)

// ParseCpeURI parses cpe22uri and set to models.CpeBase
func ParseCpeURI(uri string) (*models.CpeBase, error) {
	var wfn common.WellFormedName
	var err error
	if strings.HasPrefix(uri, "cpe:/") {
		val := strings.TrimPrefix(uri, "cpe:/")
		if strings.Contains(val, "/") {
			uri = "cpe:/" + strings.Replace(val, "/", `\/`, -1)
		}
		wfn, err = naming.UnbindURI(uri)
		if err != nil {
			return nil, err
		}
	} else {
		wfn, err = naming.UnbindFS(uri)
		if err != nil {
			return nil, err
		}
	}

	return &models.CpeBase{
		URI:             naming.BindToURI(wfn),
		FormattedString: naming.BindToFS(wfn),
		WellFormedName:  wfn.String(),
		CpeWFN: models.CpeWFN{
			Part:            fmt.Sprintf("%s", wfn.Get(common.AttributePart)),
			Vendor:          fmt.Sprintf("%s", wfn.Get(common.AttributeVendor)),
			Product:         fmt.Sprintf("%s", wfn.Get(common.AttributeProduct)),
			Version:         fmt.Sprintf("%s", wfn.Get(common.AttributeVersion)),
			Update:          fmt.Sprintf("%s", wfn.Get(common.AttributeUpdate)),
			Edition:         fmt.Sprintf("%s", wfn.Get(common.AttributeEdition)),
			Language:        fmt.Sprintf("%s", wfn.Get(common.AttributeLanguage)),
			SoftwareEdition: fmt.Sprintf("%s", wfn.Get(common.AttributeSwEdition)),
			TargetSW:        fmt.Sprintf("%s", wfn.Get(common.AttributeTargetSw)),
			TargetHW:        fmt.Sprintf("%s", wfn.Get(common.AttributeTargetHw)),
			Other:           fmt.Sprintf("%s", wfn.Get(common.AttributeOther)),
		},
	}, nil
}

// StringToFloat cast string to float64
func StringToFloat(str string) float64 {
	if len(str) == 0 {
		return 0
	}
	var f float64
	var ignorableError error
	if f, ignorableError = strconv.ParseFloat(str, 64); ignorableError != nil {
		log.Errorf("Failed to cast CVSS score. score: %s, err; %s",
			str,
			ignorableError,
		)
		f = 0
	}
	return f
}
