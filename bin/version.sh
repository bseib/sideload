#!/usr/bin/env bash

##
## This is the algorithm for bumping to the next revision number. Like 2.10.38 to 2.10.39.
##
## 1. Look at your <project>'s <version> in your project.props file.
## 2. Determine the major.minor version. Ex: if project.props has 2.10.38, then major.minor is 2.10.
## 3. Do a git fetch, to make sure all tags from the repo are brought to your local machine.
## 4. Run git tag -l "2.10.*" --sort=-version:refname, and the first item returned is the newest 2.10 version.
## 5. Add one to the revision of the previous item, e.g. 2.10.38 -> 2.10.39. Call this $VERS.
## 6. Check that your git workspace is clean. (no uncommitted files or unpushed content)
## 7. Assign the new version number $VERS to version property in project.props
## 8. Do git add and git commit of new project.props file.
## 9. Tag the code, and push the tag to the repo. (git tag -a $1 -m "version $VERS" `git rev-parse HEAD`, then git push origin $VERS)
##

## abort if anything fails, even when a function does a 'return' with non-zero
set -e

HERE=$(cd $(dirname "$0"); pwd)
ROOT=${HERE}/../
LIB="${HERE}/lib"
. ${LIB}/props.functions
. ${LIB}/git.functions
. ${LIB}/version.functions.sh

###############################################################################
## FUNCTIONS
###############################################################################

##-----------------------------------------------------------------------------
## parseCommandLine()
##
## After calling this function, the following global env vars will be set:
##   VMODE=(REPORT|BUMP|ADHOC)
##     REPORT: just report the current version in props file and exit.
##     BUMP: increment one of (+major|+minor|+revision|+build) and possibly
##           (+snapshot|-snapshot|+milestone|-milestone) according to the
##           values found in $VBUMP and $VSUFFIX respectively.
##     ADHOC: assign the ad hoc version string found in $VADHOC.
##   VBUMP=(+major|+minor|+revision|+build) bump one of these version digits
##   VSUFFIX=(+snapshot|-snapshot|+milestone|-milestone) add/del suffix
##   VADHOC=<ad hoc version string matching valid version pattern>
##
## If command line is not valid, print error and exit.
##-----------------------------------------------------------------------------
function parseCommandLine {
  local _IS_VALID_VERSION_STRING=""
  VMODE=""
  VBUMP=""
  VSUFFIX=""
  VADHOC=""

  if [ "x" = "x$1" ]; then
    VMODE="REPORT"
    VBUMP="+revision"
    ## exit this function without error
    return 0
  fi
  if [ "x--all" = "x$1" ]; then
    VMODE="REPORT"
    VBUMP="+revision"
    ## exit this function without error
    return 0
  fi

  ## The presence of (+snapshot|-snapshot|+milestone|-milestone) and/or
  ## (+major|+minor|+revision|+build) can come in either order on $1 and $2.
  ## Need to accept both ways.
  case $1 in
    +snapshot|-snapshot|+milestone|-milestone)
      VMODE="BUMP"
      VSUFFIX=$1
      case $2 in
        +major|+minor|+revision|+build)
          VBUMP=$2
          return 0
          ;;
        *)
          return 0
          ;;
      esac
      ;;
    +major|+minor|+revision|+build)
      VMODE="BUMP"
      VBUMP=$1
      case $2 in
        +snapshot|-snapshot|+milestone|-milestone)
          VSUFFIX=$2
          return 0
          ;;
        *)
          return 0
          ;;
      esac
      ;;
    *)
      ## call $_IS_VALID_VERSION_STRING=isVersionPattern() in version.functions
      isVersionPattern _IS_VALID_VERSION_STRING $1
      if [ $_IS_VALID_VERSION_STRING -eq 1 ]; then
        VMODE="ADHOC"
        VBUMP="+revision"
        VADHOC=$1
        return 0
      fi
      ;;
  esac

  ## If we get here, could not parse the adhoc version. echo usage and exit.
  echo "usage: version.sh  - report the project.props version string"
  echo "       version.sh <version-string>  - assign an adhoc version-string "
  echo "             of the format 'major.minor.revision-SUFFIX' where "
  echo "             major, minor, and revision are digits, and SUFFIX is"
  echo "             either (SNAPSHOT|MILESTONE)."
  echo "       version.sh +major - increment major number, zero others "
  echo "       version.sh +minor - increment minor number, zero rev"
  echo "       version.sh +revision - increment revision number"
  echo "       version.sh +snapshot - add -SNAPSHOT suffix"
  echo "       version.sh -snapshot - remove -SNAPSHOT suffix"
  echo "       version.sh +milestone - add -MILESTONE suffix"
  echo "       version.sh -milestone - remove -MILESTONE suffix"
  echo ""
  echo "  In all cases except the first, the new version string is assigned to"
  echo "  the project.props file, the changed properties file is added,"
  echo "  committed, and pushed to the repo on the current branch, and the repo"
  echo "  is tagged with the version string."
  exit 1
}

