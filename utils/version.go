package utils

import (
	"github.com/Masterminds/semver"
	"github.com/projectriri/bot-gateway/types"
	log "github.com/sirupsen/logrus"
	"strings"
)

func CheckIfVersionSatisfy(v string, constraint string) bool {
	v1, err := semver.NewVersion(v)
	if err != nil {
		log.Errorf("bad version %s", v)
		return false
	}
	// Constraints example.
	constraints, err := semver.NewConstraint(constraint)
	if err != nil {
		log.Errorf("bad version constraint %s", constraint)
		return false
	}
	return constraints.Check(v1)
}

func CheckIfFormatSatisfy(fmt types.Format, fmtConstraint types.Format) bool {
	return strings.ToLower(fmt.API) == strings.ToLower(fmtConstraint.API) &&
		strings.ToLower(fmt.Method) == strings.ToLower(fmtConstraint.Method) &&
		strings.ToLower(fmt.Protocol) == strings.ToLower(fmtConstraint.Protocol) &&
		CheckIfVersionSatisfy(fmt.Version, fmtConstraint.Version)
}
