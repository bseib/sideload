
## this file can be sourced into scripts to have access to version functions

function isVersionPattern {
  local __resultvar=$1
  local VER1=$2
  local ISMATCH=0
  [[ $VER1 =~ ^([0-9]+)\.([0-9]+)(\.([0-9]+))?(\.([0-9]+))?(-.+)? ]] && ISMATCH=1
  if [ $ISMATCH -eq 0 ]; then
    #echo "isVersionPattern(${VER1}) => 0"
    # Can't parse version string $VER1
    eval $__resultvar="'0'"; return 0
  fi
  #echo "isVersionPattern(${VER1}) => 1"
  eval $__resultvar="'1'"; return 0
}

##-----------------------------------------------------------------------------
## ($MAJOR,$MINOR,$REVISION,$BUILD,$SUFFIX) = versParseVersionString($VERS)
##
## Parse the $VERS string into its components and return them into a list of
## variables. If a revision or build digit is not present, or a suffix is not
## present, the corresponding returned variable will be empty.
##
## Examples:
##   versParseVersionString("2.10") ==> ( '2', '10', '', '', '' )
##   versParseVersionString("2.10.4") ==> ( '2', '10', '4', '', '' )
##   versParseVersionString("2.10.4-SNAPSHOT") ==> ( '2', '10', '4', '', '-SNAPSHOT' )
##   versParseVersionString("2.10.4.5") ==> ( '2', '10', '4', '5', '' )
##-----------------------------------------------------------------------------
function versParseVersionString {
  local VER1=$1
  local __resultMajor=$2
  local __resultMinor=$3
  local __resultRevision=$4
  local __resultBuild=$5
  local __resultSuffix=$6
  local ISMATCH=0
  [[ $VER1 =~ ^([0-9]+)\.([0-9]+)(\.([0-9]+))?(\.([0-9]+))?(-.+)? ]] && ISMATCH=1
  if [ $ISMATCH -eq 0 ]; then
    echo "Can't parse version string $VER1"
    exit 1
  fi

  eval $__resultMajor="'${BASH_REMATCH[1]}'";
  eval $__resultMinor="'${BASH_REMATCH[2]}'";
  eval $__resultRevision="'${BASH_REMATCH[4]}'";
  eval $__resultBuild="'${BASH_REMATCH[6]}'";
  eval $__resultSuffix="'${BASH_REMATCH[7]}'";
  return 0
}

##-----------------------------------------------------------------------------
## $BUMPED_VERS=versBumpVersion($FROM_VERS, $VBUMP, $VSUFFIX)
##
## FROM_VERS is known to parse as a version string
## VBUMP is in (+major|+minor|+revision|+build)
## VSUFFIX is in (+snapshot|-snapshot|+milestone|-milestone)
##-----------------------------------------------------------------------------
function versBumpVersion {
  local __resultvar=$1
  local _FROM_VERS=$2
  local _VBUMP=$3
  local _VSUFFIX=$4
  local _MAJOR
  local _MINOR
  local _REVISION
  local _BUILD
  local _SUFFIX

  ## break the version string into its components
  versParseVersionString $_FROM_VERS _MAJOR _MINOR _REVISION _BUILD _SUFFIX

  ## update the appropriate number slot and maybe cascade zeros
  local VCMD=${VBUMP}
  if [ "+major" = "${VCMD}" ]; then
    [ "x" = "x${_MAJOR}" ] && _MAJOR=0
    _MAJOR=$(($_MAJOR + 1))
    VCMD="0minor"
  fi
  if [ "0minor" = "${VCMD}" ] && [ "x" != "x${_MINOR}" ]; then
    _MINOR="0"
    VCMD="0revision"
  fi
  if [ "+minor" = "${VCMD}" ]; then
    [ "x" = "x${_MINOR}" ] && _MINOR=0
    _MINOR=$(($_MINOR + 1))
    VCMD="0revision"
  fi
  if [ "0revision" = "${VCMD}" ] && [ "x" != "x${_REVISION}" ]; then
    _REVISION="0"
    VCMD="0build"
  fi
  if [ "+revision" = "${VCMD}" ]; then
    [ "x" = "x${_REVISION}" ] && _REVISION=0
    _REVISION=$(($_REVISION + 1))
    VCMD="0build"
  fi
  if [ "0build" = "${VCMD}" ] && [ "x" != "x${_BUILD}" ]; then
    _BUILD="0"
  fi
  if [ "+build" = "${VCMD}" ]; then
    [ "x" = "x${_BUILD}" ] && _BUILD=0
    _BUILD=$(($_BUILD + 1))
  fi

  ## alter the suffix accordingly
  if [ "+snapshot" = "${_VSUFFIX}" ]; then
    _SUFFIX="-SNAPSHOT"
  fi
  if [ "-snapshot" = "${_VSUFFIX}" ]; then
    _SUFFIX=""
  fi
  if [ "+milestone" = "${_VSUFFIX}" ]; then
    _SUFFIX="-MILESTONE"
  fi
  if [ "-milestone" = "${_VSUFFIX}" ]; then
    _SUFFIX=""
  fi

  ## tack on the dots if there is a number for that slot
  if [ "x" != "x${_MINOR}" ]; then
    _MINOR=".${_MINOR}"
  fi
  if [ "x" != "x${_REVISION}" ]; then
    _REVISION=".${_REVISION}"
  fi
  if [ "x" != "x${_BUILD}" ]; then
    _BUILD=".${_BUILD}"
  fi

  ## build the new version string
  local NEWVERS="${_MAJOR}${_MINOR}${_REVISION}${_BUILD}${_SUFFIX}"
  eval $__resultvar="'${NEWVERS}'";
}