##-----------------------------------------------------------------------------
## $TAG_SEARCH_PATTERN=determineTagSearchPattern($PROP_VERS, $VBUMP)
##
## Given these args:
##   PROP_VERS=the version from the project.props file
##   VBUMP=(+major|+minor|+revision|+build) which version digit to bump
##
## Determine the string pattern to use when calling `git tag -l <pattern>`
## to get a list of tags that would be next in line if the caller is intending
## to bump the $VBUMP digit.
##
## Example 1:
##   PROP_VERS=2.10.12
##   VBUMP=+revision
##   git tags = (2.10.9, 2.10.10, 2.10.11, 2.10.12, 2.10.13)
##
##   The user wants the next 2.10.x version because $VBUMP is +revision. Get
##   list of git tags by calling `git tag -l "2.10.*" --sort=-version:refname`
##
##   Therefore, the returned TAG_SEARCH_PATTERN is the string "2.10.*"
##   The proposed next version change will be 2.10.12 --> 2.10.14.
##
## Example 2:
##   PROP_VERS=2.10.12
##   VBUMP=+minor
##   git tags = (2.11.0, 2.11.1, 2.11.3)
##
##   The user wants the next 2.11.x version because $VBUMP is +minor. Get
##   list of git tags by calling `git tag -l "2.11.*" --sort=-version:refname`
##
##   Therefore, the returned TAG_SEARCH_PATTERN is the string "2.11.*"
##   The proposed next version change will be 2.10.12 --> 2.11.4.
##
## Example 3:
##   PROP_VERS=2.10.12
##   VBUMP=+major
##   git tags = (3.0.0, 3.0.1, 3.1.0, 3.1.1)
##
##   The user wants the next 3.x.x version because $VBUMP is +major. Get
##   list of git tags by calling `git tag -l "3.*" --sort=-version:refname`
##
##   Therefore, the returned TAG_SEARCH_PATTERN is the string "3.*"
##   The proposed next version change will be 2.10.12 --> 3.1.2.
##
## Example 4:
##   PROP_VERS=2.10.12
##   VBUMP=+revision
##   git tags = ()
##   Therefore, the returned TAG_SEARCH_PATTERN is the string "2.10.*"
##   The proposed next version change will be 2.10.12 --> 2.10.13.
##
## Example 5:
##   PROP_VERS=2.10.12
##   VBUMP=+minor
##   git tags = ()
##   Therefore, the returned TAG_SEARCH_PATTERN is the string "2.11.*"
##   The proposed next version change will be 2.10.12 --> 2.11.0.
##
## Example 6:
##   PROP_VERS=2.10.12
##   VBUMP=+major
##   git tags = ()
##   Therefore, the returned TAG_SEARCH_PATTERN is the string "3.*"
##   The proposed next version change will be 2.10.12 --> 3.0.0.
##
##-----------------------------------------------------------------------------
function determineTagSearchPattern {
  local __resultvar=$1
  local _PROP_VERS=$2
  local _VBUMP=$3
  local _MAJOR=''
  local _MINOR=''
  local _REVISION=''
  local _BUILD=''
  local _SUFFIX=''
  local _PATTERN=''
  ## ($MAJOR,$MINOR,$REVISION,$BUILD,$SUFFIX) = versParseVersionString($VERS)
  versParseVersionString ${_PROP_VERS} _MAJOR _MINOR _REVISION _BUILD _SUFFIX

  case ${_VBUMP} in
    +build)
      _PATTERN="${_MAJOR}.${_MINOR}.${_REVISION}.*"
      ;;
    +revision)
      _PATTERN="${_MAJOR}.${_MINOR}.*"
      ;;
    +minor)
      [ "x" = "x${_MINOR}" ] && _MINOR=0
      _MINOR=$(($_MINOR + 1))
      _PATTERN="${_MAJOR}.${_MINOR}.*"
      ;;
    +major)
      [ "x" = "x${_MAJOR}" ] && _MAJOR=0
      _MAJOR=$(($_MAJOR + 1))
      _PATTERN="${_MAJOR}.*"
      ;;
    *)
      _PATTERN="${_MAJOR}.${_MINOR}.*"
      ;;
  esac
  eval $__resultvar="'$_PATTERN'"
}

