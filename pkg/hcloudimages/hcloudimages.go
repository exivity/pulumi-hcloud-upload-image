package hcloudimages

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/apricote/hcloud-upload-image/hcloudimages"
)

// Static error variables
var (
	ErrHcloudTokenRequired     = errors.New("hcloudToken is required")
	ErrImageURLRequired        = errors.New("imageUrl is required")
	ErrUnsupportedCompression  = errors.New("unsupported compression format")
	ErrUnsupportedImageFormat  = errors.New("unsupported image format")
	ErrUnsupportedArchitecture = errors.New("unsupported architecture")
	ErrServerTypeNotFound      = errors.New("server type not found")
)

// UploadedImage represents a Pulumi resource for uploading custom images to Hetzner Cloud
type UploadedImage struct{}

// UploadedImageArgs defines the input arguments for uploading an image
type UploadedImageArgs struct {
	// HcloudToken is the Hetzner Cloud API token
	HcloudToken string `pulumi:"hcloudToken" provider:"secret"`

	// ImageURL is the URL to download the image from (mutually exclusive with ImageReader)
	ImageURL *string `pulumi:"imageUrl,optional"`

	// ImageCompression describes the compression of the image file
	ImageCompression *string `pulumi:"imageCompression,optional"`

	// ImageFormat describes the format of the image file (raw or qcow2)
	ImageFormat *string `pulumi:"imageFormat,optional"`

	// ImageSize can be optionally set to validate that the image can be written to the server
	ImageSize *int64 `pulumi:"imageSize,optional"`

	// Architecture should match the architecture of the Image (x86 or arm)
	Architecture string `pulumi:"architecture"`

	// ServerType can be optionally set to override the default server type
	ServerType *string `pulumi:"serverType,optional"`

	// Description is an optional description for the resulting image
	Description *string `pulumi:"description,optional"`

	// Labels will be added to the resulting image
	Labels map[string]string `pulumi:"labels,optional"`
}

func (args *UploadedImageArgs) Annotate(a infer.Annotator) {
	a.Describe(&args.HcloudToken, "The Hetzner Cloud API token.")
	a.Describe(&args.ImageURL, "The URL to download the image from. Must be publicly accessible.")
	a.Describe(&args.ImageCompression, "The compression format of the image. Supported: 'none', 'bz2', 'xz'. Defaults to 'none'.")
	a.Describe(&args.ImageFormat, "The format of the image. Supported: 'raw', 'qcow2'. Defaults to 'raw'.")
	a.Describe(&args.ImageSize, "Optional size validation for the image in bytes.")
	a.Describe(&args.Architecture, "The architecture of the image. Supported: 'x86', 'arm'.")
	a.Describe(&args.ServerType, "Optional server type to use for the temporary server. If not specified, a default will be chosen based on architecture.")
	a.Describe(&args.Description, "Optional description for the resulting image.")
	a.Describe(&args.Labels, "Labels to add to the resulting image. These can be used to filter images later.")

	a.SetDefault(&args.ImageCompression, "none")
	a.SetDefault(&args.ImageFormat, "raw")
}

// UploadedImageState represents the state of an uploaded image resource
type UploadedImageState struct {
	UploadedImageArgs

	// ImageID is the ID of the created Hetzner Cloud image
	ImageID int64 `pulumi:"imageId"`

	// ImageName is the name of the created image
	ImageName string `pulumi:"imageName"`

	// Created is the creation timestamp
	Created string `pulumi:"created"`

	// DiskSize is the disk size of the image in GB
	DiskSize int `pulumi:"diskSize"`

	// OSFlavor is the OS flavor of the image
	OSFlavor string `pulumi:"osFlavor"`

	// OSVersion is the OS version of the image
	OSVersion string `pulumi:"osVersion"`

	// Status is the current status of the image
	Status string `pulumi:"status"`

	// Type is the type of the image
	Type string `pulumi:"type"`
}

func (state *UploadedImageState) Annotate(a infer.Annotator) {
	a.Describe(&state.ImageID, "The ID of the created Hetzner Cloud image.")
	a.Describe(&state.ImageName, "The name of the created image.")
	a.Describe(&state.Created, "The creation timestamp of the image.")
	a.Describe(&state.DiskSize, "The disk size of the image in GB.")
	a.Describe(&state.OSFlavor, "The OS flavor of the image.")
	a.Describe(&state.OSVersion, "The OS version of the image.")
	a.Describe(&state.Status, "The current status of the image.")
	a.Describe(&state.Type, "The type of the image.")
}

