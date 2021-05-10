package db

import (
	"os"
	"testing"

	"github.com/kotakanbe/go-cve-dictionary/models"
)

func TestMain(m *testing.M) {
	// log.Initialize("/tmp", true, os.Stderr)
	code := m.Run()
	os.Exit(code)
}

// move to rdb_test.go
func TestMakeVersionConstraint(t *testing.T) {
	var testdata = []struct {
		cpe        models.Cpe
		constraint string
	}{
		{
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					CpeWFN: models.CpeWFN{
						Part:            "a",
						Vendor:          "cisco",
						Product:         "node-jose",
						Version:         "",
						Update:          "",
						Edition:         "",
						Language:        "",
						SoftwareEdition: "",
						TargetSW:        "",
						TargetHW:        "",
						Other:           "",
					},
					VersionStartExcluding: "",
					VersionStartIncluding: "",
					VersionEndExcluding:   "0.11.0",
					VersionEndIncluding:   "",
				},
			},
			constraint: "< 0.11.0",
		},
		{
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					CpeWFN: models.CpeWFN{
						Part:    "a",
						Vendor:  "cisco",
						Product: "node-jose",
					},
					VersionEndIncluding: "0.11.0",
				},
			},
			constraint: "<= 0.11.0",
		},
		{
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					CpeWFN: models.CpeWFN{
						Part:    "a",
						Vendor:  "cisco",
						Product: "node-jose",
					},
					VersionStartExcluding: "0.10.0",
					VersionEndIncluding:   "0.11.0",
				},
			},
			constraint: "> 0.10.0, <= 0.11.0",
		},
		{
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					CpeWFN: models.CpeWFN{
						Part:    "a",
						Vendor:  "cisco",
						Product: "node-jose",
					},
					VersionStartIncluding: "0.10.0",
					VersionEndExcluding:   "0.11.0",
				},
			},
			constraint: ">= 0.10.0, < 0.11.0",
		},
		{
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					CpeWFN: models.CpeWFN{
						Part:    "a",
						Vendor:  "cisco",
						Product: "node-jose",
					},
					VersionStartIncluding: "",
					VersionEndExcluding:   "",
				},
			},
			constraint: "",
		},
	}

	for i, tt := range testdata {
		constraint := makeVersionConstraint(tt.cpe)
		if tt.constraint != constraint {
			t.Errorf("[%d] expected %s, actual %s", i, tt.constraint, constraint)
		}
	}
}

