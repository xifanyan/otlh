package otlh

import (
	"strconv"
)

type OptionType int

const (
	QUERYPARAM OptionType = iota + 1
	FILE
	BODY
)

type Options interface {
	optionType() OptionType
	options() map[string]string
}

type ListOptions struct {
	pageSize   int
	pageNumber int
	sort       string
	filterTerm string
	filterName string
}

func NewListOptions() *ListOptions {
	return &ListOptions{}
}

func (opts *ListOptions) optionType() OptionType {
	return QUERYPARAM
}

func (opts *ListOptions) WithPageSize(pageSize int) *ListOptions {
	opts.pageSize = pageSize
	return opts
}

func (opts *ListOptions) WithPageNumber(pageNumber int) *ListOptions {
	opts.pageNumber = pageNumber
	return opts
}

func (opts *ListOptions) WithSort(sort string) *ListOptions {
	opts.sort = sort
	return opts
}

func (opts *ListOptions) WithFilterTerm(filterTerm string) *ListOptions {
	opts.filterTerm = filterTerm
	return opts
}

func (opts *ListOptions) WithFilterName(filterName string) *ListOptions {
	opts.filterName = filterName
	return opts
}

func (opts *ListOptions) options() map[string]string {
	params := map[string]string{}

	if opts.pageSize > 0 {
		params["page_size"] = strconv.Itoa(opts.pageSize)
	}

	if opts.pageNumber > 0 {
		params["page_number"] = strconv.Itoa(opts.pageNumber)
	}

	if opts.sort != "" {
		params["sort"] = opts.sort
	}

	if opts.filterTerm != "" {
		params["filter[term]"] = opts.filterTerm
	}

	if opts.filterName != "" {
		params["filter[name]"] = opts.filterName
	}

	return params
}

type FileOptions struct {
	files map[string]string
}

func (opts *FileOptions) optionType() OptionType {
	return FILE
}

func NewFileOptions() *FileOptions {
	return &FileOptions{files: map[string]string{}}
}

func (opts *FileOptions) WithFile(description string, path string) *FileOptions {
	opts.files[description] = path
	return opts
}

func (opts *FileOptions) options() map[string]string {
	return opts.files
}

type BodyOptions struct {
	body map[string]string
}

func NewBodyOptions() *BodyOptions {
	return &BodyOptions{body: map[string]string{}}
}

func (opts *BodyOptions) optionType() OptionType {
	return BODY
}

func (opts *BodyOptions) WithBody(body string) *BodyOptions {
	opts.body["body"] = body
	return opts
}

func (opts *BodyOptions) options() map[string]string {
	return opts.body
}