// Create uploads a new image to Hetzner Cloud
func (UploadedImage) Create( //nolint:cyclop,funlen // TODO: refactor this function
	ctx context.Context, req infer.CreateRequest[UploadedImageArgs],
) (infer.CreateResponse[UploadedImageState], error) {
	name := req.Name
	inputs := req.Inputs

	// Validate required inputs
	if inputs.HcloudToken == "" {
		return infer.CreateResponse[UploadedImageState]{}, ErrHcloudTokenRequired
	}

	if inputs.ImageURL == nil {
		return infer.CreateResponse[UploadedImageState]{}, ErrImageURLRequired
	}

	state := UploadedImageState{UploadedImageArgs: inputs}

	if req.DryRun {
		return infer.CreateResponse[UploadedImageState]{ID: name, Output: state}, nil
	}

	// Create Hetzner Cloud client
	hcloudClient := hcloud.NewClient(hcloud.WithToken(inputs.HcloudToken))
	client := hcloudimages.NewClient(hcloudClient)

	// Parse image URL
	imageURL, err := url.Parse(*inputs.ImageURL)
	if err != nil {
		return infer.CreateResponse[UploadedImageState]{}, fmt.Errorf("invalid image URL: %w", err)
	}

	// Build upload options
	uploadOpts := hcloudimages.UploadOptions{
		ImageURL: imageURL,
	}

	// Set compression
	if inputs.ImageCompression != nil {
		switch *inputs.ImageCompression {
		case "bz2":
			uploadOpts.ImageCompression = hcloudimages.CompressionBZ2
		case "xz":
			uploadOpts.ImageCompression = hcloudimages.CompressionXZ
		case "none", "":
			uploadOpts.ImageCompression = hcloudimages.CompressionNone
		default:
			return infer.CreateResponse[UploadedImageState]{}, fmt.Errorf("%w: %s", ErrUnsupportedCompression, *inputs.ImageCompression)
		}
	}

	// Set image format
	if inputs.ImageFormat != nil {
		switch *inputs.ImageFormat {
		case "qcow2":
			uploadOpts.ImageFormat = hcloudimages.FormatQCOW2
		case "raw", "":
			uploadOpts.ImageFormat = hcloudimages.FormatRaw
		default:
			return infer.CreateResponse[UploadedImageState]{}, fmt.Errorf("%w: %s", ErrUnsupportedImageFormat, *inputs.ImageFormat)
		}
	}

	// Set image size
	if inputs.ImageSize != nil {
		uploadOpts.ImageSize = *inputs.ImageSize
	}

	// Set architecture
	switch inputs.Architecture {
	case "x86":
		uploadOpts.Architecture = hcloud.ArchitectureX86
	case "arm":
		uploadOpts.Architecture = hcloud.ArchitectureARM
	default:
		return infer.CreateResponse[UploadedImageState]{}, fmt.Errorf("%w: %s", ErrUnsupportedArchitecture, inputs.Architecture)
	}

	// Set server type if specified
	if inputs.ServerType != nil {
		serverType, _, err := hcloudClient.ServerType.GetByName(ctx, *inputs.ServerType)
		if err != nil {
			return infer.CreateResponse[UploadedImageState]{}, fmt.Errorf("failed to get server type: %w", err)
		}
		if serverType == nil {
			return infer.CreateResponse[UploadedImageState]{}, fmt.Errorf("%w: %s", ErrServerTypeNotFound, *inputs.ServerType)
		}
		uploadOpts.ServerType = serverType
	}

	// Set description
	if inputs.Description != nil {
		uploadOpts.Description = inputs.Description
	}

	// Set labels
	if inputs.Labels != nil {
		uploadOpts.Labels = inputs.Labels
	}

	// Upload the image
	image, err := client.Upload(ctx, uploadOpts)
	if err != nil {
		return infer.CreateResponse[UploadedImageState]{}, fmt.Errorf("failed to upload image: %w", err)
	}

	// Populate state with image information
	state.ImageID = image.ID
	state.ImageName = image.Name
	state.Created = image.Created.String()
	state.DiskSize = int(image.DiskSize)
	state.OSFlavor = image.OSFlavor
	state.OSVersion = image.OSVersion
	state.Status = string(image.Status)
	state.Type = string(image.Type)

	return infer.CreateResponse[UploadedImageState]{
		ID:     strconv.FormatInt(image.ID, 10),
		Output: state,
	}, nil
}

