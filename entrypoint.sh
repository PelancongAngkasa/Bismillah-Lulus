#!/bin/sh
# filepath: holodeckb2b-7.0.0-A/entrypoint.sh

copy_if_empty() {
  src="$1"
  dest="$2"
  if [ -d "$dest" ] && [ -z "$(ls -A "$dest")" ]; then
    echo "Copying default files from $src to $dest"
    cp -rT "$src" "$dest"
  fi
}

copy_if_empty "/opt/holodeckb2b/_default/repository" "/opt/holodeckb2b/repository"
copy_if_empty "/opt/holodeckb2b/_default/data/msg_in" "/opt/holodeckb2b/data/msg_in"
copy_if_empty "/opt/holodeckb2b/_default/data/msg_out" "/opt/holodeckb2b/data/msg_out"
copy_if_empty "/opt/holodeckb2b/_default/conf" "/opt/holodeckb2b/conf"
copy_if_empty "/opt/holodeckb2b/_default/examples/pmodes" "/opt/holodeckb2b/examples/pmodes"
copy_if_empty "/opt/holodeckb2b/_default/logs" "/opt/holodeckb2b/logs"

# Jalankan perintah utama container (ganti sesuai aplikasi Anda)
exec "$@"