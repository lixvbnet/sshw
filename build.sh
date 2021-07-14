OS_LIST=(darwin linux windows)
ARCH_LIST=(amd64)

buildTarget() {
  dist_dir=$1
  cmd=$2
  os=$3
  arch=$4
  version=$5
  flags=$6

  folder=$cmd-$os-$arch-$version
  folderPath=$dist_dir/$folder

  echo "---------------------------------------------"
  build_cmd="go build -o $folderPath/$cmd $flags"
  echo "$build_cmd"
  env GOOS="$os" GOARCH="$arch" sh -c "$build_cmd"
  tar -zcf "$folderPath.tar.gz" -C "$dist_dir" "$folder"
}

#dist_dir=dist
#cmd=sshw
#version=1.0
#flags="-ldflags=\"-X 'main.CMD=$cmd' -X 'main.Version=$version'\""

dist_dir=$1
cmd=$2
version=$3
flags=$4


for os in "${OS_LIST[@]}"; do
  for arch in "${ARCH_LIST[@]}"; do
    buildTarget "$dist_dir" "$cmd" "$os" "$arch" "$version" "$flags"
  done
done
