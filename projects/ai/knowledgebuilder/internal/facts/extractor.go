package facts


// Step 6: Fact Extraction Layer

// Every analyzer emits:

// type Fact struct {
// }

// Examples:

// {
//   "type":"service",
//   "value":"Payment Service"
// }
// {
//   "type":"error",
//   "value":"Kafka timeout"
// }
// {
//   "type":"dependency",
//   "value":"Payment Service -> Kafka"
// }

// This becomes your common language.


// Step 7: Correlation Engine

// Merge duplicate facts.

// Example:

// Image:

// Payment Service

// Log:

// [payment-service]

// Resolve:

// {
//   "canonical":"Payment Service"
// }

// Use simple fuzzy matching first.

// Libraries:

// sahilm/fuzzy
// agnivade/levenshtein