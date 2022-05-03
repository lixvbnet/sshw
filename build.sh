OS_LIST=(darwin linux windows)
ARCH_LIST=(amd64)

buildTarget() {
  dist_dir=$1
  name=$2
  os=$3
  arch=$4
  flags=$5

  version=$(grep "\bVersion\b" version.go | grep -E -o "(\d+)((\.{1}\d+)*)(\.{0})")
  folder=$name-$version-$os-$arch
  folderPath=$dist_dir/$folder

  echo "---------------------------------------------"
  build_cmd="go build -o $folderPath/$name $flags"
  echo "$build_cmd"
  env GOOS="$os" GOARCH="$arch" sh -c "$build_cmd"
  tar -zcf "$folderPath.tar.gz" -C "$dist_dir" "$folder"
}

#dist_dir=dist
#name=$(basename $PWD)
#flags="-ldflags=\"-X 'main.CMD=$name'\""

dist_dir=$1
name=$2
flags=$3


for os in "${OS_LIST[@]}"; do
  for arch in "${ARCH_LIST[@]}"; do
    buildTarget "$dist_dir" "$name" "$os" "$arch" "$flags"
  done
done
