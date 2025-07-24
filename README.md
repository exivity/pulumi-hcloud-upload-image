# Pulumi Hetzner Cloud Upload Image Provider

A Pulumi provider for uploading custom images to Hetzner Cloud using [hcloud-upload-image](https://github.com/apricote/hcloud-upload-image).

This provider allows you to upload disk images to Hetzner Cloud and create snapshots that can be used to create servers.

## Development

### Building the Provider

To build the provider locally:

```bash
make build
```

### Installing the Plugin Locally

For local testing, you can install the plugin directly:

```bash
make install-plugin
```

### Generating SDKs

To generate SDKs for all supported languages (automatically runs in CI on PRs):

```bash
make gen-sdk
```

You can also generate SDKs for specific languages:

- TypeScript: `make gen-sdk-typescript`
- Python: `make gen-sdk-python`
- Go: `make gen-sdk-go`
- C#: `make gen-sdk-csharp`
- Java: `make gen-sdk-java`

### Versioning

To create a new version, update the version in the `version` file. The current version format is semantic versioning (e.g., `v0.0.1`).

## Usage Examples

### Go

```go
package main

import (
    "errors"
    "fmt"
    "time"

    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/exivity/pulumi-hcloud-upload-image/sdk/go/pulumi-hcloud-upload-image/hcloudimages"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        snapshot, err := hcloudimages.NewUploadedImage(ctx, "my-image", &hcloudimages.UploadedImageArgs{
            Description:      pulumi.Sprintf("Custom image - %s", time.Now().Format(time.RFC3339)),
            HcloudToken:      pulumi.String("your-hetzner-token"),
            Architecture:     pulumi.String("x86"),
            ImageUrl:         pulumi.String("https://example.com/my-image.raw.xz"),
            ImageCompression: pulumi.StringPtr("xz"),
            ServerType:       pulumi.String("cx11"),
            Labels: pulumi.StringMap{
                "environment": pulumi.String("production"),
                "stack":       pulumi.String(ctx.Stack()),
                "project":     pulumi.String(ctx.Project()),
            },
        }, pulumi.IgnoreChanges([]string{"description"}))
        
        if err != nil {
            return err
        }

        ctx.Export("imageId", snapshot.ImageId)
        ctx.Export("imageName", snapshot.ImageName)
        return nil
    })
}
```

### TypeScript

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as hcloud from "@exivity/pulumi-hcloud-upload-image";

const snapshot = new hcloud.hcloudimages.UploadedImage("my-image", {
    description: `Custom image - ${new Date().toISOString()}`,
    hcloudToken: "your-hetzner-token",
    architecture: "x86",
    imageUrl: "https://example.com/my-image.raw.xz",
    imageCompression: "xz",
    serverType: "cx11",
    labels: {
        environment: "production",
        stack: pulumi.getStack(),
        project: pulumi.getProject(),
    },
}, {
    ignoreChanges: ["description"],
});

export const imageId = snapshot.imageId;
export const imageName = snapshot.imageName;
```

### Python

```python
import pulumi
import pulumi_hcloud_upload_image as hcloud

snapshot = hcloud.hcloudimages.UploadedImage("my-image",
    description=f"Custom image - {pulumi.get_project()}-{pulumi.get_stack()}",
    hcloud_token="your-hetzner-token",
    architecture="x86",
    image_url="https://example.com/my-image.raw.xz",
    image_compression="xz",
    server_type="cx11",
    labels={
        "environment": "production",
        "stack": pulumi.get_stack(),
        "project": pulumi.get_project(),
    },
    opts=pulumi.ResourceOptions(ignore_changes=["description"])
)

pulumi.export("image_id", snapshot.image_id)
pulumi.export("image_name", snapshot.image_name)
```

### C #

```csharp
using System.Collections.Generic;
using Pulumi;
using HcloudUploadImage = HcloudUploadImage.HcloudUploadImage;

return await Deployment.RunAsync(() =>
{
    var snapshot = new HcloudUploadImage.Hcloudimages.UploadedImage("my-image", new()
    {
        Description = $"Custom image - {Deployment.Instance.ProjectName}-{Deployment.Instance.StackName}",
        HcloudToken = "your-hetzner-token",
        Architecture = "x86",
        ImageUrl = "https://example.com/my-image.raw.xz",
        ImageCompression = "xz",
        ServerType = "cx11",
        Labels = new Dictionary<string, string>
        {
            ["environment"] = "production",
            ["stack"] = Deployment.Instance.StackName,
            ["project"] = Deployment.Instance.ProjectName,
        },
    }, new CustomResourceOptions
    {
        IgnoreChanges = { "description" },
    });

    return new Dictionary<string, object?>
    {
        ["imageId"] = snapshot.ImageId,
        ["imageName"] = snapshot.ImageName,
    };
});
```

## Resource Properties

### UploadedImage

#### Required Arguments

- `architecture` (string): The architecture of the image. Supported values: 'x86', 'arm'
- `hcloudToken` (string): The Hetzner Cloud API token

#### Optional Arguments

- `description` (string): Optional description for the resulting image
- `imageCompression` (string): The compression format of the image. Supported values: 'none', 'bz2', 'xz'. Defaults to 'none'
- `imageFormat` (string): The format of the image. Supported values: 'raw', 'qcow2'. Defaults to 'raw'
- `imageSize` (number): Optional size validation for the image in bytes
- `imageUrl` (string): The URL to download the image from. Must be publicly accessible
- `labels` (map): Labels to add to the resulting image. These can be used to filter images later
- `serverType` (string): Optional server type to use for the temporary server. If not specified, a default will be chosen based on architecture

#### Outputs

- `created` (string): The creation timestamp of the image
- `diskSize` (number): The disk size of the image in GB
- `imageId` (number): The ID of the created Hetzner Cloud image
- `imageName` (string): The name of the created image
- `osFlavor` (string): The OS flavor of the image
- `osVersion` (string): The OS version of the image
- `status` (string): The current status of the image
- `type` (string): The type of the image

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and ensure SDK generation works: `make gen-sdk`
5. Submit a pull request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
