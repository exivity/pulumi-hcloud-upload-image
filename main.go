package main

import (
	"context"
	"fmt"
	"os"

	"github.com/linuxluigi/pulumi-hcloud-upload-image/pkg/hcloudimages"
	"github.com/pulumi/pulumi-go-provider/infer"
)

func main() {
	p, err := infer.NewProviderBuilder().
		WithNamespace("hcloud-upload-image").
		WithResources(
			infer.Resource(hcloudimages.UploadedImage{}),
		).
		WithConfig(
			infer.Config(&hcloudimages.Config{}),
		).
		Build()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	p.Run(context.Background(), "hcloud-upload-image", "0.1.0")
}
