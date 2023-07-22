package main

import (
	"context"
	"fmt"
	"os"

	platformFormat "github.com/containerd/containerd/platforms"

	"dagger.io/dagger"
)

// list of platforms to execute on
var platforms = []dagger.Platform{
	"linux/amd64", // a.k.a. x86_64
	"linux/arm64", // a.k.a. aarch64
}

// util that returns the architecture of the provided platform
func architectureOf(platform dagger.Platform) string {
	return platformFormat.MustParse(string(platform)).Architecture
}

//CONST="-X github.com/juli3nk/dotfiles/version.Version=${VERSION} -X github.com/juli3nk/dotfiles/version.GitCommit=${COMMIT} -X github.com/juli3nk/dotfiles/version.GitState=${GITSTATE} -X github.com/juli3nk/dotfiles/version.BuildDate=$(date +%s)"

func main() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	project := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"ci/"},
	})

	platformVariants := make([]*dagger.Container, 0, len(platforms))
	for _, platform := range platforms {
		//echo " Building ${VERSION} from ${COMMIT} on ${ARCH}"

		// initialize this container with the platform
		ctr := client.Container()
		ctr = ctr.From("golang:1-alpine")

		ctr = ctr.WithDirectory("/src", project)

		ctr = ctr.WithDirectory("/output", client.Directory())

		// ensure the binary will be statically linked and thus executable
		// in the final image
		ctr = ctr.WithEnvVariable("CGO_ENABLED", "0")

		// configure the go compiler to use cross-compilation targeting the
		// desired platform
		ctr = ctr.WithEnvVariable("GOOS", "linux")
		ctr = ctr.WithEnvVariable("GOARCH", architectureOf(platform))

		ctr = ctr.WithWorkdir("/src")

		ctr = ctr.WithExec([]string{"apk", "--update", "add",
			"ca-certificates",
			"gcc",
			"git",
			"musl-dev",
		})

		ctr = ctr.WithExec([]string{
			"go",
			"build",
			"-ldflags",
			"-linkmode external -extldflags -static -s -w",
			"-o", fmt.Sprintf("/output/dotfiles-%s", architectureOf(platform)),
		})

		// select the output directory
		outputDir := ctr.Directory("/output")
	}
}
