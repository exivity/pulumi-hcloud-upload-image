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
		WithDescription("A Pulumi provider for uploading custom images to Hetzner Cloud using https://github.com/apricote/hcloud-upload-image").
		WithGoImportPath("github.com/linuxluigi/pulumi-hcloud-upload-image/sdk/go/pulumi-hcloud-upload-image").
		WithLicense("MIT License").
		WithResources(
			infer.Resource(hcloudimages.UploadedImage{}),
		).
		Build()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = p.Run(context.Background(), "hcloud-upload-image", "0.1.0")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Provider initialized successfully")
	os.Exit(0)
}
