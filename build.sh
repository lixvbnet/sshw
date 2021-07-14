OS_LIST=(darwin linux windows)
ARCH_LIST=(amd64)

buildTarget() {
  dist_dir=$1
  name=$2
  os=$3
  arch=$4
  version=$5
  flags=$6

  folder=$name-$os-$arch-$version
  folderPath=$dist_dir/$folder

  echo "---------------------------------------------"
  build_cmd="go build -o $folderPath/$name $flags"
  echo "$build_cmd"
  env GOOS="$os" GOARCH="$arch" sh -c "$build_cmd"
  tar -zcf "$folderPath.tar.gz" -C "$dist_dir" "$folder"
}

#dist_dir=dist
#name=sshw
#version=1.0
#flags="-ldflags=\"-X 'main.CMD=$name' -X 'main.Version=$version'\""

dist_dir=$1
name=$2
version=$3
flags=$4


for os in "${OS_LIST[@]}"; do
  for arch in "${ARCH_LIST[@]}"; do
    buildTarget "$dist_dir" "$name" "$os" "$arch" "$version" "$flags"
  done
done