// Read retrieves the current state of the image
func (UploadedImage) Read(
	ctx context.Context, req infer.ReadRequest[UploadedImageArgs, UploadedImageState],
) (infer.ReadResponse[UploadedImageArgs, UploadedImageState], error) {
	imageID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.ReadResponse[UploadedImageArgs, UploadedImageState]{}, fmt.Errorf("invalid image ID: %w", err)
	}

	// Validate required inputs
	if req.Inputs.HcloudToken == "" {
		return infer.ReadResponse[UploadedImageArgs, UploadedImageState]{}, ErrHcloudTokenRequired
	}

	// Create Hetzner Cloud client
	hcloudClient := hcloud.NewClient(hcloud.WithToken(req.Inputs.HcloudToken))

	// Get the image
	image, _, err := hcloudClient.Image.GetByID(ctx, imageID)
	if err != nil {
		return infer.ReadResponse[UploadedImageArgs, UploadedImageState]{}, fmt.Errorf("failed to get image: %w", err)
	}

	if image == nil {
		// Image doesn't exist anymore
		return infer.ReadResponse[UploadedImageArgs, UploadedImageState]{}, nil
	}

	// Update state with current image information
	state := req.State
	state.ImageID = image.ID
	state.ImageName = image.Name
	state.Created = image.Created.String()
	state.DiskSize = int(image.DiskSize)
	state.OSFlavor = image.OSFlavor
	state.OSVersion = image.OSVersion
	state.Status = string(image.Status)
	state.Type = string(image.Type)

	return infer.ReadResponse[UploadedImageArgs, UploadedImageState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  state,
	}, nil
}

// Delete removes the image from Hetzner Cloud.
func (UploadedImage) Delete(ctx context.Context, req infer.DeleteRequest[UploadedImageState]) (infer.DeleteResponse, error) {
	// Create Hetzner Cloud client
	hcloudClient := hcloud.NewClient(hcloud.WithToken(req.State.HcloudToken))

	imageID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("invalid image ID: %w", err)
	}

	image := &hcloud.Image{ID: imageID}
	_, err = hcloudClient.Image.Delete(ctx, image)
	if err != nil {
		if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
			return infer.DeleteResponse{}, nil // Image already deleted
		}
		return infer.DeleteResponse{}, fmt.Errorf("failed to delete image: %w", err)
	}

	return infer.DeleteResponse{}, nil
}

// Update handles updates to the image resource
func (UploadedImage) Update(
	ctx context.Context, req infer.UpdateRequest[UploadedImageArgs, UploadedImageState],
) (infer.UpdateResponse[UploadedImageState], error) {
	// Most properties can't be updated after creation, so we'll need to replace
	// Only labels and description can potentially be updated

	imageID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.UpdateResponse[UploadedImageState]{}, fmt.Errorf("invalid image ID: %w", err)
	}

	// Validate required inputs
	if req.Inputs.HcloudToken == "" {
		return infer.UpdateResponse[UploadedImageState]{}, ErrHcloudTokenRequired
	}

	// Create Hetzner Cloud client
	hcloudClient := hcloud.NewClient(hcloud.WithToken(req.Inputs.HcloudToken))

	// Update the image with new labels and description
	updateOpts := hcloud.ImageUpdateOpts{}

	if req.Inputs.Description != nil {
		updateOpts.Description = req.Inputs.Description
	}

	if req.Inputs.Labels != nil {
		updateOpts.Labels = req.Inputs.Labels
	}

	image, _, err := hcloudClient.Image.Update(ctx, &hcloud.Image{ID: imageID}, updateOpts)
	if err != nil {
		return infer.UpdateResponse[UploadedImageState]{}, fmt.Errorf("failed to update image: %w", err)
	}

	// Update state
	state := req.State
	state.UploadedImageArgs = req.Inputs
	state.ImageID = image.ID
	state.ImageName = image.Name
	state.Created = image.Created.String()
	state.DiskSize = int(image.DiskSize)
	state.OSFlavor = image.OSFlavor
	state.OSVersion = image.OSVersion
	state.Status = string(image.Status)
	state.Type = string(image.Type)

	return infer.UpdateResponse[UploadedImageState]{
		Output: state,
	}, nil
}

