#!/usr/bin/env bash
# shellcheck disable=SC2034

iso_name="LambdaOS"
iso_publisher="LambdaOS Project <https://lambdaos.dev>"
iso_application="LambdaOS — The TUI-First Linux Distribution"

# Version resolution: env var → exact tag → describe → fallback
if [[ -n "${LAMBDAOS_VERSION:-}" ]]; then
    iso_version="${LAMBDAOS_VERSION}"
elif tag=$(git describe --tags --exact-match 2>/dev/null); then
    iso_version="${tag#v}"
else
    iso_version="$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')"
fi

iso_label="Lambda_OS_${iso_version//./}"
install_dir="arch"
buildmodes=('iso')
bootmodes=('bios.syslinux'
    'uefi.systemd-boot')
pacman_conf="pacman.conf"
airootfs_image_type="squashfs"
airootfs_image_tool_options=('-comp' 'xz' '-Xbcj' 'x86' '-b' '1M' '-Xdict-size' '1M')
bootstrap_tarball_compression=('zstd' '-c' '-T0' '--auto-threads=logical' '--long' '-19')
file_permissions=(
    ["/etc/shadow"]="0:0:400"
    ["/root"]="0:0:750"
    ["/root/.automated_script.sh"]="0:0:755"
    ["/root/.gnupg"]="0:0:700"
    ["/usr/local/bin/choose-mirror"]="0:0:755"
    ["/usr/local/bin/Installation_guide"]="0:0:755"
    ["/usr/local/bin/livecd-sound"]="0:0:755"
    ["/etc/sudoers.d/liveuser"]="0:0:0440"
)
