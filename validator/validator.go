package validator

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/reugn/go-quartz/quartz"
	"gopkg.in/ldap.v2"
)

func LowerCase(value interface{}, key cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	m := value.(string)
	low := strings.ToLower(m)

	if m != low {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid value",
			Detail:   fmt.Sprintf("%s should be lowercase", key),
		})
	}

	return diags
}

func Cron(value interface{}, _ cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	err := quartz.ValidateCronExpression(value.(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Cron expression",
			Detail:   fmt.Sprintf("%s is not a valid cron: %s", value, err),
		})
	}

	return diags
}

func CronLength(value interface{}, _ cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	cron := strings.Split(value.(string), " ")

	if len(cron) < 6 || len(cron) > 7 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Cron expression, value should be between 6 and 7 parts long",
			Detail:   fmt.Sprintf("%s is not a valid cron", value),
		})
	}

	return diags
}

var CommaSeperatedList = validation.ToDiagFunc(
	validation.StringMatch(regexp.MustCompile(`.+(?:,.+)*`), "must be comma separated string"),
)

var RepoKey = validation.AllDiag(
	validation.ToDiagFunc(
		validation.StringLenBetween(1, 64),
	),
	validation.ToDiagFunc(
		validation.StringDoesNotContainAny(" !@#$%^&*()+={}[]:;<>,/?~`|\\"),
	),
)

