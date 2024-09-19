detect_os() {
    unameOut="$(uname -s)"
    case "${unameOut}" in
        Linux*)     OS=linux;;
        Darwin*)    OS=darwin;;
        CYGWIN*)    OS=windows;;
        MINGW*)     OS=windows;;
        *)          OS="UNKNOWN:${unameOut}"
    esac
    echo "${OS}"
}

detect_arch() {
    archOut="$(uname -m)"
    case "${archOut}" in
        x86_64)    ARCH=amd64;;
        i686)      ARCH=386;;
        i386)      ARCH=386;;
        armv7l)    ARCH=arm;;
        aarch64)   ARCH=arm64;;
        arm64)     ARCH=arm64;;
        *)         ARCH="UNKNOWN:${archOut}"
    esac
    echo "${ARCH}"
}