// Diff determines what changes are needed
func (UploadedImage) Diff(
	ctx context.Context, req infer.DiffRequest[UploadedImageArgs, UploadedImageState],
) (infer.DiffResponse, error) {
	diff := map[string]p.PropertyDiff{}

	// Check if properties that require replacement have changed
	// Only imageUrl and architecture changes require replacement
	if stringPtrNotEqual(req.Inputs.ImageURL, req.State.ImageURL) {
		diff["imageUrl"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if req.Inputs.Architecture != req.State.Architecture {
		diff["architecture"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	// Other properties can be updated in place
	if req.Inputs.HcloudToken != req.State.HcloudToken {
		diff["hcloudToken"] = p.PropertyDiff{Kind: p.Update}
	}
	if stringPtrNotEqual(req.Inputs.ImageCompression, req.State.ImageCompression) {
		diff["imageCompression"] = p.PropertyDiff{Kind: p.Update}
	}
	if stringPtrNotEqual(req.Inputs.ImageFormat, req.State.ImageFormat) {
		diff["imageFormat"] = p.PropertyDiff{Kind: p.Update}
	}
	if req.Inputs.ImageSize != req.State.ImageSize {
		diff["imageSize"] = p.PropertyDiff{Kind: p.Update}
	}
	if stringPtrNotEqual(req.Inputs.ServerType, req.State.ServerType) {
		diff["serverType"] = p.PropertyDiff{Kind: p.Update}
	}

	// Labels and description can be updated in place
	if !mapsEqual(req.Inputs.Labels, req.State.Labels) {
		diff["labels"] = p.PropertyDiff{Kind: p.Update}
	}

	if stringPtrNotEqual(req.Inputs.Description, req.State.Description) {
		diff["description"] = p.PropertyDiff{Kind: p.Update}
	}

	return infer.DiffResponse{
		DeleteBeforeReplace: false,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// Annotate provides documentation for the resource
func (r *UploadedImage) Annotate(a infer.Annotator) {
	a.Describe(r, "Uploads a custom disk image to Hetzner Cloud and creates a snapshot that can be used to create servers.")
}

// Helper functions
func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func stringPtrNotEqual(a, b *string) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil || b == nil {
		return true
	}
	return *a != *b
}

// CleanupFunction provides a function to clean up any leftover resources
type CleanupFunction struct{}

type CleanupFunctionArgs struct {
	HcloudToken string `pulumi:"hcloudToken" provider:"secret"`
}

type CleanupFunctionResult struct {
	Message string `pulumi:"message"`
}

func (CleanupFunction) Invoke(
	ctx context.Context, req infer.FunctionRequest[CleanupFunctionArgs],
) (infer.FunctionResponse[CleanupFunctionResult], error) {
	if req.Input.HcloudToken == "" {
		return infer.FunctionResponse[CleanupFunctionResult]{}, ErrHcloudTokenRequired
	}

	// Create Hetzner Cloud client
	hcloudClient := hcloud.NewClient(hcloud.WithToken(req.Input.HcloudToken))
	client := hcloudimages.NewClient(hcloudClient)

	// Clean up temporary resources
	err := client.CleanupTempResources(ctx)
	if err != nil {
		return infer.FunctionResponse[CleanupFunctionResult]{}, fmt.Errorf("failed to cleanup temporary resources: %w", err)
	}

	return infer.FunctionResponse[CleanupFunctionResult]{
		Output: CleanupFunctionResult{
			Message: "Successfully cleaned up temporary resources",
		},
	}, nil
}

func (r *CleanupFunction) Annotate(a infer.Annotator) {
	a.Describe(r, "Cleans up any temporary resources (servers, SSH keys) that may have been left over from failed upload operations.")

	var args CleanupFunctionArgs
	var result CleanupFunctionResult

	a.Describe(&args.HcloudToken, "The Hetzner Cloud API token.")
	a.Describe(&result.Message, "A message indicating the result of the cleanup operation.")
}
