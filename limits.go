package router

// Limits bounds construction, dispatch and URL-generation inputs. Every field
// must be positive; start with DefaultLimits and adjust trusted startup budgets.
// WithLimits validates positivity but imposes no ceiling on caller-selected
// maxima. Limits are retained by value and do not bound callback execution.
// Byte budgets count bytes, not runes; limits reject rather than truncate input.
type Limits struct {
	// MaxRoutes caps retained flattened descriptors, counting each mount once.
	MaxRoutes int
	// MaxGroups caps successful group callbacks, including empty/nested groups.
	MaxGroups int
	// MaxGroupDepth caps nested child-builder depth.
	MaxGroupDepth int
	// MaxMethodsPerRoute caps explicitly registered methods per descriptor.
	MaxMethodsPerRoute int
	// MaxMethodBytes caps each registered or dispatched method token in bytes.
	MaxMethodBytes int
	// MaxWildcardsPerRoute caps combined host/path wildcard occurrences; the
	// path count also includes the {$} end marker.
	MaxWildcardsPerRoute int
	// MaxWildcardNameBytes caps each wildcard/generation parameter name.
	MaxWildcardNameBytes int
	// MaxPatternBytes caps paths, composed group prefixes and mount patterns.
	MaxPatternBytes int
	// MaxHostBytes caps registered route/group host patterns, not trusted bases.
	MaxHostBytes int
	// MaxNameBytes caps route/group names, middleware names and exclusions.
	MaxNameBytes int
	// MaxSourceBytes caps registered source labels before retention.
	MaxSourceBytes int
	// MaxOperationBytes caps retained operation identifiers.
	MaxOperationBytes int
	// MaxDocumentationBytes caps retained descriptive route text.
	MaxDocumentationBytes int
	// MaxMetadataEntries caps composed metadata entry counts before merging.
	MaxMetadataEntries int
	// MaxMetadataKeyBytes caps each non-empty metadata key.
	MaxMetadataKeyBytes int
	// MaxMetadataValueBytes caps each metadata value; empty values are valid.
	MaxMetadataValueBytes int
	// MaxMiddleware caps middleware collections, exclusions and resolved depth.
	MaxMiddleware int
	// MaxRequestTargetBytes caps RequestURI, decoded/raw/escaped path and raw
	// query components plus escaped-path/query bytes; excess dispatch returns 414.
	MaxRequestTargetBytes int
	// MaxURLParameters caps compiled wildcard names, generation descriptors and
	// total supplied segment values, including remainder segments.
	MaxURLParameters int
	// MaxURLParameterBytes caps cumulative raw parameter-name and value bytes.
	MaxURLParameterBytes int
	// MaxQueryValues caps query key count and total entries; an empty value list
	// consumes one entry even though encoding emits no pair.
	MaxQueryValues int
	// MaxQueryBytes caps cumulative raw query key/value bytes before encoding.
	MaxQueryBytes int
	// MaxGeneratedURLBytes caps final escaped relative/absolute output bytes.
	MaxGeneratedURLBytes int
}

// DefaultLimits returns conservative production budgets.
func DefaultLimits() Limits {
	return Limits{
		MaxRoutes:             1_024,
		MaxGroups:             64,
		MaxGroupDepth:         8,
		MaxMethodsPerRoute:    16,
		MaxMethodBytes:        32,
		MaxWildcardsPerRoute:  16,
		MaxWildcardNameBytes:  64,
		MaxPatternBytes:       2_048,
		MaxHostBytes:          255,
		MaxNameBytes:          128,
		MaxSourceBytes:        256,
		MaxOperationBytes:     128,
		MaxDocumentationBytes: 4_096,
		MaxMetadataEntries:    32,
		MaxMetadataKeyBytes:   64,
		MaxMetadataValueBytes: 256,
		MaxMiddleware:         32,
		MaxRequestTargetBytes: 8_192,
		MaxURLParameters:      32,
		MaxURLParameterBytes:  4_096,
		MaxQueryValues:        128,
		MaxQueryBytes:         4_096,
		MaxGeneratedURLBytes:  8_192,
	}
}

func (l Limits) valid() bool {
	return l.MaxRoutes > 0 && l.MaxGroups > 0 && l.MaxGroupDepth > 0 &&
		l.MaxMethodsPerRoute > 0 && l.MaxMethodBytes > 0 &&
		l.MaxWildcardsPerRoute > 0 && l.MaxWildcardNameBytes > 0 &&
		l.MaxPatternBytes > 0 && l.MaxHostBytes > 0 &&
		l.MaxNameBytes > 0 && l.MaxSourceBytes > 0 &&
		l.MaxOperationBytes > 0 && l.MaxDocumentationBytes > 0 &&
		l.MaxMetadataEntries > 0 &&
		l.MaxMetadataKeyBytes > 0 && l.MaxMetadataValueBytes > 0 &&
		l.MaxMiddleware > 0 && l.MaxRequestTargetBytes > 0 &&
		l.MaxURLParameters > 0 &&
		l.MaxURLParameterBytes > 0 && l.MaxQueryValues > 0 &&
		l.MaxQueryBytes > 0 && l.MaxGeneratedURLBytes > 0
}
