# Pulumi Hetzner Cloud Upload Image Provider

A Pulumi provider for uploading custom images to Hetzner Cloud.

## Building the Provider

```bash
make build
```

This will create the provider binary at `bin/pulumi-resource-hcloud-upload-image`.

## Using the Provider

### 1. Install the Provider

Copy the binary to your PATH or to a location where Pulumi can find it:

```bash
# Option 1: Copy to a directory in your PATH
cp bin/pulumi-resource-hcloud-upload-image /usr/local/bin/

# Option 2: Add the bin directory to your PATH
export PATH=$PATH:/path/to/pulumi-hcloud-upload-image/bin
```

### 2. Generate SDKs (Optional)

To use the provider in languages other than YAML, generate SDKs:

```bash
# Generate all SDKs
pulumi package gen-sdk ./bin/pulumi-resource-hcloud-upload-image

# Generate specific language SDK
pulumi package gen-sdk ./bin/pulumi-resource-hcloud-upload-image --language typescript
pulumi package gen-sdk ./bin/pulumi-resource-hcloud-upload-image --language python
pulumi package gen-sdk ./bin/pulumi-resource-hcloud-upload-image --language go
```

### 3. Using in a Pulumi Project

#### Pulumi YAML Example

Create a `Pulumi.yaml` file:

```yaml
name: my-hcloud-image-project
runtime: yaml
description: Upload custom image to Hetzner Cloud

resources:
  my-custom-image:
    type: hcloud-upload-image:index:UploadedImage
    properties:
      imageUrl: "https://example.com/my-custom-image.qcow2"
      imageFormat: "qcow2"
      imageCompression: "gzip"
      architecture: "x86"
      description: "My custom image"
      labels:
        environment: "production"
        project: "my-project"

outputs:
  imageId: ${my-custom-image.imageId}
  imageName: ${my-custom-image.imageName}
```

#### TypeScript Example

```typescript
import * as hcloud from "@pulumi/hcloud-upload-image";

const customImage = new hcloud.UploadedImage("my-custom-image", {
    imageUrl: "https://example.com/my-custom-image.qcow2",
    imageFormat: "qcow2",
    imageCompression: "gzip",
    architecture: "x86",
    description: "My custom image",
    labels: {
        environment: "production",
        project: "my-project",
    },
});

export const imageId = customImage.imageId;
export const imageName = customImage.imageName;
```

#### Python Example

```python
import pulumi_hcloud_upload_image as hcloud

custom_image = hcloud.UploadedImage("my-custom-image",
    image_url="https://example.com/my-custom-image.qcow2",
    image_format="qcow2",
    image_compression="gzip",
    architecture="x86",
    description="My custom image",
    labels={
        "environment": "production",
        "project": "my-project",
    }
)

pulumi.export("image_id", custom_image.image_id)
pulumi.export("image_name", custom_image.image_name)
```

### 4. Configure the Provider

Set your Hetzner Cloud API token:

```bash
# Set as environment variable
export HCLOUD_TOKEN="your-hetzner-cloud-api-token"

# Or configure in Pulumi config
pulumi config set hcloud-upload-image:hcloudToken "your-hetzner-cloud-api-token" --secret
```

### 5. Deploy

```bash
pulumi up
```

## Resource Properties

### UploadedImage

#### Inputs

- `imageUrl` (string, optional): URL to download the image from
- `imageCompression` (string, optional): Compression of the image file
- `imageFormat` (string, optional): Format of the image file (raw or qcow2)
- `imageSize` (int64, optional): Size validation for the image
- `architecture` (string, required): Architecture of the image (x86 or arm)
- `serverType` (string, optional): Override default server type
- `description` (string, optional): Description for the resulting image
- `labels` (map[string]string, optional): Labels to add to the image

#### Outputs

- `imageId` (int64): ID of the created Hetzner Cloud image
- `imageName` (string): Name of the created image
- `created` (string): Creation timestamp
- `diskSize` (int): Disk size of the image in GB
- `osFlavor` (string): OS flavor of the image
- `osVersion` (string): OS version of the image
- `status` (string): Current status of the image
- `type` (string): Type of the image

## Development

### Prerequisites

- Go 1.21+
- Pulumi CLI

### Commands

```bash
make download    # Download dependencies
make build      # Build the provider
make test       # Run tests
make lint       # Run linting
make clean      # Clean build artifacts
```