function sortVersions {
  local  __resultvar=$1
  local VER1=$2
  local VER2=$3

  ## handle cases where one or the other is empty or doesn't parse
  local _ISMATCH1=0
  local _ISMATCH2=0
  isVersionPattern _ISMATCH1 $VER1
  isVersionPattern _ISMATCH2 $VER2
  #echo "_ISMATCH1=${_ISMATCH1} _ISMATCH2=${_ISMATCH2}"
  if [ $_ISMATCH1 -eq 1 ] && [ $_ISMATCH2 -eq 0 ]; then
    ## choose $VER1
    eval $__resultvar="'$VER1'"; return 0
  fi
  if [ $_ISMATCH1 -eq 0 ] && [ $_ISMATCH2 -eq 1 ]; then
    ## choose $VER2
    eval $__resultvar="'$VER2'"; return 0
  fi

  local __topIndex
  indexVersions __topIndex $VER1 $VER2
  #echo "__topIndex=${__topIndex} VER1=${VER1} VER2=${VER2}"
  if [[ ${__topIndex} -lt 0 ]]; then
    eval $__resultvar="'$VER2'"; return 0
  else
    eval $__resultvar="'$VER1'"; return 0
  fi
}

function indexVersions {
	local  __resultvar=$1
	local VER1=$2
	local VER2=$3

	local ISMATCH=0
	[[ $VER1 =~ ^([0-9]+)\.([0-9]+)(\.([0-9]+))?(\.([0-9]+))?(-.+)? ]] && ISMATCH=1
	if [ $ISMATCH -eq 0 ]; then
		echo "Can't parse version string $VER1"
		exit 1
	fi
	local MAJOR1=${BASH_REMATCH[1]}
	local MINOR1=${BASH_REMATCH[2]}
	local REVISION1=${BASH_REMATCH[4]}
	local BUILD1=${BASH_REMATCH[6]}
	[[ "x" = "x${MAJOR1}" ]] && MAJOR1=0
	[[ "x" = "x${MINOR1}" ]] && MINOR1=0
	[[ "x" = "x${REVISION1}" ]] && REVISION1=0
	[[ "x" = "x${BUILD1}" ]] && BUILD1=0

	ISMATCH=0
	[[ $VER2 =~ ^([0-9]+)\.([0-9]+)(\.([0-9]+))?(\.([0-9]+))?(-.+)? ]] && ISMATCH=1
	if [ $ISMATCH -eq 0 ]; then
		echo "Can't parse version string $VER2"
		exit 1
	fi
	local MAJOR2=${BASH_REMATCH[1]}
	local MINOR2=${BASH_REMATCH[2]}
	local REVISION2=${BASH_REMATCH[4]}
	local BUILD2=${BASH_REMATCH[6]}
	[[ "x" = "x${MAJOR2}" ]] && MAJOR2=0
	[[ "x" = "x${MINOR2}" ]] && MINOR2=0
	[[ "x" = "x${REVISION2}" ]] && REVISION2=0
	[[ "x" = "x${BUILD2}" ]] && BUILD2=0

	if [[ "$MAJOR1" -gt "$MAJOR2" ]]; then
		eval $__resultvar="'1'"; return 0
	fi
	if [[ "$MAJOR1" -lt "$MAJOR2" ]]; then
		eval $__resultvar="'-1'"; return 0
	fi
	if [[ "$MINOR1" -gt "$MINOR2" ]]; then
		eval $__resultvar="'1'"; return 0
	fi
	if [[ "$MINOR1" -lt "$MINOR2" ]]; then
		eval $__resultvar="'-1'"; return 0
	fi
	if [[ "$REVISION1" -gt "$REVISION2" ]]; then
		eval $__resultvar="'1'"; return 0
	fi
	if [[ "$REVISION1" -lt "$REVISION2" ]]; then
		eval $__resultvar="'-1'"; return 0
	fi
	if [[ "$BUILD1" -gt "$BUILD2" ]]; then
		eval $__resultvar="'1'"; return 0
	fi
	if [[ "$BUILD1" -lt "$BUILD2" ]]; then
		eval $__resultvar="'-1'"; return 0
	fi
	eval $__resultvar="'0'"; return 0
}
