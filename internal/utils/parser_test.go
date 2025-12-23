package utils

import (
	"reflect"
	"testing"
)

func TestSliceMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected []Section
	}{
		{
			name: "basic sections with instructions and lists",
			markdown: `## SECTION1
(This is an instruction)
- key: value
- Simple item

## SECTION2
(Another instruction)
- name: John`,
			expected: []Section{
				{
					Header: "SECTION1",
					Content: []interface{}{
						"(This is an instruction)",
						KeyValue{Key: "key", Value: "value"},
						"Simple item",
					},
				},
				{
					Header: "SECTION2",
					Content: []interface{}{
						"(Another instruction)",
						KeyValue{Key: "name", Value: "John"},
					},
				},
			},
		},
		{
			name: "empty sections",
			markdown: `## EMPTY1

## EMPTY2`,
			expected: []Section{
				{Header: "EMPTY1", Content: []interface{}{}},
				{Header: "EMPTY2", Content: []interface{}{}},
			},
		},
		{
			name: "mixed content with filtering",
			markdown: `# TITLE
Some paragraph text

## VALID_SECTION
(Valid instruction)
- valid: item

Regular paragraph should be ignored

## ANOTHER_SECTION
- just a string`,
			expected: []Section{
				{
					Header: "VALID_SECTION",
					Content: []interface{}{
						"(Valid instruction)",
						KeyValue{Key: "valid", Value: "item"},
					},
				},
				{
					Header: "ANOTHER_SECTION",
					Content: []interface{}{
						"just a string",
					},
				},
			},
		},
		{
			name: "key-value with colon in value",
			markdown: `## URL_SECTION
- url: https://example.com:8080/path
- time: 12:30:45`,
			expected: []Section{
				{
					Header: "URL_SECTION",
					Content: []interface{}{
						KeyValue{Key: "url", Value: "https://example.com:8080/path"},
						KeyValue{Key: "time", Value: "12:30:45"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseSections(tt.markdown)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseSections() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestKeyValueParsing(t *testing.T) {
	markdown := `## TEST
- simple string
- key: value
- complex: value with: colons`

	result := ParseSections(markdown)

	if len(result) != 1 {
		t.Fatalf("Expected 1 section, got %d", len(result))
	}

	section := result[0]
	if len(section.Content) != 3 {
		t.Fatalf("Expected 3 content items, got %d", len(section.Content))
	}

	// Check simple string
	if section.Content[0] != "simple string" {
		t.Errorf("Expected 'simple string', got %v", section.Content[0])
	}

	// Check key-value pair
	kv, ok := section.Content[1].(KeyValue)
	if !ok {
		t.Errorf("Expected KeyValue, got %T", section.Content[1])
	}
	if kv.Key != "key" || kv.Value != "value" {
		t.Errorf("Expected key='key', value='value', got key='%s', value='%s'", kv.Key, kv.Value)
	}

	// Check complex key-value with colons in value
	kv2, ok := section.Content[2].(KeyValue)
	if !ok {
		t.Errorf("Expected KeyValue, got %T", section.Content[2])
	}
	if kv2.Key != "complex" || kv2.Value != "value with: colons" {
		t.Errorf("Expected key='complex', value='value with: colons', got key='%s', value='%s'", kv2.Key, kv2.Value)
	}
}
