// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package kinesis

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// ListTags lists kinesis service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(ctx context.Context, conn kinesisiface.KinesisAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &kinesis.ListTagsForStreamInput{
		StreamName: aws.String(identifier),
	}

	output, err := conn.ListTagsForStreamWithContext(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns kinesis service tags.
func Tags(tags tftags.KeyValueTags) []*kinesis.Tag {
	result := make([]*kinesis.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &kinesis.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from kinesis service tags.
func KeyValueTags(ctx context.Context, tags []*kinesis.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// UpdateTags updates kinesis service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(ctx context.Context, conn kinesisiface.KinesisAPI, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		for _, removedTags := range removedTags.Chunks(10) {
			input := &kinesis.RemoveTagsFromStreamInput{
				StreamName: aws.String(identifier),
				TagKeys:    aws.StringSlice(removedTags.IgnoreAWS().Keys()),
			}

			_, err := conn.RemoveTagsFromStreamWithContext(ctx, input)

			if err != nil {
				return fmt.Errorf("untagging resource (%s): %w", identifier, err)
			}
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		for _, updatedTags := range updatedTags.Chunks(10) {
			input := &kinesis.AddTagsToStreamInput{
				StreamName: aws.String(identifier),
				Tags:       aws.StringMap(updatedTags.IgnoreAWS().Map()),
			}

			_, err := conn.AddTagsToStreamWithContext(ctx, input)

			if err != nil {
				return fmt.Errorf("tagging resource (%s): %w", identifier, err)
			}
		}
	}

	return nil
}

func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags interface{}, newTags interface{}) error {
	return UpdateTags(ctx, meta.(*conns.AWSClient).KinesisConn(), identifier, oldTags, newTags)
}
