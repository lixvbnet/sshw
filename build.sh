OS_LIST=(darwin linux windows)
ARCH_LIST=(amd64)

buildTarget() {
  dist_dir=$1
  name=$2
  os=$3
  arch=$4
  flags=$5

  # grep line 'const Version = ""'
  version=$(grep -E "const\b\s*\bVersion\s*=" "$main_dir/main.go" | grep -E -o "(\d+)((\.{1}\d+)*)(\.{0})")
  folder=$name-$version-$os-$arch
  folderPath=$dist_dir/$folder

  echo "---------------------------------------------"
  filename=$name
  if [ "$os" = "windows" ]; then
    filename="$name.exe"
  fi
  build_cmd="go build -o $folderPath/$filename $flags $main_dir"
  echo "$build_cmd"
  env GOOS="$os" GOARCH="$arch" sh -c "$build_cmd"
  tar -zcf "$folderPath.tar.gz" -C "$dist_dir" "$folder"
}

#dist_dir=dist
#name=$(basename $PWD)
#flags="-ldflags=\"-X 'main.CMD=$name'\""
#main_dir="./cmd/$name"
#main_dir="."

dist_dir=$1
name=$2
flags=$3
main_dir=$4


for os in "${OS_LIST[@]}"; do
  for arch in "${ARCH_LIST[@]}"; do
    buildTarget "$dist_dir" "$name" "$os" "$arch" "$flags" "$main_dir"
  done
done
