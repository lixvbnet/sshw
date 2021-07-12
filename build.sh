OS_LIST=(darwin linux windows)
ARCH_LIST=(amd64)

buildTarget() {
  dist_dir=$1
  cmd=$2
  os=$3
  arch=$4
  version=$5

  folder=$cmd-$os-$arch-$version
  folderPath=$dist_dir/$folder

  echo "building target $folderPath"
  env GOOS=$os GOARCH=$arch sh -c "go build -o $folderPath/$cmd"
  tar -zcf "$folderPath.tar.gz" -C $dist_dir "$folder"
}

#dist_dir=dist
#cmd=sshw
#version=1.0

dist_dir=$1
cmd=$2
version=$3

for os in "${OS_LIST[@]}"; do
  for arch in "${ARCH_LIST[@]}"; do
    buildTarget $dist_dir $cmd $os $arch $version $version
  done
done