// move to rdb_test.go
func TestMatch(t *testing.T) {
	var testdata = []struct {
		uri   string
		cpe   models.Cpe
		match bool
		err   bool
	}{
		//0
		{
			uri: "cpe:/a:oracle:vm_virtualbox:5.1.1",
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					URI: "cpe:/a:oracle:vm_virtualbox",
					CpeWFN: models.CpeWFN{
						Part:            "a",
						Vendor:          "oracle",
						Product:         "vm_virtualbox",
						Version:         "",
						Update:          "",
						Edition:         "",
						Language:        "",
						SoftwareEdition: "",
						TargetSW:        "",
						TargetHW:        "",
						Other:           "",
					},
					VersionStartIncluding: "5.1.0",
					VersionEndExcluding:   "5.1.32",
				},
			},
			match: true,
		},
		//1
		{
			uri: "cpe:/a:oracle:vm_virtualbox:5.0.9",
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					URI: "cpe:/a:oracle:vm_virtualbox",
					CpeWFN: models.CpeWFN{
						Part:    "a",
						Vendor:  "oracle",
						Product: "vm_virtualbox",
					},
					VersionStartIncluding: "5.1.0",
					VersionEndExcluding:   "5.1.32",
				},
			},
			match: false,
		},
		//2
		{
			uri: "cpe:/a:oracle:vm_virtualbox:5.1.0",
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					URI: "cpe:/a:oracle:vm_virtualbox",
					CpeWFN: models.CpeWFN{
						Part:    "a",
						Vendor:  "oracle",
						Product: "vm_virtualbox",
					},
					VersionStartIncluding: "5.1.0",
					VersionEndExcluding:   "5.1.32",
				},
			},
			match: true,
		},
		//3
		{
			uri: "cpe:/a:oracle:vm_virtualbox:5.2.32",
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					URI: "cpe:/a:oracle:vm_virtualbox",
					CpeWFN: models.CpeWFN{
						Part:    "a",
						Vendor:  "oracle",
						Product: "vm_virtualbox",
					},
					VersionStartIncluding: "5.1.0",
					VersionEndExcluding:   "5.1.32",
				},
			},
			match: false,
		},
		//4
		{
			uri: "cpe:/a:oracle:vm_virtualbox:5.1.31",
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					URI: "cpe:/a:oracle:vm_virtualbox",
					CpeWFN: models.CpeWFN{
						Part:    "a",
						Vendor:  "oracle",
						Product: "vm_virtualbox",
					},
					VersionStartIncluding: "5.1.0",
					VersionEndExcluding:   "5.1.32",
				},
			},
			match: true,
		},
		//5
		{
			uri: "cpe:/a:oracle:vm_virtualbox:5.1.31",
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					URI: "cpe:/a:oracle:vm_virtualbox:5.1.31",
					CpeWFN: models.CpeWFN{
						Part:    "a",
						Vendor:  "oracle",
						Product: "vm_virtualbox",
						Version: "5.1.31",
					},
				},
			},
			match: true,
		},
		//6 superset
		{
			uri: "cpe:/o:microsoft:windows_7",
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					URI: "cpe:/o:microsoft:windows_7::sp1",
					CpeWFN: models.CpeWFN{
						Part:    "o",
						Vendor:  "microsoft",
						Product: "windows_7",
						Version: "",
						Update:  "sp1",
					},
				},
			},
			match: true,
		},
		//7 subset
		{
			uri: "cpe:/o:microsoft:windows_7::sp2",
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					URI: "cpe:/o:microsoft:windows_7",
					CpeWFN: models.CpeWFN{
						Part:    "o",
						Vendor:  "microsoft",
						Product: "windows_7",
					},
				},
			},
			match: true,
		},
		{
			uri: "cpe:/o:microsoft:windows_10",
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					URI: "cpe:/o:microsoft:windows_10:-",
					CpeWFN: models.CpeWFN{
						Part:    "o",
						Vendor:  "microsoft",
						Product: "windows_10",
						Version: "-",
					},
				},
			},
			match: true,
		},
		//9 Cisco Versioning
		{
			uri: "cpe:/a:cisco:ios:15.2%282%29eb",
			cpe: models.Cpe{
				CpeBase: models.CpeBase{
					URI: "cpe:/a:cisco:ios:15.2%282%29ec",
					CpeWFN: models.CpeWFN{
						Part:    "a",
						Vendor:  "cisco",
						Product: "ios",
						Version: `15\.2\(2\)ec`,
					},
				},
			},
			match: false,
			err:   true,
		},
	}

	for i, tt := range testdata {
		match, err := match(tt.uri, tt.cpe)
		if !tt.err && err != nil {
			t.Errorf("[%d] err: %s", i, err)
		}

		if tt.match != match {
			t.Errorf("[%d] expected: %t, actual: %t", i, tt.match, match)
		}
	}
}

func TestMatchProductVendor(t *testing.T) {
	var testdata = []struct {
		uri     string
		affects []models.Affect
		match   bool
	}{
		{
			uri: "cpe:/o:cisco:nx-os:6.0%282%29a6%288%29",
			affects: []models.Affect{
				{
					Vendor:  "cisco",
					Product: "nx-os",
					Version: "6.0(2)a6(8)",
				},
			},
			match: true,
		},
		{
			uri: "cpe:/o:cisco:nx-os:6.0%282%29a6%288%29",
			affects: []models.Affect{
				{
					Vendor:  "cisco",
					Product: "nx-os",
					Version: "7.0(2)a6(8)",
				},
			},
			match: false,
		},
	}
	for i, tt := range testdata {
		match, err := matchExactByAffects(tt.uri, tt.affects)
		if err != nil {
			t.Errorf("[%d] err: %s", i, err)
		}

		if tt.match != match {
			t.Errorf("[%d] expected: %t, actual: %t", i, tt.match, match)
		}
	}
}
