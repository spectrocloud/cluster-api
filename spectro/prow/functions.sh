print_step() {
	text_val=$1
	set +x
	echo " "
	echo "###################################################
#  ${text_val}
###################################################"
	echo " "
	set -x
}


set_image_tag() {
	IMG_TAG="latest"

	if [[ ${JOB_TYPE} == 'presubmit' ]]; then
	    VERSION_SUFFIX="-dev"
	    IMG_LOC='pr'
	    IMG_TAG=${PULL_NUMBER}
        fi
	if [[ ${JOB_TYPE} == 'periodic' ]]; then
	    VERSION_SUFFIX="-$(date +%m%d%y)"
	    IMG_LOC='daily'
	    IMG_TAG=$(date +%Y%m%d.%H%M)
	fi
	if [[ ${SPECTRO_RELEASE} ]] && [[ ${SPECTRO_RELEASE} == "yes" ]]; then
	    export VERSION_SUFFIX=""
	    IMG_LOC='release'
	    IMG_TAG=$(make get-version)
	fi
  export IMG_LOC
	export PROD_BUILD_ID
	export IMG_TAG
	export VERSION_SUFFIX
}

create_images() {
	print_step "Create and Push the images"
	make release
}

build_vendor_manifest() {
  print_step "Build vendor manifest for cluster-api"
  ../run.sh
}

IMG_LOC="cluster-api/release"
export IMG_LOC
export REGISTRY=${DOCKER_REGISTRY_PUBLIC}/${IMG_LOC}
export STAGING_REGISTRY=${DOCKER_REGISTRY_PUBLIC}/${IMG_LOC}
export PROD_REGISTRY=${DOCKER_REGISTRY_PUBLIC}/${IMG_LOC}