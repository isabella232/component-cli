// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0
package process

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

const processorTimeout = 60 * time.Second

type resourceProcessingPipelineImpl struct {
	processors []ResourceStreamProcessor
}

func (p *resourceProcessingPipelineImpl) Process(ctx context.Context, cd cdv2.ComponentDescriptor, res cdv2.Resource) (*cdv2.ComponentDescriptor, cdv2.Resource, error) {
	infile, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, cdv2.Resource{}, fmt.Errorf("unable to create temporary infile: %w", err)
	}

	if err := WriteProcessorMessage(cd, res, nil, infile); err != nil {
		return nil, cdv2.Resource{}, fmt.Errorf("unable to write: %w", err)
	}

	for _, proc := range p.processors {
		outfile, err := p.process(ctx, infile, proc)
		if err != nil {
			return nil, cdv2.Resource{}, err
		}

		infile = outfile
	}
	defer infile.Close()

	if _, err := infile.Seek(0, io.SeekStart); err != nil {
		return nil, cdv2.Resource{}, fmt.Errorf("unable to seek to beginning of file: %w", err)
	}

	processedCD, processedRes, blobreader, err := ReadProcessorMessage(infile)
	if err != nil {
		return nil, cdv2.Resource{}, fmt.Errorf("unable to read output data: %w", err)
	}
	if blobreader != nil {
		defer blobreader.Close()
	}

	return processedCD, processedRes, nil
}

func (p *resourceProcessingPipelineImpl) process(ctx context.Context, infile *os.File, proc ResourceStreamProcessor) (*os.File, error) {
	defer infile.Close()

	if _, err := infile.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("unable to seek to beginning of input file: %w", err)
	}

	outfile, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary outfile: %w", err)
	}

	inreader := infile
	outwriter := outfile

	ctx, cancelfunc := context.WithTimeout(ctx, processorTimeout)
	defer cancelfunc()

	if err := proc.Process(ctx, inreader, outwriter); err != nil {
		return nil, fmt.Errorf("unable to process resource: %w", err)
	}

	return outfile, nil
}

// NewResourceProcessingPipeline returns a new ResourceProcessingPipeline
func NewResourceProcessingPipeline(processors ...ResourceStreamProcessor) ResourceProcessingPipeline {
	p := resourceProcessingPipelineImpl{
		processors: processors,
	}
	return &p
}