##-----------------------------------------------------------------------------
## $VBUMP=versDemoteBump($REPO_VERS, $VBUMP)
##
## If REPO_VERS is not empty and VBUMP is either +major or +minor, then we know
## that we found version tags on the repo that are already in the desired
## "bumped" range. We know this because of what determineTagSearchPattern()
## returns based on the same inputs. So we need to lower VBUMP ambitions.
##-----------------------------------------------------------------------------
function demoteBump {
  local __resultvar=$1
  local _REPO_VERS=$2
  local _VBUMP=$3

  local _NEW_VBUMP=${_VBUMP}
  if [ "x" != "x${_REPO_VERS}" ]; then
    case ${_VBUMP} in
      +minor)
        _NEW_VBUMP="+revision"
        ;;
      +major)
        _NEW_VBUMP="+minor"
        ;;
    esac
  fi
  eval $__resultvar="'$_NEW_VBUMP'"
}


##-----------------------------------------------------------------------------
## assignVersion($BUMPED_VERS, $PROP_VERS, $REPO_VERS)
##
## Propose that we upgrade from PROP_VERS/REPO_VERS to BUMPED_VERS. If user
## responds affirmative, then check that our git workspace is clean,
## alter the project.props file version, add the file to git and commit and
## push. Then tag the repo with BUMPED_VERS and push it.
##-----------------------------------------------------------------------------
function assignVersion {
  local _BUMPED_VERS=$1
	local _PROP_VERS=$2
  local _REPO_VERS=$3

  ## propose
  if [ "x" != "x${_REPO_VERS}" ]; then
		echo "project.props: ${_PROP_VERS}, repo: ${_REPO_VERS}  -->  ${_BUMPED_VERS}"
  else
		echo "project.props: ${_PROP_VERS}  -->  ${_BUMPED_VERS}"
  fi

  ## confirm to proceed.
  read -p "Update project.props, commit, tag, and push to repo? (y/N) " -r REPLY;
  if [ "xy" != "x$REPLY" ]; then
      echo "whew! that was close!"
      exit 1;
  fi

  ## got a clean workspace?
  gitIsClean

  ## update the project.props file
  propsSetVersion ${_BUMPED_VERS}
	propsCommitVersion

  ## watch commands as they execute. set +xv negates this
  set -xv

  ## commit the updated project.props file and push it
  git commit -a -m "version ${_BUMPED_VERS}"
  git push --set-upstream origin $(git rev-parse --abbrev-ref HEAD)

  ## create version tag
  git tag -a ${_BUMPED_VERS} -m "version ${_BUMPED_VERS}" `git rev-parse HEAD`
  git push origin ${_BUMPED_VERS}

}