var ProjectKey = validation.ToDiagFunc(
	validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9\-]{1,31}$`), "project_key must be 2 - 32 lowercase alphanumeric and hyphen characters"),
)

var validLicenseTypes = []string{
	"0BSD",
	"AAL",
	"Abstyles",
	"Adobe-2006",
	"Adobe-Glyph",
	"ADSL",
	"AFL-1.1",
	"AFL-1.2",
	"AFL-2.0",
	"AFL-2.1",
	"AFL-3.0",
	"Afmparse",
	"AGPL-1.0",
	"AGPL-3.0",
	"AGPL-3.0-only",
	"AGPL-3.0-or-later",
	"Aladdin",
	"AMDPLPA",
	"AML",
	"AMPAS",
	"ANTLR-PD",
	"Apache-1.0",
	"Apache-1.1",
	"Apache-2.0",
	"APAFML",
	"APL-1.0",
	"APSL-1.0",
	"APSL-1.1",
	"APSL-1.2",
	"APSL-2.0",
	"Artistic-1.0",
	"Artistic-1.0-cl8",
	"Artistic-1.0-Perl",
	"Artistic-2.0",
	"Atlassian End User License Agreement 3.0",
	"Attribution",
	"Bahyph",
	"Barr",
	"Beerware",
	"BitTorrent-1.0",
	"BitTorrent-1.1",
	"Borceux",
	"Bouncy-Castle",
	"BSD",
	"BSD 2-Clause",
	"BSD 3-Clause",
	"BSD-1-Clause",
	"BSD-2-Clause",
	"BSD-2-Clause-FreeBSD",
	"BSD-2-Clause-NetBSD",
	"BSD-2-Clause-Patent",
	"BSD-3-Clause",
	"BSD-3-Clause-Attribution",
	"BSD-3-Clause-Clear",
	"BSD-3-Clause-LBNL",
	"BSD-3-Clause-No-Nuclear-License",
	"BSD-3-Clause-No-Nuclear-License-2014",
	"BSD-3-Clause-No-Nuclear-Warranty",
	"BSD-4-Clause",
	"BSD-4-Clause-UC",
	"BSD-Protection",
	"BSD-Source-Code",
	"BSL-1.0",
	"bzip2-1.0.5",
	"bzip2-1.0.6",
	"CA-TOSL-1.1",
	"Caldera",
	"CATOSL-1.1",
	"CC-BY-1.0",
	"CC-BY-2.0",
	"CC-BY-2.5",
	"CC-BY-3.0",
	"CC-BY-4.0",
	"CC-BY-NC-1.0",
	"CC-BY-NC-2.0",
	"CC-BY-NC-2.5",
	"CC-BY-NC-3.0",
	"CC-BY-NC-4.0",
	"CC-BY-NC-ND-1.0",
	"CC-BY-NC-ND-2.0",
	"CC-BY-NC-ND-2.5",
	"CC-BY-NC-ND-3.0",
	"CC-BY-NC-ND-4.0",
	"CC-BY-NC-SA-1.0",
	"CC-BY-NC-SA-2.0",
	"CC-BY-NC-SA-2.5",
	"CC-BY-NC-SA-3.0",
	"CC-BY-NC-SA-4.0",
	"CC-BY-ND-1.0",
	"CC-BY-ND-2.0",
	"CC-BY-ND-2.5",
	"CC-BY-ND-3.0",
	"CC-BY-ND-4.0",
	"CC-BY-SA-1.0",
	"CC-BY-SA-2.0",
	"CC-BY-SA-2.5",
	"CC-BY-SA-3.0",
	"CC-BY-SA-4.0",
	"CC0-1.0",
	"CCAG-2.5",
	"CDDL-1.0",
	"CDDL-1.1",
	"CDLA-Permissive-1.0",
	"CDLA-Sharing-1.0",
	"CeCILL-1",
	"CECILL-1.0",
	"CECILL-1.1",
	"CeCILL-2",
	"CECILL-2.0",
	"CECILL-2.1",
	"CeCILL-2.1",
	"CeCILL-B",
	"CECILL-B",
	"CeCILL-C",
	"CECILL-C",
	"ClArtistic",
	"CNRI-Jython",
	"CNRI-Python",
	"CNRI-Python-GPL-Compatible",
	"Codehaus",
	"Condor-1.1",
	"Copyfree",
	"CPAL-1.0",
	"CPL-1.0",
	"CPOL-1.02",
	"Crossword",
	"CrystalStacker",
	"CUA-OPL-1.0",
	"CUAOFFICE-1.0",
	"Cube",
	"curl",
	"D-FSL-1.0",
	"Day",
	"Day-Addendum",
	"diffmark",
	"DOC",
	"Dotseqn",
	"DSDP",
	"dvipdfm",
	"ECL-1.0",
	"ECL-2.0",
	"ECL2",
	"eCos-2.0",
	"EFL-1.0",
	"EFL-2.0",
	"eGenix",
	"Eiffel-2.0",
	"Entessa",
	"Entessa-1.0",
	"EPL-1.0",
	"EPL-2.0",
	"ErlPL-1.1",
	"EUDatagrid",
	"EUDATAGRID",
	"EUPL-1.0",
	"EUPL-1.1",
	"EUPL-1.2",
	"Eurosym",
	"Facebook-Platform",
	"Fair",
	"Frameworx-1.0",
	"FreeImage",
	"FSFAP",
	"FSFUL",
	"FSFULLR",
	"FTL",
	"GFDL-1.1",
	"GFDL-1.1-only",
	"GFDL-1.1-or-later",
	"GFDL-1.2",
	"GFDL-1.2-only",
	"GFDL-1.2-or-later",
	"GFDL-1.3",
	"GFDL-1.3-only",
	"GFDL-1.3-or-later",
	"Giftware",
	"GL2PS",
	"Glide",
	"Glulxe",
	"gnuplot",
	"Go",
	"GPL-1.0",
	"GPL-1.0+",
	"GPL-1.0-only",
	"GPL-1.0-or-later",
	"GPL-2.0",
	"GPL-2.0+",
	"GPL-2.0+CE",
	"GPL-2.0-only",
	"GPL-2.0-or-later",
	"GPL-2.0-with-autoconf-exception",
	"GPL-2.0-with-bison-exception",
	"GPL-2.0-with-classpath-exception",
	"GPL-2.0-with-font-exception",
	"GPL-2.0-with-GCC-exception",
	"GPL-3.0",
	"GPL-3.0+",
	"GPL-3.0-only",
	"GPL-3.0-or-later",
	"GPL-3.0-with-autoconf-exception",
	"GPL-3.0-with-GCC-exception",
	"gSOAP-1.3b",
	"HaskellReport",
	"Historical",
	"HPND",
	"HSQLDB",
	"IBM-pibs",
	"IBMPL-1.0",
	"ICU",
	"IJG",
	"ImageMagick",
	"iMatix",
	"Imlib2",
	"Info-ZIP",
	"Intel",
	"Intel-ACPI",
	"Interbase-1.0",
	"IPA",
	"IPAFont-1.0",
	"IPL-1.0",
	"ISC",
	"IU-Extreme-1.1.1",
	"JA-SIG",
	"JasPer-2.0",
	"JSON",
	"JTA-Specification-1.0.1B",
	"JTidy",
	"LAL-1.2",
	"LAL-1.3",
	"Latex2e",
	"Leptonica",
	"LGPL-2.0",
	"LGPL-2.0+",
	"LGPL-2.0-only",
	"LGPL-2.0-or-later",
	"LGPL-2.1",
	"LGPL-2.1+",
	"LGPL-2.1-only",
	"LGPL-2.1-or-later",
	"LGPL-3.0",
	"LGPL-3.0+",
	"LGPL-3.0-only",
	"LGPL-3.0-or-later",
	"LGPLLR",
	"Libpng",
	"libtiff",
	"LiLiQ-P-1.1",
	"LiLiQ-R-1.1",
	"LiLiQ-Rplus-1.1",
	"LPL-1.0",
	"LPL-1.02",
	"LPPL-1.0",
	"LPPL-1.1",
	"LPPL-1.2",
	"LPPL-1.3a",
	"LPPL-1.3c",
	"Lucent-1.02",
	"MakeIndex",
	"MirOS",
	"MIT",
	"MIT-advertising",
	"MIT-CMU",
	"MIT-enna",
	"MIT-feh",
	"MITNFA",
	"Motosoto",
	"Motosoto-0.9.1",
	"mpich2",
	"MPL-1.0",
	"MPL-1.1",
	"MPL-2.0",
	"MPL-2.0-no-copyleft-exception",
	"MS-ASP-NET-COMPONENT-RTW",
	"MS-ASP-NET-MVC-3-UPDATE-EULA",
	"MS-ASP-NET-WEB-PAGES-2-EULA",
	"MS-DOT-NET-LIBRARY",
	"MS-DOT-NET-LIBRARY-EULA",
	"MS-DOT-NET-LIBRARY-NON-REDISTRIBUTABLE",
	"MS-PL",
	"MS-RL",
	"MS-RSL",
	"MTLL",
	"Multics",
	"Mup",
	"NASA-1.3",
	"Naumen",
	"NAUMEN",
	"NBPL-1.0",
	"NCSA",
	"Net-SNMP",
	"NetCDF",
	"Nethack",
	"Newsletr",
	"NGPL",
	"NLOD-1.0",
	"NLPL",
	"Nokia",
	"Nokia-1.0a",
	"NOSL",
	"NOSL-3.0",
	"Noweb",
	"NPL-1.0",
	"NPL-1.1",
	"NPOSL-3.0",
	"NRL",
	"NTP",
	"Nunit",
	"NUnit-2.6.3",
	"NUnit-Test-Adapter-2.6.3",
	"OCCT-PL",
	"OCLC-2.0",
	"ODbL-1.0",
	"OFL-1.0",
	"OFL-1.1",
	"OGTSL",
	"OLDAP-1.1",
	"OLDAP-1.2",
	"OLDAP-1.3",
	"OLDAP-1.4",
	"OLDAP-2.0",
	"OLDAP-2.0.1",
	"OLDAP-2.1",
	"OLDAP-2.2",
	"OLDAP-2.2.1",
	"OLDAP-2.2.2",
	"OLDAP-2.3",
	"OLDAP-2.4",
	"OLDAP-2.5",
	"OLDAP-2.6",
	"OLDAP-2.7",
	"OLDAP-2.8",
	"OML",
	"Openfont-1.1",
	"Opengroup",
	"OpenLDAP",
	"OpenSSL",
	"OPL-1.0",
	"OSET-PL-2.1",
	"OSL-1.0",
	"OSL-1.1",
	"OSL-2.0",
	"OSL-2.1",
	"OSL-3.0",
	"PDDL-1.0",
	"PHP-3.0",
	"PHP-3.01",
	"Plexus",
	"PostgreSQL",
	"psfrag",
	"psutils",
	"Public Domain",
	"Public Domain - SUN",
	"Python-2.0",
	"Python-2.1.1",
	"Qhull",
	"QPL-1.0",
	"QTPL-1.0",
	"Rdisc",
	"Real-1.0",
	"RHeCos-1.1",
	"RicohPL",
	"RPL-1.1",
	"RPL-1.5",
	"RPSL-1.0",
	"RSA-MD",
	"RSCPL",
	"Ruby",
	"SAX-PD",
	"Saxpath",
	"Scala",
	"SCEA",
	"Sendmail",
	"SGI-B-1.0",
	"SGI-B-1.1",
	"SGI-B-2.0",
	"SimPL-2.0",
	"SISSL",
	"SISSL-1.2",
	"Sleepycat",
	"SMLNJ",
	"SMPPL",
	"SNIA",
	"Spencer-86",
	"Spencer-94",
	"Spencer-99",
	"SPL-1.0",
	"StandardML-NJ",
	"SugarCRM-1.1.3",
	"SUNPublic-1.0",
	"SWL",
	"Sybase-1.0",
	"TCL",
	"TCP-wrappers",
	"TMate",
	"TORQUE-1.1",
	"TOSL",
	"TPL",
	"Unicode-DFS-2015",
	"Unicode-DFS-2016",
	"Unicode-TOU",
	"Unlicense",
	"UoI-NCSA",
	"UPL-1.0",
	"Vim",
	"VIM License",
	"VOSTROM",
	"VovidaPL-1.0",
	"VSL-1.0",
	"W3C",
	"W3C-19980720",
	"W3C-20150513",
	"Watcom-1.0",
	"Wsuipa",
	"WTFPL",
	"wxWindows",
	"X11",
	"Xerox",
	"XFree86-1.1",
	"xinetd",
	"Xnet",
	"xpp",
	"XSkat",
	"YPL-1.0",
	"YPL-1.1",
	"Zed",
	"Zend-2.0",
	"Zimbra-1.3",
	"Zimbra-1.4",
	"ZLIB",
	"Zlib",
	"zlib-acknowledgement",
	"ZPL-1.1",
	"ZPL-2.0",
	"ZPL-2.1",
}

var LicenseType = validation.ToDiagFunc(validation.StringInSlice(validLicenseTypes, false))

func IsEmail(address interface{}, _ cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	_, err := mail.ParseAddress(address.(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Email",
			Detail:   fmt.Sprintf("%s is not a valid address: %s", address, err),
		})
	}

	return diags
}

// All Updated version of the Terraform's original validation func:
// https://github.com/hashicorp/terraform-plugin-sdk/blob/main/helper/validation/meta.go#L32
//
// All returns a SchemaValidateFunc which tests if the provided value
// passes all provided SchemaValidateFunc
func All(validators ...schema.SchemaValidateDiagFunc) schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		var allDiags diag.Diagnostics
		for _, validator := range validators {
			validatorDiags := validator(i, p)
			allDiags = append(allDiags, validatorDiags...)
		}
		return allDiags
	}
}

// StringIsNotEmpty Updated version of the Terraform's original validation func:
// https://github.com/hashicorp/terraform-plugin-sdk/blob/main/helper/validation/strings.go#L14
func StringIsNotEmpty(i interface{}, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	attr := p[len(p)-1].(cty.GetAttrStep)

	v, ok := i.(string)
	if !ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid string",
			Detail:   fmt.Sprintf("expected type of %q to be string", attr.Name),
		})
		return diags
	}

	if v == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid string",
			Detail:   fmt.Sprintf("expected %q to not be an empty string", attr.Name),
		})
		return diags
	}

	return diags
}

// Updated version of the Terraform's original validation func:
// https://github.com/hashicorp/terraform-plugin-sdk/blob/5adf5f1c4bf0aee7bb31d51cd7f016e81a39e3e5/helper/validation/strings.go#L132
//
// StringInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and matches the value of an element in the valid slice
// will test with in lower case if ignoreCase is true
func StringInSlice(ignoreCase bool, valid ...string) schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		v, ok := i.(string)
		attr := p[len(p)-1].(cty.GetAttrStep)

		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid string",
				AttributePath: p,
				Detail:        fmt.Sprintf("expected type of %q to be string", attr.Name),
			})
			return diags
		}

		for _, str := range valid {
			if v == str || (ignoreCase && strings.EqualFold(v, str)) {
				return diags
			}
		}

		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid string",
			AttributePath: p,
			Detail:        fmt.Sprintf("expected %s to be one of %v, got %s", attr.Name, valid, v),
		})
		return diags
	}
}

func StringIsNotURL(value interface{}, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	attr := p[len(p)-1].(cty.GetAttrStep)

	_, errs := validation.IsURLWithHTTPorHTTPS(value, attr.Name)
	if len(errs) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid string",
			AttributePath: p,
			Detail:        fmt.Sprintf("expected %q not to be a valid url, got %v", attr.Name, value),
		})
	}

	return diags
}

// IntAtLeast Updated version of the Terraform's original validation func:
// https://github.com/hashicorp/terraform-plugin-sdk/blob/5adf5f1c4bf0aee7bb31d51cd7f016e81a39e3e5/helper/validation/int.go#L31
//
// IntAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at least min (inclusive)
func IntAtLeast(min int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		v, ok := i.(int)
		attr := p[len(p)-1].(cty.GetAttrStep)

		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid integer",
				AttributePath: p,
				Detail:        fmt.Sprintf("expected type of %q to be integer", attr.Name),
			})
			return diags
		}

		if v < min {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid integer",
				AttributePath: p,
				Detail:        fmt.Sprintf("expected %s to be at least (%d), got %d", attr.Name, min, v),
			})
			return diags
		}

		return diags
	}
}

func LdapDn(value interface{}, _ cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	_, err := ldap.ParseDN(value.(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid LDAP Domain Name",
			Detail:   fmt.Sprintf("%s is not a valid LDAP Domain Name: %s", value, err),
		})
	}

	return diags
}

func LdapFilter(value interface{}, _ cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	_, err := ldap.CompileFilter(value.(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid LDAP Filter",
			Detail:   fmt.Sprintf("%s is not a valid LDAP Filter: %s", value, err),
		})
	}

	return diags
}

// CheckImportState is used in ImportStateCheck if ImportState is set to `true`.
// IdAttribute is the field used in d.SetId() in Create function to set a resource ID.
func CheckImportState(resourceId string, idAttribute string) func(states []*terraform.InstanceState) error {
	return func(states []*terraform.InstanceState) error {
		if len(states) == 0 {
			return fmt.Errorf("no import state")
		}
		instanceState := states[0]
		if instanceState.ID != resourceId {
			return fmt.Errorf("incorrect state ID: %s", instanceState.ID)
		}
		if instanceState.Attributes[idAttribute] != resourceId {
			return fmt.Errorf("incorrect state attribute '%s': %s", idAttribute, instanceState.Attributes[idAttribute])
		}
		return nil
	}
}