###############################################################################
## THE PROGRAM STARTS HERE
###############################################################################

## Parse the command line to see what the user wants to do
parseCommandLine $1 $2
#echo "VMODE=${VMODE} VBUMP=${VBUMP} VSUFFIX=${VSUFFIX} VADHOC=${VADHOC}"
#exit 1;

## No matter what VMODE is, we need to grab current version information
## call $PROP_VERS=propsGetCurrentVersion() in props.functions
propsGetCurrentVersion PROP_VERS

## call $TAG_SEARCH_PATTERN=determineTagSearchPattern($PROP_VERS, $VBUMP)
determineTagSearchPattern TAG_SEARCH_PATTERN "${PROP_VERS}" "${VBUMP}"

## call $REPO_VERS=gitGetNewestVersionTag($TAG_SEARCH_PATTERN) in git.functions
gitGetNewestVersionTag REPO_VERS "${TAG_SEARCH_PATTERN}"

## Do one of three things: just report the current version(s), bump a version
## number up a notch, or assign an ad hoc version string.
case ${VMODE} in
  REPORT)
    ## If version is in sync between props and repo, just report the raw
    ## version number with no adornment (so that it is parsable by other
    ## programs) and exit 0. But if anything is out of sync, then report who
    ## has what version and exit 1.

    ## case where everything is in sync
    if [[ "${PROP_VERS}" = "${REPO_VERS}" ]]; then
      ## everything is in sync
      echo "${PROP_VERS}"
      exit 0
    fi

    ## okay, something is out of sync. Tell us who has what and exit 1.
    ## first, get a handle on what's going on with our repo version. If it
    ## came back empty, then there is no tag matching the pattern we are
    ## looking for. So just say that.
    REPORT_REPO_VERS="${REPO_VERS}"
    if [ "x" = "x$REPO_VERS" ]; then
      REPORT_REPO_VERS="no tags matching ${TAG_SEARCH_PATTERN}"
    fi
		echo "project.props: ${PROP_VERS}, repo: ${REPORT_REPO_VERS}"
  	exit 1
    ;;
  BUMP)
    ## If we are bumping +major or +minor and the REPO_VERS is already in
    ## the bumped range, then demote our VBUMP ambitions. For example,
    ## the following is wrong:
    ##   $ version.sh +minor
    ##   project.props: 2.10.38, repo: 2.11.2-SNAPSHOT  -->  2.12.0-SNAPSHOT
    ##
    ## The following is correct:
    ##   $ version.sh +minor
    ##   project.props: 2.10.38, repo: 2.11.2-SNAPSHOT  -->  2.11.3-SNAPSHOT
    ##
    ## Because the REPO_VERS is "already there", then demote +minor to
    ## be +revision and operate on the REPO_VERS.
    ##
    ## call $VBUMP=demoteBump($VBUMP, $REPO_VERS)
    ## REPO_VERS might be empty, so quote it
    demoteBump VBUMP "${REPO_VERS}" "${VBUMP}"

    ## are we going to bump from the props vers? or latest repo vers?
    ## sort them, choosing the larger version.
    ## call $FROM_VERS=sortVersions(${PROP_VERS}, ${REPO_VERS}) in version.functions
    sortVersions FROM_VERS "${PROP_VERS}" "${REPO_VERS}"

    ## call $BUMPED_VERS=versBumpVersion($FROM_VERS, $VBUMP, $VSUFFIX) in version.functions
    versBumpVersion BUMPED_VERS "${FROM_VERS}" "${VBUMP}" "${VSUFFIX}"

    ## assign the new version
    assignVersion "${BUMPED_VERS}" "${PROP_VERS}" "${REPO_VERS}"
    exit 0
    ;;
  ADHOC)
    ## assign the new version
    assignVersion "${VADHOC}" "${PROP_VERS}" "${REPO_VERS}"
    exit 0
    ;;
  *)
    echo "Don't know what to do when VMODE=${VMODE}"
    exit 1
    ;;
esac
